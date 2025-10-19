package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() *httptest.Server {
	// Initialize the repository, service, and handler
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)

	// Create a test HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", handler.HandleBooks)
	mux.HandleFunc("/api/books/", handler.HandleBooks)

	return httptest.NewServer(mux)
}

func TestGetAllBooksEmpty(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/books", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var books []*Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(books) != 0 {
		t.Errorf("Expected empty array; got %d books", len(books))
	}
}

func TestCreateBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a new book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", resp.Status)
	}

	var createdBook Book
	if err := json.NewDecoder(resp.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if createdBook.ID == "" {
		t.Error("Expected book to have an ID")
	}
	if createdBook.Title != book.Title {
		t.Errorf("Expected book title %s; got %s", book.Title, createdBook.Title)
	}
}

func TestCreateBookInvalid(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a book with missing required fields
	book := &Book{
		// Title intentionally missing
		Author: "John Doe",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request; got %v", resp.Status)
	}
}

func TestGetBookByID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First create a book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, _ := json.Marshal(book)
	resp, _ := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)

	var createdBook Book
	json.NewDecoder(resp.Body).Decode(&createdBook)
	resp.Body.Close()

	// Now get the book by ID
	resp, err := http.Get(fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var retrievedBook Book
	if err := json.NewDecoder(resp.Body).Decode(&retrievedBook); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if retrievedBook.ID != createdBook.ID {
		t.Errorf("Expected book ID %s; got %s", createdBook.ID, retrievedBook.ID)
	}
	if retrievedBook.Title != book.Title {
		t.Errorf("Expected book title %s; got %s", book.Title, retrievedBook.Title)
	}
}

func TestGetBookByIDNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/books/nonexistent", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", resp.Status)
	}
}

func TestUpdateBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First create a book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, _ := json.Marshal(book)
	resp, _ := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)

	var createdBook Book
	json.NewDecoder(resp.Body).Decode(&createdBook)
	resp.Body.Close()

	// Now update the book
	updatedBook := createdBook
	updatedBook.Description = "Updated description"

	updatedBookJSON, _ := json.Marshal(updatedBook)
	req, _ := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID),
		bytes.NewBuffer(updatedBookJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var returnedBook Book
	if err := json.NewDecoder(resp.Body).Decode(&returnedBook); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if returnedBook.Description != updatedBook.Description {
		t.Errorf("Expected description %s; got %s", updatedBook.Description, returnedBook.Description)
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	book := &Book{
		ID:            "nonexistent",
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, _ := json.Marshal(book)
	req, _ := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/books/nonexistent", server.URL),
		bytes.NewBuffer(bookJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", resp.Status)
	}
}

func TestDeleteBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First create a book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, _ := json.Marshal(book)
	resp, _ := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)

	var createdBook Book
	json.NewDecoder(resp.Body).Decode(&createdBook)
	resp.Body.Close()

	// Now delete the book
	req, _ := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID),
		nil,
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	// Verify the book was deleted
	resp, _ = http.Get(fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID))
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found after deletion; got %v", resp.Status)
	}
	resp.Body.Close()
}

func TestDeleteBookNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/books/nonexistent", server.URL),
		nil,
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", resp.Status)
	}
}

func TestSearchBooksByAuthor(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create several books
	books := []*Book{
		{
			Title:         "The Go Programming Language",
			Author:        "Alan A. A. Donovan and Brian W. Kernighan",
			PublishedYear: 2015,
			ISBN:          "978-0134190440",
			Description:   "The definitive guide to programming in Go",
		},
		{
			Title:         "Go in Action",
			Author:        "William Kennedy",
			PublishedYear: 2015,
			ISBN:          "978-1617291784",
			Description:   "An introduction to Go",
		},
		{
			Title:         "The C Programming Language",
			Author:        "Brian W. Kernighan and Dennis Ritchie",
			PublishedYear: 1988,
			ISBN:          "978-0131103627",
			Description:   "The definitive guide to C",
		},
	}

	for _, book := range books {
		bookJSON, _ := json.Marshal(book)
		resp, _ := http.Post(
			fmt.Sprintf("%s/api/books", server.URL),
			"application/json",
			bytes.NewBuffer(bookJSON),
		)
		resp.Body.Close()
	}

	// Search for books by Kernighan
	resp, err := http.Get(fmt.Sprintf("%s/api/books/search?author=Kernighan", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var foundBooks []*Book
	if err := json.NewDecoder(resp.Body).Decode(&foundBooks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(foundBooks) != 2 {
		t.Errorf("Expected 2 books; got %d", len(foundBooks))
	}
}

func TestSearchBooksByTitle(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create several books
	books := []*Book{
		{
			Title:         "The Go Programming Language",
			Author:        "Alan A. A. Donovan and Brian W. Kernighan",
			PublishedYear: 2015,
			ISBN:          "978-0134190440",
			Description:   "The definitive guide to programming in Go",
		},
		{
			Title:         "Go in Action",
			Author:        "William Kennedy",
			PublishedYear: 2015,
			ISBN:          "978-1617291784",
			Description:   "An introduction to Go",
		},
		{
			Title:         "The C Programming Language",
			Author:        "Brian W. Kernighan and Dennis Ritchie",
			PublishedYear: 1988,
			ISBN:          "978-0131103627",
			Description:   "The definitive guide to C",
		},
	}

	for _, book := range books {
		bookJSON, _ := json.Marshal(book)
		resp, _ := http.Post(
			fmt.Sprintf("%s/api/books", server.URL),
			"application/json",
			bytes.NewBuffer(bookJSON),
		)
		resp.Body.Close()
	}

	// Search for Go books
	resp, err := http.Get(fmt.Sprintf("%s/api/books/search?title=Go", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var foundBooks []*Book
	if err := json.NewDecoder(resp.Body).Decode(&foundBooks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(foundBooks) != 2 {
		t.Errorf("Expected 2 books; got %d", len(foundBooks))
	}
}

func TestSearchBooksByTitleNoResults(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create several books
	books := []*Book{
		{
			Title:         "The Go Programming Language",
			Author:        "Alan A. A. Donovan and Brian W. Kernighan",
			PublishedYear: 2015,
			ISBN:          "978-0134190440",
			Description:   "The definitive guide to programming in Go",
		},
		{
			Title:         "Go in Action",
			Author:        "William Kennedy",
			PublishedYear: 2015,
			ISBN:          "978-1617291784",
			Description:   "An introduction to Go",
		},
	}

	for _, book := range books {
		bookJSON, _ := json.Marshal(book)
		resp, _ := http.Post(
			fmt.Sprintf("%s/api/books", server.URL),
			"application/json",
			bytes.NewBuffer(bookJSON),
		)
		resp.Body.Close()
	}

	// Search for Python books
	resp, err := http.Get(fmt.Sprintf("%s/api/books/search?title=Python", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var foundBooks []*Book
	if err := json.NewDecoder(resp.Body).Decode(&foundBooks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(foundBooks) != 0 {
		t.Errorf("Expected 0 books; got %d", len(foundBooks))
	}
}
