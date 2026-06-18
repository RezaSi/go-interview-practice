// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Book represents a book in the database
type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
	ISBN          string `json:"isbn"`
	Description   string `json:"description"`
}

// BookRepository defines the operations for book data access
type BookRepository interface {
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Create(book *Book) error
	Update(id string, book *Book) error
	Delete(id string) error
	SearchByAuthor(author string) ([]*Book, error)
	SearchByTitle(title string) ([]*Book, error)
}

// InMemoryBookRepository implements BookRepository using in-memory storage
type InMemoryBookRepository struct {
	books map[string]*Book
	mu    sync.RWMutex
}

// NewInMemoryBookRepository creates a new in-memory book repository
func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books: make(map[string]*Book),
	}
}

// Implement BookRepository methods for InMemoryBookRepository
// ...

func (r *InMemoryBookRepository) GetAll() ([]*Book, error) {
	bookSlice := []*Book{}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, book := range r.books {
		bookSlice = append(bookSlice, book)
	}

	return bookSlice, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if b, exists := r.books[id]; !exists {
		return nil, errors.New("Book not found")
	} else {
		return b, nil
	}
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if book.ISBN == "" || book.Author == "" || book.Title == "" || book.Description == "" {
		return errors.New("Missing Information")
	}

	book.ID = book.ISBN
	if _, exists := r.books[book.ID]; exists {
		return errors.New("Book ID already exists")
	} else {
		r.books[book.ID] = book
		return nil
	}
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return errors.New("Book not found")
	} else {
		r.books[id] = book
		return nil
	}
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return errors.New("Book not found")
	} else {
		delete(r.books, id)
		return nil
	}
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bookSlice := []*Book{}

	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Author), strings.ToLower(author)) {
			bookSlice = append(bookSlice, book)
		}
	}
	return bookSlice, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bookSlice := []*Book{}

	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(title)) {
			bookSlice = append(bookSlice, book)
		}
	}

	return bookSlice, nil
}

// BookService defines the business logic for book operations
type BookService interface {
	GetAllBooks() ([]*Book, error)
	GetBookByID(id string) (*Book, error)
	CreateBook(book *Book) error
	UpdateBook(id string, book *Book) error
	DeleteBook(id string) error
	SearchBooksByAuthor(author string) ([]*Book, error)
	SearchBooksByTitle(title string) ([]*Book, error)
}

// DefaultBookService implements BookService
type DefaultBookService struct {
	repo BookRepository
}

// NewBookService creates a new book service
func NewBookService(repo BookRepository) *DefaultBookService {
	return &DefaultBookService{
		repo: repo,
	}
}

// Implement BookService methods for DefaultBookService
// ...

func (dbs *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return dbs.repo.GetAll()
}

func (dbs *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return dbs.repo.GetByID(id)
}

func (dbs *DefaultBookService) CreateBook(book *Book) error {
	return dbs.repo.Create(book)
}

func (dbs *DefaultBookService) UpdateBook(id string, book *Book) error {
	return dbs.repo.Update(id, book)
}

func (dbs *DefaultBookService) DeleteBook(id string) error {
	return dbs.repo.Delete(id)
}

func (dbs *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	return dbs.repo.SearchByAuthor(author)
}

func (dbs *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	return dbs.repo.SearchByTitle(title)
}

// BookHandler handles HTTP requests for book operations
type BookHandler struct {
	Service BookService
}

// NewBookHandler creates a new book handler
func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{
		Service: service,
	}
}

func writeJSONReponse(w http.ResponseWriter, data any, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSONReponse(w, ErrorResponse{
		StatusCode: status,
		Error:      msg,
	}, status)
}

// HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	switch {
	case strings.HasPrefix(path, "/api/books/search") && method == http.MethodGet:
		if author := r.URL.Query().Get("author"); author != "" {
			res, err := h.Service.SearchBooksByAuthor(author)
			if err != nil {
				writeError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSONReponse(w, res, http.StatusOK)
			return
		}
		if title := r.URL.Query().Get("title"); title != "" {
			res, err := h.Service.SearchBooksByTitle(title)
			if err != nil {
				writeError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSONReponse(w, res, http.StatusOK)
			return
		}
		writeError(w, "missing search parameter", http.StatusBadRequest)
		return

	case path == "/api/books" && method == http.MethodGet:
		books, err := h.Service.GetAllBooks()
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSONReponse(w, books, http.StatusOK)
		return

	case strings.HasPrefix(path, "/api/books/") && method == http.MethodGet:
		book, err := h.Service.GetBookByID(strings.TrimPrefix(path, "/api/books/"))
		if err != nil {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSONReponse(w, book, http.StatusOK)
		return

	case path == "/api/books" && method == http.MethodPost:
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.Service.CreateBook(&book); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSONReponse(w, book, http.StatusCreated)
		return

	case strings.HasPrefix(path, "/api/books/") && method == http.MethodPut:
		id := strings.TrimPrefix(path, "/api/books/")
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.Service.UpdateBook(id, &book); err != nil {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSONReponse(w, book, http.StatusOK)
		return

	case strings.HasPrefix(path, "/api/books/") && method == http.MethodDelete:
		id := strings.TrimPrefix(path, "/api/books/")
		if err := h.Service.DeleteBook(id); err != nil {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSONReponse(w, "", http.StatusOK)
		return

	default:
		writeError(w, "invalid endpoint", http.StatusNotFound)
		return
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helper functions
// ...

func main() {
	// Initialize the repository, service, and handler
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	// Create a new router and register endpoints
	http.HandleFunc("/api/books", handler.HandleBooks)
	http.HandleFunc("/api/books/", handler.HandleBooks)

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

