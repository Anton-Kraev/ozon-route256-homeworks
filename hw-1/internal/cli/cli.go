package cli

type Module interface {
}

type Deps struct {
	Module Module
}

type CLI struct {
	Deps
	availableCommands []command
}

func NewCLI(d Deps) CLI {
	return CLI{
		Deps:              d,
		availableCommands: commandsList,
	}
}

func (c CLI) Run() error {
	return nil
}
