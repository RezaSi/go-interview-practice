# Hints for GORM Advanced Queries Challenge

## General Tips

1. **Understand the data model** - This challenge involves Users, Posts, and Likes with complex relationships.

2. **Use aggregations effectively** - Many functions require counting, grouping, and mathematical operations.

3. **Optimize for performance** - Use proper indexing, limit result sets, and avoid N+1 queries.

4. **Handle edge cases** - Consider empty results, invalid inputs, and error conditions.

## Function-Specific Hints

### ConnectDB()
- Use `gorm.Open()` with SQLite driver
- Auto-migrate all models (User, Post, Like)
- Return the database connection

### GetTopUsersByPostCount()
- Use `Joins()` to join users with posts
- Use `Group()` to group by user
- Use `Select()` to include count in results
- Use `Order()` to sort by post count descending
- Use `Limit()` to restrict results

### GetPostsByCategoryWithUserInfo()
- Use `Where()` to filter by category
- Use `Preload("User")` to load user information
- Implement pagination with `Offset()` and `Limit()`
- Use `Count()` to get total count for pagination
- Handle the case where category doesn't exist

### GetUserEngagementStats()
- Calculate multiple metrics in one function
- Use `Count()` for post count
- Use `Joins()` and `Count()` for likes received
- Use `Select("AVG(view_count)")` for average views
- Return a map with all statistics

### GetPopularPostsByLikes()
- Use `Joins()` to join posts with likes
- Use `Where()` to filter by time period
- Use `Group()` to group by post
- Use `Order()` to sort by like count
- Use `Limit()` to restrict results

### GetCountryUserStats()
- Use `Select()` with aggregation functions
- Use `Group("country")` to group by country
- Use `Scan()` to map results to struct
- Return slice of maps with country statistics

### SearchPostsByContent()
- Use `Where()` with `LIKE` operator
- Search in both title and content fields
- Use `OR` conditions for multiple fields
- Use `Limit()` to restrict results
- Consider case-insensitive search

### GetUserRecommendations()
- Find users with similar interests (categories)
- Use subqueries to find users who like similar posts
- Exclude the current user from results
- Use `Limit()` to restrict recommendations
- Consider using `DISTINCT` to avoid duplicates

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