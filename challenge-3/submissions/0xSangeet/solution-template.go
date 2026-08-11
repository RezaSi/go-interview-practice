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
	for index, value := range m.Employees {
	    if value.ID == id {
	        m.Employees = append(m.Employees[:index], m.Employees[index+1:]...)
	    }
	}
}

func (m *Manager) GetAverageSalary() float64 {
	var sum float64
	var count int
	for _, val := range m.Employees {
	    sum += val.Salary
	    count += 1
	}
	if count != 0 {
	    return sum/float64(count)
	}
	return 0
}

func (m *Manager) FindEmployeeByID(id int) *Employee {
	for _, value := range m.Employees {
	    if value.ID == id {
	        return &value
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
