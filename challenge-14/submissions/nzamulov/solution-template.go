package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	"io"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/credentials/insecure"
)

// Protocol Buffer definitions (normally would be in .proto files)
// For this challenge, we'll define them as Go structs

// User represents a user in the system
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

// Product represents a product in the catalog
type Product struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Inventory int32   `json:"inventory"`
}

// Order represents an order in the system
type Order struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Total     float64 `json:"total"`
}

// UserService interface
type UserService interface {
	GetUser(ctx context.Context, userID int64) (*User, error)
	ValidateUser(ctx context.Context, userID int64) (bool, error)
}

// ProductService interface
type ProductService interface {
	GetProduct(ctx context.Context, productID int64) (*Product, error)
	CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error)
}

// UserServiceServer implements the UserService
type UserServiceServer struct {
    mu    sync.RWMutex
	users map[int64]*User
}

// NewUserServiceServer creates a new UserServiceServer
func NewUserServiceServer() *UserServiceServer {
	users := map[int64]*User{
		1: {ID: 1, Username: "alice", Email: "alice@example.com", Active: true},
		2: {ID: 2, Username: "bob", Email: "bob@example.com", Active: true},
		3: {ID: 3, Username: "charlie", Email: "charlie@example.com", Active: false},
	}
	return &UserServiceServer{users: users}
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, userID int64) (*User, error) {
    s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return user, nil
}

// ValidateUser checks if a user exists and is active
func (s *UserServiceServer) ValidateUser(ctx context.Context, userID int64) (bool, error) {
    s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	if !exists {
		return false, status.Errorf(codes.NotFound, "user not found")
	}
	return user.Active, nil
}

// ProductServiceServer implements the ProductService
type ProductServiceServer struct {
    mu       sync.RWMutex
	products map[int64]*Product
}

// NewProductServiceServer creates a new ProductServiceServer
func NewProductServiceServer() *ProductServiceServer {
	products := map[int64]*Product{
		1: {ID: 1, Name: "Laptop", Price: 999.99, Inventory: 10},
		2: {ID: 2, Name: "Phone", Price: 499.99, Inventory: 20},
		3: {ID: 3, Name: "Headphones", Price: 99.99, Inventory: 0},
	}
	return &ProductServiceServer{products: products}
}

// GetProduct retrieves a product by ID
func (s *ProductServiceServer) GetProduct(ctx context.Context, productID int64) (*Product, error) {
    s.mu.RLock()
	defer s.mu.RUnlock()
	product, exists := s.products[productID]
	if !exists {
	    return nil, status.Errorf(codes.NotFound, "product not found")
	}
	return product, nil
}

// CheckInventory checks if a product is available in the requested quantity
func (s *ProductServiceServer) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
    s.mu.RLock()
	defer s.mu.RUnlock()
    product, exists := s.products[productID]
    if !exists {
        return false, status.Errorf(codes.NotFound, "product not found")
    }
    return product.Inventory >= quantity, nil
}

