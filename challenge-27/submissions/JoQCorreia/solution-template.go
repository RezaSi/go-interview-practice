package generics

import ("errors"
"slices")

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
	values []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{values: make([]T, 0)}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	
	s.values = append(s.values, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
	
	var zero T
	if len(s.values) == 0 {
		return zero, errors.New("The stack is empty!")
	}

	lastElem := len(s.values) - 1
	element := s.values[lastElem]
	s.values = s.values[:lastElem]
	return element, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
	
	var zero T

	if len(s.values) == 0 {
		return zero, errors.New("The stack is empty!")
	}
	
	lastElem := len(s.values) - 1
	element := s.values[lastElem]
	
	return element, nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	size := len(s.values)
	return size
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
    if s.Size() != 0 {
        return false
    }
	return true
}

//
// 3. Generic Queue
//

// Queue is a generic First-In-First-Out (FIFO) data structure
type Queue[T any] struct {
	values []T
}

// NewQueue creates a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{values: make([]T, 0)}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q. values = append(q.values, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	
	if len(q.values) == 0 {
	    return zero, errors.New("The queue is empty!")
	}
	firstVal := q.values[0]
	q.values = q.values [1:]
	return firstVal, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
	
    var zero T
	
	if len(q.values) == 0 {
	    return zero, errors.New("The queue is empty!")
	}
	firstVal := q.values[0]

	return firstVal, nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
    	size := len(q.values)
	return size
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
    if q.Size() != 0 {
        return false
    }
	return true
}

//
// 4. Generic Set
//

// Set is a generic collection of unique elements
type Set[T comparable] struct {
	values []T
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{values: make([]T, 0)}
}

// Add adds an element to the set if it's not already present
func (s *Set[T]) Add(value T) {

	if slices.Contains(s.values, value) {
		return
	}

	s.values = append(s.values, value)

}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	if slices.Contains(s.values, value) {
		ind := slices.Index(s.values, value)
		s.values = slices.Delete(s.values, ind, (ind + 1))
	}
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {

	if slices.Contains(s.values, value) {
		return true
	}
	return false
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {

	return len(s.values)
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
	newSlice := s.values
	return newSlice
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	s3 := NewSet[T]()
	s3.values = s1.values

	for _, itr := range s2.values {
		s3.Add(itr)
	}

	return s3
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	s3 := NewSet[T]()

	for _, val1 := range s1.values {

		if slices.Contains(s2.values, val1) {
			s3.values = append(s3.values, val1)
		}

	}

	return s3
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	s3 := NewSet[T]()

	for _, val1 := range s1.values {

		if !slices.Contains(s2.values, val1) {
			s3.values = append(s3.values, val1)
		}

	}

	return s3
}
//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var s2 []T

	for _, itr := range slice {

		if predicate(itr) == true {
			s2 = append(s2, itr)
		}

	}

	return s2
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	var s2 []U

	for _ , itr := range slice {

		s2 = append(s2, mapper(itr))
	}


	return s2
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
    
	s2 := initial

	if len(slice) > 0 {
		for _, itr := range slice {
		    
			s2 = reducer(s2, itr)
		}
		return s2
	}

	return initial
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {

	return slices.Contains(slice, element)
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {

	return slices.Index(slice, element)
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	
	numbers := slice
	
	for i, num := range numbers {
		var j int
		for _, loopN := range numbers {

			if num == loopN && i != j {

				numbers = slices.Delete(numbers, j, (j + 1))
				j--
			}
			j++

		}
		j = 0
	}
	return numbers
	
}