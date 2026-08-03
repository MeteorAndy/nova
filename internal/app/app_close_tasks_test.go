package app

import (
	"context"
	"testing"

	"denova/config"
	"denova/internal/agent"
)

func TestAppCloseStopsAgentAndInteractiveRuns(t *testing.T) {
	agentStarted := make(chan struct{})
	interactiveStarted := make(chan struct{})
	agentTask := blockingLifecycleTask(agentStarted)
	interactiveTask := blockingLifecycleTask(interactiveStarted)
	t.Cleanup(func() {
		agentTask.Abort()
		interactiveTask.Abort()
	})

	application := &App{
		cfg: &config.Config{},
		agentTaskRuns: map[string]*agentTaskRun{
			"agent": {task: agentTask},
		},
		interactiveTaskRuns: map[string]*interactiveTaskRun{
			"story": {task: interactiveTask},
		},
	}
	<-agentStarted
	<-interactiveStarted

	application.Close()

	if agentTask.Status() != TaskAborted || !agentTask.Finished() {
		t.Fatalf("Agent task remained active after App.Close: status=%s finished=%t", agentTask.Status(), agentTask.Finished())
	}
	if interactiveTask.Status() != TaskAborted || !interactiveTask.Finished() {
		t.Fatalf("interactive task remained active after App.Close: status=%s finished=%t", interactiveTask.Status(), interactiveTask.Finished())
	}
}

func blockingLifecycleTask(started chan<- struct{}) *Task {
	return NewTask(func(ctx context.Context, _ *Task, _ func(agent.Event)) {
		close(started)
		<-ctx.Done()
	})
}
