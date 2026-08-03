package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/ut"

	"denova/internal/book"
	"denova/internal/taskcenter"
)

func TestNovelImportAppearsInTaskCenterAsImportExport(t *testing.T) {
	application := newTestApplication(t)
	server := NewServer(application, "0")

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("book_title", "测试书"); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("split_strategy", "builtin"); err != nil {
		t.Fatal(err)
	}
	part, err := writer.CreateFormFile("file", "novel.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("第一章 测试\n内容一行\n第二章 测试\n内容二行\n")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	importResp := ut.PerformRequest(
		server.engine.Engine,
		http.MethodPost,
		"/api/books/import-novel",
		&ut.Body{Body: bytes.NewReader(body.Bytes()), Len: body.Len()},
		ut.Header{Key: "Content-Type", Value: writer.FormDataContentType()},
	)
	if importResp.Code != http.StatusOK {
		t.Fatalf("import status = %d body=%s", importResp.Code, importResp.Body.String())
	}
	var result book.NovelImportResult
	decodeResponse(t, importResp.Body.Bytes(), &result)
	if result.Workspace == "" || result.ChapterCount != 2 {
		t.Fatalf("unexpected import result: %#v", result)
	}

	tasksResp := performJSONRequest(t, server, http.MethodGet, "/api/tasks", nil)
	if tasksResp.Code != http.StatusOK {
		t.Fatalf("task center status = %d body=%s", tasksResp.Code, tasksResp.Body.String())
	}
	var tasks taskcenter.ListResult
	decodeResponse(t, tasksResp.Body.Bytes(), &tasks)
	var importTask *taskcenter.Task
	for i := range tasks.Tasks {
		if tasks.Tasks[i].Type == taskcenter.TaskTypeImportExport {
			importTask = &tasks.Tasks[i]
			break
		}
	}
	if importTask == nil {
		t.Fatalf("task center missing import/export task: %#v", tasks.Tasks)
	}
	if importTask.Status != taskcenter.StatusCompleted {
		t.Fatalf("import task status = %q, want %q", importTask.Status, taskcenter.StatusCompleted)
	}
	if importTask.Project.Path != result.Workspace {
		t.Fatalf("import task project path = %q, want %q", importTask.Project.Path, result.Workspace)
	}
	if importTask.Recovery.Kind != taskcenter.RecoveryKind("import_export") ||
		importTask.Recovery.Workspace != result.Workspace {
		t.Fatalf("import task recovery = %#v", importTask.Recovery)
	}
	if importTask.Error != "" {
		t.Fatalf("completed import task should not carry an error: %#v", importTask)
	}
}
