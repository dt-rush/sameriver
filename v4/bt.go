package sameriver

import "fmt"

type BTCompletionPredicate func(*BTNode) bool

type BTNode struct {
	Parent   *BTNode
	Children []*BTNode

	Failed   bool
	IsFailed func(self *BTNode) bool

	Name string
	// decorators are just strings - either plains strings or... TODO: a DSL?
	Decorators []string
	// if non-nil, this is a composite node; the selector will tell which
	// child to select for running.
	// (for Sequence, this is { return ChildrenCompleted })
	Selector func(self *BTNode) int
	// whenever Done() is called on an action its parent's ChildrenCompleted
	// counter is incremented. When this counter in the parent increments,
	// if it causes its CompletionPredicate to return true, then it too will
	// consider itself done and increment the counter in its parent, then
	// checking if its done to continue to percolate up.
	CompletedChildren int
	// for a node that never ends, this is always false
	// for a Sequence node, this is n.ChildrenCompleted == len(n.Children)
	// for an Any node, this is n.ChildrenCompleted > 0
	// for an All node, this is n.ChildrenCompleted == len(n.Children)
	CompletionPredicate func(self *BTNode) bool
}

type BTDecorator func(*BTNode) bool

// named with a string so we can modularly reuse subtrees by string name
type BehaviourTree struct {
	Name string
	Root *BTNode
	// current state is the path that's active down to its lowest node, an action
	state *BTExecState
}

// at any time a BT has a pathway down to an action
type BTExecState struct {
	Path   string
	Action *BTNode
}

func NewBehaviourTree(name string, root *BTNode) *BehaviourTree {
	bt := &BehaviourTree{
		Name: name,
		Root: root,
	}
	// set the parent of all nodes by descent
	var setChildren func(node *BTNode)
	setChildren = func(node *BTNode) {
		for _, ch := range node.Children {
			setChildren(ch)
			ch.Parent = node
		}
	}
	setChildren(bt.Root)
	return bt
}

// stores the database of named trees and decorators needed to run a tree
type BTRunner struct {
	// if we reach a string node with no children, it is potentially just
	// an atomic action (some FSM will run it, it's a single animation, etc.)
	// OR
	// it can be a reference to a named tree in the system, so this way
	// we can get simple reusability of subtrees
	trees map[string]*BehaviourTree

	// the runner has a set of decorators that it can honour - a decorator is
	// just a string for now, but really, it should be a parsed DSL
	decorators map[string]BTDecorator
}

func (btr *BTRunner) ExecuteBT(e *Entity, bt *BehaviourTree) *BTExecState {
	state := &BTExecState{}
	if bt.Root == nil {
		return state
	}
	node := bt.Root
	dotPath := func(s string) {
		if state.Path == "" {
			state.Path = s
		} else {
			state.Path += "." + s
		}
	}
	for node != nil {
		if len(node.Decorators) > 0 {
			for _, dstr := range node.Decorators {
				if dec, ok := btr.decorators[dstr]; ok {
					// run the decorator (it can transform the node in any way,
					// add children etc., write to blackboards, etc.)
					// and if it returns false, it failed. Mark this node as
					// failed and
					if !dec(node) {
						node.Failed = true
						parent := node.Parent
						// percolate failure up
						for parent != nil {
							if parent.IsFailed(parent) {
								parent.Failed = true
								parent = parent.Parent
								continue
							} else {
								break
							}
						}
						// re-execute
						return btr.ExecuteBT(e, bt)
					}
				} else {
					panic(fmt.Sprintf("Unknown decorator: %s", dstr))
				}
			}
		}
		if node.Selector != nil {
			dotPath(node.Name)
			childIndex := node.Selector(node)
			if childIndex >= 0 && childIndex < len(node.Children) {
				child := node.Children[childIndex]
				dotPath(child.Name)
				node = child
				continue
			}
		} else {
			if tree, ok := btr.trees[node.Name]; ok {
				node = tree.Root
				continue
			} else {
				state.Action = node
				break
			}
		}
	}
	bt.state = state
	return bt.state
}

func (bt *BehaviourTree) Done() {
	if bt.state == nil {
		return
	}

	node := bt.state.Action.Parent
	for node != nil {
		node.CompletedChildren++
		if node.CompletionPredicate != nil {
			if node.CompletionPredicate(node) {
				node = node.Parent
			} else {
				break
			}
		} else {
			node = node.Parent
		}
	}
}
