package sameriver

import "fmt"

// Decorators are used to modify or control the execution of their child nodes.
type BTDecorator struct {
	Name string
	Impl func(self *BTNode) bool
}

// BTNode: A struct that represents a node in the Behavior Tree, either an action
// or a composite node.
type BTNode struct {
	Tree *BehaviourTree

	Parent   *BTNode
	Children []*BTNode

	State map[string]any

	Init func(self *BTNode)

	Failed   bool
	IsFailed func(self *BTNode) bool

	Name string
	// decorators are just strings that index a BTDecorator in the BTRunner's map.
	// either plains strings or... TODO: a DSL?
	Decorators []string
	// the runs counter from the behaviour tree of in which run we last ran
	// all the decorators to success (prevent recalculation when a parent node
	// wanted to run the decorator first before running)
	decoratorsSucceededInRun int
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
	Complete            bool
	CompletionPredicate func(self *BTNode) bool
	WhenDone            func(self *BTNode)
	WhenChildDone       func(self *BTNode)
}

func (n *BTNode) SetChildren(children []*BTNode) {
	n.Children = children
	for _, ch := range children {
		ch.Parent = n
		if ch.Init != nil {
			ch.State = make(map[string]any)
			ch.Init(ch)
		}
	}
}

// when the currently-active action finishes, inform the tree
func (n *BTNode) Done() {
	// set flag
	n.Complete = true
	// callback
	if n.WhenDone != nil {
		n.WhenDone(n)
	}
	// percolate up
	p := n.Parent
	if p != nil && p.WhenChildDone != nil {
		p.WhenChildDone(p)
	}
	// percolate up
	if p != nil {
		p.CompletedChildren++
		// skip over those with nil CompletionPredicate up to the next
		// parent that has a completion predicate
		if p.CompletionPredicate != nil && p.CompletionPredicate(p) {
			p.Done()
		}
	}
}

func (n *BTNode) SetFailed() {
	n.Failed = true
	n.Tree.FailedNodeSet[n] = true
	parent := n.Parent
	// percolate failure up
	for parent != nil {
		if parent.IsFailed(parent) {
			parent.Failed = true
			n.Tree.FailedNodeSet[parent] = true
			parent = parent.Parent
			continue
		} else {
			break
		}
	}
}

// at any time a BT has a pathway down to an action
type BTExecState struct {
	Path   string
	Action *BTNode
}

type BehaviourTree struct {
	// how many times we've run this bad boy
	run int
	// named with a string so we can modularly reuse subtrees by string name
	Name          string
	Root          *BTNode
	FailedNodeSet map[*BTNode]bool
	// current state is the path that's active down to its lowest node, an action
	state *BTExecState
}

func NewBehaviourTree(name string, root *BTNode) *BehaviourTree {
	bt := &BehaviourTree{
		Name:          name,
		Root:          root,
		FailedNodeSet: make(map[*BTNode]bool),
	}
	// recursively iterate the tree and do 3 things:
	// call node.Init(),
	// set the Parent reference
	// set the Tree reference
	var refChildren func(node *BTNode)
	refChildren = func(node *BTNode) {
		if node.Init != nil {
			node.State = make(map[string]any)
			node.Init(node)
		}
		node.Tree = bt
		for _, ch := range node.Children {
			ch.Parent = node
			refChildren(ch)
		}
	}
	refChildren(bt.Root)
	return bt
}

func (bt *BehaviourTree) ResetFailed() {
	for node := range bt.FailedNodeSet {
		node.Failed = false
		delete(bt.FailedNodeSet, node)
	}
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
	decorators map[string]func(*BTNode) bool
}

func NewBTRunner() *BTRunner {
	return &BTRunner{
		trees:      make(map[string]*BehaviourTree),
		decorators: make(map[string]func(self *BTNode) bool),
	}
}

func (btr *BTRunner) RegisterDecorators(decorators []BTDecorator) {
	for _, d := range decorators {
		btr.decorators[d.Name] = d.Impl
	}
}

func (btr *BTRunner) RunDecorators(node *BTNode) bool {
	for _, dstr := range node.Decorators {
		if dec, ok := btr.decorators[dstr]; ok {
			// run the decorator (it can transform the node in any way,
			// add children etc., write to blackboards, etc.)
			// and if it returns false, it failed. Mark this node as
			// failed and
			if !dec(node) {
				node.SetFailed()
				return false
			}
		} else {
			panic(fmt.Sprintf("Unknown decorator: %s", dstr))
		}
	}
	node.decoratorsSucceededInRun = node.Tree.run
	return true
}

func (btr *BTRunner) ExecuteBT(e *Entity, bt *BehaviourTree) *BTExecState {
	bt.run++
	bt.ResetFailed()
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
	// go til we reach the bottom
	for node != nil {
		// every node we visit, is on the path
		// how beautiful
		dotPath(node.Name)

		if node.Complete {
			// when we reach something that's done, our path ends in a dot.
			// you should detect this and know your tree ran out of things to do,
			// it ended up in a state where it doesn't have some other path,
			// the thing it wants to do is actually already done.
			state.Path += "."
			state.Action = node
			return state
		}
		// this is a parasitic bit of code that is not directly related to
		// execution. We just happent to be traversing the tree so, it's a good place
		// to do it. Tell the node who it lives in, if it doesn't know yet
		// (it might have been spawned in as a child of a node by a decorator,
		// eg. a GOAP planner)
		if node.Tree == nil {
			// if ya don't know, now ya know
			node.Tree = bt
		}
		// run decorators (unless we ran them already this run, by a parent probing)
		if len(node.Decorators) > 0 && node.decoratorsSucceededInRun != bt.run {
			if !btr.RunDecorators(node) {
				// we failed.
				return nil
			}
		}
		if node.Selector != nil {
			childIndex := node.Selector(node)
			if childIndex >= 0 && childIndex < len(node.Children) {
				child := node.Children[childIndex]
				node = child
				continue
			} else {
				return nil
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
