package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user document in MongoDB
type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Age   int                `bson:"age" json:"age"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Age   int    `json:"age" bson:"age"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Age   int    `json:"age,omitempty" bson:"age,omitempty"`
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// UserService handles user-related database operations
type UserService struct {
	Collection *mongo.Collection
}

func main() {
	// TODO: Connect to MongoDB
	// TODO: Get collection reference
	// TODO: Create UserService instance
	// TODO: Test CRUD operations
}

// CreateUser creates a new user in the database
func (us *UserService) CreateUser(ctx context.Context, req CreateUserRequest) Response {
	// TODO: Validate input data (name, email, age)
	// TODO: Create User with auto-generated ObjectID
	// TODO: Insert user into MongoDB collection
	// TODO: Return success response with created user
	return Response{
		Success: false,
		Error:   "CreateUser not implemented",
		Code:    500,
	}
}

// GetUser retrieves a user by ID from the database
func (us *UserService) GetUser(ctx context.Context, userID string) Response {
	// TODO: Convert userID string to ObjectID
	// TODO: Find user in database by ID
	// TODO: Handle user not found case
	// TODO: Return user data
	return Response{
		Success: false,
		Error:   "GetUser not implemented",
		Code:    500,
	}
}

// UpdateUser updates an existing user in the database
func (us *UserService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) Response {
	// TODO: Convert userID to ObjectID
	// TODO: Update user with $set operator
	// TODO: Check if user was found and modified
	// TODO: Return success response
	return Response{
		Success: false,
		Error:   "UpdateUser not implemented",
		Code:    500,
	}
}

// DeleteUser removes a user from the database
func (us *UserService) DeleteUser(ctx context.Context, userID string) Response {
	// TODO: Convert userID to ObjectID
	// TODO: Delete user from database
	// TODO: Check if user was found and deleted
	// TODO: Return success response
	return Response{
		Success: false,
		Error:   "DeleteUser not implemented",
		Code:    500,
	}
}

// ListUsers retrieves all users from the database
func (us *UserService) ListUsers(ctx context.Context) Response {
	// TODO: Find all users in collection
	// TODO: Iterate through cursor and decode results
	// TODO: Return users array
	return Response{
		Success: false,
		Error:   "ListUsers not implemented",
		Code:    500,
	}
}

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	// TODO: Create client options with URI
	// TODO: Connect to MongoDB
	// TODO: Test connection with Ping
	// TODO: Return client or error
	return nil, fmt.Errorf("ConnectMongoDB not implemented")
}

// Helper function to validate user input
func validateUser(req CreateUserRequest) error {
	// TODO: Check if name is not empty
	// TODO: Check if email is not empty
	// TODO: Check if age is positive
	return nil
}
