# Hints for GORM CRUD Operations Challenge

## Hint 1: Database Connection Setup

Start with the database connection - Make sure your `ConnectDB()` function properly connects to SQLite and auto-migrates the User model.

Use `gorm.Open()` with SQLite driver, call `AutoMigrate(&User{})` to create the table, and return the database connection and any error.

```go
import "gorm.io/driver/sqlite"

func ConnectDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    err = db.AutoMigrate(&User{})
    return db, err
}
```

## Hint 2: Creating Users

Use `db.Create(user)` to insert the user. Check for errors after the operation. The user's ID will be automatically set after creation.

```go
func CreateUser(db *gorm.DB, user *User) error {
    result := db.Create(user)
    return result.Error
}
```

## Hint 3: Reading Users by ID

Use `db.First(&user, id)` to find a user by ID. Handle the case where user doesn't exist and return a pointer to the user.

```go
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
```

## Hint 4: Reading All Users

Use `db.Find(&users)` to get all users. Return a slice of users and handle empty results (return empty slice, not nil).

```go
func GetAllUsers(db *gorm.DB) ([]User, error) {
    var users []User
    result := db.Find(&users)
    return users, result.Error
}
```

## Hint 5: Updating Users

Use `db.Save(user)` to update the user. Make sure the user has a valid ID and handle the case where user doesn't exist.

```go
func UpdateUser(db *gorm.DB, user *User) error {
    result := db.Save(user)
    return result.Error
}
```

## Hint 6: Deleting Users

Use `db.Delete(&User{}, id)` to delete by ID. Handle the case where user doesn't exist and return appropriate error messages.

```go
func DeleteUser(db *gorm.DB, id uint) error {
    result := db.Delete(&User{}, id)
    return result.Error
}
```

## Common Patterns

### Database Connection
```go
func ConnectDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    err = db.AutoMigrate(&User{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

### Error Handling
```go
func CreateUser(db *gorm.DB, user *User) error {
    result := db.Create(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
```

### Not Found Handling
```go
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
```

## Validation and Constraints

### Model Validation
Your User model has these constraints:
- `Name`: Required (not null)
- `Email`: Required and unique
- `Age`: Must be greater than 0

### Testing Validation
The tests will check:
- Creating user with invalid age (negative)
- Creating user with duplicate email
- All CRUD operations work correctly

## Common Mistakes to Avoid

1. **Not handling errors** - Always check for errors after database operations
2. **Returning nil instead of empty slice** - For `GetAllUsers()`, return empty slice if no users
3. **Not checking if user exists** - For update/delete operations, verify user exists first
4. **Forgetting to auto-migrate** - Make sure to call `AutoMigrate()` in `ConnectDB()`
5. **Not using pointers** - Return pointers to User structs, not values

## Testing Tips

1. **Clean up after tests** - Always clean up test data
2. **Test edge cases** - Test with invalid data, non-existent users, etc.
3. **Verify constraints** - Make sure validation works correctly

## Debugging

1. **Enable GORM logging** to see SQL queries:
```go
db = db.Debug()
```

2. **Check table structure** after migration:
```go
// Verify table exists
assert.True(t, db.Migrator().HasTable(&User{}))
```

3. **Check data in database**:
```go
// Print all users
var users []User
db.Find(&users)
for _, user := range users {
    fmt.Printf("User: %+v\n", user)
}
```

## Performance Considerations

1. **Use appropriate methods** - Use `First()` for single records, `Find()` for multiple
2. **Handle large datasets** - Consider pagination for large result sets
3. **Use transactions** - For multiple related operations

## Useful GORM Methods

- `db.Create()` - Create records
- `db.First()` - Get first record
- `db.Find()` - Get multiple records
- `db.Save()` - Update records
- `db.Delete()` - Delete records
- `db.Where()` - Filter results
- `db.AutoMigrate()` - Migrate models

## SQLite Specific Notes

- SQLite is used for this challenge, so some SQL syntax might be different from other databases
- SQLite has good support for all basic operations
- Use `gorm.io/driver/sqlite` for the driver

## Final Tips

1. **Read the tests carefully** - They show exactly what your functions should do
2. **Start simple** - Get basic CRUD working first, then add validation
3. **Test incrementally** - Test each function as you implement it
4. **Use the learning resources** - Check the GORM documentation for detailed examples

## Common Error Messages

- `UNIQUE constraint failed` - Email already exists
- `CHECK constraint failed` - Age is invalid
- `record not found` - User doesn't exist
- `database is locked` - SQLite file access issue

## Code Structure Example

```go
package main

import (
    "time"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    Age       int       `gorm:"check:age > 0"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func ConnectDB() (*gorm.DB, error) {
    // TODO: Implement database connection
    return nil, nil
}

func CreateUser(db *gorm.DB, user *User) error {
    // TODO: Implement user creation
    return nil
}

// ... other functions
```

Remember to implement each function step by step and test thoroughly! 