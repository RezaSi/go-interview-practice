package main

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
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
	    return nil, err
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
	    return nil, err
	}
	return db, err
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user *User) error {
	// TODO: Implement user creation
	return db.Create(user).Error
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	// TODO: Implement user retrieval by ID
	var user User
	result := db.First(&user, id)
	return &user, result.Error
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
	// TODO: Implement retrieval of all users
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
	    return nil, result.Error
	}
	return users, result.Error
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
	// TODO: Implement user update
    if err := db.First(&User{}, user.ID).Error; err != nil {
        return err
    }
	return db.Model(&user).Updates(*user).Error
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
	// TODO: Implement user deletion
	if err := db.First(&User{}, id).Error; err != nil {
        return err
    }
	return db.Delete(&User{}, id).Error
}
