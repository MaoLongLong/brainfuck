package internal

import (
	"container/list"
	"errors"
)

type Stack struct {
	data *list.List
}

func NewStack() *Stack {
	return &Stack{list.New()}
}

func (s *Stack) Push(v int) {
	s.data.PushBack(v)
}

func (s *Stack) Pop() (int, error) {
	if e := s.data.Back(); e != nil {
		s.data.Remove(e)
		return e.Value.(int), nil
	}
	return 0, errors.New("empty stack")
}

func (s *Stack) Empty() bool {
	return s.data.Len() == 0
}
