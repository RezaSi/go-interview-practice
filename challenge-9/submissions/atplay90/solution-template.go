// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"errors"
	"log"
	//"maps"
	"net/http"
	//"slices"
	"sync"
	"strings"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
	out := make([]*Book, 0, len(r.books))
	for _, b := range r.books {
		copied := *b
		out = append(out, &copied)
	}
	return out, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, exists := r.books[id]
	if !exists {
		return nil, errors.New("Item not found")
	}
	copied := *value
	return &copied, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	if book == nil {
		return errors.New("book is nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.books[book.ID]
	if exists {
		return errors.New("Item already exists")
	}
	// here we need to ensure ID exists (not empty)
	copied := *book
	r.books[copied.ID] = &copied
	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	if book == nil {
		return errors.New("book is nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	value, exists := r.books[id]
	if !exists {
		return errors.New("Item not found")
	}
	value.Title = book.Title
	value.Author = book.Author
	value.ISBN = book.ISBN
	value.PublishedYear = book.PublishedYear
	value.Description = book.Description
	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return errors.New("Item not found")
	}

	delete(r.books, id)
	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*Book, 0)
	for _, b := range r.books {
		if strings.Contains(b.Author, author) {
			copied := *b
			out = append(out, &copied)
		}
	}

	return out, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*Book, 0)
	for _, b := range r.books {
		if strings.Contains(b.Title, title) {
			copied := *b
			out = append(out, &copied)
		}
	}
	return out, nil
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

func validateBook(book *Book) error {
	if book == nil {
		return errors.New("book is nil")
	}
	if book.Title == "" {
		return errors.New("title is required")
	}
	if book.Author == "" {
		return errors.New("author is required")
	}
	return nil
}

func newBookID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Implement BookService methods for DefaultBookService
func (s *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return s.repo.GetAll()
}

func (s *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return s.repo.GetByID(id)
}

func (s *DefaultBookService) CreateBook(book *Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
	if book.ID == "" {
		id, err := newBookID()
		if err != nil {
			return err
		}
		book.ID = id
	}
	return s.repo.Create(book)
}

func (s *DefaultBookService) UpdateBook(id string, book *Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
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

// HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/books")
	path = strings.Trim(path, "/")
	switch {
	case path == "" && r.Method == http.MethodGet:
		h.listBooks(w)
	case path == "" && r.Method == http.MethodPost:
		h.createBook(w, r)
	case path == "search" && r.Method == http.MethodGet:
		h.searchBooks(w, r)
	case path != "" && r.Method == http.MethodGet:
		h.getBook(w, path)
	case path != "" && r.Method == http.MethodPut:
		h.updateBook(w, r, path)
	case path != "" && r.Method == http.MethodDelete:
		h.deleteBook(w, path)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *BookHandler) listBooks(w http.ResponseWriter) {
	books, err := h.Service.GetAllBooks()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.Service.CreateBook(&book); err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) getBook(w http.ResponseWriter, id string) {
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, id string) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.Service.UpdateBook(id, &book); err != nil {
		writeServiceError(w, err)
		return
	}
	updated, err := h.Service.GetBookByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, id string) {
	if err := h.Service.DeleteBook(id); err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *BookHandler) searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	author := query.Get("author")
	title := query.Get("title")
	var (
		books []*Book
		err   error
	)
	if author != "" {
		books, err = h.Service.SearchBooksByAuthor(author)
	} else {
		books, err = h.Service.SearchBooksByTitle(title)
	}
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, books)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helper functions
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{
		StatusCode: status,
		Error:      msg,
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "Item not found":
		writeError(w, http.StatusNotFound, err.Error())
	case "title is required", "author is required", "book is nil", "Item already exists":
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
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
