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

	var idx int = -1

	for i, emp := range m.Employees {
		if emp.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return
	}
	
	if idx == 0 {
		m.Employees = m.Employees[idx+1:]
	} else if idx == len(m.Employees)-1 {
		m.Employees = m.Employees[:idx]
	} else {
		prev := m.Employees[:idx]
		after := m.Employees[idx+1:]

		prev = append(prev, after...)

		m.Employees = prev
	}

}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {

	tot := 0.0
	for _, emp := range m.Employees {
		tot += (emp.Salary)
	}

	var sz int = len(m.Employees)

	if sz == 0 {
		return 0
	}

	avg := tot / float64(sz)

	return avg
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {

	for i, ep := range m.Employees {
		if ep.ID == id {
			return &m.Employees[i]
		}
	}

	return nil
}

func main() {
	manager := Manager{}
	manager.AddEmployee(
		Employee{
			ID:     1,
			Name:   "Alice",
			Age:    30,
			Salary: 70000,
		},
	)
	manager.AddEmployee(
		Employee{
			ID:     2,
			Name:   "Bob",
			Age:    25,
			Salary: 65000,
		},
	)
	manager.RemoveEmployee(1)
	averageSalary := manager.GetAverageSalary()
	employee := manager.FindEmployeeByID(2)

	fmt.Printf("Average Salary: %f\n", averageSalary)
	if employee != nil {
		fmt.Printf("Employee found: %+v\n", *employee)
	}
}
