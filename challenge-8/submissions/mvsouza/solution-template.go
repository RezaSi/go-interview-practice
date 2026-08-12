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
	Username     string
	MessageChan  chan string
	Disconnected bool
	mu           sync.Mutex
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.Disconnected {
		return
	}

	select {
	case c.MessageChan <- message:
	default:
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	msg, ok := <-c.MessageChan
	if !ok {
		return ""
	}
	return msg
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

	if _, hasName := s.clients[username]; hasName {
		return nil, ErrUsernameAlreadyTaken
	}
	newClient := Client{
		Username:     username,
		MessageChan:  make(chan string, 10), // Adding a small buffer as good practice, but non-blocking select works too
		Disconnected: false,
	}
	s.clients[username] = &newClient
	return &newClient, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[client.Username]; exists {
		client.mu.Lock()
		client.Disconnected = true
		close(client.MessageChan)
		client.mu.Unlock()
		delete(s.clients, client.Username)
	}
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	messageFormated := fmt.Sprintf("%s: %s", sender.Username, message)
	for _, client := range s.clients {
		if client.Username != sender.Username {
			client.Send(messageFormated)
		}
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	s.mu.Lock()
	r, ok := s.clients[recipient]
	s.mu.Unlock()

	if !ok {
		return ErrRecipientNotFound
	}

	sender.mu.Lock()
	senderDisconnected := sender.Disconnected
	sender.mu.Unlock()

	if senderDisconnected {
		return ErrClientDisconnected
	}

	messageFormated := fmt.Sprintf("%s: %s", sender.Username, message)
	r.Send(messageFormated)
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
