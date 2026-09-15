package generics

import "errors"
import "fmt"

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
	// TODO: Implement this method
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
	data []T
	lastIdx int
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
    return &Stack[T]{
        lastIdx: -1,
    }
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.data = append(s.data, value)
	s.lastIdx++
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	var value T
	
	if s.lastIdx == -1 || len(s.data) == 0 {
	    return value, fmt.Errorf("empty stack")
	}
	
	value = s.data[s.lastIdx]
	s.data = s.data[:len(s.data)-1]
	s.lastIdx--
	return value, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	var value T
	
	if s.lastIdx == -1 || len(s.data) == 0 {
	    return value, fmt.Errorf("empty stack")
	}
	
	value = s.data[s.lastIdx]
	return value, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.data)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
    return len(s.data) == 0
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	data []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
    return &Queue[T]{
    }
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.data = append(q.data, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
    var value T
    
    if len(q.data) == 0 {
        return value, fmt.Errorf("empty queue")
    }
    value = q.data[0]
    q.data = q.data[1:]
    return value, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
    var value T
    
    if len(q.data) == 0 {
        return value, fmt.Errorf("empty queue")
    }
    value = q.data[0]
    return value, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
    return len(q.data)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
    return len(q.data) == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	data map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
    return &Set[T]{
        data: make(map[T]struct{}),
    }
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
    _, ex := s.data[value]
    if !ex {
        s.data[value] = struct{}{}
    }
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.data, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
    _, ex := s.data[value]
    return ex
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
    return len(s.data)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
    result := []T{}
    for key := range s.data {
        result = append(result, key)
    }
    return result
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
    result := NewSet[T]()
    
    for key := range s1.data {
        result.data[key] = struct{}{}
    }
    for key := range s2.data {
        _, ex := s1.data[key]
        if !ex {
            result.data[key] = struct{}{}
        }
    }
    return result
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
    result := NewSet[T]()

    for key := range s1.data {
        _, ex := s2.data[key]
        if ex {
            result.data[key] = struct{}{}
        }
    }
    return result
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
    result := NewSet[T]()
    
    for key := range s1.data {
        _, ex := s2.data[key]
        if !ex {
            result.data[key] = struct{}{}
        }
    }
    return result
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := []T{}
    
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
    result := []U{}
    
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
        if v == element{
            return true
        }
    }
    return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for i, v := range slice {
        if v == element{
            return i
        }
    }
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	result := []T{}
	set := make(map[T]struct{})
	
	for _, v := range slice{
	     _, ex := set[v]
	    if !ex {
	        result = append(result, v)
	        set[v] = struct{}{}
	    }
	}
	return result
}
