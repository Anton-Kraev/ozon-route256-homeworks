package module

type Storage interface {
}

type Deps struct {
	Storage Storage
}

type Module struct {
	Deps
}

func NewModule(d Deps) Module {
	return Module{Deps: d}
}
