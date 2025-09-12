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
			// Remove emplyee by slicing
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
	for _, emp := range m.Employees {
		total += emp.Salary
	}
	return total / float64(len(m.Employees))
}

func (m *Manager) FindEmployeeByID(id int) *Employee {
	for _, emp := range m.Employees {
		if emp.ID == id {
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

	fmt.Printf("Average Salary: %.2f\n", averageSalary)
	if employee != nil {
		fmt.Printf("Employee Found: %+v\n", *employee)
	} else {
		fmt.Println("Employee not found")
	}
}
