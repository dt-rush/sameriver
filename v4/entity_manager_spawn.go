package sameriver

import (
	"fmt"
	"strconv"
	"strings"
)

// get the current number of requests in the channel and only process
// them. More may continue to pile up. They'll get processed next time.
func (m *EntityManager) processSpawnChannel() {
	n := len(m.spawnSubscription.C)
	for i := 0; i < n; i++ {
		// get the request from the channel
		e := <-m.spawnSubscription.C
		spec := e.Data.(map[string]any)
		m.Spawn(spec)
	}
}

func (m *EntityManager) Spawn(spec map[string]any) *Entity {

	if spec == nil {
		spec = make(map[string]any)
	}

	var active bool
	var uniqueTag string
	var tags []string
	var componentSpecs map[ComponentID]any
	var customComponentSpecs map[ComponentID]any
	var customComponentsImpl map[ComponentID]CustomContiguousComponent
	var logics map[string](func(e *Entity, dt_ms float64))
	var funcs map[string](func(e *Entity, params any) any)
	var mind map[string]any

	// type assert spec vars

	if _, ok := spec["active"]; ok {
		active = spec["active"].(bool)
	} else {
		active = true
	}

	if _, ok := spec["uniqueTag"]; ok {
		uniqueTag = spec["uniqueTag"].(string)
	} else {
		uniqueTag = ""
	}

	if _, ok := spec["tags"]; ok {
		tags = spec["tags"].([]string)
	} else {
		tags = []string{}
	}

	if _, ok := spec["components"]; ok {
		componentSpecs = spec["components"].(map[ComponentID]any)
	} else {
		componentSpecs = make(map[ComponentID]any)
	}

	if _, ok := spec["customComponents"]; ok {
		customComponentsImpl = spec["customComponentsImpl"].(map[ComponentID]CustomContiguousComponent)
		customComponentSpecs = spec["customComponents"].(map[ComponentID]any)
	}

	if _, ok := spec["logics"]; ok {
		logics = spec["logics"].(map[string](func(e *Entity, dt_ms float64)))
	} else {
		logics = make(map[string](func(e *Entity, dt_ms float64)))
	}

	if _, ok := spec["funcs"]; ok {
		funcs = spec["funcs"].(map[string](func(e *Entity, params any) any))
	} else {
		funcs = make(map[string](func(e *Entity, params any) any))
	}

	if _, ok := spec["mind"]; ok {
		mind = spec["mind"].(map[string]any)
	} else {
		mind = make(map[string]any)
	}

	// add empty generictags (overwriting the spec? boo hoo,
	// shouldn't be specifying tags this way anyway)
	componentSpecs[GENERICTAGS] = NewTagList()

	return m.doSpawn(
		active,
		uniqueTag,
		tags,
		m.components.makeCustomComponentSet(componentSpecs, customComponentSpecs, customComponentsImpl),
		logics,
		funcs,
		mind,
	)
}

func (m *EntityManager) QueueSpawn(spec map[string]any) {
	if len(m.spawnSubscription.C) >= EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
		go func() {
			m.spawnSubscription.C <- Event{"spawn-request", spec}
		}()
	} else {
		m.spawnSubscription.C <- Event{"spawn-request", spec}
	}
}

// given a list of components, spawn an entity with the default values
// returns the Entity (used to spawn an entity for which we *want* the
// token back)

func (m *EntityManager) doSpawn(
	active bool,
	uniqueTag string,
	tags []string,
	components ComponentSet,
	logics map[string](func(e *Entity, dt_ms float64)),
	funcs map[string](func(e *Entity, params any) any),
	mind map[string]any,
) *Entity {

	// get an ID for the entity
	// if maximum entity count reached, resize arrays then allocate id
	if total, _ := m.NumEntities(); total == m.MaxEntities() {
		logWarning("Reached entity capacity: %d"+
			"; expanding component tables, system storage, and id pool", m.MaxEntities())
		m.ExpandEntityTables()
	}
	e := m.entityIDAllocator.allocateID()

	e.World = m.w
	// copy the data into the component storage for each component
	m.components.applyComponentSet(e, components)
	// create (if doesn't exist) entitiesWithTag lists for each tag
	m.TagEntity(e, tags...)
	// apply the unique tag if provided
	if uniqueTag != "" {
		if _, ok := m.uniqueEntities[uniqueTag]; ok {
			errorMsg := fmt.Sprintf("requested to spawn unique "+
				"entity for %s, but %s already exists", uniqueTag, uniqueTag)
			panic(errorMsg)
		}
		m.TagEntity(e, uniqueTag)
		m.uniqueEntities[uniqueTag] = e
	}
	// add logics
	e.Logics = make(map[string]*LogicUnit)
	for name, f := range logics {
		split := strings.Split(name, ",")
		if len(split) == 1 {
			e.AddLogic(name, f)
		} else if len(split) == 2 {
			fName := split[0]
			period, err := strconv.Atoi(split[1])
			if err != nil {
				Logger.Printf("[ERROR] failed to Atoi the specified ms period of a scheduled logic (%s)", name)
				panic(err)
			}
			e.AddLogicWithSchedule(fName, f, float64(period))
		} else {
			panic("malformed logic name! wants <name> or <name>,<ms_schedule>")
		}
	}
	// create funcset
	closureFuncs := make(map[string](func(params any) any))
	for name, f := range funcs {
		closureF := func(params any) any {
			return f(e, params)
		}
		closureFuncs[name] = closureF
	}
	e.funcs = NewFuncSet(closureFuncs)
	// add mind
	e.mind = mind
	// set entity active and notify entity is active
	m.setActiveState(e, active)
	// return Entity
	return e
}

func (m *EntityManager) ExpandEntityTables() {
	n := len(m.entityIDAllocator.currentEntities) / 2
	m.entityIDAllocator.expand(n)
	m.components.expand(n)
	for _, s := range m.w.systems {
		s.Expand(n)
	}
}
