package engine

type LogicUnit struct {
	Name        string
	f           EntityLogicFunc
	StopChannel chan bool
}

// Create a new LogicUnit instance
func NewLogicUnit(Name string, f EntityLogicFunc) LogicUnit {
	return LogicUnit{
		Name,
		f,
		make(chan bool)}
}
