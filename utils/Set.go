package utils

type Set struct {
	m map[string]bool
}

func NewSet() *Set {
	return &Set{
		make(map[string]bool),
	}
}

func (s *Set) addAll(values []string) {
	for _, value := range values {
		s.m[value] = true
	}
}

func (s *Set) Push(value string) {
	s.m[value] = true
}

func (s *Set) pop(value string) {
	delete(s.m, value)
}

func (s *Set) remove(values []string) {
	for _, value := range values {
		delete(s.m, value)
	}
}

func (s *Set) reset() {
	s.m = make(map[string]bool)
}

func (s *Set) Array() []string {
	var result []string
	for key, _ := range s.m {
		result = append(result, key)
	}
	return result
}
