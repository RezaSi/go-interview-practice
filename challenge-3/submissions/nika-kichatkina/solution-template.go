package main

import "fmt"

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
	employee := -1
	for i, v := range m.Employees {
		if v.ID == id {
			employee = i
			break
		}
	}
	if employee != -1 {m.Employees = append(m.Employees[:employee], m.Employees[employee+1:]...)}
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	sum := 0.0

	for _, v := range m.Employees {
		sum += v.Salary
	}
	result:= 0.000000
	if sum != 0.0 {result = sum / float64(len(m.Employees))}
	return result
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	employee:= Employee{}
	for _, v := range m.Employees {
		if v.ID == id {
			employee = v
			break
		}
	}
    if employee != (Employee{}) {
		return &employee
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