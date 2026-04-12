package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	db, err := ConnectDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test that tables are created
	assert.True(t, db.Migrator().HasTable(&User{}))
	assert.True(t, db.Migrator().HasTable(&Post{}))
	assert.True(t, db.Migrator().HasTable(&Tag{}))

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()
	os.Remove("test.db")
}

func TestCreateUserWithPosts(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		Posts: []Post{
			{Title: "First Post", Content: "Content 1"},
			{Title: "Second Post", Content: "Content 2"},
		},
	}

	err := CreateUserWithPosts(db, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Verify user was created
	var foundUser User
	db.First(&foundUser, user.ID)
	assert.Equal(t, "John Doe", foundUser.Name)

	// Verify posts were created
	var posts []Post
	db.Where("user_id = ?", user.ID).Find(&posts)
	assert.Len(t, posts, 2)
}

func TestGetUserWithPosts(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create user with posts
	user := &User{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Posts: []Post{
			{Title: "Jane's Post", Content: "Jane's content"},
		},
	}
	CreateUserWithPosts(db, user)

	// Test retrieval
	retrievedUser, err := GetUserWithPosts(db, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", retrievedUser.Name)
	assert.Len(t, retrievedUser.Posts, 1)
	assert.Equal(t, "Jane's Post", retrievedUser.Posts[0].Title)
}

func TestCreatePostWithTags(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create user first
	user := &User{Name: "Author", Email: "author@example.com"}
	CreateUserWithPosts(db, user)

	post := &Post{
		Title:   "Tagged Post",
		Content: "Post with tags",
		UserID:  user.ID,
	}

	tagNames := []string{"golang", "gorm", "database"}

	err := CreatePostWithTags(db, post, tagNames)
	assert.NoError(t, err)
	assert.NotZero(t, post.ID)

	// Verify post was created
	var foundPost Post
	db.First(&foundPost, post.ID)
	assert.Equal(t, "Tagged Post", foundPost.Title)

	// Verify tags were created and associated
	var tags []Tag
	db.Find(&tags)
	assert.Len(t, tags, 3)

	// Verify association
	var postWithTags Post
	db.Preload("Tags").First(&postWithTags, post.ID)
	assert.Len(t, postWithTags.Tags, 3)
}

func TestGetPostsByTag(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create user
	user := &User{Name: "Author", Email: "author@example.com"}
	CreateUserWithPosts(db, user)

	// Create posts with tags
	post1 := &Post{Title: "Go Post", Content: "Go content", UserID: user.ID}
	post2 := &Post{Title: "Database Post", Content: "DB content", UserID: user.ID}

	CreatePostWithTags(db, post1, []string{"golang", "programming"})
	CreatePostWithTags(db, post2, []string{"database", "gorm"})

	// Test retrieval by tag
	posts, err := GetPostsByTag(db, "golang")
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Go Post", posts[0].Title)

	posts, err = GetPostsByTag(db, "database")
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Database Post", posts[0].Title)
}

func TestAddTagsToPost(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create user and post
	user := &User{Name: "Author", Email: "author@example.com"}
	CreateUserWithPosts(db, user)

	post := &Post{Title: "New Post", Content: "Content", UserID: user.ID}
	CreatePostWithTags(db, post, []string{"initial"})

	// Add more tags
	err := AddTagsToPost(db, post.ID, []string{"additional", "new-tag"})
	assert.NoError(t, err)

	// Verify all tags are associated
	var postWithTags Post
	db.Preload("Tags").First(&postWithTags, post.ID)
	assert.Len(t, postWithTags.Tags, 3)

	tagNames := make([]string, len(postWithTags.Tags))
	for i, tag := range postWithTags.Tags {
		tagNames[i] = tag.Name
	}
	assert.Contains(t, tagNames, "initial")
	assert.Contains(t, tagNames, "additional")
	assert.Contains(t, tagNames, "new-tag")
}

func TestGetPostWithUserAndTags(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Create user and post with tags
	user := &User{Name: "Author", Email: "author@example.com"}
	CreateUserWithPosts(db, user)

	post := &Post{Title: "Complete Post", Content: "Complete content", UserID: user.ID}
	CreatePostWithTags(db, post, []string{"tag1", "tag2"})

	// Test retrieval with user and tags
	retrievedPost, err := GetPostWithUserAndTags(db, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Complete Post", retrievedPost.Title)
	assert.Equal(t, "Author", retrievedPost.User.Name)
	assert.Len(t, retrievedPost.Tags, 2)
}

func TestErrorHandling(t *testing.T) {
	db, _ := ConnectDB()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove("test.db")
	}()

	// Test getting non-existent user
	_, err := GetUserWithPosts(db, 999)
	assert.Error(t, err)

	// Test getting non-existent post
	_, err = GetPostWithUserAndTags(db, 999)
	assert.Error(t, err)

	// Test adding tags to non-existent post
	err = AddTagsToPost(db, 999, []string{"tag"})
	assert.Error(t, err)
}
