package engine

type LogicUnit struct {
	Name    string
	F       func()
	Active  bool
	WorldID int
}
