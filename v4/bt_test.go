package sameriver

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBTSimple(t *testing.T) {
	w := testingWorld()

	btr := NewBTRunner()

	btr.RegisterDecorators([]BTDecorator{
		BTDecorator{
			Name: "planPlant",
			Impl: func(self *BTNode) bool {
				// mock GOAP
				self.SetChildren([]*BTNode{
					{Name: "getHandplow"},
					{Name: "goToField"},
					{Name: "doHandplow"},
				})
				return true
			},
		},
	})

	// Create and add villagerRoot tree
	villagerRoot := NewBehaviourTree(
		"villagerRoot",
		&BTNode{
			Name: "Utility",
			Selector: func(self *BTNode) int {
				// Implement your logic for selecting the Utility node
				return 3 // Select the "plant" node for testing purposes
			},
			IsFailed: func(self *BTNode) bool {
				return false
			},
			CompletionPredicate: func(self *BTNode) bool {
				return false
			},
			Children: []*BTNode{
				{Name: "fight-flight"},
				{Name: "rest"},
				{Name: "eat"},
				{Name: "plant"},
				{Name: "harvest"},
				{Name: "craft"},
				{Name: "leisure"},
				{Name: "religion"},
			},
		},
	)
	btr.trees["villagerRoot"] = villagerRoot

	// Create and add plant tree
	plant := NewBehaviourTree(
		"plant",
		&BTNode{
			Name:       "Sequence",
			Decorators: []string{"planPlant"},
			Selector: func(self *BTNode) int {
				return self.CompletedChildren
			},
			IsFailed: func(self *BTNode) bool {
				for _, child := range self.Children {
					if child.Failed {
						return true
					}
				}
				return false
			},
			CompletionPredicate: func(self *BTNode) bool {
				return self.CompletedChildren == len(self.Children)
			},
			Children: nil,
		},
	)
	btr.trees["plant"] = plant

	// Execute the behavior tree and check the result
	e := w.Spawn(nil)

	result := btr.ExecuteBT(e, villagerRoot)
	expectedPath := "Utility.plant.Sequence.getHandplow"

	Logger.Printf("BT descent path: %s", result.Path)

	if result.Path != expectedPath {
		t.Errorf("Expected result: %s, got: %s", expectedPath, result.Path)
	}

	for i := 0; i < 5; i++ {
		if result != nil {
			result.Action.Done()
		}
		result = btr.ExecuteBT(e, villagerRoot)
		if result != nil {
			Logger.Printf("BT descent path: %s", result.Path)
		} else {
			Logger.Printf("Nil")
		}
	}

}

func TestBTAnyNodeFailure(t *testing.T) {
	w := testingWorld()

	btr := NewBTRunner()

	btr.RegisterDecorators([]BTDecorator{
		BTDecorator{
			Name: "fail",
			Impl: func(self *BTNode) bool {
				return false
			},
		},
		BTDecorator{
			Name: "pass",
			Impl: func(self *BTNode) bool {
				return true
			},
		},
	})

	// Create and add anyRoot tree
	anyRoot := NewBehaviourTree(
		"anyRoot",
		&BTNode{
			Name: "Any",
			Selector: func(self *BTNode) int {
				// Implement your logic for selecting the Any node
				perm := rand.Perm(len(self.Children))
				for _, i := range perm {
					child := self.Children[i]
					if btr.RunDecorators(child) {
						return i
					}
				}
				return -1
			},
			IsFailed: func(self *BTNode) bool {
				// failed if they all failed
				failed := true
				for _, ch := range self.Children {
					failed = failed && ch.Failed
				}
				return failed
			},
			CompletionPredicate: func(self *BTNode) bool {
				// complete if one has run
				return self.CompletedChildren > 0
			},
			Children: []*BTNode{
				{Name: "fail1", Decorators: []string{"fail"}},
				{Name: "fail2", Decorators: []string{"fail"}},
				{Name: "fail3", Decorators: []string{"fail"}},
				{Name: "fail4", Decorators: []string{"fail"}},
				{Name: "fail5", Decorators: []string{"fail"}},
				{Name: "fail6", Decorators: []string{"fail"}},
				{Name: "fail7", Decorators: []string{"fail"}},
				{Name: "fail8", Decorators: []string{"fail"}},
				{Name: "fail9", Decorators: []string{"fail"}},
				{Name: "success", Decorators: []string{"pass"}},
			},
		},
	)
	btr.trees["anyRoot"] = anyRoot

	// Execute the behavior tree and check the result
	e := w.Spawn(nil)

	result := btr.ExecuteBT(e, anyRoot)
	expectedPath := "Any.success"

	Logger.Printf("BT descent path: %s", result.Path)

	if result.Path != expectedPath {
		t.Errorf("Expected result: %s, got: %s", expectedPath, result.Path)
	}
}

