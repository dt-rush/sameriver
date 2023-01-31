package engine

type FuncSet struct {
	funcs map[string](func(interface{}) interface{})
}

func NewFuncSet() *FuncSet {
	fs := &FuncSet{}
	fs.funcs = make(map[string](func(interface{}) interface{}))
	return fs
}

func (fs *FuncSet) Add(name string, f func(interface{}) interface{}) {
	fs.funcs[name] = f
}

func (fs *FuncSet) Remove(name string) {
	delete(fs.funcs, name)
}

func (fs *FuncSet) Has(name string) bool {
	_, ok := fs.funcs[name]
	return ok
}
