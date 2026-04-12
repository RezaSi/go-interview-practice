// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid format")
)

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
func (b *InMemoryBookRepository) GetAll() ([]*Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	books := []*Book{}

	for _, b := range b.books {
		books = append(books, b)
	}
	return books, nil
}

func (b *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	book, exist := b.books[id]

	if !exist {
		return nil, ErrNotFound
	}

	return book, nil
}

func (b *InMemoryBookRepository) Create(book *Book) error {

	b.mu.Lock()
	defer b.mu.Unlock()

	b.books[book.ID] = book

	return nil
}

func (b *InMemoryBookRepository) Update(id string, book *Book) error {

	b.mu.Lock()
	defer b.mu.Unlock()

	_, exist := b.books[id]

	if !exist {
		return ErrNotFound
	}
	b.books[id] = book
	return nil
}

func (b *InMemoryBookRepository) Delete(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, exist := b.books[id]

	if !exist {
		return ErrNotFound
	}
	delete(b.books, id)
	return nil
}

func (b *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	books := make([]*Book, 0)
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, book := range b.books {
		if strings.Contains(strings.ToLower(book.Author), strings.ToLower(author)) {
			books = append(books, book)
		}
	}

	return books, nil
}

func (b *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	books := make([]*Book, 0)
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, book := range b.books {
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
func (bs *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return bs.repo.GetAll()
}

func (bs *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return bs.repo.GetByID(id)
}

func (bs *DefaultBookService) CreateBook(book *Book) error {
	err := bs.validateBook(book)

	if err != nil {
		return ErrInvalid
	}

	book.ID = uuid.New().String()

	return bs.repo.Create(book)
}

func (bs *DefaultBookService) UpdateBook(id string, book *Book) error {
	return bs.repo.Update(id, book)
}

func (bs *DefaultBookService) DeleteBook(id string) error {
	return bs.repo.Delete(id)
}

func (bs *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	if author == "" {
		return nil, ErrInvalid
	}
	return bs.repo.SearchByAuthor(author)
}

func (bs *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	if title == "" {
		return nil, ErrInvalid
	}
	return bs.repo.SearchByTitle(title)
}

func (bs *DefaultBookService) validateBook(book *Book) error {
	if book == nil {
		return errors.New("book can't be empty")
	}

	if book.Title == "" {
		return errors.New("title can't be empty")
	}

	if book.Author == "" {
		return errors.New("author can't be empty")
	}

	if book.PublishedYear < 0 {
		return errors.New("invalid published year")
	}

	if book.ISBN == "" {
		return errors.New("isbn can't be empty")
	}

	if book.Description == "" {
		return errors.New("description can't be empty")
	}
	return nil
}

// BookHandler handles HTTP requests for book operations
type BookHandler struct {
	Service BookService
}

// GET /api/books: Get all books
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /api/books/{id}: Get a specific book by ID
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/books/")

	book, err := h.Service.GetBookByID(id)

	if err == ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// POST /api/books: Create a new book
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.Service.CreateBook(&book)
	if err == ErrInvalid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// PUT /api/books/{id}: Update an existing book
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	id := strings.TrimPrefix(r.URL.Path, "/api/books/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "invalid id",
		})
		return
	}
	err := json.NewDecoder(r.Body).Decode(&book)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.Service.UpdateBook(id, &book)
	if err == ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

}

// DELETE /api/books/{id}: Delete a book
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/books/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "invalid id",
		})
		return
	}
	err := h.Service.DeleteBook(id)
	if err == ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GET /api/books/search?author={author}: Search books by author
// GET /api/books/search?title={title}: Search books by title
func (h *BookHandler) SearchBook(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var books []*Book
	var err error

	switch {
	case query.Get("author") != "":
		books, err = h.Service.SearchBooksByAuthor(query.Get("author"))
	case query.Get("title") != "":
		books, err = h.Service.SearchBooksByTitle(query.Get("title"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "invalid query params",
		})
		return
	}
	if err == ErrInvalid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

// NewBookHandler creates a new book handler
func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{
		Service: service,
	}
}

// // HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this method to handle all book endpoints
	// Use the path and method to determine the appropriate action
	// Call the service methods accordingly
	// Return appropriate status codes and JSON responses
	w.Header().Set("Content-Type", "application/json")
	path, method := r.URL.Path, r.Method
	switch {
	case strings.HasPrefix(path, "/api/books/search") && method == http.MethodGet:
		h.SearchBook(w, r)
	case path == "/api/books" && method == http.MethodGet:
		h.GetAllBooks(w, r)
	case path == "/api/books" && method == http.MethodPost:
		h.CreateBook(w, r)
	case strings.HasPrefix(path, "/api/books/") && method == http.MethodGet:
		h.GetBookByID(w, r)
	case strings.HasPrefix(path, "/api/books/") && method == http.MethodPut:
		h.UpdateBook(w, r)
	case strings.HasPrefix(path, "/api/books/") && method == http.MethodDelete:
		h.DeleteBook(w, r)
	}

}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helper functions
// func

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
