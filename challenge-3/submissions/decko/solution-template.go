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
    position := -1
    for i, employee := range m.Employees {
        if employee.ID == id {
            position = i
            break
        }
    }
    if position != -1 {
        m.Employees = append(m.Employees[:position], m.Employees[position+1:]...)
    }
}
// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	var salary_sum float64
	n_employees := float64(len(m.Employees))
	
	if n_employees == 0 {
	    return 0.0
	}
	
    for _, employee := range m.Employees {
        salary_sum += employee.Salary
    } 
    
    avg_salary := salary_sum / n_employees
	return avg_salary
}

// FindEmployeeByID finds and returns an employee by their ID.
func (m *Manager) FindEmployeeByID(id int) *Employee {
    for _, employee := range m.Employees {
        if employee.ID == id {
            return &employee
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
