package sameriver

import (
	"encoding/json"
	"fmt"
	"sort"
)

type ItemSpec struct {
	Name        string
	DisplayName string
	FlavourText string
	Properties  map[string]int
	Tags        []string
	Count       int
}

func NewItem(name, displayName, flavourText string, properties map[string]int, Tags []string) *ItemSpec {
	return &ItemSpec{
		Name:        name,
		DisplayName: displayName,
		FlavourText: flavourText,
		Properties:  properties,
		Tags:        Tags,
		Count:       1,
	}
}

func (i *ItemSpec) copyOf() *ItemSpec {
	result := *i
	result.Properties = make(map[string]int)
	for key := range i.Properties {
		result.Properties[key] = i.Properties[key]
	}
	return &result
}

func ItemFromArchetype(arch *ItemSpec, nameSuffix, displayName, flavourText string) ItemSpec {
	i := ItemSpec{
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

func (i *ItemSpec) PropertiesForDisplay() []string {
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

func (i *ItemSpec) AddTags(Tags ...string) {
	for _, t := range Tags {
		i.Tags = append(i.Tags, t)
	}
}

func (i *ItemSpec) HasTag(tag string) bool {
	for _, t := range i.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (i *ItemSpec) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
