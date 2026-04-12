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

func (r *InMemoryBookRepository) GetAll() ([]*Book, error) {
	books := []*Book{}
	r.mu.Lock()
	for _, book := range r.books {
		books = append(books, book)
	}
	r.mu.Unlock()
	return books, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, book := range r.books {
		if book.ID == id {
			return book, nil
		}
	}
	return nil, errors.New("Book not found")
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if book.Author == "" || book.Description == "" || book.ISBN == "" || book.PublishedYear == 0 || book.Title == "" {
		return errors.New("Book is missing data")
	}
	if len(book.ID) > 0 {
		if _, exists := r.books[book.ID]; exists {
			return errors.New("Book already exists")
		}
	} else {
		book.ID = uuid.NewString()
	}
	r.books[book.ID] = book
	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	book.ID = id
	if _, ok := r.books[id]; !ok {
		return errors.New("Book doesn't exist")
	} else {
		r.books[id] = book
	}
	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[id]; !ok {
		return errors.New("Book doesn't exist")
	} else {
		delete(r.books, id)
	}
	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	results := []*Book{}
	for _, book := range r.books {
		if strings.Contains(book.Author, author) {
			results = append(results, book)
		}
	}
	return results, nil
}

func (r *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	results := []*Book{}
	for _, book := range r.books {
		if strings.Contains(book.Title, title) {
			results = append(results, book)
		}
	}
	return results, nil
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

func (bs *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return bs.repo.GetAll()
}

func (bs *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return bs.repo.GetByID(id)
}

func (bs *DefaultBookService) CreateBook(book *Book) error {
	return bs.repo.Create(book)
}

func (bs *DefaultBookService) UpdateBook(id string, book *Book) error {
	return bs.repo.Update(id, book)
}

func (bs *DefaultBookService) DeleteBook(id string) error {
	return bs.repo.Delete(id)
}

func (bs *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	return bs.repo.SearchByAuthor(author)
}

func (bs *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	return bs.repo.SearchByTitle(title)
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
	router := mux.NewRouter()
	router.HandleFunc("/api/books", h.getAllBooks).Methods("GET")
	router.HandleFunc("/api/books", h.createBook).Methods("POST")
	router.HandleFunc("/api/books/search", h.searchBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.getBookByID).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", h.deleteBook).Methods("DELETE")

	router.ServeHTTP(w, r)
}

func (h *BookHandler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	if books, err := h.Service.GetAllBooks(); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
	} else {
		writeJsonResponse(w, http.StatusOK, books)
	}
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
	} else {
		if err := h.Service.CreateBook(&book); err != nil {
			writeJsonError(w, http.StatusBadRequest, err)
		} else {
			writeJsonResponse(w, http.StatusCreated, book)
		}
	}
}

func (h *BookHandler) searchBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	title := r.URL.Query().Get("title")

	if author != "" {
		h.searchBooksByAuthor(w, r, author)
	} else if title != "" {
		h.searchBooksByTitle(w, r, title)
	} else {
		writeJsonError(w, http.StatusInternalServerError, errors.New("Neither author nor title set"))
	}
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
	} else {
		if err := h.Service.UpdateBook(id, &book); err != nil {
			writeJsonError(w, http.StatusNotFound, err)
		} else {
			writeJsonResponse(w, http.StatusOK, book)
		}
	}
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.Service.DeleteBook(id); err != nil {
		writeJsonError(w, http.StatusNotFound, err)
	} else {
		writeJsonResponse(w, http.StatusOK, nil)
	}
}

func (h *BookHandler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if book, err := h.Service.GetBookByID(id); err != nil {
		writeJsonError(w, http.StatusNotFound, err)
	} else {
		writeJsonResponse(w, http.StatusOK, book)
	}
}

func (h *BookHandler) searchBooksByAuthor(w http.ResponseWriter, r *http.Request, author string) {
	if books, err := h.Service.SearchBooksByAuthor(author); err != nil {
		writeJsonError(w, http.StatusNotFound, err)
	} else {
		writeJsonResponse(w, http.StatusOK, books)
	}
}

func (h *BookHandler) searchBooksByTitle(w http.ResponseWriter, r *http.Request, title string) {
	if books, err := h.Service.SearchBooksByTitle(title); err != nil {
		writeJsonError(w, http.StatusNotFound, err)
	} else {
		writeJsonResponse(w, http.StatusOK, books)
	}
}

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

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

func main() {
	// Initialize the repository, service, and handler
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	http.HandleFunc("/api/books", handler.HandleBooks)
	http.HandleFunc("/api/books/", handler.HandleBooks)

	// Start the server
	log.Println("Server starting on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
