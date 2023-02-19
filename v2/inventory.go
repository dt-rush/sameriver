package sameriver

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Inventory struct {
	Items []*Item
}

func NewInventory() Inventory {
	return Inventory{
		Items: make([]*Item, 0),
	}
}

func (i *Inventory) copyOf() Inventory {
	i2 := NewInventory()
	i2.Items = make([]*Item, len(i.Items))
	for ix, item := range i.Items {
		i2.Items[ix] = item.copyOf()
	}
	return i2
}

func (i *Inventory) ItemsForDisplay() []string {
	result := make([]string, 0)
	for _, item := range i.Items {
		displayStr := fmt.Sprintf("%s x %d", item.Archetype.DisplayName, item.Count)
		result = append(result, displayStr)
	}
	sort.Strings(result)
	return result
}

func (i *Inventory) Delete(item *Item) {
	for ix := 0; ix < len(i.Items); ix++ {
		if i.Items[ix] == item {
			i.Items = append(i.Items[:ix], i.Items[ix+1:]...)
		}
	}
}

func (i *Inventory) Debit(item *Item) *Item {
	retrieved := item.copyOf()
	item.Count -= 1
	retrieved.Count = 1
	if item.Count <= 0 {
		i.Delete(item)
	}
	return retrieved
}

func (i *Inventory) DebitN(item *Item, count int) *Item {
	retrieved := item.copyOf()
	item.Count -= count
	retrieved.Count = count
	if item.Count <= 0 {
		i.Delete(item)
	}
	return retrieved
}

func (i *Inventory) DebitAll(item *Item) *Item {
	retrieved := item.copyOf()
	i.Delete(item)
	return retrieved
}

func (i *Inventory) DebitName(name string) *Item {
	retrieved := i.NameFilter(name)[0]
	return i.Debit(retrieved)
}

func (i *Inventory) DebitNName(name string, n int) *Item {
	retrieved := i.NameFilter(name)[0]
	return i.DebitN(retrieved, n)
}

func (i *Inventory) DebitAllName(name string) *Item {
	retrieved := i.NameFilter(name)[0]
	return i.DebitAll(retrieved)
}

func (i *Inventory) Credit(item *Item) {
	i.Items = append(i.Items, item.copyOf())
}

func (i *Inventory) CreditN(item *Item, n int) {
	toAppend := item.copyOf()
	toAppend.Count *= n
	i.Items = append(i.Items, toAppend)
}

func (i *Inventory) Filter(predicate func(*Item) bool) []*Item {
	results := make([]*Item, 0)
	for _, item := range i.Items {
		if predicate(item) {
			results = append(results, item)
		}
	}
	return results
}

func (i *Inventory) NameFilter(name string) []*Item {
	predicate := func(i *Item) bool {
		return i.Archetype.Name == name
	}
	return i.Filter(predicate)
}

func (i *Inventory) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}
