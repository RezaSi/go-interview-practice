package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Posts     []Post `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
	err = db.AutoMigrate(&User{}, &Post{}, &Tag{})
	return db, err
}

func CreateUserWithPosts(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	var user User
	err := db.Preload("Posts").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	// Create the post first
	if err := db.Create(post).Error; err != nil {
		return err
	}
	// Find or create tags and associate them
	for _, name := range tagNames {
		var tag Tag
		db.FirstOrCreate(&tag, Tag{Name: name})
		db.Model(post).Association("Tags").Append(&tag)
	}
	return nil
}

// GetPostsByTag retrieves all posts that have a specific tag
func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	var posts []Post
	err := db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.name = ?", tagName).
		Find(&posts).Error
	return posts, err
}

// AddTagsToPost adds tags to an existing post
func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}
	for _, name := range tagNames {
		var tag Tag
		db.FirstOrCreate(&tag, Tag{Name: name})
		db.Model(&post).Association("Tags").Append(&tag)
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
