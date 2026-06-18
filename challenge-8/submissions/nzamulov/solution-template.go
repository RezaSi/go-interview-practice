package challenge8

import (
	"errors"
	"sync"
)

// Client represents a connected chat client
type Client struct {
	// Hint: username, message channel, mutex, disconnected flag
	username     string
	msgs         chan string
	disconnected bool
}

func NewClient(username string) *Client {
	return &Client{
		username: username,
		msgs:     make(chan string),
	}
}

func (c *Client) IsDisconnected() bool {
	return c.disconnected
}

func (c *Client) Disconnect() {
	close(c.msgs)
	c.disconnected = true
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.msgs <- message
	// Hint: thread-safe, non-blocking send
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// Hint: read from channel, handle closed channel
	return <-c.msgs
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// Hint: clients map, mutex
	clients map[string]*Client
	m       sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]*Client),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// Hint: check username, create client, add to map
	s.m.Lock()
	defer s.m.Unlock()
	if _, found := s.clients[username]; found {
		return nil, ErrUsernameAlreadyTaken
	}
	c := NewClient(username)
	s.clients[username] = c
	return c, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// Hint: remove from map, close channels
	if client.IsDisconnected() {
		return
	}
	s.m.Lock()
	defer s.m.Unlock()
	client.Disconnect()
	delete(s.clients, client.username)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// Hint: format message, send to all clients
	if sender.IsDisconnected() {
		return
	}
	s.m.RLock()
	defer s.m.RUnlock()
	for _, client := range s.clients {
		if client.IsDisconnected() {
			continue
		}
		if client.username == sender.username {
			continue
		}
		client.Send(message)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// Hint: find recipient, check errors, send message
	if sender.IsDisconnected() {
		return ErrClientDisconnected
	}
	s.m.RLock()
	defer s.m.RUnlock()
	client, found := s.clients[recipient]
	if !found {
		return ErrRecipientNotFound
	}
	if client.IsDisconnected() {
		return ErrClientDisconnected
	}
	client.Send(message)
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
