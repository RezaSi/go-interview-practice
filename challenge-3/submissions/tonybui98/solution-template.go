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
	l := len(m.Employees)
	if m.Employees[0].ID == id {
		m.Employees = m.Employees[1:]
	} else if m.Employees[l-1].ID == id {
		m.Employees = m.Employees[:l-1]
	} else {
		i := 0
		for k, v := range m.Employees {
			if v.ID == id {
				i = k
				break
			}
		}
		if i > 0 {
			m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
		}
	}
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	avg := 0.0
	count := 0
	for _, v := range m.Employees {
		avg += v.Salary
		count++
	}
	if count == 0 {
		return 0.0
	}
	return avg / float64(count)
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	for _, v := range m.Employees {
		if v.ID == id {
			return &v
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
