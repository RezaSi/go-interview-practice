package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
	ISBN          string `json:"isbn"`
	Description   string `json:"description"`
}
type BookRepository interface {
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Create(book *Book) error
	Update(id string, book *Book) error
	Delete(id string) error
	SearchByAuthor(author string) ([]*Book, error)
	SearchByTitle(title string) ([]*Book, error)
}
type InMemoryBookRepository struct {
	books map[string]*Book
	mu    sync.RWMutex
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{books: make(map[string]*Book)}
}

var (
	ErrBookRepositoryIdNotFound = errors.New("no book with this ID was found")
	ErrBookRepositoryCantCreate = errors.New("book is invalid, cannot create book")
	ErrInvalidJSON              = errors.New("invalid JSON")
)

func validateBook(book *Book) error {
	if book == nil {
		return fmt.Errorf("%w: book payload is required", ErrBookRepositoryCantCreate)
	}
	if book.Title == "" {
		return fmt.Errorf("%w: title is empty", ErrBookRepositoryCantCreate)
	}
	if book.Author == "" {
		return fmt.Errorf("%w: author is empty", ErrBookRepositoryCantCreate)
	}
	if book.PublishedYear > time.Now().Year() {
		return fmt.Errorf("%w: published year cannot be too far in the future", ErrBookRepositoryCantCreate)
	}
	if book.PublishedYear <= 0 {
		return fmt.Errorf("%w: published year must be positive", ErrBookRepositoryCantCreate)
	}
	if book.ISBN != "" {
		normalizedISBN := strings.ReplaceAll(book.ISBN, "-", "")
		isbnLen := len(normalizedISBN)
		if isbnLen != 10 && isbnLen != 13 {
			return fmt.Errorf("%w: ISBN must be exactly 10 or 13 characters", ErrBookRepositoryCantCreate)
		}
		for i, ch := range normalizedISBN {
			// ISBN-10 can have 'X' as the last character (checksum = 10)
			if ch == 'X' || ch == 'x' {
				if isbnLen == 10 && i == 9 {
					continue
				}
				return fmt.Errorf("%w: 'X' is only valid as the last character of ISBN-10", ErrBookRepositoryCantCreate)
			}
			if ch < '0' || ch > '9' {
				return fmt.Errorf("%w: ISBN must contain only digits (optionally separated by dashes)", ErrBookRepositoryCantCreate)
			}
		}
	}
	return nil
}
func (d *InMemoryBookRepository) GetAll() ([]*Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	books := make([]*Book, 0, len(d.books))
	for _, v := range d.books {
		copy := *v
		books = append(books, &copy)
	}
	return books, nil
}

func (d *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if book, exists := d.books[id]; exists {
		copy := *book
		return &copy, nil
	}
	return nil, ErrBookRepositoryIdNotFound
}

func (d *InMemoryBookRepository) Create(book *Book) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := validateBook(book); err != nil {
		return err
	}
	id := uuid.New().String()
	book.ID = id
	storedBook := *book
	d.books[id] = &storedBook
	return nil
}

func (d *InMemoryBookRepository) Update(id string, book *Book) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := validateBook(book); err != nil {
		return err
	}
	_, exists := d.books[id]
	if !exists {
		return ErrBookRepositoryIdNotFound
	}
	book.ID = id
	storedBook := *book
	d.books[id] = &storedBook
	return nil
}

func (d *InMemoryBookRepository) Delete(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, exists := d.books[id]
	if !exists {
		return ErrBookRepositoryIdNotFound
	}
	delete(d.books, id)
	return nil
}

func (d *InMemoryBookRepository) SearchBy(predicate func(*Book) bool) ([]*Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	books := make([]*Book, 0)
	for _, book := range d.books {
		if predicate(book) {
			copy := *book
			books = append(books, &copy)
		}
	}
	return books, nil
}

func (d *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	return d.SearchBy(func(book *Book) bool {
		return strings.Contains(strings.ToLower(book.Author), strings.ToLower(author))
	})
}

func (d *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	return d.SearchBy(func(book *Book) bool {
		return strings.Contains(strings.ToLower(book.Title), strings.ToLower(title))
	})
}

type BookService interface {
	GetAllBooks() ([]*Book, error)
	GetBookByID(id string) (*Book, error)
	CreateBook(book *Book) error
	UpdateBook(id string, book *Book) error
	DeleteBook(id string) error
	SearchBooksByAuthor(author string) ([]*Book, error)
	SearchBooksByTitle(title string) ([]*Book, error)
}

type DefaultBookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *DefaultBookService {
	return &DefaultBookService{repo: repo}
}
func (d *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return d.repo.GetAll()
}
func (d *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return d.repo.GetByID(id)
}
func (d *DefaultBookService) CreateBook(book *Book) error {
	return d.repo.Create(book)
}
func (d *DefaultBookService) UpdateBook(id string, book *Book) error {
	return d.repo.Update(id, book)
}
func (d *DefaultBookService) DeleteBook(id string) error {
	return d.repo.Delete(id)
}
func (d *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	return d.repo.SearchByAuthor(author)
}
func (d *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	return d.repo.SearchByTitle(title)
}
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	encoded, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(encoded); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := ErrorResponse{Error: message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func (h *BookHandler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.GetAllBooks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1 MB limit
	var book Book
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON.Error())
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON.Error())
		return
	}
	if err := h.Service.CreateBook(&book); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1 MB limit
	var book Book
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON.Error())
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON.Error())
		return
	}
	if err := h.Service.UpdateBook(id, &book); err != nil {
		if errors.Is(err, ErrBookRepositoryIdNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.Service.DeleteBook(id); err != nil {
		if errors.Is(err, ErrBookRepositoryIdNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Book deleted successfully"})
}

func (h *BookHandler) getBookById(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		if errors.Is(err, ErrBookRepositoryIdNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) searchBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	title := r.URL.Query().Get("title")
	switch {
	case author != "":
		books, err := h.Service.SearchBooksByAuthor(author)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, books)
	case title != "":
		books, err := h.Service.SearchBooksByTitle(title)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, books)
	default:
		writeError(w, http.StatusBadRequest, "Missing search parameter: author or title")
	}
}

type BookHandler struct {
	Service BookService
}

func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{
		Service: service,
	}
}

func getIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) != 3 {
		return ""
	}

	return parts[2]
}

// this function signature is part of the assignment signature and cannot be deleted
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	if path == "" {
		path = "/"
	}

	switch {
	case path == "/api/books":
		switch r.Method {
		case http.MethodGet:
			h.getAllBooks(w, r)
		case http.MethodPost:
			h.createBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case path == "/api/books/search":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.searchBooks(w, r)

	case strings.HasPrefix(path, "/api/books/"):
		switch r.Method {
		case http.MethodGet:
			h.getBookById(w, r)

		case http.MethodPut:
			h.updateBook(w, r)

		case http.MethodDelete:
			h.deleteBook(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	default:
		http.NotFound(w, r)
	}
}

// cannot use mux for this assignment
func main() {
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/books", handler.HandleBooks)
	mux.HandleFunc("/api/books/", handler.HandleBooks)

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, map[string]string{
				"status": "ok",
			})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Fatal(http.ListenAndServe(":8083", mux))
}