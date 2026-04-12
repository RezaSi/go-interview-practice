package generics

import (
	"errors"
	"fmt"
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
	data []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	// TODO: Implement this function
	return &Stack[T]{data: make([]T, 0)}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.data = append(s.data, value)
	// TODO: Implement this method
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method
	var zero T

	if s.IsEmpty() {
		return zero, fmt.Errorf("stack is empty")
	}

	lastIndex := s.Size() - 1
	lastValue := s.data[lastIndex]
	s.data = s.data[:lastIndex]

	return lastValue, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	// TODO: Implement this method
	var zero T

	if s.IsEmpty() {
		return zero, fmt.Errorf("stack is empty")
	}

	return s.data[s.Size()-1], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	// TODO: Implement this method
	return len(s.data)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	// TODO: Implement this method
	return s.Size() == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type Queue[T any] struct {
	// TODO: Add necessary fields
	Head *Node[T]
	Len  int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Implement this function
	return &Queue[T]{}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	// TODO: Implement this method

	newNode := &Node[T]{Value: value}
	if q.Head == nil {
		q.Head = newNode
	} else {
		node := q.Head
		for node.Next != nil {
			node = node.Next
		}
		node.Next = newNode
	}
	q.Len++
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	// TODO: Implement this method
	var zero T

	if q.IsEmpty() {
		return zero, fmt.Errorf("the queue is empty")
	}

	node := q.Head
	nodeNext := q.Head.Next
	q.Head = nodeNext

	q.Len--

	return node.Value, nil
}

func (q *Queue[T]) Print() {

	if q.IsEmpty() {
		fmt.Println("queue is empty")
	}

	nodeValue := q.Head
	index := 0
	for nodeValue != nil {
		fmt.Printf("%v: %v\n", index, nodeValue.Value)
		nodeValue = nodeValue.Next
		index++
	}
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	// TODO: Implement this method
	var zero T

	if q.IsEmpty() {
		return zero, fmt.Errorf("the queue is empty")
	}

	return q.Head.Value, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	// TODO: Implement this method
	return q.Len
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	// TODO: Implement this method
	return q.Size() == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	// TODO: Add necessary fields
	data map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	// TODO: Implement this function
	return &Set[T]{data: make(map[T]struct{})}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	// TODO: Implement this method
	s.data[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	// TODO: Implement this method
	delete(s.data, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	// TODO: Implement this method
	if _, ok := s.data[value]; !ok {
		return false
	}
	return true
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	// TODO: Implement this method
	return len(s.data)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	// TODO: Implement this method
	result := make([]T, 0, s.Size())

	for val := range s.data {
		result = append(result, val)
	}

	return result
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	result := NewSet[T]()

	for val := range s1.data {
		result.Add(val)
	}

	for val := range s2.data {
		result.Add(val)
	}

	return result
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	result := NewSet[T]()

	for val := range s1.data {
		if s2.Contains(val) {
			result.Add(val)
		}
	}

	return result
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function

	result := NewSet[T]()

	for val := range s1.data {
		if !s2.Contains(val) {
			result.Add(val)
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

	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	// TODO: Implement this function
	result := make([]U, 0, len(slice))

	for _, v := range slice {
		result = append(result, mapper(v))
	}

	return result
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	// TODO: Implement this function

	result := initial

	for _, v := range slice {
		result = reducer(result, v)
	}

	return result
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	// TODO: Implement this function

	for _, v := range slice {
		if v == element {
			return true
		}
	}

	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	// TODO: Implement this function

	for i, v := range slice {
		if v == element {
			return i
		}
	}

	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	// TODO: Implement this function

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
