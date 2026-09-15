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

func (m *Manager) AddEmployee(e Employee) {
	m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
    if len(m.Employees) == 0 {
	    return
    }

	for i := range m.Employees {
	    if m.Employees[i].ID == id {
	        m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
	        return
	    }
	}
}

func (m *Manager) GetAverageSalary() float64 {
	if len(m.Employees) == 0 {
	    return 0
	}
	
	total := 0.0
	count := len(m.Employees)
	
	for _, v := range m.Employees {
	    total += v.Salary
	}
	
	return total / float64(count)
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	if len(m.Employees) == 0 {
	    return nil
	}
	
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
