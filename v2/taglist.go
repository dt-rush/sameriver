package sameriver

import (
	"sort"

	"encoding/json"
)

type TagList struct {
	tags map[string]bool
	// sorted slice representation is computed lazily
	dirty bool
	slice []string
}

func NewTagList() TagList {
	l := TagList{}
	l.tags = make(map[string]bool)
	return l
}

func (l *TagList) Has(tags ...string) bool {
	ok := true
	for _, tag := range tags {
		_, has := l.tags[tag]
		ok = ok && has
	}
	return ok
}

func (l *TagList) Add(tags ...string) {
	if l.tags == nil {
		l.tags = make(map[string]bool, 1)
	}
	for _, t := range tags {
		l.tags[t] = true
	}
	l.dirty = true
}

func (l *TagList) MergeIn(l2 TagList) {
	for t, _ := range l2.tags {
		l.tags[t] = true
	}
	l.dirty = true
}

func (l *TagList) Remove(tag string) {
	delete(l.tags, tag)
	l.dirty = true
}

func (l *TagList) CopyOf() TagList {
	tagsCopy := make(map[string]bool, len(l.tags))
	for tag, _ := range l.tags {
		tagsCopy[tag] = true
	}
	return TagList{
		tags:  tagsCopy,
		dirty: l.dirty,
		slice: l.slice,
	}
}

func (l *TagList) AsSlice() []string {
	if !l.dirty && (l.slice != nil) {
		return l.slice
	} else {
		slice := make([]string, 0, len(l.tags))
		for tag, _ := range l.tags {
			slice = append(slice, tag)
		}
		sort.Strings(slice)
		l.slice = slice
		l.dirty = false
		return l.slice
	}
}

func (l *TagList) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.AsSlice())
}
