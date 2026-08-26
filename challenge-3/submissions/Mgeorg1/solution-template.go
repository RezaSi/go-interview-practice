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
	Employees []Employee
	rwm       sync.RWMutex
}

func (m *Manager) find(id int) *Employee {
	for i := range len(m.Employees) {
		if m.Employees[i].ID == id {
			return &m.Employees[i]
		}
	}
	return nil
}

// AddEmployee adds a new employee to the manager's list.
func (m *Manager) AddEmployee(e Employee) {
	m.rwm.Lock()
	defer m.rwm.Unlock()

	/*	existEmployee := m.find(e.ID)
		if existEmployee != nil {
		    return
		}*/
	m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
	m.rwm.Lock()
	defer m.rwm.Unlock()
	keep := 0
	for i := range len(m.Employees) {
		if m.Employees[i].ID != id {
			m.Employees[keep] = m.Employees[i]
			keep++
		}
	}
	m.Employees = m.Employees[:keep]
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	m.rwm.RLock()
	defer m.rwm.RUnlock()

	if len(m.Employees) == 0 {
		return 0
	}

	sum := float64(0)
	for _, v := range m.Employees {
		sum += v.Salary
	}

	return sum / float64(len(m.Employees))
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	m.rwm.RLock()
	defer m.rwm.RUnlock()

	employee := m.find(id)
	if employee == nil {
		return employee
	}
	snapshot := *employee
	return &snapshot
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
