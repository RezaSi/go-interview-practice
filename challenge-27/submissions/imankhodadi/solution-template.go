package generics

import (
	"errors"
)

var ErrEmptyCollection = errors.New("collection is empty")

// 1. Generic Pair
// Pair represents a generic pair of values of potentially different types
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

// 2. Generic Stack
type Stack[T any] struct {
	elements []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{elements: make([]T, 0)}
}

func (s *Stack[T]) Push(element T) {
	s.elements = append(s.elements, element)
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.elements) == 0 {
		return zero, ErrEmptyCollection
	}
	lastIndex := len(s.elements) - 1
	element := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]
	return element, nil
}

func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.elements) == 0 {
		return zero, ErrEmptyCollection
	}
	return s.elements[len(s.elements)-1], nil
}
func (s *Stack[T]) Size() int {
	return len(s.elements)
}
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

// 3. Generic Queue
type Queue[T any] struct {
	elements []T
	head     int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{elements: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(value T) {
	q.elements = append(q.elements, value)
}

func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if q.head >= len(q.elements) {
		return zero, ErrEmptyCollection
	}
	element := q.elements[q.head]
	var zeroVal T
	q.elements[q.head] = zeroVal
	q.head++
	if q.head > len(q.elements)/2 {
		q.elements = append([]T(nil), q.elements[q.head:]...)
		q.head = 0
	}
	return element, nil
}
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if q.head >= len(q.elements) {
		return zero, ErrEmptyCollection
	}
	return q.elements[q.head], nil
}

func (q *Queue[T]) Size() int {
	return len(q.elements) - q.head
}

func (q *Queue[T]) IsEmpty() bool {
	return q.Size() == 0
}

// 4. Generic Set
type Set[T comparable] struct {
	elements map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{elements: make(map[T]struct{})}
}

func (s *Set[T]) Add(element T) {
	if s.elements == nil {
		s.elements = make(map[T]struct{})
	}
	s.elements[element] = struct{}{}
}

func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}

func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

func (s *Set[T]) Size() int {
	return len(s.elements)
}

func (s *Set[T]) Elements() []T {
	result := make([]T, 0, len(s.elements))
	for element := range s.elements {
		result = append(result, element)
	}
	return result
}

func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for element := range s1.elements {
		result.Add(element)
	}
	for element := range s2.elements {
		result.Add(element)
	}
	return result
}

func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for element := range s1.elements {
		if s2.Contains(element) {
			result.Add(element)
		}
	}
	return result
}
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s1.elements {
		if !s2.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// 5. Generic Utility Functions
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filtered := make([]T, 0)
	for _, x := range slice {
		if predicate(x) {
			filtered = append(filtered, x)
		}
	}
	return filtered
}

func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	var s U = initial
	for _, x := range slice {
		s = reducer(s, x)
	}
	return s
}
func Contains[T comparable](slice []T, element T) bool {
	for _, x := range slice {
		if x == element {
			return true
		}
	}
	return false
}

func FindIndex[T comparable](slice []T, element T) int {
	for i, x := range slice {
		if x == element {
			return i
		}
	}
	return -1
}

func RemoveDuplicates[T comparable](slice []T) []T {
	unique := make([]T, 0)
	seen := make(map[T]bool)
	for _, x := range slice {
		if !seen[x] {
			seen[x] = true
			unique = append(unique, x)
		}
	}
	return unique
}
