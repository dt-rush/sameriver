package sameriver

type FuncSet struct {
	funcs map[string](func(params any) any)
}

func NewFuncSet(funcs map[string](func(params any) any)) *FuncSet {
	fs := &FuncSet{}
	if funcs == nil {
		funcs = make(map[string](func(params any) any))
	}
	fs.funcs = funcs
	return fs
}

func (fs *FuncSet) Add(name string, f func(params any) any) {
	fs.funcs[name] = f
}

func (fs *FuncSet) Remove(name string) {
	delete(fs.funcs, name)
}

func (fs *FuncSet) Has(name string) bool {
	_, ok := fs.funcs[name]
	return ok
}
