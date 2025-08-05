package main

import (
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
	return db.Create(&user).Error
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
    var user User
	result := db.First(&user, id)
	if result.Error != nil {
	    return nil, result.Error
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
    var users []User
    result := db.Find(&users)
	return users, result.Error
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
    var existing User
    err := db.First(&existing, user.ID).Error
    if err != nil {
        return err
    }
	return db.Save(&user).Error
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
    var existing User
    err := db.First(&existing, id).Error
    if err != nil {
        return err
    }
	return db.Delete(&User{}, id).Error
}
