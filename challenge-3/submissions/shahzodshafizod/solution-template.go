package main

import (
	"fmt"
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
func (m *Manager) AddEmployee(e Employee) { // O(1) or O(n)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Employees = append(m.Employees, e) // O(1) or O(n)
	m.totalSalary += e.Salary
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) bool { // O(n)
	m.mu.Lock()
	defer m.mu.Unlock()
	for index := range m.Employees {
		if m.Employees[index].ID == id {
			m.totalSalary -= m.Employees[index].Salary
			m.Employees = append(m.Employees[:index], m.Employees[index+1:]...)
			return true
		}
	}
	return false
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 { // O(1)
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.Employees) == 0 {
		return 0
	}
	return m.totalSalary / float64(len(m.Employees))
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee { // O(n)
	m.mu.RLock()
	defer m.mu.RUnlock()
	for index := range m.Employees {
		if m.Employees[index].ID == id {
			return &m.Employees[index]
		}
	}
	return nil
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
