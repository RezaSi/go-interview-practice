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
	M  []T
	li int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		M:  []T{},
		li: -1,
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.li += 1
	s.M = append(s.M, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.M) == 0 {
		return zero, ErrEmptyCollection
	}
	value := s.M[s.li]
	s.M = s.M[0:s.li]
	if s.li >= 0 {
		s.li -= 1
	}
	return value, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.M) == 0 {
		return zero, ErrEmptyCollection
	}

	return s.M[s.li], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.M)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return len(s.M) == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	Q []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		Q: []T{},
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.Q = append(q.Q, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T

	if len(q.Q) == 0 {
		return zero, ErrEmptyCollection
	}

	value := q.Q[0]
	q.Q = q.Q[1:]

	return value, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if len(q.Q) == 0 {
		return zero, ErrEmptyCollection
	}
	value := q.Q[0]

	return value, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.Q)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return len(q.Q) == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	S map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		S: map[T]struct{}{},
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	s.S[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.S, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	_, ok := s.S[value]
	return ok
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.S)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	result := make([]T, 0, len(s.S))
	for k, _ := range s.S {
		result = append(result, k)
	}
	return result
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	union := map[T]struct{}{}
	for k, _ := range s1.S {
		union[k] = struct{}{}
	}
	for k, _ := range s2.S {
		union[k] = struct{}{}
	}

	return &Set[T]{
		S: union,
	}
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	isec := map[T]struct{}{}

	for k, _ := range s1.S {
		if _, ok := s2.S[k]; ok {
			isec[k] = struct{}{}
		}
	}

	return &Set[T]{S: isec}
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	diff := map[T]struct{}{}

	for k, _ := range s1.S {
		if _, ok := s2.S[k]; !ok {
			diff[k] = struct{}{}
		}
	}
	/*
		for k, _ := range s2.S {
			if _, ok := s1.S[k]; !ok {
				diff[k] = struct{}{}
			}
		}
	*/

	return &Set[T]{S: diff}
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	res := []T{}
	for _, v := range slice {
		if predicate(v) {
			res = append(res, v)
		}
	}
	return res
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	res := []U{}
	for _, v := range slice {
		res = append(res, mapper(v))
	}
	return res
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	var result U

	result = initial
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
	exists := map[T]struct{}{}
	result := []T{}
	for _, v := range slice {
		if _, ok := exists[v]; !ok {
			result = append(result, v)
		}
		exists[v] = struct{}{}
	}

	return result
}
