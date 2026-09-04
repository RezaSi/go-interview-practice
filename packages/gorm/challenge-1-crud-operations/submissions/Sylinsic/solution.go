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
	res := db.Create(user)

	return res.Error
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user *User
	
	res := db.First(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}
	
	return user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	
	res := db.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	if users == nil {
		return []User{}, nil
	}
	
	return users, nil
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
	res := db.Where("ID = ?", user.ID).Updates(user)

	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
	res := db.Delete(&User{}, id)

	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
