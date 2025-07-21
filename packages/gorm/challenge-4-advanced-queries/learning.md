# Learning GORM Advanced Queries

## Overview

Advanced queries in GORM allow you to perform complex database operations, aggregations, and analytics. This challenge focuses on mastering advanced querying techniques to build efficient and powerful data retrieval systems.

## Query Building

### 1. Chain Methods
GORM allows you to chain query methods for complex queries:

```go
var users []User
db.Where("age > ?", 18).
   Where("country = ?", "USA").
   Order("created_at DESC").
   Limit(10).
   Find(&users)
```

### 2. Preloading with Conditions
Preload related data with specific conditions:

```go
var users []User
db.Preload("Posts", "is_published = ?", true).
   Preload("Posts.Tags").
   Find(&users)
```

## Aggregations

### 1. Count
Count records with conditions:

```go
var count int64
db.Model(&User{}).Where("country = ?", "USA").Count(&count)

// Count with distinct
db.Model(&User{}).Distinct("country").Count(&count)
```

### 2. Sum, Average, Min, Max
Perform mathematical aggregations:

```go
var total float64
db.Model(&Product{}).Select("SUM(price)").Scan(&total)

var avgPrice float64
db.Model(&Product{}).Select("AVG(price)").Scan(&avgPrice)

var maxPrice float64
db.Model(&Product{}).Select("MAX(price)").Scan(&maxPrice)
```

### 3. Group By
Group results by specific fields:

```go
type CountryStats struct {
    Country string
    Count   int64
    AvgAge  float64
}

var stats []CountryStats
db.Model(&User{}).
   Select("country, COUNT(*) as count, AVG(age) as avg_age").
   Group("country").
   Scan(&stats)
```

## Complex Queries

### 1. Subqueries
Use subqueries for complex filtering:

```go
// Find users who have more than 5 posts
var users []User
db.Where("id IN (?)", 
    db.Model(&Post{}).
       Select("user_id").
       Group("user_id").
       Having("COUNT(*) > ?", 5)).
   Find(&users)
```

### 2. Joins
Perform complex joins:

```go
var results []struct {
    UserName  string
    PostCount int64
    LikeCount int64
}

db.Table("users").
   Select("users.username, COUNT(DISTINCT posts.id) as post_count, COUNT(likes.id) as like_count").
   Joins("LEFT JOIN posts ON users.id = posts.user_id").
   Joins("LEFT JOIN likes ON posts.id = likes.post_id").
   Group("users.id, users.username").
   Scan(&results)
```

### 3. Raw SQL
Use raw SQL for complex queries:

```go
var users []User
db.Raw(`
    SELECT u.*, COUNT(p.id) as post_count 
    FROM users u 
    LEFT JOIN posts p ON u.id = p.user_id 
    WHERE u.country = ? 
    GROUP BY u.id 
    HAVING post_count > ? 
    ORDER BY post_count DESC
`, "USA", 5).Scan(&users)
```

## Pagination

### 1. Offset and Limit
Implement pagination:

```go
func GetPaginatedPosts(db *gorm.DB, page, pageSize int) ([]Post, int64, error) {
    var posts []Post
    var total int64
    
    // Get total count
    db.Model(&Post{}).Count(&total)
    
    // Get paginated results
    offset := (page - 1) * pageSize
    err := db.Offset(offset).Limit(pageSize).Find(&posts).Error
    
    return posts, total, err
}
```

### 2. Cursor-Based Pagination
For better performance with large datasets:

```go
func GetPostsByCursor(db *gorm.DB, cursor uint, limit int) ([]Post, error) {
    var posts []Post
    err := db.Where("id > ?", cursor).
              Order("id ASC").
              Limit(limit).
              Find(&posts).Error
    return posts, err
}
```

## Full-Text Search

### 1. LIKE Queries
Simple text search:

```go
var posts []Post
db.Where("title LIKE ? OR content LIKE ?", 
    "%"+query+"%", "%"+query+"%").Find(&posts)
```

### 2. Full-Text Search (MySQL)
For more advanced search capabilities:

```go
var posts []Post
db.Where("MATCH(title, content) AGAINST(? IN BOOLEAN MODE)", query).Find(&posts)
```

## Query Optimization

### 1. Select Specific Fields
Only select needed fields:

```go
var users []User
db.Select("id, username, email").Find(&users)
```

### 2. Use Indexes
Ensure proper indexing for frequently queried fields:

