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
	    First: first,
	    Second: second,
	}
}

// Swap returns a new pair with the elements swapped
func (p Pair[T, U]) Swap() Pair[U, T] {
	return Pair[U, T]{
	    First: p.Second,
	    Second: p.First,
	}
}

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	storage []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
	    storage: make([]T, 0),
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.storage = append(s.storage, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method
	var zero T
	if len(s.storage) == 0 {
	    return zero, ErrEmptyCollection
	}
	
	zero = s.storage[len(s.storage) - 1]
	s.storage = s.storage[:len(s.storage)-1]
	
	return zero, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.storage) == 0 {
	    return zero, ErrEmptyCollection
	}
	
	zero = s.storage[len(s.storage) - 1]
	
	return zero, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.storage)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return len(s.storage) == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	storage []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
	    storage: make([]T, 0),
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.storage = append(q.storage, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if len(q.storage) == 0 {
	    return zero, ErrEmptyCollection
	}
	
	zero = q.storage[0]
	q.storage = q.storage[1:]
	
	return zero, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	var zero T
	if len(q.storage) == 0 {
	    return zero, ErrEmptyCollection
	}
	
	zero = q.storage[0]
	
	return zero, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.storage)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return len(q.storage) == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	storage []T
	contains map[T]bool
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
	    storage: make([]T, 0, 1000),
	    contains: make(map[T]bool),
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	if !s.contains[value] {
	    s.storage = append(s.storage, value)
	    s.contains[value] = true
	}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	if !s.contains[value] {
        return
    }
    
    delete(s.contains, value)
    
    newStorage := make([]T, 0, len(s.storage)-1)
    for _, v := range s.storage {
        if v != value {
            newStorage = append(newStorage, v)
        }
    }
    s.storage = newStorage
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	return s.contains[value]
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.storage)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
    out := make([]T, len(s.storage))
    copy(out, s.storage)
    return out
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	
	for _, val := range s1.Elements() {
	    newSet.Add(val)
	}
	
	for _, val := range s2.Elements() {
	    newSet.Add(val)
	}
	
	return newSet
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	
	for _, val := range s1.Elements() {
	    if s2.Contains(val) {
    	    newSet.Add(val)
	    }
	}
	
	return newSet
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	newSet := NewSet[T]()
	
	for _, val := range s1.Elements() {
	    if !s2.Contains(val) {
    	    newSet.Add(val)
	    }
	}
	
	return newSet
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	new := make([]T, 0, len(slice))
	for _, val := range slice {
	    if predicate(val) {
	        new = append(new, val)
	    }
	}
	return new
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	mapped := make([]U, 0, len(slice))
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
	for _, val := range slice {
	    if val == element {
	        return true
	    }
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for i, val := range slice {
	    if val == element {
	        return i
	    }
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	set := NewSet[T]()
	for _, val := range slice {
	    set.Add(val)
	}
	return set.Elements()
}
