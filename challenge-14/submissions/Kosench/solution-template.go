package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	user, exists := s.users[userID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return user, nil
}

// ValidateUser checks if a user exists and is active
func (s *UserServiceServer) ValidateUser(ctx context.Context, userID int64) (bool, error) {
	user, exists := s.users[userID]
	if !exists {
		return false, status.Errorf(codes.NotFound, "user not found")
	}
	return user.Active, nil
}

// ProductServiceServer implements the ProductService
type ProductServiceServer struct {
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
	product, exist := s.products[productID]
	if !exist {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}
	return product, nil
}

// CheckInventory checks if a product is available in the requested quantity
func (s *ProductServiceServer) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
	product, exists := s.products[productID]
	if !exists {
		return false, status.Errorf(codes.NotFound, "product not found")
	}
	return product.Inventory >= quantity, nil
}

// gRPC method handlers for UserService
func (s *UserServiceServer) GetUserRPC(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, err := s.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{User: user}, nil
}

func (s *UserServiceServer) ValidateUserRPC(ctx context.Context, req *ValidateUserRequest) (*ValidateUserResponse, error) {
	valid, err := s.ValidateUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &ValidateUserResponse{Valid: valid}, nil
}

// gRPC method handlers for ProductService
func (s *ProductServiceServer) GetProductRPC(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	product, err := s.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return &GetProductResponse{Product: product}, nil
}

func (s *ProductServiceServer) CheckInventoryRPC(ctx context.Context, req *CheckInventoryRequest) (*CheckInventoryResponse, error) {
	available, err := s.CheckInventory(ctx, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &CheckInventoryResponse{Available: available}, nil
}

// Request/Response types (normally generated from .proto)
type GetUserRequest struct {
	UserId int64 `json:"user_id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type ValidateUserRequest struct {
	UserId int64 `json:"user_id"`
}

type ValidateUserResponse struct {
	Valid bool `json:"valid"`
}

type GetProductRequest struct {
	ProductId int64 `json:"product_id"`
}

type GetProductResponse struct {
	Product *Product `json:"product"`
}

type CheckInventoryRequest struct {
	ProductId int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CheckInventoryResponse struct {
	Available bool `json:"available"`
}

// OrderService handles order creation
type OrderService struct {
	mu            sync.Mutex
	userClient    UserService
	productClient ProductService
	orders        map[int64]*Order
	nextOrderID   int64
	userConn      *grpc.ClientConn
	productConn   *grpc.ClientConn
}

// NewOrderService creates a new OrderService
func NewOrderService(userClient UserService,
	productClient ProductService,
	userConn, productConn *grpc.ClientConn) *OrderService {
	return &OrderService{
		userClient:    userClient,
		productClient: productClient,
		orders:        make(map[int64]*Order),
		nextOrderID:   1,
		userConn:      userConn,
		productConn:   productConn,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, userID, productID int64, quantity int32) (*Order, error) {
	// TODO: Implement this method
	// Hint: 1. Validate user, 2. Get product and check inventory, 3. Create order

	valid, err := s.userClient.ValidateUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}
	if !valid {
		return nil, status.Errorf(codes.PermissionDenied, "user is not active")
	}

	// 2. Get product
	product, err := s.productClient.GetProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 3. Check inventory
	available, err := s.productClient.CheckInventory(ctx, productID, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to check inventory: %w", err)
	}
	if !available {
		return nil, status.Errorf(codes.ResourceExhausted, "insufficient inventory")
	}

	// 4. Create order
	s.mu.Lock()
	defer s.mu.Unlock()
	order := &Order{
		ID:        s.nextOrderID,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Total:     product.Price * float64(quantity),
	}
	s.orders[order.ID] = order
	s.nextOrderID++

	return order, nil
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

// Close closes the gRPC connections
func (s *OrderService) Close() error {
	var err error
	if s.userConn != nil {
		if closeErr := s.userConn.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if s.productConn != nil {
		if closeErr := s.productConn.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; %v", err, closeErr)
			} else {
				err = closeErr
			}
		}
	}
	return err
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
			http.Error(w, "invalid user id", http.StatusBadRequest)
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
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

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
	// TODO: Implement this function
	// Hint: create listener, gRPC server with interceptor, register service, serve

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))
	productServer := NewProductServiceServer()

	// Register HTTP handlers for gRPC methods
	mux := http.NewServeMux()

	mux.HandleFunc("/product/get", func(w http.ResponseWriter, r *http.Request) {
		productIDStr := r.URL.Query().Get("id")
		productID, _ := strconv.ParseInt(productIDStr, 10, 64)

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
		productID, _ := strconv.ParseInt(productIDStr, 10, 64)
		quantityStr := r.URL.Query().Get("quantity")
		quantity, _ := strconv.ParseInt(quantityStr, 10, 32)

		available, err := productServer.CheckInventory(r.Context(), productID, int32(quantity))
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"available": available})
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
	// TODO: Implement this function
	// Hint: create gRPC connections with interceptors, create clients, return OrderService
	userConn, err := grpc.NewClient(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AuthInterceptor))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	productConn, err := grpc.NewClient(
		productServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AuthInterceptor))
	if err != nil {
		userConn.Close()
		return nil, fmt.Errorf("failed to connect to product service %w", err)
	}

	userClient := NewUserServiceClient(userConn)
	productClient := NewProductServiceClient(productConn)

	return NewOrderService(userClient, productClient, userConn, productConn), nil
}

// Client implementations
type UserServiceClient struct {
	baseURL string
}

func NewUserServiceClient(conn *grpc.ClientConn) UserService {
	// Extract address from connection for HTTP calls
	// In a real gRPC implementation, this would use the connection directly
	target := conn.Target()
	return &UserServiceClient{baseURL: "http://" + target}
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID int64) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/user/get?id=%d", c.baseURL, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *UserServiceClient) ValidateUser(ctx context.Context, userID int64) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/user/validate?id=%d", c.baseURL, userID))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "user not found")
	}

	var result map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result["valid"], nil
}

type ProductServiceClient struct {
	conn    *grpc.ClientConn
	baseURL string
}

func NewProductServiceClient(conn *grpc.ClientConn) ProductService {
	target := conn.Target()
	baseURL := "http://" + target
	return &ProductServiceClient{
		conn:    conn,
		baseURL: baseURL}
}

func (c *ProductServiceClient) GetProduct(ctx context.Context, productID int64) (*Product, error) {
	// TODO: Implement gRPC client call
	// Hint: make gRPC call to GetProductRPC method
	resp, err := http.Get(fmt.Sprintf("%s/product/get?id=%d", c.baseURL, productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *ProductServiceClient) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
	// TODO: Implement gRPC client call
	// Hint: make gRPC call to CheckInventoryRPC method
	resp, err := http.Get(fmt.Sprintf("%s/product/check-inventory?id=%d&quantity=%d", c.baseURL, productID, quantity))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "product not found")
	}

	var result map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result["available"], nil
}

// gRPC service registration helpers
func RegisterUserServiceServer(s *grpc.Server, srv *UserServiceServer) {
	// In a real implementation, this would be generated code
	// For this challenge, we'll manually handle the registration
}

func RegisterProductServiceServer(s *grpc.Server, srv *ProductServiceServer) {
	// In a real implementation, this would be generated code
	// For this challenge, we'll manually handle the registration
}

func main() {
	// Example usage:
	fmt.Println("Challenge 14: Microservices with gRPC")
	fmt.Println("Implement the TODO methods to make the tests pass!")
}
