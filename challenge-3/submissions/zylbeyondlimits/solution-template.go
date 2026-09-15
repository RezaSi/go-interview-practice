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

// RemoveEmployee 推荐做法：原地删除
func (m *Manager) RemoveEmployee(id int) {
    for i, e := range m.Employees {
        if e.ID == id {
            // 使用切片拼接实现原地删除，无需申请 newEmployees 内存
            m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
            return // 找到并删除后直接返回
        }
    }
}

// GetAverageSalary calculates the average salary of all employees.
func (m *Manager) GetAverageSalary() float64 {
	// TODO: Implement this method
	if len(m.Employees) == 0 {
	    return 0
	}
	sum := 0.0
	for _, employee := range(m.Employees) {
	    sum = sum + employee.Salary
	}
	return sum / float64(len(m.Employees))
}

// FindEmployeeByID 返回真实的内存地址
func (m *Manager) FindEmployeeByID(id int) *Employee {
    for i := range m.Employees {
        if m.Employees[i].ID == id {
            return &m.Employees[i] // 返回切片中元素的真实地址
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
