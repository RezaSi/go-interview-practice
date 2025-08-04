// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type ServiceHandler int

const (
	NotFound ServiceHandler = iota
	GetAllBooks
	CreateBook
	GetBookByID
	UpdateBook
	DeleteBook
	Search
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

	// this is not valid in older go versions
	// books := slices.Collect(maps.Values(r.books))

	// old way
	books := make([]*Book, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}

	return books, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if book, ok := r.books[id]; ok {
		return book, nil
	}

	return nil, nil
}

func (r *InMemoryBookRepository) Create(book *Book) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// use ISDN as id key
	idISBN := book.ISBN
	book.ID = idISBN

	// check if book already exists in repo
	if _, ok := r.books[idISBN]; ok {
		return fmt.Errorf("book already in store ISBN:%s", book.ISBN)
	}

	// add book to repo
	r.books[idISBN] = book

	return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	existingBook, ok := r.books[id]
	if !ok {
		return fmt.Errorf("cannot find book, id:%s", id)
	}

	// update fields
	existingBook.Title = book.Title
	existingBook.Author = book.Author
	existingBook.PublishedYear = book.PublishedYear
	existingBook.Description = book.Description

	// change of ISBN triggers a new id
	if existingBook.ISBN != book.ISBN {
		existingBook.ISBN = book.ISBN
		existingBook.ID = book.ISBN
		r.books[book.ISBN] = existingBook
		delete(r.books, id)
	}

	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	delete(r.books, id)
	return nil
}

func (r *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := []*Book{}
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

	books := []*Book{}
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
	serviceHandler := getServiceHandler(r.URL.Path, r.Method)

	// Call the service methods accordingly
	switch serviceHandler {
	case GetAllBooks:
		books, err := h.Service.GetAllBooks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return appropriate status codes and JSON responses
		writeJsonResponse(w, books)
	case CreateBook:
		var book Book

		// extract book from POST body
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate book fields
		if invalid(book) {
			err := fmt.Errorf("incomplete book details, missing Title/Author/Year/ISBN/Desc")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// create book adding an id
		if err := h.Service.CreateBook(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return appropriate status codes and JSON responses
		w.WriteHeader(http.StatusCreated)
		writeJsonResponse(w, book)
	case GetBookByID:
		// extract id from path
		reID := regexp.MustCompile(`\d+\-\d+$`)
		id := reID.FindString(r.URL.Path)

		book, err := h.Service.GetBookByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if book == nil {
			err := fmt.Errorf("book not found, ISBN:%s", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Return appropriate status codes and JSON responses
		w.WriteHeader(http.StatusOK)
		writeJsonResponse(w, book)

	case UpdateBook:
		var book Book

		// extract book from POST body
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate book fields
		if invalid(book) {
			err := fmt.Errorf("incomplete book details, missing Title/Author/Year/ISBN/Desc")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// extract id from path
		reID := regexp.MustCompile(`\d+\-\d+$`)
		id := reID.FindString(r.URL.Path)

		if id != book.ID {
			err := fmt.Errorf("book id in URL different to book id in body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// verify book exists
		existingBook, err := h.Service.GetBookByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if existingBook == nil {
			err := fmt.Errorf("book not found, ISBN:%s", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.Service.UpdateBook(id, &book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return appropriate status codes and JSON responses
		w.WriteHeader(http.StatusOK)
		writeJsonResponse(w, existingBook)

	case DeleteBook:
		// extract id from path
		reID := regexp.MustCompile(`\d+\-\d+$`)
		id := reID.FindString(r.URL.Path)

		// verify book exists
		existingBook, err := h.Service.GetBookByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if existingBook == nil {
			err := fmt.Errorf("book not found, ISBN:%s", id)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.Service.DeleteBook(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return appropriate status codes and JSON responses
		w.WriteHeader(http.StatusOK)

	case Search:
		books := []*Book{}
		query := r.URL.Query()

		author := query.Get("author")
		if author != "" {
			books, err := h.Service.SearchBooksByAuthor(author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Return appropriate status codes and JSON responses
			writeJsonResponse(w, books)
			return
		}

		title := query.Get("title")
		if title != "" {
			var err error
			books, err := h.Service.SearchBooksByTitle(title)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Return appropriate status codes and JSON responses
			writeJsonResponse(w, books)
			return
		}

		// Return appropriate status codes and JSON responses
		writeJsonResponse(w, books)
		return

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helper functions
func getServiceHandler(path, method string) ServiceHandler {
	if path == "/api/books" {
		switch method {
		case "GET":
			return GetAllBooks
		case "POST":
			return CreateBook
		default:
			return NotFound
		}
	}

	reGetBookByID := regexp.MustCompile(`^/api/books/\d+\-\d+$`)
	if reGetBookByID.MatchString(path) {
		switch method {
		case "GET":
			return GetBookByID
		case "PUT":
			return UpdateBook
		case "DELETE":
			return DeleteBook
		default:
			return NotFound
		}
	}

	if path == "/api/books/search" {
		return Search
	}

	return NotFound
}

func writeJsonResponse(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// return results
	fmt.Fprintf(w, "%+v\n", string(response))
}

func invalid(book Book) bool {
	return strings.TrimSpace(book.Title) == "" ||
		strings.TrimSpace(book.Author) == "" ||
		book.PublishedYear == 0 ||
		strings.TrimSpace(book.ISBN) == "" ||
		strings.TrimSpace(book.Description) == ""
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
