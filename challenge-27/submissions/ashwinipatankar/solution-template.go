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
	return Pair[T, U]{first, second}
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
	head, tail *Node[T]
}

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	// TODO: Implement this function
	return &Stack[T]{}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	// TODO: Implement this method
	node := &Node[T]{Value: value, Next: s.head}

	s.head = node
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method

	var zero T
	if s.head == nil && s.tail == nil {
		return zero, ErrEmptyCollection
	}

	zero = s.head.Value
	s.head = s.head.Next

	return zero, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	
	var zero T
	if s.head == nil && s.tail == nil {
		return zero, ErrEmptyCollection
	}

	zero = s.head.Value

	return zero, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	// TODO: Implement this method
	var totalElements int

	for node := s.head; node != nil; node = node.Next {
		totalElements++
	}

	return totalElements
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	// TODO: Implement this method

	if s.head == nil && s.tail == nil {
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
	head, tail *Node[T]
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Implement this function
	return &Queue[T]{}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	// TODO: Implement this method
	n := &Node[T]{value, nil}
	if q.head == nil {
		q.head = n
		q.tail = n
	} else {
		q.tail.Next = n
		q.tail = n
	}
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	// TODO: Implement this method
	var zero T
	if q.head == nil {
		return zero, ErrEmptyCollection
	}

	zero = q.head.Value
	q.head = q.head.Next

	return zero, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	// TODO: Implement this method
	var zero T
	if q.head == nil {
		return zero, ErrEmptyCollection
	}

	zero = q.head.Value

	return zero, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	// TODO: Implement this method
	qSize := 0
	for node := q.head; node != nil; node = node.Next {
		qSize++
	}

	return qSize
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	// TODO: Implement this method
	if q.head == nil && q.tail == nil || q.Size() == 0{
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
	data map[T]bool
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	// TODO: Implement this function
	return &Set[T]{data: make(map[T]bool)}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	// TODO: Implement this method
	s.data[value] = true
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	// TODO: Implement this method
	delete(s.data, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	// TODO: Implement this method
	if _, ok := s.data[value]; ok {
		return true
	}

	return false
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	// TODO: Implement this method
	return len(s.data)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	// TODO: Implement this method
	data := []T{}

	for key := range s.data {
		data = append(data, key)
	}

	return data
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	newSet := NewSet[T]()

	for _, value := range s1.Elements() {
		newSet.Add(value)
	}

	for _, value := range s2.Elements() {
		newSet.Add(value)
	}

	return newSet
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	newSet := NewSet[T]()

	for _, value := range s1.Elements() {
		if s2.Contains(value) {
			newSet.Add(value)
		}
	}

	return newSet
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	newSet := NewSet[T]()

	for _, value := range s1.Elements() {
		if !s2.Contains(value) {
			newSet.Add(value)
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
	s := []T{}
	for i := range slice {
		if predicate(slice[i]) {
			s = append(s, slice[i])
		}
	}

	return s
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	// TODO: Implement this function
	s := []U{}
	for i := range slice {
		s = append(s, mapper(slice[i]))
	}

	return s
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	// TODO: Implement this function
	for i := range slice {
		initial = reducer(initial, slice[i])
	}

	return initial
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	// TODO: Implement this function
	for i := range slice {
		if slice[i] == element {
			return true
		}
	}

	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	// TODO: Implement this function
	for i := range slice {
		if slice[i] == element {
			return i
		}
	}

	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	// TODO: Implement this function
	s := []T{}
	sm := make(map[T]bool, len(slice))
	for i := range slice {
		if _, ok := sm[slice[i]]; ok {
			continue
		} else {
			sm[slice[i]] = true
			s = append(s, slice[i])
		}
	}

	return s
}

