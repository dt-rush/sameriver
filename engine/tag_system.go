/**
  * 
  * 
  * 
  * 
**/



package engine



type TagSystem struct {
	// two one-way maps support a many-to-many relationship
	// tag -> []IDs
	tag_entities map[string]([]int)
	// ID -> []tag
	entity_tags map[int]([]string)
	// used to size the lists pointed-to
	capacity int
}

func (m *TagSystem) Init (capacity int) {
	m.tag_entities = make (map[string]([]int))
	m.entity_tags = make (map[int]([]string))
	m.capacity = capacity
}

func (m *TagSystem) TagEntity (id int, tag string) {
	_, et_ok := m.entity_tags [id]
	_, te_ok := m.tag_entities [tag]
	if ! et_ok {
		m.entity_tags [id] = make ([]string, 0)
	}
	if ! te_ok {
		m.tag_entities [tag] = make ([]int, 0)
	}
	m.entity_tags [id] = append (m.entity_tags [id], tag)
	m.tag_entities [tag] = append (m.tag_entities [tag], id)
}


