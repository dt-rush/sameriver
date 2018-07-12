package utils

import (
	"go.uber.org/atomic"
)

type IDGenerator struct {
	x atomic.Uint32
}

func (g *IDGenerator) Gen() int {
	return int(g.x.Inc())
}

var IDGEN_OBJ = IDGenerator{}
var IDGEN = IDGEN_OBJ.Gen
