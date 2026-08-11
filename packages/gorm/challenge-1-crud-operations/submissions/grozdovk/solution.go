package main

import (
	"time"
	"database/sql" 

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	_ "modernc.org/sqlite"
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
	sqlDB, err := sql.Open("sqlite", "test.db")
    if err != nil {
        return nil, err
    }
    
    db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    err = db.AutoMigrate(&User{})
    return db, err
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	return result.Error
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
    var user User
	result:= db.First(&user, id)
	if result.Error!= nil{
	    return nil, result.Error
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
    var users []User
	result:= db.Find(&users)
	if result.Error!= nil{
	    return nil, result.Error
	}
	return users, nil
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
    var existingUser User
    result := db.First(&existingUser, user.ID)
    if result.Error != nil {
        return result.Error 
    }
	result = db.Save(user)
	return result.Error
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
    var existingUser User
    result := db.First(&existingUser, id)
    if result.Error != nil {
        return result.Error
    }
	result = db.Delete(&User{}, id )
	return result.Error
}
