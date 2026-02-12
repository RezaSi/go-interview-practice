// Package main contains the implementation for Challenge 9: RESTful Book Management API
package main

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

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

// BookRepository
// =======================================================================

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
	books   map[string]*Book
	idIndex map[string]string // map[bookID]hashKey
	mu      sync.RWMutex
	cnt     uint64
}

// NewInMemoryBookRepository creates a new in-memory book repository
func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books:   make(map[string]*Book),
		idIndex: make(map[string]string),
		cnt:     1,
	}
}

var (
	ErrBookRepositoryEmpty      = errors.New("not a single book was found")
	ErrBookRepositoryIdNotFound = errors.New("no book with this ID was found")
	ErrBookRepositoryCantCreate = errors.New("book is invalid, cannot create book")
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

	if book.ISBN != "" {
		if len(book.ISBN) < 10 {
			return fmt.Errorf("%w: ISBN too short", ErrBookRepositoryCantCreate)
		}
	}

	return nil
}

func createHashByBook(book *Book) string {
	title := strings.ToLower(strings.TrimSpace(book.Title))
	author := strings.ToLower(strings.TrimSpace(book.Author))
	isbn := strings.ToLower(strings.TrimSpace(book.ISBN))

	input := fmt.Sprintf("%s|%s|%d|%s",
		title,
		author,
		book.PublishedYear,
		isbn,
	)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (d *InMemoryBookRepository) GetAll() ([]*Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var books []*Book

	/*if len(d.books) == 0 {
		return books, ErrBookRepositoryEmpty
	}*/

	for _, v := range d.books {
		books = append(books, v)
	}

	return books, nil
}

func (d *InMemoryBookRepository) GetByID(id string) (*Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if hashKey, exists := d.idIndex[id]; exists {
		return d.books[hashKey], nil
	}
	return nil, ErrBookRepositoryIdNotFound
}

func (d *InMemoryBookRepository) Create(book *Book) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := validateBook(book); err != nil {
		return err
	}

	book.ID = fmt.Sprintf("%d", d.cnt)
	hashBook := createHashByBook(book)

	if _, exists := d.books[hashBook]; exists {
		return fmt.Errorf("%w: there is a similar book", ErrBookRepositoryCantCreate)
	}

	d.books[hashBook] = book
	d.idIndex[book.ID] = hashBook
	d.cnt++

	return nil
}

func (d *InMemoryBookRepository) Update(id string, book *Book) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := validateBook(book); err != nil {
		return err
	}

	oldHash, exists := d.idIndex[id]
	if !exists {
		return ErrBookRepositoryIdNotFound
	}

	book.ID = id
	newHash := createHashByBook(book)

	if oldHash != newHash {
		if _, exists := d.books[newHash]; exists {
			return fmt.Errorf("%w: there is a similar book", ErrBookRepositoryCantCreate)
		}
	}

	delete(d.books, oldHash)
	d.books[newHash] = book
	d.idIndex[id] = newHash

	return nil
}

func (d *InMemoryBookRepository) Delete(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	hashKey, exists := d.idIndex[id]
	if !exists {
		return ErrBookRepositoryIdNotFound
	}

	delete(d.books, hashKey)
	delete(d.idIndex, id)

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

	/*if len(books) == 0 {
		return nil, ErrBookRepositoryEmpty
	}*/

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

// BookService
// =======================================================================

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

// BookHandler
// =======================================================================

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

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// LoggingMiddleware логирует запросы, ответы и ошибки
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		wrapped := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           new(bytes.Buffer),
		}

		// Обрабатываем запрос
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Логируем всю информацию
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		//log.Printf("[HEADERS] %v", r.Header)

		if len(requestBody) > 0 && len(requestBody) < 1024 { // Логируем body только если он не слишком большой
			log.Printf("[BODY] %s", string(requestBody))
		}

		log.Printf("[RESPONSE] Status: %d, Time: %v", wrapped.statusCode, duration)

		if wrapped.statusCode >= 400 {
			responseBody := wrapped.body.String()
			if len(responseBody) > 0 {
				var errorResp ErrorResponse
				if err := json.Unmarshal([]byte(responseBody), &errorResp); err == nil {
					log.Printf("[ERROR] %s", errorResp.Error)
				} else {
					log.Printf("[ERROR_RESPONSE] %s", responseBody)
				}
			}
			log.Printf("[ERROR_DETAILS] Method: %s, Path: %s, Status: %d",
				r.Method, r.URL.Path, wrapped.statusCode)
		}

		log.Printf("----------------------------------------")
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(errorResponse)
}

