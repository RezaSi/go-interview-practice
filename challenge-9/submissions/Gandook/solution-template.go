// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
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

func (imbr *InMemoryBookRepository) GetAll() ([]*Book, error) {
	books := make([]*Book, 0)

	imbr.mu.RLock()
	defer imbr.mu.RUnlock()

	for _, book := range imbr.books {
		books = append(books, book)
	}

	return books, nil
}

func (imbr *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	imbr.mu.RLock()
	defer imbr.mu.RUnlock()

	if book, ok := imbr.books[id]; ok {
		return book, nil
	} else {
		return nil, errors.New("book not found")
	}
}

func (imbr *InMemoryBookRepository) Create(book *Book) error {
	imbr.mu.Lock()
	defer imbr.mu.Unlock()

	if book.ISBN == "" {
		return errors.New("missing info: ISBN is required")
	}

	book.ID = book.ISBN
	if _, ok := imbr.books[book.ID]; ok {
		return errors.New("this book already exists")
	} else {
		imbr.books[book.ID] = book
		return nil
	}
}

func (imbr *InMemoryBookRepository) Update(id string, book *Book) error {
	imbr.mu.Lock()
	defer imbr.mu.Unlock()

	if _, ok := imbr.books[id]; !ok {
		return errors.New("book not found")
	} else {
		imbr.books[id] = book
		return nil
	}
}

func (imbr *InMemoryBookRepository) Delete(id string) error {
	imbr.mu.Lock()
	defer imbr.mu.Unlock()

	if _, ok := imbr.books[id]; !ok {
		return errors.New("book not found")
	} else {
		delete(imbr.books, id)
		return nil
	}
}

func (imbr *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	books := make([]*Book, 0)

	imbr.mu.RLock()
	defer imbr.mu.RUnlock()

	for _, book := range imbr.books {
		if strings.Contains(book.Author, author) {
			books = append(books, book)
		}
	}

	return books, nil
}

func (imbr *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	books := make([]*Book, 0)

	imbr.mu.RLock()
	defer imbr.mu.RUnlock()

	for _, book := range imbr.books {
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

func (dbs *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return dbs.repo.GetAll()
}

func (dbs *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return dbs.repo.GetByID(id)
}

func (dbs *DefaultBookService) CreateBook(book *Book) error {
	return dbs.repo.Create(book)
}

func (dbs *DefaultBookService) UpdateBook(id string, book *Book) error {
	return dbs.repo.Update(id, book)
}

func (dbs *DefaultBookService) DeleteBook(id string) error {
	return dbs.repo.Delete(id)
}

func (dbs *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	return dbs.repo.SearchByAuthor(author)
}

func (dbs *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	return dbs.repo.SearchByTitle(title)
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
	w.Header().Set("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	switch r.Method {
	case "GET":
		if strings.HasPrefix(r.URL.Path, "/api/books") {
			theRest := r.URL.Path[10:]
			if theRest == "" {
				bookPointers, err := h.Service.GetAllBooks()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				encodingErr := jsonEncoder.Encode(bookPointers)
				if encodingErr != nil {
					return
				}
			} else {
				if strings.HasPrefix(theRest, "/search") {
					query := r.URL.RawQuery
					if strings.HasPrefix(query, "author") {
						bookAuthor := query[7:]
						results, err := h.Service.SearchBooksByAuthor(bookAuthor)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						encodingErr := jsonEncoder.Encode(results)
						if encodingErr != nil {
							return
						}
					} else if strings.HasPrefix(query, "title") {
						bookTitle := query[6:]
						results, err := h.Service.SearchBooksByTitle(bookTitle)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						encodingErr := jsonEncoder.Encode(results)
						if encodingErr != nil {
							return
						}
					} else {
						http.Error(w, "invalid search query", http.StatusBadRequest)
						return
					}
				} else {
					bookID := theRest[1:]
					book, err := h.Service.GetBookByID(bookID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusNotFound)
					}

					encodingErr := jsonEncoder.Encode(book)
					if encodingErr != nil {
						return
					}
				}
			}
		} else {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
	case "POST":
		if r.URL.Path == "/api/books" {
			var bookPointer *Book
			err := jsonDecoder.Decode(&bookPointer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			addingErr := h.Service.CreateBook(bookPointer)
			if addingErr != nil {
				http.Error(w, addingErr.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			encodingErr := jsonEncoder.Encode(*bookPointer)
			if encodingErr != nil {
				return
			}
		} else {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
	case "PUT":
		if strings.HasPrefix(r.URL.Path, "/api/books/") {
			var bookPointer *Book
			err := jsonDecoder.Decode(&bookPointer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			updatingErr := h.Service.UpdateBook(r.URL.Path[11:], bookPointer)
			if updatingErr != nil {
				http.Error(w, updatingErr.Error(), http.StatusNotFound)
				return
			}

			encodingErr := jsonEncoder.Encode(*bookPointer)
			if encodingErr != nil {
				return
			}
		} else {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
	case "DELETE":
		if strings.HasPrefix(r.URL.Path, "/api/books/") {
			deletingErr := h.Service.DeleteBook(r.URL.Path[11:])
			if deletingErr != nil {
				http.Error(w, deletingErr.Error(), http.StatusNotFound)
				return
			}
		} else {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
	}
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

	// Create a new router and register endpoints
	http.HandleFunc("/api/books", handler.HandleBooks)
	http.HandleFunc("/api/books/", handler.HandleBooks)

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
