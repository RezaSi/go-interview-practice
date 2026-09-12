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
	return Pair[U, T] {
	    First: p.Second,
	    Second: p.First,
	}
}

//
// 2. Generic Stack
//

// Stack is a generic Last-In-First-Out (LIFO) data structure
type Stack[T any] struct {
	elems []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.elems = append(s.elems, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
    var zero T
    
	if s.IsEmpty() {
	    return zero, ErrEmptyCollection
	}
	
	idx := len(s.elems)-1
	top := s.elems[idx]
	
	s.elems[idx] = zero
	s.elems = s.elems[:idx]
	
	return top, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	if s.IsEmpty() {
	    var zero T
	    return zero, ErrEmptyCollection
	}
	
	elem := s.elems[len(s.elems)-1]
	
	return elem, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.elems)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elems) == 0 
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	elems []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.elems = append(q.elems, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
    var zero T
    
	if q.IsEmpty() {
	    return zero, ErrEmptyCollection
	}
	
	old := q.elems
	front := old[0]
	old[0] = zero
	q.elems = old[1:]
	
	if len(q.elems) == 0 {
	    q.elems = nil
	}
 	
	return front, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	if q.IsEmpty() {
	    var zero T
	    return zero, ErrEmptyCollection
	}
	
	return q.elems[0], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.elems)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	return len(q.elems) == 0
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	elems map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
	    elems: make(map[T]struct{}),
	}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {
    s.elems[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
    delete(s.elems, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
	_, exists := s.elems[value]
	return exists
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.elems)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
    if s == nil || len(s.elems) == 0 {
        return nil
    }
    
 	list := make([]T, len(s.elems))
	i := 0
	for e := range s.elems {
	    list[i] = e
	    i++
	}
	return list
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
    maxSize := s1.Size() + s1.Size()
    
	union := &Set[T]{
	    elems: make(map[T]struct{}, maxSize),
	}
	
	for k := range s1.elems {
	    union.elems[k] = struct{}{}
	}
	for k := range s2.elems {
	    union.elems[k] = struct{}{}
	}
	
	return union
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
    var big, small *Set[T]
    if s2.Size() > s1.Size() {
        big, small = s2, s1
    } else {
        big, small = s1, s2
    }
    
	inter := &Set[T]{
	    elems: make(map[T]struct{}, small.Size()),
	}
	
	for k := range small.elems {
	    if _, exists := big.elems[k]; exists {
	        inter.elems[k] = struct{}{}
	    }
	}
	
	return inter
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
    if s1 == nil || s1.Size() == 0 {
        return NewSet[T]()
    }
    
    if s2 == nil || s2.Size() == 0 {
        return Union(s1, NewSet[T]())
    }
    
	diff := &Set[T]{
	    elems: make(map[T]struct{}, s1.Size()),
	}
	
	for k := range s1.elems {
	    if _, exists := s2.elems[k]; !exists {
	        diff.elems[k] = struct{}{}
	    }
	}
	
	return diff
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filtered := make([]T, 0, len(slice) / 10)
	for _, v := range slice {
	    if predicate(v) {
	        filtered = append(filtered, v)
	    }
	}
	
	return filtered
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	res := make([]U, len(slice))
	for i, v := range slice {
	    res[i] = mapper(v)
	}
	return res
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
	    if element == v {
	        return true
	    }
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for i, v := range slice {
	    if element == v {
	        return i
	    }
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	newSlice := make([]T, 0, len(slice))
	unique := make(map[T]struct{}, len(slice))
	
	for _, v := range slice {
	    if _, exists := unique[v]; !exists {
	        unique[v] = struct{}{}
	        newSlice = append(newSlice, v)
	    }
	}
	
	return newSlice
}
