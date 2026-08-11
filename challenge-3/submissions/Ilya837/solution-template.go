package main

import (
	"fmt"
	"sort"
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
	sort.Slice(m.Employees, func(i, j int) bool {
		return m.Employees[i].ID < m.Employees[j].ID
	})
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
	ind, _ := sort.Find(len(m.Employees), func(i int) int {
		return id - m.Employees[i].ID
	})

	if ind == len(m.Employees) || m.Employees[ind].ID != id {
		return
	} else {
		m.Employees = append(m.Employees[:ind], m.Employees[ind+1:]...)
		return
	}

}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	res := 0.0

	l := float64(len(m.Employees))

	for i := 0; i < len(m.Employees); i++ {
		res += m.Employees[i].Salary / l
	}

	return res
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {

	ind, _ := sort.Find(len(m.Employees), func(i int) int {
		return id - m.Employees[i].ID
	})

	if ind == len(m.Employees) || m.Employees[ind].ID != id {
		return nil
	} else {
		return &m.Employees[ind]
	}
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
