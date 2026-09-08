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
	storage []T
	top     int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	// TODO: Implement this function
	return &Stack[T]{storage: []T{}, top: 0}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	// TODO: Implement this method
	s.storage = append(s.storage, value)
	s.top++
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method
	var zero T

	if s.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	s.top--

	element := s.storage[s.top]

	s.storage[s.top] = zero

	return element, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	// TODO: Implement this method
	var zero T

	if s.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	element := s.storage[s.top-1]

	return element, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	// TODO: Implement this method
	return s.top
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	// TODO: Implement this method
	return s.top == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	// TODO: Add necessary fields
	storage []T
	head    int
	tail    int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Implement this function
	return &Queue[T]{
		storage: []T{},
		head:    0,
		tail:    0,
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	// TODO: Implement this method
	q.storage = append(q.storage, value)
	q.tail++
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	// TODO: Implement this method
	var zero T
	if q.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	value := q.storage[q.head]
	q.head++

	return value, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	// TODO: Implement this method
	var zero T

	if q.IsEmpty() {
		return zero, ErrEmptyCollection
	}

	value := q.storage[q.head]

	return value, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	// TODO: Implement this method
	return q.tail - q.head
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	// TODO: Implement this method
	return q.tail == q.head
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
	return &Set[T]{items: make(map[T]bool)}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	// TODO: Implement this method
	s.items[value] = true
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	// TODO: Implement this method
	if _, ok := s.items[value]; ok {
		delete(s.items, value)
	}
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	return s.items[value]
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	// TODO: Implement this method
	return len(s.items)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	// TODO: Implement this method
	elements := []T{}
	for keys := range s.items {
		elements = append(elements, keys)
	}
	return elements
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	newSet := NewSet[T]()

	for k := range s1.items {
		newSet.Add(k)
	}

	for k := range s2.items {
		newSet.Add(k)
	}
	return newSet
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	newSet := NewSet[T]()

	for k1 := range s1.items {
		if _, exists := s2.items[k1]; exists {
			newSet.Add(k1)
		}
	}

	return newSet
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	newSet := NewSet[T]()

	for k1 := range s1.items {
		if _, exists := s2.items[k1]; !exists {
			newSet.Add(k1)
		}
	}

	return newSet
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// TODO: Implement this function
	filtered := []T{}

	for _, item := range slice {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	// TODO: Implement this function

	mapped := make([]U, len(slice))

	for i, item := range slice {
		mapped[i] = mapper(item)
	}

	return mapped
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	// TODO: Implement this function

	for _, item := range slice {
		initial = reducer(initial, item)
	}

	return initial
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	// TODO: Implement this function
	contains := false

	for _, item := range slice {
		if item == element {
			contains = true
		}
	}

	return contains
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	// TODO: Implement this function
	idx := -1

	for i, item := range slice {
		if item == element {
			idx = i
		}
	}

	return idx
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	// TODO: Implement this function
	items := make(map[T]bool)

	newSlice := []T{}

	for _, item := range slice {
		if !items[item] {
			newSlice = append(newSlice, item)
			items[item] = true
		}
	}

	return newSlice
}
