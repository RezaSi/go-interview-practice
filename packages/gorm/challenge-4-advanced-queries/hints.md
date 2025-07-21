# Hints for GORM Advanced Queries Challenge

## Hint 1: Database Connection & Data Model

This challenge involves Users, Posts, and Likes with complex relationships. Use `gorm.Open()` with SQLite driver and auto-migrate all models:

```go
func ConnectDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    err = db.AutoMigrate(&User{}, &Post{}, &Like{})
    return db, err
}
```

## Hint 2: Top Users by Post Count

Use aggregations with joins, grouping, and ordering:

```go
func GetTopUsersByPostCount(db *gorm.DB, limit int) ([]User, error) {
    var users []User
    err := db.Joins("LEFT JOIN posts ON users.id = posts.user_id").
        Group("users.id").
        Order("COUNT(posts.id) DESC").
        Limit(limit).
        Find(&users).Error
    return users, err
}
```

## Hint 3: Posts by Category with Pagination

Use `Where()` to filter, `Preload()` for user info, and implement pagination:

```go
func GetPostsByCategoryWithUserInfo(db *gorm.DB, category string, page, pageSize int) ([]Post, int64, error) {
    var posts []Post
    var total int64
    
    query := db.Where("category = ?", category)
    query.Model(&Post{}).Count(&total)
    
    offset := (page - 1) * pageSize
    err := query.Preload("User").Offset(offset).Limit(pageSize).Find(&posts).Error
    
    return posts, total, err
}
```

## Hint 4: User Engagement Statistics

Calculate multiple metrics in one function using different query methods:

```go
func GetUserEngagementStats(db *gorm.DB, userID uint) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Post count
    var postCount int64
    db.Model(&Post{}).Where("user_id = ?", userID).Count(&postCount)
    stats["post_count"] = postCount
    
    // Likes received
    var likesReceived int64
    db.Model(&Like{}).Joins("JOIN posts ON likes.post_id = posts.id").
        Where("posts.user_id = ?", userID).Count(&likesReceived)
    stats["likes_received"] = likesReceived
    
    // Average views
    var avgViews float64
    db.Model(&Post{}).Select("AVG(view_count)").Where("user_id = ?", userID).Scan(&avgViews)
    stats["average_views"] = avgViews
    
    return stats, nil
}
```

## Hint 5: Popular Posts by Likes in Time Period

Use joins with time filtering and aggregation:

```go
func GetPopularPostsByLikes(db *gorm.DB, days int, limit int) ([]Post, error) {
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
```

## Hint 6: Country User Statistics

Use `Select()` with aggregation functions and `Group()`:

```go
func GetCountryUserStats(db *gorm.DB) ([]map[string]interface{}, error) {
    var results []struct {
        Country   string
        UserCount int64
        AvgAge    float64
    }
    
    err := db.Model(&User{}).
        Select("country, COUNT(*) as user_count, AVG(age) as avg_age").
        Group("country").
        Scan(&results).Error
    
    var stats []map[string]interface{}
    for _, result := range results {
        stat := map[string]interface{}{
            "country":    result.Country,
            "user_count": result.UserCount,
            "avg_age":    result.AvgAge,
        }
        stats = append(stats, stat)
    }
    
    return stats, err
}
```

## Hint 7: Search Posts by Content

Use `Where()` with `LIKE` operator for multiple fields:

```go
func SearchPostsByContent(db *gorm.DB, query string, limit int) ([]Post, error) {
    var posts []Post
    searchPattern := "%" + query + "%"
    
    err := db.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
        Limit(limit).
        Find(&posts).Error
    
    return posts, err
}
```

## Hint 8: User Recommendations

Use subqueries to find users with similar interests:

```go
func GetUserRecommendations(db *gorm.DB, userID uint, limit int) ([]User, error) {
    var users []User
    
    // Find users who liked posts in similar categories as the current user
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
```

## Query Patterns

### Aggregation Queries
```go
// Count with grouping
var results []struct {
    UserID    uint
    PostCount int64
}
db.Model(&Post{}).
   Select("user_id, COUNT(*) as post_count").
   Group("user_id").
   Order("post_count DESC").
   Scan(&results)
```

### Complex Joins
```go
// Join multiple tables
var posts []Post
db.Joins("User").
   Joins("LEFT JOIN likes ON posts.id = likes.post_id").
   Where("posts.category = ?", category).
   Group("posts.id").
   Having("COUNT(likes.id) > ?", minLikes).
   Find(&posts)
```

### Subqueries
```go
// Use subquery for filtering
var users []User
db.Where("id IN (?)", 
    db.Model(&Post{}).
       Select("user_id").
       Group("user_id").
       Having("COUNT(*) > ?", 5)).
   Find(&users)
```

