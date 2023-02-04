package sameriver

import (
	"bytes"
	"fmt"
	"math"
	"unsafe"
)

// the actual cell data structure is a cellDimension x cellDimension array of
// entities
type SpatialHashTable [][][]*Entity

// used to compute the spatial hash tables given a list of entities
type SpatialHashSystem struct {
	// the world entities live in
	w *World
	// spatialEntities is an UpdatedEntityList of entities who have position
	// and hitbox components
	spatialEntities *UpdatedEntityList
	// basic data members needed to divide the world into cells and
	// store the entity data in each cell
	GridX int
	GridY int
	Table SpatialHashTable
}

func NewSpatialHashSystem(cellX int, cellY int) *SpatialHashSystem {
	h := SpatialHashSystem{
		GridX: cellX,
		GridY: cellY,
	}
	h.Table = make([][][]*Entity, cellX)
	// for each column (x) in the cell
	for x := 0; x < cellX; x++ {
		h.Table[x] = make([][]*Entity, cellY)
		// for each cell in the row (y)
		for y := 0; y < cellY; y++ {
			h.Table[x][y] = make([]*Entity, 0, MAX_ENTITIES)
		}
	}
	return &h
}

func (s *SpatialHashSystem) GetComponentDeps() []string {
	return []string{"Vec2D,Position", "Vec2D,Box"}
}

func (s *SpatialHashSystem) LinkWorld(w *World) {
	s.w = w
	// get a list of spatial entities
	s.spatialEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray("spatial",
			w.em.components.BitArrayFromNames([]string{"Position", "Box"})))
}

func (h *SpatialHashSystem) Update(dt_ms float64) {
	// clear any old data and run the computation
	h.clearTable()
	h.scanAndInsertEntities()
}

func (h *SpatialHashSystem) clearTable() {
	// NOTE: we "clear" the slice by setting its length to 0 (capacity remains
	// allocated, this will cause a negligible memory "waste" if entities
	// cluster in a cell but never somewhere else. Maybe this could matter if
	// MAX_ENTITIES eventually clustered in each cell, but that's unlikely)
	for x := 0; x < h.GridX; x++ {
		for y := 0; y < h.GridY; y++ {
			cell := &h.Table[x][y]
			*cell = (*cell)[:0]
		}
	}
}

// find out how many cells the box centered at pos spans in x and y
func (h *SpatialHashSystem) CellRangeOfRect(pos, box Vec2D) (cellX0, cellX1, cellY0, cellY1 int) {
	cellSizeX := h.w.Width / float64(h.GridX)
	cellSizeY := h.w.Height / float64(h.GridY)
	clamp := func(x, min, max int) int {
		return int(math.Min(float64(max), math.Max(float64(x), float64(min))))
	}
	cellX0 = clamp(int(pos.X/cellSizeX), 0, h.GridX)
	cellX1 = clamp(int((pos.X+box.X)/cellSizeX), 0, h.GridX)
	cellY0 = clamp(int(pos.Y/cellSizeY), 0, h.GridY)
	cellY1 = clamp(int((pos.Y+box.Y)/cellSizeY), 0, h.GridY)
	return cellX0, cellX1, cellY0, cellY1
}

// used to iterate the entities and send them to the right cells
func (h *SpatialHashSystem) scanAndInsertEntities() {
	for _, e := range h.spatialEntities.entities {
		pos := e.GetVec2D("Position")
		box := e.GetVec2D("Box")

		// walk through each cell the entity touches by
		// starting in the bottom-left and walking cell by cell
		// through each row to the top-right
		cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(pos.ShiftedCenterToBottomLeft(*box), *box)
		for x := cellX0; x <= cellX1; x++ {
			for y := cellY0; y <= cellY1; y++ {
				if x < 0.0 || x > h.GridX-1 ||
					y < 0.0 || y > h.GridY-1 {
					continue
				}
				cell := &h.Table[x][y]
				*cell = append(*cell, e)
			}
		}
	}
}

// get a *copy* of the current table which is safe to hold onto, mutate, etc.
func (h *SpatialHashSystem) TableCopy() SpatialHashTable {
	t2 := make(SpatialHashTable, h.GridX)
	for x := 0; x < h.GridX; x++ {
		t2[x] = make([][]*Entity, h.GridX)
		for y := 0; y < h.GridY; y++ {
			t2[x][y] = make([]*Entity, len(h.Table[x][y]))
			copy(t2[x][y], h.Table[x][y])
		}
	}
	return t2
}

func (h *SpatialHashSystem) GetCellPosAndBox(x, y int) (pos, box Vec2D) {
	cellSizeX := h.w.Width / float64(h.GridX)
	cellSizeY := h.w.Height / float64(h.GridY)
	pos = Vec2D{
		float64(x) * cellSizeX,
		float64(y) * cellSizeY}
	box = Vec2D{
		cellSizeX,
		cellSizeY}
	return pos, box
}

