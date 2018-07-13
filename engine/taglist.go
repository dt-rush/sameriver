package engine

type TagList struct {
	Tags map[string]bool
}

func (l *TagList) Has(tags ...string) bool {
	ok := true
	for _, tag := range tags {
		_, has := l.Tags[tag]
		ok = ok && has
	}
	return ok
}

func (l *TagList) Add(tag string) {
	if l.Tags == nil {
		l.Tags = make(map[string]bool, 1)
	}
	l.Tags[tag] = true
}

func (l *TagList) Remove(tag string) {
	delete(l.Tags, tag)
}

func (l *TagList) Copy() *TagList {
	tagsCopy := make(map[string]bool, len(l.Tags))
	for tag, _ := range l.Tags {
		tagsCopy[tag] = true
	}
	return &TagList{Tags: tagsCopy}
}

func (l *TagList) ToSlice() []string {
	slice := make([]string, 0, len(l.Tags))
	for tag, _ := range l.Tags {
		slice = append(slice, tag)
	}
	return slice
}
