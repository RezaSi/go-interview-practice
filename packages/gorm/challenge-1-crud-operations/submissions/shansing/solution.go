package main

import (
	"errors"
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
	// Implement database connection
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
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
	// Implement retrieval of all users
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
	// Implement user update
	return db.Transaction(func(tx *gorm.DB) error {
		if u := tx.First(&User{}, user.ID); u.Error != nil {
			if errors.Is(u.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user with id %d not found", user.ID)
			}
			return fmt.Errorf("user with id %d error", user.ID)
		}
		result := tx.Save(user)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
	// Implement user deletion
	return db.Transaction(func(tx *gorm.DB) error {
		if u := tx.First(&User{}, id); u.Error != nil {
			if errors.Is(u.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user with id %d not found", id)
			}
			return fmt.Errorf("user with id %d error", id)
		}
		result := tx.Delete(&User{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("user with id %d not found", id)
		}
		return nil
	})
}
