package workerpool

import (
	"context"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/tasks"
	"sync"
)

type WorkerPool struct {
	tasks    chan tasks.Task
	results  chan tasks.TaskResult
	workersN int
	wg       sync.WaitGroup
}

func NewWorkerPool(workersN, tasksN int) *WorkerPool {
	return &WorkerPool{
		workersN: workersN,
		tasks:    make(chan tasks.Task, tasksN),
		results:  make(chan tasks.TaskResult, tasksN),
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	wp.wg.Add(wp.workersN)

	for i := 0; i < wp.workersN; i++ {
		go wp.worker()
	}

	go func() {
		<-ctx.Done()
		close(wp.tasks)
	}()

	go func() {
		wp.wg.Wait()
		close(wp.results)
	}()
}

func (wp *WorkerPool) AddTask(taskID int, task func() (string, error)) {
	wp.tasks <- tasks.Task{ID: taskID, Execute: task}
}

func (wp *WorkerPool) GetTaskResult() tasks.TaskResult {
	return <-wp.results
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for task := range wp.tasks {
		res, err := task.Execute()
		wp.results <- tasks.TaskResult{TaskID: task.ID, Result: res, Error: err}
	}
}
