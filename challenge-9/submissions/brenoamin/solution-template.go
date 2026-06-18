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
	"github.com/gorilla/mux"
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
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, ok := r.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if book.Title == "" || book.Author == "" || book.ISBN == "" || book.PublishedYear == 0 || book.Description == "" {
		return errors.New("book is missing required data")
	}
	if _, ok := r.books[book.ID]; ok {
		return errors.New("book already exists")
	}
	if book.ID == "" {
		book.ID = uuid.NewString()
	}
	r.books[book.ID] = book
	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[id] = book
	return nil
}
func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.books, id)
	return nil
}
func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	books := make([]*Book, 0, len(r.books))
	for _, book := range r.books {
		if strings.Contains(book.Author, author) {
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
		if strings.Contains(book.Title, title) {
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

func (s *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return s.repo.GetAll()
}
func (s *DefaultBookService) GetBookByID(id string) (*Book, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return book, nil
}
func (s *DefaultBookService) CreateBook(book *Book) error {
	return s.repo.Create(book)
}

func (s *DefaultBookService) UpdateBook(id string, book *Book) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("book not found")
	}
	return s.repo.Update(id, book)
}

func (s *DefaultBookService) DeleteBook(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("book not found")
	}
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
	// Use the path and method to determine the appropriate action
	// Call the service methods accordingly
	// Return appropriate status codes and JSON responses
	router := mux.NewRouter()
	router.HandleFunc("/api/books", h.getAllBooks).Methods("GET")
	router.HandleFunc("/api/books", h.createBook).Methods("POST")
	router.HandleFunc("/api/books/search", h.searchBooksBy).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.getBookById).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", h.deleteBook).Methods("DELETE")

	router.ServeHTTP(w, r)
}

func (h *BookHandler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.GetAllBooks()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
	}
	writeJsonResponse(w, http.StatusOK, books)
}
func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
	} else {
		if err := h.Service.CreateBook(&book); err != nil {
			writeJsonError(w, http.StatusBadRequest, err)
		}
		writeJsonResponse(w, http.StatusCreated, book)
	}
}

func (h *BookHandler) searchBooksBy(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")

	if title == "" && author == "" {
		writeJsonError(w, http.StatusBadRequest, errors.New("provide title or author"))
	}
	if title != "" {
		h.searchBooksByTitle(w, title)
	} else if author != "" {
		h.searchBooksByAuthor(w, author)
	}
}

func (h *BookHandler) searchBooksByTitle(w http.ResponseWriter, title string) {
	books, err := h.Service.SearchBooksByTitle(title)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
		return
	}
	writeJsonResponse(w, http.StatusOK, books)
}

func (h *BookHandler) searchBooksByAuthor(w http.ResponseWriter, author string) {
	books, err := h.Service.SearchBooksByAuthor(author)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
		return
	}
	writeJsonResponse(w, http.StatusOK, books)
}

func (h *BookHandler) getBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	book, err := h.Service.GetBookByID(id)
	if err != nil {
		writeJsonError(w, http.StatusNotFound, err)
	}
	writeJsonResponse(w, http.StatusOK, book)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	book := &Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
	}

	err = h.Service.UpdateBook(book.ID, book)
	if err != nil {
		writeJsonError(w, http.StatusNotFound, err)
	}
	writeJsonResponse(w, http.StatusOK, book)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.Service.DeleteBook(id)
	if err != nil {
		writeJsonError(w, http.StatusNotFound, err)
	}
	writeJsonResponse(w, http.StatusOK, nil)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helper functions
// ...

func writeJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJsonError(w http.ResponseWriter, statusCode int, err error) {
	writeJsonResponse(w, statusCode, ErrorResponse{
		StatusCode: statusCode,
		Error:      err.Error(),
	})
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
