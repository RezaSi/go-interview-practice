package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert" // Added for testing
	"gorm.io/gorm"
)

func TestConnectDB(t *testing.T) {
	db, err := ConnectDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test that all tables are created
	assert.True(t, db.Migrator().HasTable(&User{}))
	assert.True(t, db.Migrator().HasTable(&Company{}))
	assert.True(t, db.Migrator().HasTable(&Post{}))

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()
	os.Remove("test.db")
}

func TestCreateUser(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	err := CreateUser(ctx, db, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Verify user was created using traditional API for verification
	var foundUser User
	db.First(&foundUser, user.ID)
	assert.Equal(t, "John Doe", foundUser.Name)
	assert.Equal(t, "john@example.com", foundUser.Email)
	assert.Equal(t, 25, foundUser.Age)
}

func TestGetUserByID(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create a user first
	user := &User{Name: "Jane Doe", Email: "jane@example.com", Age: 30}
	CreateUser(ctx, db, user)

	// Test GetUserByID
	foundUser, err := GetUserByID(ctx, db, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, "Jane Doe", foundUser.Name)
	assert.Equal(t, "jane@example.com", foundUser.Email)
	assert.Equal(t, 30, foundUser.Age)

	// Test with non-existent ID
	_, err = GetUserByID(ctx, db, 999)
	assert.Error(t, err)
}

func TestUpdateUserAge(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create a user first
	user := &User{Name: "Bob Smith", Email: "bob@example.com", Age: 25}
	CreateUser(ctx, db, user)

	// Update user age
	err := UpdateUserAge(ctx, db, user.ID, 30)
	assert.NoError(t, err)

	// Verify update
	updatedUser, _ := GetUserByID(ctx, db, user.ID)
	assert.Equal(t, 30, updatedUser.Age)
}

func TestDeleteUser(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create a user first
	user := &User{Name: "Charlie Brown", Email: "charlie@example.com", Age: 35}
	CreateUser(ctx, db, user)

	// Delete user
	err := DeleteUser(ctx, db, user.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = GetUserByID(ctx, db, user.ID)
	assert.Error(t, err)
}

func TestCreateUsersInBatches(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	users := []User{
		{Name: "User1", Email: "user1@example.com", Age: 20},
		{Name: "User2", Email: "user2@example.com", Age: 25},
		{Name: "User3", Email: "user3@example.com", Age: 30},
		{Name: "User4", Email: "user4@example.com", Age: 35},
		{Name: "User5", Email: "user5@example.com", Age: 40},
	}

	err := CreateUsersInBatches(ctx, db, users, 3)
	assert.NoError(t, err)

	// Verify all users were created
	for _, user := range users {
		assert.NotZero(t, user.ID)
	}

	// Count total users
	var count int64
	db.Model(&User{}).Count(&count)
	assert.Equal(t, int64(5), count)
}

func TestFindUsersByAgeRange(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create test users
	users := []User{
		{Name: "Young User", Email: "young@example.com", Age: 18},
		{Name: "Middle User1", Email: "middle1@example.com", Age: 25},
		{Name: "Middle User2", Email: "middle2@example.com", Age: 30},
		{Name: "Old User", Email: "old@example.com", Age: 50},
	}
	CreateUsersInBatches(ctx, db, users, 4)

	// Test age range query
	foundUsers, err := FindUsersByAgeRange(ctx, db, 20, 35)
	assert.NoError(t, err)
	assert.Len(t, foundUsers, 2)

	// Verify correct users found
	ages := make([]int, len(foundUsers))
	for i, user := range foundUsers {
		ages[i] = user.Age
	}
	assert.Contains(t, ages, 25)
	assert.Contains(t, ages, 30)
}

func TestUpsertUser(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Test insert (first time)
	user := &User{Name: "New User", Email: "new@example.com", Age: 25}
	err := UpsertUser(ctx, db, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Test update (same email)
	userUpdate := &User{Name: "Updated User", Email: "new@example.com", Age: 30}
	err = UpsertUser(ctx, db, userUpdate)
	assert.NoError(t, err)

	// Verify only one record exists and it's updated
	var count int64
	db.Model(&User{}).Where("email = ?", "new@example.com").Count(&count)
	assert.Equal(t, int64(1), count)

	// Check the name was updated
	var foundUser User
	db.Where("email = ?", "new@example.com").First(&foundUser)
	assert.Equal(t, "Updated User", foundUser.Name)
	assert.Equal(t, 30, foundUser.Age)
}

func TestCreateUserWithResult(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	user := &User{Name: "Result User", Email: "result@example.com", Age: 25}
	rowsAffected, err := CreateUserWithResult(ctx, db, user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
	assert.NotZero(t, user.ID)
}

func TestGetUsersWithCompany(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create company
	company := &Company{Name: "Tech Corp", Industry: "Technology", FoundedYear: 2020}
	db.Create(company)

	// Create users with company
	companyID := company.ID
	users := []User{
		{Name: "Employee1", Email: "emp1@example.com", Age: 25, CompanyID: &companyID},
		{Name: "Employee2", Email: "emp2@example.com", Age: 30, CompanyID: &companyID},
		{Name: "Freelancer", Email: "free@example.com", Age: 35}, // No company
	}
	CreateUsersInBatches(ctx, db, users, 3)

	// Test getting users with company
	usersWithCompany, err := GetUsersWithCompany(ctx, db)
	assert.NoError(t, err)
	assert.Len(t, usersWithCompany, 2) // Only users with company

	// Verify company data is loaded
	for _, user := range usersWithCompany {
		assert.NotNil(t, user.Company)
		assert.Equal(t, "Tech Corp", user.Company.Name)
	}
}

func TestGetUsersWithPosts(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create user
	user := &User{Name: "Blogger", Email: "blogger@example.com", Age: 30}
	CreateUser(ctx, db, user)

	// Create posts
	posts := []Post{
		{Title: "Post 1", Content: "Content 1", UserID: user.ID},
		{Title: "Post 2", Content: "Content 2", UserID: user.ID},
		{Title: "Post 3", Content: "Content 3", UserID: user.ID},
		{Title: "Post 4", Content: "Content 4", UserID: user.ID},
	}
	for _, post := range posts {
		db.Create(&post)
	}

	// Test preloading with limit
	users, err := GetUsersWithPosts(ctx, db, 2)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Len(t, users[0].Posts, 2) // Limited to 2 posts
}

func TestGetUserWithPostsAndCompany(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create company
	company := &Company{Name: "Blog Corp", Industry: "Media", FoundedYear: 2019}
	db.Create(company)

	// Create user with company
	companyID := company.ID
	user := &User{Name: "Writer", Email: "writer@example.com", Age: 28, CompanyID: &companyID}
	CreateUser(ctx, db, user)

	// Create posts
	posts := []Post{
		{Title: "Article 1", Content: "Content 1", UserID: user.ID},
		{Title: "Article 2", Content: "Content 2", UserID: user.ID},
	}
	for _, post := range posts {
		db.Create(&post)
	}

	// Test multiple preloads
	foundUser, err := GetUserWithPostsAndCompany(ctx, db, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Len(t, foundUser.Posts, 2)
	assert.NotNil(t, foundUser.Company)
	assert.Equal(t, "Blog Corp", foundUser.Company.Name)
}

func TestSearchUsersInCompany(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create companies
	tech := &Company{Name: "Tech Solutions", Industry: "Technology", FoundedYear: 2020}
	finance := &Company{Name: "Finance Corp", Industry: "Finance", FoundedYear: 2018}
	db.Create(tech)
	db.Create(finance)

	// Create users
	techID := tech.ID
	financeID := finance.ID
	users := []User{
		{Name: "Tech Worker 1", Email: "tech1@example.com", Age: 25, CompanyID: &techID},
		{Name: "Tech Worker 2", Email: "tech2@example.com", Age: 30, CompanyID: &techID},
		{Name: "Finance Worker", Email: "finance1@example.com", Age: 35, CompanyID: &financeID},
	}
	CreateUsersInBatches(ctx, db, users, 3)

	// Test searching users in specific company
	techUsers, err := SearchUsersInCompany(ctx, db, "Tech Solutions")
	assert.NoError(t, err)
	assert.Len(t, techUsers, 2)

	financeUsers, err := SearchUsersInCompany(ctx, db, "Finance Corp")
	assert.NoError(t, err)
	assert.Len(t, financeUsers, 1)
}

func TestGetTopActiveUsers(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	ctx := context.Background()

	// Create users
	users := []User{
		{Name: "Active User", Email: "active@example.com", Age: 25},
		{Name: "Moderate User", Email: "moderate@example.com", Age: 30},
		{Name: "Inactive User", Email: "inactive@example.com", Age: 35},
	}
	CreateUsersInBatches(ctx, db, users, 3)

	// Create posts
	posts := []Post{
		// Active user - 3 posts
		{Title: "Post 1", Content: "Content 1", UserID: users[0].ID},
		{Title: "Post 2", Content: "Content 2", UserID: users[0].ID},
		{Title: "Post 3", Content: "Content 3", UserID: users[0].ID},
		// Moderate user - 1 post
		{Title: "Post 4", Content: "Content 4", UserID: users[1].ID},
		// Inactive user - 0 posts
	}
	for _, post := range posts {
		db.Create(&post)
	}

	// Test getting top active users
	topUsers, err := GetTopActiveUsers(ctx, db, 2)
	assert.NoError(t, err)
	assert.Len(t, topUsers, 2)

	// Verify order (most active first)
	assert.Equal(t, "Active User", topUsers[0].Name)
	assert.Equal(t, "Moderate User", topUsers[1].Name)
}

func TestContextTimeout(t *testing.T) {
	db, _ := ConnectDB()
	defer cleanup(db)

	// Test with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	user := &User{Name: "Timeout User", Email: "timeout@example.com", Age: 25}
	err := CreateUser(ctx, db, user)

	// Should get context deadline exceeded error
	assert.Error(t, err)
}

// Helper function to clean up database
func cleanup(db *gorm.DB) {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}
}
