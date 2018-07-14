package engine

import (
	"reflect"
)

func fullZeroedComponentSet() ComponentSet {
	cs := ComponentSet{}
	v := reflect.ValueOf(&cs).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		val := reflect.New(f.Type().Elem())
		f.Set(val)
	}
	return cs
}
