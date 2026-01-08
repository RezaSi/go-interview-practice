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
	"github.com/gorilla/mux"
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
	ErrBookRepositoryIdNotFound  = errors.New("no book with this ID was found")
	ErrBookRepositoryCantCreate  = errors.New("book is invalid, cannot create book")
	ErrStatusInternalServerError = errors.New("Internal server error")
	ErrInvalidJSON               = errors.New("Invalid JSON")
)

func validateBook(book *Book) error {
	if book.Title == "" {
		return fmt.Errorf("%w: title is empty", ErrBookRepositoryCantCreate)
	}
	if book.Author == "" {
		return fmt.Errorf("%w: author is empty", ErrBookRepositoryCantCreate)
	}
	if book.PublishedYear <= 0 {
		return fmt.Errorf("%w: published year must be positive", ErrBookRepositoryCantCreate)
	}
	// if book.ISBN != "" {
	// 	isbnLen := len(book.ISBN)
	// 	if isbnLen != 10 && isbnLen != 13 {
	// 		return fmt.Errorf("%w: ISBN must be exactly 10 or 13 characters", ErrBookRepositoryCantCreate)
	// 	}
	// }
	return nil
}
func (d *InMemoryBookRepository) GetAll() ([]*Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var books []*Book
	for _, v := range d.books {
		books = append(books, v)
	}
	return books, nil
}

func (d *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if book, exists := d.books[id]; exists {
		return book, nil
	}
	return nil, ErrBookRepositoryIdNotFound
}

func (d *InMemoryBookRepository) Create(book *Book) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := validateBook(book); err != nil {
		return err
	}
	book.ID = uuid.New().String()
	d.books[book.ID] = book
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
	d.books[id] = book
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
			books = append(books, book)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
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
	vars := mux.Vars(r)
	id := vars["id"]
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON.Error())
		return
	}
	if err := h.Service.UpdateBook(id, &book); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.Service.DeleteBook(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Book deleted successfully"})
}

func (h *BookHandler) searchBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	title := r.URL.Query().Get("title")
	vars := mux.Vars(r)
	id := vars["id"]
	switch {
	case author != "":
		h.searchBooksByAuthor(w, r, author)
	case title != "":
		h.searchBooksByTitle(w, r, title)
	case id != "":
		h.searchBooksById(w, r, id)
	default:
		writeError(w, http.StatusBadRequest, "Missing search parameter: id, author, title")
	}
}

func (h *BookHandler) searchBooksById(w http.ResponseWriter, r *http.Request, id string) {
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) searchBooksByAuthor(w http.ResponseWriter, r *http.Request, author string) {
	books, err := h.Service.SearchBooksByAuthor(author)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) searchBooksByTitle(w http.ResponseWriter, r *http.Request, title string) {
	books, err := h.Service.SearchBooksByTitle(title)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, books)
}

type BookHandler struct {
	Service BookService
}

func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{Service: service}
}

// this function signature is part of the assignment signature and cannot be deleted
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.HandleFunc("/api/books", h.getAllBooks).Methods("GET")
	router.HandleFunc("/api/books", h.createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", h.searchBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", h.deleteBook).Methods("DELETE")
	router.ServeHTTP(w, r)
}

func main() {
 	repo := NewInMemoryBookRepository()
 	service := NewBookService(repo)
 	handler := NewBookHandler(service)
 
	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", handler.HandleBooks)
	mux.HandleFunc("/api/books/", handler.HandleBooks)

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {

 		if r.Method == "GET" {
 			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
 		}
 	})
	if err := http.ListenAndServe(":8081", mux); err != nil {
 		log.Fatalf("Failed to start server: %v", err)
 	}
 }




