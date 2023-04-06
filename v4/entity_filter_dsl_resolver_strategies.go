package sameriver

import "strings"

type IdentifierResolver interface {
	Resolve(identifier string) any
}

type EntityResolver struct {
	e *Entity
}

func (er *EntityResolver) Resolve(identifier string) any {
	parts := strings.SplitN(identifier, ".", 2)

	switch parts[0] {
	case "self":
		return er.e
	case "mind":
		if len(parts) > 1 {
			key := parts[1]
			return er.e.GetMind(key)
		}
	case "bb":
		if len(parts) > 1 {
			bbParts := strings.SplitN(parts[1], ".", 2)
			if len(bbParts) > 1 {
				bbname := bbParts[0]
				key := bbParts[1]
				return er.e.World.Blackboard(bbname).Get(key)
			}
		}
	}

	return nil
}

type WorldResolver struct {
	w *World
}

func (wr *WorldResolver) Resolve(identifier string) any {
	parts := strings.SplitN(identifier, ".", 2)

	if parts[0] == "bb" {
		if len(parts) > 1 {
			bbParts := strings.SplitN(parts[1], ".", 2)
			if len(bbParts) > 1 {
				bbname := bbParts[0]
				key := bbParts[1]
				return wr.w.Blackboard(bbname).Get(key)
			}
		}
	}

	return nil
}
