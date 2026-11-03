// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json/v2"
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

func (r *InMemoryBookRepository) GetAll() ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	books := slices.Collect(maps.Values(r.books))
	return books, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, ok := r.books[id]
	if !ok {
		return nil, fmt.Errorf("There is no book with ID '%s'", id)
	}
	return book, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	if err := validate(book); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	uuid := uuid.NewString()
	book.ID = uuid
	r.books[uuid] = book
	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	if err := validate(book); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.books[id]
	if !ok {
		return fmt.Errorf("There is no book with ID '%s'", id)
	}
	book.ID = existing.ID
	r.books[id] = book
	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.books[id]
	if !ok {
		return fmt.Errorf("There is no book with ID '%s'", id)
	}
	delete(r.books, id)
	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []*Book
	authorLower := strings.ToLower(author)
	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Author), authorLower) {
			results = append(results, book)
		}
	}
	return results, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []*Book
	titleLower := strings.ToLower(title)
	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Title), titleLower) {
			results = append(results, book)
		}
	}
	return results, nil
}

func validate(book *Book) error {
	if book.Title == "" {
		return fmt.Errorf("Invalid title: '%s'", book.Title)
	}
	if book.Author == "" {
		return fmt.Errorf("Invalid author: '%s'", book.Author)
	}
	if book.ISBN == "" {
		return fmt.Errorf("Invalid ISBN: '%s'", book.ISBN)
	}
	return nil
}

// NewInMemoryBookRepository creates a new in-memory book repository
func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books: make(map[string]*Book),
	}
}

// Implement BookRepository methods for InMemoryBookRepository
// ...

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

func (s *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return s.repo.GetAll()
}

func (s *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return s.repo.GetByID(id)
}

func (s *DefaultBookService) CreateBook(book *Book) error {
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

// NewBookService creates a new book service
func NewBookService(repo BookRepository) *DefaultBookService {
	return &DefaultBookService{
		repo: repo,
	}
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
	w.Header().Set("Content-Type", "application/json")
	// Since `r.PathValue` does not work with the tests I'll use a more hacky approach
	pathSplit := strings.Split(r.URL.Path, "/")
	lastFrag := pathSplit[len(pathSplit)-1]
	switch method := r.Method; method {
	case "GET":
		if len(pathSplit) == 4 && lastFrag == "search" {
			author := r.URL.Query().Get("author")
			title := r.URL.Query().Get("title")
			if author != "" {
				books, err := h.Service.SearchBooksByAuthor(author)
				if err != nil {
					http.Error(w, "Not found", http.StatusNotFound)
					return
				}
				json.MarshalWrite(w, books)
				return
			} else if title != "" {
				books, err := h.Service.SearchBooksByTitle(title)
				if err != nil {
					http.Error(w, "Not found", http.StatusNotFound)
					return
				}
				json.MarshalWrite(w, books)
				return
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
		}

		if len(pathSplit) == 4 && lastFrag != "" {
			book, err := h.Service.GetBookByID(lastFrag)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			json.MarshalWrite(w, book)
			return
		}

		books, err := h.Service.GetAllBooks()
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		json.MarshalWrite(w, books)
		return
	case "POST":
		var book Book
		defer r.Body.Close()
		err := json.UnmarshalRead(r.Body, &book)
		if err != nil {
			http.Error(w, "Malformed body", http.StatusBadRequest)
			return
		}
		err = h.Service.CreateBook(&book)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.MarshalWrite(w, book)
		return
	case "PUT":
		if len(pathSplit) == 4 && lastFrag != "" {
			var book Book
			defer r.Body.Close()
			err := json.UnmarshalRead(r.Body, &book)
			if err != nil {
				http.Error(w, "Malformed body", http.StatusBadRequest)
				return
			}
			err = h.Service.UpdateBook(lastFrag, &book)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			json.MarshalWrite(w, book)
			return
		}
	case "DELETE":
		if len(pathSplit) == 4 && lastFrag != "" {
			err := h.Service.DeleteBook(lastFrag)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
		}
	default:
		http.Error(w, "Invalid method: "+method, http.StatusMethodNotAllowed)
		w.Header().Del("Content-Type")
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
	// HandleBooks is given by the exercise, but I would prefer separate functions
	http.HandleFunc("GET /api/books", handler.HandleBooks)
	http.HandleFunc("POST /api/books", handler.HandleBooks)
	http.HandleFunc("PUT /api/books/{id}", handler.HandleBooks)
	http.HandleFunc("DELETE /api/books/{id}", handler.HandleBooks)

	// Start the server
	addr := ":8080"
	log.Println("Server starting on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
