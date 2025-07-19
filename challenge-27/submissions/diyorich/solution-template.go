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
	stack []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{stack: []T{}}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
    if len(s.stack) > 0 {
        element := s.stack[len(s.stack)-1]
        s.stack = s.stack[:len(s.stack)-1]
        return element, nil
    }
	
	var zero T
	return zero, ErrEmptyCollection
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	if len(s.stack) > 0 {
	    return s.stack[len(s.stack)-1], nil
	}
	
	var zero T
	return zero, ErrEmptyCollection
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.stack)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
    if len(s.stack) != 0 {
        return false
    }
    
	return true
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	q []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{q: []T{}}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.q = append(q.q, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	if len(q.q) > 0 {
	    elem := q.q[0]
	    q.q = q.q[1:]
	    return elem, nil
	}
	
	var zero T
	return zero, ErrEmptyCollection
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	if len(q.q) > 0 {
	    return q.q[0], nil
	}
	var zero T
	return zero, ErrEmptyCollection
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.q)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	if len(q.q) > 0 {
	    return false
	}
	return true
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
    store map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{store: make(map[T]struct{})}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	s.store[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.store, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	if _, ok := s.store[value]; ok {
	    return true
	}
	return false
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.store)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	var elems []T
	for key, _ := range s.store{
	    elems = append(elems, key)
	    
	}
	
    return elems
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	for key, _ := range s1.store {
	    s2.store[key] = struct{}{}
	}
	return s2
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	intersection := &Set[T]{store: make(map[T]struct{})}
	
	for key, _ := range s1.store {
	    if _, ok := s2.store[key]; ok {
	        intersection.store[key] = struct{}{}
	    }
	}
	
	return intersection
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
    diff := &Set[T]{store: make(map[T]struct{})}
    for key, _ := range s1.store {
	    if _, ok := s2.store[key]; !ok {
	        diff.store[key] = struct{}{}
	    }
	} 
    
	return diff
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filtered := []T{}
	
	for _, value := range slice {
	    if ok := predicate(value); ok {
	        filtered = append(filtered, value)
	    }
	}
	
	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	results := []U{}
	
	for _, value := range slice {
	    res := mapper(value)
	    results = append(results, res)
	}
	
	
	return results
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	for _, value := range slice {
	   initial = reducer(initial, value)
	}
	
	return initial
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	for _, value := range slice {
	    if value == element {
	        return true
	    }
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for index, value := range slice {
	    if value == element {
	        return index
	    }
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	cache := make(map[T]struct{})
	res := []T{}
	
	for _, value := range slice {
	    if _, ok := cache[value]; !ok {
	        res = append(res, value)
	    }
	    cache[value] = struct{}{}
	}
	
	return res
}
