package main

import "fmt"

// is this what I need to do for docstrings
type Employee struct {
	ID     int
	Name   string
	Age    int
	Salary float64
}

// contains a list of Employees under their charge
type Manager struct {
	Employees []Employee
}

// adds an employee to the manager's list
func (m *Manager) AddEmployee(e Employee) {
    m.Employees = append(m.Employees, e)
}

// removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
    for i, emp := range m.Employees {
        if (emp.ID == id) {
           m.Employees[i] = m.Employees[len(m.Employees)-1];
           m.Employees = m.Employees[:len(m.Employees)-1]
        }
    }
}

// calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
    slice := m.Employees
    if (len(slice) == 0) {
        return 0
    }
    var res float64
    for i := range slice {
        res += slice[i].Salary
    }
    res /= float64(len(slice))
    return res
}

// finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
    for i,emp := range m.Employees {
        if (emp.ID == id) {
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
