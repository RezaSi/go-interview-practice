// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"sync"
)

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)

// Client represents a connected chat client
type Client struct {
	Username string
	Messages chan string
	server   *ChatServer
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.Messages <- message
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	msg := <-c.Messages
	return msg
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients    map[string]*Client
	connect    chan *Client
	disconnect chan *Client
	mutex      sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	chatSrv := &ChatServer{
		clients:    make(map[string]*Client),
		connect:    make(chan *Client),
		disconnect: make(chan *Client),
	}
	return chatSrv
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mutex.RLock()
	if _, ok := s.clients[username]; ok {
		s.mutex.RUnlock()
		return nil, ErrUsernameAlreadyTaken
	}
	s.mutex.RUnlock()

	c := &Client{Username: username, Messages: make(chan string, 1), server: s}
	s.mutex.Lock()
	s.clients[username] = c
	s.mutex.Unlock()

	return c, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.clients, client.Username)
	close(client.Messages)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	for _, client := range s.clients {
		client.Messages <- message
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	recipientClient, ok := s.clients[recipient]
	if !ok {
		return ErrRecipientNotFound
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if _, ok := s.clients[sender.Username]; !ok {
		return ErrClientDisconnected
	}

	msg := fmt.Sprintf("message from %s: %s", sender.Username, message)
	recipientClient.Messages <- msg

	return nil
}
