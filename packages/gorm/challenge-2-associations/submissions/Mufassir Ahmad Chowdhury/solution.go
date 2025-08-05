package main

import (
	"time"

    "gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
	    return nil, err
	}
	err = db.AutoMigrate(&User{}, &Post{}, &Tag{})
	return db, err
}

// CreateUserWithPosts creates a new user with associated posts
func CreateUserWithPosts(db *gorm.DB, user *User) error {
	return db.Create(&user).Error
}

// GetUserWithPosts retrieves a user with all their posts preloaded
func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	var user User
	result := db.Preload("Posts").First(&user, userID)
	if result.Error != nil {
	    return nil, result.Error
	}
	return &user, nil
}

// CreatePostWithTags creates a new post with specified tags
func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	tx := db.Begin()
	defer func() {
	    if r := recover(); r != nil {
	        tx.Rollback()
	    }
	} ()
	
	if err := tx.Create(&post).Error; err != nil {
	    tx.Rollback()
	    return err
	}
	for _, name := range tagNames {
	    var tag Tag
	    if err := tx.FirstOrCreate(&tag, Tag{Name: name}).Error; err != nil {
	        tx.Rollback()
	        return err
	    }
	    if err := tx.Model(&post).Association("Tags").Append(&tag); err != nil {
	        tx.Rollback()
	        return err
	    }
	}
	if err := tx.Commit().Error; err != nil {
	    tx.Rollback()
	    return err
	}
	return nil
}

// GetPostsByTag retrieves all posts that have a specific tag
func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	// TODO: Implement posts retrieval by tag
	var posts []Post
	results := db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
        Joins("JOIN tags ON post_tags.tag_id = tags.id").
        Where("tags.name = ?", tagName).
        Find(&posts)
	return posts, results.Error
}

// AddTagsToPost adds tags to an existing post
func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
	    return err
	}
	tx := db.Begin()
	defer func() {
	    if r := recover(); r != nil {
	        tx.Rollback()
	    }
	} ()
	for _, name := range tagNames {
	    var tag Tag
	    
	    if err := tx.FirstOrCreate(&tag, Tag{Name: name}).Error; err != nil {
	        tx.Rollback()
	        return err
	    }
	    if err := tx.Model(&post).Association("Tags").Append(&tag); err != nil {
	        tx.Rollback()
	        return err
	    }
	}
	if err := tx.Commit().Error; err != nil {
	    tx.Rollback()
	    return err
	}
	return nil
}

// GetPostWithUserAndTags retrieves a post with user and tags preloaded
func GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error) {
	var post Post
	if err := db.Preload("User").Preload("Tags").First(&post, postID).Error; err != nil {
	    return nil, err
	}
	return &post, nil
}
