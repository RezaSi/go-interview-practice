package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id" binding:"-"`
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,gt=0,lte=150"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// for test-
var (
	users         []User
	nextID        int
	globalService *DefaultUserService
)

// Repository
// =======================================================================
type UserRepository interface {
	GetAll() ([]*User, error)
	GetByID(id int) (*User, error)
	Create(book *User) error
	Update(id int, book *User) error
	Delete(id int) error
	SearchByName(name string) ([]*User, error)
}

type InMemoryUserRepository struct {
	users   map[string]*User
	idIndex map[int]string
	mu      sync.RWMutex
	cnt     int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:   make(map[string]*User),
		idIndex: make(map[int]string),
		cnt:     1,
	}
}

var (
	ErrUserRepositoryEmpty      = errors.New("not a single user was found")
	ErrUserRepositoryIdNotFound = errors.New("no user with this ID was found")
	ErrUserRepositoryCantCreate = errors.New("user is invalid, cannot create user")
)

func (d *InMemoryUserRepository) createHashByUser(user *User) string {
	name := strings.ToLower(strings.TrimSpace(user.Name))
	email := strings.ToLower(strings.TrimSpace(user.Email))

	input := fmt.Sprintf("%s|%s|%d",
		name,
		email,
		user.Age,
	)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (d *InMemoryUserRepository) validateUser(user *User) error {
	if user.Name == "" {
		return fmt.Errorf("%w: name is empty", ErrUserRepositoryCantCreate)
	}

	if user.Email == "" {
		return fmt.Errorf("%w: email is empty", ErrUserRepositoryCantCreate)
	}

	if user.Age <= 0 {
		return fmt.Errorf("%w: age must be positive", ErrUserRepositoryCantCreate)
	}

	return nil
}

func (d *InMemoryUserRepository) GetAll() ([]*User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var users []*User

	for _, v := range d.users {
		users = append(users, v)
	}

	return users, nil
}

func (d *InMemoryUserRepository) GetByID(id int) (*User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if hashKey, exists := d.idIndex[id]; exists {
		return d.users[hashKey], nil
	}
	return nil, ErrUserRepositoryIdNotFound
}

func (d *InMemoryUserRepository) Create(user *User) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.validateUser(user); err != nil {
		return err
	}

	user.ID = d.cnt
	hashUser := d.createHashByUser(user)

	if _, exists := d.users[hashUser]; exists {
		return fmt.Errorf("%w: there is a similar user", ErrUserRepositoryCantCreate)
	}

	d.users[hashUser] = user
	d.idIndex[user.ID] = hashUser
	d.cnt++

	return nil
}

func (d *InMemoryUserRepository) Update(id int, user *User) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.validateUser(user); err != nil {
		return err
	}

	oldHash, exists := d.idIndex[id]
	if !exists {
		return ErrUserRepositoryIdNotFound
	}

	user.ID = id
	newHash := d.createHashByUser(user)

	if oldHash != newHash {
		if _, exists := d.users[newHash]; exists {
			return fmt.Errorf("%w: there is a similar user", ErrUserRepositoryCantCreate)
		}
	}

	delete(d.users, oldHash)
	d.users[newHash] = user
	d.idIndex[id] = newHash

	return nil
}

func (d *InMemoryUserRepository) Delete(id int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	hashKey, exists := d.idIndex[id]
	if !exists {
		return ErrUserRepositoryIdNotFound
	}

	delete(d.users, hashKey)
	delete(d.idIndex, id)

	return nil
}

func (d *InMemoryUserRepository) SearchBy(predicate func(*User) bool) ([]*User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	users := make([]*User, 0)
	for _, user := range d.users {
		if predicate(user) {
			users = append(users, user)
		}
	}

	return users, nil
}

func (d *InMemoryUserRepository) SearchByName(name string) ([]*User, error) {
	return d.SearchBy(func(user *User) bool {
		return strings.Contains(strings.ToLower(user.Name), strings.ToLower(name))
	})
}

// Service
// =======================================================================
type UserService interface {
	GetAllUsers() ([]*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
	UpdateUser(id int, book *User) error
	DeleteUser(id int) error
	SearchUsersByName(name string) ([]*User, error)
}

type DefaultUserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *DefaultUserService {
	return &DefaultUserService{
		repo: repo,
	}
}

func (d *DefaultUserService) GetAllUsers() ([]*User, error) {
	return d.repo.GetAll()
}

func (d *DefaultUserService) GetUserByID(id int) (*User, error) {
	return d.repo.GetByID(id)
}

func (d *DefaultUserService) CreateUser(book *User) error {
	return d.repo.Create(book)
}

func (d *DefaultUserService) UpdateUser(id int, book *User) error {
	return d.repo.Update(id, book)
}

func (d *DefaultUserService) DeleteUser(id int) error {
	return d.repo.Delete(id)
}

func (d *DefaultUserService) SearchUsersByName(name string) ([]*User, error) {
	return d.repo.SearchByName(name)
}

// Middleware
// =======================================================================
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("Response: %d\n", c.Writer.Status())
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			log.Printf("Errors: %v\n",
				Response{
					Success: false,
					Error:   "Internal server error",
					Code:    http.StatusInternalServerError,
				})
		}
	}
}

func OverrideHtmlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" {
			method := c.Request.FormValue("_method")
			fmt.Printf("Original method: POST, Override method: %s\n", method)
			if method == "PUT" || method == "DELETE" {
				c.Request.Method = method
				fmt.Printf("Method overridden to: %s\n", method)
			}
		}
		c.Next()
	}
}

// Helper
// =======================================================================
func sendJSONResponse(c *gin.Context, status int, success bool, data interface{}, message string) {
	response := Response{
		Success: success,
		Data:    data,
		Message: message,
		Code:    status,
	}
	c.JSON(status, response)
}

func sendErrorResponse(c *gin.Context, status int, message string) {
	sendJSONResponse(c, status, false, nil, message)
}

// Handler
// =======================================================================
func getAllUsers(c *gin.Context) {

	users, err := getServiceData().GetAllUsers()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		response := Response{
			Success: true,
			Data:    users,
			Message: "Users retrieved successfully",
			Code:    http.StatusOK,
		}
		c.HTML(http.StatusOK, "users", gin.H{
			"title":    "All Users",
			"response": response,
		})
		return
	}

	sendJSONResponse(c, http.StatusOK, true, users, "Users retrieved successfully")
}

func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := getServiceData().GetUserByID(id)
	if err != nil {
		sendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		response := Response{
			Success: true,
			Data:    []*User{user},
			Message: "Users retrieved successfully",
			Code:    http.StatusOK,
		}
		c.HTML(http.StatusOK, "users", gin.H{
			"title":    "All Users",
			"response": response,
		})
		return
	}

	sendJSONResponse(c, http.StatusOK, true, user, "Users retrieved successfully")
}

func createUser(c *gin.Context) {
	var user User

	if err := c.ShouldBind(&user); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid user data: "+err.Error())
		return

	}

	if err := getServiceData().CreateUser(&user); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/users")

	sendJSONResponse(c, http.StatusCreated, true, user, "User Create successfully")
}

func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user User
	if err := c.ShouldBind(&user); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid user data: "+err.Error())
		return

	}

	if err := getServiceData().UpdateUser(id, &user); err != nil {
		sendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	sendJSONResponse(c, http.StatusOK, true, user, "User Update successfully")
}

func deleteUser(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := getServiceData().DeleteUser(id); err != nil {
		sendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	sendJSONResponse(c, http.StatusOK, true, id, "User Delete successfully")

}

func searchUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		sendErrorResponse(c, http.StatusBadRequest, "Name parameter is required")
		return
	}

	users, err := getServiceData().SearchUsersByName(name)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		response := Response{
			Success: true,
			Data:    users,
			Message: "Users retrieved successfully",
			Code:    http.StatusOK,
		}
		c.HTML(http.StatusOK, "users", gin.H{
			"title":    "All Users",
			"response": response,
		})
		return
	}

	sendJSONResponse(c, http.StatusOK, true, users, "Users retrieved successfully")
}

// X. Main
// =======================================================================
// for test
func getServiceData() *DefaultUserService {
	if globalService == nil {
		repo := NewInMemoryUserRepository()
		globalService = NewUserService(repo)
		/*users = []User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
			{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
		}
		nextID = 4

		for _, user := range users {
			globalService.CreateUser(&user)
		}*/

		globalService.CreateUser(&User{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30})
		globalService.CreateUser(&User{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25})
		globalService.CreateUser(&User{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35})
		nextID = 4
	}

	return globalService
}

func main() {

	//repo := NewInMemoryUserRepository()
	//service := NewUserService(repo)

	//service.CreateUser(&User{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30})
	//service.CreateUser(&User{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25})
	//service.CreateUser(&User{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35})

	r := gin.Default()

	//r.Use(OverrideHtmlMiddleware())
	r.Use(LoggingMiddleware())
	r.Use(ErrorMiddleware())

	tmpl := template.Must(template.New("users").Parse(usersHTML))
	r.SetHTMLTemplate(tmpl)

	r.GET("/", func(c *gin.Context) {
		getAllUsers(c)
	})

	r.GET("/users", func(c *gin.Context) {
		getAllUsers(c)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		getUserByID(c)
	})

	r.POST("/users", func(c *gin.Context) {
		createUser(c)
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		updateUser(c)
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		deleteUser(c)
	})

	r.GET("/users/search", func(c *gin.Context) {
		searchUsers(c)
	})

	log.Println("http://localhost:8085/users")
	r.Run(":8085")
}

