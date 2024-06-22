package workerpool

import (
	"fmt"
	"sync"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/tasks"
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
		tasks:    make(chan tasks.Task, tasksN),
		results:  make(chan tasks.TaskResult, tasksN),
		workersN: workersN,
		done:     make(chan struct{}),
	}
}

func (wp *WorkerPool) Run() {
	wp.wg.Add(wp.workersN)

	for i := 1; i <= wp.workersN; i++ {
		go wp.worker(i)
	}

	go wp.resultLogger()

	go func() {
		wp.wg.Wait()
		close(wp.results)
	}()
}

func (wp *WorkerPool) AddTask(taskID int, command string, task func() (string, error)) {
	wp.tasks <- tasks.Task{ID: taskID, Command: command, Execute: task}
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

func (wp *WorkerPool) Shutdown() {
	close(wp.tasks)
}

func (wp *WorkerPool) Done() <-chan struct{} {
	return wp.done
}

func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		fmt.Printf("\n%d) started: %s\n", task.ID, task.Command)

		res, err := task.Execute()
		wp.results <- tasks.TaskResult{TaskID: task.ID, Result: res, Error: err}

		// stop some workers if workersN decreased
		if workerID > wp.workersN {
			return
		}
	}
}

func (wp *WorkerPool) resultLogger() {
	for task := range wp.results {
		if task.Error != nil {
			fmt.Printf("\n%d) error: %v\n", task.TaskID, task.Error)
		} else {
			fmt.Printf("\n%d) ok %s\n", task.TaskID, task.Result)
		}
	}

	// done if all results have been processed
	wp.done <- struct{}{}
}
