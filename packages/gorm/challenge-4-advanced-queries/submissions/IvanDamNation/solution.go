package main

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User represents a user in the social media system
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Age       int    `gorm:"not null"`
	Country   string `gorm:"not null"`
	CreatedAt time.Time
	Posts     []Post `gorm:"foreignKey:UserID"`
	Likes     []Like `gorm:"foreignKey:UserID"`
}

// Post represents a social media post
type Post struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Content     string `gorm:"type:text"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
	Category    string `gorm:"not null"`
	ViewCount   int    `gorm:"default:0"`
	IsPublished bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Likes       []Like `gorm:"foreignKey:PostID"`
}

// Like represents a user's like on a post
type Like struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	PostID    uint `gorm:"not null"`
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
	CreatedAt time.Time
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// ConnectDB establishes a connection to the SQLite database with auto-migration
func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &Post{}, &Like{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetTopUsersByPostCount retrieves users with the most posts
func GetTopUsersByPostCount(db *gorm.DB, limit int) ([]User, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit value")
	}

	var users []User
	res := db.
		Joins("LEFT JOIN posts ON posts.user_id = users.id").
		Group("users.id").
		Order("COUNT(posts.id) DESC").
		Limit(limit).
		Find(&users)

	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

// GetPostsByCategoryWithUserInfo retrieves posts by category with pagination and user info
func GetPostsByCategoryWithUserInfo(db *gorm.DB, category string, page, pageSize int) ([]Post, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, fmt.Errorf("invalid value of page or page size")
	}

	tx := db.Begin()
	defer tx.Rollback()

	var total int64
	err := tx.
		Model(&Post{}).
		Where("category = ?", category).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	var posts []Post
	offset := (page - 1) * pageSize
	err = tx.
		Where("category = ?", category).
		Preload("User").
		Order("created_at DESC, id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, tx.Commit().Error
}

// GetUserEngagementStats calculates engagement statistics for a user
func GetUserEngagementStats(db *gorm.DB, userID uint) (map[string]any, error) {
	stats := make(map[string]any, 4)

	tx := db.Begin()
	defer tx.Rollback()

	var user User
	err := tx.First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	var totalPosts int
	err = tx.
		Model(&Post{}).
		Where("user_id = ?", userID).
		Select("COUNT(*)").
		Scan(&totalPosts).Error

	if err != nil {
		return nil, err
	}

	var totalLikesGiven int64
	err = tx.
		Model(&Like{}).
		Where("user_id = ?", userID).
		Count(&totalLikesGiven).Error

	if err != nil {
		return nil, err
	}

	var totalLikesReceived int64
	err = tx.
		Model(&Like{}).
		Joins("LEFT JOIN posts ON posts.id = likes.post_id").
		Where("posts.user_id = ?", userID).
		Count(&totalLikesReceived).Error

	if err != nil {
		return nil, err
	}

	var avgViews float64
	err = tx.
		Model(&Post{}).
		Where("user_id = ?", userID).
		Select("COALESCE(AVG(view_count), 0)").
		Scan(&avgViews).Error

	if err != nil {
		return nil, err
	}

	stats["total_posts"] = totalPosts
	stats["total_likes_given"] = totalLikesGiven
	stats["total_likes_received"] = totalLikesReceived
	stats["average_views_per_post"] = avgViews

	return stats, tx.Commit().Error
}

// GetPopularPostsByLikes retrieves popular posts by likes in a time period
func GetPopularPostsByLikes(db *gorm.DB, days int, limit int) ([]Post, error) {
	if days <= 0 || limit <= 0 {
		return nil, fmt.Errorf("invalid parameter for days or limit")
	}

	targetDate := time.Now().AddDate(0, 0, -days)

	var posts []Post
	err := db.
		Model(&Post{}).
		Where("posts.created_at >= ?", targetDate).
		Joins("LEFT JOIN likes ON posts.id = likes.post_id").
		Group("posts.id").
		Order("COUNT(likes.id) DESC").
		Limit(limit).
		Preload("Likes").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// GetCountryUserStats retrieves user statistics grouped by country
func GetCountryUserStats(db *gorm.DB) ([]map[string]any, error) {
	var stats []map[string]any
	err := db.
		Model(&User{}).
		Select("country, COUNT(*) AS user_count").
		Group("country").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// SearchPostsByContent searches posts by content using full-text search
func SearchPostsByContent(db *gorm.DB, query string, limit int) ([]Post, error) {
	if query == "" || limit <= 0 {
		return nil, fmt.Errorf("empty query or invalid limit value")
	}

	var posts []Post
	err := db.
		Where("content LIKE ? ESCAPE '\\'", "%"+likeEscaper.Replace(query)+"%").
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// GetUserRecommendations retrieves user recommendations based on similar interests
func GetUserRecommendations(db *gorm.DB, userID uint, limit int) ([]User, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit value")
	}

	likedPostsUsersSubQuery := db.
		Model(&Like{}).
		Select("user_id").
		Where("post_id IN (?)", db.
			Model(&Like{}).
			Select("post_id").
			Where("user_id = ?", userID),
		)

	if likedPostsUsersSubQuery.Error != nil {
		return nil, likedPostsUsersSubQuery.Error
	}

	var recommendations []User
	err := db.
		Where("id IN (?) AND id != ?", likedPostsUsersSubQuery, userID).
		Limit(limit).
		Find(&recommendations).Error

	if err != nil {
		return nil, err
	}

	return recommendations, nil
}
