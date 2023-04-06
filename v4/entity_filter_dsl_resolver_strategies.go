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

	entityAccess := func(entity *Entity, bracket string, accessor string) any {
		switch bracket {
		case "[":
			ct := entity.World.em.components
			componentID, ok := ct.stringsRev[accessor]
			if !ok {
				logDSLError("Component %s doesn't exist for DSL expression \"%s\"", accessor, identifier)
			}
			return entity.GetVal(componentID)
		case "<":
			key := accessor
			state := entity.GetIntMap(STATE)
			if !state.Has(key) {
				logDSLError("Entity %s doesn't have state key %s to resolve DSL expression \"%s\"", entity, accessor, identifier)
			}
			return entity.GetIntMap(STATE).Get(key)
		}
		return nil
	}

	valueOrEntityAccess := func(value any) any {
		bracket := ""
		switch {
		case strings.Contains(identifier, "["):
			bracket = "["
		case strings.Contains(identifier, "<"):
			bracket = "<"
		default:
			// else if no access notation, this is just a value, return early
			return value
		}
		// parse the actual access notation by splitting at the bracket
		split := strings.Split(identifier, bracket)
		object, accessor := split[0], split[1]
		// chop the hanging close bracket off the end >:3
		accessor = accessor[:len(accessor)-1]
		// the type of the value any we got as param must be an entity if we're
		// using entity access notation, eg. self for self[position]
		// must be an entity, or mind.bestFriend for mind.bestFriend[mood]
		// must be an entity
		entity, entityOk := value.(*Entity)
		if !entityOk {
			logDSLError("for expression %s, what appears to be entity access notation did not have an entity as its object (%s is not an entity)", identifier, object)
		}
		return entityAccess(entity, bracket, accessor)
	}

	switch parts[0] {
	case "self":
		if len(parts) > 1 {
			return valueOrEntityAccess(er.e)
		}
		return er.e
	case "mind":
		if len(parts) > 1 {
			key := parts[1]
			return valueOrEntityAccess(er.e.GetMind(key))
		}
	case "bb":
		if len(parts) > 1 {
			bbParts := strings.SplitN(parts[1], ".", 2)
			if len(bbParts) > 1 {
				bbname := bbParts[0]
				key := bbParts[1]
				return valueOrEntityAccess(er.e.World.Blackboard(bbname).Get(key))
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
