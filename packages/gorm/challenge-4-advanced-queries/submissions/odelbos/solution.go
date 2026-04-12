package main

import (
	"time"
	"fmt"
	"errors"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
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

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &Post{}, &Like{}); err != nil {
		return nil, err
	}
	return db, nil
}

func GetTopUsersByPostCount(db *gorm.DB, limit int) ([]User, error) {
	var users []User
	err := db.Model(&User{}).
		Select("users.*, COUNT(posts.id) as post_count").
		Joins("LEFT JOIN posts ON users.id = posts.user_id").
		Group("users.id").
		Order("post_count DESC").
		Limit(limit).
		Scan(&users).Error
	return users, err
}

func GetPostsByCategoryWithUserInfo(db *gorm.DB, category string, page, pageSize int) ([]Post, int64, error) {
    if page <= 0 {
		return nil, 0, errors.New("page cannot be negative or equals to 0")
    }
    if pageSize <= 0 {
		return nil, 0, errors.New("page size cannot be negative or equals to 0")
    }

	var posts []Post
	var total int64
	offset := (page - 1) * pageSize

	q := db.Model(&Post{}).Where("category=?", category).Where("is_published=?", true)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

func GetUserEngagementStats(db *gorm.DB, userID uint) (map[string]interface{}, error) {
    stats := make(map[string]interface{})

    var user User
    if err := db.First(&user, userID).Error; err != nil {
        return nil, err
    }

    var totalPosts int64
    var totalLikesReceived int64
    var totalLikesGiven int64
    var avgPostViews float64
    var mostLikedPost Post
    var topCategory string

    db.Model(&Post{}).Where("user_id=? AND is_published=?", userID, true).Count(&totalPosts)

    db.Model(&Like{}).
        Joins("JOIN posts ON likes.post_id = posts.id").
        Where("posts.user_id=?", userID).
        Count(&totalLikesReceived)

    db.Model(&Like{}).Where("user_id=?", userID).Count(&totalLikesGiven)

    db.Model(&Post{}).
        Where("user_id=?", userID).
        Select("COALESCE(AVG(view_count), 0)").
        Scan(&avgPostViews)

    type Result struct {
        PostID uint
        Count  int
    }
    var result Result
    db.Table("likes").
        Joins("JOIN posts ON likes.post_id=posts.id").
        Where("posts.user_id=?", userID).
        Select("posts.id as post_id, COUNT(likes.id) as count").
        Group("posts.id").
        Order("count DESC").
        Limit(1).
        Scan(&result)

    if result.PostID != 0 {
        db.First(&mostLikedPost, result.PostID)
    }

    db.Model(&Post{}).
        Where("user_id = ?", userID).
        Select("category").
        Group("category").
        Order("COUNT(*) DESC").
        Limit(1).
        Scan(&topCategory)

    stats["user"] = user
    stats["total_posts"] = totalPosts
    stats["total_likes_received"] = totalLikesReceived
    stats["total_likes_given"] = totalLikesGiven
    stats["average_views_per_post"] = fmt.Sprintf("%.2f", avgPostViews)
    stats["top_category"] = topCategory
    if mostLikedPost.ID != 0 {
        stats["most_liked_post"] = mostLikedPost
    } else {
        stats["most_liked_post"] = nil
    }

    return stats, nil
}

func GetPopularPostsByLikes(db *gorm.DB, days int, limit int) ([]Post, error) {
	var posts []Post
	since := time.Now().AddDate(0, 0, -days)

	err := db.Model(&Post{}).
		Select("posts.*, COUNT(likes.id) as like_count").
		Joins("LEFT JOIN likes ON posts.id = likes.post_id").
		Where("posts.created_at >= ?", since).
		Group("posts.id").
		Order("like_count DESC").
		Limit(limit).
		Scan(&posts).Error
	return posts, err
}

func GetCountryUserStats(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := db.Model(&User{}).
		Select("country, COUNT(*) as user_count, AVG(age) as avg_age").
		Group("country").
		Order("user_count DESC").
		Find(&results).Error
	return results, err
}

func SearchPostsByContent(db *gorm.DB, query string, limit int) ([]Post, error) {
	var posts []Post
	pattern := fmt.Sprintf("%%%s%%", query)
	err := db.Where("content LIKE ?", pattern).Limit(limit).Find(&posts).Error
	return posts, err
}

func GetUserRecommendations(db *gorm.DB, userID uint, limit int) ([]User, error) {
	categories := db.Model(&Like{}).
		Select("DISTINCT posts.category").
		Joins("JOIN posts ON likes.post_id = posts.id").
		Where("likes.user_id=?", userID)

	var users []User
	result := db.Model(&User{}).
		Distinct().
		Joins("JOIN likes ON users.id = likes.user_id").
		Joins("JOIN posts ON likes.post_id = posts.id").
		Where("users.id <> ?", userID).
		Where("posts.category IN (?)", categories).
		Limit(limit).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, u := range(users) {
		fmt.Printf("user: %d\n", u.ID)
	}
	return users, nil
}
