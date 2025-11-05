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
	for i, e := range m.Employees {
	    if(e.ID == id){
	        m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
	        break;
	    }
	}
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
    
    countEmployees := len(m.Employees)
    sumSalary := float64(0)
    
    if(countEmployees == 0){
        return sumSalary
    }
    
	for _, e := range m.Employees {
	    sumSalary += e.Salary
	}
	
	avgSalary := sumSalary / float64(len(m.Employees))
	
	return avgSalary
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	for _, e := range m.Employees {
	    if(e.ID == id){
	        return &e
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
