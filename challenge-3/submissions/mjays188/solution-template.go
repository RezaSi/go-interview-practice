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
	// TODO: Implement this method
	m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
	// TODO: Implement this method
	var idxOfEmpToRemove = 0
	isEmployeeWithGivenIdPresent := false
	for i:=0; i < len(m.Employees); i++ {
	    if m.Employees[i].ID == id {
	        idxOfEmpToRemove = i
	        isEmployeeWithGivenIdPresent = true
	        break
	    }
	}
	if !isEmployeeWithGivenIdPresent {
	    return
	}
	m.Employees = append(m.Employees[:idxOfEmpToRemove], m.Employees[idxOfEmpToRemove+1:]...)
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	// TODO: Implement this method
	var totalSal float64 = 0
	var avgSal float64 = 0
	if len(m.Employees) == 0 {
	    return avgSal
	}
	for i:=0; i < len(m.Employees); i++ {
	    totalSal += m.Employees[i].Salary
	}
	avgSal = totalSal / float64(len(m.Employees))
	return avgSal
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	// TODO: Implement this method
	for i:=0; i < len(m.Employees); i++ {
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
