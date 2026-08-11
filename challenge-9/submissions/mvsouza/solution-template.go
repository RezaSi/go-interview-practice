// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"log"
	"maps"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const NotFoundMessage string = "Not Found"

var NotFound error = errors.New(NotFoundMessage)

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

func (s *InMemoryBookRepository) GetAll() ([]*Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return slices.Collect(maps.Values(s.books)), nil
}

func (s *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if b, ok := s.books[id]; ok {
		return b, nil
	}
	return nil, NotFound
}

func (s *InMemoryBookRepository) Create(book *Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	if _, ok := s.books[id]; !ok {
		book.ID = id
		s.books[id] = book
		return nil
	}
	return NotFound
}

func (s *InMemoryBookRepository) Update(id string, book *Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; ok {
		s.books[id] = book
		return nil
	}
	return NotFound
}

func (s *InMemoryBookRepository) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; ok {
		delete(s.books, id)
		return nil
	}
	return NotFound
}

func (s *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if author == "" {
		return nil, errors.New("search can't be empty")
	}
	var filteredBooks []*Book
	for _, book := range s.books {
		if strings.Contains(book.Author, author) {
			filteredBooks = append(filteredBooks, book)
		}
	}
	return filteredBooks, nil
}

func (s *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if title == "" {
		return nil, errors.New("search can't be empty")
	}
	var filteredBooks []*Book
	for _, book := range s.books {
		if strings.Contains(book.Title, title) {
			filteredBooks = append(filteredBooks, book)
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
	Router  *mux.Router
}

// NewBookHandler creates a new book handler
func NewBookHandler(service BookService) *BookHandler {
	router := mux.NewRouter()
	h := BookHandler{
		Service: service,
		Router:  router,
	}

	router.HandleFunc("/api/books", h.GetAllBooks).Methods("GET")
	router.HandleFunc("/api/books", h.CreateBook).Methods("POST")
	router.HandleFunc("/api/books/search", h.SearchBook).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.GetBookById).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.DeleteBook).Methods("DELETE")
	router.HandleFunc("/api/books/{id}", h.EditBook).Methods("PUT")
	return &h
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.GetAllBooks()
	if err == nil {
		err = writeJsonReponse(w, books)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	book, err := h.Service.GetBookByID(id)
	if err == nil {
		err = writeJsonReponse(w, book)
	} else if errors.Is(err, NotFound) {
		http.NotFound(w, r)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if book.Title == "" {
		http.Error(w, "title cannot be empty", http.StatusBadRequest)
		return
	}
	err = h.Service.CreateBook(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BookHandler) EditBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	id := mux.Vars(r)["id"]
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if book.Title == "" {
		http.Error(w, "title cannot be empty", http.StatusBadRequest)
		return
	}
	if err == nil {
		err = h.Service.UpdateBook(id, &book)
	}
	if err != nil && errors.Is(err, NotFound) {
		http.NotFound(w, r)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := h.Service.DeleteBook(id)
	if err != nil && errors.Is(err, NotFound) {
		http.NotFound(w, r)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BookHandler) SearchBook(w http.ResponseWriter, r *http.Request) {
	var books []*Book
	var err error
	author := r.URL.Query().Get("author")
	title := r.URL.Query().Get("title")
	if len(author) > 0 {
		books, err = h.Service.SearchBooksByAuthor(author)
	} else {
		books, err = h.Service.SearchBooksByTitle(title)
	}
	if err == nil {
		err = writeJsonReponse(w, books)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

func writeJsonReponse(w http.ResponseWriter, o any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(o)
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
