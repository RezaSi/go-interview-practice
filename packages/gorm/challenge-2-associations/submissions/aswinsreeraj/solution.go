package main

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

// User represents a user in the blog system
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Posts     []Post `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Post represents a blog post
type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	UserID    uint   `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	Tags      []Tag  `gorm:"many2many:post_tags;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Tag represents a tag for categorizing posts
type Tag struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Posts []Post `gorm:"many2many:post_tags;"`
}

// ConnectDB establishes a connection to the SQLite database and auto-migrates the models
func ConnectDB() (*gorm.DB, error) {
	// TODO: Implement database connection with auto-migration
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
	    return nil, err
	}
	err = db.AutoMigrate(&User{}, &Post{}, &Tag{})
	return db, err
}

// CreateUserWithPosts creates a new user with associated posts
func CreateUserWithPosts(db *gorm.DB, user *User) error {
	// TODO: Implement user creation with posts
	result := db.Create(user)
	return result.Error
}

// GetUserWithPosts retrieves a user with all their posts preloaded
func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	// TODO: Implement user retrieval with posts
	var user User
	result := db.Preload("Posts").First(&user, userID)
	return &user, result.Error
}

// CreatePostWithTags creates a new post with specified tags
func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	// TODO: Implement post creation with tags
	if err := db.Create(post).Error; err != nil {
	    return err
	}
	for _, name := range tagNames {
	    var tag Tag
	    db.FirstOrCreate(&tag, Tag{Name: name})
	    db.Model(post).Association("Tags").Append(&tag)
	}
	return nil
}

// GetPostsByTag retrieves all posts that have a specific tag
func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	// TODO: Implement posts retrieval by tag
	var posts []Post
	err := db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
	            Joins("JOIN tags ON post_tags.tag_id = tags.id").
	            Where("tags.name = ?", tagName).
	            Find(&posts).Error
	return posts, err
}

// AddTagsToPost adds tags to an existing post
func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
	// TODO: Implement adding tags to existing post
	var post Post
	result := db.First(&post, postID)
	if result.Error != nil {
	    return result.Error
	}
	for _, tagName := range tagNames {
	    var tag Tag
	    db.FirstOrCreate(&tag, Tag{Name: tagName})
	    if err := db.Model(&post).Association("Tags").Append(&tag); err != nil {
	     return err
	 }   
	}
	return nil
}

// GetPostWithUserAndTags retrieves a post with user and tags preloaded
func GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error) {
	// TODO: Implement post retrieval with user and tags
	var post Post
	err := db.Preload("User").Preload("Tags").First(&post, postID).Error
	if err != nil {
	    return nil, err
	}
	return &post, nil
}
