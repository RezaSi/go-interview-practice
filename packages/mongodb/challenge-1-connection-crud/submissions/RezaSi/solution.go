package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// CreateUser creates a new user in the database
func (us *UserService) CreateUser(ctx context.Context, req CreateUserRequest) Response {
	// Validate input data
	if req.Name == "" {
		return Response{
			Success: false,
			Error:   "Name cannot be empty",
			Code:    400,
		}
	}

	if req.Email == "" {
		return Response{
			Success: false,
			Error:   "Email cannot be empty",
			Code:    400,
		}
	}

	// Basic email format validation
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return Response{
			Success: false,
			Error:   "Invalid email format",
			Code:    400,
		}
	}

	if req.Age <= 0 {
		return Response{
			Success: false,
			Error:   "Age must be greater than 0",
			Code:    400,
		}
	}

	if req.Age > 150 {
		return Response{
			Success: false,
			Error:   "Age must be realistic (150 or less)",
			Code:    400,
		}
	}

	// Create user with auto-generated ObjectID
	user := User{
		ID:    primitive.NewObjectID(),
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	// Insert the user into MongoDB collection
	_, err := us.Collection.InsertOne(ctx, user)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to create user: " + err.Error(),
			Code:    500,
		}
	}

	// Return success response with created user data
	return Response{
		Success: true,
		Data:    &user,
		Message: "User created successfully",
		Code:    201,
	}
}

// GetUser retrieves a user by ID from the database
func (us *UserService) GetUser(ctx context.Context, userID string) Response {
	// Validate userID is not empty
	if userID == "" {
		return Response{
			Success: false,
			Error:   "User ID is required",
			Code:    400,
		}
	}

	// Convert userID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    400,
		}
	}

	// Create filter using the ObjectID
	filter := bson.M{"_id": objectID}

	// Use FindOne to get the user from database
	var user User
	err = us.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Response{
				Success: false,
				Error:   "User not found",
				Code:    404,
			}
		}
		return Response{
			Success: false,
			Error:   "Database error: " + err.Error(),
			Code:    500,
		}
	}

	// Return success response with user data
	return Response{
		Success: true,
		Data:    &user,
		Code:    200,
	}
}

// UpdateUser updates an existing user in the database
func (us *UserService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) Response {
	// Validate userID is not empty
	if userID == "" {
		return Response{
			Success: false,
			Error:   "User ID is required",
			Code:    400,
		}
	}

	// Convert userID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    400,
		}
	}

	// Validate email format if provided
	if req.Email != "" && (!strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".")) {
		return Response{
			Success: false,
			Error:   "Invalid email format",
			Code:    400,
		}
	}

	// Create filter using the ObjectID
	filter := bson.M{"_id": objectID}

	// Create update document with $set operator
	update := bson.M{"$set": req}

	// Use UpdateOne to update the user
	result, err := us.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to update user: " + err.Error(),
			Code:    500,
		}
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return Response{
			Success: false,
			Error:   "User not found or no changes made",
			Code:    404,
		}
	}

	// Return appropriate response
	return Response{
		Success: true,
		Message: "User updated successfully",
		Code:    200,
	}
}

// DeleteUser removes a user from the database
func (us *UserService) DeleteUser(ctx context.Context, userID string) Response {
	// Validate userID is not empty
	if userID == "" {
		return Response{
			Success: false,
			Error:   "User ID is required",
			Code:    400,
		}
	}

	// Convert userID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    400,
		}
	}

	// Create filter using the ObjectID
	filter := bson.M{"_id": objectID}

	// Use DeleteOne to remove the user
	result, err := us.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to delete user: " + err.Error(),
			Code:    500,
		}
	}

	// Check if any document was deleted
	if result.DeletedCount == 0 {
		return Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		}
	}

	// Return deletion confirmation
	return Response{
		Success: true,
		Message: "User deleted successfully",
		Code:    200,
	}
}

// ListUsers retrieves all users from the database
func (us *UserService) ListUsers(ctx context.Context) Response {
	// Use Find with empty filter to get all users
	cursor, err := us.Collection.Find(ctx, bson.M{})
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to retrieve users: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx) // Don't forget to close the cursor

	// Iterate through cursor to decode all users
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode users: " + err.Error(),
			Code:    500,
		}
	}

	// Handle empty collection gracefully
	if users == nil {
		users = []User{}
	}

	// Return success response with users array
	return Response{
		Success: true,
		Data:    users,
		Message: fmt.Sprintf("Retrieved %d users", len(users)),
		Code:    200,
	}
}

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	// Create client options with the provided URI
	clientOptions := options.Client().ApplyURI(uri)

	// Create context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB using mongo.Connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection with Ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// Example usage and testing
func main() {
	// Example connection string - replace with your MongoDB URI
	mongoURI := "mongodb://localhost:27017"

	// Connect to MongoDB
	client, err := ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	// Get collection reference
	collection := client.Database("user_management").Collection("users")

	// Create user service
	userService := &UserService{Collection: collection}

	// Example operations
	ctx := context.Background()

	// Create a user
	createReq := CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	fmt.Println("Creating user...")
	createResp := userService.CreateUser(ctx, createReq)
	fmt.Printf("Create Response: %+v\n", createResp)

	if createResp.Success {
		user := createResp.Data.(*User)
		userID := user.ID.Hex()

		// Get the user
		fmt.Println("\nGetting user...")
		getResp := userService.GetUser(ctx, userID)
		fmt.Printf("Get Response: %+v\n", getResp)

		// Update the user
		fmt.Println("\nUpdating user...")
		updateReq := UpdateUserRequest{
			Age: 31,
		}
		updateResp := userService.UpdateUser(ctx, userID, updateReq)
		fmt.Printf("Update Response: %+v\n", updateResp)

		// List all users
		fmt.Println("\nListing all users...")
		listResp := userService.ListUsers(ctx)
		fmt.Printf("List Response: %+v\n", listResp)

		// Delete the user
		fmt.Println("\nDeleting user...")
		deleteResp := userService.DeleteUser(ctx, userID)
		fmt.Printf("Delete Response: %+v\n", deleteResp)
	}
}
