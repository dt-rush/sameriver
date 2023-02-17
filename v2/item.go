package sameriver

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Item struct {
	Name        string
	DisplayName string
	FlavourText string
	Properties  map[string]int
	Tags        []string
	Count       int
}

func NewItem(name, displayName, flavourText string, properties map[string]int, Tags []string) *Item {
	return &Item{
		Name:        name,
		DisplayName: displayName,
		FlavourText: flavourText,
		Properties:  properties,
		Tags:        Tags,
		Count:       1,
	}
}

func (i *Item) copyOf() *Item {
	result := *i
	result.Properties = make(map[string]int)
	for key := range i.Properties {
		result.Properties[key] = i.Properties[key]
	}
	return &result
}

func ItemFromArchetype(arch *Item, nameSuffix, displayName, flavourText string) Item {
	i := Item{
		Name:        fmt.Sprintf("%s-%s", arch.Name, nameSuffix),
		DisplayName: displayName,
		FlavourText: flavourText,
		Count:       1,
	}
	i.Properties = make(map[string]int)
	for k, v := range arch.Properties {
		i.Properties[k] = v
	}
	i.Tags = make([]string, len(arch.Tags))
	copy(i.Tags, arch.Tags)
	return i
}

func (i *Item) PropertiesForDisplay() []string {
	result := make([]string, 0)
	for k, v := range i.Properties {
		var displayStr string
		if v == 1 {
			displayStr = k
		} else {
			displayStr = fmt.Sprintf("%s %d", k, v)
		}
		result = append(result, displayStr)
	}
	sort.Strings(result)
	return result
}

func (i *Item) AddTags(Tags ...string) {
	for _, t := range Tags {
		i.Tags = append(i.Tags, t)
	}
}

func (i *Item) HasTag(tag string) bool {
	for _, t := range i.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (i *Item) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
