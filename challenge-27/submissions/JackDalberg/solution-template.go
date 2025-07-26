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

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	Array []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{[]T{}}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.Array = append(s.Array, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	length := len(s.Array)
	if length == 0 {
	    return zero, errors.New("Cannot pop from empty stack")
	}
	zero = s.Array[length-1]
	s.Array = s.Array[:length-1]
	return zero, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	length := len(s.Array)
	if length == 0 {
	    return zero, errors.New("Cannot peek empty stack")
	}
	zero = s.Array[length-1]
	return zero, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.Array)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	if s.Size() == 0{
	    return true
	}
	return false
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	Array []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{[]T{}}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.Array = append(q.Array, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	length := len(q.Array)
	if length == 0 {
	    return zero, errors.New("Cannot dequeue from empty queue")
	}
	zero = q.Array[0]
	if length == 1{
	    q.Array = []T{}
	}else{
	    q.Array = append([]T{}, q.Array[1:]...)
	}
	return zero, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	length := len(q.Array)
	if length == 0 {
	    return zero, errors.New("Cannot dequeue from empty queue")
	}
	zero = q.Array[0]
	return zero, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.Array)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	if q.Size() == 0{
	    return true
	}
	return false
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	Map map[T]bool // unimportant what it maps to
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{map[T]bool{}}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	s.Map[value] = true
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.Map, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	return s.Map[value]
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.Map)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
    var keys []T
    for key, _ := range s.Map {
        keys = append(keys, key)
    }
	return keys
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	for _, element := range append(s1.Elements(), s2.Elements()...) {
	    newSet.Add(element)
	}
	return newSet
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	for _, element := range s1.Elements() {
	    if s2.Contains(element){
	        newSet.Add(element)
	    }
	}
	return newSet
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	for _, element := range s1.Elements() {
	    if !s2.Contains(element){
	        newSet.Add(element)
	    }
	}
	return newSet
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filtered := []T{}
	for _, val := range slice {
	    if predicate(val){
	        filtered = append(filtered, val)
	    }
	}
	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	mapped := []U{}
	for _, val := range slice {
	    mapped = append(mapped, mapper(val))
	}
	return mapped
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
	for _, val :=range slice {
	    if val == element{
	        return true
	    }
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for idx, val := range slice{
	    if val == element {
	        return idx
	    }
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	noDupes := []T{}
	for _, val := range slice {
	    if !Contains(noDupes, val){
	        noDupes = append(noDupes, val)
	    }
	}
	return noDupes
}
