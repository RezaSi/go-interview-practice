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

func (m *Manager) RemoveEmployee(id int) {
	for i, emp := range m.Employees {
		if emp.ID == id {
			m.Employees[i] = m.Employees[len(m.Employees)-1]
			m.Employees = m.Employees[:len(m.Employees)-1]
			return
		}
	}
}

func (m *Manager) GetAverageSalary() float64 {
	if len(m.Employees) == 0 {
		return 0.0
	}
	s := 0.0
	for _, emp := range m.Employees {
		s += emp.Salary
	}
	return s / float64(len(m.Employees))
}

func (m *Manager) FindEmployeeByID(id int) *Employee {
	for ind, emp := range m.Employees {
		if emp.ID == id {
			return &m.Employees[ind]
		}
	}
	return nil
}

func main() {
	manager := Manager{}
	manager.AddEmployee(Employee{ID: 1, Name: "Alice", Age: 30, Salary: 70000})
	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})
	fmt.Printf("Average Salary: %f\n", manager.GetAverageSalary())
	manager.RemoveEmployee(1)
	employee := manager.FindEmployeeByID(2)
	if employee != nil {
		fmt.Printf("Employee found: %+v\n", *employee)
	}
	manager.RemoveEmployee(2)
	fmt.Println(manager.Employees)
}