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
	m.Employees = append(m.Employees, e)
}

// RemoveEmployee removes an employee by ID from the manager's list.
func (m *Manager) RemoveEmployee(id int) {
	// TODO: Implement this method
	for i := 0; i < len(m.Employees) - 1; i++{
	    if m.Employees[i].ID == id {
	        m.Employees = append(m.Employees[:i], m.Employees[i + 1:]...)
	        break
	    }else{
	        continue
	    }
	}
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	// TODO: Implement this method
	var sum = 0
	for i := 0; i < len(m.Employees); i++{
        sum += int(m.Employees[i].Salary)
	}
	if sum == 0 {
	    return float64(0)
	}
	sum /= len(m.Employees)
	
	return float64(sum)
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
	// TODO: Implement this method
		for i := 0; i < len(m.Employees) - 1; i++{
	    if m.Employees[i].ID == id {
	       e := &Employee{
	           m.Employees[i].ID,
	           m.Employees[i].Name,
	           m.Employees[i].Age,
	           m.Employees[i].Salary,
	       }
	       return e
	    }else{
	        continue
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
