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
    items []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
	    items: make([]T, 0),
	}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
    s.items = append(s.items, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
    length := len(s.items)
    
    if length > 0 {
        top := s.items[length-1]
        s.items = s.items[:length-1]
        return top, nil
    }
    
	var zero T
	return zero, ErrEmptyCollection
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
    length := len(s.items)
    
    if length > 0 {
        return s.items[length-1], nil
    }
    
	var zero T
	return zero, ErrEmptyCollection
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.items)
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
    items []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
	    items: make([]T, 0),
	}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
    q.items = append(q.items, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	
    if len(q.items) > 0 {
        front := q.items[0]
        q.items[0] = zero
        q.items = q.items[1:]
        return front, nil
    }
    
	return zero, ErrEmptyCollection
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
    if len(q.items) > 0 {
        return q.items[0], nil
    }
    
	var zero T
	return zero, ErrEmptyCollection
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.items)
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
    items map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
	    items: make(map[T]struct{}),
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
    if s.items == nil {
        s.items = make(map[T]struct{})
    }
    
    s.items[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
    if _, ok := s.items[value]; ok {
        delete(s.items, value)
    }
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
    _, ok := s.items[value]
    
    return ok
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.items)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
    var i int
    elements := make([]T, len(s.items))
    
    for k, _ := range s.items {
        elements[i] = k
        i++
    }
    
	return elements
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
    union := make(map[T]struct{})
    
    for k, _ := range s1.items {
        union[k] = struct{}{}
    }
    
    for k, _ := range s2.items {
        union[k] = struct{}{}
    }
    
	return &Set[T]{items: union}
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
    intersection := make(map[T]struct{})
    
    for k, _ := range s1.items {
        if _, ok := s2.items[k]; ok {
            intersection[k] = struct{}{}
        }
    }
    
	return &Set[T]{items: intersection}
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
    diff := make(map[T]struct{})
    
    for k, _ := range s1.items {
        if _, ok := s2.items[k]; !ok {
            diff[k] = struct{}{}
        }
    }
    
	return &Set[T]{items: diff}
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
    filtered := make([]T, 0)
    
    for _, val := range slice {
        if predicate(val) {
            filtered = append(filtered, val)
        }
    }
    
	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
    mapped := make([]U, len(slice))
    
    for i, val := range slice {
        mapped[i] = mapper(val)
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
    seen := make(map[T]struct{})
    unique := make([]T, 0)
    
    for _, val := range slice {
        if _, ok := seen[val]; !ok {
            unique = append(unique, val)
        }
        
        seen[val] = struct{}{}
    }
    
	return unique
}
