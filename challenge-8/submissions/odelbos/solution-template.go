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
	username     string
	incoming     chan string
	outgoing     chan string
	disconnect   chan struct{}
	disconnected bool
	mu           sync.RWMutex
}

// Send sends a message to the client (non-blocking)
func (c *Client) Send(message string) {
	if c.disconnected {
		return
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	select {
	case c.incoming <- message:
	default:
		// Do not block
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	if msg, ok := <-c.incoming; ok {
		return msg
	}
	return ""
}

func (c *Client) do_disconnect() {
	if c.disconnected {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.incoming)
	close(c.disconnect)
	c.disconnected = true
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

	if _, ok := s.clients[username]; ok {
		return nil, ErrUsernameAlreadyTaken
	}

	client := &Client{
		username:   username,
		incoming:   make(chan string, 100),
		outgoing:   make(chan string, 100),
		disconnect: make(chan struct{}),
	}
	s.clients[username] = client

	go s.handleClient(client)

	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[client.username]; ! ok {
		return
	}

	client.do_disconnect()
	delete(s.clients, client.username)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg := fmt.Sprintf("%s: %s", sender.username, message)
	for _, client := range(s.clients) {
		if client.username != sender.username {
			client.Send(msg)
		}
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	if sender.disconnected {
		return ErrClientDisconnected
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	target, ok := s.clients[recipient]
	if ! ok {
		return ErrRecipientNotFound
	}
	if target.disconnected {
		return ErrClientDisconnected
	}

	msg := fmt.Sprintf("(pm) %s: %s", sender.username, message)
	target.Send(msg)
	return nil
}

// handleClient processes outgoing messages and disconnection for a client
func (s *ChatServer) handleClient(client *Client) {
	for {
		select {
		case msg := <-client.outgoing:
			s.Broadcast(client, msg)
		case <-client.disconnect:
			s.Disconnect(client)
			return
		}
	}
}
