package sameriver

/*
This file implements the identifier resolver functionality for the
Entity Filter DSL. The parser recognizes expressions that are structured
like a series of function calls, but it does not differentiate between the
different types of notation that can be used to reference different types of
values. Instead, it just identifies a series of functions and their arguments,
like "F(x, y) && G(z); H(q)". The identifier resolver is responsible for
interpreting these expressions (x, y, z, q) and determining the appropriate value
to pass to the implementation predicate/sort funcs.

To accomplish this, the identifier resolver supports several different types
of notation.

"self" refers to the current entity

"mind.field" and "bb.village4.huntingParty" allow the user to look up values
in the entity's mind or a named blackboard, respectively.

Entity-access notation, like

"bb.village4.fisherman[position]" (component access)

or

"mind.friend<mood>" (state access),

... allows the user to access specific components or state values associated with
the entity.

This file contains two resolver types: EntityResolver and
WorldResolver, which implement the IdentifierResolver interface and provide
functionality for resolving identifiers in the context of an entity or the
entire world, respectively.
*/

import "strings"

type IdentifierResolver interface {
	Resolve(identifier string) any
}

type EntityResolver struct {
	e *Entity
}

type WorldResolver struct {
	w *World
}

func valueOrEntityAccess(value any, identifier string) any {
	bracket := ""
	switch {
	case strings.ContainsRune(identifier, '['):
		bracket = "["
	case strings.ContainsRune(identifier, '<'):
		bracket = "<"
	default:
		return value
	}

	split := strings.Split(identifier, bracket)
	object, accessor := split[0], split[1]
	accessor = accessor[:len(accessor)-1]

	entity, entityOk := value.(*Entity)
	if !entityOk {
		logDSLError("for expression %s, what appears to be entity access notation did not have an entity as its object (%s is not an entity)", identifier, object)
		return nil
	}

	switch bracket {
	case "[":
		ct := entity.World.em.components
		componentID, ok := ct.stringsRev[accessor]
		if !ok {
			logDSLError("Component %s doesn't exist for DSL expression \"%s\"", accessor, identifier)
			return nil
		}
		return entity.GetVal(componentID)
	case "<":
		key := accessor
		state := entity.GetIntMap(STATE)
		if !state.Has(key) {
			logDSLError("Entity %s doesn't have state key %s to resolve DSL expression \"%s\"", entity, accessor, identifier)
			return nil
		}
		return entity.GetIntMap(STATE).Get(key)
	}

	return nil
}

func (er *EntityResolver) Resolve(identifier string) any {
	parts := strings.SplitN(identifier, ".", 2)

	switch parts[0] {
	case "x":
		if len(parts) > 1 {
			// what do we do here? how do we do valueOrEntityAccess on x? we don't have it
		}
		// TODO: what do we return here? we don't have access to x?
	case "self":
		if len(parts) > 1 {
			return valueOrEntityAccess(er.e, identifier)
		}
		return er.e
	case "mind":
		if len(parts) > 1 {
			key := parts[1]
			return valueOrEntityAccess(er.e.GetMind(key), identifier)
		}
	case "bb":
		if len(parts) > 1 {
			bbParts := strings.SplitN(parts[1], ".", 2)
			if len(bbParts) > 1 {
				bbname := bbParts[0]
				key := bbParts[1]
				return valueOrEntityAccess(er.e.World.Blackboard(bbname).Get(key), identifier)
			}
		}
	}

	return nil
}

func (wr *WorldResolver) Resolve(identifier string) any {
	parts := strings.SplitN(identifier, ".", 2)

	if parts[0] == "bb" {
		if len(parts) > 1 {
			bbParts := strings.SplitN(parts[1], ".", 2)
			if len(bbParts) > 1 {
				bbname := bbParts[0]
				key := bbParts[1]
				return valueOrEntityAccess(wr.w.Blackboard(bbname).Get(key), identifier)
			}
		}
	}

	return nil
}
