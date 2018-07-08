package engine

type TagList struct {
	Tags []string
}

func (l *TagList) Has(tag string) bool {
	for _, t := range l.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// TODO: ensure idempotent
func (l *TagList) Add(tag string) {
	l.Tags = append(l.Tags, tag)
}

// TODO: ensure idempotent
func (l *TagList) Remove(tag string) {
	removeStringFromSlice(&l.Tags, tag)
}

func (l *TagList) Copy() *TagList {
	tagsCopy := make([]string, len(l.Tags))
	copy(tagsCopy, l.Tags)
	return &TagList{Tags: tagsCopy}
}
