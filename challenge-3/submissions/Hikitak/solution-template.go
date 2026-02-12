package main

import (
    "fmt"
    "slices"
)

type Employee struct {
	ID     int
	Name   string
	Age    int
	Salary float64
}

type Manager struct {
	Employees []Employee
}

// AddEmployee adds a new employee to the manager's list.
func (m *Manager) AddEmployee(e Employee) {
	m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
    i := slices.IndexFunc(m.Employees, func(e Employee) bool { return e.ID == id })
    if i < 0 || i >= len(m.Employees) {
        return 
    }
	m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	sum := 0.0
	if len(m.Employees) == 0 {
	    return sum
	}
	
	for _, e := range m.Employees {
	    sum += e.Salary
	}
	return sum / float64(len(m.Employees))
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
    i := slices.IndexFunc(m.Employees, func(e Employee) bool { return e.ID == id })
    if i<0 || i>=len(m.Employees) {
        return nil
    }
	return &m.Employees[slices.IndexFunc(m.Employees, func(e Employee) bool { return e.ID == id })]
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
