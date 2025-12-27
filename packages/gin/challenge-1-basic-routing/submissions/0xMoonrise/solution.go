package main

import (
  "encoding/json"
  "errors"
  "io"
  "net/http"
  "strconv"
  "strings"

  "github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
  ID    int    `json:"id"`
  Name  string `json:"name"`
  Email string `json:"email"`
  Age   int    `json:"age"`
}

// Response represents a standard API response
type Response struct {
  Success bool        `json:"success"`
  Data    interface{} `json:"data,omitempty"`
  Message string      `json:"message,omitempty"`
  Error   string      `json:"error,omitempty"`
  Code    int         `json:"code,omitempty"`
}

// In-memory storage
var users = []User{
  {ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
  {ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
  {ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}

var nextID = 4

func main() {
  r := gin.Default()

  r.GET("/users", getAllUsers)
  r.GET("/users/search", searchUsers)
  r.GET("/users/:id", getUserByID)
  r.POST("/users", createUser)
  r.PUT("/users/:id", updateUser)
  r.DELETE("/users/:id", deleteUser)

  r.Run(":8080")
}

func getAllUsers(c *gin.Context) {
  c.JSON(http.StatusOK, Response{
    Success: true,
    Data:    users,
    Message: "Users retrieved successfully",
    Code:    http.StatusOK,
  })
}

func getUserByID(c *gin.Context) {
  paramId := c.Param("id")
  id, err := strconv.Atoi(paramId)

  if err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Code:    http.StatusBadRequest,
      Error:   "Cannot process the request",
    })
    return
  }

  user, index := findUserByID(id)
  if index == -1 {
    c.JSON(http.StatusNotFound, Response{
      Success: false,
      Code:    http.StatusNotFound,
      Error:   "User id not found.",
    })
    return
  }

  c.JSON(http.StatusOK, Response{
    Success: true,
    Message: "User retrive successfully",
    Data:    user,
    Code:    http.StatusOK,
  })

}

// createUser handles POST /users
func createUser(c *gin.Context) {
  jsonData, err := io.ReadAll(c.Request.Body)

  if err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Error:   "Cannot process the JSON data",
      Code:    http.StatusBadRequest,
    })
    return
  }

  user := User{}
  json.Unmarshal(jsonData, &user)
  err = validateUser(user)

  if err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Error:   "Cannot create the user",
      Code:    http.StatusBadRequest,
    })
    return
  }

  user.ID = nextID
  nextID += 1

  users = append(users, user)

  c.JSON(http.StatusCreated, Response{
    Success: true,
    Message: "User has been created.",
    Code:    http.StatusCreated,
    Data:    user,
  })

}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {

  paramId := c.Param("id")
  id, err := strconv.Atoi(paramId)

  if err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Code:    http.StatusBadRequest,
      Error:   "Cannot process the request",
    })
    return
  }

  user, index := findUserByID(id)

  if index == -1 {
    c.JSON(http.StatusNotFound, Response{
      Success: false,
      Code:    http.StatusNotFound,
      Error:   "User id not found.",
    })
    return
  }
    c.JSON(http.StatusNotFound, Response{
      Success: false,
      Code:    http.StatusNotFound,
      Error:   "User id not found.",
    })
    return
  }

  jsonData, err := io.ReadAll(c.Request.Body)

  if err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Error:   "Cannot process the JSON data",
      Code:    http.StatusBadRequest,
    })
    return
  }

  newUser := User{}
  if err = json.Unmarshal(jsonData, &newUser); err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Error:   "Invalid JSON format",
      Code:    http.StatusBadRequest,
    })
    return
  }
  err = validateUser(newUser)

  if err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Error:   "Cannot update the user",
      Code:    http.StatusBadRequest,
    })
    return
  }

  newUser.ID = user.ID
  *user = newUser

  c.JSON(http.StatusOK, Response{
    Success: true,
    Message: "The user has been successfully updated",
    Code:    http.StatusOK,
    Data:    newUser,
  })

}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {

  paramId := c.Param("id")
  id, err := strconv.Atoi(paramId)

  if err != nil {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Code:    http.StatusBadRequest,
      Error:   "Cannot process the request",
    })
    return
  }

  _, index := findUserByID(id)

  if index == -1 {
    c.JSON(http.StatusNotFound, Response{
      Success: false,
      Code:    http.StatusNotFound,
      Error:   "User not found",
    })
    return
  }

  users = append(users[:index], users[index+1:]...)

  c.JSON(http.StatusOK, Response{
    Success: true,
    Code:    http.StatusOK,
    Message: "The user has been deleted",
  })

}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {

  search := c.Query("name")

  if search == "" {
    c.JSON(http.StatusBadRequest, Response{
      Success: false,
      Code:    http.StatusBadRequest,
      Error:   "Missing parameter",
    })
    return
  }

  var founds []User

  for i := 0; i < len(users); i++ {
    if strings.Contains(strings.ToLower(users[i].Name), strings.ToLower(search)) {
      founds = append(founds, users[i])
    }
  }

  if len(founds) == 0 {
    c.JSON(http.StatusOK, Response{
      Success: true,
      Code:    http.StatusOK,
      Error:   "User searched not found",
      Data:    []User{},
    })
    return
  }

  c.JSON(http.StatusOK, Response{
    Success: true,
    Code:    http.StatusOK,
    Message: "User found",
    Data:    founds,
  })

}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
  for i := 0; i < len(users); i++ {
    user := users[i]
    if user.ID == id {
      return &users[i], i
    }
  }
  return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {

  if user.Name == "" {
    return errors.New("Missing 'Name' field")
  }

  if user.Email == "" {
    return errors.New("Missing 'Email' field")
  }

  if !strings.Contains(user.Email, "@") {
    return errors.New("Email invalid format")
  }

  return nil
}