# Challenge 2: Associations & Relationships

Build a **Blog System** using GORM that demonstrates database relationships and associations between models.

## Challenge Requirements

Create a Go application that implements:

1. **One-to-Many Relationship** - Users can have multiple posts
2. **Many-to-Many Relationship** - Posts can have multiple tags
3. **Association Operations** - Create, query, and manage related data
4. **Preloading** - Efficiently load related data

## Data Models

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    Posts     []Post    `gorm:"foreignKey:UserID"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Post struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"not null"`
    Content     string    `gorm:"type:text"`
    UserID      uint      `gorm:"not null"`
    User        User      `gorm:"foreignKey:UserID"`
    Tags        []Tag     `gorm:"many2many:post_tags;"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Tag struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"unique;not null"`
    Posts []Post `gorm:"many2many:post_tags;"`
}
```

## Required Functions

Implement these functions:
- `ConnectDB() (*gorm.DB, error)` - Database connection with auto-migration
- `CreateUserWithPosts(db *gorm.DB, user *User) error` - Create user with posts
- `GetUserWithPosts(db *gorm.DB, userID uint) (*User, error)` - Get user with posts
- `CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error` - Create post with tags
- `GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error)` - Get posts by tag
- `AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error` - Add tags to existing post
- `GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error)` - Get post with user and tags

## Testing Requirements

Your solution must pass tests for:
- Creating users with associated posts
- Creating posts with multiple tags
- Querying users with their posts (preloading)
- Querying posts by tag
- Adding tags to existing posts
- Loading posts with user and tag associations
- Proper foreign key constraints and relationships 