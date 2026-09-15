package main

import (
	"context"
	"fmt"
	"strings"
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
	// TODO: Implement database connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		// log.Fatalf("Error occured:" err)
		return nil, err
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	return db, err
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user *User) error {
	// TODO: Implement user creation
	if err := db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("email already exists")
		}
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	// TODO: Implement user retrieval by ID
	var user User
	if err := db.WithContext(context.Background()).Model(&User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
	// TODO: Implement retrieval of all users
	var users []User
	err := db.WithContext(context.Background()).Model(&User{}).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, err
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
	// TODO: Implement user update
	if (*user).Age < 0 {
		return fmt.Errorf("age cannot be negative")
	}

	if (*user).Email == "" || strings.Contains((*user).Email, "@") == false {
		return fmt.Errorf("invalid email format")
	}
	user_id := (*user).ID
	result := db.WithContext(context.Background()).Model(&user).Where("id = ?", user_id).Updates(*user)
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return result.Error
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
	// TODO: Implement user deletion
	if id < 1 {
		return fmt.Errorf("invalid user ID")
	}
	result := db.WithContext(context.Background()).Where("id = ?", id).Delete(&User{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return result.Error
}
