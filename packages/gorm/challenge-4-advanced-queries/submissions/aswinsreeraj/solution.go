package main

import (
	"time"
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

// ConnectDB establishes a connection to the SQLite database with auto-migration
func ConnectDB() (*gorm.DB, error) {
	// TODO: Implement database connection with auto-migration
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
	    return nil, err
	}
	err = db.AutoMigrate(&User{}, &Post{}, &Like{})
	return db, err
}

// GetTopUsersByPostCount retrieves users with the most posts
func GetTopUsersByPostCount(db *gorm.DB, limit int) ([]User, error) {
	// TODO: Implement top users by post count aggregation
	var users []User
	err := db.Joins("LEFT JOIN posts ON users.id = posts.user_id").
	        Group("users.id").
	        Order("COUNT(posts.id) DESC").
	        Limit(limit).
	        Find(&users).Error
	return users, err
}

// GetPostsByCategoryWithUserInfo retrieves posts by category with pagination and user info
func GetPostsByCategoryWithUserInfo(db *gorm.DB, category string, page, pageSize int) ([]Post, int64, error) {
	// TODO: Implement paginated posts retrieval with user info
	if page < 0 {
	    return nil, 0, errors.New("Error")
	}
	var posts []Post
	var total int64
	query := db.Where("category = ?", category)
	query.Model(&Post{}).Count(&total)
	offset := (page - 1) * pageSize
	err := query.Preload("User").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}

// GetUserEngagementStats calculates engagement statistics for a user
func GetUserEngagementStats(db *gorm.DB, userID uint) (map[string]interface{}, error) {
	// TODO: Implement user engagement statistics
	stats := make(map[string]interface{})
	var user User
    if err := db.First(&user, userID).Error; err != nil {
        return nil, err
    }
    var postCount int64
    db.Model(&Post{}).Where("user_id = ?", userID).Count(&postCount)
    stats["total_posts"] = postCount
   
    var likesReceived int64
    db.Model(&Like{}).Joins("JOIN posts ON likes.post_id = posts.id").
        Where("posts.user_id = ?", userID).Count(&likesReceived)
    stats["total_likes_received"] = likesReceived
    
    var likesGiven int64
    db.Model(&Like{}).Joins("JOIN users ON likes.user_id = users.id").
        Where("users.id = ?", userID).Count(&likesGiven)
    stats["total_likes_given"] = likesGiven
    
    var avgViews float64
    db.Model(&Post{}).Select("AVG(view_count)").Where("user_id = ?", userID).Scan(&avgViews)
    stats["average_views_per_post"] = avgViews
    
    return stats, nil
}

// GetPopularPostsByLikes retrieves popular posts by likes in a time period
func GetPopularPostsByLikes(db *gorm.DB, days int, limit int) ([]Post, error) {
	// TODO: Implement popular posts by likes
	var posts []Post
	cutoffDate := time.Now().AddDate(0, 0, -days)
	err := db.Joins("LEFT JOIN likes ON posts.id = likes.post_id").
	        Where("posts.created_at >= ?", cutoffDate).
	        Group("posts.id").
	        Order("COUNT(likes.id) DESC").
	        Limit(limit).
	        Find(&posts).Error
	return posts, err
}

// GetCountryUserStats retrieves user statistics grouped by country
func GetCountryUserStats(db *gorm.DB) ([]map[string]interface{}, error) {
	// TODO: Implement country-based user statistics
	var results []struct {
	    Country string
	    UserCount int64
	    AvgAge float64
	}
	err := db.Model(&User{}).
	        Select("country, COUNT(*) as user_count, AVG(age) as avg_age").
	        Group("country").
	        Scan(&results).Error
	var stats []map[string]interface{}
	for _, result := range results {
	    stat := map[string]interface{}{
	        "country": result.Country,
	        "user_count": result.UserCount,
	        "avg_age": result.AvgAge,
	    }
	    stats = append(stats, stat)
	}
	return stats, err
}

// SearchPostsByContent searches posts by content using full-text search
func SearchPostsByContent(db *gorm.DB, query string, limit int) ([]Post, error) {
	// TODO: Implement full-text search
	var posts []Post
	searchPattern := "%" + query + "%"
	err := db.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
	        Limit(limit).
	        Find(&posts).Error
	return posts, err
}

// GetUserRecommendations retrieves user recommendations based on similar interests
func GetUserRecommendations(db *gorm.DB, userID uint, limit int) ([]User, error) {
	// TODO: Implement user recommendations algorithm
	var users []User
	err := db.Where("id != ? AND id IN (?)", userID,
	        db.Model(&Like{}).
	            Select("DISTINCT likes.user_id").
	            Joins("JOIN posts ON likes.post_id = posts.id").
	            Joins("JOIN posts p2 ON p2.category = posts.category").
	            Joins("JOIN likes l2 ON l2.post_id = p2.id").
	            Where("l2.user_id = ?", userID)).
	       Limit(limit).
	       Find(&users).Error
	return users, err
}
