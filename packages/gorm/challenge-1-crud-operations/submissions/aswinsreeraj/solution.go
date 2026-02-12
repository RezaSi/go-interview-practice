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
    return db, err
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user *User) error {
	// TODO: Implement user creation
	result := db.Create(&user)
	return result.Error
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
	return users, result.Error
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
	// TODO: Implement user update
	var user2 User
	result := db.First(&user2)
	if result.RowsAffected == 0 {
	    return result.Error
	} else {
	    return db.Save(&user).Error   
	}
	return result.Error
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
	// TODO: Implement user deletion
	var user2 User
	result := db.First(&user2, id)
	if result.RowsAffected == 0 {
	    return result.Error
	} 
	result = db.Delete(&User{}, id)
	return result.Error
}
