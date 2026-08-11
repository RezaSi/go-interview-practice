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

func (m *Manager) employeeIndex(id int) int {
    for idx, emp := range m.Employees {
        if emp.ID == id {
            return idx
        }
    }
    return -1
}

// AddEmployee adds a new employee to the manager's list.
func (m *Manager) AddEmployee(e Employee) {
    m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
    var idx = m.employeeIndex(id)
   
    if idx != -1 {
        m.Employees[idx] = m.Employees[len(m.Employees)-1]
        m.Employees = m.Employees[:len(m.Employees)-1]
    }
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	if len(m.Employees) > 0 {
    	var sum = 0.00
    	for _, emp := range m.Employees {
            sum = sum + emp.Salary
        }
        return sum/float64(len(m.Employees))
	    
	} else {
	    return 0.00
	}
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
    var idx = m.employeeIndex(id)
    
    if idx > -1 {
        return &m.Employees[idx]
    } else {
        return nil
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
