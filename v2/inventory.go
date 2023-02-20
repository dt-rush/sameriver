package sameriver

import (
	"encoding/json"
	"fmt"
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

func (i *Inventory) DebitNWithPreference(n int, stack *Item, leastOrMostDegraded int) *Item {
	if n == stack.Count {
		i.Delete(stack)
		return stack
	} else {
		return stack.DebitStack(n, leastOrMostDegraded)
	}
}

func (i *Inventory) DebitN(n int, stack *Item) *Item {
	return i.DebitNWithPreference(n, stack, ITEM_MOST_DEGRADED)
}

func (i *Inventory) DebitAll(stack *Item) *Item {
	i.Delete(stack)
	return stack
}

func (i *Inventory) DebitByFilter(predicate func(*Item) bool) []*Item {
	items := i.Filter(predicate)
	for _, it := range items {
		i.DebitAll(it)
	}
	return items
}

func (i *Inventory) DebitNByFilter(n int, predicate func(*Item) bool) []*Item {
	count := i.Count(predicate)
	if count < n {
		panic(fmt.Sprintf("Tried to Debit %d items from inv, but only has %d", n, count))
	}
	stacks := make([]*Item, 0)
	remaining := n
	for _, it := range i.Filter(predicate) {
		if it.Count > remaining {
			stacks = append(stacks, i.DebitN(n, it))
		} else {
			stacks = append(stacks, i.DebitAll(it))
		}
	}
	return stacks
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
	// set the inv to this inv
	stack.inv = i
}

func (i *Inventory) GetNByName(inv *Inventory, n int, name string) {
	debited := inv.DebitNByFilter(n,
		func(s *Item) bool { return s.GetArchetype().Name == name })
	for _, it := range debited {
		i.Credit(it)
	}
}

func (i *Inventory) GetNByFilter(inv *Inventory, n int, predicate func(*Item) bool) {
	for _, it := range inv.DebitNByFilter(n, predicate) {
		i.Credit(it)
	}
}

func (i *Inventory) GetAllByName(inv *Inventory, name string) {
	for _, it := range inv.NameFilter(name) {
		i.Credit(inv.DebitAll(it))
	}
}

func (i *Inventory) GetAllByFilter(inv *Inventory, predicate func(*Item) bool) {
	for _, it := range inv.Filter(predicate) {
		i.Credit(inv.DebitAll(it))
	}
}

func (i *Inventory) GetAll(inv *Inventory) {
	for _, it := range inv.Stacks {
		i.Credit(inv.DebitAll(it))
	}
}

func (i *Inventory) CountName(name string) int {
	n := 0
	for _, stack := range i.NameFilter(name) {
		n += stack.Count
	}
	return n
}

func (i *Inventory) Count(predicate func(*Item) bool) int {
	n := 0
	for _, stack := range i.Stacks {
		if predicate(stack) {
			n += stack.Count
		}
	}
	return n
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

func (i *Inventory) ContainsName(name string) bool {
	return i.Contains(func(it *Item) bool {
		return it.GetArchetype().Name == name
	})
}

func (i *Inventory) Contains(predicate func(*Item) bool) bool {
	for _, item := range i.Stacks {
		if predicate(item) {
			return true
		}
	}
	return false
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
