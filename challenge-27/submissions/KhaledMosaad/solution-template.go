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
	data    []T
	size    int
	pointer int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	// TODO: Implement this function
	return &Stack[T]{
		data:    []T{},
		size:    0,
		pointer: -1,
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	// TODO: Implement this method
	s.data = append(s.data, value)
	s.size++
	s.pointer++
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method
	if s.pointer == -1 {
		return *new(T), ErrEmptyCollection
	}
	value := s.data[s.pointer]
	s.pointer--
	s.size--
	return value, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	if s.pointer == -1 {
		return *new(T), ErrEmptyCollection
	}
	value := s.data[s.pointer]
	return value, nil
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
	// TODO: Add necessary fields

	size    int
	data    []T
	pointer int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Implement this function
	return &Queue[T]{
		size:    0,
		pointer: 0,
		data:    []T{},
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	// TODO: Implement this method
	q.data = append(q.data, value)
	q.size++
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	// TODO: Implement this method

	var value T
	if q.size == 0 {
		return value, ErrEmptyCollection
	}

	value = q.data[q.pointer]
	q.pointer++
	q.size--
	return value, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var value T
	if q.size == 0 {
		return value, ErrEmptyCollection
	}

	value = q.data[q.pointer]
	return value, nil
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
	data map[T]bool
	size int
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	// TODO: Implement this function
	return &Set[T]{
		data: map[T]bool{},
		size: 0,
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	// TODO: Implement this method
	if _, ok := s.data[value]; !ok {
		s.data[value] = true
		s.size++
	}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	// TODO: Implement this method
	if _, ok := s.data[value]; ok {
		delete(s.data, value)
		s.size--
	}
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	if _, ok := s.data[value]; ok {
		return true
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

	var result []T
	for key := range s.data {
		result = append(result, key)
	}
	return result
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	result := NewSet[T]()
	for key := range s1.data {
		result.Add(key)
	}

	for key := range s2.data {
		result.Add(key)
	}
	return result
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	union := Union(s1, s2)
	result := NewSet[T]()
	for key := range union.data {
		if s1.Contains(key) && s2.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	result := NewSet[T]()
	for key := range s1.data {
		if !s2.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// TODO: Implement this function

	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	// TODO: Implement this function

	var result []U
	for _, item := range slice {
		result = append(result, mapper(item))
	}
	return result
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	// TODO: Implement this function

	var result U = initial
	for _, item := range slice {
		result = reducer(result, item)
	}

	return result
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
	for i, item := range slice {
		if item == element {
			return i
		}
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	// TODO: Implement this function

	uniq := make(map[T]bool)
	for _, item := range slice {
		uniq[item] = true
	}
	var result []T

	for _, item := range slice {
		if ok := uniq[item]; ok {
			result = append(result, item)
			uniq[item] = false
		}
	}
	return result
}
