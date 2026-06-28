// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
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

func (r *InMemoryBookRepository) GetAll() ([]*Book, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    return slices.Collect(maps.Values(r.books)), nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    book, exists := r.books[id]
    if !exists {
        return nil, errors.New(fmt.Sprintf("Error, cannot find book with ID %v", id))
    }
    return book, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.books[book.ID]; exists {
        return errors.New(fmt.Sprintf("Error, cannot create book with ID %v since it already exists.", book.ID))
    }
    
	r.books[book.ID] = book
	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.books[id]; !exists {
        return errors.New(fmt.Sprintf("Error, cannot edit book with ID %v since it does not exist.", id))
    }
    
	r.books[id] = book
	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.books[id]; !exists {
        return errors.New(fmt.Sprintf("Error, cannot delete book with ID %v since it does not exist.", id))
    }
    
	delete(r.books, id)
	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    books := slices.Collect(maps.Values(r.books))
    filteredBooks := []*Book{}
    
    for _, b := range(books) {
        if (strings.Contains(b.Author, author)) {
            filteredBooks = append(filteredBooks, b)
        }
    }
    
    return filteredBooks, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    books := slices.Collect(maps.Values(r.books))
    filteredBooks := []*Book{}
    
    for _, b := range(books) {
        if (strings.Contains(b.Title, title)) {
            filteredBooks = append(filteredBooks, b)
        }
    }
    
    return filteredBooks, nil
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

func (s *DefaultBookService) GetAllBooks() ([]*Book, error) {
    return s.repo.GetAll()
}

func (s *DefaultBookService) GetBookByID(id string) (*Book, error) {
    return s.repo.GetByID(id)
}

func (s *DefaultBookService) CreateBook(book *Book) error {
    missingFields := []string{}
	if book.ID == "" {
	    missingFields = append(missingFields, "id")
	} 
	if book.Title == ""{
	    missingFields = append(missingFields, "title")
	}
	if book.Author == ""{
	    missingFields = append(missingFields, "author")
	}
	if book.PublishedYear == 0 {
	    missingFields = append(missingFields, "published_year")
	}
	if book.ISBN == ""{
	    missingFields = append(missingFields, "isbn")
	}
	if book.Description == ""{
	    missingFields = append(missingFields, "description")
	}
	if len(missingFields) != 0 {
        return errors.New(fmt.Sprintf("Error, field(s) %v are invalid.", missingFields))
	}

    return s.repo.Create(book)
}

func (s *DefaultBookService) UpdateBook(id string, book *Book) error {
    return s.repo.Update(id, book)
}

func (s *DefaultBookService) DeleteBook(id string) error {
    return s.repo.Delete(id)
}

func (s *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
    return s.repo.SearchByAuthor(author)
}

func (s *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
    return s.repo.SearchByTitle(title)
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

func (h *BookHandler) handleGetAllBooks (w http.ResponseWriter, r *http.Request) {
    books, err := h.Service.GetAllBooks()
    if err != nil {
        http.Error(w, "Cannot get all books", http.StatusInternalServerError)
        return
    }
    
    writeJSONResponse(w, books, http.StatusOK)
}

func (h *BookHandler) handleGetBookById (w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    book, err := h.Service.GetBookByID(id)
    if err != nil {
        http.Error(w, fmt.Sprintf("Cannot get book by ID. %v", err), http.StatusNotFound)
        return
    }
    
    writeJSONResponse(w, book, http.StatusOK)
}

func (h *BookHandler) handleCreateBook (w http.ResponseWriter, r *http.Request) {
    book := extractBookFromResponseBody(w, r)
    book.ID = uuid.New().String()
    
    if err := h.Service.CreateBook(book); err != nil {
        http.Error(w, fmt.Sprintf("Could not create book request payload. %v", err), http.StatusBadRequest)
        return
    }
    writeJSONResponse(w, book, http.StatusCreated)
}

func (h *BookHandler) handleUpdateBook (w http.ResponseWriter, r *http.Request) {
    book := extractBookFromResponseBody(w, r)
    
    if err := h.Service.UpdateBook(book.ID, book); err != nil {
        http.Error(w, fmt.Sprintf("Cannot update book with ID %v. %v", book.ID, err), http.StatusNotFound)
        return
    }
    
    writeJSONResponse(w, book, http.StatusOK)
}

func (h *BookHandler) handleDeleteBook (w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if err := h.Service.DeleteBook(id); err != nil {
        http.Error(w, fmt.Sprintf("Cannot delete book with ID %v. %v", id, err), http.StatusNotFound)
        return
    }
    
    writeJSONResponse(w, nil, http.StatusOK)
}

func (h *BookHandler) handleSearchBooks (w http.ResponseWriter, r *http.Request) {
    author := r.URL.Query().Get("author")
    if author != "" {
        books, err := h.Service.SearchBooksByAuthor(author)
        if err != nil {
            http.Error(w, fmt.Sprintf("Cannot search books by author %v. %v", author, err), http.StatusNotFound)
            return
        }
        writeJSONResponse(w, books, http.StatusOK)
        return
    }
    
    title := r.URL.Query().Get("title")
    if title != "" {
        books, err := h.Service.SearchBooksByTitle(title)
        if err != nil {
            http.Error(w, fmt.Sprintf("Cannot search books by title %v. %v", title, err), http.StatusNotFound)
            return
        }
        
        writeJSONResponse(w, books, http.StatusOK)
        return
    }
}

// HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /api/books", h.handleGetAllBooks)
	mux.HandleFunc("GET /api/books/{id}", h.handleGetBookById)
	mux.HandleFunc("POST /api/books", h.handleCreateBook)
	mux.HandleFunc("PUT /api/books/{id}", h.handleUpdateBook)
	mux.HandleFunc("DELETE /api/books/{id}", h.handleDeleteBook)
	mux.HandleFunc("GET /api/books/search", h.handleSearchBooks)

    mux.ServeHTTP(w, r)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Printf("JSON encoding failed: %v", err) 
    }
}

func extractBookFromResponseBody(w http.ResponseWriter, r *http.Request) *Book {
    var book Book
    
    decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()
    if err := decoder.Decode(&book); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return nil
    }
    return &book
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