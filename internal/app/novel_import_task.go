package app

import (
	"context"
	"fmt"
	"strings"

	"denova/internal/agent"
	"denova/internal/book"
)

// NovelImportTaskRequest carries everything the background import task needs.
// The handler owns upload parsing and split-regex inference options.
type NovelImportTaskRequest struct {
	Filename    string
	Data        []byte
	Title       string
	Author      string
	Description string
	Options     book.NovelImportOptions
}

// StartNovelImportTask runs a confirmed novel import in the recoverable task
// pipeline and keeps the latest result visible in the task center.
func (a *App) StartNovelImportTask(ctx context.Context, req NovelImportTaskRequest) (*Task, error) {
	layered, err := a.Settings()
	if err != nil {
		return nil, err
	}
	if layered.Paths.DenovaDir == "" {
		return nil, fmt.Errorf("缺少数据目录")
	}
	novaDir := layered.Paths.DenovaDir
	title := strings.TrimSpace(req.Title)
	task := NewTask(func(ctx context.Context, task *Task, emit func(agent.Event)) {
		defer a.clearActiveNovelImportTask(task)
		emit(agent.Event{Type: "progress", Data: novelImportTaskProgress{Step: "uploaded"}})
		preview, err := book.PreviewNovelImport(req.Filename, req.Data, req.Options)
		if err != nil {
			a.RecordNovelImportResult("", req.Filename, err)
			emit(agent.Event{Type: "error", Data: novelImportTaskError{Error: err.Error()}})
			return
		}
		if title == "" {
			title = preview.Title
		}
		emit(agent.Event{Type: "progress", Data: novelImportTaskProgress{Step: "create_book"}})
		workspace, meta, err := a.CreateBook(ctx, novaDir, title, req.Author, req.Description)
		if err != nil {
			a.RecordNovelImportResult("", title, err)
			emit(agent.Event{Type: "error", Data: novelImportTaskError{Error: err.Error()}})
			return
		}
		emit(agent.Event{Type: "progress", Data: novelImportTaskProgress{Step: "importing"}})
		importPreview, paths, err := book.ImportNovelToWorkspace(workspace, req.Filename, req.Data, req.Options)
		if err != nil {
			a.RecordNovelImportResult(workspace, title, err)
			emit(agent.Event{Type: "error", Data: novelImportTaskError{Error: err.Error()}})
			return
		}
		a.RecordNovelImportResult(workspace, title, nil)
		emit(agent.Event{Type: "done", Data: book.NovelImportResult{
			Workspace:    workspace,
			BookMeta:     &meta,
			Title:        importPreview.Title,
			ChapterCount: importPreview.ChapterCount,
			TotalChars:   importPreview.TotalChars,
			ChapterPaths: paths,
		}})
	})
	a.mu.Lock()
	a.activeNovelImportTask = task
	a.activeNovelImportTitle = req.Filename
	a.mu.Unlock()
	return task, nil
}

func (a *App) clearActiveNovelImportTask(task *Task) {
	a.mu.Lock()
	if a.activeNovelImportTask == task {
		a.activeNovelImportTask = nil
		a.activeNovelImportTitle = ""
	}
	a.mu.Unlock()
}

type novelImportTaskProgress struct {
	Step string `json:"step"`
}

type novelImportTaskError struct {
	Error string `json:"error"`
}
