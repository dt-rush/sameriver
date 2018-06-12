package main

func (c *WorldMapCell) Neighbors() []*WorldMapCell {
	neighbors := make([]*WorldMapCell, 0)
	for dy := -1; dy <= 1; dy++ {
		if c.pos.y+dy < 0 ||
			c.pos.y+dy > WORLD_CELLHEIGHT-1 {
			continue
		}
		for dx := -1; dx <= 1; dx++ {
			if c.pos.x+dx < 0 ||
				c.pos.x+dx > WORLD_CELLWIDTH-1 {
				continue
			}
			neighbors = append(neighbors,
				&c.m.cells[c.pos.y+dy][c.pos.x+dx])
		}
	}
	return neighbors
}
