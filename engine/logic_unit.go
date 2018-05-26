package engine

type LogicUnit struct {
	f           EntityLogicFunc
	Name        string
	StopChannel chan bool
}

// Create a new LogicUnit instance
func NewLogicUnit(Name string, f EntityLogicFunc) LogicUnit {
	return LogicUnit{
		f,
		Name,
		make(chan bool)}
}
