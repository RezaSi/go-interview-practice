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

	"github.com/google/uuid"
)

// Custom errors
var (
	ErrNotFound   = errors.New("Not found")
	ErrBadParam   = errors.New("Bad parameter")
	ErrNotAllowed = errors.New("Not allowed")
)

// Book represents a book in the database
// Repository level Model
type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
	ISBN          string `json:"isbn"`
	Description   string `json:"description"`
}

// BookRepository defines the operations for book data access
// Service leyer interface. Repository depends on
type BookRepository interface {
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Create(book *Book) error
	Update(id string, book *Book) error
	Delete(id string) error
	SearchByAuthor(author string) ([]*Book, error)
	SearchByTitle(title string) ([]*Book, error)
}

// InMemoryBookRepository implements BookRepository interface using in-memory storage
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

// Implement BookRepository interface methods for InMemoryBookRepository entity
// GetAll returns slice of pointers on Book or error
func (br *InMemoryBookRepository) GetAll() ([]*Book, error) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	books := make([]*Book, 0, len(br.books))
	for _, b := range br.books {
		books = append(books, b)
	}

	return books, nil
}

// GetByID returns pointer on Book or error
func (br *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	// Parameter validation
	if id == "" {

		return nil, fmt.Errorf("Id parameter: %w", ErrBadParam)
	}

	br.mu.RLock()
	defer br.mu.RUnlock()

	book, ok := br.books[id]
	if !ok {

		return nil, fmt.Errorf("Get book by id: %w", ErrNotFound)
	}

	return book, nil
}

// Create creates book in repository or returns error
func (br *InMemoryBookRepository) Create(book *Book) error {
	// Parametr validation
	if book == nil {

		return fmt.Errorf("Create parametr book: %w", ErrBadParam)
	}
	if book.Title == "" {

		return fmt.Errorf("Create parameter book.Title: %w", ErrBadParam)
	}

	// Get new UUID for new book
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("UUID: %w", err)
	}

	// Set new UUID for new book
	book.ID = id.String()

	// Store new book with new UUID under mutex
	br.mu.Lock()
	defer br.mu.Unlock()

	br.books[book.ID] = book

	return nil
}

// Update updates existing book in repository or returns error
func (br *InMemoryBookRepository) Update(id string, book *Book) error {
	// Parameter validation
	if id == "" {

		return fmt.Errorf("Update id: %w", ErrBadParam)
	}
	if book == nil {

		return fmt.Errorf("Update book: %w", ErrBadParam)
	}

	br.mu.Lock()
	defer br.mu.Unlock()

	oldBook, ok := br.books[id]
	if !ok {

		return fmt.Errorf("Update book: %w", ErrNotFound)
	}
	// Update fields if not blanc
	if book.Title != "" {
		oldBook.Title = book.Title
	}
	if book.Author != "" {
		oldBook.Author = book.Author
	}
	if book.PublishedYear != 0 {
		oldBook.PublishedYear = book.PublishedYear
	}
	if book.ISBN != "" {
		oldBook.ISBN = book.ISBN
	}
	if book.Description != "" {
		oldBook.Description = book.Description
	}

	br.books[id] = oldBook

	return nil
}

// Delete delets book from repository or returns error
func (br *InMemoryBookRepository) Delete(id string) error {
	// Parameter validation
	if id == "" {

		return fmt.Errorf("Delete id: %w", ErrBadParam)
	}
	br.mu.Lock()
	defer br.mu.Unlock()

	if _, ok := br.books[id]; !ok {
		return fmt.Errorf("Delete book: %w", ErrNotFound)
	}
	delete(br.books, id)

	return nil
}

// SearchByAuthor returns slce of pointers of Book or error
func (br *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	// Parameter validation
	if author == "" {
		return nil, fmt.Errorf("SearchByAuthor author: %w", ErrBadParam)
	}

	var res []*Book
	br.mu.RLock()
	defer br.mu.RUnlock()

	for _, v := range br.books {
		if strings.Contains(v.Author, author) {
			res = append(res, v)
		}
	}

	return res, nil
}

