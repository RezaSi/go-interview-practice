package main

import (
	"time"
	"database/sql" 
    "errors"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	_ "modernc.org/sqlite"
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
	sqlDB, err := sql.Open("sqlite", "test.db")
    if err != nil {
        return nil, err
    }
    
    db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    err = db.AutoMigrate(&User{},&Post{},&Tag{})
    return db, err
}

// CreateUserWithPosts creates a new user with associated posts
func CreateUserWithPosts(db *gorm.DB, user *User) error {
    
	result := db.Create(&user)
	
	return result.Error
}

// GetUserWithPosts retrieves a user with all their posts preloaded
func GetUserWithPosts(db *gorm.DB, userID uint) (*User, error) {
	var user User
	result := db.Preload("Posts").First(&user, userID)
	if result.Error!=nil{
	    return nil, result.Error
	}
	return &user, nil
}

// CreatePostWithTags creates a new post with specified tags
func CreatePostWithTags(db *gorm.DB, post *Post, tagNames []string) error {
	tx := db.Begin()
	if tx.Error != nil {
    return tx.Error
}
    ShouldRollback := true
	defer func(){
	    if ShouldRollback{
	        tx.Rollback()
	    }
	}()
	var tags []Tag
	for _, tagName := range tagNames{
	    var tag Tag
	    if err:= tx.Where("name = ?", tagName).First(&tag).Error; err!= nil{
	        if errors.Is(err, gorm.ErrRecordNotFound){
	            tag = Tag{Name: tagName}
	            if err:= tx.Create(&tag).Error; err!=nil{
	                tx.Rollback()
	                return err
	            }
	        } else {
	            tx.Rollback()
	            return err
	        }
	    }
	    tags = append(tags, tag)
	}
	post.Tags = tags
	if err:= tx.Create(post).Error; err!= nil{
	    tx.Rollback()
	    return err
	}
	if err:= tx.Commit().Error; err!=nil{
	    return err
	}
	ShouldRollback = false
	return nil
}

// GetPostsByTag retrieves all posts that have a specific tag
func GetPostsByTag(db *gorm.DB, tagName string) ([]Post, error) {
	var tag Tag
	if err:= db.Where("name = ?", tagName).First(&tag).Error; err!=nil{
	    return nil, err
	}
	var posts []Post
	if err:= db.Model(&tag).
	            Preload("User").
	            Preload("Tags").
	            Association("Posts").
	            Find(&posts); err!=nil{
	                return nil,err
	            }
	return posts, nil
}

// AddTagsToPost adds tags to an existing post
func AddTagsToPost(db *gorm.DB, postID uint, tagNames []string) error {
    tx := db.Begin()
	if tx.Error != nil {
    return tx.Error
}
    shouldRollback:= true
	defer func(){
	    if shouldRollback{
	        tx.Rollback()
	    }
	}()
	var tags []Tag
	for _, tagName := range tagNames{
	    var tag Tag
	    if err:= tx.Where("name = ?", tagName).First(&tag).Error; err!= nil{
	        if errors.Is(err, gorm.ErrRecordNotFound){
	            tag = Tag{Name: tagName}
	            if err:= tx.Create(&tag).Error; err!=nil{
	                tx.Rollback()
	                return err
	            }
	        } else {
	            tx.Rollback()
	            return err
	        }
	    }
	    tags = append(tags, tag)
	}
	var post Post
	if err:= tx.First(&post, postID).Error; err!= nil{
	   tx.Rollback()
	   return err
	}
	if err:= tx.Model(&post).Association("Tags").Append(tags); err!=nil{
	    tx.Rollback()
	    return err
	}
	if err:= tx.Commit().Error; err!=nil{
	    return err
	}
	shouldRollback = false
	return nil
}

// GetPostWithUserAndTags retrieves a post with user and tags preloaded
func GetPostWithUserAndTags(db *gorm.DB, postID uint) (*Post, error) {
    var post Post
	if err:= db.Preload("User").
	            Preload("Tags").
	            First(&post, postID).Error; err!=nil{
	       return nil, err
	   }
	return &post, nil
}