func TestBTOrderedAnyNodeFailure(t *testing.T) {
	w := testingWorld()

	btr := NewBTRunner()

	btr.RegisterDecorators([]BTDecorator{
		BTDecorator{
			Name: "fail",
			Impl: func(self *BTNode) bool {
				return false
			},
		},
		BTDecorator{
			Name: "pass",
			Impl: func(self *BTNode) bool {
				return true
			},
		},
	})

	// Create and add orderedAnyRoot tree
	orderedAnyRoot := NewBehaviourTree(
		"orderedAnyRoot",
		&BTNode{
			// OrderedAny runs the first of its children whose decorators pass
			Name: "OrderedAny",
			Selector: func(self *BTNode) int {
				// Implement your logic for selecting the OrderedAny node
				for i, child := range self.Children {
					if btr.RunDecorators(child) {
						return i
					}
				}
				return -1
			},
			IsFailed: func(self *BTNode) bool {
				// failed if they all failed
				failed := true
				for _, ch := range self.Children {
					failed = failed && ch.Failed
				}
				return failed
			},
			CompletionPredicate: func(self *BTNode) bool {
				// complete if one has run
				return self.CompletedChildren > 0
			},
			Children: []*BTNode{
				{Name: "fail1", Decorators: []string{"fail"}},
				{Name: "fail2", Decorators: []string{"fail"}},
				{Name: "fail3", Decorators: []string{"fail"}},
				{Name: "successA", Decorators: []string{"pass"}},
				{Name: "fail5", Decorators: []string{"fail"}},
				{Name: "fail6", Decorators: []string{"fail"}},
				{Name: "fail7", Decorators: []string{"fail"}},
				{Name: "fail8", Decorators: []string{"fail"}},
				{Name: "fail9", Decorators: []string{"fail"}},
				{Name: "successB", Decorators: []string{"pass"}},
			},
		},
	)
	btr.trees["orderedAnyRoot"] = orderedAnyRoot

	// Execute the behavior tree and check the result
	e := w.Spawn(nil)

	result := btr.ExecuteBT(e, orderedAnyRoot)
	expectedPath := "OrderedAny.successA"

	Logger.Printf("BT descent path: %s", result.Path)

	if result.Path != expectedPath {
		t.Errorf("Expected result: %s, got: %s", expectedPath, result.Path)
	}
}

func TestBTAllNode(t *testing.T) {
	w := testingWorld()

	btr := NewBTRunner()

	btr.RegisterDecorators([]BTDecorator{
		BTDecorator{
			Name: "fail",
			Impl: func(self *BTNode) bool {
				return false
			},
		},
		BTDecorator{
			Name: "pass",
			Impl: func(self *BTNode) bool {
				return true
			},
		},
	})

	// Create and add allRoot tree
	allRoot := NewBehaviourTree(
		"allRoot",
		&BTNode{
			// All composite node runs all of its children in a random order
			Name: "All",
			Init: func(self *BTNode) {
				self.State["perm"] = rand.Perm(len(self.Children))
			},
			Selector: func(self *BTNode) int {
				perm := self.State["perm"].([]int)
				for _, idx := range perm {
					child := self.Children[idx]
					if !child.Complete {
						return idx
					}
				}
				return -1
			},
			IsFailed: func(self *BTNode) bool {
				for _, child := range self.Children {
					if child.Failed {
						return true
					}
				}
				return false
			},
			CompletionPredicate: func(self *BTNode) bool {
				for _, child := range self.Children {
					if !child.Complete {
						return false
					}
				}
				self.Init(self)
				return true
			},
			Children: []*BTNode{
				{Name: "successA", Decorators: []string{"pass"}},
				{Name: "successB", Decorators: []string{"pass"}},
				{Name: "successC", Decorators: []string{"pass"}},
				{Name: "successD", Decorators: []string{"pass"}},
			},
		},
	)
	btr.trees["allRoot"] = allRoot

	// the test itself
	e := w.Spawn(nil)

	expectedNodes := []string{"successA", "successB", "successC", "successD"}

	executedNodes := map[string]bool{}
	for i := 0; i < 10; i++ {
		result := btr.ExecuteBT(e, allRoot)
		if result != nil {
			Logger.Println(result.Path)
			result.Action.Done()
			executedNodes[result.Action.Name] = true
		} else {
			Logger.Println("nil")
		}
	}

	// Check if all child nodes were executed
	for _, expectedNode := range expectedNodes {
		if !executedNodes[expectedNode] {
			t.Errorf("Expected %s to be in the path, but it was not.", expectedNode)
		}
	}
}

