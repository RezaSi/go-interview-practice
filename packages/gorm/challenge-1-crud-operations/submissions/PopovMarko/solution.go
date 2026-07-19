package main

import (
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

// ConnectDB opens a connection to the SQLite database file "test.db" and runs
// the automatic migration for the User model, creating the table if it does not
// already exist. It returns the ready-to-use *gorm.DB handle, or an error if the
// connection or migration fails.
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

// CreateUser inserts the given user into the database. On success the user's
// auto-generated fields (ID, CreatedAt, UpdatedAt) are populated in place. It
// returns an error if the insert violates a constraint (e.g. duplicate email
// or non-positive age) or otherwise fails.
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

// GetUserByID retrieves the user with the given primary key. It returns a
// pointer to the found user, or a nil user and an error if no matching record
// exists (gorm.ErrRecordNotFound) or the query fails.
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	user := User{}
	res := db.First(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

// GetAllUsers retrieves every user stored in the database. It returns the slice
// of users (empty if the table has no rows) together with any error from the
// query.
func GetAllUsers(db *gorm.DB) ([]User, error) {
	users := []User{}
	return users, db.Find(&users).Error
}

// UpdateUser updates the non-zero fields of the given user, identified by its
// ID. Because it uses a struct-based update, only fields set to a non-zero value
// are written. It returns an error if the update fails, or "user not found" if
// no record matches the user's ID.
func UpdateUser(db *gorm.DB, user *User) error {
	res := db.Model(user).Updates(user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// DeleteUser removes the user with the given primary key from the database. It
// returns an error if the delete fails, or "user not found" if no record
// matched the given ID.
func DeleteUser(db *gorm.DB, id uint) error {
	res := db.Delete(&User{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
