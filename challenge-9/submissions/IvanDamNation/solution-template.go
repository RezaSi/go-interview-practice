// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	maxBodySize    = 1 * 1024 * 1024
	minPublishYear = 848
)

var (
	ErrBookNotFound    = errors.New("book not found")
	ErrIDAlreadyExists = errors.New("ID is not unique")
	ErrInvalidInput    = errors.New("invalid input")
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
func (r *InMemoryBookRepository) GetAll() ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]*Book, 0, len(r.books))
	for _, b := range r.books {
		copyB := *b
		books = append(books, &copyB)
	}

	return books, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.books[id]
	if !exists {
		return nil, ErrBookNotFound
	}

	copyB := *book
	return &copyB, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[book.ID]; exists {
		return ErrIDAlreadyExists
	}

	copyB := *book
	r.books[book.ID] = &copyB

	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return ErrBookNotFound
	}

	book.ID = id

	copyB := *book
	r.books[id] = &copyB

	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return ErrBookNotFound
	}

	delete(r.books, id)
	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, b := range r.books {
		if strings.Contains(strings.ToLower(b.Author), strings.ToLower(author)) {
			count++
		}
	}

	srchRes := make([]*Book, 0, count)
	for _, b := range r.books {
		if strings.Contains(strings.ToLower(b.Author), strings.ToLower(author)) {
			copyB := *b
			srchRes = append(srchRes, &copyB)
		}
	}

	return srchRes, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, b := range r.books {
		if strings.Contains(strings.ToLower(b.Title), strings.ToLower(title)) {
			count++
		}
	}

	srchRes := make([]*Book, 0, count)
	for _, b := range r.books {
		if strings.Contains(strings.ToLower(b.Title), strings.ToLower(title)) {
			copyB := *b
			srchRes = append(srchRes, &copyB)
		}
	}

	return srchRes, nil
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
func (s *DefaultBookService) GetAllBooks() ([]*Book, error) {
	books, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *DefaultBookService) GetBookByID(id string) (*Book, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *DefaultBookService) CreateBook(book *Book) error {
	if book.Title == "" || book.Author == "" || book.ISBN == "" || book.PublishedYear < minPublishYear || book.PublishedYear > time.Now().Year() {
		return ErrInvalidInput
	}

	const maxRetries = 3
	var err error

	for range maxRetries {
		book.ID = uuid.New().String()
		err = s.repo.Create(book)
		if errors.Is(err, ErrIDAlreadyExists) {
			continue
		}
		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultBookService) UpdateBook(id string, book *Book) error {
	if book.Title == "" || book.Author == "" || book.ISBN == "" || book.PublishedYear < minPublishYear || book.PublishedYear > time.Now().Year() {
		return ErrInvalidInput
	}

	err := s.repo.Update(id, book)
	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultBookService) DeleteBook(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	books, err := s.repo.SearchByAuthor(author)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	books, err := s.repo.SearchByTitle(title)
	if err != nil {
		return nil, err
	}

	return books, nil
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
	path := strings.TrimPrefix(r.URL.Path, "/api/books")
	path = strings.Trim(path, "/")

	var segments []string
	if path != "" {
		segments = strings.Split(path, "/")
	}

	// /api/books
	if len(segments) == 0 {
		switch r.Method {
		case http.MethodGet:
			books, err := h.Service.GetAllBooks()
			if err != nil {
				respondWithError(w, ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
				return
			}
			respondWithJSON(w, http.StatusOK, books)
			return

		case http.MethodPost:
			var book Book
			r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
			defer r.Body.Close()

			if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
				respondWithError(w, ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Error:      err.Error(),
				})
				return
			}

			err := h.Service.CreateBook(&book)
			if err != nil {
				if errors.Is(err, ErrInvalidInput) {
					respondWithError(w, ErrorResponse{
						StatusCode: http.StatusBadRequest,
						Error:      err.Error(),
					})
					return
				}
				respondWithError(w, ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
				return
			}

			respondWithJSON(w, http.StatusCreated, book)
			return

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}

	// /api/books/search
	if len(segments) == 1 && segments[0] == "search" {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		var books []*Book
		var err error

		if author := query.Get("author"); author != "" {
			books, err = h.Service.SearchBooksByAuthor(author)
		} else if title := query.Get("title"); title != "" {
			books, err = h.Service.SearchBooksByTitle(title)
		} else {
			respondWithError(w, ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "Missing 'author' or 'title' query parameter",
			})
			return
		}

		if err != nil {
			respondWithError(w, ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
			return
		}

		respondWithJSON(w, http.StatusOK, books)
		return
	}

	// /api/books/{id}
	if len(segments) == 1 {
		id := segments[0]

		switch r.Method {
		case http.MethodGet:
			book, err := h.Service.GetBookByID(id)
			if err != nil {
				if errors.Is(err, ErrBookNotFound) {
					respondWithError(w, ErrorResponse{
						StatusCode: http.StatusNotFound,
						Error:      err.Error(),
					})
					return
				}
				respondWithError(w, ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
				return
			}
			respondWithJSON(w, http.StatusOK, book)
			return

		case http.MethodPut:
			var book Book
			r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
			defer r.Body.Close()

			if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
				respondWithError(w, ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Error:      "Invalid JSON payload",
				})
				return
			}

			err := h.Service.UpdateBook(id, &book)
			if err != nil {
				if errors.Is(err, ErrInvalidInput) {
					respondWithError(w, ErrorResponse{
						StatusCode: http.StatusBadRequest,
						Error:      err.Error(),
					})
					return
				}
				if errors.Is(err, ErrBookNotFound) {
					respondWithError(w, ErrorResponse{
						StatusCode: http.StatusNotFound,
						Error:      err.Error(),
					})
					return
				}
				respondWithError(w, ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
				return
			}

			respondWithJSON(w, http.StatusOK, &book)
			return

		case http.MethodDelete:
			if err := h.Service.DeleteBook(id); err != nil {
				if errors.Is(err, ErrBookNotFound) {
					respondWithError(w, ErrorResponse{
						StatusCode: http.StatusNotFound,
						Error:      err.Error(),
					})
					return
				}
				respondWithError(w, ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
				return
			}
			respondWithJSON(w, http.StatusOK, nil)
			return

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError)
		}
	}
}

func respondWithError(w http.ResponseWriter, errResp ErrorResponse) {
	if errResp.StatusCode == 0 {
		errResp.StatusCode = http.StatusInternalServerError
	}
	respondWithJSON(w, errResp.StatusCode, errResp)
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