### Pagination
```go
func GetPaginatedResults(db *gorm.DB, page, pageSize int) ([]Post, int64, error) {
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

## Performance Optimization

### Use Indexes
```go
// Add indexes to your models
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex"`
    Country  string `gorm:"index"`
}

type Post struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint   `gorm:"index"`
    Category  string `gorm:"index"`
    CreatedAt time.Time `gorm:"index"`
}
```

### Avoid N+1 Queries
```go
// Good: Use preloading
var users []User
db.Preload("Posts").Find(&users)

// Bad: N+1 queries
var users []User
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)
}
```

### Limit Result Sets
```go
// Always limit large result sets
db.Limit(100).Find(&users)

// Use cursor-based pagination for large datasets
db.Where("id > ?", cursor).Limit(50).Find(&posts)
```

## Error Handling

1. **Check for errors** after each database operation
2. **Handle empty results** - return empty slices, not nil
3. **Validate input parameters** - check for valid page numbers, limits, etc.
4. **Handle database errors** - connection issues, constraint violations, etc.

## Testing Strategies

### Setup Test Data
```go
func setupTestData(db *gorm.DB) {
    // Create users
    users := []User{
        {Username: "user1", Email: "user1@test.com", Age: 25, Country: "USA"},
        {Username: "user2", Email: "user2@test.com", Age: 30, Country: "Canada"},
    }
    for i := range users {
        db.Create(&users[i])
    }
    
    // Create posts
    posts := []Post{
        {Title: "Post 1", Content: "Content 1", UserID: users[0].ID, Category: "tech"},
        {Title: "Post 2", Content: "Content 2", UserID: users[0].ID, Category: "sports"},
    }
    for i := range posts {
        db.Create(&posts[i])
    }
    
    // Create likes
    likes := []Like{
        {UserID: users[1].ID, PostID: posts[0].ID},
        {UserID: users[0].ID, PostID: posts[1].ID},
    }
    for i := range likes {
        db.Create(&likes[i])
    }
}
```

### Test Aggregations
```go
// Test top users by post count
users, err := GetTopUsersByPostCount(db, 3)
assert.NoError(t, err)
assert.Len(t, users, 2) // Only 2 users have posts
assert.Equal(t, "user1", users[0].Username) // user1 has 2 posts
```

## Common Patterns

### Time-Based Filtering
```go
// Filter by time period
db.Where("created_at >= ?", time.Now().AddDate(0, 0, -days))
```

### Full-Text Search
```go
// Simple LIKE search
db.Where("title LIKE ? OR content LIKE ?", 
    "%"+query+"%", "%"+query+"%")
```

### Conditional Queries
```go
// Build queries conditionally
query := db.Model(&Post{})
if category != "" {
    query = query.Where("category = ?", category)
}
if userID != 0 {
    query = query.Where("user_id = ?", userID)
}
query.Find(&posts)
```

## Debugging Tips

1. **Enable GORM logging**:
```go
db = db.Debug()
```

2. **Check query results**:
```go
// Print query results
var count int64
db.Model(&User{}).Count(&count)
fmt.Printf("Total users: %d\n", count)
```

3. **Verify relationships**:
```go
// Check if associations are loaded
user := User{}
db.Preload("Posts").First(&user, userID)
fmt.Printf("User has %d posts\n", len(user.Posts))
```

## Common Mistakes to Avoid

1. **Not using preloading** - This leads to N+1 query problems
2. **Forgetting to limit results** - Can cause performance issues
3. **Not handling empty results** - Return empty slices, not nil
4. **Not using transactions** - For complex operations involving multiple tables
5. **Not optimizing queries** - Use proper indexes and query patterns

## Useful GORM Methods

- `db.Joins()` - Join tables
- `db.Preload()` - Preload related data
- `db.Group()` - Group results
- `db.Having()` - Filter grouped results
- `db.Select()` - Select specific fields or aggregations
- `db.Count()` - Count records
- `db.Offset()` / `db.Limit()` - Pagination
- `db.Order()` - Sort results
- `db.Scan()` - Scan results to struct

## SQLite Specific Notes

- SQLite doesn't support full-text search like MySQL
- Use `LIKE` for text search
- Some complex aggregations might be slower
- Be careful with large datasets

## Final Tips

1. **Start with simple queries** - Get basic functionality working first
2. **Test with small datasets** - Verify logic before scaling up
3. **Use the learning resources** - Check GORM documentation for examples
4. **Profile your queries** - Use GORM debug mode to see actual SQL
5. **Consider caching** - For frequently accessed data

## Performance Checklist

- [ ] Use proper indexes
- [ ] Limit result sets
- [ ] Use preloading to avoid N+1 queries
- [ ] Use transactions for complex operations
- [ ] Optimize aggregation queries
- [ ] Handle edge cases and errors
- [ ] Test with realistic data volumes 