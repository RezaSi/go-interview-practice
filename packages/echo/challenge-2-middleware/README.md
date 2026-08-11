# Challenge 2: Middleware & Request/Response Handling

Build an **Enhanced Blog API** using Echo that demonstrates advanced middleware patterns.

## Challenge Requirements

You need to implement the following middleware:

1. **Custom Logging Middleware** - Log all requests with timing and request IDs
2. **Authentication Middleware** - Protect certain routes with API keys  
3. **CORS Middleware** - Handle cross-origin requests properly
4. **Rate Limiting Middleware** - Limit requests per IP (100 per minute)
5. **Request ID Middleware** - Add unique request IDs to each request
6. **Error Handling Middleware** - Centralized error management with consistent responses

## API Endpoints

### Public Endpoints
- `GET /health` - Health check

### Protected Endpoints (Require API Key)
- `GET /posts` - Get all blog posts (paginated)
- `POST /posts` - Create a new blog post
- `GET /posts/:id` - Get specific post by ID
- `PUT /posts/:id` - Update existing post
- `DELETE /posts/:id` - Delete post

## Data Structure

```go
type Post struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    Author    string `json:"author"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type PostResponse struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
    Post    *Post  `json:"post,omitempty"`
}

type PostListResponse struct {
    Status string `json:"status"`
    Count  int    `json:"count"`
    Posts  []Post `json:"posts"`
}

type ErrorResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

type HealthResponse struct {
    Status    string `json:"status"`
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
    RequestID string `json:"request_id"`
}
```

## Request/Response Examples

**GET /health**
```json
{
    "status": "healthy",
    "message": "Blog API is running",
    "timestamp": "2024-01-15T10:30:00Z",
    "request_id": "req-123456"
}
```

**POST /posts** (Request body)
```json
{
    "title": "Getting Started with Echo",
    "content": "Echo is a high-performance web framework...",
    "author": "John Doe"
}
```

## Testing Requirements

Your solution must pass tests for:
- Health check works without authentication
- All protected endpoints require valid API key
- Request ID middleware adds unique IDs to all requests
- Logging middleware captures request details and timing
- CORS headers are properly set
- Rate limiting blocks excessive requests
- Error handling middleware provides consistent error responses
- All CRUD operations work correctly for blog posts
- Proper HTTP status codes for all scenarios