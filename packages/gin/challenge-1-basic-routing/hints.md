# Hints for Challenge 1: Basic Routing with Gin

## Hint 1: Setting up the Basic Router

Start with the basic Gin router setup:

```go
router := gin.Default()
```

The `gin.Default()` creates a router with default middleware (logger and recovery). For production you can use `gin.New()` for a clean router.

## Hint 2: Define Your Route Handlers

Create handler functions for each endpoint:

```go
func getUsers(c *gin.Context) {
    // Return the users slice as JSON
    c.JSON(200, users)
}

func createUser(c *gin.Context) {
    // Bind JSON to a User struct and add to users slice
}
```

## Hint 3: Route Structure Patterns

Use the HTTP method functions to define routes:

```go
router.GET("/users", getUsers)
router.POST("/users", createUser)
router.GET("/users/:id", getUserByID)
router.PUT("/users/:id", updateUser)
router.DELETE("/users/:id", deleteUser)
```

## Hint 4: Handling URL Parameters

For routes with parameters like `/users/:id`, access them using:

```go
func getUserByID(c *gin.Context) {
    id := c.Param("id")
    // Convert id to int and find user
    userID, _ := strconv.Atoi(id)
    for _, user := range users {
        if user.ID == userID {
            c.JSON(200, user)
            return
        }
    }
    c.JSON(404, gin.H{"error": "User not found"})
}
```

## Hint 5: Binding JSON Input

For POST/PUT requests, bind the JSON body to your struct:

```go
func createUser(c *gin.Context) {
    var newUser User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    newUser.ID = len(users) + 1
    users = append(users, newUser)
    c.JSON(201, newUser)
}
```

## Hint 6: Starting the Server

Don't forget to start the server on the specified port:

```go
router.Run(":8080")
``` 