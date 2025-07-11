# Hints for RESTful Book Management API

## Hint 1: HTTP Handler Structure
Use `http.HandlerFunc` for each endpoint and route them with a multiplexer:
```go
func (h *BookHandler) SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/books", h.handleBooks)
    mux.HandleFunc("/api/books/", h.handleBookByID)
    mux.HandleFunc("/api/books/search", h.handleSearch)
    return mux
}
```

## Hint 2: Method-based Routing
Handle different HTTP methods in your handler:
```go
func (h *BookHandler) handleBooks(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.getAllBooks(w, r)
    case http.MethodPost:
        h.createBook(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
```

## Hint 3: JSON Response Helper
Create a helper function for JSON responses:
```go
func writeJSONResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

## Hint 4: Request Body Parsing
Parse JSON request bodies:
```go
func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    // Validate and create book
}
```

## Hint 5: In-Memory Repository
Implement the repository with a map:
```go
type InMemoryBookRepository struct {
    books map[string]*Book
    mutex sync.RWMutex
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    book, exists := r.books[id]
    if !exists {
        return nil, errors.New("book not found")
    }
    return book, nil
}
```

## Hint 6: URL Parameter Extraction
Extract ID from URL path:
```go
func extractIDFromPath(path string) string {
    parts := strings.Split(path, "/")
    if len(parts) >= 4 {
        return parts[3] // /api/books/{id}
    }
    return ""
}
```

## Hint 7: Query Parameter Handling
Handle search parameters:
```go
func (h *BookHandler) handleSearch(w http.ResponseWriter, r *http.Request) {
    author := r.URL.Query().Get("author")
    title := r.URL.Query().Get("title")
    
    if author != "" {
        books, err := h.Service.SearchBooksByAuthor(author)
        // handle result
    } else if title != "" {
        books, err := h.Service.SearchBooksByTitle(title)
        // handle result
    }
}
```

## Hint 8: Input Validation
Validate required fields:
```go
func validateBook(book *Book) error {
    if book.Title == "" {
        return errors.New("title is required")
    }
    if book.Author == "" {
        return errors.New("author is required")
    }
    if book.PublishedYear <= 0 {
        return errors.New("published year must be positive")
    }
    return nil
}
```

## Hint 9: Search Implementation
Implement case-insensitive search:
```go
func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    var results []*Book
    lowerAuthor := strings.ToLower(author)
    
    for _, book := range r.books {
        if strings.Contains(strings.ToLower(book.Author), lowerAuthor) {
            results = append(results, book)
        }
    }
    return results, nil
}
```

## Hint 10: Error Response Structure
Create consistent error responses:
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

func writeErrorResponse(w http.ResponseWriter, message string, status int) {
    response := ErrorResponse{
        Error:   http.StatusText(status),
        Message: message,
    }
    writeJSONResponse(w, response, status)
} 