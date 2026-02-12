// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBookNotFound  = errors.New("book not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDuplicateBook = errors.New("book already exists")
)

// ============================================
// MODELS
// ============================================

// Book represents a book in the database
type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
	ISBN          string `json:"isbn"`
	Description   string `json:"description"`
}

// ============================================
// REPOSITORY
// ============================================

// BookRepository defines the operations for book data access
type BookRepository interface {
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Create(book *Book) error
	Update(id string, book *Book) error
	Delete(id string) error
	SearchByAuthor(author string) ([]*Book, error)
	SearchByTitle(title string) ([]*Book, error)
	SearchByISBN(isbn string) ([]*Book, error)
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

	var books []*Book
	for _, book := range r.books {
		bookCopy := *book
		books = append(books, &bookCopy)
	}

	return books, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if book, ok := r.books[id]; ok {
		bookCopy := *book
		return &bookCopy, nil
	}

	return nil, ErrBookNotFound
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.books[book.ID]; exist {
		return ErrDuplicateBook
	}

	r.books[book.ID] = book
	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	if book == nil {
		return errors.New("book cannot be nil")
	}

	if id == "" {
		return errors.New("id cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[id]; !exists {
		return ErrBookNotFound
	}

	r.books[id] = book
	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.books[id]; !exist {
		return ErrBookNotFound
	}

	delete(r.books, id)
	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	var result []*Book

	author = strings.ToLower(strings.TrimSpace(author))

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Author), author) {
			bookCopy := *book
			result = append(result, &bookCopy)
		}
	}

	return result, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	var result []*Book

	title = strings.ToLower(strings.TrimSpace(title))

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, book := range r.books {
		if strings.Contains(strings.ToLower(book.Title), title) {
			bookCopy := *book
			result = append(result, &bookCopy)
		}
	}

	return result, nil
}

func (r *InMemoryBookRepository) SearchByISBN(isbn string) ([]*Book, error) {
	var result []*Book

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, book := range r.books {
		if strings.EqualFold(book.ISBN, isbn) {
			bookCopy := *book
			result = append(result, &bookCopy)
		}
	}

	return result, nil
}

// ============================================
// SERVICE
// ============================================

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

func (d *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return d.repo.GetAll()
}

func (d *DefaultBookService) GetBookByID(id string) (*Book, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	return d.repo.GetByID(id)
}

func (d *DefaultBookService) CreateBook(book *Book) error {
	if book == nil {
		return errors.New("book cannot be nil")
	}

	if strings.TrimSpace(book.Title) == "" {
		return errors.New("title is required")
	}

	if strings.TrimSpace(book.Author) == "" {
		return errors.New("author is required")
	}

	if strings.TrimSpace(book.ISBN) == "" {
		return errors.New("ISBN is required")
	}

	currentYear := time.Now().Year()
	if book.PublishedYear < 1000 || book.PublishedYear > currentYear {
		return fmt.Errorf("published year must be between 1000 and %d", currentYear)
	}

	existingBooks, err := d.repo.SearchByISBN(book.ISBN)
	if err != nil {
		return fmt.Errorf("failed to check for duplicates: %w", err)
	}

	if len(existingBooks) > 0 {
		return errors.New("book with this ISBN already exists")
	}

	book.ID = uuid.New().String()

	return d.repo.Create(book)
}

func (d *DefaultBookService) UpdateBook(id string, book *Book) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	if book == nil {
		return errors.New("book cannot be nil")
	}

	if strings.TrimSpace(book.Title) == "" {
		return errors.New("title is required")
	}

	if strings.TrimSpace(book.Author) == "" {
		return errors.New("author is required")
	}

	if strings.TrimSpace(book.ISBN) == "" {
		return errors.New("ISBN is required")
	}

	currentYear := time.Now().Year()
	if book.PublishedYear < 1000 || book.PublishedYear > currentYear {
		return fmt.Errorf("published year must be between 1000 and %d", currentYear)
	}

	existingBook, err := d.repo.GetByID(id)
	if err != nil {
		return err
	}

	if existingBook.ISBN != book.ISBN {
		booksWithISBN, err := d.repo.SearchByISBN(book.ISBN)
		if err != nil {
			return fmt.Errorf("failed to check for ISBN duplicates: %w", err)
		}

		if len(booksWithISBN) > 0 {
			return errors.New("book with this ISBN already exists")
		}
	}

	book.ID = id

	return d.repo.Update(id, book)
}

