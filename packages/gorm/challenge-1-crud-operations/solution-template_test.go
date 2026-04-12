package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	db, err := ConnectDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test that table is created
	assert.True(t, db.Migrator().HasTable(&User{}))

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()
	os.Remove("test.db")
}

func TestCreateUser(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	err := CreateUser(db, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Verify user was created
	var foundUser User
	db.First(&foundUser, user.ID)
	assert.Equal(t, "John Doe", foundUser.Name)
	assert.Equal(t, "john@example.com", foundUser.Email)
	assert.Equal(t, 25, foundUser.Age)
}

func TestGetUserByID(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create a user first
	user := &User{Name: "Jane Doe", Email: "jane@example.com", Age: 30}
	CreateUser(db, user)

	// Test retrieval
	retrievedUser, err := GetUserByID(db, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", retrievedUser.Name)
	assert.Equal(t, "jane@example.com", retrievedUser.Email)
	assert.Equal(t, 30, retrievedUser.Age)
}

func TestGetUserByIDNotFound(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Test getting non-existent user
	_, err := GetUserByID(db, 999)
	assert.Error(t, err)
}

func TestGetAllUsers(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create multiple users
	users := []User{
		{Name: "User 1", Email: "user1@example.com", Age: 25},
		{Name: "User 2", Email: "user2@example.com", Age: 30},
		{Name: "User 3", Email: "user3@example.com", Age: 35},
	}

	for i := range users {
		CreateUser(db, &users[i])
	}

	// Test retrieval of all users
	retrievedUsers, err := GetAllUsers(db)
	assert.NoError(t, err)
	assert.Len(t, retrievedUsers, 3)

	// Verify all users are present
	emails := make([]string, len(retrievedUsers))
	for i, user := range retrievedUsers {
		emails[i] = user.Email
	}
	assert.Contains(t, emails, "user1@example.com")
	assert.Contains(t, emails, "user2@example.com")
	assert.Contains(t, emails, "user3@example.com")
}

func TestUpdateUser(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create a user
	user := &User{Name: "Original Name", Email: "original@example.com", Age: 25}
	CreateUser(db, user)

	// Update user
	user.Name = "Updated Name"
	user.Age = 30
	err := UpdateUser(db, user)
	assert.NoError(t, err)

	// Verify update
	var updatedUser User
	db.First(&updatedUser, user.ID)
	assert.Equal(t, "Updated Name", updatedUser.Name)
	assert.Equal(t, 30, updatedUser.Age)
	assert.Equal(t, "original@example.com", updatedUser.Email) // Email should remain unchanged
}

func TestUpdateUserNotFound(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Try to update non-existent user
	user := &User{ID: 999, Name: "Test", Email: "test@example.com", Age: 25}
	err := UpdateUser(db, user)
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create a user
	user := &User{Name: "To Delete", Email: "delete@example.com", Age: 25}
	CreateUser(db, user)
	userID := user.ID

	// Delete user
	err := DeleteUser(db, userID)
	assert.NoError(t, err)

	// Verify user is deleted
	var deletedUser User
	err = db.First(&deletedUser, userID).Error
	assert.Error(t, err) // Should return error as user doesn't exist
}

func TestDeleteUserNotFound(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Try to delete non-existent user
	err := DeleteUser(db, 999)
	assert.Error(t, err)
}

func TestUserValidation(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Test creating user with invalid age
	invalidUser := &User{
		Name:  "Invalid User",
		Email: "invalid@example.com",
		Age:   -5, // Invalid age
	}

	err := CreateUser(db, invalidUser)
	assert.Error(t, err) // Should fail due to age constraint
}

func TestUniqueEmailConstraint(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create first user
	user1 := &User{Name: "User 1", Email: "same@example.com", Age: 25}
	err := CreateUser(db, user1)
	assert.NoError(t, err)

	// Try to create second user with same email
	user2 := &User{Name: "User 2", Email: "same@example.com", Age: 30}
	err = CreateUser(db, user2)
	assert.Error(t, err) // Should fail due to unique email constraint
}

func TestCRUDOperations(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create
	user := &User{Name: "Test User", Email: "test@example.com", Age: 25}
	err := CreateUser(db, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Read
	retrievedUser, err := GetUserByID(db, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", retrievedUser.Name)

	// Update
	user.Name = "Updated Test User"
	err = UpdateUser(db, user)
	assert.NoError(t, err)

	// Verify update
	updatedUser, err := GetUserByID(db, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Test User", updatedUser.Name)

	// Delete
	err = DeleteUser(db, user.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = GetUserByID(db, user.ID)
	assert.Error(t, err)
}
