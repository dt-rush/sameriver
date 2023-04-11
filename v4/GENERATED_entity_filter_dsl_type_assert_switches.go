package sameriver

func (e *EntityFilterDSLEvaluator) predicateSignatureAssertSwitch(f any, argsTyped []any) func(*Entity) bool {
	var result func(*Entity) bool
	switch fTyped := f.(type) {
	case func(bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool))
	case func(int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int))
	case func(string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string))
	case func(*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity))
	case func([]*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity))
	case func(*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D))
	case func([]*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D))
	case func(bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool))
	case func(int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool))
	case func(string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool))
	case func(*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool))
	case func([]*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool))
	case func(*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool))
	case func([]*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool))
	case func(bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int))
	case func(int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int))
	case func(string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int))
	case func(*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int))
	case func([]*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int))
	case func(*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int))
	case func([]*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int))
	case func(bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string))
	case func(int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string))
	case func(string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string))
	case func(*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string))
	case func([]*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string))
	case func(*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string))
	case func([]*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string))
	case func(bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity))
	case func(int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity))
	case func(string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity))
	case func(*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity))
	case func([]*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity))
	case func(*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity))
	case func([]*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity))
	case func(bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity))
	case func(int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity))
	case func(string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity))
	case func(*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity))
	case func([]*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity))
	case func(*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity))
	case func([]*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity))
	case func(bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D))
	case func(int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D))
	case func(string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D))
	case func(*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D))
	case func([]*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D))
	case func(*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D))
	case func([]*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D))
	case func(bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D))
	case func(int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D))
	case func(string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D))
	case func(*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D))
	case func([]*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D))
	case func(*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D))
	case func([]*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D))
	case func(bool, bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(int, bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(string, bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(*Entity, bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(bool))
	case func([]*Entity, bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(*Vec2D, bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(bool))
	case func([]*Vec2D, bool, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(bool, int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(bool))
	case func(int, int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(bool))
	case func(string, int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(bool))
	case func(*Entity, int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(bool))
	case func([]*Entity, int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(bool))
	case func(*Vec2D, int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(bool))
	case func([]*Vec2D, int, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(bool))
	case func(bool, string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(bool))
	case func(int, string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(bool))
	case func(string, string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(bool))
	case func(*Entity, string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(bool))
	case func([]*Entity, string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(bool))
	case func(*Vec2D, string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(bool))
	case func([]*Vec2D, string, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(bool))
	case func(bool, *Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(int, *Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(string, *Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(*Entity, *Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func([]*Entity, *Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(*Vec2D, *Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func([]*Vec2D, *Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(bool, []*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(int, []*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(string, []*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(*Entity, []*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func([]*Entity, []*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(*Vec2D, []*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func([]*Vec2D, []*Entity, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(bool, *Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(int, *Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(string, *Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(*Entity, *Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func([]*Entity, *Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(*Vec2D, *Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func([]*Vec2D, *Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(bool, []*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(int, []*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(string, []*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(*Entity, []*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func([]*Entity, []*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(*Vec2D, []*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func([]*Vec2D, []*Vec2D, bool) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(bool, bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(int))
	case func(int, bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(int))
	case func(string, bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(int))
	case func(*Entity, bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(int))
	case func([]*Entity, bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(int))
	case func(*Vec2D, bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(int))
	case func([]*Vec2D, bool, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(int))
	case func(bool, int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(int))
	case func(int, int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(int))
	case func(string, int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(int))
	case func(*Entity, int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(int))
	case func([]*Entity, int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(int))
	case func(*Vec2D, int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(int))
	case func([]*Vec2D, int, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(int))
	case func(bool, string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(int))
	case func(int, string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(int))
	case func(string, string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(int))
	case func(*Entity, string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(int))
	case func([]*Entity, string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(int))
	case func(*Vec2D, string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(int))
	case func([]*Vec2D, string, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(int))
	case func(bool, *Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(int, *Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(string, *Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(*Entity, *Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func([]*Entity, *Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(*Vec2D, *Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func([]*Vec2D, *Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(bool, []*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(int, []*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(string, []*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(*Entity, []*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func([]*Entity, []*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(*Vec2D, []*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func([]*Vec2D, []*Entity, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(bool, *Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(int, *Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(string, *Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(*Entity, *Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func([]*Entity, *Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(*Vec2D, *Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func([]*Vec2D, *Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(bool, []*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(int, []*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(string, []*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(*Entity, []*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func([]*Entity, []*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(*Vec2D, []*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func([]*Vec2D, []*Vec2D, int) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(bool, bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(string))
	case func(int, bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(string))
	case func(string, bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(string))
	case func(*Entity, bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(string))
	case func([]*Entity, bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(string))
	case func(*Vec2D, bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(string))
	case func([]*Vec2D, bool, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(string))
	case func(bool, int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(string))
	case func(int, int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(string))
	case func(string, int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(string))
	case func(*Entity, int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(string))
	case func([]*Entity, int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(string))
	case func(*Vec2D, int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(string))
	case func([]*Vec2D, int, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(string))
	case func(bool, string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(string))
	case func(int, string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(string))
	case func(string, string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(string))
	case func(*Entity, string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(string))
	case func([]*Entity, string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(string))
	case func(*Vec2D, string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(string))
	case func([]*Vec2D, string, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(string))
	case func(bool, *Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(int, *Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(string, *Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(*Entity, *Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func([]*Entity, *Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(*Vec2D, *Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func([]*Vec2D, *Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(bool, []*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(int, []*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(string, []*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(*Entity, []*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func([]*Entity, []*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(*Vec2D, []*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func([]*Vec2D, []*Entity, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(bool, *Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(int, *Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(string, *Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(*Entity, *Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func([]*Entity, *Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(*Vec2D, *Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func([]*Vec2D, *Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(bool, []*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(int, []*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(string, []*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(*Entity, []*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func([]*Entity, []*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(*Vec2D, []*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func([]*Vec2D, []*Vec2D, string) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(bool, bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(int, bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(string, bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(*Entity, bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func([]*Entity, bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(*Vec2D, bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func([]*Vec2D, bool, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(bool, int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(int, int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(string, int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(*Entity, int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func([]*Entity, int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(*Vec2D, int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func([]*Vec2D, int, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(bool, string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(int, string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(string, string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(*Entity, string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func([]*Entity, string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(*Vec2D, string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func([]*Vec2D, string, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(bool, *Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(int, *Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(string, *Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(*Entity, *Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func([]*Entity, *Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(*Vec2D, *Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func([]*Vec2D, *Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(bool, []*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(int, []*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(string, []*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(*Entity, []*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func([]*Entity, []*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(*Vec2D, []*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func([]*Vec2D, []*Entity, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(bool, *Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(int, *Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(string, *Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(*Entity, *Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func([]*Entity, *Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(*Vec2D, *Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func([]*Vec2D, *Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(bool, []*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(int, []*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(string, []*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(*Entity, []*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func([]*Entity, []*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(*Vec2D, []*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func([]*Vec2D, []*Vec2D, *Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(bool, bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(int, bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(string, bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(*Entity, bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func([]*Entity, bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(*Vec2D, bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func([]*Vec2D, bool, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(bool, int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(int, int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(string, int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(*Entity, int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func([]*Entity, int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(*Vec2D, int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func([]*Vec2D, int, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(bool, string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(int, string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(string, string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(*Entity, string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func([]*Entity, string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(*Vec2D, string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func([]*Vec2D, string, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(bool, *Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(int, *Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(string, *Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(*Entity, *Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func([]*Entity, *Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(*Vec2D, *Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func([]*Vec2D, *Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(bool, []*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(int, []*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(string, []*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(*Entity, []*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func([]*Entity, []*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(*Vec2D, []*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func([]*Vec2D, []*Entity, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(bool, *Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(int, *Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(string, *Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(*Entity, *Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Entity, *Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(*Vec2D, *Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Vec2D, *Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(bool, []*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(int, []*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(string, []*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(*Entity, []*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Entity, []*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(*Vec2D, []*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Vec2D, []*Vec2D, []*Entity) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(bool, bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(int, bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(string, bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(*Entity, bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func([]*Entity, bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(*Vec2D, bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, bool, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(bool, int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(int, int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(string, int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(*Entity, int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func([]*Entity, int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(*Vec2D, int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, int, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(bool, string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(int, string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(string, string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(*Entity, string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func([]*Entity, string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(*Vec2D, string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, string, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(bool, *Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(int, *Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(string, *Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(*Entity, *Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func([]*Entity, *Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(*Vec2D, *Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, *Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(bool, []*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(int, []*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(string, []*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(*Entity, []*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func([]*Entity, []*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(*Vec2D, []*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, []*Entity, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(bool, *Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(int, *Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(string, *Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Entity, *Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Entity, *Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Vec2D, *Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, *Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(bool, []*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(int, []*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(string, []*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Entity, []*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Entity, []*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Vec2D, []*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, []*Vec2D, *Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(bool, bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(int, bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(string, bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(*Entity, bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func([]*Entity, bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, bool, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(bool, int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(int, int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(string, int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(*Entity, int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func([]*Entity, int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, int, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(bool, string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(int, string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(string, string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(*Entity, string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func([]*Entity, string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, string, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(bool, *Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(int, *Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(string, *Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(*Entity, *Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Entity, *Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, *Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, *Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(bool, []*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(int, []*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(string, []*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(*Entity, []*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Entity, []*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, []*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, []*Entity, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(bool, *Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(int, *Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(string, *Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Entity, *Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Entity, *Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, *Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, *Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(bool, []*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(int, []*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(string, []*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Entity, []*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Entity, []*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, []*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, []*Vec2D, []*Vec2D) func(*Entity) bool:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	default:
		panic("No case in either engine or user-registered signatures for the given func. Use EFDSL.RegisterUserPredicateSignatureAsserter()")
	}

	return result
}

func (e *EntityFilterDSLEvaluator) sortSignatureAssertSwitch(f any, argsTyped []any) func(xs []*Entity) func(i, j int) int {
	var result func(xs []*Entity) func(i, j int) int
	switch fTyped := f.(type) {
	case func(bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool))
	case func(int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int))
	case func(string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string))
	case func(*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity))
	case func([]*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity))
	case func(*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D))
	case func([]*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D))
	case func(bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool))
	case func(int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool))
	case func(string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool))
	case func(*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool))
	case func([]*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool))
	case func(*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool))
	case func([]*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool))
	case func(bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int))
	case func(int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int))
	case func(string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int))
	case func(*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int))
	case func([]*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int))
	case func(*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int))
	case func([]*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int))
	case func(bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string))
	case func(int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string))
	case func(string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string))
	case func(*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string))
	case func([]*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string))
	case func(*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string))
	case func([]*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string))
	case func(bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity))
	case func(int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity))
	case func(string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity))
	case func(*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity))
	case func([]*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity))
	case func(*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity))
	case func([]*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity))
	case func(bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity))
	case func(int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity))
	case func(string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity))
	case func(*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity))
	case func([]*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity))
	case func(*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity))
	case func([]*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity))
	case func(bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D))
	case func(int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D))
	case func(string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D))
	case func(*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D))
	case func([]*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D))
	case func(*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D))
	case func([]*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D))
	case func(bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D))
	case func(int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D))
	case func(string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D))
	case func(*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D))
	case func([]*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D))
	case func(*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D))
	case func([]*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D))
	case func(bool, bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(int, bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(string, bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(*Entity, bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(bool))
	case func([]*Entity, bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(*Vec2D, bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(bool))
	case func([]*Vec2D, bool, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(bool))
	case func(bool, int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(bool))
	case func(int, int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(bool))
	case func(string, int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(bool))
	case func(*Entity, int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(bool))
	case func([]*Entity, int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(bool))
	case func(*Vec2D, int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(bool))
	case func([]*Vec2D, int, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(bool))
	case func(bool, string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(bool))
	case func(int, string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(bool))
	case func(string, string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(bool))
	case func(*Entity, string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(bool))
	case func([]*Entity, string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(bool))
	case func(*Vec2D, string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(bool))
	case func([]*Vec2D, string, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(bool))
	case func(bool, *Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(int, *Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(string, *Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(*Entity, *Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func([]*Entity, *Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(*Vec2D, *Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func([]*Vec2D, *Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(bool))
	case func(bool, []*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(int, []*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(string, []*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(*Entity, []*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func([]*Entity, []*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(*Vec2D, []*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func([]*Vec2D, []*Entity, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(bool))
	case func(bool, *Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(int, *Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(string, *Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(*Entity, *Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func([]*Entity, *Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(*Vec2D, *Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func([]*Vec2D, *Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(bool))
	case func(bool, []*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(int, []*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(string, []*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(*Entity, []*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func([]*Entity, []*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(*Vec2D, []*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func([]*Vec2D, []*Vec2D, bool) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(bool))
	case func(bool, bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(int))
	case func(int, bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(int))
	case func(string, bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(int))
	case func(*Entity, bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(int))
	case func([]*Entity, bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(int))
	case func(*Vec2D, bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(int))
	case func([]*Vec2D, bool, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(int))
	case func(bool, int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(int))
	case func(int, int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(int))
	case func(string, int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(int))
	case func(*Entity, int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(int))
	case func([]*Entity, int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(int))
	case func(*Vec2D, int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(int))
	case func([]*Vec2D, int, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(int))
	case func(bool, string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(int))
	case func(int, string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(int))
	case func(string, string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(int))
	case func(*Entity, string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(int))
	case func([]*Entity, string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(int))
	case func(*Vec2D, string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(int))
	case func([]*Vec2D, string, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(int))
	case func(bool, *Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(int, *Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(string, *Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(*Entity, *Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func([]*Entity, *Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(*Vec2D, *Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func([]*Vec2D, *Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(int))
	case func(bool, []*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(int, []*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(string, []*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(*Entity, []*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func([]*Entity, []*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(*Vec2D, []*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func([]*Vec2D, []*Entity, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(int))
	case func(bool, *Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(int, *Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(string, *Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(*Entity, *Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func([]*Entity, *Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(*Vec2D, *Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func([]*Vec2D, *Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(int))
	case func(bool, []*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(int, []*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(string, []*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(*Entity, []*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func([]*Entity, []*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(*Vec2D, []*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func([]*Vec2D, []*Vec2D, int) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(int))
	case func(bool, bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(string))
	case func(int, bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(string))
	case func(string, bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(string))
	case func(*Entity, bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(string))
	case func([]*Entity, bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(string))
	case func(*Vec2D, bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(string))
	case func([]*Vec2D, bool, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(string))
	case func(bool, int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(string))
	case func(int, int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(string))
	case func(string, int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(string))
	case func(*Entity, int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(string))
	case func([]*Entity, int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(string))
	case func(*Vec2D, int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(string))
	case func([]*Vec2D, int, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(string))
	case func(bool, string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(string))
	case func(int, string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(string))
	case func(string, string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(string))
	case func(*Entity, string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(string))
	case func([]*Entity, string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(string))
	case func(*Vec2D, string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(string))
	case func([]*Vec2D, string, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(string))
	case func(bool, *Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(int, *Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(string, *Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(*Entity, *Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func([]*Entity, *Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(*Vec2D, *Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func([]*Vec2D, *Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(string))
	case func(bool, []*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(int, []*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(string, []*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(*Entity, []*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func([]*Entity, []*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(*Vec2D, []*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func([]*Vec2D, []*Entity, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(string))
	case func(bool, *Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(int, *Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(string, *Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(*Entity, *Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func([]*Entity, *Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(*Vec2D, *Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func([]*Vec2D, *Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(string))
	case func(bool, []*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(int, []*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(string, []*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(*Entity, []*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func([]*Entity, []*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(*Vec2D, []*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func([]*Vec2D, []*Vec2D, string) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(string))
	case func(bool, bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(int, bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(string, bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(*Entity, bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func([]*Entity, bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(*Vec2D, bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func([]*Vec2D, bool, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Entity))
	case func(bool, int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(int, int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(string, int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(*Entity, int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func([]*Entity, int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(*Vec2D, int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func([]*Vec2D, int, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(*Entity))
	case func(bool, string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(int, string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(string, string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(*Entity, string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func([]*Entity, string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(*Vec2D, string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func([]*Vec2D, string, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(*Entity))
	case func(bool, *Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(int, *Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(string, *Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(*Entity, *Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func([]*Entity, *Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(*Vec2D, *Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func([]*Vec2D, *Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Entity))
	case func(bool, []*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(int, []*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(string, []*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(*Entity, []*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func([]*Entity, []*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(*Vec2D, []*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func([]*Vec2D, []*Entity, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Entity))
	case func(bool, *Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(int, *Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(string, *Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(*Entity, *Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func([]*Entity, *Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(*Vec2D, *Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func([]*Vec2D, *Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Entity))
	case func(bool, []*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(int, []*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(string, []*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(*Entity, []*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func([]*Entity, []*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(*Vec2D, []*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func([]*Vec2D, []*Vec2D, *Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Entity))
	case func(bool, bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(int, bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(string, bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(*Entity, bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func([]*Entity, bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(*Vec2D, bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func([]*Vec2D, bool, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Entity))
	case func(bool, int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(int, int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(string, int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(*Entity, int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func([]*Entity, int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(*Vec2D, int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func([]*Vec2D, int, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Entity))
	case func(bool, string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(int, string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(string, string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(*Entity, string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func([]*Entity, string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(*Vec2D, string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func([]*Vec2D, string, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Entity))
	case func(bool, *Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(int, *Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(string, *Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(*Entity, *Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func([]*Entity, *Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(*Vec2D, *Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func([]*Vec2D, *Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Entity))
	case func(bool, []*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(int, []*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(string, []*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(*Entity, []*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func([]*Entity, []*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(*Vec2D, []*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func([]*Vec2D, []*Entity, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Entity))
	case func(bool, *Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(int, *Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(string, *Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(*Entity, *Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Entity, *Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(*Vec2D, *Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Vec2D, *Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Entity))
	case func(bool, []*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(int, []*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(string, []*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(*Entity, []*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Entity, []*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(*Vec2D, []*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func([]*Vec2D, []*Vec2D, []*Entity) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Entity))
	case func(bool, bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(int, bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(string, bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(*Entity, bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func([]*Entity, bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(*Vec2D, bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, bool, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].(*Vec2D))
	case func(bool, int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(int, int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(string, int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(*Entity, int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func([]*Entity, int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(*Vec2D, int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, int, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].(*Vec2D))
	case func(bool, string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(int, string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(string, string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(*Entity, string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func([]*Entity, string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(*Vec2D, string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, string, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].(*Vec2D))
	case func(bool, *Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(int, *Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(string, *Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(*Entity, *Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func([]*Entity, *Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(*Vec2D, *Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, *Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].(*Vec2D))
	case func(bool, []*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(int, []*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(string, []*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(*Entity, []*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func([]*Entity, []*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(*Vec2D, []*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, []*Entity, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].(*Vec2D))
	case func(bool, *Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(int, *Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(string, *Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Entity, *Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Entity, *Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Vec2D, *Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, *Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].(*Vec2D))
	case func(bool, []*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(int, []*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(string, []*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Entity, []*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Entity, []*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(*Vec2D, []*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func([]*Vec2D, []*Vec2D, *Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].(*Vec2D))
	case func(bool, bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(int, bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(string, bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(*Entity, bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func([]*Entity, bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, bool, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(bool), argsTyped[2].([]*Vec2D))
	case func(bool, int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(int, int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(string, int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(*Entity, int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func([]*Entity, int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, int, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(int), argsTyped[2].([]*Vec2D))
	case func(bool, string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(int, string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(string, string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(*Entity, string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func([]*Entity, string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, string, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(string), argsTyped[2].([]*Vec2D))
	case func(bool, *Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(int, *Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(string, *Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(*Entity, *Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Entity, *Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, *Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, *Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Entity), argsTyped[2].([]*Vec2D))
	case func(bool, []*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(int, []*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(string, []*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(*Entity, []*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Entity, []*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, []*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, []*Entity, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Entity), argsTyped[2].([]*Vec2D))
	case func(bool, *Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(int, *Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(string, *Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Entity, *Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Entity, *Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, *Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, *Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].(*Vec2D), argsTyped[2].([]*Vec2D))
	case func(bool, []*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(bool), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(int, []*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(int), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(string, []*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(string), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Entity, []*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Entity, []*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Entity), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func(*Vec2D, []*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].(*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	case func([]*Vec2D, []*Vec2D, []*Vec2D) func(xs []*Entity) func(i, j int) int:
		result = fTyped(argsTyped[0].([]*Vec2D), argsTyped[1].([]*Vec2D), argsTyped[2].([]*Vec2D))
	default:
		panic("No case in either engine or user-registered signatures for the given func. Use EFDSL.RegisterUserSortSignatureAsserter()")
	}

	return result
}
