package engine

type EntityLogicUnit struct {
	active bool
	f      func()
}
