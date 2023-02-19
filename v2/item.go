package sameriver

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Item struct {
	Archetype  *ItemArchetype
	properties map[string]int
	Tags       TagList
	Count      int
	// used for lazy computation
	// computed lazily
	totalPropertiesDirty bool
	totalProperties      map[string]int `json:"-"` // (and not included in json representation)
	// computed lazily
	propertiesForDisplayDirty bool
	propertiesForDisplay      []string
}

func (i *Item) copyOf() *Item {
	c := &Item{
		Archetype:  i.Archetype,
		properties: make(map[string]int),
		Tags:       i.Tags.CopyOf(),
		Count:      i.Count,

		totalPropertiesDirty:      i.totalPropertiesDirty,
		totalProperties:           i.totalProperties,
		propertiesForDisplayDirty: i.propertiesForDisplayDirty,
		propertiesForDisplay:      i.propertiesForDisplay,
	}
	for k, v := range i.properties {
		c.properties[k] = v
	}
	return c
}

func (i *Item) SetProperty(k string, v int) {
	i.properties[k] = v
	i.totalPropertiesDirty = true
	i.propertiesForDisplayDirty = true
}

func (i *Item) GetProperty(k string) int {
	if v, ok := i.properties[k]; ok {
		return v
	} else {
		return i.Archetype.Properties[k]
	}
}

func (i *Item) GetTotalProperties() map[string]int {
	if !i.totalPropertiesDirty && i.totalProperties != nil {
		return i.totalProperties
	} else {
		i.totalProperties = make(map[string]int)
		for k, v := range i.Archetype.Properties {
			i.totalProperties[k] = v
		}
		for k, v := range i.properties {
			i.totalProperties[k] = v
		}
		i.totalPropertiesDirty = false
		return i.totalProperties
	}
}

func (i *Item) PropertiesForDisplay() []string {
	if !i.propertiesForDisplayDirty && i.propertiesForDisplay != nil {
		return i.propertiesForDisplay
	} else {
		result := make([]string, 0)
		totalProperties := i.GetTotalProperties()
		for k, v := range totalProperties {
			var displayStr string
			displayStr = fmt.Sprintf("%s %d", k, v)
			result = append(result, displayStr)
		}
		sort.Strings(result)
		i.propertiesForDisplay = result
		i.propertiesForDisplayDirty = false
		return i.propertiesForDisplay
	}
}

func (i *Item) TagsForDisplay() []string {
	return i.Tags.AsSlice()
}

func (i *Item) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
