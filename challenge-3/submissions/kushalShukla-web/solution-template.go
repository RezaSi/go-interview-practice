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

func (m *Manager) GetEmployee() {

}

func (m *Manager) AddEmployee(em Employee) {
	m.Employees = append(m.Employees, em)
}

func (m *Manager) RemoveEmployee(id int) {
	var tempValue []Employee
	for i := 0; i < len(m.Employees); i++ {
		if m.Employees[i].ID == id {
			continue
		} else {
			tempValue = append(tempValue, m.Employees[i])
		}
	}
	m.Employees = tempValue
}

func (m *Manager) GetAverageSalary() float64 {
	if len(m.Employees) == 0 {
		return 0
	}
	var averagesalary float64
	for i := 0; i < len(m.Employees); i++ {
		averagesalary += m.Employees[i].Salary
	}

	averagesalary = averagesalary / float64(len(m.Employees))
	return averagesalary
}

func (m *Manager) FindEmployeeByID(id int) *Employee {
	for i := 0; i < len(m.Employees); i++ {
		if m.Employees[i].ID == id {
			return &m.Employees[i]
		}
	}
	return nil
}

func main() {
	manager := Manager{}
	manager.AddEmployee(Employee{ID: 1, Name: "Kushal", Age: 23, Salary: 4500})
	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})
	manager.RemoveEmployee(1)
	averagesalary := manager.GetAverageSalary()
	fmt.Println("The average salary is ", averagesalary)
	employee := manager.FindEmployeeByID(2)
	if employee != nil {
		fmt.Printf("Employee found: %+v\n", *employee)
	}
}