var (
	ErrStatusInternalServerError = errors.New("Internal server error")
	ErrInvalidJSON               = errors.New("Invalid JSON")
)

// HandleBooks processes the book-related endpoints
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.Use(LoggingMiddleware)

	router.HandleFunc("/api/books", h.getAllBooks).Methods("GET")
	router.HandleFunc("/api/books", h.createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", h.searchBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", h.deleteBook).Methods("DELETE")

	router.ServeHTTP(w, r)
}

// getAllBooks handles GET /api/books
func (h *BookHandler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.GetAllBooks()
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, books)
}

// createBook handles POST /api/books
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

// updateBook handles PUT /api/books/{id}
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

// deleteBook handles DELETE /api/books/{id}
func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.Service.DeleteBook(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Book deleted successfully"})
}

// searchBooks handles GET /api/books/search, handles GET /api/books/{id}
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

// searchBooksById handles search by author
func (h *BookHandler) searchBooksById(w http.ResponseWriter, r *http.Request, id string) {
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, book)
}

// searchBooksByAuthor handles search by author
func (h *BookHandler) searchBooksByAuthor(w http.ResponseWriter, r *http.Request, author string) {
	books, err := h.Service.SearchBooksByAuthor(author)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, books)
}

// searchBooksByTitle handles search by title
func (h *BookHandler) searchBooksByTitle(w http.ResponseWriter, r *http.Request, title string) {
	books, err := h.Service.SearchBooksByTitle(title)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, books)
}

