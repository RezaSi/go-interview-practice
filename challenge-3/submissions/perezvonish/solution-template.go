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
	if m.Employees == nil {
		m.Employees = []Employee{}
	}
	m.Employees = append(m.Employees, e)
}

func (m *Manager) RemoveEmployee(id int) {
	for i, em := range m.Employees {
		if em.ID == id {
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
	for _, em := range m.Employees {
		total += em.Salary
	}

	return total / float64(len(m.Employees))
}

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

	fmt.Printf("Average Salary: %.2f\n", averageSalary)
	if employee != nil {
		fmt.Printf("Employee found: %+v\n", *employee)
	} else {
		fmt.Println("Employee not found")
	}
}
