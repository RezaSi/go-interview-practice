package generics

import "errors"

// ============================================================================
// 1. Generic Pair
// ============================================================================

type Pair[T, U any] struct {
	First  T
	Second U
}

func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{First: p.Second, Second: p.First}
}

// ============================================================================
// 2. Generic Stack
// ============================================================================

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.items) == 0 {
		return zero, errors.New("stack is empty")
	}
	index := len(s.items) - 1
	value := s.items[index]
	s.items = s.items[:index]
	return value, nil
}

func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.items) == 0 {
		return zero, errors.New("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// ============================================================================
// 3. Generic Queue
// ============================================================================

type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(value T) {
	q.items = append(q.items, value)
}

func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, errors.New("queue is empty")
	}
	value := q.items[0]
	q.items = q.items[1:]
	return value, nil
}

func (q *Queue[T]) Front() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, errors.New("queue is empty")
	}
	return q.items[0], nil
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// ============================================================================
// 4. Generic Set
// ============================================================================

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

func (s *Set[T]) Add(value T) {
	s.items[value] = struct{}{}
}

func (s *Set[T]) Remove(value T) {
	delete(s.items, value)
}

func (s *Set[T]) Contains(value T) bool {
	_, exists := s.items[value]
	return exists
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) Elements() []T {
	elements := make([]T, 0, len(s.items))
	for k := range s.items {
		elements = append(elements, k)
	}
	return elements
}

func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s1.items {
		result.Add(k)
	}
	for k := range s2.items {
		result.Add(k)
	}
	return result
}

func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s1.items {
		if s2.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s1.items {
		if !s2.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// ============================================================================
// 5. Generic Utility Functions
// ============================================================================

func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = mapper(v)
	}
	return result
}

func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func FindIndex[T comparable](slice []T, element T) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, v := range slice {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}