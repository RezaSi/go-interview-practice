package generics

import (
	"errors"
)

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
	return Pair[T, U]{first, second}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{p.Second, p.First}
}

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	values []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	size := len(s.values)
	if size == 0 {
		return zero, ErrEmptyCollection
	}
	val := s.values[size-1]
	s.values[size-1] = zero
	s.values = s.values[:size-1]
	return val, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	size := len(s.values)
	if size == 0 {
		return zero, ErrEmptyCollection
	}
	return s.values[size-1], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.values)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	values []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.values = append(q.values, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	size := len(q.values)
	if size == 0 {
		return zero, ErrEmptyCollection
	}
	val := q.values[0]
	q.values[0] = zero
	q.values = q.values[1:]
	return val, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	size := len(q.values)
	if size == 0 {
		return zero, ErrEmptyCollection
	}
	return q.values[0], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.values)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return len(q.values) == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}),
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	s.elements[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.elements, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	_, exists := s.elements[value]
	return exists
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	elems := make([]T, 0, len(s.elements))
	for k := range s.elements {
		elems = append(elems, k)
	}
	return elems
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	for k := range s1.elements {
		newSet.Add(k)
	}
	for k := range s2.elements {
		newSet.Add(k)
	}
	return newSet
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	// Optimize: Iterate over the smaller set first to minimize contains checks
	smaller, larger := s1, s2
	if s1.Size() > s2.Size() {
		smaller, larger = s2, s1
	}
	for k := range smaller.elements {
		if larger.Contains(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	for k := range s1.elements {
		if !s2.Contains(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	newSlice := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	newMap := make([]U, len(slice))
	for i, v := range slice {
		newMap[i] = mapper(v)
	}
	return newMap
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	newValue := initial
	for _, v := range slice {
		newValue = reducer(newValue, v)
	}
	return newValue
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
	seen := make(map[T]struct{}, len(slice))
	newSlice := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
