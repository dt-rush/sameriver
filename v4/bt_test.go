package sameriver

import (
	"math/rand"
	"testing"
)

func TestBTConstruction(t *testing.T) {
	w := testingWorld()

	btr := NewBTRunner()

	btr.RegisterDecorators([]BTDecorator{
		BTDecorator{
			Name: "planPlant",
			Impl: func(node *BTNode) bool {
				// mock GOAP
				node.SetChildren([]*BTNode{
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
				return false
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
				return false
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
