package sameriver

import (
	"encoding/json"
	"fmt"
	"math"
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
	result = append(result, i.Stacks...)
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

func (i *Inventory) DebitNName(n int, name string) []*Item {
	return i.DebitNFilter(n, func(it *Item) bool {
		return it.GetArchetype().Name == name
	})
}

func (i *Inventory) DebitFilter(predicate func(*Item) bool) []*Item {
	items := i.Filter(predicate)
	for _, it := range items {
		i.DebitAll(it)
	}
	return items
}

func (i *Inventory) DebitNFilter(n int, predicate func(*Item) bool) []*Item {
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

func (i *Inventory) DebitAllFilter(predicate func(*Item) bool) []*Item {
	results := make([]*Item, 0)
	for _, s := range i.Filter(predicate) {
		results = append(results, i.DebitAll(s))
	}
	return results
}

func (i *Inventory) DebitTags(tags ...string) *Item {
	return i.DebitNTags(1, tags...)[0]
}

func (i *Inventory) DebitNTags(n int, tags ...string) []*Item {
	return i.DebitNFilter(n, func(it *Item) bool {
		return it.HasTags(tags...)
	})
}

func (i *Inventory) DebitAllTags(tags ...string) []*Item {
	stacks := make([]*Item, 0)
	for _, it := range i.FilterTags(tags...) {
		stacks = append(stacks, i.DebitAll(it))
	}
	return stacks
}

func (i *Inventory) DebitAllName(name string) []*Item {
	stacks := make([]*Item, 0)
	for _, it := range i.FilterName(name) {
		stacks = append(stacks, i.DebitAll(it))
	}
	return stacks
}

func (i *Inventory) Credit(stack *Item) {
	if stack.Count == 0 {
		// a stack of zero items was probably created by a random number generator; ignore
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
	debited := inv.DebitNFilter(n,
		func(s *Item) bool { return s.GetArchetype().Name == name })
	for _, it := range debited {
		i.Credit(it)
	}
}

func (i *Inventory) GetNByFilter(inv *Inventory, n int, predicate func(*Item) bool) {
	for _, it := range inv.DebitNFilter(n, predicate) {
		i.Credit(it)
	}
}

func (i *Inventory) GetAllByName(inv *Inventory, name string) {
	for _, it := range inv.FilterName(name) {
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

func (i *Inventory) setCount(n int, filtered []*Item) {
	if n == 0 {
		for _, it := range filtered {
			i.DebitAll(it)
		}
	}

	count := 0
	for _, s := range filtered {
		count += s.Count
	}
	if count == n {
		return
	}

	diff := n - count
	for ix := 0; diff != 0; ix++ {
		it := filtered[ix]
		if diff < 0 {
			// decrease this stack as much as we can to fulfill diff
			newCount := int(math.Max(0, float64(it.Count+diff)))
			change := it.Count - newCount
			it.SetCount(newCount)
			diff += change
		} else if diff > 0 {
			// simply set this first stack to the count needed to fulfill diff
			newCount := it.Count + diff
			it.SetCount(newCount)
			diff = 0
		}
	}
}

func (i *Inventory) SetCountName(n int, archetype string) {
	filtered := i.FilterName(archetype)
	if len(filtered) == 0 {
		panic(fmt.Sprintf("Can't SetCountName(%d, %s) since no items matched that archetype!", n, archetype))
	}
	i.setCount(n, filtered)
}

func (i *Inventory) SetCountTags(n int, tags ...string) {
	filtered := i.FilterTags(tags...)
	if len(filtered) == 0 {
		panic(fmt.Sprintf("Can't SetCountTags(%d, %s) since no items matched tags!", n, tags))
	}
	i.setCount(n, filtered)
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

func (i *Inventory) CountName(name string) int {
	n := 0
	for _, stack := range i.Stacks {
		if stack.GetArchetype().Name == name {
			n += stack.Count
		}
	}
	return n
}

func (i *Inventory) CountTags(tags ...string) int {
	n := 0
	for _, stack := range i.FilterTags(tags...) {
		n += stack.Count
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

func (i *Inventory) FilterTags(tags ...string) []*Item {
	check := func(s *Item) bool {
		for _, t := range tags {
			if !s.Tags.Has(t) {
				return false
			}
		}
		return true
	}
	results := make([]*Item, 0)
	for _, item := range i.Stacks {
		if check(item) {
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

func (i *Inventory) FilterName(name string) []*Item {
	predicate := func(it *Item) bool {
		return it.GetArchetype().Name == name
	}
	return i.Filter(predicate)
}

func (i *Inventory) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
