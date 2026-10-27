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
	db, err := gorm.Open(sqlite.Open("test.db?_foreign_keys=on"), &gorm.Config{})
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
	err := db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUserWithPosts retrieves a user with all their posts preloaded
func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	var user User
	res := db.Preload("Posts").First(&user, userID)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

// CreatePostWithTags creates a new post with specified tags
func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var tags []Tag
		var err error

		for _, t := range tagNames {
			var tag Tag
			err = tx.Where("name = ?", t).FirstOrCreate(&tag, Tag{Name: t}).Error
			if err != nil {
				return err
			}
			tags = append(tags, tag)
		}

		post.Tags = tags
		err = tx.Create(post).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// GetPostsByTag retrieves all posts that have a specific tag
func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	var posts []Post
	err := db.
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.name = ?", tagName).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// AddTagsToPost adds tags to an existing post
func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
	var post Post
	var err error

	err = db.Preload("Tags").First(&post, postID).Error
	if err != nil {
		return err
	}

	var tags []Tag
	for _, t := range tagNames {
		var tag Tag
		err = db.Where("name = ?", t).FirstOrCreate(&tag, Tag{Name: t}).Error
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	err = db.Model(&post).Association("Tags").Append(&tags)
	if err != nil {
		return err
	}

	return nil
}

// GetPostWithUserAndTags retrieves a post with user and tags preloaded
func GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error) {
	var post Post
	err := db.Preload("User").Preload("Tags").First(&post, postID).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}
