// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"io"
	
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

// InMemoryBookRepository implements BookRepository using in-memory storage
type InMemoryBookRepository struct {
	books map[string]*Book
	mu    sync.RWMutex
}

var errBookAlreadyExist = errors.New("book already exist")
var errBookDoesNotExist = errors.New("book does not exist")

// NewInMemoryBookRepository creates a new in-memory book repository
func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books: make(map[string]*Book),
	}
}

func (r *InMemoryBookRepository) GetAll() ([]*Book, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    books := []*Book{}
    for _, book := range r.books {
        books = append(books, book)
    }
    return books, nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*Book, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    if book, exists := r.books[id]; !exists {
        return nil, errBookDoesNotExist
    } else {
        return book, nil
    }
}

func (r *InMemoryBookRepository) Create(book *Book) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.books[book.ID]; exists {
        return errBookAlreadyExist
    }
    r.books[book.ID] = book
    return nil
}

func (r *InMemoryBookRepository) Update(id string, book *Book) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.books[id]; !exists {
        return errBookDoesNotExist
    }
    r.books[id] = book
    return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.books[id]; !exists {
        return errBookDoesNotExist
    }
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
	switch r.Method {
	    case http.MethodGet: {
	        if strings.HasPrefix(r.URL.Path, "/api/books/search") {
	            query := r.URL.Query()

	            if author := query.Get("author"); author != "" {
	                books, _ := h.Service.SearchBooksByAuthor(author)
	                writeJsonResponse(w, books)
	                return
	            }

	            if title := query.Get("title"); title != "" {
	                books, _ := h.Service.SearchBooksByTitle(title)
	                writeJsonResponse(w, books)
	                return
	            }
	        } else if r.URL.Path == "/api/books" {
	            books, err := h.Service.GetAllBooks()

	            if err != nil {
	                w.WriteHeader(http.StatusInternalServerError)
	                return
    	        }

    	        writeJsonResponse(w, books)
	            return
	        } else if strings.HasPrefix(r.URL.Path, "/api/books/") {
    	        book, err := h.Service.GetBookByID(getID(r.URL.Path))

    	        if err != nil {
    	            w.WriteHeader(http.StatusNotFound)
    	            return
    	        }

    	        writeJsonResponse(w, book)
    	        return
	        }
	        w.WriteHeader(http.StatusBadRequest)
	    }
	    case http.MethodPost: {
	        if r.URL.Path != "/api/books" {
	            w.WriteHeader(http.StatusBadRequest)
	            return
	        }

            book, err := readBookFromRequest(r)
	        if err != nil {
	            w.WriteHeader(http.StatusBadRequest)
        		return
	        }
            
            if book.Title == "" || book.Author == "" || book.PublishedYear == 0 || book.ISBN == "" {
                w.WriteHeader(http.StatusBadRequest)
	            return
            }
            
            book.ID = uuid.New().String()

            if err = h.Service.CreateBook(book); err != nil {
                w.WriteHeader(http.StatusConflict)
        		return
            }
            
            w.WriteHeader(http.StatusCreated)
            writeJsonResponse(w, book)
	    }
	    case http.MethodPut: {
	        if !strings.HasPrefix(r.URL.Path, "/api/books/") {
	            w.WriteHeader(http.StatusBadRequest)
                return
	        }
	        
	        book, err := readBookFromRequest(r)
	        if err != nil {
	            w.WriteHeader(http.StatusBadRequest)
        		return
	        }
	        
	        pathID := getID(r.URL.Path)
            
            if book.Title == "" || book.Author == "" || book.PublishedYear == 0 || book.ISBN == "" {
                w.WriteHeader(http.StatusBadRequest)
	            return
            }
            
            if book.ID != "" && book.ID != pathID {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            book.ID = pathID
            
            if err = h.Service.UpdateBook(pathID, book); err != nil {
                w.WriteHeader(http.StatusNotFound)
        		return
            }

            writeJsonResponse(w, book)
	    }
	    case http.MethodDelete: {
	        if !strings.HasPrefix(r.URL.Path, "/api/books/") {
	            w.WriteHeader(http.StatusBadRequest)
                return
	        }
	        
	        if err := h.Service.DeleteBook(getID(r.URL.Path)); err != nil {
	            w.WriteHeader(http.StatusNotFound)
        		return
	        }
	        
	        w.WriteHeader(http.StatusOK)
	    }
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// Helpers
func writeJsonResponse(w http.ResponseWriter, response interface{}) {
	bytes, err := json.Marshal(response)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(bytes)
	}
}

func getID(path string) string {
    return strings.Replace(path, "/api/books/", "", 1)
}

func readBookFromRequest(r *http.Request) (*Book, error) {
    defer r.Body.Close()

    bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var book Book
    if err = json.Unmarshal(bodyBytes, &book); err != nil {
		return nil, err
    }
    
    return &book, nil
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