package main

import (
	"fmt"
	"sort"
	"sync"
)

type Employee struct {
	ID     int
	Name   string
	Age    int
	Salary float64
}

type Manager struct {
	Employees   []Employee
	totalSalary float64
	mu          sync.RWMutex
}

// AddEmployee adds a new employee to the manager's list.
func (m *Manager) AddEmployee(e Employee) { // Time: O(n)
	m.mu.Lock()
	defer m.mu.Unlock()
	// sort.Slice(m.Employees, func(i, j int) bool {
	// 	return m.Employees[i].ID < m.Employees[j].ID
	// })
	idx := m.search(e.ID)
	if idx == -1 {
		idx = len(m.Employees)
	}
	m.Employees = append(m.Employees, Employee{})
	copy(m.Employees[idx+1:], m.Employees[idx:])
	m.Employees[idx] = e
	m.totalSalary += e.Salary
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) bool { // Time: O(n)
	m.mu.Lock()
	defer m.mu.Unlock()
	index := m.search(id) // O(log n)
	if index != -1 {
		m.totalSalary -= m.Employees[index].Salary
		m.Employees = append(m.Employees[:index], m.Employees[index+1:]...) // O(n)
		return true
	}
	return false
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 { // Time: O(1)
	m.mu.RLock()
	defer m.mu.RUnlock()
	if n := len(m.Employees); n > 0 {
		return m.totalSalary / float64(n)
	}
	return 0
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee { // Time: O(log n)
	m.mu.RLock()
	defer m.mu.RUnlock()
	index := m.search(id)
	if index != -1 {
		return &m.Employees[index]
	}
	return nil
}

func (m *Manager) search(id int) int { // Time: O(log n)
	index := sort.Search(len(m.Employees), func(i int) bool {
		return m.Employees[i].ID >= id
	}) // sort.Search returns an index where a new item could be inserted, so we do not receive -1 or an error
	if index < len(m.Employees) && m.Employees[index].ID == id {
		return index
	}
	return -1
}

func main() {
	manager := Manager{}
	manager.AddEmployee(Employee{ID: 1, Name: "Alice", Age: 30, Salary: 70000})
	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})
	manager.RemoveEmployee(1)
	averageSalary := manager.GetAverageSalary()
	employee := manager.FindEmployeeByID(2)

	fmt.Printf("Average Salary: %f\n", averageSalary)
	if employee != nil {
		fmt.Printf("Employee found: %+v\n", *employee)
	}
}
