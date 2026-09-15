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
    l := len(m.Employees)
    
    if cap(m.Employees)-l < 1 {
        m.Employees = append(make([]Employee, 0, l+1), m.Employees...)
    }
    
    m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
	employeeIndex := m.findEmployeeIndexByID(id)
	
	if employeeIndex != -1 {
	    m.Employees = append(m.Employees[:employeeIndex],m.Employees[employeeIndex+1:]...)
	}
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	tot := 0.0
	n := len(m.Employees)
	
	if n == 0 {
	    return tot
	}
	
	for _,employee := range m.Employees {
	    tot += employee.Salary
	}
	
	return tot / float64(n)
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	employeeIndex := m.findEmployeeIndexByID(id)
    if employeeIndex == -1 {
	   return nil
    }
    
	return &m.Employees[employeeIndex]
}

func (m *Manager) findEmployeeIndexByID(id int) int {
    return slices.IndexFunc(m.Employees, func(e Employee) bool { return e.ID == id })
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