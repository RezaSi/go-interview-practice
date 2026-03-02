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
	return Pair[T, U]{first, second}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	return NewPair(p.Second, p.First)
}

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	Data []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	if s == nil {
		return
	}
	s.Data = append(s.Data, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.Data) == 0 {
		return zero, ErrEmptyCollection
	}
	res := s.Data[len(s.Data)-1]
	s.Data[len(s.Data)-1] = zero
	s.Data = s.Data[:len(s.Data)-1]
	return res, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.Data) == 0 {
		return zero, ErrEmptyCollection
	}
	zero = s.Data[len(s.Data)-1]
	return zero, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	if s == nil {
		return 0
	}
	return len(s.Data)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	if s == nil {
		return false
	}
	return len(s.Data) == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	Data []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	if q == nil {
		return
	}
	q.Data = append(q.Data, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	var res T
	if len(q.Data) == 0 {
		return zero, ErrEmptyCollection
	}
	res = q.Data[0]
	q.Data[0] = zero
	q.Data = q.Data[1:]
	return res, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	var res T
	if len(q.Data) == 0 {
		return zero, ErrEmptyCollection
	}
	res = q.Data[0]
	return res, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	if q == nil {
		return 0
	}
	return len(q.Data)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	if q == nil {
		return false
	}
	return len(q.Data) == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	Data map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		Data: make(map[T]struct{}),
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	if s == nil {
		return
	}
	if s.Data == nil {
		s.Data = make(map[T]struct{})
	}
	s.Data[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	if s == nil {
		return
	}
	delete(s.Data, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	if s == nil {
		return false
	}
	_, res := s.Data[value]
	return res
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	if s == nil {
		return 0
	}
	return len(s.Data)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	if s == nil {
		return nil
	}
	var res []T
	for e := range s.Data {
		res = append(res, e)
	}
	return res
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	if s1 != nil {
		for k, v := range s1.Data {
			res.Data[k] = v
		}
	}
	if s2 != nil {
		for k, v := range s2.Data {
			res.Data[k] = v
		}
	}
	return res
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	if s1 == nil || s2 == nil {
		return res
	}
	for k := range s1.Data {
		if _, ok := s2.Data[k]; ok {
			res.Data[k] = struct{}{}
		}
	}
	return res
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	if s1 == nil {
		return res
	}
	for k := range s1.Data {
		if s2 == nil {
			res.Data[k] = struct{}{}
			continue
		}
		if _, ok := s2.Data[k]; !ok {
			res.Data[k] = struct{}{}
		}
	}
	return res
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	res := []T{}
	for _, t := range slice {
		if predicate(t) {
			res = append(res, t)
		}
	}
	return res
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	res := []U{}
	for _, t := range slice {
		res = append(res, mapper(t))
	}
	return res
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	res := initial
	for _, s := range slice {
		res = reducer(res, s)
	}
	return res
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	return slices.Contains(slice, element)
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	return slices.IndexFunc(slice, func(t T) bool {
		if t == element {
			return true
		}
		return false
	})
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	tmp := NewSet[T]()
	res := []T{}
	for _, t := range slice {
		if _, ok := tmp.Data[t]; !ok {
			tmp.Data[t] = struct{}{}
			res = append(res, t)
		}
	}
	return res
}
