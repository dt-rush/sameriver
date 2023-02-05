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

type SpatialHasher struct {
	// spatialEntities is an UpdatedEntityList of entities who have position
	// and hitbox components
	SpatialEntities *UpdatedEntityList
	// basic data members needed to divide the world into cells and
	// store the entity data in each cell
	GridX     int
	GridY     int
	CellSizeX float64
	CellSizeY float64
	// table of cells, GridX x GridY, that holds the entities
	Table SpatialHashTable
}

func NewSpatialHasher(gridX, gridY int, w *World) *SpatialHasher {
	h := &SpatialHasher{
		GridX:     gridX,
		GridY:     gridY,
		CellSizeX: w.Width / float64(gridX),
		CellSizeY: w.Height / float64(gridY),
	}
	h.Table = make([][][]*Entity, gridX)
	// for each column (x)
	for x := 0; x < gridX; x++ {
		h.Table[x] = make([][]*Entity, gridY)
		// for each cell in the row (y)
		for y := 0; y < gridY; y++ {
			h.Table[x][y] = make([]*Entity, 0, MAX_ENTITIES/4)
		}
	}
	// get spatial entities from world
	h.SpatialEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray("spatial",
			w.em.components.BitArrayFromNames([]string{"Position", "Box"})))

	return h
}

func (h *SpatialHasher) Entities(x, y int) []*Entity {
	return h.Table[x][y]
}

func (h *SpatialHasher) ClearTable() {
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

func (h *SpatialHasher) CellRangeOfRect(pos, box Vec2D) (cellX0, cellX1, cellY0, cellY1 int) {
	clamp := func(x, min, max int) int {
		return int(math.Min(float64(max), math.Max(float64(x), float64(min))))
	}
	// TODO: don't clamp here, test for out of bounds and `continue` in receivers
	cellX0 = clamp(int(pos.X/h.CellSizeX), 0, h.GridX)
	cellX1 = clamp(int((pos.X+box.X)/h.CellSizeX), 0, h.GridX)
	cellY0 = clamp(int(pos.Y/h.CellSizeY), 0, h.GridY)
	cellY1 = clamp(int((pos.Y+box.Y)/h.CellSizeY), 0, h.GridY)
	return cellX0, cellX1, cellY0, cellY1
}

// find out how many cells the box centered at pos spans in x and y
// used to iterate the entities and send them to the right cells
func (h *SpatialHasher) ScanAndInsertEntities() {
	for _, e := range h.SpatialEntities.entities {
		pos := e.GetVec2D("Position")
		box := e.GetVec2D("Box")

		// walk through each cell the entity touches by
		// starting in the bottom-left and walking cell by cell
		// through each row to the top-right
		cellX0, cellX1, cellY0, cellY1 := h.CellRangeOfRect(pos.ShiftedCenterToBottomLeft(*box), *box)
		for x := cellX0; x <= cellX1; x++ {
			for y := cellY0; y <= cellY1; y++ {
				if x < 0 || x > h.GridX-1 ||
					y < 0 || y > h.GridY-1 {
					continue
				}
				cell := &h.Table[x][y]
				*cell = append(*cell, e)
			}
		}
	}
}

// get a *copy* of the current table which is safe to hold onto, mutate, etc.
func (h *SpatialHasher) TableCopy() SpatialHashTable {
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

func (h *SpatialHasher) GetCellPosAndBox(x, y int) (pos, box Vec2D) {
	pos = Vec2D{
		float64(x) * h.CellSizeX,
		float64(y) * h.CellSizeY}
	box = Vec2D{
		h.CellSizeX,
		h.CellSizeY}
	return pos, box
}

// calculate which cells are within the distance d from the closest
// point on the box centered at pos (imagine a rounded-corner box
// extending d past the limits of the box)
func (h *SpatialHasher) CellsWithinDistance(pos, box Vec2D, d float64) [][2]int {
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
func (h *SpatialHasher) CellsWithinDistanceApprox(pos, box Vec2D, d float64) [][2]int {
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
func (h *SpatialHasher) EntitiesWithinDistanceApprox(pos, box Vec2D, d float64) []*Entity {
	results := make([]*Entity, 0)
	found := make(map[int]*Entity)
	cells := h.CellsWithinDistanceApprox(pos, box, d)
	// for each cell, note that we found the entity in the map
	// initial naive slice append was quadruple-counting entities that
	// sat on the intersection of four cells, etc.
	for _, cell := range cells {
		x := cell[0]
		y := cell[1]
		for _, e := range h.Table[x][y] {
			found[e.ID] = e
		}
	}
	for _, e := range found {
		results = append(results, e)
	}
	return results
}

func (h *SpatialHasher) EntitiesWithinDistance(pos, box Vec2D, d float64) []*Entity {
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
func (h *SpatialHasher) String() string {
	var buffer bytes.Buffer
	size := int(unsafe.Sizeof(h.Table))
	buffer.WriteString("[\n")
	for x := 0; x < h.GridX; x++ {
		for y := 0; y < h.GridY; y++ {
			cell := h.Table[x][y]
			size += int(unsafe.Sizeof(&Entity{})) * cap(cell)
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %.64s...", x, y,
				EntitySliceToString(cell)))
			if !(y == h.GridY-1 && x == h.GridX-1) {
				buffer.WriteString(",\n")
			}
		}
	}
	buffer.WriteString(fmt.Sprintf("] (using %d bytes)", size))
	return buffer.String()
}
