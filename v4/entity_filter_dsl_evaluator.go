package sameriver

type EntityFilterDSLEvaluator struct {
	predicates map[string](func(args []string, resolver IdentifierResolver) func(*Entity) bool)
	sorts      map[string](func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) bool)
}

func NewEntityFilterDSLEvaluator(
	predicates map[string](func(args []string, resolver IdentifierResolver) func(*Entity) bool),
	sorts map[string](func(args []string, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) bool),
) *EntityFilterDSLEvaluator {
	return &EntityFilterDSLEvaluator{
		predicates: predicates,
		sorts:      sorts,
	}
}

func (e *EntityFilterDSLEvaluator) Evaluate(n *Node, resolver IdentifierResolver) (filter func(*Entity) bool, sort func(xs []*Entity) func(i, j int) bool) {
	if n.Type != NodeExpr {
		panic("Node type must be NodeExpr")
	}

	predicateNode := n.Children[0]
	filter = e.evaluatePredicate(predicateNode, resolver)

	if len(n.Children) > 1 {
		sortNode := n.Children[1]
		sort = e.evaluateSort(sortNode, resolver)
	}

	return filter, sort
}

func (e *EntityFilterDSLEvaluator) evaluatePredicate(n *Node, resolver IdentifierResolver) func(*Entity) bool {
	if n.Type == NodePredicateExpr {
		predicates := make([]func(*Entity) bool, 0, len(n.Children))
		for _, child := range n.Children {
			predicates = append(predicates, e.evaluatePredicate(child, resolver))
		}
		return func(entity *Entity) bool {
			for _, pred := range predicates {
				if !pred(entity) {
					return false
				}
			}
			return true
		}
	} else if n.Type == NodeNot {
		predicate := e.evaluatePredicate(n.Children[0], resolver)
		return func(entity *Entity) bool {
			return !predicate(entity)
		}
	} else if n.Type == NodeFunction {
		args := make([]string, 0, len(n.Children))
		for _, child := range n.Children {
			args = append(args, child.Value)
		}
		return e.predicates[n.Value](args, resolver)
	} else if n.Type == NodeAnd || n.Type == NodeOr {
		left := e.evaluatePredicate(n.Children[0], resolver)
		right := e.evaluatePredicate(n.Children[1], resolver)
		// TODO: how does this work for P && Q && R ?
		// or P && Q || R for that matter?
		if n.Type == NodeAnd {
			return func(entity *Entity) bool {
				return left(entity) && right(entity)
			}
		} else {
			return func(entity *Entity) bool {
				return left(entity) || right(entity)
			}
		}
	}
	panic("Invalid node type for predicate")
}

func (e *EntityFilterDSLEvaluator) evaluateSort(n *Node, resolver IdentifierResolver) func(xs []*Entity) func(i, j int) bool {
	if n.Type != NodeSortExpr {
		panic("Node type must be NodeSortExpr")
	}

	functionNode := n.Children[0]
	args := make([]string, 0, len(functionNode.Children))
	for _, child := range functionNode.Children {
		args = append(args, child.Value)
	}
	return e.sorts[functionNode.Value](args, resolver)
}
