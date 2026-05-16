// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"sync"
	// Add any other necessary imports
)

// Client
// =========================================
// Client represents a connected chat client
type Client struct {
	username     string      //username
	inbox        chan string //channel for incoming messages
	disconnected bool        //flag to prevent Send write to closed inbox
	mu           sync.Mutex
}

// Client constructor returns *Client
func newClient(username string) *Client {
	return &Client{
		username: username,
		inbox:    make(chan string, 150),
	}

}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// Check if user disconnected under users mutex
	c.mu.Lock()
	if c.disconnected {
		c.mu.Unlock()
		return
	}
	inbox := c.inbox
	c.mu.Unlock()
	defer func() {
		// Recover if channel was closed between check and close
		recover()
	}()
	inbox <- message
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// Blocking reading from inbox channel
	result := <-c.inbox
	return result
}

// ChatServer
// =========================================
// ChatServer manages client connections and message routing
type ChatServer struct {
	users map[string]*Client
	mu    sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		users: make(map[string]*Client),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// Check if user exist
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isConnected(username) {
		return nil, ErrUsernameAlreadyTaken
	}
	// Create a new user
	user := newClient(username)
	s.users[username] = user
	return user, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// Check if user exist
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isCurrentClient(client) {
		// Gracefull disconnection set flag, close channel, delete user
		client.mu.Lock()
		client.disconnected = true
		close(client.inbox)
		client.mu.Unlock()
		delete(s.users, client.username)
	}
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// Check if sender exist
	s.mu.Lock()
	if !s.isCurrentClient(sender) || message == "" {
		s.mu.Unlock()
		return
	}
	// Send formatted message to all users exept sender
	recipients := make([]*Client, 0, len(s.users)-1)
	for name, user := range s.users {
		if name != sender.username {
			recipients = append(recipients, user)
		}
	}
	s.mu.Unlock()
	// Formating message
	fstring := fmt.Sprintf("{%s}: %s", sender.username, message)
	// Send Broadcast message to users
	for _, user := range recipients {
		user.Send(fstring)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// Check if sender exist
	s.mu.Lock()
	if !s.isCurrentClient(sender) {
		s.mu.Unlock()
		return ErrClientDisconnected
	}
	// Check if recipient exist
	if !s.isConnected(recipient) {
		s.mu.Unlock()
		return ErrRecipientNotFound
	}
	// Check if message not empty
	if message == "" {
		s.mu.Unlock()
		return ErrBadParam
	}
	// Format the message
	fstring := fmt.Sprintf("{%s}->{%s}:%s", sender.username, recipient, message)
	// Send the formatted string
	target := s.users[recipient]
	s.mu.Unlock()
	target.Send(fstring)
	return nil
}

// isConnected returns true if user connected, otherwise false
func (s *ChatServer) isConnected(username string) bool {
	if _, exist := s.users[username]; exist {
		return true
	}
	return false
}

// isCurrentClient returns true if Client not nil and equal to client in storage
func (s *ChatServer) isCurrentClient(client *Client) bool {
	if client == nil {
		return false
	}
	current, ok := s.users[client.username]
	return ok && client == current
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	ErrBadParam             = errors.New("bad or missing parameter")
)