// HTML
// =======================================================================

const usersHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f5f5f5;
            padding: 20px;
        }
        
        .container {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 30px;
            max-width: 1400px;
            margin: 0 auto;
        }
        
        .section {
            background: white;
            padding: 25px;
            border-radius: 12px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        
        h1 {
            color: #2c3e50;
            margin-bottom: 20px;
            text-align: center;
            font-size: 28px;
        }
        
        h2 {
            color: #34495e;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #3498db;
            font-size: 22px;
        }
        
        h3 {
            color: #2c3e50;
            margin: 15px 0 10px 0;
            font-size: 18px;
        }
        
        .users-table {
            width: 100%;
            border-collapse: collapse;
            margin: 15px 0;
        }
        
        .users-table th,
        .users-table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ecf0f1;
        }
        
        .users-table th {
            background-color: #34495e;
            color: white;
            font-weight: 600;
        }
        
        .users-table tr:hover {
            background-color: #f8f9fa;
        }
        
        .user-id {
            font-weight: bold;
            color: #2c3e50;
        }
        
        .user-email {
            color: #7f8c8d;
            font-size: 14px;
        }
        
        .user-age {
            text-align: center;
            background-color: #ecf0f1;
            border-radius: 15px;
            padding: 4px 8px;
            font-size: 14px;
        }
        
        .form-group {
            margin-bottom: 15px;
        }
        
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: 600;
            color: #2c3e50;
            font-size: 14px;
        }
        
        input, button, select {
            width: 100%;
            padding: 12px;
            border: 2px solid #bdc3c7;
            border-radius: 6px;
            font-size: 14px;
        }
        
        input:focus {
            outline: none;
            border-color: #3498db;
            box-shadow: 0 0 5px rgba(52, 152, 219, 0.3);
        }
        
        button {
            background: linear-gradient(135deg, #3498db, #2980b9);
            color: white;
            border: none;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s ease;
            margin-top: 10px;
        }
        
        button:hover {
            background: linear-gradient(135deg, #2980b9, #3498db);
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.2);
        }
        
        .success {
            color: #27ae60;
            background-color: #d5f4e6;
            padding: 12px;
            border-radius: 6px;
            margin: 15px 0;
            border-left: 4px solid #27ae60;
        }
        
        .error {
            color: #e74c3c;
            background-color: #fadbd8;
            padding: 12px;
            border-radius: 6px;
            margin: 15px 0;
            border-left: 4px solid #e74c3c;
        }
        
        .form-section {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            margin: 15px 0;
            border: 1px solid #e9ecef;
        }
        
        .quick-actions {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 10px;
            margin: 15px 0;
        }
        
        .quick-btn {
            padding: 8px;
            font-size: 13px;
        }

        .btn-create {
            background: linear-gradient(135deg, #27ae60, #229954);
        }

        .btn-create:hover {
            background: linear-gradient(135deg, #229954, #27ae60);
        }
    </style>
</head>
<body>
    <h1>{{.title}}</h1>
    
    <div class="container">
        <!-- Левая колонка: Список пользователей -->
        <div class="section">
            <h2>Users List</h2>
            
            {{if .response.Success}}
                <p class="success">✅ {{.response.Message}}</p>
                {{if .response.Data}}
                    <table class="users-table">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>Name</th>
                                <th>Email</th>
                                <th>Age</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .response.Data}}
                                <tr>
                                    <td class="user-id">{{.ID}}</td>
                                    <td>{{.Name}}</td>
                                    <td class="user-email">{{.Email}}</td>
                                    <td class="user-age">{{.Age}}</td>
                                </tr>
                            {{end}}
                        </tbody>
                    </table>
                {{else}}
                    <p>No users found</p>
                {{end}}
            {{else}}
                <p class="error">❌ {{.response.Error}}</p>
            {{end}}
        </div>

        <!-- Правая колонка: Формы -->
        <div class="section">
            <h2>Test Forms</h2>
            
            <!-- Quick Actions -->
            <div class="quick-actions">
                <button onclick="location.href='/users'" class="quick-btn">View All Users</button>
                <button onclick="location.href='/users/1'" class="quick-btn">View User #1</button>
            </div>

            <!-- Get Users -->
            <div class="form-section">
                <h3>📋 Get Users</h3>
                <form action="/users" method="GET">
                    <button type="submit">Get All Users</button>
                </form>
                
                <form id="getUserForm" method="GET" onsubmit="updateFormAction('getUserForm', '/users/')">
                    <div class="form-group">
                        <label for="getUserFormId">Get User by ID:</label>
                        <input type="number" id="getUserFormId" placeholder="Enter User ID" required>
                    </div>
                    <button type="submit">Get User</button>
                </form>

                <form action="/users/search" method="GET">
                    <div class="form-group">
                        <label for="searchName">Search by Name:</label>
                        <input type="text" id="searchName" name="name" placeholder="Enter name">
                    </div>
                    <button type="submit">Search Users</button>
                </form>
            </div>

            <!-- Create User -->
            <div class="form-section">
                <h3>➕ Create User</h3>
                <div class="form-group" visibility: hidden>
                    <label for="createId">ID:</label>
                    <input type="number" id="createId" placeholder="Enter ID" value=1 required>
                </div>
                <div class="form-group">
                    <label for="createName">Name:</label>
                    <input type="text" id="createName" placeholder="Enter name" required>
                </div>
                <div class="form-group">
                    <label for="createEmail">Email:</label>
                    <input type="email" id="createEmail" placeholder="Enter email" required>
                </div>
                <div class="form-group">
                    <label for="createAge">Age:</label>
                    <input type="number" id="createAge" placeholder="Enter age" required>
                </div>
                <button type="button" onclick="createUserFetch()" class="btn-create">Create User</button>
            </div>

            <!-- Update User -->
            <div class="form-section">
                <h3>✏️ Update User</h3>
                <div class="form-group">
                    <label for="updateUserId">User ID:</label>
                    <input type="number" id="updateUserId" placeholder="User ID to update" required>
                </div>
                <div class="form-group">
                    <label for="updateName">New Name:</label>
                    <input type="text" id="updateName" placeholder="New name" required>
                </div>
                <div class="form-group">
                    <label for="updateEmail">New Email:</label>
                    <input type="email" id="updateEmail" placeholder="New email" required>
                </div>
                <div class="form-group">
                    <label for="updateAge">New Age:</label>
                    <input type="number" id="updateAge" placeholder="New age" required>
                </div>
                <button type="button" onclick="updateUserFetch()">Update User</button>
            </div>

            <!-- Delete User -->
            <div class="form-section">
                <h3>🗑️ Delete User</h3>
                <div class="form-group">
                    <label for="fetchUserId">User ID to delete:</label>
                    <input type="number" id="fetchUserId" placeholder="Enter User ID" required>
                </div>
                <button onclick="deleteUserFetch()">Delete User</button>
            </div>
        </div>
    </div>

    <script>
        function updateFormAction(formId, urlSuffix) {
            const form = document.getElementById(formId);
            const id = document.getElementById(formId + 'Id').value;
            if (id) {
                form.action = urlSuffix + id;
            }
        }

        function createUserFetch() {
            const id = document.getElementById('createId').value;
            const name = document.getElementById('createName').value;
            const email = document.getElementById('createEmail').value;
            const age = document.getElementById('createAge').value;

            if (id && name && email && age) {
                const userData = {
					ID: 0,
                    Name: name,
                    Email: email,
                    Age: parseInt(age)
                };

                fetch('/users', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(userData)
                })
                .then(response => {
                    if (response.redirected) {
                        window.location.href = response.url;
                    } else if (response.ok) {
                        return response.json();
                    } else {
                        throw new Error('Network response was not ok');
                    }
                })
                .then(data => {
                    console.log('User created:', data);
                    alert('User created successfully!');
                    window.location.reload();
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('Error creating user: ' + error.message);
                });
            } else {
                alert('Please fill all fields');
            }
        }

        function updateUserFetch() {
            const userId = document.getElementById('updateUserId').value;
            const name = document.getElementById('updateName').value;
            const email = document.getElementById('updateEmail').value;
            const age = document.getElementById('updateAge').value;

            if (userId && name && email && age) {
                const userData = {
                    ID: parseInt(userId),
                    Name: name,
                    Email: email,
                    Age: parseInt(age)
                };

                fetch('/users/' + userId, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(userData)
                })
                .then(response => {
                    if (response.redirected) {
                        window.location.href = response.url;
                    } else if (response.ok) {
                        return response.json();
                    } else {
                        throw new Error('Network response was not ok');
                    }
                })
                .then(data => {
                    console.log('User updated:', data);
                    alert('User updated successfully!');
                    window.location.reload();
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('Error updating user: ' + error.message);
                });
            } else {
                alert('Please fill all fields');
            }
        }

        function deleteUserFetch() {
            const userId = document.getElementById('fetchUserId').value;
            if (userId) {
                fetch('/users/' + userId, {
                    method: 'DELETE',
                })
                .then(response => {
                    if (response.redirected) {
                        window.location.href = response.url;
                    } else {
                        return response.json();
                    }
                })
                .then(data => {
                    console.log('User deleted:', data);
                    window.location.reload();
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('Error deleting user: ' + error.message);
                });
            }
        }
    </script>
</body>
</html>`
