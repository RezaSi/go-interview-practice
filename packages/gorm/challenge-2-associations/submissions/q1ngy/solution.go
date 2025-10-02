package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	// Implement database connection with auto-migration
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Post{}, &Tag{})
	db.Debug()
	return db, nil
}

// CreateUserWithPosts creates a new user with associated posts
func CreateUserWithPosts(db *gorm.DB, user *User) error {
	// Implement user creation with posts
	db.Create(user)
	return nil
}

// GetUserWithPosts retrieves a user with all their posts preloaded
func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	// Implement user retrieval with posts
	user := &User{}
	result := db.Preload("Posts").First(user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// CreatePostWithTags creates a new post with specified tags
func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	// Implement post creation with tags
	tags := make([]Tag, 0, len(tagNames))
	for _, t := range tagNames {
		tags = append(tags, Tag{Name: t})
	}
	post.Tags = tags
	db.Create(post)
	return nil
}

// GetPostsByTag retrieves all posts that have a specific tag
func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	// Implement posts retrieval by tag
	var tag Tag
	if err := db.Where("name = ?", tagName).First(&tag); err.Error != nil {
		return nil, err.Error
	}
	var posts []Post
	if err := db.Model(&tag).Association("Posts").Find(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// AddTagsToPost adds tags to an existing post
func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
	// Implement adding tags to existing post
	var post Post
	result := db.Preload("Tags").Find(&post, postID)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	var tags []Tag
	for _, name := range tagNames {
		tags = append(tags, Tag{Name: name})
	}
	db.Model(&post).Association("Tags").Append(&tags)
	return nil
}

// GetPostWithUserAndTags retrieves a post with user and tags preloaded
func GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error) {
	// Implement post retrieval with user and tags
	var post Post
	result := db.Preload("User").Preload("Tags").Find(&post, postID)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &post, nil
}
