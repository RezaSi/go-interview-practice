# Hints for Challenge 14: Microservices with gRPC

## Hint 1: gRPC Service Implementation Structure
Start by implementing the basic service structure with proper error handling:
```go
import (
    "context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type userServiceServer struct {
    users map[string]*User // in-memory store for demo
}

func (s *userServiceServer) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
    if req.UserId == "" {
        return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
    }
    
    user, exists := s.users[req.UserId]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "user not found")
    }
    
    return &GetUserResponse{User: user}, nil
}
```

## Hint 2: gRPC Status Codes for Business Logic
Use appropriate gRPC status codes for different error conditions:
```go
func (s *userServiceServer) ValidateUser(ctx context.Context, req *ValidateUserRequest) (*ValidateUserResponse, error) {
    user, exists := s.users[req.UserId]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "user not found")
    }
    
    if !user.IsActive {
        return nil, status.Errorf(codes.PermissionDenied, "user is not active")
    }
    
    return &ValidateUserResponse{IsValid: true}, nil
}

func (s *productServiceServer) CheckInventory(ctx context.Context, req *CheckInventoryRequest) (*CheckInventoryResponse, error) {
    product, exists := s.products[req.ProductId]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "product not found")
    }
    
    if product.Stock < req.Quantity {
        return nil, status.Errorf(codes.ResourceExhausted, "insufficient inventory")
    }
    
    return &CheckInventoryResponse{Available: true}, nil
}
```

## Hint 3: Setting Up gRPC Server with Interceptors
Create servers with logging and authentication interceptors:
```go
func StartUserService(port string) (*grpc.Server, error) {
    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return nil, err
    }
    
    server := grpc.NewServer(
        grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
            LoggingInterceptor,
            AuthInterceptor,
        )),
    )
    
    userService := &userServiceServer{
        users: make(map[string]*User),
    }
    // Register your proto service here
    RegisterUserServiceServer(server, userService)
    
    go func() {
        server.Serve(lis)
    }()
    
    return server, nil
}
```

## Hint 4: Implementing Interceptors
Create interceptors for cross-cutting concerns:
```go
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    
    log.Printf("gRPC call: %s started", info.FullMethod)
    
    resp, err := handler(ctx, req)
    
    duration := time.Since(start)
    log.Printf("gRPC call: %s completed in %v", info.FullMethod, duration)
    
    return resp, err
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // Skip auth for certain methods
    if info.FullMethod == "/health/check" {
        return handler(ctx, req)
    }
    
    // Extract auth metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
    }
    
    authHeaders := md.Get("authorization")
    if len(authHeaders) == 0 {
        return nil, status.Errorf(codes.Unauthenticated, "authorization header not provided")
    }
    
    // Validate token (simplified)
    if authHeaders[0] != "Bearer valid-token" {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token")
    }
    
    return handler(ctx, req)
}
```

## Hint 5: gRPC Client Implementation
Create clients that connect to the services:
```go
type UserServiceClient struct {
    conn   *grpc.ClientConn
    client UserServiceClient // from generated proto
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID string) (*User, error) {
    // Add auth metadata
    ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer valid-token")
    
    req := &GetUserRequest{UserId: userID}
    resp, err := c.client.GetUser(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return resp.User, nil
}

func ConnectToServices(userServiceAddr, productServiceAddr string) (*UserServiceClient, *ProductServiceClient, error) {
    // Connect to User Service
    userConn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect to user service: %w", err)
    }
    
    userClient := &UserServiceClient{
        conn:   userConn,
        client: NewUserServiceClient(userConn),
    }
    
    // Similar for product service...
    
    return userClient, productClient, nil
}
```

## Hint 6: Service Orchestration
Implement the order service that coordinates multiple services:
```go
type OrderService struct {
    userClient    *UserServiceClient
    productClient *ProductServiceClient
}

func (s *OrderService) CreateOrder(ctx context.Context, userID, productID string, quantity int) (*Order, error) {
    // Step 1: Validate user
    _, err := s.userClient.ValidateUser(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("user validation failed: %w", err)
    }
    
    // Step 2: Check inventory
    _, err = s.productClient.CheckInventory(ctx, productID, quantity)
    if err != nil {
        return nil, fmt.Errorf("inventory check failed: %w", err)
    }
    
    // Step 3: Create order
    order := &Order{
        Id:        generateOrderID(),
        UserId:    userID,
        ProductId: productID,
        Quantity:  quantity,
        Status:    "created",
        CreatedAt: time.Now().Unix(),
    }
    
    return order, nil
}
```

## Hint 7: Error Handling Across Services
Handle both gRPC errors and business logic errors:
```go
func handleServiceError(err error) error {
    if err == nil {
        return nil
    }
    
    // Check if it's a gRPC status error
    if status, ok := status.FromError(err); ok {
        switch status.Code() {
        case codes.NotFound:
            return fmt.Errorf("resource not found: %s", status.Message())
        case codes.PermissionDenied:
            return fmt.Errorf("access denied: %s", status.Message())
        case codes.ResourceExhausted:
            return fmt.Errorf("resource exhausted: %s", status.Message())
        default:
            return fmt.Errorf("service error: %s", status.Message())
        }
    }
    
    // Handle other errors
    return fmt.Errorf("unexpected error: %w", err)
}
```

## Key gRPC Concepts:
- **Status Codes**: Use appropriate codes (NotFound, PermissionDenied, etc.)
- **Interceptors**: Chain interceptors for cross-cutting concerns
- **Metadata**: Pass authentication and context information
- **Error Handling**: Distinguish between transport and business logic errors
- **Service Discovery**: Connect clients to running services
- **Context Propagation**: Pass context across service boundaries 