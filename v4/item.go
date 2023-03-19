package sameriver

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Item struct {
	sys             *ItemSystem
	inv             *Inventory
	Archetype       string
	DisplayStr      string
	Properties      map[string]float64
	Tags            TagList
	Count           int
	Degradations    []float64 `json:",omitempty"`
	degradationRate float64

	// computed lazily
	propertiesForDisplayDirty bool
	propertiesForDisplay      []string
}

func (i *Item) CopyOf() *Item {
	c := &Item{
		sys:             i.sys,
		inv:             i.inv,
		Archetype:       i.Archetype,
		Properties:      make(map[string]float64),
		Tags:            i.Tags.CopyOf(),
		Count:           i.Count,
		degradationRate: i.degradationRate,

		propertiesForDisplayDirty: i.propertiesForDisplayDirty,
		propertiesForDisplay:      i.propertiesForDisplay,
	}
	for k, v := range i.Properties {
		c.Properties[k] = v
	}
	copy(c.Degradations, i.Degradations)
	return c
}

func (i *Item) GetArchetype() *ItemArchetype {
	return i.sys.Archetypes[i.Archetype]
}

func (i *Item) SetProperty(k string, v float64) {
	if k == "degradation" {
		for ix := range i.Degradations {
			i.Degradations[ix] = v
		}
	} else {
		i.Properties[k] = v
	}
	i.propertiesForDisplayDirty = true
}

func (i *Item) GetProperty(k string) float64 {
	if k == "degradation" {
		return i.Degradations[0]
	} else {
		return i.Properties[k]
	}
}

const (
	ITEM_MOST_DEGRADED = iota
	ITEM_LEAST_DEGRADED
)

func (i *Item) DebitStack(n int, leastOrMostDegraded int) *Item {
	if n <= 0 {
		panic(fmt.Sprintf("Tried to Debit stack %s for %d items", i.DisplayStr, n))
	}
	if n > i.Count {
		panic(fmt.Sprintf("Tried to Debit stack of %d items for %d", i.Count, n))
	}
	// new stack
	result := i.CopyOf()
	// update counts
	i.Count -= n
	result.Count = n
	// reevaluate display strings
	i.reevaluateDisplayStr()
	result.reevaluateDisplayStr()
	// divide Degradations (they are sorted in increasing order)
	if i.Degradations != nil {
		full := i.Degradations
		result.Degradations = make([]float64, n)
		switch leastOrMostDegraded {
		case ITEM_MOST_DEGRADED:
			i.Degradations = full[:len(full)-n]
			copy(result.Degradations, full[len(full)-n:])
		case ITEM_LEAST_DEGRADED:
			i.Degradations = full[n+1:]
			copy(result.Degradations, full[:n])
		}
	}
	return result
}

func (i *Item) CreditStack(stack *Item) {
	if stack.Count == 0 {
		panic(fmt.Sprintf("Tried to Credit using a stack %s (count 0)", stack.DisplayStr))
	}
	i.Count += stack.Count
	i.reevaluateDisplayStr()
	i.Degradations = append(i.Degradations, stack.Degradations...)
	sort.Float64s(i.Degradations)
}

/*
in the case of growth, we stretch the degradations slice
and fill in the gaps with their nearest left value

000000444        len 9
________________ len 16

16/9 = 1.777777

skip ahead by 1.777777, int()'d

0_0_00_0_0_44_4_

and fill in

0000000000044444

in the case of shrinking, we sample the list down into the new one

000444999 len 9
_____     len 5

9/5 = 1.8

skip ahead by 1.8, int()'d

v v vv v
000444999

00449
*/
func (i *Item) SetCount(n int) {
	if n == i.Count {
		return
	}
	if n == 0 {
		i.Degradations = []float64{}
		i.Count = 0
		return
	}
	if i.Tags.Has("degrades") {
		sampleIx := func(ix int, factor float64) int {
			return int(float64(ix) * factor)
		}
		newDegradations := make([]float64, n)
		if n > i.Count {
			factor := float64(n) / float64(i.Count)
			for ix := 0; ix < i.Count; ix++ {
				// find our sample point and the next one (defines the gap between)
				samplePoint := sampleIx(ix, factor)
				nextPoint := sampleIx(ix+1, factor)
				newDegradations[samplePoint] = i.Degradations[ix]
				// fill in the space til the next
				for jx := samplePoint + 1; jx < nextPoint && jx < n; jx++ {
					newDegradations[jx] = i.Degradations[ix]
				}
			}
		} else if n < i.Count {
			factor := float64(i.Count) / float64(n)
			for ix := 0; ix < n; ix++ {
				samplePoint := sampleIx(ix, factor)
				newDegradations[ix] = i.Degradations[samplePoint]
			}
		}
	}
	i.Count = n
}

func (i *Item) PropertiesForDisplay() []string {

	formatFloatForDisplay := func(value float64) string {
		var formattedValue string
		if value == float64(int(value)) {
			formattedValue = fmt.Sprintf("%d", int(value))
		} else if float64(int(value*10)/10) == 0 {
			formattedValue = fmt.Sprintf("%.1f", value)
			if formattedValue == "0.0" {
				formattedValue = "0.0+"
			} else if formattedValue == "-0.0" {
				formattedValue = "0.0-"
			}
		} else {
			formattedValue = fmt.Sprintf("%.1f", value)
		}
		return formattedValue
	}

	if !i.propertiesForDisplayDirty && i.propertiesForDisplay != nil {
		return i.propertiesForDisplay
	} else {
		result := make([]string, 0)
		for k, v := range i.Properties {
			displayStr := fmt.Sprintf("%s %s", k, formatFloatForDisplay(v))
			result = append(result, displayStr)
		}
		// include degradation value in properties, even though it's not actually
		// in the properties map (this is so that stacks can merge ignoring
		// differences in degradations)
		if i.Tags.Has("degrades") {
			if i.Count == 1 {
				str := fmt.Sprintf("degradation %d", int(i.Degradations[0]))
				result = append(result, str)
			} else if i.Count > 1 {
				str := fmt.Sprintf("degradations %d-%d",
					int(i.Degradations[0]),
					int(i.Degradations[len(i.Degradations)-1]))
				result = append(result, str)
			}
		}
		sort.Strings(result)
		i.propertiesForDisplay = result
		i.propertiesForDisplayDirty = false
		return i.propertiesForDisplay
	}
}

func (i *Item) propertiesAndTagsMatch(other *Item) bool {
	for k, v := range i.Properties {
		if other.Properties[k] != v {
			return false
		}
	}
	if i.Tags.Length() != other.Tags.Length() {
		return false
	} else {
		for _, t := range i.Tags.AsSlice() {
			if !other.Tags.Has(t) {
				return false
			}
		}
	}
	return true
}

func (i *Item) HasProperty(k string) bool {
	_, ok := i.Properties[k]
	return ok
}

func (i *Item) DisplayName() string {
	return i.sys.Archetypes[i.Archetype].DisplayName
}

func (i *Item) reevaluateDisplayStr() {
	i.DisplayStr = fmt.Sprintf("%s x %d", i.DisplayName(), i.Count)
}

func (i *Item) TagsForDisplay() []string {
	return i.Tags.AsSlice()
}

func (i *Item) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
