package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Age       int    `gorm:"check:age > 0"` //Built-in Validation
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	return result.Error
}
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id) // Find by primary key
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	return users, result.Error
}

func UpdateUser(db *gorm.DB, user *User) error {
	_, err := GetUserByID(db, user.ID)
	if err != nil {
		return err
	}
	result := db.Save(user) // Update by primary key
	return result.Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	_, err := GetUserByID(db, id)
	if err != nil {
		return err
	}
	result := db.Delete(&User{}, id) // Delete by primary key
	return result.Error
}