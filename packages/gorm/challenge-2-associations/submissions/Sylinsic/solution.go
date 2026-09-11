package main

import (
	"errors"
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
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateUserWithPosts creates a new user with associated posts
func CreateUserWithPosts(db *gorm.DB, user *User) error {
	res := db.Create(user)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("Error adding user to database")
	}

	return nil
}

// GetUserWithPosts retrieves a user with all their posts preloaded
func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	var user User
	if err := db.Preload("Posts").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// CreatePostWithTags creates a new post with specified tags
func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	tx := db.Begin(nil)
	defer tx.Rollback()

	post.Tags = []Tag{}
	for _,tagName := range tagNames {
		post.Tags = append(post.Tags, Tag{
			Name: tagName,
		})
	}
	
	if err := tx.Create(post).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// GetPostsByTag retrieves all posts that have a specific tag
func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	var tag Tag

	if err := db.Preload("Posts").Find(&tag, "Name = ?", tagName).Error; err != nil {
		return nil, err
	}

	return tag.Posts, nil
}

// AddTagsToPost adds tags to an existing post
func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
	post := Post{}
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}
	
	tags := []Tag{}
	for _,tagName := range tagNames {
		tags = append(tags, Tag{Name: tagName})
	}

	if err := db.Model(&post).Association("Tags").Append(tags); err != nil {
		return err
	}

	return nil
}

// GetPostWithUserAndTags retrieves a post with user and tags preloaded
func GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error) {
	post := Post{}
	if err := db.Preload("Tags").Preload("User").First(&post, postID).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
