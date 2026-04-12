# Learning GORM Associations

## Overview

GORM associations allow you to define relationships between different models, making it easy to work with related data. This challenge focuses on understanding and implementing various types of associations in GORM.

## Types of Associations

### 1. One-to-One
A one-to-one relationship where each record in one table is associated with exactly one record in another table.

```go
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string
    Profile  Profile `gorm:"foreignKey:UserID"`
}

type Profile struct {
    ID     uint   `gorm:"primaryKey"`
    UserID uint   `gorm:"unique"`
    Bio    string
}
```

### 2. One-to-Many
A one-to-many relationship where one record can be associated with multiple records in another table.

```go
type User struct {
    ID    uint    `gorm:"primaryKey"`
    Name  string
    Posts []Post  `gorm:"foreignKey:UserID"`
}

type Post struct {
    ID     uint   `gorm:"primaryKey"`
    Title  string
    UserID uint
    User   User   `gorm:"foreignKey:UserID"`
}
```

### 3. Many-to-Many
A many-to-many relationship where multiple records can be associated with multiple records in another table.

```go
type Post struct {
    ID    uint   `gorm:"primaryKey"`
    Title string
    Tags  []Tag  `gorm:"many2many:post_tags;"`
}

type Tag struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Posts []Post `gorm:"many2many:post_tags;"`
}
```

## Key Concepts

### Foreign Keys
- Use `gorm:"foreignKey:FieldName"` to specify the foreign key field
- The foreign key should reference the primary key of the related model
- GORM automatically handles foreign key constraints

### Preloading
Preloading allows you to load related data efficiently:

```go
// Load user with posts
var user User
db.Preload("Posts").First(&user, userID)

// Load post with user and tags
var post Post
db.Preload("User").Preload("Tags").First(&post, postID)
```

### Association Mode
GORM provides different association modes for creating related records:

```go
// Create user with posts
user := User{
    Name: "John",
    Posts: []Post{
        {Title: "First Post"},
        {Title: "Second Post"},
    },
}
db.Create(&user) // Creates user and posts in a transaction
```

## Best Practices

### 1. Use Preloading for Performance
Always use preloading when you need related data to avoid N+1 query problems:

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

### 2. Handle Association Errors
Always check for errors when working with associations:

```go
if err := db.Create(&user).Error; err != nil {
    // Handle error
}
```

### 3. Use Transactions for Complex Operations
Use transactions when creating multiple related records:

```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Commit().Error; err != nil {
    return err
}
```

## Common Patterns

### Creating Related Records
```go
// Method 1: Using association mode
user := User{Name: "John"}
post := Post{Title: "My Post", User: user}
db.Create(&post)

// Method 2: Using Association
user := User{Name: "John"}
db.Create(&user)
post := Post{Title: "My Post"}
db.Model(&user).Association("Posts").Append(&post)
```

### Querying Related Data
```go
// Get user with posts
var user User
db.Preload("Posts").First(&user, userID)

// Get posts by user
var posts []Post
db.Where("user_id = ?", userID).Find(&posts)

// Get users who have posts
var users []User
db.Preload("Posts").Where("EXISTS (SELECT 1 FROM posts WHERE posts.user_id = users.id)").Find(&users)
```

### Updating Associations
```go
// Replace all posts for a user
db.Model(&user).Association("Posts").Replace([]Post{
    {Title: "New Post 1"},
    {Title: "New Post 2"},
})

// Add posts to user
db.Model(&user).Association("Posts").Append(&Post{Title: "New Post"})

// Remove posts from user
db.Model(&user).Association("Posts").Delete(&Post{Title: "Old Post"})
```

## Advanced Features

### Polymorphic Associations
GORM supports polymorphic associations for more complex relationships:

```go
type Comment struct {
    ID        uint   `gorm:"primaryKey"`
    Content   string
    UserID    uint
    User      User
    CommentableType string
    CommentableID   uint
}

type Post struct {
    ID       uint      `gorm:"primaryKey"`
    Title    string
    Comments []Comment `gorm:"polymorphic:Commentable;"`
}
```

### Self-Referential Associations
For hierarchical data structures:

```go
type Category struct {
    ID       uint       `gorm:"primaryKey"`
    Name     string
    ParentID *uint
    Parent   *Category  `gorm:"foreignKey:ParentID"`
    Children []Category `gorm:"foreignKey:ParentID"`
}
```

## Resources

- [GORM Associations Documentation](https://gorm.io/docs/associations.html)
- [GORM Preloading](https://gorm.io/docs/preload.html)
- [GORM Many to Many](https://gorm.io/docs/many_to_many.html)
- [GORM Polymorphic Associations](https://gorm.io/docs/polymorphic_association.html)

## Practice Exercises

1. Create a blog system with users, posts, and comments
2. Implement a product catalog with categories and tags
3. Build a social network with users, posts, and likes
4. Create an e-commerce system with orders, products, and customers

These exercises will help you master GORM associations and understand how to design effective database relationships. 