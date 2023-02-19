package sameriver

import (
	"encoding/json"
)

type ItemArchetype struct {
	Name        string
	DisplayName string
	FlavourText string
	Properties  map[string]int
	Tags        TagList
	Entity      map[string]any
}

func (i *ItemArchetype) copyOf() *ItemArchetype {
	result := *i
	result.Properties = make(map[string]int)
	for key := range i.Properties {
		result.Properties[key] = i.Properties[key]
	}
	return &result
}

func (i *ItemArchetype) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
