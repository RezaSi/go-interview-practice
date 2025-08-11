# Challenge 2: Middleware & Request/Response Handling

Build an **Enhanced Blog API** using Fiber that demonstrates advanced middleware patterns.

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
- `GET /ping` - Health check
- `GET /articles` - Get all articles (paginated)
- `GET /articles/:id` - Get article by ID

### Protected Endpoints (Require API Key: `X-API-Key` header)
- `POST /articles` - Create new article  
- `PUT /articles/:id` - Update article
- `DELETE /articles/:id` - Delete article
- `GET /admin/stats` - Get API usage statistics

**Valid API Keys:** `admin-key-123` (full access), `user-key-456` (read-only protected routes)

## Data Structures

```go
type Article struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Message   string      `json:"message,omitempty"`
    Error     string      `json:"error,omitempty"`
    RequestID string      `json:"request_id,omitempty"`
}
```

## Testing Requirements

Your solution must handle:
- Proper middleware execution order
- Request ID generation and propagation
- Rate limiting enforcement (reject after 100 requests/minute)
- CORS headers for cross-origin requests
- Authentication validation for protected routes
- Centralized error handling with consistent response format
- Request logging with timing information