// calculate which cells are within the distance d from the closest
// point on the box centered at pos (imagine a rounded-corner box
// extending d past the limits of the box)
func (h *SpatialHashSystem) CellsWithinDistance(pos, box Vec2D, d float64) [][2]int {
	cells := make([][2]int, 0)

	// first approximate which cells might be valid by simply
	// extending the box by +d in each direction
	approximatorPos := pos.ShiftedCenterToBottomLeft(box).Sub(Vec2D{d, d})
	approximatorBox := box.Add(Vec2D{2 * d, 2 * d})
	cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(approximatorPos, approximatorBox)
	candidateCells := make([][2]int, 0)
	for x := cellX0; x <= cellX1; x++ {
		for y := cellY0; y <= cellY1; y++ {
			candidateCells = append(candidateCells, [2]int{x, y})
		}
	}
	for _, cellXY := range candidateCells {
		cellPos, cellBox := h.GetCellPosAndBox(cellXY[0], cellXY[1])
		if RectWithinDistanceOfRect(cellPos, cellBox, pos, box, d) {
			cells = append(cells, cellXY)
		}
	}
	return cells
}

// extend the box +d on all sides and return the cells it touches
// (NOTE: the corners will slightly over-estimate since they should
// truly be rounded)
// but it's a faster calculation
func (h *SpatialHashSystem) CellsWithinApproxDistance(pos, box Vec2D, d float64) [][2]int {
	approximatorPos := pos.ShiftedCenterToBottomLeft(box).Sub(Vec2D{d, d})
	approximatorBox := box.Add(Vec2D{2 * d, 2 * d})
	cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(approximatorPos, approximatorBox)
	candidateCells := make([][2]int, 0)
	for x := cellX0; x <= cellX1; x++ {
		for y := cellY0; y <= cellY1; y++ {
			candidateCells = append(candidateCells, [2]int{x, y})
		}
	}
	return candidateCells
}

// uses the approx distance since it's faster. Overestimates slightly diagonally.
func (h *SpatialHashSystem) EntitiesWithinDistanceApprox(pos, box Vec2D, d float64) []*Entity {
	results := make([]*Entity, 0)
	cells := h.CellsWithinApproxDistance(pos, box, d)
	// for each cell, append its entities to results
	for _, cell := range cells {
		x := cell[0]
		y := cell[1]
		results = append(results, h.Table[x][y]...)
	}
	return results
}

func (h *SpatialHashSystem) EntitiesWithinDistance(pos, box Vec2D, d float64) []*Entity {
	candidates := h.EntitiesWithinDistanceApprox(pos, box, d)
	results := make([]*Entity, 0)
	for _, e := range candidates {
		ePos := *e.GetVec2D("Position")
		eBox := *e.GetVec2D("Box")
		if RectWithinDistanceOfRect(
			pos.ShiftedCenterToBottomLeft(box), box,
			ePos.ShiftedCenterToBottomLeft(eBox), eBox,
			d) {
			results = append(results, e)
		}
	}
	return results
}

// turn a SpatialHashTable into a String representation (NOTE: do *NOT* call
// this on a pointer returned from CurrentTablePointer unless you can be sure
// that you have not called Update more than once - it does not
// lock the table it reads, and if you call Update twice, you may
// start to write to the table as this function reads it). Usually best to call
// on a Copy()
func (table *SpatialHashTable) String() string {
	var buffer bytes.Buffer
	w := len(*table)
	h := len((*table)[0])
	size := int(unsafe.Sizeof(*table))
	buffer.WriteString("[\n")
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cell := (*table)[x][y]
			size += int(unsafe.Sizeof(&Entity{})) * len(cell)
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %.64s...", x, y,
				EntitySliceToString(cell)))
			if !(y == h-1 && x == w-1) {
				buffer.WriteString(",\n")
			}
		}
	}
	buffer.WriteString(fmt.Sprintf("] (using %d bytes)", size))
	return buffer.String()
}

/*
func (h *SpatialHashSystem) EntitiesWithinDistance(pos, box Vec2D, d float64) []*Entity {
	// algorithm from stackoverflow user Nick Alger
	// https://stackoverflow.com/a/65107290
	boxDist := func(aMin, aMax, bMin, bMax Vec2D) float64 {
		entrywiseMaxZero := func(vec Vec2D) Vec2D {
			return Vec2D{
				math.Max(0, vec.X),
				math.Max(0, vec.Y),
			}
		}
		euclidNorm := func(vec Vec2D) float64 {
			return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y)
		}
		u := entrywiseMaxZero(aMin.Sub(bMax))
		v := entrywiseMaxZero(bMin.Sub(aMax))
		unorm := euclidNorm(u)
		vnorm := euclidNorm(v)
		return math.Sqrt(unorm*unorm + vnorm*vnorm)
	}
	inDist := func(j *Entity) bool {
		// place the boxes in space according to the position
		iBox := box
		jBox := j.GetVec2D("Box")
		iPos := pos
		jPos := j.GetVec2D("Position")
		// lower-left and upper-right corners
		iMin := iPos.ShiftedCenterToBottomLeft(&iBox)
		iMax := iMin.Add(iBox)
		jMin := jPos.ShiftedCenterToBottomLeft(jBox)
		jMax := jMin.Add(*jBox)
		return boxDist(iMin, iMax, jMin, jMax) < d
	}

	candidates := h.EntitiesPotentiallyWithinDistance(pos, box, d)
	results := make([]*Entity, 0)
	for _, e := range candidates {
		if inDist(e) {
			results = append(results, e)
		}
	}
	return results
}
*/
