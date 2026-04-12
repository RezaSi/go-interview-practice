package generics

import "errors"

// ErrEmptyCollection is returned when an operation cannot be performed on an empty collection
var ErrEmptyCollection = errors.New("collection is empty")

//
// 1. Generic Pair
//

// Pair represents a generic pair of values of potentially different types
type Pair[T, U any] struct {
	First  T
	Second U
}

// NewPair creates a new pair with the given values
func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{
		First:  first,
		Second: second,
	}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{
		First:  p.Second,
		Second: p.First,
	}
}

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	elements []T
	size     int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		elements: make([]T, 0),
		size:     0,
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	if s.size == len(s.elements) {
		s.elements = append(s.elements, value)
	} else {
		s.elements[s.size] = value
	}
	s.size++
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if s.size == 0 {
		return zero, errors.New("stack is empty")
	}
	s.size--
	return s.elements[s.size], nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if s.size == 0 {
		return zero, errors.New("stack is empty")
	}
	return s.elements[s.size-1], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return s.size
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return s.size == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	elements []T
	l, r     int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		elements: make([]T, 0),
		l:        0,
		r:        0,
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	if q.r == len(q.elements) {
		q.elements = append(q.elements, value)
	} else {
		q.elements[q.r] = value
	}
	q.r++
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if q.l == q.r {
		return zero, errors.New("queue is empty")
	}
	q.l++
	return q.elements[q.l-1], nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if q.l == q.r {
		return zero, errors.New("queue is empty")
	}
	return q.elements[q.l], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return q.r - q.l
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return q.l == q.r
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	m map[T]bool
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]bool)}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	s.m[value] = true
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	if _, ok := s.m[value]; ok {
		delete(s.m, value)
	}
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	_, ok := s.m[value]
	return ok
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.m)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	elements := make([]T, 0)
	for k := range s.m {
		elements = append(elements, k)
	}
	return elements
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	var union = NewSet[T]()
	for k := range s1.m {
		union.Add(k)
	}
	for k := range s2.m {
		union.Add(k)
	}
	return union
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	intersection := NewSet[T]()
	for k := range s1.m {
		if s2.Contains(k) {
			intersection.Add(k)
		}
	}
	return intersection
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	diff := NewSet[T]()
	for k := range s1.m {
		if !s2.Contains(k) {
			diff.Add(k)
		}
	}
	return diff
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filtered := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	mapped := make([]U, 0)
	for _, v := range slice {
		mapped = append(mapped, mapper(v))
	}
	return mapped
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	m := make(map[T]bool)
	result := make([]T, 0)
	for _, v := range slice {
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}
