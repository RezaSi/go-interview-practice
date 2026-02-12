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
	return Pair[T, U]{First: first, Second: second}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	// TODO: Implement this method
	return Pair[U, T]{First: p.Second, Second: p.First}
}

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	// TODO: Add necessary fields
	items []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	// TODO: Implement this function
	return &Stack[T]{items: []T{}}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	// TODO: Implement this method
	s.items = append(s.items, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method

	if len(s.items) == 0 {
		var zero T
		return zero, ErrEmptyCollection
	}

	var removed T
	removed = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1] // [5,8,9]
	return removed, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	// TODO: Implement this method
	if len(s.items) == 0 {
		var zero T
		return zero, ErrEmptyCollection
	}

	top := s.items[len(s.items)-1]
	return top, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	// TODO: Implement this method
	return len(s.items)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	// TODO: Implement this method

	if len(s.items) == 0 {
		return true
	}

	return false
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	// TODO: Add necessary fields
	items []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Implement this function
	return &Queue[T]{items: []T{}}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	// TODO: Implement this method
	q.items = append(q.items, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	// TODO: Implement this method
	var zero T

	if len(q.items) == 0 {
		return zero, ErrEmptyCollection
	}

	removed := q.items[0]
	q.items = q.items[1:]

	return removed, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	// TODO: Implement this method
	var zero T

	if len(q.items) == 0 {
		return zero, ErrEmptyCollection
	}

	front := q.items[0]
	return front, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	// TODO: Implement this method
	return len(q.items)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	// TODO: Implement this method
	if len(q.items) == 0 {
		return true
	}
	return false
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	// TODO: Add necessary fields
	items map[T]bool
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	// TODO: Implement this function
	return &Set[T]{items: map[T]bool{}}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	// TODO: Implement this method
	if s.items == nil {
		s.items = make(map[T]bool)
	}
	s.items[value] = true
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	// TODO: Implement this method
	delete(s.items, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	// TODO: Implement this method
	_, exists := s.items[value]
	if exists {
		return true
	}
	return false
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	// TODO: Implement this method
	return len(s.items)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	// TODO: Implement this method

	if len(s.items) == 0 {
		return []T{}
	}

	var res []T
	for k, _ := range s.items {
		res = append(res, k)
	}

	return res
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	s3 := NewSet[T]()

	for k, _ := range s1.items {
		s3.Add(k)
	}

	for k, _ := range s2.items {
		s3.Add(k)
	}

	return s3
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	s3 := NewSet[T]()

	for k, _ := range s1.items {
		if s2.Contains(k) {
			s3.Add(k)
		}
	}

	return s3
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	s3 := NewSet[T]()

	for k, _ := range s1.items {
		if !s2.Contains(k) {
			s3.Add(k)
		}
	}
	return s3
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// TODO: Implement this function

	var res []T
	for _, v := range slice {
		if predicate(v) {
			res = append(res, v)
		}
	}
	return res
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	// TODO: Implement this function

	var res []U
	for _, v := range slice {
		res = append(res, mapper(v))
	}
	return res
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	// TODO: Implement this function

	for _, v := range slice {
		initial = reducer(initial, v)
	}

	return initial
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	// TODO: Implement this function

	for _, e := range slice {
		if e == element {
			return true
		}
	}

	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	// TODO: Implement this function

	for i, e := range slice {
		if e == element {
			return i
		}
	}

	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	// TODO: Implement this function

	m := map[T]bool{}
	var res []T

	for _, e := range slice {
		if m[e] == true {
			continue
		}
		res = append(res, e)
		m[e] = true
	}

	return res
}
