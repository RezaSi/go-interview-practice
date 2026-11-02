package main

import (
	"context"
	"fmt"
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

// ConnectDB establishes a connection to the SQLite database and auto-migrates models
func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &Company{}, &Post{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateUser creates a new user using GORM's generics API
func CreateUser(ctx context.Context, db *gorm.DB, user *User) error {
	err := gorm.G[User](db).Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by ID using generics API
func GetUserByID(ctx context.Context, db *gorm.DB, id uint) (*User, error) {
	user, err := gorm.G[User](db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserAge updates a user's age using generics API
func UpdateUserAge(ctx context.Context, db *gorm.DB, userID uint, age int) error {
	rows, err := gorm.G[User](db).Where("id = ?", userID).Update(ctx, "age", age)
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user by ID using generics API
func DeleteUser(ctx context.Context, db *gorm.DB, userID uint) error {
	rows, err := gorm.G[User](db).Where("id = ?", userID).Delete(ctx)
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	if err != nil {
		return err
	}

	return nil
}

// CreateUsersInBatches creates multiple users in batches using generics API
func CreateUsersInBatches(ctx context.Context, db *gorm.DB, users []User, batchSize int) error {
	if batchSize <= 0 {
		return fmt.Errorf("invalid batch size value")
	}

	if len(users) == 0 {
		return fmt.Errorf("empty users slice")
	}

	err := gorm.G[User](db).CreateInBatches(ctx, &users, batchSize)
	if err != nil {
		return err
	}

	return nil
}

// FindUsersByAgeRange finds users within an age range using generics API
func FindUsersByAgeRange(ctx context.Context, db *gorm.DB, minAge, maxAge int) ([]User, error) {
	if minAge < 0 || maxAge < 0 {
		return nil, fmt.Errorf("invalid min or max age value")
	}

	if minAge > maxAge {
		return nil, fmt.Errorf("min age value is higher than max age value")
	}

	users, err := gorm.G[User](db).Where("age BETWEEN ? AND ?", minAge, maxAge).Find(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpsertUser creates or updates a user handling conflicts using OnConflict
func UpsertUser(ctx context.Context, db *gorm.DB, user *User) error {
	err := gorm.G[User](db, clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "age", "updated_at"}),
	}).Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserWithResult creates a user and returns result metadata
func CreateUserWithResult(ctx context.Context, db *gorm.DB, user *User) (int64, error) {
	res := gorm.WithResult()

	err := gorm.G[User](db, res).Create(ctx, user)
	if err != nil {
		return 0, err
	}

	if res.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return res.RowsAffected, nil
}

// GetUsersWithCompany retrieves users with their company information using enhanced joins
func GetUsersWithCompany(ctx context.Context, db *gorm.DB) ([]User, error) {
	users, err := gorm.G[User](db).Joins(
		clause.LeftJoin.Association("Company"),
		func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error { return nil },
	).Find(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersWithPosts retrieves users with their posts using enhanced preloading
func GetUsersWithPosts(ctx context.Context, db *gorm.DB, limit int) ([]User, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit value")
	}

	users, err := gorm.G[User](db).Preload("Posts", func(db gorm.PreloadBuilder) error {
		db.LimitPerRecord(limit)
		return nil
	}).Find(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserWithPostsAndCompany retrieves a user with both posts and company preloaded
func GetUserWithPostsAndCompany(ctx context.Context, db *gorm.DB, userID uint) (*User, error) {
	user, err := gorm.G[User](db).
		Where("id = ?", userID).
		Preload("Posts", func(db gorm.PreloadBuilder) error { return nil }).
		Preload("Company", func(db gorm.PreloadBuilder) error { return nil }).
		First(ctx)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SearchUsersInCompany finds users working in a specific company using join with filters
func SearchUsersInCompany(ctx context.Context, db *gorm.DB, companyName string) ([]User, error) {
	users, err := gorm.G[User](db).
		Joins(clause.InnerJoin.Association("Company"), func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
			db.Where("Company.name = ?", companyName)
			return nil
		}).Find(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetTopActiveUsers retrieves users with the most posts using complex joins and grouping
func GetTopActiveUsers(ctx context.Context, db *gorm.DB, limit int) ([]User, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit value")
	}

	users, err := gorm.G[User](db).
		Joins(
			clause.JoinTarget{
				Type:  clause.LeftJoin,
				Table: clause.Table{Name: "posts"},
			},
			func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
				db.Where("users.id = posts.user_id")
				return nil
			},
		).
		Group("users.id").Order("COUNT(posts.id) DESC").
		Limit(limit).
		Find(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}
