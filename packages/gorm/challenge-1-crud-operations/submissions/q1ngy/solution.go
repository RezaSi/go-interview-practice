package main

import (
	"errors"
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
	// Implement database connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	return db, nil
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user *User) error {
	// Implement user creation
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	// Implement user retrieval by ID
	user := &User{}
	result := db.First(user, id)
	if result.Error != nil {
		return &User{}, result.Error
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
	// Implement retrieval of all users
	var users []User
	db.Find(&users)
	return users, nil
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
	// Implement user update
	result := db.Model(user).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not fount")
	}
	return nil
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
	// Implement user deletion
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not fount")
	}
	return nil
}
