// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// Client represents a connected chat client
type Client struct {
	username     string
	message      chan string
	disconnected atomic.Bool
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	if c.disconnected.Load() {
		return
	}
	c.message <- message
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	return <-c.message
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{clients: make(map[string]*Client)}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.clients[username]
	if ok {
		return nil, ErrUsernameAlreadyTaken
	}
	client := Client{username: username, message: make(chan string)}
	s.clients[username] = &client
	return &client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	client.disconnected.Store(true)
	close(client.message)
	delete(s.clients, client.username)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	if sender.disconnected.Load() {
		return
	}
	s.mu.RLock()
	recipients := make([]*Client, 0, len(s.clients))
	for _, recipient := range s.clients {
		if recipient.username == sender.username {
			continue
		}
		recipients = append(recipients, recipient)
	}
	s.mu.RUnlock()
	messageWithSender := fmt.Sprintf("[%s] %s", sender.username, message)
	for _, recipient := range recipients {
		recipient.Send(messageWithSender)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	if sender.disconnected.Load() {
		return ErrClientDisconnected
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	recipientClient, ok := s.clients[recipient]
	if !ok {
		return ErrRecipientNotFound
	}
	if recipientClient.disconnected.Load() {
		return ErrClientDisconnected
	}
	messageWithSender := fmt.Sprintf("[%s] %s", sender.username, message)
	recipientClient.Send(messageWithSender)
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)
