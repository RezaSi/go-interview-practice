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
	// TODO: Implement this method
	m.Employees = append(m.Employees,Employee{
	    ID: e.ID,
	    Name: e.Name, 
	    Age: e.Age,
	    Salary: e.Salary,
	})
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
	// TODO: Implement this method
	for index, item := range m.Employees {
	    if (item.ID == id) {
	        m.Employees = append( m.Employees[:index],m.Employees[index + 1:]...)
	    }
	}
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	sum := 0.000000 
	if (len(m.Employees) < 1) {
	    return sum
	}
	for _, item := range m.Employees {
	    sum += item.Salary
	}
	return sum / float64(len(m.Employees))
}

// FindEmployeeByID finds and returns an employee by their ID.
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

	fmt.Printf("Average Salary: %f\n", averageSalary)
	if employee != nil {
		fmt.Printf("Employee found: %+v\n", *employee)
	}
}
