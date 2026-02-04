package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestConnectDB(t *testing.T) {
	db, err := ConnectDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test that tables are created
	assert.True(t, db.Migrator().HasTable(&User{}))
	assert.True(t, db.Migrator().HasTable(&Post{}))
	assert.True(t, db.Migrator().HasTable(&Like{}))

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()
	os.Remove("test.db")
}

func setupTestData(db *gorm.DB) {
	// Create test users
	users := []User{
		{Username: "user1", Email: "user1@test.com", Age: 25, Country: "USA"},
		{Username: "user2", Email: "user2@test.com", Age: 30, Country: "Canada"},
		{Username: "user3", Email: "user3@test.com", Age: 28, Country: "USA"},
		{Username: "user4", Email: "user4@test.com", Age: 35, Country: "UK"},
	}
	for i := range users {
		db.Create(&users[i])
	}

	// Create test posts
	posts := []Post{
		{Title: "Post 1", Content: "Content about technology", UserID: users[0].ID, Category: "tech", ViewCount: 100},
		{Title: "Post 2", Content: "Content about sports", UserID: users[0].ID, Category: "sports", ViewCount: 50},
		{Title: "Post 3", Content: "Content about food", UserID: users[1].ID, Category: "food", ViewCount: 75},
		{Title: "Post 4", Content: "Content about travel", UserID: users[1].ID, Category: "travel", ViewCount: 200},
		{Title: "Post 5", Content: "Content about music", UserID: users[2].ID, Category: "music", ViewCount: 150},
		{Title: "Post 6", Content: "Content about movies", UserID: users[3].ID, Category: "entertainment", ViewCount: 300},
	}
	for i := range posts {
		db.Create(&posts[i])
	}

	// Create test likes
	likes := []Like{
		{UserID: users[1].ID, PostID: posts[0].ID},
		{UserID: users[2].ID, PostID: posts[0].ID},
		{UserID: users[3].ID, PostID: posts[0].ID},
		{UserID: users[0].ID, PostID: posts[1].ID},
		{UserID: users[2].ID, PostID: posts[1].ID},
		{UserID: users[0].ID, PostID: posts[2].ID},
		{UserID: users[3].ID, PostID: posts[2].ID},
		{UserID: users[1].ID, PostID: posts[3].ID},
		{UserID: users[2].ID, PostID: posts[3].ID},
		{UserID: users[3].ID, PostID: posts[3].ID},
	}
	for i := range likes {
		db.Create(&likes[i])
	}
}

func TestGetTopUsersByPostCount(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	setupTestData(db)

	// Test getting top users by post count
	users, err := GetTopUsersByPostCount(db, 3)
	assert.NoError(t, err)
	assert.Len(t, users, 3)

	// Verify users are ordered by post count (descending)
	// user1 and user2 should have 2 posts each, user3 and user4 should have 1 each
	assert.Equal(t, "user1", users[0].Username) // or user2, both have 2 posts
	assert.Equal(t, "user2", users[1].Username) // or user1, both have 2 posts
}

func TestGetPostsByCategoryWithUserInfo(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	setupTestData(db)

	// Test paginated posts retrieval
	posts, total, err := GetPostsByCategoryWithUserInfo(db, "tech", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Post 1", posts[0].Title)
	assert.Equal(t, "user1", posts[0].User.Username)

	// Test pagination
	posts, total, err = GetPostsByCategoryWithUserInfo(db, "sports", 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, posts, 1)
}

func TestGetUserEngagementStats(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	setupTestData(db)

	// Get first user
	var user User
	db.First(&user)

	// Test engagement statistics
	stats, err := GetUserEngagementStats(db, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Verify stats contain expected keys
	assert.Contains(t, stats, "total_posts")
	assert.Contains(t, stats, "total_likes_received")
	assert.Contains(t, stats, "total_likes_given")
	assert.Contains(t, stats, "average_views_per_post")
}

func TestGetPopularPostsByLikes(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	setupTestData(db)

	// Test popular posts by likes
	posts, err := GetPopularPostsByLikes(db, 30, 5)
	assert.NoError(t, err)
	assert.Len(t, posts, 5)

	// Verify posts are ordered by like count (descending)
	// Post 3 should have 3 likes, Post 1 should have 3 likes, etc.
	assert.GreaterOrEqual(t, len(posts[0].Likes), len(posts[1].Likes))
}

func TestGetCountryUserStats(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	setupTestData(db)

	// Test country-based user statistics
	stats, err := GetCountryUserStats(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, stats)

	// Verify stats contain expected countries
	countries := make([]string, 0)
	for _, stat := range stats {
		if country, ok := stat["country"].(string); ok {
			countries = append(countries, country)
		}
	}
	assert.Contains(t, countries, "USA")
	assert.Contains(t, countries, "Canada")
	assert.Contains(t, countries, "UK")
}

func TestSearchPostsByContent(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	setupTestData(db)

	// Test content search
	posts, err := SearchPostsByContent(db, "technology", 5)
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Post 1", posts[0].Title)

	// Test search with multiple results
	posts, err = SearchPostsByContent(db, "Content", 5)
	assert.NoError(t, err)
	// There are 6 posts containing "Content" but we queried with a limit of 5
	assert.Len(t, posts, 5)
}

func TestGetUserRecommendations(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	setupTestData(db)

	// Get first user
	var user User
	db.First(&user)

	// Test user recommendations
	recommendations, err := GetUserRecommendations(db, user.ID, 3)
	assert.NoError(t, err)
	assert.Len(t, recommendations, 2)

	// Verify recommendations don't include the user themselves
	for _, rec := range recommendations {
		assert.NotEqual(t, user.ID, rec.ID)
	}
}

func TestErrorHandling(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Test with non-existent user
	_, err := GetUserEngagementStats(db, 999)
	assert.Error(t, err)

	// Test with non-existent category
	_, _, err = GetPostsByCategoryWithUserInfo(db, "non-existent", 1, 10)
	assert.NoError(t, err) // Should return empty results, not error

	// Test with invalid pagination
	_, _, err = GetPostsByCategoryWithUserInfo(db, "tech", -1, 10)
	assert.Error(t, err)
}
