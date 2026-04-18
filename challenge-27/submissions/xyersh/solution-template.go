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
	array []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		array: make([]T, 0),
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.array = append(s.array, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var value T

	if len(s.array) == 0 {
		return value, ErrEmptyCollection
	}

	value = s.array[len(s.array)-1]
	s.array = s.array[:len(s.array)-1]

	return value, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var value T
	if len(s.array) == 0 {
		return value, ErrEmptyCollection
	}

	return s.array[len(s.array)-1], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.array)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	if len(s.array) == 0 {
		return true
	}
	return false
}

// ---------------------------
// 3. Generic Queue
// ---------------------------

type node[T any] struct {
	value T
	next  *node[T]
}

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	head *node[T]
	tail *node[T]
	size int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {

	return &Queue[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	newNode := &node[T]{
		value: value,
		next:  nil,
	}

	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}
	q.size++
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var result T
	if q.head == nil {
		return result, ErrEmptyCollection
	} else {
		result = q.head.value
		q.head = q.head.next
	}
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return result, nil

}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var result T
	if q.head == nil {
		return result, ErrEmptyCollection
	} else {
		result = q.head.value
	}
	return result, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {

	return q.size
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	if q.head == nil {
		return true
	}
	return false
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	values map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {

	return &Set[T]{
		values: make(map[T]struct{}),
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	s.values[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.values, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	_, ok := s.values[value]
	return ok
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {

	return len(s.values)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	result := make([]T, 0, len(s.values))
	for key, _ := range s.values {
		result = append(result, key)
	}
	return result
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for key, _ := range s1.values {
		result.Add(key)
	}
	for key, _ := range s2.values {
		result.Add(key)
	}
	return result
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for key, _ := range s1.values {
		if s2.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	result := NewSet[T]()
	for key, _ := range s1.values {
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
	result := make([]T, 0,  len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, 0, len(slice))
	for _, v := range slice {
		result = append(result, mapper(v))
	}
	return result
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	
	for _, v := range slice {
		initial = reducer(initial, v)
	}
	return initial
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
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
