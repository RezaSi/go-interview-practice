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

var ErrBookNotFound = errors.New("book not found")

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
	books := []*Book{}
	r.mu.RLock()
	defer r.mu.RUnlock()
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
		return nil, ErrBookNotFound
	}

	return book, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.Lock()
	r.books[book.ID] = book
	r.mu.Unlock()

	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; !ok {
		return ErrBookNotFound
	}

	r.books[id] = book

	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; !ok {
		return ErrBookNotFound
	}

	delete(r.books, id)

	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	books := []*Book{}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, book := range r.books {
		if strings.Contains(book.Author, author) {
			books = append(books, book)
		}
	}

	return books, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	books := []*Book{}

	r.mu.RLock()
	defer r.mu.RUnlock()
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

func (s *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return s.repo.GetAll()
}

func (s *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return s.repo.GetByID(id)
}

func (s *DefaultBookService) CreateBook(book *Book) error {
	book.ID = uuid.NewString()
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

// HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this method to handle all book endpoints
	// Use the path and method to determine the appropriate action
	// Call the service methods accordingly
	// Return appropriate status codes and JSON responses

	switch r.Method {
	case http.MethodGet:
		h.HandleGetBooks(w, r)

	case http.MethodPost:
		h.HandleCreateBook(w, r)

	case http.MethodDelete:
		h.HandleDeleteBook(w, r)

	case http.MethodPut:
		h.HandleUpdateBook(w, r)
	}
}

func (h *BookHandler) HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.Contains(path, "search") {
		h.HandleSearchBook(w, r)
		return
	}

	if path == "/api/books" {
		books, err := h.Service.GetAllBooks()
		if err != nil {
			WriteInternalServerError(w, err)
			return
		}
		WriteRespone(w, books)
		return
	}

	id := GetPathID(r)
	if id == "" {
		WriteBadRequestError(w, "path parameter is rquired")
		return
	}
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			WriteNotFoundError(w, err)
			return
		}
		WriteInternalServerError(w, err)
		return
	}

	WriteRespone(w, book)
}

func (h *BookHandler) HandleSearchBook(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.RawQuery, "author") && strings.Contains(r.URL.RawQuery, "title") {
		WriteBadRequestError(w, "one query parameter must be provided")
		return
	}

	if strings.Contains(r.URL.RawQuery, "author") {
		author := r.URL.Query().Get("author")
		if author == "" {
			WriteBadRequestError(w, "empty query author parameter")
			return
		}

		books, err := h.Service.SearchBooksByAuthor(author)
		if err != nil {
			if errors.Is(err, ErrBookNotFound) {
				WriteNotFoundError(w, err)
				return
			}
			WriteInternalServerError(w, err)
		}

		WriteRespone(w, books)
		return
	}

	if strings.Contains(r.URL.RawQuery, "title") {
		title := r.URL.Query().Get("title")
		if title == "" {
			WriteBadRequestError(w, "empty query author parameter")
			return
		}

		books, err := h.Service.SearchBooksByTitle(title)
		if err != nil {
			if errors.Is(err, ErrBookNotFound) {
				WriteNotFoundError(w, err)
				return
			}
			WriteInternalServerError(w, err)
		}

		WriteRespone(w, books)
		return
	}
}

func (h *BookHandler) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	book, err := GetBookFromBody(r)
	if err != nil {
		WriteBadRequestError(w, err.Error())
		return
	}

	ValidateBook(w, book)

	err = h.Service.CreateBook(book)
	if err != nil {
		WriteInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	WriteRespone(w, book)
}

func (h *BookHandler) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id := GetPathID(r)
	if id == "" {
		WriteBadRequestError(w, "path parameter is rquired")
		return
	}

	book, err := GetBookFromBody(r)
	if err != nil {
		WriteBadRequestError(w, err.Error())
		return
	}

	ValidateBook(w, book)

	if err := h.Service.UpdateBook(id, book); err != nil {
		if errors.Is(err, ErrBookNotFound) {
			WriteNotFoundError(w, err)
			return
		}
		WriteInternalServerError(w, err)
		return
	}
	WriteRespone(w, book)
}

func (h *BookHandler) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id := GetPathID(r)
	if id == "" {
		WriteBadRequestError(w, "path parameter is rquired")
		return
	}

	err := h.Service.DeleteBook(id)
	if errors.Is(err, ErrBookNotFound) {
		WriteNotFoundError(w, err)
		return
	}
	if err != nil {
		WriteInternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func WriteRespone(w http.ResponseWriter, data any) {
	json.NewEncoder(w).Encode(data)
}

func GetPathID(r *http.Request) string {
	id := strings.TrimPrefix(r.URL.Path, "/api/books")
	if id == "" || id == "/" {
		return ""
	}
	id = id[1:]
	return id
}

func GetBookFromBody(r *http.Request) (*Book, error) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return &book, nil
}

func ValidateBook(w http.ResponseWriter, book *Book) {
	if book.Author == "" {
		WriteBadRequestError(w, "author is missing")
		return
	}

	if book.Title == "" {
		WriteBadRequestError(w, "title is missing")
		return
	}

	if book.ISBN == "" {
		WriteBadRequestError(w, "ISBN is missing")
		return
	}

	if book.PublishedYear == 0 {
		WriteBadRequestError(w, "published year is missing")
		return
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

func WriteError(w http.ResponseWriter, errRes ErrorResponse) {
	w.WriteHeader(errRes.StatusCode)
	json.NewEncoder(w).Encode(errRes)
}

func WriteBadRequestError(w http.ResponseWriter, err string) {
	WriteError(w, ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Error:      err,
	})
}

func WriteNotFoundError(w http.ResponseWriter, err error) {
	WriteError(w, ErrorResponse{
		StatusCode: http.StatusNotFound,
		Error:      err.Error(),
	})
}

func WriteInternalServerError(w http.ResponseWriter, err error) {
	WriteError(w, ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Error:      err.Error(),
	})
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
