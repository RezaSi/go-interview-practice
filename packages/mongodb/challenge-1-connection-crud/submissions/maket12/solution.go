package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Age   int                `bson:"age" json:"age"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

type UserService struct {
	Collection *mongo.Collection
}

func validateUser(name, email string, age int) error {
	if name == "" {
		return errors.New("Name cannot be empty")
	}
	if email == "" {
		return errors.New("Email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("Invalid email format")
	}
	if age <= 0 {
		return errors.New("Age must be greater than 0")
	}
	if age > 150 {
		return errors.New("Age must be realistic")
	}
	return nil
}

func (us *UserService) CreateUser(ctx context.Context, req CreateUserRequest) Response {
	if err := validateUser(req.Name, req.Email, req.Age); err != nil {
		return Response{Success: false, Message: "failed to create user", Error: err.Error(), Code: 400}
	}

	user := User{
		ID:    primitive.NewObjectID(),
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	_, err := us.Collection.InsertOne(ctx, user)
	if err != nil {
    return Response{
        Success: false,
        Message: "failed to create user",
        Error:   "Failed to create user: " + err.Error(), 
        Code:    500,
    }
}

	return Response{Success: true, Data: user, Message: "User created successfully", Code: 201}
}

func (us *UserService) GetUser(ctx context.Context, userID string) Response {
	if userID == "" {
		return Response{Success: false, Message: "failed to get user", Error: "User ID is required", Code: 400}
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Response{Success: false, Message: "failed to get user", Error: "Invalid user ID format", Code: 400}
	}

	var user User
	err = us.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Response{Success: false, Error: "User not found", Code: 404}
		}
		return Response{Success: false, Error: "Database error: " + err.Error(), Code: 500}
	}

	return Response{Success: true, Data: user, Code: 200}
}

func (us *UserService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) Response {
	if userID == "" {
		return Response{Success: false, Error: "User ID is required", Code: 400}
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Response{Success: false, Error: "Invalid user ID format", Code: 400}
	}

	if req.Email != "" && !strings.Contains(req.Email, "@") {
		return Response{Success: false, Error: "Invalid email format", Code: 400}
	}

	update := bson.M{"$set": bson.M{"name": req.Name, "email": req.Email, "age": req.Age}}
	res, err := us.Collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return Response{Success: false, Error: err.Error(), Code: 500}
	}

	if res.MatchedCount == 0 {
		return Response{Success: false, Error: "User not found", Code: 404}
	}

	return Response{Success: true, Message: "User updated successfully", Code: 200}
}

func (us *UserService) DeleteUser(ctx context.Context, userID string) Response {
	if userID == "" {
		return Response{Success: false, Error: "User ID is required", Code: 400}
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return Response{Success: false, Error: "Invalid user ID format", Code: 400}
	}

	res, err := us.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return Response{Success: false, Error: err.Error(), Code: 500}
	}

	if res.DeletedCount == 0 {
		return Response{Success: false, Error: "User not found", Code: 404}
	}

	return Response{Success: true, Message: "User deleted successfully", Code: 200}
}

func (us *UserService) ListUsers(ctx context.Context) Response {
	cursor, err := us.Collection.Find(ctx, bson.M{})
	if err != nil {
		return Response{Success: false, Error: "Failed to retrieve users", Code: 500}
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		return Response{Success: false, Error: err.Error(), Code: 500}
	}

	return Response{
		Success: true,
		Data:    users,
		Message: fmt.Sprintf("Retrieved %d users", len(users)),
		Code:    200,
	}
}

func main() {}