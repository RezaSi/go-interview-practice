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
    m.Employees = slices.DeleteFunc(m.Employees, func (e Employee) bool {
        return e.ID == id
    })
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
    numEmployees := len(m.Employees)
    
    if numEmployees == 0 {
        return 0.0
    }
    
    salarySum := 0.0
    for _, e := range m.Employees {
        salarySum += e.Salary
    }
    
    return salarySum / float64(numEmployees)
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
    i := slices.IndexFunc(m.Employees, func (e Employee) bool {
        return e.ID == id
    })
    
    if i == -1 {
        return nil
    }
    
    return &m.Employees[i]
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
