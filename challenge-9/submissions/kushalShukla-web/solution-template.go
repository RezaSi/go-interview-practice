package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Book struct {
	ID            string
	Title         string
	Author        string
	PublishedYear int
	ISBN          string
	Description   string
}

type InMemoryBookRepository struct {
	books map[string]*Book
}
type DefaultBookService struct {
	repo BookRepository
}

// BookHandler handles HTTP requests for book operations
type BookHandler struct {
	Service BookService
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

type BookRepository interface {
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Create(book *Book) error
	Update(id string, book *Book) error
	Delete(id string) error
	SearchByAuthor(Author string) ([]*Book, error)
	SearchByTitle(title string) ([]*Book, error)
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books: make(map[string]*Book),
	}
}

func NewBookService(repo BookRepository) *DefaultBookService {
	return &DefaultBookService{
		repo: repo,
	}
}

func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{
		Service: service,
	}
}

func (x *InMemoryBookRepository) GetAll() ([]*Book, error) {
	var books []*Book
	for _, value := range x.books {
		books = append(books, value)
	}
	return books, nil
}

func (x *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	if x.books[id] == nil {
		return nil, fmt.Errorf("There is no ID")
	}
	return x.books[id], nil
}

func (x *InMemoryBookRepository) Delete(id string) error {
	if _, ok := x.books[id]; !ok {
		return fmt.Errorf("book not found: %s", id)
	}
	delete(x.books, id)
	return nil
}
func (x *InMemoryBookRepository) Create(book *Book) error {
	if book == nil {
		return fmt.Errorf("Error while creating a book")
	}
	book.ID = strconv.Itoa(len(x.books) + 1)
	x.books[strconv.Itoa((len(x.books) + 1))] = book

	return nil
}
func (x *InMemoryBookRepository) Update(id string, book *Book) error {
	_, ok := x.books[id]
	if !ok {
		return fmt.Errorf("book not found")
	}
	book.ID = id
	x.books[id] = book
	return nil
}
func (x *InMemoryBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	var res []*Book
	if strings.TrimSpace(author) == "" {
		return res, nil
	}
	a := strings.ToLower(author)
	for _, v := range x.books {
		if strings.Contains(strings.ToLower(v.Author), a) {
			res = append(res, v)
		}
	}
	return res, nil // return empty slice (not error) when none found
}

func (x *InMemoryBookRepository) SearchByTitle(title string) ([]*Book, error) {
	var res []*Book
	if strings.TrimSpace(title) == "" {
		return res, nil
	}
	t := strings.ToLower(title)
	for _, v := range x.books {
		if strings.Contains(strings.ToLower(v.Title), t) {
			res = append(res, v)
		}
	}
	return res, nil
}

func (x *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return x.repo.GetAll()
}
func (x *DefaultBookService) GetBookByID(id string) (*Book, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id empty")
	}
	return x.repo.GetByID(id)
}
func (x *DefaultBookService) CreateBook(book *Book) error {
	if book == nil {
		return errors.New("book is nil")
	}
	if strings.TrimSpace(book.Title) == "" {
		return errors.New("title required")
	}
	if strings.TrimSpace(book.Author) == "" {
		return errors.New("author required")
	}
	return x.repo.Create(book)
}
func (x *DefaultBookService) UpdateBook(id string, book *Book) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id empty")
	}
	if book == nil {
		return errors.New("book is nil")
	}
	if strings.TrimSpace(book.Title) == "" {
		return errors.New("title required")
	}
	if strings.TrimSpace(book.Author) == "" {
		return errors.New("author required")
	}
	return x.repo.Update(id, book)
}

func (x *DefaultBookService) DeleteBook(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id empty")
	}
	return x.repo.Delete(id)
}

func (x *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	if strings.TrimSpace(author) == "" {
		return nil, errors.New("author empty")
	}
	return x.repo.SearchByAuthor(author)
}

func (x *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title empty")
	}
	return x.repo.SearchByTitle(title)
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/api/books":
		books, err := h.Service.GetAllBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		if len(books) == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		json.NewEncoder(w).Encode(books)
		return
	case r.Method == http.MethodPost && r.URL.Path == "/api/books":
		var b Book
		json.NewDecoder(r.Body).Decode(&b)
		err := h.Service.CreateBook(&b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		json.NewEncoder(w).Encode(&b)
		return
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/books/search"):
		author := r.URL.Query().Get("author")
		title := r.URL.Query().Get("title")

		if author != "" {
			res, err := h.Service.SearchBooksByAuthor(author)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return
		}

		if title != "" {
			res, err := h.Service.SearchBooksByTitle(title)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "no search query"})
		return

	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/books/"):
		id := strings.TrimPrefix(r.URL.Path, "/api/books/")
		book, err := h.Service.GetBookByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		}
		json.NewEncoder(w).Encode(&book)
		return
	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/books/"):
		id := strings.TrimPrefix(r.URL.Path, "/api/books/")
		var b Book
		b.ID = id
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		err := h.Service.UpdateBook(id, &b)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		updated, err := h.Service.GetBookByID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(updated)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/books/"):
		id := strings.TrimPrefix(r.URL.Path, "/api/books/")
		err := h.Service.DeleteBook(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Book deleted successfully",
		})
		return
	}
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