func TestBTRandomPriorityLoopNode(t *testing.T) {
	w := testingWorld()

	btr := NewBTRunner()

	// Create and add randomRoot tree
	randomRoot := NewBehaviourTree(
		"randomRoot",
		&BTNode{
			Name: "Random",
			Selector: func(self *BTNode) int {
				return rand.Intn(len(self.Children))
			},
			CompletionPredicate: func(self *BTNode) bool {
				return false
			},
			Children: []*BTNode{
				{Name: "actionA"},
				{Name: "actionB"},
				{Name: "actionC"},
			},
		},
	)
	btr.trees["randomRoot"] = randomRoot

	// Create and add priorityRoot tree
	priorityRoot := NewBehaviourTree(
		"priorityRoot",
		&BTNode{
			Name: "Priority",
			Selector: func(self *BTNode) int {
				minPriority := math.MaxFloat64
				selectedIdx := -1
				for idx, child := range self.Children {
					if !child.Complete && !child.Failed {
						if child.State["priority"].(float64) < minPriority {
							minPriority = child.State["priority"].(float64)
							selectedIdx = idx
						}
					}
				}
				return selectedIdx
			},
			CompletionPredicate: func(self *BTNode) bool {
				for _, child := range self.Children {
					if child.Complete {
						return true
					}
				}
				return false
			},
			Children: []*BTNode{
				{Name: "actionD", State: map[string]any{"priority": 1.0}},
				{Name: "actionE", State: map[string]any{"priority": 2.0}},
				{Name: "actionF", State: map[string]any{"priority": 3.0}},
			},
		},
	)
	btr.trees["priorityRoot"] = priorityRoot

	// Create and add loopRoot tree
	loopRoot := NewBehaviourTree(
		"loopRoot",
		&BTNode{
			Name: "Loop",
			State: map[string]any{
				"N": 3,
			},
			Selector: func(self *BTNode) int {
				if self.State["currentIndex"] == nil {
					self.State["currentIndex"] = 0
				}
				return self.State["currentIndex"].(int)
			},
			CompletionPredicate: func(self *BTNode) bool {
				if _, ok := self.State["N"]; !ok {
					// if no N, we loop forever
					return false
				}
				// else, have we done N full sets?
				loops := self.CompletedChildren / len(self.Children)
				return loops >= self.State["N"].(int)
			},
			WhenChildDone: func(self *BTNode) {
				// increment modulo
				self.State["currentIndex"] = (self.State["currentIndex"].(int) + 1) % len(self.Children)
				// when we wrap around, turn all back to Done false
				if self.State["currentIndex"].(int) == 0 {
					for _, ch := range self.Children {
						ch.Complete = false
					}
				}
			},
			Children: []*BTNode{
				{Name: "actionG"},
				{Name: "actionH"},
				{Name: "actionI"},
			},
		},
	)
	btr.trees["loopRoot"] = loopRoot

	// the test itself
	e := w.Spawn(nil)

	// Helper function to run the behavior tree and collect executed nodes
	var executedNodes map[string]bool
	runAndCollectExecutedNodes := func(treeName string, iterations int) {
		for i := 0; i < iterations; i++ {
			result := btr.ExecuteBT(e, btr.trees[treeName])
			if result != nil {
				Logger.Println(result.Path)
				result.Action.Done()
				executedNodes[result.Action.Name] = true
			} else {
				Logger.Println("nil")
			}
		}
	}

	// Run each tree using the helper function

	// run random 10 times cause the odds are crazy
	randomPassed := false
	for i := 0; i < 10; i++ {
		executedNodes = make(map[string]bool)
		runAndCollectExecutedNodes("randomRoot", 10)
		if executedNodes["actionA"] && executedNodes["actionB"] && executedNodes["actionC"] {
			randomPassed = true
			break
		}
	}
	if !randomPassed {
		t.Error("Failed to run all required nodes after 10 tries")
	}

	runAndCollectExecutedNodes("priorityRoot", 1)
	runAndCollectExecutedNodes("loopRoot", 4*len(loopRoot.Root.Children))

	fmt.Printf("Executed nodes: %v", executedNodes)

	expectedNodes := []string{"actionD", "actionG", "actionH", "actionI"}
	for _, node := range expectedNodes {
		assert.Contains(t, executedNodes, node, "expected node %q to be executed", node)
	}
}

/*

TODO:

this is actually a decorator
Cooldown: This composite node adds a cooldown period to its child, preventing the child from being executed until the cooldown has expired. This can help prevent entities from repeatedly performing certain actions too quickly, which could be unrealistic or unbalanced.

Parallel: This composite node executes its children concurrently, allowing for simultaneous actions. For example, an entity could be moving and attacking at the same time. Note that parallel execution may be more complex to implement depending on your game engine.
*/
