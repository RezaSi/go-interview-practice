// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
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
	r.mu.RLock()
	defer r.mu.RUnlock()
	books := make([]*Book, 0, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}
	return books, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	book, ok := r.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[book.ID]; ok {
		return errors.New("book already exists")
	}
	id, _ := uuid.NewUUID()
	book.ID = id.String()
	r.books[id.String()] = book
	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; ok {
		r.books[id] = book
		return nil
	}
	return errors.New("book already exists")
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; ok {
		delete(r.books, id)
		return nil
	}
	return errors.New("book not found")
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	books := make([]*Book, 0, len(r.books))
	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Author), strings.ToLower(author)) {
			books = append(books, book)
		}
	}
	return books, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	books := make([]*Book, 0, len(r.books))
	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(title)) {
			books = append(books, book)
		}
	}
	return books, nil
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
func (d DefaultBookService) GetAllBooks() ([]*Book, error) {
	resp, err := d.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d DefaultBookService) GetBookByID(id string) (*Book, error) {
	resp, err := d.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d DefaultBookService) CreateBook(book *Book) error {
	err := d.repo.Create(book)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultBookService) UpdateBook(id string, book *Book) error {
	err := d.repo.Update(id, book)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultBookService) DeleteBook(id string) error {
	err := d.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	resp, err := d.repo.SearchByAuthor(author)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	resp, err := d.repo.SearchByTitle(title)
	if err != nil {
		return nil, err
	}
	return resp, nil
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

// HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this method to handle all book endpoints
	// Use the path and method to determine the appropriate action
	// Call the service methods accordingly
	// Return appropriate status codes and JSON responses
	switch r.Method {
	case http.MethodGet:
		author := r.URL.Query().Get("author")
		title := r.URL.Query().Get("title")
		if author != "" {
			books, err := h.Service.SearchBooksByAuthor(author)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
			writeJSONResponse(w, books, http.StatusOK)
		} else if title != "" {
			books, err := h.Service.SearchBooksByTitle(title)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
			writeJSONResponse(w, books, http.StatusOK)
		}
		id := extractIDFromPath(r.URL.Path)
		if id != "" {
			book, err := h.Service.GetBookByID(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
			writeJSONResponse(w, book, http.StatusOK)
		}

		books, err := h.Service.GetAllBooks()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		writeJSONResponse(w, books, http.StatusOK)
	case http.MethodPost:
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		err := validateBookRequest(book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err = h.Service.CreateBook(&book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		writeJSONResponse(w, book, http.StatusCreated)
	case http.MethodPut:
		id := extractIDFromPath(r.URL.Path)
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err := validateBookRequest(book)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err = h.Service.UpdateBook(id, &book)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		writeJSONResponse(w, book, http.StatusOK)
	case http.MethodDelete:
		id := extractIDFromPath(r.URL.Path)
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		err := h.Service.DeleteBook(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		writeJSONResponse(w, id, http.StatusOK)
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helper functions
// ...
func writeJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func extractIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3] // /api/books/{id}
	}
	return ""
}

func validateBookRequest(book Book) error {
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