// HTMLInterface
// =======================================================================
// serveHTMLInterface serves the HTML testing interface
func serveHTMLInterface(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Book API Tester</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: Arial, sans-serif; line-height: 1.6; padding: 20px; background: #f4f4f4; }
        .container { max-width: 1200px; margin: 0 auto; }
        h1 { text-align: center; margin-bottom: 30px; color: #333; }
        .section { background: white; padding: 20px; margin-bottom: 20px; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        h2 { color: #2c3e50; margin-bottom: 15px; border-bottom: 2px solid #eee; padding-bottom: 10px; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; color: #555; }
        input, textarea, select { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
        textarea { height: 80px; resize: vertical; }
        button { background: #3498db; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; }
        button:hover { background: #2980b9; }
        .response { margin-top: 15px; padding: 15px; background: #f8f9fa; border: 1px solid #ddd; border-radius: 4px; max-height: 300px; overflow-y: auto; }
        pre { white-space: pre-wrap; word-wrap: break-word; }
        .success { color: #27ae60; }
        .error { color: #e74c3c; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        @media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
    </style>
</head>
<body>
    <div class="container">
        <h1>📚 Book API Tester</h1>
        
        <div class="grid">
            <!-- Create Book -->
            <div class="section">
                <h2>➕ Create Book</h2>
                <form id="createForm">
                    <div class="form-group">
                        <label>Title:</label>
                        <input type="text" name="title" value="The Go Programming Language" required>
                    </div>
                    <div class="form-group">
                        <label>Author:</label>
                        <input type="text" name="author" value="Alan A. A. Donovan and Brian W. Kernighan" required>
                    </div>
                    <div class="form-group">
                        <label>Published Year:</label>
                        <input type="number" name="published_year" value="2015" required>
                    </div>
                    <div class="form-group">
                        <label>ISBN:</label>
                        <input type="text" name="isbn" value="978-0134190440">
                    </div>
                    <div class="form-group">
                        <label>Description:</label>
                        <textarea name="description">The definitive guide to programming in Go</textarea>
                    </div>
                    <button type="submit">Create Book</button>
                </form>
                <div class="response" id="createResponse"></div>
            </div>

            <!-- Get All Books -->
            <div class="section">
                <h2>📋 Get All Books</h2>
                <button onclick="getAllBooks()">Get All Books</button>
                <div class="response" id="getAllResponse"></div>
            </div>

            <!-- Get Book by ID -->
            <div class="section">
                <h2>🔍 Get Book by ID</h2>
                <div class="form-group">
                    <label>Book ID:</label>
                    <input type="text" id="getById" placeholder="Enter book ID" value="1">
                </div>
                <button onclick="getBookById()">Get Book</button>
                <div class="response" id="getByIdResponse"></div>
            </div>

            <!-- Update Book -->
            <div class="section">
                <h2>✏️ Update Book</h2>
                <form id="updateForm">
                    <div class="form-group">
                        <label>Book ID:</label>
                        <input type="text" name="id" required>
                    </div>
                    <div class="form-group">
                        <label>Title:</label>
                        <input type="text" name="title" required>
                    </div>
                    <div class="form-group">
                        <label>Author:</label>
                        <input type="text" name="author" required>
                    </div>
                    <div class="form-group">
                        <label>Published Year:</label>
                        <input type="number" name="published_year" required>
                    </div>
                    <div class="form-group">
                        <label>ISBN:</label>
                        <input type="text" name="isbn">
                    </div>
                    <div class="form-group">
                        <label>Description:</label>
                        <textarea name="description"></textarea>
                    </div>
                    <button type="submit">Update Book</button>
                </form>
                <div class="response" id="updateResponse"></div>
            </div>

            <!-- Delete Book -->
            <div class="section">
                <h2>🗑️ Delete Book</h2>
                <div class="form-group">
                    <label>Book ID:</label>
                    <input type="text" id="deleteId" placeholder="Enter book ID">
                </div>
                <button onclick="deleteBook()">Delete Book</button>
                <div class="response" id="deleteResponse"></div>
            </div>

            <!-- Search Books -->
            <div class="section">
                <h2>🔎 Search Books</h2>
                <div class="form-group">
                    <label>Search by Author:</label>
                    <input type="text" id="searchAuthor" placeholder="Enter author name" value="Alan A. A. Donovan and Brian W. Kernighan">
                </div>
                <button onclick="searchByAuthor()">Search by Author</button>
                
                <div class="form-group" style="margin-top: 15px;">
                    <label>Search by Title:</label>
                    <input type="text" id="searchTitle" placeholder="Enter book title" value="The Go Programming Language">
                </div>
                <button onclick="searchByTitle()">Search by Title</button>
                
                <div class="response" id="searchResponse"></div>
            </div>
        </div>
    </div>

    <script>
        const API_BASE = '/api/books';
        
        function formatResponse(data) {
            return JSON.stringify(data, null, 2);
        }

        function showResponse(elementId, data, isError = false) {
            const div = document.getElementById(elementId);
            div.innerHTML = '<pre class="' + (isError ? 'error' : 'success') + '">' + 
                           formatResponse(data) + '</pre>';
        }

        function handleError(error, elementId) {
            console.error('Error:', error);
            showResponse(elementId, { error: error.error || 'Unknown error' }, true);
        }

        document.getElementById('createForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const book = {
                title: formData.get('title'),
                author: formData.get('author'),
                published_year: parseInt(formData.get('published_year')),
                isbn: formData.get('isbn'),
                description: formData.get('description')
            };

            try {
                const response = await fetch(API_BASE, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(book)
                });
                const data = await response.json();
                if (response.ok) {
                    showResponse('createResponse', data);
                    e.target.reset();
                } else {
                    handleError(data, 'createResponse');
                }
            } catch (error) {
                handleError(error, 'createResponse');
            }
        });

        async function getAllBooks() {
            try {
                const response = await fetch(API_BASE);
                const data = await response.json();
                if (response.ok) {
                    showResponse('getAllResponse', data);
                } else {
                    handleError(data, 'getAllResponse');
                }
            } catch (error) {
                handleError(error, 'getAllResponse');
            }
        }

        async function getBookById() {
            const id = document.getElementById('getById').value;
            if (!id) {
                showResponse('getByIdResponse', { error: 'Please enter book ID' }, true);
                return;
            }

            try {
                const response = await fetch(API_BASE + '/' + id);
                const data = await response.json();
                if (response.ok) {
                    showResponse('getByIdResponse', data);
                } else {
                    handleError(data, 'getByIdResponse');
                }
            } catch (error) {
                handleError(error, 'getByIdResponse');
            }
        }

        document.getElementById('updateForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const id = formData.get('id');
            const book = {
                title: formData.get('title'),
                author: formData.get('author'),
                published_year: parseInt(formData.get('published_year')),
                isbn: formData.get('isbn'),
                description: formData.get('description')
            };

            try {
                const response = await fetch(API_BASE + '/' + id, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(book)
                });
                const data = await response.json();
                if (response.ok) {
                    showResponse('updateResponse', data);
                } else {
                    handleError(data, 'updateResponse');
                }
            } catch (error) {
                handleError(error, 'updateResponse');
            }
        });

        async function deleteBook() {
            const id = document.getElementById('deleteId').value;
            if (!id) {
                showResponse('deleteResponse', { error: 'Please enter book ID' }, true);
                return;
            }

            try {
                const response = await fetch(API_BASE + '/' + id, { method: 'DELETE' });
                const data = await response.json();
                if (response.ok) {
                    showResponse('deleteResponse', data);
                    document.getElementById('deleteId').value = '';
                } else {
                    handleError(data, 'deleteResponse');
                }
            } catch (error) {
                handleError(error, 'deleteResponse');
            }
        }

        async function searchByAuthor() {
            const author = document.getElementById('searchAuthor').value;
            if (!author) {
                showResponse('searchResponse', { error: 'Please enter author name' }, true);
                return;
            }

            try {
                const response = await fetch(API_BASE + '/search?author=' + encodeURIComponent(author));
                const data = await response.json();
                if (response.ok) {
                    showResponse('searchResponse', data);
                } else {
                    handleError(data, 'searchResponse');
                }
            } catch (error) {
                handleError(error, 'searchResponse');
            }
        }

        async function searchByTitle() {
            const title = document.getElementById('searchTitle').value;
            if (!title) {
                showResponse('searchResponse', { error: 'Please enter book title' }, true);
                return;
            }

            try {
                const response = await fetch(API_BASE + '/search?title=' + encodeURIComponent(title));
                const data = await response.json();
                if (response.ok) {
                    showResponse('searchResponse', data);
                } else {
                    handleError(data, 'searchResponse');
                }
            } catch (error) {
                handleError(error, 'searchResponse');
            }
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// MAIN
// =======================================================================

var content embed.FS

func main() {
	// Initialize the repository, service, and handler
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	// Create a simple HTTP server without gorilla/mux in main
	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", handler.HandleBooks)
	mux.HandleFunc("/api/books/", handler.HandleBooks)
	mux.HandleFunc("/", serveHTMLInterface)
	mux.HandleFunc("/index.html", serveHTMLInterface)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		}
	})

	// Start the server
	log.Println("Server starting on :8085")
	log.Println("Open http://localhost:8085 in your browser to test the API")

	if err := http.ListenAndServe(":8085", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
