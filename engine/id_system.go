/**
  * 
  * 
  * 
  * 
**/



package engine


type IDSystem struct {
	// records all ids given out
	ids []int

	// used to gen
	seed int
}


func (s *IDSystem) Init (capacity int) {
	s.ids = make ([]int, 0)
	s.seed = -1
}


func (s *IDSystem) Gen () int {
	s.seed++
	s.ids = append (s.ids, s.seed)
	return s.seed
}

func (s *IDSystem) GetIDs () []int {
	return s.ids
}

func (s *IDSystem) NumberOfIDs () int {
	return len (s.ids)
}
