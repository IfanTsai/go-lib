package set

var exist = struct{}{}

type Set struct {
	m map[interface{}]struct{}
}

func NewSet(items ...interface{}) *Set {
	s := &Set{
		m: make(map[interface{}]struct{}),
	}

	for _, item := range items {
		s.Add(item)
	}

	return s
}

func (s *Set) Add(item interface{}) {
	s.m[item] = exist
}

func (s *Set) AddSet(addSet *Set) {
	for item := range addSet.m {
		s.Add(item)
	}
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]

	return ok
}

func (s *Set) Remove(item interface{}) {
	delete(s.m, item)
}

func (s *Set) Clear() {
	if !s.IsEmpty() {
		s.m = make(map[interface{}]struct{})
	}
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Set) Compare(compareSet *Set) bool {
	if s.Size() != compareSet.Size() {
		return false
	}

	for item := range s.m {
		if !compareSet.Contains(item) {
			return false
		}
	}

	return true
}

func (s *Set) ConvertSlice() []interface{} {
	items := make([]interface{}, s.Size())

	i := 0
	for item := range s.m {
		items[i] = item
		i++
	}

	return items
}
