package generics

import (
	"cmp"
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
	Values []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		Values: []T{},
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.Values = append(s.Values, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T

	if s.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	l := len(s.Values)
	v := s.Values[l-1]

	s.Values = slices.Delete(s.Values, l-1, l)

	return v, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T

	if s.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	return s.Values[len(s.Values)-1], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.Values)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return len(s.Values) == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	Values []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		Values: []T{},
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.Values = append(q.Values, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T

	if q.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	v := q.Values[0]
	q.Values = slices.Delete(q.Values, 0, 1)

	return v, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T

	if q.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	return q.Values[0], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.Values)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return len(q.Values) == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	Values []T
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		Values: []T{},
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	if !s.Contains(value) {
		s.Values = append(s.Values, value)
	}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	i := slices.Index(s.Values, value)
	if i >= 0 {
		s.Values = slices.Delete(s.Values, i, i+1)
	}
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	return slices.Contains(s.Values, value)
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.Values)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	return slices.Clone(s.Values)
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// allocate result slice with enough capacity
	res := make([]T, 0, s1.Size()+s2.Size())

	// fill result with contents of first slice
	res = append(res, s1.Values...)

	// append elements from 2nd slice if not already members of first slice
	for _, v := range s2.Values {
		if !s1.Contains(v) {
			res = append(res, v)
		}
	}

	return &Set[T]{
		Values: res,
	}
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	// init empty result slice
	res := []T{}

	// iter thru s1 capturing values in s2
	for _, v := range s1.Values {
		if slices.Contains(s2.Values, v) {
			res = append(res, v)
		}
	}

	return &Set[T]{
		Values: res,
	}
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	// init empty result slice
	res := []T{}

	// iter thru s1 capturing values not in s2
	for _, v := range s1.Values {
		if !slices.Contains(s2.Values, v) {
			res = append(res, v)
		}
	}

	return &Set[T]{
		Values: res,
	}
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	res := make([]T, 0, len(slice))

	for _, v := range slice {
		if predicate(v) {
			res = append(res, v)
		}
	}

	return res
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	res := make([]U, 0, len(slice))

	for _, v := range slice {
		res = append(res, mapper(v))
	}

	return res
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	res := initial

	for _, v := range slice {
		res = reducer(res, v)
	}

	return res
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
func RemoveDuplicates[T cmp.Ordered](slice []T) []T {
	res := slices.Clone(slice)
	slices.Sort(res)
	res = slices.Compact(res)

	return res
}