// SearchByTitle returns slice of pointers of Book or error
func (br *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	// Parameter validation
	if title == "" {

		return nil, fmt.Errorf("SearchByTitle title: %w", ErrBadParam)
	}

	var res []*Book

	br.mu.RLock()
	defer br.mu.RUnlock()

	for _, v := range br.books {
		if strings.Contains(v.Title, title) {
			res = append(res, v)
		}
	}

	return res, nil
}

// BookService defines the business logic for book operations
// Transport layer interface. Service depends on
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
// Service layer
type DefaultBookService struct {
	repo BookRepository
}

// NewBookService creates a new book service
func NewBookService(repo BookRepository) *DefaultBookService {

	return &DefaultBookService{
		repo: repo,
	}
}

// There is no book domain struct in this implimentation of the Service layer
// Implement BookService interface methods for DefaultBookService entity
// Service layer
// GetAllBooks returns slice of all books domains or error
func (bs *DefaultBookService) GetAllBooks() ([]*Book, error) {
	books, err := bs.repo.GetAll()
	if err != nil {

		return nil, fmt.Errorf("GetAllBooks: %w", err)
	}

	return books, nil
}

// GetBookByID returns pointer on Book or error
func (bs *DefaultBookService) GetBookByID(id string) (*Book, error) {
	// parameter validation
	if id == "" {

		return nil, fmt.Errorf("GetBookByID id: %w", ErrBadParam)
	}
	book, err := bs.repo.GetByID(id)
	if err != nil {

		return nil, fmt.Errorf("GetBookByID: %w", err)
	}

	return book, nil
}

// CreateBook creates Book domain and send it in reporitory layer to save in repository
func (bs *DefaultBookService) CreateBook(book *Book) error {
	// Parameter validation
	if book == nil {

		return fmt.Errorf("CreateBook book: %w", ErrBadParam)
	}

	if err := bs.repo.Create(book); err != nil {

		return fmt.Errorf("CreateBook: %w", err)
	}

	return nil
}

// UpdateBook creates Book domain for updating and send in repository layer to update
func (bs *DefaultBookService) UpdateBook(id string, book *Book) error {
	// Parameter validaion
	if id == "" {

		return fmt.Errorf("UpdateBook id: %w", ErrBadParam)
	}
	if book == nil {

		return fmt.Errorf("UpdateBook book: %w", ErrBadParam)
	}

	if err := bs.repo.Update(id, book); err != nil {
		return fmt.Errorf("UpdateBook: %w", err)
	}

	return nil
}

// DeleteBook send id to repository layer to delete book from repository
func (bs *DefaultBookService) DeleteBook(id string) error {
	// Parameter validation
	if id == "" {

		return fmt.Errorf("DeleteBook id: %w", ErrBadParam)
	}

	if err := bs.repo.Delete(id); err != nil {

		return fmt.Errorf("DeleteBook: %w", err)
	}

	return nil
}

// SearchBooksByAuthor send author to repository layer and returns slice of
// pointers on book or error
func (bs *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	// Parameter validation
	if author == "" {
		return nil, fmt.Errorf("SearchBooksByAuthor author: %w", ErrBadParam)
	}
	books, err := bs.repo.SearchByAuthor(author)
	if err != nil {

		return nil, fmt.Errorf("SearchBooksByAuthor: %w", err)
	}

	return books, nil
}

// SearchBooksByTitle send title to repository layer and returns slice of
// pointers on book or error
func (bs *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	// Parameter validation
	if title == "" {

		return nil, fmt.Errorf("SearchBooksByTitle title: %w", ErrBadParam)
	}
	books, err := bs.repo.SearchByTitle(title)
	if err != nil {

		return nil, fmt.Errorf("SearchBooksByTitle: %w", err)
	}

	return books, nil
}

