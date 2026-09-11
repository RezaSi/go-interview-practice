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
	return Pair[T, U]{first, second}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{p.Second, p.First}
}

// 2. Generic Stack
var InitCap = 8

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	arr []T
	top int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{make([]T, InitCap), -1}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	if s.top+1 >= len(s.arr) {
		newArr := make([]T, len(s.arr)*2)
		copy(newArr, s.arr)
		s.arr = newArr
	}
	s.top++
	s.arr[s.top] = value
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if s.top < 0 {
		return zero, ErrEmptyCollection
	}
	value := s.arr[s.top]
	s.arr[s.top] = zero
	s.top--
	return value, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if s.top < 0 {
		return zero, ErrEmptyCollection
	}
	return s.arr[s.top], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return s.top + 1
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return s.top < 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	arr   []T
	front int
	back  int
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{make([]T, InitCap), 0, -1}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	if q.back+1 >= len(q.arr) {
		newArr := make([]T, len(q.arr)*2)
		copy(newArr, q.arr[q.front:q.back+1])
		q.arr = newArr
	}
	q.back++
	q.arr[q.back] = value
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if q.front > q.back {
		return zero, ErrEmptyCollection
	}
	value := q.arr[q.front]
	q.arr[q.front] = zero
	q.front++
	return value, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if q.front > q.back {
		return zero, ErrEmptyCollection
	}
	return q.arr[q.front], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return q.back - q.front + 1
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return q.front > q.back
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	ma map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{make(map[T]struct{}, InitCap)}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	//if _, ok := s.ma[value]; ok {
	//	return
	//}
	s.ma[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.ma, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	_, ok := s.ma[value]
	return ok
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.ma)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	arr := make([]T, 0, len(s.ma))
	for k := range s.ma {
		arr = append(arr, k)
	}
	return arr
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	union := NewSet[T]()
	for k := range s1.ma {
		union.Add(k)
	}
	for k := range s2.ma {
		union.Add(k)
	}
	return union
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	intersection := NewSet[T]()
	var smallSet, bigSet *Set[T]
	if s1.Size() <= s2.Size() {
		smallSet = s1
		bigSet = s2
	} else {
		smallSet = s2
		bigSet = s1
	}
	for k := range smallSet.ma {
		if bigSet.Contains(k) {
			intersection.Add(k)
		}
	}
	return intersection
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	difference := NewSet[T]()
	for k := range s1.ma {
		if !s2.Contains(k) {
			difference.Add(k)
		}
	}
	return difference
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filtered := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	mapped := make([]U, 0, len(slice))
	for _, v := range slice {
		mapped = append(mapped, mapper(v))
	}
	return mapped
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	reduced := initial
	for _, v := range slice {
		reduced = reducer(reduced, v)
	}
	return reduced
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
	arr := make([]T, 0, len(slice))
	showed := NewSet[T]()
	for _, v := range slice {
		if !showed.Contains(v) {
			arr = append(arr, v)
		}
		showed.Add(v)
	}
	return arr
}
