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
	// TODO: Implement user creation
	result := db.Create(user)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	// TODO: Implement user retrieval by ID

	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) ([]User, error) {
	// TODO: Implement retrieval of all users

	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// UpdateUser updates an existing user's information
func UpdateUser(db *gorm.DB, user *User) error {
    // 我们利用刚才写好的 GetUserByID 或者直接查询
    var tempUser User
    // 这里的 Select("id") 是优化，只查 ID 省流量，不查也行
    check := db.First(&tempUser, user.ID)
    
    if check.Error != nil {
        // 如果查不到（Result is gorm.ErrRecordNotFound），直接返回错误
        return check.Error
    }

    // 用户存在，才执行 Save
    result := db.Save(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

// DeleteUser removes a user from the database
func DeleteUser(db *gorm.DB, id uint) error {
	// TODO: Implement user deletion
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// 如果受影响行数是 0，说明 ID 不存在
		return gorm.ErrRecordNotFound
	}

	return nil
}