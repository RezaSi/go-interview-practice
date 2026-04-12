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

// ConnectDB establishes a connection to the SQLite database and migrates the User schema
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
	result := db.Create(user)
	return result.Error
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
	result := db.Model(&User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteUser removes a user from the database by their ID.
func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("user not found for deletion")
	}
	return result.Error
}
