package workerpool

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/tasks"
	"sync"
)

type WorkerPool struct {
	tasks    chan tasks.Task
	results  chan tasks.TaskResult
	workersN int
	wg       sync.WaitGroup
	mu       sync.Mutex
	done     chan struct{}
}

func NewWorkerPool(workersN, tasksN int) *WorkerPool {
	return &WorkerPool{
		workersN: workersN,
		tasks:    make(chan tasks.Task, tasksN),
		results:  make(chan tasks.TaskResult, tasksN),
		done:     make(chan struct{}),
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	wp.wg.Add(wp.workersN)

	for i := 1; i <= wp.workersN; i++ {
		go wp.worker(i)
	}

	go wp.resultLogger()

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

func (wp *WorkerPool) SetNumWorkers(workersN int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.workersN > workersN {
		wp.workersN = workersN
	}

	for wp.workersN < workersN {
		wp.workersN++
		wp.wg.Add(1)

		go wp.worker(wp.workersN)
	}
}

func (wp *WorkerPool) Done() <-chan struct{} {
	return wp.done
}

func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		fmt.Printf("%d) started\n", task.ID)

		res, err := task.Execute()
		wp.results <- tasks.TaskResult{TaskID: task.ID, Result: res, Error: err}

		if workerID > wp.workersN {
			return
		}
	}
}

func (wp *WorkerPool) resultLogger() {
	for task := range wp.results {
		if task.Error != nil {
			fmt.Printf("%d) error: %v\n", task.TaskID, task.Error)
		} else {
			fmt.Printf("%d) ok: %s\n", task.TaskID, task.Result)
		}
	}

	wp.done <- struct{}{}
}
