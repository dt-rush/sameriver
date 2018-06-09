package engine

import (
	"sync"
)

type TagList struct {
	tags  []string
	Mutex sync.RWMutex
}

func (l *TagList) Has(tag string) bool {
	for _, t := range l.tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (l *TagList) Add(tag string) {
	l.tags = append(l.tags, tag)
}

func (l *TagList) Remove(tag string) {
	removeStringFromSlice(&l.tags, tag)
}

func (l *TagList) Copy() TagList {
	tagsCopy := make([]string, len(l.tags))
	copy(tagsCopy, l.tags)
	return TagList{tags: tagsCopy}
}
