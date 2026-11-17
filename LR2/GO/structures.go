package main

// Свой динамический массив — аналог Array<T> из C++.
type Array[T any] struct {
	data []T
	size int
}

func (a *Array[T]) Size() int { return a.size }

func (a *Array[T]) Push(value T) {
	if a.size == len(a.data) {
		newCapacity := len(a.data) * 2
		if newCapacity == 0 {
			newCapacity = 1
		}
		newData := make([]T, newCapacity)
		copy(newData, a.data[:a.size])
		a.data = newData
	}
	a.data[a.size] = value
	a.size++
}

func (a *Array[T]) Get(index int) T { return a.data[index] }

func (a *Array[T]) Set(index int, value T) { a.data[index] = value }

func (a *Array[T]) Remove(index int) {
	if index < 0 || index >= a.size {
		return
	}
	copy(a.data[index:], a.data[index+1:a.size])
	a.size--
	var zero T
	a.data[a.size] = zero
}

func (a *Array[T]) Clear() { a.data = nil; a.size = 0 }

// Односвязный список — аналог ForwardList<T> из C++.
type Node[T comparable] struct {
	value T
	next  *Node[T]
}

type ForwardList[T comparable] struct{ head *Node[T] }

func (l *ForwardList[T]) PushFront(value T) {
	l.head = &Node[T]{value: value, next: l.head}
}

func (l *ForwardList[T]) PushBack(value T) {
	node := &Node[T]{value: value}
	if l.head == nil {
		l.head = node
		return
	}
	current := l.head
	for current.next != nil {
		current = current.next
	}
	current.next = node
}

func (l *ForwardList[T]) DelValue(value T) bool {
	var previous *Node[T]
	for current := l.head; current != nil; current = current.next {
		if current.value == value {
			if previous == nil {
				l.head = current.next
			} else {
				previous.next = current.next
			}
			return true
		}
		previous = current
	}
	return false
}

// Стек использует те же узлы, что односвязный список.
type Stack[T comparable] struct{ head *Node[T] }

func (s *Stack[T]) Push(value T) { s.head = &Node[T]{value: value, next: s.head} }
func (s *Stack[T]) Empty() bool  { return s.head == nil }
func (s *Stack[T]) Pop() (T, bool) {
	if s.head == nil {
		var zero T
		return zero, false
	}
	value := s.head.value
	s.head = s.head.next
	return value, true
}