// Request/Response types (normally generated from .proto)
type GetUserRequest struct {
	UserId int64 `json:"user_id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

// OrderService handles order creation
type OrderService struct {
	userClient    UserService
	productClient ProductService
	orders        map[int64]*Order
	nextOrderID   int64
	mu sync.Mutex
}

// NewOrderService creates a new OrderService
func NewOrderService(userClient UserService, productClient ProductService) *OrderService {
	return &OrderService{
		userClient:    userClient,
		productClient: productClient,
		orders:        make(map[int64]*Order),
		nextOrderID:   1,
	}
}

// CreateOrder creates a new order
// Note: Inventory check is not atomic with order creation.
// In production, implement atomic inventory reservation.
func (s *OrderService) CreateOrder(ctx context.Context, userID, productID int64, quantity int32) (*Order, error) {
	active, err := s.userClient.ValidateUser(ctx, userID)
	if err != nil {
	    return nil, err
	}
	if !active {
	    return nil, status.Errorf(codes.Unavailable, "user is not active")
	}
	hasNeeded, err := s.productClient.CheckInventory(ctx, productID, quantity)
	if err != nil {
	    return nil, err
	}
	if !hasNeeded {
	    return nil, status.Errorf(codes.Unavailable, "product in needed quantity does not exist")
	}
	product, err := s.productClient.GetProduct(ctx, productID)
	if err != nil {
	    return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	order := Order{
	    ID: s.nextOrderID,
	    UserID: userID,
	    ProductID: productID,
	    Quantity: quantity,
	    Total: product.Price * float64(quantity),
	}
	s.nextOrderID++
	s.orders[order.ID] = &order
	return &order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(orderID int64) (*Order, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
	order, exists := s.orders[orderID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}
	return order, nil
}

// LoggingInterceptor is a server interceptor for logging
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Request received: %s", info.FullMethod)
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Request completed: %s in %v", info.FullMethod, time.Since(start))
	return resp, err
}

// AuthInterceptor is a client interceptor for authentication
func AuthInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add auth token to metadata
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer token123")
	return invoker(ctx, method, req, reply, cc, opts...)
}

// StartUserService starts the user service on the given port
func StartUserService(port string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))
	userServer := NewUserServiceServer()

	// Register HTTP handlers for gRPC methods
	mux := http.NewServeMux()
	mux.HandleFunc("/user/get", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
		    http.Error(w, "invalid user ID", http.StatusBadRequest)
		    return
		}

		user, err := userServer.GetUser(r.Context(), userID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("/user/validate", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
		    http.Error(w, "invalid user ID", http.StatusBadRequest)
		    return
		}

		valid, err := userServer.ValidateUser(r.Context(), userID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"valid": valid})
	})

	go func() {
		log.Printf("User service HTTP server listening on %s", port)
		if err := http.Serve(lis, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return s, nil
}

// StartProductService starts the product service on the given port
func StartProductService(port string) (*grpc.Server, error) {
    lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))
	productServer := NewProductServiceServer()

	// Register HTTP handlers for gRPC methods
	mux := http.NewServeMux()
	mux.HandleFunc("/product/get", func(w http.ResponseWriter, r *http.Request) {
		productIDStr := r.URL.Query().Get("id")
		productID, err := strconv.ParseInt(productIDStr, 10, 64)
		if err != nil {
		    http.Error(w, "invalid product ID", http.StatusBadRequest)
		    return
		}

		product, err := productServer.GetProduct(r.Context(), productID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	})

	mux.HandleFunc("/product/check-inventory", func(w http.ResponseWriter, r *http.Request) {
		productIDStr := r.URL.Query().Get("id")
		productID, err := strconv.ParseInt(productIDStr, 10, 64)
		if err != nil {
		    http.Error(w, "invalid product ID", http.StatusBadRequest)
		    return
		}
		
		quantityStr := r.URL.Query().Get("quantity")
		quantity, err := strconv.ParseInt(quantityStr, 10, 32)
		if err != nil {
		    http.Error(w, "invalid quantity", http.StatusBadRequest)
		    return
		}

		valid, err := productServer.CheckInventory(r.Context(), productID, int32(quantity))
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(valid)
	})
	
	go func() {
		log.Printf("Product service HTTP server listening on %s", port)
		if err := http.Serve(lis, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return s, nil
}

// Connect to both services and return an OrderService
func ConnectToServices(userServiceAddr, productServiceAddr string) (*OrderService, error) {
	userServiceConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to user service: %v", err)
	}
	userService := NewUserServiceClient(userServiceConn)

	productServiceConn, err := grpc.Dial(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
	    userServiceConn.Close()
		return nil, fmt.Errorf("Failed to connect to product service: %v", err)
	}
	productService := NewProductServiceClient(productServiceConn)

	return NewOrderService(userService, productService), nil
}

// Client implementations
type UserServiceClient struct {
	baseURL string
}

func NewUserServiceClient(conn *grpc.ClientConn) UserService {
	return &UserServiceClient{baseURL: fmt.Sprintf("http://%s", conn.Target())}
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID int64) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/user/get?id=%d", c.baseURL, userID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *UserServiceClient) ValidateUser(ctx context.Context, userID int64) (bool, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/user/validate?id=%d", c.baseURL, userID), nil)
	if err != nil {
		return false, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "user not found")
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result["valid"], nil
}

type ProductServiceClient struct {
	baseURL string
}

func NewProductServiceClient(conn *grpc.ClientConn) ProductService {
	return &ProductServiceClient{baseURL: fmt.Sprintf("http://%s", conn.Target())}
}

func (c *ProductServiceClient) GetProduct(ctx context.Context, productID int64) (*Product, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/product/get?id=%d", c.baseURL, productID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *ProductServiceClient) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/product/check-inventory?id=%d&quantity=%d", c.baseURL, productID, quantity), nil)
	if err != nil {
		return false, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "product not found")
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
		return false, err
	}

	bodyStr := strings.TrimSpace(string(bodyBytes))
	boolValue, err := strconv.ParseBool(bodyStr)
	if err != nil {
		return false, err
	}

	return boolValue, nil
}

func main() {
	// Example usage:
	fmt.Println("Challenge 14: Microservices with gRPC")
	fmt.Println("Implement the TODO methods to make the tests pass!")
}
