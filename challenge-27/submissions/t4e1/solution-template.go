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
	elements []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	// TODO: Implement this function
	return &Stack[T]{elements: []T{}}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	// TODO: Implement this method
	s.elements = append(s.elements, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	// TODO: Implement this method
	var zero T
	if len(s.elements) == 0 {
		return zero, ErrEmptyCollection
	}
	zero = s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]

	return zero, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	// TODO: Implement this method
	var zero T
	if len(s.elements) == 0 {
		return zero, ErrEmptyCollection
	}
	zero = s.elements[len(s.elements)-1]

	return zero, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	// TODO: Implement this method
	return len(s.elements)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	// TODO: Implement this method
	if len(s.elements) != 0 {
		return false
	}
	return true
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	// TODO: Add necessary fields
	elements []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	// TODO: Implement this function
	return &Queue[T]{elements: []T{}}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	// TODO: Implement this method
	q.elements = append(q.elements, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	// TODO: Implement this method
	var zero T
	if len(q.elements) == 0 {
		return zero, ErrEmptyCollection
	}
	zero = q.elements[0]
	q.elements = q.elements[1:]

	return zero, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	// TODO: Implement this method
	var zero T
	if len(q.elements) == 0 {
		return zero, ErrEmptyCollection
	}
	zero = q.elements[0]

	return zero, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	// TODO: Implement this method
	return len(q.elements)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	// TODO: Implement this method
	if len(q.elements) != 0 {
		return false
	}
	return true
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	// TODO: Add necessary fields
	elements map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	// TODO: Implement this function
	return &Set[T]{elements: make(map[T]struct{})}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
	// TODO: Implement this method
	s.elements[value] = struct{}{} 
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	// TODO: Implement this method
	if _, exists := s.elements[value]; exists {
		delete(s.elements, value)
	}
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	// TODO: Implement this method
	if _, exists := s.elements[value]; exists {
		return true
	}

	return false
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	// TODO: Implement this method
	return len(s.elements)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	// TODO: Implement this method
	sv := make([]T, 0, len(s.elements))
	for k, _ := range s.elements {
		sv = append(sv, k)
	}

	return sv
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	// TODO: Implement this function
	us := make(map[T]struct{})
	for k, _ := range s1.elements {
		us[k] = struct{}{}
	}

	for k, _ := range s2.elements {
		us[k] = struct{}{}
	}

	return &Set[T]{elements: us}
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
    interSet := make(map[T]struct{})
    for k1 := range s1.elements {
        if _, exists := s2.elements[k1]; exists {
            interSet[k1] = struct{}{}
        }
    }
    return &Set[T]{elements: interSet}
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
    diff := make(map[T]struct{})
    for k := range s1.elements {
        if _, exists := s2.elements[k]; !exists {
            diff[k] = struct{}{}
        }
    }
    return &Set[T]{elements: diff}
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
    	acc := initial
    	for _, v := range slice {
        	acc = reducer(acc, v)
    	}
    	return acc
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
    	seen := make(map[T]struct{})
    	result := make([]T, 0, len(slice))
    	for _, v := range slice {
        	if _, exists := seen[v]; !exists {
            	seen[v] = struct{}{}
            	result = append(result, v)
        	}
    	}
    	return result
	}
