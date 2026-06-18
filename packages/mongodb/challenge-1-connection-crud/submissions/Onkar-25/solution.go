package main

import (
	"context"
    "fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"log"
	"errors"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
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
	client,err := ConnectMongoDB("mongodb://localhost:27017")
	defer client.Disconnect(context.Background())
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    
	// TODO: Get collection reference
	 collection := client.Database("user_management").Collection("users")
    userService := &UserService{Collection: collection}
   // Test your implementation
    ctx := context.Background()
    // Create user
    createReq := CreateUserRequest{
        Name:  "John",
        Email: "john@example1.com",
        Age:   30,
    }
    resp := userService.CreateUser(ctx, createReq)
    fmt.Printf("Create: %+v\n", resp)
	// TODO: Create UserService instance
	// TODO: Test CRUD operations
}

// CreateUser creates a new user in the database
func (us *UserService) CreateUser(ctx context.Context, req CreateUserRequest) Response {
	// TODO: Validate input data (name, email, age)
	err := validateUser(req)
	if err != nil{
	    return Response{
		Success: false,
		Error:   err.Error(),
		Code:    400,
	    }
	}
	// TODO: Create User with auto-generated ObjectID
	var user User
	user.ID = primitive.NewObjectID()
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age
	// TODO: Insert user into MongoDB collection
	_, err = us.Collection.InsertOne(ctx, user)
    if err != nil {
        return Response{
		Success: false,
		Error:   "Failed to create user",
		Code:    500,
	    }
    }
	// TODO: Return success response with created user
	return Response{
		Success: true,
		Message : "User created successfully",
		Code:    201,
	}
}

// GetUser retrieves a user by ID from the database
func (us *UserService) GetUser(ctx context.Context, userID string) Response {
	// TODO: Convert userID string to ObjectID
	if userID == ""{
	     return Response{
		Success: false,
		Error:   "User ID is required",
		Code:    400,
	    }
	}
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil{
	    return Response{
		Success: false,
		Error:   "Invalid user ID format",
		Code:    400,
	    }
	}
	
	// TODO: Find user in database by ID
	var user User
	err = us.Collection.FindOne(ctx, bson.M{"ID":id}).Decode(&user)
	if err!= nil{
	    if err == mongo.ErrNoDocuments{
	        return Response{
		        Success: false,
		        Error:   "User not found",
		        Code:    404,
	           }
	    }
	    
	    return Response{
			Success: false,
			Error:   "Database error:",
			Code:    500,
		}
	}
	// TODO: Handle user not found case
	// TODO: Return user data
	return Response{
		Success: true,
	    Data :  user,
		Code:    200,
	}
}

// UpdateUser updates an existing user in the database
func (us *UserService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) Response {
	// TODO: Convert userID to ObjectID
	if userID == ""{
	     return Response{
		    Success: false,
		    Error:   "User ID is required",
		    Code:    400,
	    }
	}
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil{
	    return Response{
		Success: false,
		Error:   "Invalid user ID format",
		Code:    400,
	    }
	}
	// TODO: Update user with $set operator
	update := bson.M{
	    "$set":bson.M{
	        "name":req.Name,
	        "email":req.Email,
	        "age": req.Age,
	    },
	}
	// TODO: Check if user was found and modified
	result,err := us.Collection.UpdateOne(ctx,bson.M{"userID":id},update)
	if err != nil{
	    return Response{
			Success: false,
			Error:   "Invalid email format",
			Code:    400,
		}
	}
	
	if result.MatchedCount == 0{
	    return Response{
			Success: false,
			Error:   "User not found",
			Code:    404,
		}
	}
	// TODO: Return success response
	return Response{
		Success: true,
		Message : "User updated successfully",
		Code:    200,
	}
}

// DeleteUser removes a user from the database
func (us *UserService) DeleteUser(ctx context.Context, userID string) Response {
	// TODO: Convert userID to ObjectID
	if userID == ""{
	     return Response{
		    Success: false,
		    Error:   "User ID is required",
		    Code:    400,
	    }
	}
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil{
	    return Response{
		Success: false,
		Error:   "Invalid user ID format",
		Code:    400,
	    }
	}
	// TODO: Delete user from database
	result , err := us.Collection.DeleteOne(ctx,bson.M{"userID":id})
	if err != nil{
	     return Response{
		    Success: false,
		    Error:   "Database Error",
		    Code:    500,
	    }
	}
	// TODO: Check if user was found and deleted
	if result.DeletedCount == 0{
	     return Response{
		    Success: false,
	    	Error:   "User not found",
		    Code:    404,
	    }
	}
	// TODO: Return success response
	return Response{
		Success: true,
		Message:   "User deleted successfully",
		Code:    200,
	}
}

// ListUsers retrieves all users from the database
func (us *UserService) ListUsers(ctx context.Context) Response {
	// TODO: Find all users in collection
	cursor , err := us.Collection.Find(ctx,bson.M{})
	if err != nil{
	    return Response{
		    Success: false,
		    Error:   "Failed to retrieve users",
		    Code:    500,
	    }
	}
	defer cursor.Close(ctx)
	
	// TODO: Iterate through cursor and decode results
	var users []User
	// TODO: Return users array
	for cursor.Next(ctx){
	    var user User
	    
	    if err := cursor.Decode(&user);err!=nil{
	        return Response{
				Success: false,
				Error:   "error decoding user",
				Code:    500,
			}
	    }
	    users = append(users,user)
	}
	
	if err := cursor.Err(); err != nil {
		return Response{
			Success: false,
			Error:   "cursor error",
			Code:    500,
		}
	}
	return Response{
		Success: true,
		Message : "Retrieved",
		Data:    users,
		Code:    200,
	}
}

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	// TODO: Create client options with URI
	clientOptions:= options.Client().ApplyURI(uri)
	// TODO: Connect to MongoDB
	ctx,cancel:= context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx,clientOptions)
	if err!=nil{
	    return nil, err
	}
	// TODO: Test connection with Ping
	if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }
    
    return client, nil
}

// Helper function to validate user input
func validateUser(req CreateUserRequest) error {
	// TODO: Check if name is not empty
	if req.Name == ""{
	    return errors.New("Name cannot be empty")
	}
	// TODO: Check if email is not empty
	if req.Email == "" {
	    return errors.New("Email cannot be empty")
	}
	
	if !strings.Contains(req.Email,"@"){
	    return errors.New("Invalid email format")
	}
	// TODO: Check if age is positive
	if req.Age <= 0{
	    return errors.New("Age must be greater than 0")
	}
	
	if req.Age > 100{
	    return errors.New("Age must be realistic")
	}
	return nil
}
