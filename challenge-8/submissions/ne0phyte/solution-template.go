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
	username  string
	incoming  chan string
	mu        sync.Mutex
	connected bool
}

func NewClient(username string) *Client {
	return &Client{
		username:  username,
		incoming:  make(chan string),
		connected: true,
	}
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.mu.Lock()
	if c.connected {
		c.incoming <- message
	}
	c.mu.Unlock()
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	if msg, ok := <-c.incoming; ok {
		return msg
	}
	return ""
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients map[string]*Client
	mu      sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]*Client),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.clients[username]; exists {
		return nil, ErrUsernameAlreadyTaken
	}
	client := NewClient(username)
	s.clients[username] = client
	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	// flag client as disconnected and close channel
	client.mu.Lock()
	client.connected = false
	close(client.incoming)

	// remove client from server
	delete(s.clients, client.username)
	client.mu.Unlock()
	s.mu.Unlock()
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	msg := fmt.Sprintf("%s: %s", sender.username, message)

	s.mu.Lock()
	for _, client := range s.clients {
		client.Send(msg)
	}
	s.mu.Unlock()
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !sender.connected {
		return ErrClientDisconnected
	}
	targetClient := s.clients[recipient]
	if targetClient == nil {
		return ErrRecipientNotFound
	}
	if !targetClient.connected {
		return ErrClientDisconnected
	}
	targetClient.Send(fmt.Sprintf("PRIVATE: %s: %s", sender.username, message))
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
