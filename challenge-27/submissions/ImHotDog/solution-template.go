package generics

import (
	"errors"
	"slices"
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
	items []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	top, err := s.Peek()
	if err != nil {
		return top, ErrEmptyCollection
	}

	s.items = s.items[:len(s.items)-1]

	return top, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T

	if s.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	topIdx := len(s.items) - 1
	top := s.items[topIdx]

	return top, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	if s.Size() == 0 {
		return true
	}

	return false
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	items []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.items = append(q.items, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	front, err := q.Front()
	if err != nil {
		return front, ErrEmptyCollection
	}

	q.items = q.items[1:]

	return front, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T

	if q.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	return q.items[0], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	if q.Size() == 0 {
		return true
	}

	return false
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
	_, contains := s.elements[value]

	return contains
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	resultElements := make([]T, 0, s.Size())

	for k := range s.elements {
		resultElements = append(resultElements, k)
	}

	return resultElements
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()

	for k := range s1.elements {
		result.Add(k)
	}

	for k := range s2.elements {
		result.Add(k)
	}

	return result
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()

	var smaller *Set[T]
	var larger *Set[T]

	if s1.Size() < s2.Size() {
		smaller = s1
		larger = s2
	} else {
		smaller = s2
		larger = s1
	}

	for k := range smaller.elements {
		if larger.Contains(k) {
			result.elements[k] = struct{}{}
		}
	}

	return result
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	
	for k := range s1.elements {
		if !s2.Contains(k) {
			result.Add(k)
		}
	}
	
	return result
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, 0, len(slice))

	for _, item := range slice {
		result = append(result, mapper(item))
	}

	return result
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial

	for _, item := range slice {
		result = reducer(result, item)
	}

	return result
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	return slices.Contains(slice, element)
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	return slices.Index(slice, element)
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	seen := NewSet[T]()
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if seen.Contains(item) {
			continue
		}

		seen.Add(item)
		result = append(result, item)
	}

	return result
}