// BookHandler handles HTTP requests for book operations
// Transport layer
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
// Transport layer
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	// Pars path determin is there id or query in it
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	hasID := len(parts) == 4 && parts[3] != "" && parts[3] != "search"
	hasQuery := len(parts) == 4 && parts[3] == "search"
	// for Requests with id
	if hasID {
		id := parts[3]
		switch r.Method {
		case http.MethodGet:

			book, err := h.Service.GetBookByID(id)
			if err != nil {
				writeError(w, fmt.Errorf("HandleBooks: %w", err))

				return
			}
			writeJson(w, http.StatusOK, book)

			return

		case http.MethodPut:
			var book Book
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
				writeError(w, fmt.Errorf("HandleBooks json decode: %w", err))

				return
			}
			if err := h.Service.UpdateBook(id, &book); err != nil {
				writeError(w, fmt.Errorf("UpdateBook: %w", err))

				return
			}
			writeJson(w, http.StatusOK, book)

		case http.MethodDelete:
			if err := h.Service.DeleteBook(id); err != nil {
				writeError(w, fmt.Errorf("DeleteBook: %w", err))

				return
			}
			writeJson(w, http.StatusOK, nil)
		default:
			writeError(w, fmt.Errorf("HTTP with param method: %w", ErrNotAllowed))

			return
		}
	}
	// for requests with query
	if hasQuery {
		//search implimentation
		switch r.Method {
		case http.MethodGet:
			query := r.URL.Query()
			author := query["author"]
			title := query["title"]
			if len(author) > 0 {
				books, err := h.Service.SearchBooksByAuthor(author[0])
				if err != nil {
					writeError(w, fmt.Errorf("SearchBooksByAuthor: %w", err))
				}
				writeJson(w, http.StatusOK, books)
			}
			if len(title) > 0 {
				books, err := h.Service.SearchBooksByTitle(title[0])
				if err != nil {
					writeError(w, fmt.Errorf("SearchBooksByTitle: %w", err))
				}
				writeJson(w, http.StatusOK, books)
			}

		default:
			writeError(w, fmt.Errorf("HTTP with param method: %w", ErrNotAllowed))
		}
	}

	// For no id no query requests
	switch r.Method {
	case http.MethodGet:
		books, err := h.Service.GetAllBooks()
		if err != nil {
			writeError(w, fmt.Errorf("GetAllBooks: %w", err))

			return
		}

		writeJson(w, http.StatusOK, books)

	case http.MethodPost:
		var book Book

		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			writeError(w, err)
		}
		defer r.Body.Close()

		if err := h.Service.CreateBook(&book); err != nil {
			writeError(w, err)
		}
		writeJson(w, http.StatusCreated, book)

	default:
		writeError(w, fmt.Errorf("HTTP without param method: %w", ErrNotAllowed))
	}

}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Error      string `json:"error"`
}

// newErrorResponse returns new ErrorResponse struct
func newErrorResponse(s int, e string) ErrorResponse {
	return ErrorResponse{
		StatusCode: s,
		Error:      e,
	}
}

// Helper functions
// writeError writes RespontseWriter with http status that determined on custom error and error itself
func writeError(w http.ResponseWriter, err error) {
	var resp ErrorResponse
	switch {
	case errors.Is(err, ErrBadParam):
		resp = newErrorResponse(http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		resp = newErrorResponse(http.StatusNotFound, err.Error())
	case errors.Is(err, ErrNotAllowed):
		resp = newErrorResponse(http.StatusMethodNotAllowed, err.Error())
	default:
		resp = newErrorResponse(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp)
}

// writeJson writes ResponseWriter with encoded json response
func writeJson(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Cotnent-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		writeError(w, err)
	}
}

func main() {
	// Initialize the repository, service, and handler
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	// Create a new router and register endpoints
	http.HandleFunc("/api/books", handler.HandleBooks)
	http.HandleFunc("/api/books/{id}", handler.HandleBooks)
	http.HandleFunc("/api/books/{id}", handler.HandleBooks)

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
