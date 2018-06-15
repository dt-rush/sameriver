package main

import (
	"fmt"
)

// X, Y: 		position in grid
// From:		link to the next node in path
// G:	 		path cost
// H:	 		heuristic
// F:			G + H
type Node struct {
	X, Y   int
	From   *Node
	G      int
	H      int
	F      int
	HeapIX int
}

func (n *Node) String() string {
	return fmt.Sprintf("[%d]: %p: (%d, %d), from: %p, G: %d, H: %d, F: %d",
		n.HeapIX, n, n.X, n.Y, n.From, n.G, n.H, n.F)
}

// G:			reference to the grid of cells we're pathing over
// OH			NodeHeap ("Open Heap") used to pop off the nodes with the lowest
//					F during search (open nodes)
// O:			2d array shadowing Nodes used to keep track of whether the given
//					cell is on the open heap or not
// C:			2d Array shadowing Nodes used to keep track of the "closed list",
//					for the current calculation (the list of cells we have visited
//					and checked the neighbors of, adding them to the open heap).
//					We use the value of N to avoid having to clear this.
//					(Theoretically you could overflow if a cell is not visited
//					for 2^32-1 passes, lol)
// N: 			incremented each time we calculate (used to avoid having to
//					clear values in C)
type PathComputer struct {
	G  *Grid
	OH *NodeHeap
	O  [][]NodeOnOpenHeap
	C  [][]int
	N  int
}

// used to support O above
type NodeOnOpenHeap struct {
	Node *Node
	N    int
}

func NewPathComputer(g *Grid) *PathComputer {
	// make 2d arrays
	nodes := make([][]Node, g.W)
	o := make([][]NodeOnOpenHeap, g.W)
	c := make([][]int, g.W)
	for x := 0; x < g.W; x++ {
		nodes[x] = make([]Node, g.H)
		o[x] = make([]NodeOnOpenHeap, g.H)
		c[x] = make([]int, g.H)
	}
	// make node heap
	h := NewNodeHeap()
	return &PathComputer{
		G:  g,
		OH: h,
		O:  o,
		C:  c,
		N:  0,
	}
}

// neighbor x, y offsets
//
//                                       X
//      --------------------------------->
//     |
//     |    -1,  1     0,  1     1,  1
//     |
//     |    -1,  0               1,  0
//     |
//     |    -1, -1     0, -1     1, -1
//     |
//  Y  v
//
//
var neighborIXs = [][2]int{
	[2]int{-1, 1},
	[2]int{0, 1},
	[2]int{1, 1},
	[2]int{-1, 0},
	[2]int{1, 0},
	[2]int{-1, -1},
	[2]int{0, -1},
	[2]int{1, -1},
}

// manhattan heuristic * 100 (since we use 10 and 14 for one square straight
// or diagonal)
func (pc *PathComputer) Heuristic(p1 Position, p2 Position) int {
	dx := p1.X - p2.X
	if dx < 0 {
		dx *= -1
	}
	dy := p1.Y - p2.Y
	if dy < 0 {
		dy *= -1
	}
	return 100 * (dx + dy)
}

func (pc *PathComputer) Path(start Position, end Position) (path []Position) {
	// clear the heap which contains leftover nodes from the last calculation
	pc.OH.Clear()
	// increment N (easier than clearing arrays)
	pc.N++
	// add first node to open heap
	firstNode := Node{
		X:    start.X,
		Y:    start.Y,
		From: nil,
		G:    0,
		H:    0,
	}
	pc.OH.Add(&firstNode)
	fmt.Println(pc.OH.String())
	// while open heap has elements...
	for pc.OH.Len() > 0 {
		// pop from open heap and set as closed
		cur, err := pc.OH.Pop()
		fmt.Printf("Popped (%d, %d) with F = %d\n", cur.X, cur.Y, cur.F)
		fmt.Println(pc.OH.String())
		pc.C[cur.X][cur.Y] = pc.N
		// if err, we have exhausted all squares on open heap and found no path
		// return nil
		if err != nil {
			return nil
		}
		// if the current cell is the end, we're here. build the return list
		if cur.X == end.X && cur.Y == end.Y {
			path = make([]Position, 0)
			for cur != nil {
				path = append(path, Position{cur.X, cur.Y})
				cur = cur.From
			}
			// return the path to the user
			return path
		}
		// else, we have yet to complete the path. So:
		// for each neighbor
		for _, neighborIX := range neighborIXs {
			// get the coordinates of the cell we will check the cost to
			// by applying an offset to cur's coordinates
			x := cur.X + neighborIX[0]
			y := cur.Y + neighborIX[1]
			// continue loop to next neighbor early if not in grid
			inGrid := x >= 0 && x < pc.G.W && y >= 0 && y < pc.G.H
			if !inGrid {
				continue
			}
			// if neighbor is valid (there is no obstacle and it isn't in the
			// closed set), then...
			obstacle := pc.G.Cells[x][y] == OBSTACLE
			closed := pc.C[x][y] == pc.N
			if !obstacle && !closed {
				fmt.Printf("looking at (%d, %d)\n", x, y)
				// dist is an integer expression of the distance from
				// cur to the neighbor cell we're looking at here.
				// if either x or y offset is 0, we're moving straight,
				// so put 10. Otherwise we're moving diagonal, so put 14
				// (these are 1 and sqrt(2), but made into integers for speed)
				var dist int
				if neighborIX[0]*neighborIX[1] == 0 {
					dist = 10
				} else {
					dist = 14
				}
				// compute g, h, f for the current cell
				g := cur.G + dist
				h := pc.Heuristic(Position{x, y}, end)
				// if not on open heap, add it with "From" == cur
				if pc.O[x][y].N != pc.N {
					neighbor := Node{
						X:    x,
						Y:    y,
						From: cur,
						G:    g,
						H:    h}
					pc.O[x][y] = NodeOnOpenHeap{&neighbor, pc.N}
					pc.OH.Add(&neighbor)
				} else {
					// if it *is* on the open heap already, check to see if
					// this is a better path to that square
					// on -> "open node"
					on := pc.O[x][y].Node
					if g < on.G {
						// if the open node could be reached better by
						// this path, set the g to the new lower g, set the
						// "From" reference to cur and fix up the heap because
						// we've changed the value of one of its elements
						on.From = cur
						pc.OH.Modify(on.HeapIX, g)
					}
				}
			}
		}
	}
	return path
}
