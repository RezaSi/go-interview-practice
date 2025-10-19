package main

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
	for i, e := range m.Employees {
		if e.ID == id {
			m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
		}
	}
}

func (m *Manager) GetAverageSalary() float64 {
	if len(m.Employees) == 0 {
		return 0
	}

	sum := 0.0

	for _, employee := range m.Employees {
		sum += employee.Salary
	}

	return sum / float64(len(m.Employees))
}

func (m *Manager) FindEmployeeByID(id int) *Employee {
	for _, employee := range m.Employees {
		if employee.ID == id {
			return &employee
		}
	}

	return nil
}
