# Challenge 4: Advanced Queries

Build a **Social Media Analytics System** using GORM that demonstrates advanced querying techniques, aggregations, and complex data analysis.

## Challenge Requirements

Create a Go application that implements:

1. **Complex Queries** - Advanced filtering, sorting, and pagination
2. **Aggregations** - Group by, count, sum, average operations
3. **Subqueries** - Nested queries and correlated subqueries
4. **Raw SQL** - Custom SQL queries when needed
5. **Query Optimization** - Efficient data retrieval patterns

## Data Models

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"unique;not null"`
    Email     string    `gorm:"unique;not null"`
    Age       int       `gorm:"not null"`
    Country   string    `gorm:"not null"`
    CreatedAt time.Time
    Posts     []Post    `gorm:"foreignKey:UserID"`
    Likes     []Like    `gorm:"foreignKey:UserID"`
}

type Post struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"not null"`
    Content     string    `gorm:"type:text"`
    UserID      uint      `gorm:"not null"`
    User        User      `gorm:"foreignKey:UserID"`
    Category    string    `gorm:"not null"`
    ViewCount   int       `gorm:"default:0"`
    IsPublished bool      `gorm:"default:true"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Likes       []Like    `gorm:"foreignKey:PostID"`
}

type Like struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    PostID    uint      `gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID"`
    Post      Post      `gorm:"foreignKey:PostID"`
    CreatedAt time.Time
}
```

## Required Functions

Implement these functions:
- `ConnectDB() (*gorm.DB, error)` - Database connection with auto-migration
- `GetTopUsersByPostCount(db *gorm.DB, limit int) ([]User, error)` - Get users with most posts
- `GetPostsByCategoryWithUserInfo(db *gorm.DB, category string, page, pageSize int) ([]Post, int64, error)` - Get posts with pagination
- `GetUserEngagementStats(db *gorm.DB, userID uint) (map[string]interface{}, error)` - Get user engagement statistics
- `GetPopularPostsByLikes(db *gorm.DB, days int, limit int) ([]Post, error)` - Get popular posts by likes in time period
- `GetCountryUserStats(db *gorm.DB) ([]map[string]interface{}, error)` - Get user statistics by country
- `SearchPostsByContent(db *gorm.DB, query string, limit int) ([]Post, error)` - Search posts by content
- `GetUserRecommendations(db *gorm.DB, userID uint, limit int) ([]User, error)` - Get user recommendations based on similar interests

## Query Examples

**Top Users by Post Count:**
```sql
SELECT users.*, COUNT(posts.id) as post_count 
FROM users 
LEFT JOIN posts ON users.id = posts.user_id 
GROUP BY users.id 
ORDER BY post_count DESC 
LIMIT 10
```

**Popular Posts by Likes:**
```sql
SELECT posts.*, COUNT(likes.id) as like_count 
FROM posts 
LEFT JOIN likes ON posts.id = likes.post_id 
WHERE posts.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY posts.id 
ORDER BY like_count DESC 
LIMIT 20
```

## Testing Requirements

Your solution must pass tests for:
- Retrieving top users by post count with proper aggregation
- Paginated post retrieval with user information
- User engagement statistics calculation
- Popular posts filtering by time period and likes
- Country-based user statistics
- Full-text search functionality
- User recommendation algorithm
- Query performance and optimization 