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

//增加管理器成员
func (m *Manager) AddEmployee(e Employee) {
	m.Employees = append(m.Employees, e)
}

//根据id删除成员
func (m *Manager) RemoveEmployee(id int) {
	for i, emp := range m.Employees {
		if emp.ID == id {
			m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
			//找到后停止循环
			return
		}
	}
}

//计算员工的平均薪资
func (m *Manager) GetAverageSalary() float64 {
	if len(m.Employees) == 0 {
		return 0
	}

	var total float64

	for _, emp := range m.Employees {
		total += emp.Salary
	}
	return total / float64(len(m.Employees))
}

//根据 ID 查找员工，返回指针以支持空值处理
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
