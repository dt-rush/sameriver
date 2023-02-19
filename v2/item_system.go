package sameriver

import (
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"
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

func (i *ItemSystem) LoadArchetypesFile(filename string) {
	Logger.Printf("Loading item archetypes from %s...", filename)
	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Trying to open %s - doesn't exist", filename))
	}
	contents, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	i.LoadArchetypesJSON(contents)
}

func (i *ItemSystem) LoadArchetypesJSON(jsonStr []byte) {
	var archetypes [](map[string]*json.RawMessage)
	err := json.Unmarshal(jsonStr, &archetypes)
	if err != nil {
		panic(err)
	}
	for ix, jsonSpec := range archetypes {
		spec := make(map[string]any)
		if _, ok := jsonSpec["name"]; ok {
			var name string
			json.Unmarshal(*jsonSpec["name"], &name)
			spec["name"] = name
		} else {
			panic(fmt.Sprintf("object at index %d was missing \"name\" property", ix))
		}
		if _, ok := jsonSpec["displayName"]; ok {
			var displayName string
			json.Unmarshal(*jsonSpec["displayName"], &displayName)
			spec["displayName"] = displayName
		}
		if _, ok := jsonSpec["flavourText"]; ok {
			var flavourText string
			json.Unmarshal(*jsonSpec["flavourText"], &flavourText)
			spec["flavourText"] = flavourText
		}
		if _, ok := jsonSpec["properties"]; ok {
			var properties map[string]int
			json.Unmarshal(*jsonSpec["properties"], &properties)
			spec["properties"] = properties
		}
		if _, ok := jsonSpec["tags"]; ok {
			var tags []string
			json.Unmarshal(*jsonSpec["tags"], &tags)
			spec["tags"] = tags
		}
		i.CreateArchetype(spec)
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
