package sameriver

// used so that game code using the engine can define their own contiguous
// arrays of complex struct types and take advantage of cache-line performance
// rather than use the `Generic` component type which stores `any`
// (thus involves dynamic allocation peppered through memory that won't have
// nice contigous cache lines)
//
// the user should define the Allocate() function so that it allocates the
// contiguous memory (slice), and Name() should return the name, for example,
// if the CCC is holding StateTree objects, the name should be "StateTree"
//
// notice that the accessor method is defined for generic pointer return,
// which corresponds to the Entity.GetCustom() method defined in component_table.go
//
// so, user code should cast to the appropriate type on retrieval
//
// for example:
//
// st := e.GetCustom("StateTree").(*StateTree)
type CustomContiguousComponent interface {
	Name() string
	AllocateTable(n int)
	ExpandTable(n int)
	Get(e *Entity) any
	Set(e *Entity, x any)
}
