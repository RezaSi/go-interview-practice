// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
)

// max pending messages for a client, more messages will be dropped
const bufferSize = 1024

// Client represents a connected chat client
type Client struct {
	username string
	messages chan string
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	select {
	// sending to disconnected will panic
	case c.messages <- message:
		return
	default:
		// buffer full - message loss
		return
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	return <-c.messages
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

	if _, ok := s.clients[username]; ok {
		return nil, ErrUsernameAlreadyTaken
	}

	c := &Client{
		username: username,
		messages: make(chan string, bufferSize),
	}

	s.clients[username] = c

	return c, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[client.username]; !ok {
		return
	}

	close(client.messages)
	delete(s.clients, client.username)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	for _, c := range s.clients {
		c.Send(message)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	if _, ok := s.clients[sender.username]; !ok {
		return ErrClientDisconnected
	}

	rec, ok := s.clients[recipient]
	if !ok {
		return ErrRecipientNotFound
	}

	rec.Send(message)
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)
