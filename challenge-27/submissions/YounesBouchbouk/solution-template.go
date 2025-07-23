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
	// TODO: Implement this function
	return Pair[T, U]{
		First:  first,
		Second: second,
	}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	// TODO: Implement this method
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
	// TODO: Add necessary fields
	items []T
	size  int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	// TODO: Implement this function
	return &Stack[T]{
		items: make([]T, 0),
		size:  0,
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
	s.size++
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method
	if s.size == 0 {
		var zero T
		return zero, ErrEmptyCollection
	}
	top := s.items[s.size-1]
	s.items = s.items[:s.size-1]
	s.size--
	return top, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	// TODO: Implement this method
	if s.size == 0 {
		var zero T
		return zero, ErrEmptyCollection
	}
	return s.items[s.size-1], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	// TODO: Implement this method
	return s.size
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	// TODO: Implement this method
	return s.size == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	// TODO: Add necessary fields
	items []T
	size  int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Implement this function
	return &Queue[T]{
		items: make([]T, 0),
		size:  0,
	}

}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.items = append(q.items, value)
	q.size++
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	// TODO: Implement this method
	var zero T
	if q.size == 0 {
		return zero, ErrEmptyCollection
	}

	firstElement := q.items[0]
	q.items = q.items[1:]
	q.size--

	return firstElement, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	// TODO: Implement this method
	var zero T

	if q.size == 0 {
		return zero, ErrEmptyCollection
	}
	return q.items[0], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	// TODO: Implement this method
	return q.size
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	// TODO: Implement this method
	return q.size == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	// TODO: Add necessary fields
	items []T
	size  int
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	// TODO: Implement this function
	return &Set[T]{
		items: make([]T, 0),
		size:  0,
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	if !s.Contains(value) {
		s.items = append(s.items, value)
		s.size++
	}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {

	if s.size == 0 {
		return
	}

	if s.Contains(value) {
		for index, item := range s.items {
			if item == value {
				if index == 0 {
					s.items = s.items[1:]
				}

				if index == s.size {
					s.items = s.items[:s.size-1]
				}

				s.items = append(s.items[:index], s.items[index+1:]...)
				s.size--
				break
			}
		}
	}
	// TODO: Implement this method
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	// TODO: Implement this method
	for _, item := range s.items {
		if item == value {
			return true
		}
	}
	return false
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	// TODO: Implement this method
	return s.size
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	// TODO: Implement this method
	return s.items
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	unionSet := NewSet[T]()
	for _, item := range s1.items {
		unionSet.Add(item)
	}
	for _, item := range s2.items {
		unionSet.Add(item)
	}
	return unionSet
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	intersectionSet := NewSet[T]()
	for _, item := range s1.items {
		if s2.Contains(item) {
			intersectionSet.Add(item)
		}
	}
	return intersectionSet
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	differenceSet := NewSet[T]()
	for _, item := range s1.items {
		if !s2.Contains(item) {
			differenceSet.Add(item)
		}
	}
	return differenceSet
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// TODO: Implement this function
	newSlice := make([]T, 0)

	for _, item := range slice {
		if predicate(item) {
			newSlice = append(newSlice, item)
		}
	}
	return newSlice
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	// TODO: Implement this function
	newSlice := make([]U, 0)

	for _, item := range slice {
		newSlice = append(newSlice, mapper(item))
	}
	return newSlice
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	// TODO: Implement this function
	reducedValue := initial

	for _, item := range slice {
		reducedValue = reducer(reducedValue, item)
	}
	return reducedValue
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	// TODO: Implement this function
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	// TODO: Implement this function
	for index, item := range slice {
		if item == element {
			return index
		}
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	// TODO: Implement this function
	newSlice := make([]T, 0)

	for _, item := range slice {
		if !Contains(newSlice, item) {
			newSlice = append(newSlice, item)
		}
	}
	return newSlice
}