func (d *DefaultBookService) DeleteBook(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	return d.repo.Delete(id)
}

func (d *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	if strings.TrimSpace(author) == "" {
		return nil, errors.New("author cannot be empty")
	}

	return d.repo.SearchByAuthor(author)
}

func (d *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	return d.repo.SearchByTitle(title)
}

// ============================================
// HANDLERS
// ============================================

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
	path := r.URL.Path
	id := extractID(path)

	switch r.Method {
	case http.MethodGet:
		// Check for search queries
		if author := r.URL.Query().Get("author"); author != "" {
			h.handleSearchByAuthor(w, r)
			return
		}

		if title := r.URL.Query().Get("title"); title != "" {
			h.handleSearchByTitle(w, r)
			return
		}

		// Get by ID or get all
		if id != "" {
			h.handleGetBookByID(w, r, id)
		} else {
			h.handleGetAllBooks(w, r)
		}

	case http.MethodPost:
		if id != "" {
			respondWithError(w, http.StatusBadRequest, "ID should not be provided in URL for create operation")
			return
		}
		h.handleCreateBook(w, r)

	case http.MethodPut:
		if id == "" {
			respondWithError(w, http.StatusBadRequest, "ID is required in URL for update operation")
			return
		}
		h.handleUpdateBook(w, r, id)

	case http.MethodDelete:
		if id == "" {
			respondWithError(w, http.StatusBadRequest, "ID is required in URL for delete operation")
			return
		}
		h.handleDeleteBook(w, r, id)

	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}

func (h *BookHandler) handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.GetAllBooks()
	if err != nil {
		statusCode := mapErrorToStatusCode(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, books)
}

func (h *BookHandler) handleGetBookByID(w http.ResponseWriter, r *http.Request, id string) {
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		statusCode := mapErrorToStatusCode(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, book)
}

// handleCreateBook creates a new book
func (h *BookHandler) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Call service layer
	if err := h.Service.CreateBook(&book); err != nil {
		statusCode := mapErrorToStatusCode(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	// Return 201 Created with the created book
	respondWithJSON(w, http.StatusCreated, book)
}

// handleUpdateBook updates an existing book
func (h *BookHandler) handleUpdateBook(w http.ResponseWriter, r *http.Request, id string) {
	var book Book

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Call service layer
	if err := h.Service.UpdateBook(id, &book); err != nil {
		statusCode := mapErrorToStatusCode(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	// Return 200 OK with the updated book
	respondWithJSON(w, http.StatusOK, book)
}

// handleDeleteBook deletes a book
func (h *BookHandler) handleDeleteBook(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Service.DeleteBook(id); err != nil {
		statusCode := mapErrorToStatusCode(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	// Return 200 OK with success message
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Book deleted successfully",
	})
}

// handleSearchByAuthor searches books by author
func (h *BookHandler) handleSearchByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	books, err := h.Service.SearchBooksByAuthor(author)
	if err != nil {
		statusCode := mapErrorToStatusCode(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, books)
}

// handleSearchByTitle searches books by title
func (h *BookHandler) handleSearchByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	books, err := h.Service.SearchBooksByTitle(title)
	if err != nil {
		statusCode := mapErrorToStatusCode(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, books)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, ErrorResponse{
		StatusCode: statusCode,
		Error:      message,
	})
}

func extractID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 3 && parts[0] == "api" && parts[1] == "books" && parts[2] != "" {
		return parts[2]
	}
	return ""
}

func mapErrorToStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	errMsg := err.Error()

	// Check for specific errors
	if errors.Is(err, ErrBookNotFound) {
		return http.StatusNotFound
	}

	if errors.Is(err, ErrDuplicateBook) || strings.Contains(errMsg, "already exists") {
		return http.StatusConflict
	}

	// Validation errors
	if strings.Contains(errMsg, "required") ||
		strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "must be between") {
		return http.StatusBadRequest
	}

	// Default to internal server error
	return http.StatusInternalServerError
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
