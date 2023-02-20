package sameriver

import (
	"encoding/json"
	"sort"
	"strings"
)

type Inventory struct {
	Stacks []*Item
}

func NewInventory() *Inventory {
	i := &Inventory{
		Stacks: make([]*Item, 0),
	}
	return i
}

func (i *Inventory) CopyOf() *Inventory {
	i2 := NewInventory()
	i2.Stacks = make([]*Item, len(i.Stacks))
	for ix, stack := range i.Stacks {
		i2.Stacks[ix] = stack.CopyOf()
	}
	return i2
}

func (i *Inventory) ItemsForDisplay() []*Item {
	return i.StacksForDisplay()
}

func (i *Inventory) StacksForDisplay() []*Item {
	result := make([]*Item, 0)
	for _, item := range i.Stacks {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].DisplayName(), result[j].DisplayName()) == -1
	})
	return result
}

func (i *Inventory) Delete(stack *Item) {
	for ix := 0; ix < len(i.Stacks); ix++ {
		if i.Stacks[ix] == stack {
			i.Stacks = append(i.Stacks[:ix], i.Stacks[ix+1:]...)
		}
	}
}

func (i *Inventory) DebitWithPreference(stack *Item, leastOrMostDegraded int) (result *Item) {
	result = stack.DebitStack(1, leastOrMostDegraded)
	if stack.Count == 0 {
		i.Delete(stack)
	}
	return result
}

func (i *Inventory) Debit(stack *Item) (result *Item) {
	return i.DebitWithPreference(stack, ITEM_MOST_DEGRADED)
}

func (i *Inventory) DebitNWithPreference(stack *Item, n int, leastOrMostDegraded int) *Item {
	if n == stack.Count {
		i.Delete(stack)
		return stack
	} else {
		return stack.DebitStack(n, leastOrMostDegraded)
	}
}

func (i *Inventory) DebitN(stack *Item, n int) *Item {
	return i.DebitNWithPreference(stack, n, ITEM_MOST_DEGRADED)
}

func (i *Inventory) DebitAll(stack *Item) *Item {
	i.Delete(stack)
	return stack
}

func (i *Inventory) Credit(stack *Item) {
	if stack.Count == 0 {
		Logger.Printf("[WARNING] trying to Credit a stack (archetype %s) with count 0! Doing nothing.", stack.Archetype)
		return
	}
	// try to find a stack that this can be put in
	for _, s := range i.Stacks {
		if stack.Archetype == s.Archetype &&
			stack.propertiesAndTagsMatch(s) {
			s.CreditStack(stack)
			return
		}
	}
	// else we just append
	i.Stacks = append(i.Stacks, stack)
	stack.inv = i
}

func (i *Inventory) Filter(predicate func(*Item) bool) []*Item {
	results := make([]*Item, 0)
	for _, item := range i.Stacks {
		if predicate(item) {
			results = append(results, item)
		}
	}
	return results
}

func (i *Inventory) NameFilter(name string) []*Item {
	predicate := func(i *Item) bool {
		return i.GetArchetype().Name == name
	}
	return i.Filter(predicate)
}

func (i *Inventory) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
