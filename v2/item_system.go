package sameriver

import (
	"fmt"
)

type ItemSystem struct {
	w            *World
	ItemEntities *UpdatedEntityList
	Archetypes   map[string]*ItemArchetype
}

func NewItemSystem() *ItemSystem {
	return &ItemSystem{
		Archetypes: make(map[string]*ItemArchetype),
	}
}

func (i *ItemSystem) registerArchetype(arch *ItemArchetype) {
	i.Archetypes[arch.Name] = arch
}

func (i *ItemSystem) CreateArchetype(spec map[string]any) {

	var name, displayName, flavourText string
	var properties map[string]int
	var tags []string

	if _, ok := spec["name"]; ok {
		name = spec["name"].(string)
	} else {
		panic("Must supply \"name\" to CreateArchetype")
	}
	if _, ok := spec["displayName"]; ok {
		displayName = spec["displayName"].(string)
	} else {
		displayName = name
	}
	if _, ok := spec["flavourText"]; ok {
		flavourText = spec["flavourText"].(string)
	} else {
		flavourText = ""
	}
	if _, ok := spec["properties"]; ok {
		properties = spec["properties"].(map[string]int)
	}
	if _, ok := spec["tags"]; ok {
		tags = spec["tags"].([]string)
	}

	tagList := NewTagList()
	tagList.Add(tags...)
	a := &ItemArchetype{
		Name:        name,
		DisplayName: displayName,
		FlavourText: flavourText,
		Properties:  properties,
		Tags:        tagList,
	}
	i.registerArchetype(a)
}

func (i *ItemSystem) CreateSubArchetype(spec map[string]any) {
	var parent string
	var name, displayName, flavourText string
	var properties map[string]int
	var tagDiff []string

	if _, ok := spec["parent"]; ok {
		parent = spec["parent"].(string)
		if _, ok := i.Archetypes[parent]; !ok {
			panic(fmt.Sprintf("Archetype %s not found at time needed", parent))
		}
	} else {
		panic("Must supply \"parent\" to CreateSubArchetype")
	}
	if _, ok := spec["name"]; ok {
		name = spec["name"].(string)
	} else {
		panic("Must supply \"name\" to CreateSubArchetype")
	}
	if _, ok := spec["displayName"]; ok {
		displayName = spec["displayName"].(string)
	} else {
		displayName = i.Archetypes[parent].DisplayName
	}
	if _, ok := spec["flavourText"]; ok {
		flavourText = spec["flavourText"].(string)
	} else {
		flavourText = i.Archetypes[parent].FlavourText
	}
	if _, ok := spec["properties"]; ok {
		properties = spec["properties"].(map[string]int)
	} else {
		properties = make(map[string]int)
	}
	if _, ok := spec["tagDiff"]; ok {
		tagDiff = spec["tagDiff"].([]string)
	}

	a := &ItemArchetype{
		Name:        name,
		DisplayName: displayName,
		FlavourText: flavourText,
	}
	// copy in parent arch properties
	a.Properties = make(map[string]int)
	for k, v := range i.Archetypes[parent].Properties {
		a.Properties[k] = v
	}
	// supplied properties shadow those of the arch parent
	for k, v := range properties {
		a.Properties[k] = v
	}
	// copy parent arch tags then apply the tagdiff specifications
	tags := i.Archetypes[parent].Tags.CopyOf()
	for _, tSpec := range tagDiff {
		op, t := tSpec[0:1], tSpec[1:]
		if op == "+" {
			tags.Add(t)
		} else if op == "-" {
			tags.Remove(t)
		}
	}
	a.Tags = tags

	i.registerArchetype(a)
}

func (i *ItemSystem) CreateItem(spec map[string]any) *Item {
	var arch string
	var properties map[string]int
	var tags []string
	var count int
	if _, ok := spec["archetype"]; ok {
		arch = spec["archetype"].(string)
		if _, ok := i.Archetypes[arch]; !ok {
			panic(fmt.Sprintf("Trying to create item for archetype that isn't created yet: %s", arch))
		}
	} else {
		panic("Must specify \"archetype\" in CreateItem()")
	}
	if _, ok := spec["properties"]; ok {
		properties = spec["properties"].(map[string]int)
	} else {
		properties = make(map[string]int)
	}
	if _, ok := spec["tags"]; ok {
		tags = spec["tags"].([]string)
	}
	if _, ok := spec["count"]; ok {
		count = spec["count"].(int)
	} else {
		count = 1
	}

	tagList := NewTagList()
	tagList.Add(tags...)
	tagList.MergeIn(i.Archetypes[arch].Tags)

	return &Item{
		Archetype:  i.Archetypes[arch],
		Properties: properties,
		Tags:       tagList,
		Count:      count,
	}
}

// System funcs

func (i *ItemSystem) GetComponentDeps() []string {
	return []string{"Bool, Item"}
}

func (i *ItemSystem) LinkWorld(w *World) {
	i.w = w

	i.ItemEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromTag("item"))
}

func (i *ItemSystem) Update(dt_ms float64) {
	// nil?
}
