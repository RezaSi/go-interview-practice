package main

import (
	"testing"
)

func TestAddEmployee(t *testing.T) {
	manager := Manager{}

	manager.AddEmployee(Employee{ID: 1, Name: "Alice", Age: 30, Salary: 70000})
	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})

	if len(manager.Employees) != 2 {
		t.Errorf("Expected 2 employees, got %d", len(manager.Employees))
	}

	// Edge Case: Add an Employee with an Existing ID
	manager.AddEmployee(Employee{ID: 1, Name: "Derek", Age: 40, Salary: 80000})
	if len(manager.Employees) != 3 {
		t.Errorf("Expected 3 employees despite duplicate ID, got %d", len(manager.Employees))
	}
}

func TestRemoveEmployee(t *testing.T) {
	manager := Manager{}

	manager.AddEmployee(Employee{ID: 1, Name: "Alice", Age: 30, Salary: 70000})
	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})

	manager.RemoveEmployee(1)
	if len(manager.Employees) != 1 {
		t.Errorf("Expected 1 employee after removing ID 1, got %d", len(manager.Employees))
	}

	// Remove Non-Existing Employee
	manager.RemoveEmployee(999)
	if len(manager.Employees) != 1 {
		t.Errorf("Expected 1 employee after attempting to remove non-existing ID, got %d", len(manager.Employees))
	}
}

func TestGetAverageSalary(t *testing.T) {
	manager := Manager{}

	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})
	manager.AddEmployee(Employee{ID: 3, Name: "Charlie", Age: 35, Salary: 75000})

	var expectedAverage float64
	expectedAverage = (65000 + 75000) / 2
	if avg := manager.GetAverageSalary(); avg != expectedAverage {
		t.Errorf("Expected average salary %f, got %f", expectedAverage, avg)
	}

	// Edge Case: GetAverageSalary with No Employees
	manager.Employees = []Employee{}

	expectedAverage = 0

	if avg := manager.GetAverageSalary(); avg != expectedAverage {
		t.Errorf("Expected average salary %f with no employees, got %f", expectedAverage, avg)
	}
}

func TestFindEmployeeByID(t *testing.T) {
	manager := Manager{}

	manager.AddEmployee(Employee{ID: 2, Name: "Bob", Age: 25, Salary: 65000})
	manager.AddEmployee(Employee{ID: 3, Name: "Charlie", Age: 35, Salary: 75000})

	employee := manager.FindEmployeeByID(2)
	if employee == nil || employee.Name != "Bob" {
		t.Errorf("Expected to find Bob, got %+v", employee)
	}

	employee = manager.FindEmployeeByID(999)
	if employee != nil {
		t.Errorf("Expected no employee, got %+v", employee)
	}
}

func TestFindEmployeeAfterRemoval(t *testing.T) {
	manager := Manager{}

	manager.AddEmployee(Employee{ID: 4, Name: "David", Age: 45, Salary: 100000})
	manager.AddEmployee(Employee{ID: 5, Name: "Eva", Age: 29, Salary: 95000})

	employee := manager.FindEmployeeByID(4)
	if employee == nil || employee.Name != "David" {
		t.Errorf("Expected to find David, got %+v", employee)
	}

	// Edge Case: Test FindEmployeeByID After Removal
	manager.RemoveEmployee(4)
	employee = manager.FindEmployeeByID(4)
	if employee != nil {
		t.Errorf("Expected no employee, got %+v", employee)
	}
}
