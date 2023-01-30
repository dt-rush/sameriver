package engine

import (
	"bytes"
	"fmt"
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
	s.spatialEntities = w.em.GetUpdatedEntityList(
		EntityFilterFromComponentBitArray("spatial",
			w.em.components.BitArrayFromNames([]string{"Position", "Box"})))
}

func (h *SpatialHashSystem) Update() {
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

// used to iterate the entities and send them to the right cells
func (h *SpatialHashSystem) scanAndInsertEntities() {
	cellSizeX := h.w.Width / float64(h.GridX)
	cellSizeY := h.w.Height / float64(h.GridX)
	for _, e := range h.spatialEntities.entities {
		// we shift the position to the bottom-left because
		// the logic is simpler to read that way
		pos := e.GetVec2D("Position")
		box := e.GetVec2D("Box")
		pos.ShiftCenterToBottomLeft(box)
		defer pos.ShiftBottomLeftToCenter(box)
		// find out how many cells the entity spans in x and y (almost always 0,
		// but we want to be thorough, and the fact that it's got a predictable
		// pattern 99% of the time means that branch prediction should help us)
		cellsWide := box.X / cellSizeX
		cellsHigh := box.Y / cellSizeY
		cellX := pos.X / cellSizeX
		cellY := pos.Y / cellSizeY
		// walk through each cell the entity touches by starting in the bottom-
		// -left and walking according to cellsHigh and cellsWide
		for ix := 0.0; ix < cellsWide+1; ix++ {
			for iy := 0.0; iy < cellsHigh+1; iy++ {
				x := int(cellX + ix)
				y := int(cellY + iy)
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