```go
// Add indexes to your models
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex"`
    Email    string `gorm:"uniqueIndex"`
    Country  string `gorm:"index"`
}
```

### 3. Avoid N+1 Queries
Use preloading to avoid multiple queries:

```go
// Good: Single query with preloading
var users []User
db.Preload("Posts").Find(&users)

// Bad: N+1 queries
var users []User
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)
}
```

## Analytics Queries

### 1. Time-Based Analytics
Analyze data over time periods:

```go
type DailyStats struct {
    Date  string
    Count int64
}

var stats []DailyStats
db.Model(&Post{}).
   Select("DATE(created_at) as date, COUNT(*) as count").
   Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
   Group("DATE(created_at)").
   Order("date ASC").
   Scan(&stats)
```

### 2. User Engagement Metrics
Calculate user engagement:

```go
func GetUserEngagement(db *gorm.DB, userID uint) map[string]interface{} {
    var stats map[string]interface{}
    
    // Get post count
    var postCount int64
    db.Model(&Post{}).Where("user_id = ?", userID).Count(&postCount)
    
    // Get total likes received
    var likesReceived int64
    db.Model(&Like{}).
       Joins("JOIN posts ON likes.post_id = posts.id").
       Where("posts.user_id = ?", userID).
       Count(&likesReceived)
    
    // Get average views per post
    var avgViews float64
    db.Model(&Post{}).
       Select("AVG(view_count)").
       Where("user_id = ?", userID).
       Scan(&avgViews)
    
    stats = map[string]interface{}{
        "total_posts":           postCount,
        "total_likes_received":  likesReceived,
        "average_views_per_post": avgViews,
    }
    
    return stats
}
```

## Advanced Patterns

### 1. Query Scopes
Create reusable query conditions:

```go
func (db *gorm.DB) ActiveUsers() *gorm.DB {
    return db.Where("is_active = ?", true)
}

func (db *gorm.DB) PublishedPosts() *gorm.DB {
    return db.Where("is_published = ?", true)
}

// Usage
var users []User
db.ActiveUsers().Find(&users)

var posts []Post
db.PublishedPosts().Preload("User").Find(&posts)
```

### 2. Query Hooks
Add custom logic to queries:

```go
func (u *User) AfterFind(tx *gorm.DB) error {
    // Calculate additional fields after query
    u.FullName = u.FirstName + " " + u.LastName
    return nil
}
```

### 3. Custom Scanners
Handle complex query results:

```go
type UserStats struct {
    UserID    uint
    Username  string
    PostCount int64
    LikeCount int64
}

func (us *UserStats) Scan(value interface{}) error {
    // Custom scanning logic
    return nil
}
```

## Performance Tips

### 1. Use Indexes
Create indexes for frequently queried fields:

```sql
CREATE INDEX idx_users_country ON users(country);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_likes_post_id ON likes(post_id);
```

### 2. Limit Result Sets
Always limit large result sets:

```go
db.Limit(100).Find(&users) // Limit to 100 results
```

### 3. Use Count Efficiently
Use count efficiently for pagination:

```go
// Good: Use count for total
var total int64
db.Model(&User{}).Count(&total)

// Bad: Load all records to count
var users []User
db.Find(&users)
total := len(users)
```

### 4. Cache Results
Cache frequently accessed data:

```go
func GetCachedUserStats(userID uint) map[string]interface{} {
    cacheKey := fmt.Sprintf("user_stats_%d", userID)
    
    // Check cache first
    if cached, found := cache.Get(cacheKey); found {
        return cached.(map[string]interface{})
    }
    
    // Calculate stats
    stats := calculateUserStats(userID)
    
    // Cache for 5 minutes
    cache.Set(cacheKey, stats, 5*time.Minute)
    
    return stats
}
```

## Resources

- [GORM Advanced Queries](https://gorm.io/docs/advanced_query.html)
- [GORM Raw SQL](https://gorm.io/docs/raw_sql.html)
- [GORM Scopes](https://gorm.io/docs/scopes.html)
- [GORM Hooks](https://gorm.io/docs/hooks.html)
- [SQL Performance Tuning](https://use-the-index-luke.com/)

## Practice Exercises

1. Build a social media analytics dashboard
2. Implement a search engine with filters
3. Create a recommendation system
4. Build a reporting system with aggregations
5. Implement real-time analytics queries

These exercises will help you master advanced querying techniques and build efficient data retrieval systems. 