package main

import (
    "fmt"
	"time"
    "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Age       int    `gorm:"check:age > 0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConnectDB establishes a connection to the SQLite database
func ConnectDB() (*gorm.DB, error) {
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Auto migrate the schema
    err = db.AutoMigrate(&User{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
    
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user *User) error {
    r := db.Create(user)
	return r.Error
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
    u := &User{}
	r := db.First(u, id)
	return u, r.Error
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
    u := make([]User, 0)
	r := db.Find(&u)
	return u, r.Error
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
	r := db.Select("*").Updates(user)
	if r.RowsAffected == 0 {
	    return fmt.Errorf("user not found")
	}

	return r.Error
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
    u := &User{}
	r := db.Delete(u, id)
	if r.RowsAffected == 0 {
	    return fmt.Errorf("user not found")
	}
	return r.Error
}
