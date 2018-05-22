/**
  *
  *
  *
  *
**/

package engine

type IDGenerator struct {
	x uint16
}

func (g *IDGenerator) Gen() uint16 {
	id := g.x
	g.x += 1
	return id
}

var IDGEN_OBJ = IDGenerator{}
var IDGEN = IDGEN_OBJ.Gen
