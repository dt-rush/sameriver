package sameriver

import (
	"math/rand"
	"testing"
)

func TestBTConstruction(t *testing.T) {
	w := testingWorld()

	// Define your decorator functions and add them to BTRunner
	decorators := map[string]BTDecorator{
		"planPlant": func(node *BTNode) bool {
			// mock GOAP
			node.Children = []*BTNode{
				{Name: "getHandplow"},
				{Name: "goToField"},
				{Name: "doHandplow"},
			}
			return true
		},
	}

	btr := &BTRunner{
		trees:      make(map[string]*BehaviourTree),
		decorators: decorators,
	}

	// Create and add villagerRoot tree
	villagerRoot := &BehaviourTree{
		Name: "villagerRoot",
		Root: &BTNode{
			Name: "Utility",
			Selector: func(self *BTNode) int {
				// Implement your logic for selecting the Utility node
				return 3 // Select the "plant" node for testing purposes
			},
			IsFailed: func(self *BTNode) bool {
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
	}
	btr.trees["villagerRoot"] = villagerRoot

	// Create and add plant tree
	plant := &BehaviourTree{
		Name: "plant",
		Root: &BTNode{
			Name:       "Sequence",
			Decorators: []string{"planPlant"},
			Selector: func(self *BTNode) int {
				return 1
			},
			IsFailed: func(self *BTNode) bool {
				for _, child := range self.Children {
					if child.Failed {
						return true
					}
				}
				return false
			},
			Children: nil,
		},
	}
	btr.trees["plant"] = plant

	// Execute the behavior tree and check the result
	e := w.Spawn(nil)

	result := btr.ExecuteBT(e, villagerRoot)
	expectedPath := "Utility.plant.Sequence.goToField"

	Logger.Printf("BT descent path: %s", result.Path)

	if result.Path != expectedPath {
		t.Errorf("Expected result: %s, got: %s", expectedPath, result.Path)
	}
}

func TestBTAnyNodeFailure(t *testing.T) {
	w := testingWorld()

	// Define your decorator functions and add them to BTRunner
	decorators := map[string]BTDecorator{
		"fail": func(node *BTNode) bool {
			return false
		},
		"pass": func(node *BTNode) bool {
			return true
		},
	}

	btr := &BTRunner{
		trees:      make(map[string]*BehaviourTree),
		decorators: decorators,
	}

	// Create and add anyRoot tree
	anyRoot := &BehaviourTree{
		Name: "anyRoot",
		Root: &BTNode{
			Name: "Any",
			Selector: func(self *BTNode) int {
				// Implement your logic for selecting the Any node
				perm := rand.Perm(len(self.Children))
				for _, i := range perm {
					child := self.Children[i]
					decoratorsPassed := true
					for _, decorator := range child.Decorators {
						decoratorFunc, ok := btr.decorators[decorator]
						if !ok || !decoratorFunc(child) {
							decoratorsPassed = false
							break
						}
					}
					if decoratorsPassed && !child.Failed {
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
	}
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

	// Define your decorator functions and add them to BTRunner
	decorators := map[string]BTDecorator{
		"fail": func(node *BTNode) bool {
			return false
		},
		"pass": func(node *BTNode) bool {
			return true
		},
	}

	btr := &BTRunner{
		trees:      make(map[string]*BehaviourTree),
		decorators: decorators,
	}

	// Create and add orderedAnyRoot tree
	orderedAnyRoot := &BehaviourTree{
		Name: "orderedAnyRoot",
		Root: &BTNode{
			Name: "OrderedAny",
			Selector: func(self *BTNode) int {
				// Implement your logic for selecting the OrderedAny node
				for i, child := range self.Children {
					decoratorsPassed := true
					for _, decorator := range child.Decorators {
						decoratorFunc, ok := btr.decorators[decorator]
						if !ok || !decoratorFunc(child) {
							decoratorsPassed = false
							break
						}
					}
					if decoratorsPassed && !child.Failed {
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
	}
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
