package engine

type LogicUnit struct {
	name    string
	f       func()
	active  bool
	worldID int
}
