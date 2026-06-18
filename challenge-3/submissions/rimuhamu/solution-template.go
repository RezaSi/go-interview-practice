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
	for i, emp := range m.Employees{
	    if emp.ID == id{
	        m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
	        break
	    }
	} 
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	var totalSalary float64
	
	for i := range m.Employees{
	    totalSalary += m.Employees[i].Salary
	}
	
	var averageSalary float64
	var totalEmployee float64
	totalEmployee = float64(len(m.Employees))
	averageSalary = totalSalary/totalEmployee
	
	if(totalEmployee < 1 ){
	    return 0.00
	}
	return averageSalary
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	for _, emp := range m.Employees{
	    if emp.ID == id{
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

	fmt.Printf("Average Salary: %f\n", averageSalary)
	if employee != nil {
		fmt.Printf("Employee found: %+v\n", *employee)
	}
}
