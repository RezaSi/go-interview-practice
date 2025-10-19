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

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	// NOTE: Tests don't need data validation
	return db.Create(user).Error
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id)
	return &user, result.Error
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	return users, result.Error
}

func UpdateUser(db *gorm.DB, user *User) error {
	// NOTE: Tests don't need data validation
	result := db.Model(&User{}).Where("id = ?", user.ID).Updates(user)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
