// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ValidateBook(book *Book) error
}

func (this *InMemoryBookRepository) ValidateBook(book *Book) error {
	if book == nil {
		return fmt.Errorf("book is nil")
	}

	if book.Title == "" {
		return fmt.Errorf("title is empty")
	}

	if book.Author == "" {
		return fmt.Errorf("author is empty")
	}

	if book.PublishedYear == 0 {
		return fmt.Errorf("published year is empty")
	}

	if book.ISBN == "" {
		return fmt.Errorf("isbn is empty")
	}

	if book.Description == "" {
		return fmt.Errorf("description is empty")
	}
	return nil
}

func (this *InMemoryBookRepository) GetAll() ([]*Book, error) {
	var books []*Book

	for _, book := range this.books {
		books = append(books, book)
	}

	return books, nil
}

func (this *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	book, ok := this.books[id]
	if !ok {
		return nil, fmt.Errorf("book not found")
	}

	return book, nil
}
func (this *InMemoryBookRepository) Create(book *Book) error {
	book.ID = fmt.Sprintf("%d", len(this.books)+1)

	this.books[book.ID] = book

	return nil
}

func (this *InMemoryBookRepository) Update(id string, book *Book) error {
	if _, ok := this.books[id]; !ok {
		return fmt.Errorf("book not found")
	}

	this.books[id] = book

	return nil
}

func (this *InMemoryBookRepository) Delete(id string) error {
	delete(this.books, id)

	return nil
}

func (this *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	var books []*Book

	for _, book := range this.books {
		if strings.Contains(book.Author, author) {
			books = append(books, book)
		}
	}

	return books, nil
}

func (this *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	var books []*Book

	for _, book := range this.books {
		if strings.Contains(book.Title, title) {
			books = append(books, book)
		}
	}

	return books, nil
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
func (this *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return this.repo.GetAll()
}

func (this *DefaultBookService) GetBookByID(id string) (*Book, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}

	return this.repo.GetByID(id)
}

func (this *DefaultBookService) CreateBook(book *Book) error {
	err := this.repo.ValidateBook(book)
	if err != nil {
		return err
	}

	return this.repo.Create(book)
}

func (this *DefaultBookService) UpdateBook(id string, book *Book) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}

	err := this.repo.ValidateBook(book)
	if err != nil {
		return err
	}

	return this.repo.Update(id, book)
}

func (this *DefaultBookService) DeleteBook(id string) error {
	return this.repo.Delete(id)
}

func (this *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	return this.repo.SearchByAuthor(author)
}

func (this *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	return this.repo.SearchByTitle(title)
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
	if r.URL.Path == "/api/books" {
		switch r.Method {
		case http.MethodGet:
			// Handle GET requests
			books, err := h.Service.GetAllBooks()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(books)

		case http.MethodPost:
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var book Book
			if err := json.Unmarshal(data, &book); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = h.Service.CreateBook(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(book)
		// Handle POST requests

		case http.MethodPut:
			w.WriteHeader(http.StatusMethodNotAllowed)

		case http.MethodDelete:
			w.WriteHeader(http.StatusMethodNotAllowed)
			// Handle DELETE requests
		}
	} else if strings.HasPrefix(r.URL.Path, "/api/books/") {
		switch r.Method {
		case http.MethodGet:
			query := r.URL.Query()
			if len(query) > 0 {
				author := query.Get("author")
				title := query.Get("title")
				if author != "" {
					books, err := h.Service.SearchBooksByAuthor(author)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(books)
				} else if title != "" {
					books, err := h.Service.SearchBooksByTitle(title)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(books)
				}

			} else {
				query := strings.Split(r.URL.Path, "/")
				id := query[len(query)-1]
				book, err := h.Service.GetBookByID(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(book)
			}

			// Handle GET requests
		case http.MethodPost:
			fmt.Println("PATH IS : ", r.URL.Path, " Method: POST")
		// Handle POST requests

		case http.MethodPut:
			query := strings.Split(r.URL.Path, "/")
			id := query[len(query)-1]

			bookInRecords, err := h.Service.GetBookByID(id)
			if err != nil {
				http.Error(w, "Book not found", http.StatusNotFound)
				return
			}

			if bookInRecords == nil {
				http.Error(w, "Book not found", http.StatusNotFound)
				return
			}

			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var book Book
			if err := json.Unmarshal(data, &book); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = h.Service.UpdateBook(id, &book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)

		// Handle PUT requests

		case http.MethodDelete:
			query := strings.Split(r.URL.Path, "/")
			id := query[len(query)-1]

			bookInRecords, err := h.Service.GetBookByID(id)
			if err != nil {
				http.Error(w, "Book not found", http.StatusNotFound)
				return
			}

			if bookInRecords == nil {
				http.Error(w, "Book not found", http.StatusNotFound)
				return
			}

			err = h.Service.DeleteBook(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			// Handle DELETE requests
		}
	}

}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
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
