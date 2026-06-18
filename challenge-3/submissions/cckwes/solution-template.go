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
	foundIndex := m.findEmployeeIndexByID(id)
	if foundIndex == -1 {
	    return
	}
	
	m.Employees = append(m.Employees[:foundIndex], m.Employees[foundIndex+1:]...)
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
    if len(m.Employees) == 0 {
        return 0
    }
    
	totalSalary := 0.0
	
	for _, v := range m.Employees {
	    totalSalary += v.Salary
	}
	
	return totalSalary / float64(len(m.Employees))
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	foundIndex := m.findEmployeeIndexByID(id)
	if foundIndex == -1 {
	    return nil
	}
	
	return &m.Employees[foundIndex]
}

func (m *Manager) findEmployeeIndexByID(id int) int {
    for index, v := range m.Employees {
        if v.ID == id {
            return index
        }
    }
    
    return -1
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