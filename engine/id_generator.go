/**
  *
  *
  *
  *
**/

package engine

type IDGenerator struct {
	// records all ids given out
	ids []int

	// used to gen
	seed int
}

func (g *IDGenerator) Init () {
	// can be tuned
	capacity := 32
	g.ids = make ([]int, capacity)
	g.seed = -1
}

func (g *IDGenerator) Gen() int {
	g.seed++
	g.ids = append (g.ids, g.seed)
	return g.seed
}

func (g *IDGenerator) GetIDs() []int {
	return g.ids
}

func (g *IDGenerator) NumberOfIDs() int {
	return len (g.ids)
}
