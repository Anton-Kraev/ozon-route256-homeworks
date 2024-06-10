package tasks

type Task struct {
	ID      int
	Command string
	Execute func() (string, error)
}

type TaskResult struct {
	TaskID int
	Result string
	Error  error
}
