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
	
	var index = -1 
	
	// First let's match the id
	for i, emp := range m.Employees {
	    if id == emp.ID {
	        index = i
	        break
	    }
	}
	
	if index == -1 {
	    return 
	}
	
	m.Employees = append(m.Employees[:index], m.Employees[index+1:]...)
	
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	// TODO: Implement this method
	
	if len(m.Employees) == 0 {
	    return 0
	}
	
	var total float64 
	
	for _, emp := range m.Employees{
	    total += emp.Salary
	}
	
	return total / float64(len(m.Employees))
	
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	// TODO: Implement this method
	
	for _, emp := range m.Employees {
	    if id == emp.ID {
	        return &emp
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
