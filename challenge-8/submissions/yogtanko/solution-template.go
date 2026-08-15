// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"sync"
	// Add any other necessary imports
)

// Client represents a connected chat client
type Client struct {
	// TODO: Implement this struct
	Username       string
	Message        chan string
	mu             sync.Mutex
	isDisconnected bool
	// Hint: username, message channel, mutex, disconnected flag
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// TODO: Implement this method
	c.mu.Lock()
	if c.isDisconnected {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()
	select {
	case c.Message <- message:
	default:
	}
	// Hint: thread-safe, non-blocking send
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// TODO: Implement this method
	// Hint: read from channel, handle closed channel
	return <-c.Message
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	clients map[string]*Client
	mu      sync.Mutex
	// Hint: clients map, mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	// TODO: Implement this function
	return &ChatServer{
		clients: make(map[string]*Client),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clients[username] != nil {
		return nil, ErrUsernameAlreadyTaken
	}
	newClient := &Client{
		Username:       username,
		isDisconnected: false,
		Message:        make(chan string, 100),
	}
	s.clients[username] = newClient
	return newClient, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// TODO: Implement this method
	// Hint: remove from map, close channels
	s.mu.Lock()
	delete(s.clients, client.Username)
	s.mu.Unlock()
	client.mu.Lock()
	close(client.Message)
	client.isDisconnected = true
	client.mu.Unlock()
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, client := range s.clients {
		go client.Send(fmt.Sprintf("%s : %s", sender.Username, message))
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	s.mu.Lock()
	defer s.mu.Unlock()
	sender.mu.Lock()
	defer sender.mu.Unlock()
	if sender.isDisconnected {
		return ErrClientDisconnected
	}
	if s.clients[recipient] != nil {
		go s.clients[recipient].Send((fmt.Sprintf("%s : %s", sender.Username, message)))
		return nil
	}
	return ErrRecipientNotFound
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
