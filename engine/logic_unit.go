package engine

type LogicUnit struct {
	active bool
	f      func()
}
