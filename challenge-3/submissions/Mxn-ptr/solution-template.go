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
	pos := func(sliceEmployees []Employee, id int) int {
		for i, v := range sliceEmployees {
			if v.ID == id {
				return i
			}
		}
		return -1
	}(m.Employees, id)

	if pos != -1 {
		m.Employees = append(m.Employees[:pos], m.Employees[pos+1:]...)
	}
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	result := float64(0)
	for _, v := range m.Employees {
		result += v.Salary
	}
	if result > 0 {
		return result / float64(len(m.Employees))
	}
	return 0
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	for i := range m.Employees {
		if m.Employees[i].ID == id {
			return &m.Employees[i]
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
