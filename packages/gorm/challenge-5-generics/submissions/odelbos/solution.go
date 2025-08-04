package main

import (
	"context"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User represents a user in the system with company association
type User struct {
	ID        uint     `gorm:"primaryKey"`
	Name      string   `gorm:"not null"`
	Email     string   `gorm:"unique;not null"`
	Age       int      `gorm:"check:age > 0"`
	CompanyID *uint    `gorm:"index"`
	Company   *Company `gorm:"foreignKey:CompanyID"`
	Posts     []Post   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Company represents a company that users can belong to
type Company struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null;unique"`
	Industry    string `gorm:"not null"`
	FoundedYear int    `gorm:"not null"`
	Users       []User `gorm:"foreignKey:CompanyID"`
	CreatedAt   time.Time
}

// Post represents a blog post by a user
type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID"`
	ViewCount int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &Company{}, &Post{}); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUser(ctx context.Context, db *gorm.DB, user *User) error {
	return gorm.G[User](db).Create(ctx, user)
}

func GetUserByID(ctx context.Context, db *gorm.DB, id uint) (*User, error) {
	user, err := gorm.G[User](db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserAge(ctx context.Context, db *gorm.DB, userID uint, age int) error {
	_, err := gorm.G[User](db).Where("id = ?", userID).Update(ctx, "age", age)
	return err
}

func DeleteUser(ctx context.Context, db *gorm.DB, userID uint) error {
	_, err := gorm.G[User](db).Where("id = ?", userID).Delete(ctx)
	return err
}

func CreateUsersInBatches(ctx context.Context, db *gorm.DB, users []User, batchSize int) error {
	return gorm.G[User](db).CreateInBatches(ctx, &users, batchSize)
}

func FindUsersByAgeRange(ctx context.Context, db *gorm.DB, minAge, maxAge int) ([]User, error) {
	return gorm.G[User](db).Where("age BETWEEN ? AND ?", minAge, maxAge).Find(ctx)
}

func UpsertUser(ctx context.Context, db *gorm.DB, user *User) error {
	return gorm.G[User](db, clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		UpdateAll: true,
		}).Create(ctx, user)
}

func CreateUserWithResult(ctx context.Context, db *gorm.DB, user *User) (int64, error) {
	result := gorm.WithResult()
	if err := gorm.G[User](db, result).Create(ctx, user); err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func GetUsersWithCompany(ctx context.Context, db *gorm.DB) ([]User, error) {
	return gorm.G[User](db).
		Joins(clause.RightJoin.Association("Company"), func(_ gorm.JoinBuilder, _, _ clause.Table) error {
			return nil
		}).
		Find(ctx)
}

func GetUsersWithPosts(ctx context.Context, db *gorm.DB, limit int) ([]User, error) {
	return gorm.G[User](db).
		Preload("Posts", func(db gorm.PreloadBuilder) error {
			db.LimitPerRecord(limit)
			return nil
		}).
		Find(ctx)
}

func GetUserWithPostsAndCompany(ctx context.Context, db *gorm.DB, userID uint) (*User, error) {
	user, err := gorm.G[User](db).
		Preload("Posts", func(db gorm.PreloadBuilder) error { return nil }).
		Preload("Company", func(db gorm.PreloadBuilder) error { return nil }).
		Where("id = ?", userID).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SearchUsersInCompany(ctx context.Context, db *gorm.DB, companyName string) ([]User, error) {
	return gorm.G[User](db).
		Joins(clause.InnerJoin.Association("Company"), func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Where("Company.name = ?", companyName)
			return nil
		}).Find(ctx)
}

func GetTopActiveUsers(ctx context.Context, db *gorm.DB, limit int) ([]User, error) {
	type result struct {
		User
		PostCount int
	}
	var rows []result
	err := db.
		WithContext(ctx).
		Model(&User{}).
		Select("users.*, COUNT(posts.id) as post_count").
		Joins("LEFT JOIN posts ON posts.user_id = users.id").
		Group("users.id").
		Order("post_count DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	users := make([]User, len(rows))
	for i, r := range(rows) {
		users[i] = r.User
	}
	return users, nil
}
