package sameriver

import (
	"encoding/json"
)

type ItemArchetype struct {
	Name        string
	DisplayName string
	FlavourText string
	Properties  map[string]float64
	Tags        TagList
	Entity      map[string]any
}

func (i *ItemArchetype) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
