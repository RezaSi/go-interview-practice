# Hints for Employee Data Management

## Hint 1: Struct Definition
You need to implement two structs. The `Employee` struct stores individual employee data, and the `Manager` struct contains a slice of employees:
```go
type Manager struct {
    Employees []Employee
}
```

## Hint 2: AddEmployee Method
Use a pointer receiver `(m *Manager)` to modify the slice. Append the new employee to the `Employees` slice:
```go
m.Employees = append(m.Employees, e)
```

## Hint 3: RemoveEmployee Method
To remove an employee by ID, find the employee in the slice and remove it. You can create a new slice excluding the employee with the matching ID.

## Hint 4: Finding Index for Removal
Loop through the slice to find the index of the employee with the matching ID:
```go
for i, emp := range m.Employees {
    if emp.ID == id {
        // Remove element at index i
    }
}
```

## Hint 5: Slice Removal Technique
Remove an element at index `i` from a slice using:
```go
m.Employees = append(m.Employees[:i], m.Employees[i+1:]...)
```

## Hint 6: GetAverageSalary Implementation
Calculate the sum of all salaries and divide by the number of employees. Handle the case when there are no employees to avoid division by zero.

## Hint 7: FindEmployeeByID Return Type
This method should return `*Employee` (a pointer). Return `nil` if the employee is not found, or return the address of the found employee using `&`.

## Hint 8: Returning Employee Pointer
When you find the employee, return a pointer to it:
```go
return &m.Employees[i]
``` 