package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/dt-rush/donkeys-qquest/engine"
	"github.com/golang-collections/go-datastructures/bitarray"
)

var N_ENTITIES = uint64(1024)
var N_COMPONENTS = uint64(16)

func BitArrayToString(arr bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := uint64(0); i < N_COMPONENTS; i++ {
		bit, _ := arr.GetBit(i)
		var val int
		if bit {
			val = 1
		} else {
			val = 0
		}
		buf.WriteString(fmt.Sprintf("%v", val))
	}
	buf.WriteString("]")
	return buf.String()
}

func CreateBitArrayQuery(indexes []int) bitarray.BitArray {
	query := bitarray.NewBitArray(uint64(N_COMPONENTS))
	for i := 0; i < len(indexes); i++ {
		query.SetBit(uint64(indexes[i]))
	}
	return query
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	func_profiler := engine.NewFuncProfiler(engine.FUNC_PROFILER_SIMPLE)
	query_profiler_id := func_profiler.RegisterFunc("query")

	entity_components := make([]bitarray.BitArray, N_ENTITIES)
	// Init
	for i := uint64(0); i < N_ENTITIES; i++ {
		entity_components[i] = bitarray.NewBitArray(uint64(N_COMPONENTS))
	}
	// Set
	for i := uint64(0); i < N_ENTITIES; i++ {
		for j := uint64(0); j < N_COMPONENTS; j++ {
			if rand.Intn(2) == 1 {
				entity_components[i].SetBit(uint64(j))
			}
		}
	}
	// Print
	for i := 0; i < len(entity_components); i++ {
		fmt.Printf("%v\n", BitArrayToString(entity_components[i]))
	}
	// Query
	q1 := CreateBitArrayQuery([]int{1, 2, 3, 4, 5})
	fmt.Printf("Running query for %s\n", BitArrayToString(q1))
	func_profiler.StartTimer(query_profiler_id)
	matching := make([]bitarray.BitArray, 0)
	for i := uint64(0); i < N_ENTITIES; i++ {
		if entity_components[i].Equals(q1) {
			fmt.Printf("%s matches\n", BitArrayToString(entity_components[i]))
			matching = append(matching, entity_components[i])
		}
	}
	func_profiler.EndTimer(query_profiler_id)
	fmt.Printf("Query on %d entities with %d components took %.3f ms\n",
		N_ENTITIES,
		N_COMPONENTS,
		func_profiler.GetAvg(query_profiler_id))
	fmt.Printf("Matching: %d\n", len(matching))
}
