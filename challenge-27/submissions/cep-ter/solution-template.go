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
	stack []T
}

// NewStack creates a new empty stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
}

// Pop removes and returns the top element from the stack
// Returns an error if the stack is empty
func (s *Stack[T]) Pop() (T, error) {
    if len(s.stack) == 0 {
        var zero T
        return zero, ErrEmptyCollection
    }
    item := s.stack[len(s.stack)-1]
    s.stack = s.stack[:len(s.stack)-1]
	
	return item, nil
}

// Peek returns the top element without removing it
// Returns an error if the stack is empty
func (s *Stack[T]) Peek() (T, error) {
    if len(s.stack) == 0 {
        var zero T
        return zero, ErrEmptyCollection
    }
	return s.stack[len(s.stack)-1], nil
}

// Size returns the number of elements in the stack
func (s *Stack[T]) Size() int {
	return len(s.stack)
}

// IsEmpty returns true if the stack contains no elements
func (s *Stack[T]) IsEmpty() bool {
	if len(s.stack)==0 {
	    return true
	}
	return false
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
	return &Queue[T]{}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
    q.items = append(q.items, value)
}

// Dequeue removes and returns the front element from the queue
// Returns an error if the queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
    if len(q.items) == 0 {
        var zero T
        return zero, ErrEmptyCollection
    }
    item := q.items[0]
    q.items = q.items[1:]
    
	return item, nil
}

// Front returns the front element without removing it
// Returns an error if the queue is empty
func (q *Queue[T]) Front() (T, error) {
    if len(q.items) == 0 {
        var zero T
        return zero, ErrEmptyCollection
    }
	return q.items[0], nil
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// IsEmpty returns true if the queue contains no elements
func (q *Queue[T]) IsEmpty() bool {
	if len(q.items)==0 {
	    return true
	}
	return false
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
	s.items[value] = struct{}{}
}

// Remove removes an element from the set if it exists
func (s *Set[T]) Remove(value T) {
	delete(s.items, value)
}

// Contains returns true if the set contains the given element
func (s *Set[T]) Contains(value T) bool {
     _, exists := s.items[value]
    return exists
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.items) 
}

// Elements returns a slice containing all elements in the set
func (s *Set[T]) Elements() []T {
    item := make([]T, 0, len(s.items))
    for it := range s.items{
        item = append(item, it)
    }
	return item
}

// Union returns a new set containing all elements from both sets
func Union[T comparable](s1, s2 *Set[T]) *Set[T] {
	s := NewSet[T]()
	for _,it1 := range s1.Elements(){
	    s.Add(it1)
	}
	for _,it2 := range s2.Elements(){
	    s.Add(it2)
	}
	return s
}

// Intersection returns a new set containing only elements that exist in both sets
func Intersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	s := NewSet[T]()
	for _,it1 := range s1.Elements(){
	   if  ok := s2.Contains(it1); ok{
       	    s.Add(it1)
	   }
	}
	return s
}

// Difference returns a new set with elements in s1 that are not in s2
func Difference[T comparable](s1, s2 *Set[T]) *Set[T] {
	s := NewSet[T]()
	for _,it1 := range s1.Elements(){
	   if  ok := s2.Contains(it1); !ok{
       	    s.Add(it1)
	   }
	}
	return s
}

//
// 5. Generic Utility Functions
//

// Filter returns a new slice containing only the elements for which the predicate returns true
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var items []T
	for _, item := range slice{
	    if predicate(item){
	        items = append(items, item)
	    }
	}
	return items
}

// Map applies a function to each element in a slice and returns a new slice with the results
func Map[T, U any](slice []T, mapper func(T) U) []U {
	var items []U
	for _, item := range slice{
        items = append(items, mapper(item))
	}
	return items
}

// Reduce reduces a slice to a single value by applying a function to each element
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	reduced := initial
	for _, item := range slice{
        reduced = reducer(reduced, item)
	}
	return reduced
}

// Contains returns true if the slice contains the given element
func Contains[T comparable](slice []T, element T) bool {
	for _, item := range slice{
        if item == element{
            return true
        }
	}
	return false
}

// FindIndex returns the index of the first occurrence of the given element or -1 if not found
func FindIndex[T comparable](slice []T, element T) int {
	for index, item := range slice{
        if item == element{
            return index
        }
	}
	return -1
}

// RemoveDuplicates returns a new slice with duplicate elements removed, preserving order
func RemoveDuplicates[T comparable](slice []T) []T {
	m := make(map[T]struct{})
	var dup []T
	for _, item := range slice{
	    m[item] = struct{}{}
	}
	for key := range m{
	    dup = append(dup, key)
	}
	return dup
}
