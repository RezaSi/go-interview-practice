package generics

import (
	"errors"
	"sync"
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
	return Pair[T, U]{First: first, Second: second}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{First: p.Second, Second: p.First}
}

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	Items []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{Items: []T{}}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.Items = append(s.Items, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T

	if s.IsEmpty() {
		return zero, ErrEmptyCollection
	}
	length := len(s.Items)
	zero = s.Items[length-1]
	s.Items = s.Items[:length-1]
	return zero, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T

	if s.IsEmpty() {
		return zero, ErrEmptyCollection

	}

	zero = s.Items[len(s.Items)-1]
	return zero, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.Items)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return s.Size() == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	Items []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		Items: []T{},
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.Items = append(q.Items, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if q.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	zero = q.Items[0]
	q.Items = q.Items[1:]
	return zero, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if q.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	zero = q.Items[0]
	return zero, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.Items)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return q.Size() == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	ItemSet map[T]struct{}
	mu      sync.Mutex
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		ItemSet: map[T]struct{}{},
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ItemSet[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ItemSet, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.ItemSet[value]
	return ok
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.ItemSet)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {

	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]T, len(s.ItemSet))
	i := 0
	for key := range s.ItemSet {
		items[i] = key
		i += 1
	}
	return items
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {

	nw_set := NewSet[T]()

	for _, val := range s1.Elements() {
		nw_set.Add(val)
	}
	for _, val := range s2.Elements() {
		nw_set.Add(val)
	}
	return nw_set
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	intr_set := NewSet[T]()
	for _, val := range s1.Elements() {
		if !s2.Contains(val) {
			continue
		}

		intr_set.Add(val)
	}
	return intr_set
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	diff_set := NewSet[T]()
	for _, val := range s1.Elements() {
		if s2.Contains(val) {
			continue
		}

		diff_set.Add(val)
	}
	return diff_set
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filtered := []T{}
	for _, val := range slice {
		if !predicate(val) {
			continue
		}
		filtered = append(filtered, val)
	}

	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	mapped_vals := make([]U, len(slice))

	for index, val := range slice {
		mapped_vals[index] = mapper(val)
	}

	return mapped_vals
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	for _, val := range slice {
		initial = reducer(initial, val)
	}
	return initial
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	for _, val := range slice {
		if val == element {
			return true
		}
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for index, val := range slice {
		if val == element {
			return index
		}
	}

	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {

	item_set := NewSet[T]()
	uniques := []T{}

	for _, val := range slice {

		if item_set.Contains(val) {
			continue
		}

		item_set.Add(val)
		uniques = append(uniques, val)
	}
	return uniques
}
