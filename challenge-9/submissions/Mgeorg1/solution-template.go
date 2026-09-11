// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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

var (
	ErrInMemoRepoIsNotInited = errors.New("InMemoryRepository: books is not initialized")
	ErrInMemoRepoNotFound    = errors.New("InMemoryRepository: book is not found")
	ErrInMemoRepoNilBRPtr    = errors.New("InMemoryRepository: book pointer parameter is nil")
)

func (br *InMemoryBookRepository) GetAll() ([]*Book, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	if br.books == nil {
		return nil, ErrInMemoRepoIsNotInited
	}

	ret := []*Book{}

	for _, v := range br.books {
		bookCopy := *v
		ret = append(ret, &bookCopy)
	}

	return ret, nil
}

func (br *InMemoryBookRepository) getByID(id string) (*Book, error) {
	if br.books == nil {
		return nil, ErrInMemoRepoIsNotInited
	}

	bookPtr, exists := br.books[id]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrInMemoRepoNotFound, id)
	}

	return bookPtr, nil
}

func (br *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	bookPtr, err := br.getByID(id)
	if err != nil {
		return nil, err
	}

	book := *bookPtr
	return &book, nil
}

func (br *InMemoryBookRepository) Create(book *Book) error {
	br.mu.Lock()
	defer br.mu.Unlock()

	if book == nil {
		return ErrInMemoRepoNilBRPtr
	}

	if br.books == nil {
		return ErrInMemoRepoIsNotInited
	}

	id := rand.Uint64()
	book.ID = strconv.FormatUint(id, 10)

	bookCopy := *book
	br.books[book.ID] = &bookCopy

	return nil
}

func (br *InMemoryBookRepository) Update(id string, book *Book) error {
	br.mu.Lock()
	defer br.mu.Unlock()

	if book == nil {
		return ErrInMemoRepoNilBRPtr
	}

	if br.books == nil {
		return ErrInMemoRepoIsNotInited
	}

	if id != book.ID {
		return fmt.Errorf(
			"InMemoryRepository: id of a book for updating is not equal to book.id: %s != %s",
			id,
			book.ID,
		)
	}

	_, exists := br.books[id]
	if !exists {
		return fmt.Errorf("%w: %s", ErrInMemoRepoNotFound, id)
	}

	bookCopy := *book
	br.books[id] = &bookCopy

	return nil
}

func (br *InMemoryBookRepository) Delete(id string) error {
	br.mu.Lock()
	defer br.mu.Unlock()

	if br.books == nil {
		return ErrInMemoRepoIsNotInited
	}

	_, exists := br.books[id]
	if !exists {
		return fmt.Errorf("%w: %s", ErrInMemoRepoNotFound, id)
	}

	delete(br.books, id)

	return nil
}

func (br *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	if br.books == nil {
		return nil, ErrInMemoRepoIsNotInited
	}

	result := []*Book{}

	for _, v := range br.books {
		if strings.Contains(v.Author, author) {
			bookCopy := *v
			result = append(result, &bookCopy)
		}
	}

	return result, nil
}

func (br *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	if br.books == nil {
		return nil, ErrInMemoRepoIsNotInited
	}

	result := []*Book{}

	for _, v := range br.books {
		if strings.Contains(v.Title, title) {
			bookCopy := *v
			result = append(result, &bookCopy)
		}
	}

	return result, nil
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
	res, err := bs.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all books error: %w", err)
	}

	return res, nil
}

func (bs *DefaultBookService) GetBookByID(id string) (*Book, error) {
	res, err := bs.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("getting book by id error: %w", err)
	}

	return res, nil
}

func (bs *DefaultBookService) CreateBook(book *Book) error {
	err := bs.repo.Create(book)
	if err != nil {
		return fmt.Errorf("creating the book error: %w", err)
	}

	return nil
}

func (bs *DefaultBookService) UpdateBook(id string, book *Book) error {
	err := bs.repo.Update(id, book)
	if err != nil {
		return fmt.Errorf("updating the book error: %w", err)
	}

	return nil
}

func (bs *DefaultBookService) DeleteBook(id string) error {
	err := bs.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("deleting the book error: %w", err)
	}

	return nil
}

func (bs *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	res, err := bs.repo.SearchByAuthor(author)
	if err != nil {
		return nil, fmt.Errorf("searching by the author error: %w", err)
	}

	return res, nil
}

func (bs *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	res, err := bs.repo.SearchByTitle(title)
	if err != nil {
		return nil, fmt.Errorf("searching by the title error: %w", err)
	}

	return res, nil
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

// HandleBooks processes all book-related endpoints.
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/books/search" {
		h.HandleSearch(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/books/") {
		h.HandleBook(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		books, err := h.Service.GetAllBooks()
		if err != nil {
			sendError(w, &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(books)

	case http.MethodPost:
		var book Book

		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			sendError(w, &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "invalid request body",
			})
			return
		}

		if book.Title == "" ||
			book.Author == "" ||
			book.ISBN == "" ||
			book.PublishedYear <= 0 {

			sendError(w, &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "missing required field",
			})
			return
		}

		if err := h.Service.CreateBook(&book); err != nil {
			sendError(w, &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(book)

	default:
		sendError(w, &ErrorResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "method not allowed",
		})
	}
}

// HandleBook processes:
// GET /api/books/{id}
// PUT /api/books/{id}
// DELETE /api/books/{id}
func (h *BookHandler) HandleBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/books/")

	if id == "" {
		sendError(w, &ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "invalid book id",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		book, err := h.Service.GetBookByID(id)
		if err != nil {
			if errors.Is(err, ErrInMemoRepoNotFound) {
				sendError(w, &ErrorResponse{
					StatusCode: http.StatusNotFound,
					Error:      err.Error(),
				})
			} else {
				sendError(w, &ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(book)

	case http.MethodPut:
		var book Book

		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			sendError(w, &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      "invalid request body",
			})
			return
		}

		book.ID = id

		if err := h.Service.UpdateBook(id, &book); err != nil {
			if errors.Is(err, ErrInMemoRepoNotFound) {
				sendError(w, &ErrorResponse{
					StatusCode: http.StatusNotFound,
					Error:      err.Error(),
				})
			} else {
				sendError(w, &ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(book)

	case http.MethodDelete:
		if err := h.Service.DeleteBook(id); err != nil {
			if errors.Is(err, ErrInMemoRepoNotFound) {
				sendError(w, &ErrorResponse{
					StatusCode: http.StatusNotFound,
					Error:      err.Error(),
				})
			} else {
				sendError(w, &ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Error:      err.Error(),
				})
			}
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		sendError(w, &ErrorResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "method not allowed",
		})
	}
}

// HandleSearch processes:
// GET /api/books/search?author={author}
// GET /api/books/search?title={title}
func (h *BookHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, &ErrorResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "method not allowed",
		})
		return
	}

	author := r.URL.Query().Get("author")
	title := r.URL.Query().Get("title")

	if author == "" && title == "" {
		sendError(w, &ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "author or title query parameter is required",
		})
		return
	}

	if author != "" && title != "" {
		sendError(w, &ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "use either author or title query parameter",
		})
		return
	}

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
		sendError(w, &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(books)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

func sendError(w http.ResponseWriter, errResp *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.StatusCode)
	_ = json.NewEncoder(w).Encode(errResp)
}

func main() {
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	http.HandleFunc("/api/books", handler.HandleBooks)
	http.HandleFunc("/api/books/", handler.HandleBooks)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
