// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	// Add any other necessary imports
)

// Client represents a connected chat client
type Client struct {
	Username     string
	Conn         net.Conn
	Outgoing     chan string
	Disconnected bool
	mu           sync.Mutex
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.mu.Lock()
	if c.Disconnected {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()

	select {
	case c.Outgoing <- message:
	default:

	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	msg, ok := <-c.Outgoing
	if !ok {
		return ""
	}
	return msg
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// Hint: clients map, mutex
	clients   map[*Client]bool
	broadcast chan string
	join      chan *Client
	leave     chan *Client
	mu        sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:   make(map[*Client]bool),
		broadcast: make(chan string),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	log.Printf("[connect] username=%s", username)

	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		if client.Username == username {
			log.Printf("[connect-failed] username=%s reason=already_taken", username)
			return nil, ErrUsernameAlreadyTaken
		}
	}

	client := &Client{
		Username:     username,
		Conn:         nil,
		Outgoing:     make(chan string, 8),
		Disconnected: false,
	}
	s.clients[client] = true
	log.Printf("[connect-success] username=%s total_clients=%d", username, len(s.clients))

	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	_, exists := s.clients[client]
	if !exists {
		s.mu.Unlock()
		return
	}
	delete(s.clients, client)
	remaining := len(s.clients)
	s.mu.Unlock()

	client.mu.Lock()
	if !client.Disconnected {
		close(client.Outgoing)
		client.Disconnected = true
	}
	client.mu.Unlock()
	log.Printf("[disconnect] username=%s remaining_clients=%d", client.Username, remaining)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		log.Printf("[broadcast-ignored] sender=%s reason=empty_message", sender.Username)
		return
	}

	sender.mu.Lock()
	if sender.Disconnected {
		log.Printf("[broadcast-ignored] sender=%s reason=disconnected", sender.Username)
		return
	}
	sender.mu.Unlock()

	formatted := fmt.Sprintf("%s, %s", sender.Username, message)

	s.mu.Lock()
	count := 0
	for client := range s.clients {
		if client != sender {
			client.Send(formatted)
			count++
		}
	}
	s.mu.Unlock()
	log.Printf("[broadcast-success] sender=%s count=%d", sender.Username, count)
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return ErrEmptyMessage
	}

	sender.mu.Lock()
	if sender.Disconnected {
		return ErrSenderDisconnected
	}
	sender.mu.Unlock()

	s.mu.Lock()
	var rC *Client
	for client := range s.clients {
		if client.Username == recipient {
			rC = client
			break
		}
	}
	s.mu.Unlock()

	if rC == nil {
		return ErrRecipientNotFound
	}

	rC.mu.Lock()
	disconnected := rC.Disconnected
	rC.mu.Unlock()
	if disconnected {
		return ErrRecipientDisconnected
	}

	formatted := fmt.Sprintf("%s, %s", sender.Username, message)
	rC.Send(formatted)

	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken  = errors.New("username already taken")
	ErrRecipientNotFound     = errors.New("recipient not found")
	ErrSenderDisconnected    = errors.New("sender disconnected")
	ErrRecipientDisconnected = errors.New("recipient disconnected")
	ErrEmptyMessage          = errors.New("empty message")
)
