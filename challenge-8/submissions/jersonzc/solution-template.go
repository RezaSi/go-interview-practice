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
	// Hint: username, message channel, mutex, disconnected flag
	username     string
	message      chan string
	disconnected bool
	mu           sync.Mutex
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// TODO: Implement this method
	// Hint: thread-safe, non-blocking send
	c.message <- message
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// TODO: Implement this method
	// Hint: read from channel, handle closed channel
	v, ok := <-c.message
	if !ok {
		return ""
	}
	return v
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	// Hint: clients map, mutex
	clients map[string]*Client
	mu      sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	// TODO: Implement this function
	return &ChatServer{
		clients: make(map[string]*Client),
		mu:      sync.Mutex{},
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, _ := range s.clients {
		if k == username {
			return nil, ErrUsernameAlreadyTaken
		}
	}

	client := &Client{
		username:     username,
		message:      make(chan string),
		disconnected: false,
		mu:           sync.Mutex{},
	}

	s.clients[username] = client

	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// TODO: Implement this method
	// Hint: remove from map, close channels
	s.mu.Lock()
	defer s.mu.Unlock()

	client.disconnected = true
	close(client.message)

	delete(s.clients, client.username)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
	if sender.disconnected {
		return
	}

	for _, client := range s.clients {
		if !client.disconnected {
			client.message <- fmt.Sprintf("From: %s, Msg: %q", sender.username, message)
		}
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	if sender.disconnected {
		return ErrClientDisconnected
	}

	client, ok := s.clients[recipient]
	if !ok {
		return ErrRecipientNotFound
	}

	if !client.disconnected {
		client.message <- fmt.Sprintf("From: %s, Msg: %q", sender.username, message)
	} else {
		return ErrClientDisconnected
	}

	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
