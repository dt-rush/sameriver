package sameriver

type Pooled interface {
	New() any
	Clear(any)
}

type Pool struct {
	// provides new objects if needed (mainly at pool init, but also on expand)
	maker func() any
	// clears objects of their data when they are checked out
	clearer func(any)
	// storage
	pool []any
	// how many objects are in the pool
	capacity int
	// indexes of those which are currently checked out
	checkOutIxs map[any]int
	// indexes in data pool
	avail []int
	// dirty flags telling whether the object at ix has been checked out and returned
	dirty map[int]bool
}

func NewPool(capacity int, maker func() any, clearer func(any)) *Pool {
	p := &Pool{
		capacity:    capacity,
		maker:       maker,
		clearer:     clearer,
		pool:        make([]any, capacity),
		avail:       make([]int, capacity),
		checkOutIxs: make(map[any]int),
		dirty:       make(map[int]bool),
	}
	for i, _ := range p.pool {
		p.pool[i] = maker()
		p.avail[i] = i
		p.dirty[i] = false
	}
	return p
}

func (p *Pool) Clear() {
	for i, _ := range p.pool {
		p.clearer(p.pool[i])
	}
}

func (p *Pool) Expand(n int) {
	emptySpace := make([]any, n)
	newAvails := make([]int, n)
	for i, _ := range emptySpace {
		emptySpace[i] = p.maker()
		newAvails[i] = p.capacity + i
	}
	p.pool = append(p.pool, emptySpace...)
	p.avail = append(p.avail, newAvails...)
	p.capacity += n
}

func (p *Pool) Checkout() any {
	if len(p.avail) == 0 {
		p.Expand(p.capacity / 2)
	}
	ix := p.avail[len(p.avail)-1]
	p.avail = p.avail[:len(p.avail)-1]
	x := p.pool[ix]
	p.checkOutIxs[x] = ix
	if p.dirty[ix] {
		p.clearer(x)
	}
	return x
}

func (p *Pool) Return(x any) {
	ix := p.checkOutIxs[x]
	p.avail = append(p.avail, ix)
	delete(p.checkOutIxs, ix)
	p.dirty[ix] = true
}

func (p *Pool) Avail() int {
	return len(p.avail)
}
