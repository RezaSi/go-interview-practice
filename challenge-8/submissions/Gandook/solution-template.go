// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
	// Add any other necessary imports
)

const CAPACITY = 1024

// Client represents a connected chat client
type Client struct {
	username           string
	msgChan            chan string
	sndMutex, rcvMutex sync.Mutex
	connected          bool
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.sndMutex.Lock()
	defer c.sndMutex.Unlock()
	c.msgChan <- message
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	c.rcvMutex.Lock()
	defer c.rcvMutex.Unlock()
	message := <-c.msgChan
	return message
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients map[string]*Client
	mtx     sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]*Client),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if client, exists := s.clients[username]; exists {
		return nil, ErrUsernameAlreadyTaken
	} else {
		client = &Client{
			username:  username,
			msgChan:   make(chan string, CAPACITY),
			connected: true,
		}
		s.clients[username] = client
		return client, nil
	}
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.clients, client.username)
	client.connected = false
	close(client.msgChan)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, client := range s.clients {
		client.Send(message)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if !sender.connected {
		return ErrClientDisconnected
	}
	if client, exists := s.clients[recipient]; !exists {
		return ErrRecipientNotFound
	} else if !client.connected {
		return ErrClientDisconnected
	} else {
		client.Send(message)
		return nil
	}
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
