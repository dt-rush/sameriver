package engine

type LogicUnit struct {
	Name   string
	Active bool
	F      func()
}
