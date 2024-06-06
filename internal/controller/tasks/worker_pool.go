package tasks

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workersN int
	tasks    chan Task
	results  chan<- TaskResult
	wg       sync.WaitGroup
}

func NewWorkerPool(workersN int, tasks chan Task, results chan<- TaskResult) *WorkerPool {
	return &WorkerPool{workersN: workersN, tasks: tasks, results: results}
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

func (wp *WorkerPool) AddTask(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for task := range wp.tasks {
		res, err := task.Execute()
		wp.results <- TaskResult{task.ID, res, err}
	}
}
