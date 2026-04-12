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

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &Post{}, &Tag{}); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUserWithPosts(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	var user User
	err := db.Preload("Posts").First(&user, userID).Error
	return &user, err
}

func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		for _, name := range tagNames {
			var tag Tag
			if err := tx.FirstOrCreate(&tag, Tag{Name: name}).Error; err != nil {
				return err
			}
			if err := tx.Model(post).Association("Tags").Append(&tag); err != nil {
				return err
			}
		}
		return nil
	})
}

func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	var posts []Post
	err := db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.name = ?", tagName).
		Preload("User").
		Preload("Tags").
		Find(&posts).Error
	return posts, err
}

func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var post Post
		if err := tx.First(&post, postID).Error; err != nil {
			return err
		}
		for _, name := range tagNames {
			var tag Tag
			if err := tx.FirstOrCreate(&tag, Tag{Name: name}).Error; err != nil {
				return err
			}
			if err := tx.Model(&post).Association("Tags").Append(&tag); err != nil {
				return err
			}
		}
		return nil
	})
}

func GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error) {
	var post Post
	err := db.Preload("User").Preload("Tags").First(&post, postID).Error
	return &post, err
}
