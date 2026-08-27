// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"sync"
)

// Client represents a connected chat client
type Client struct {
	mu           sync.Mutex
	username     string
	inMessage    chan string
	outMessage   chan string
	disconnected bool
}

func (c *Client) Username() string {
	return c.username
}

func (c *Client) CloseChannels() {
	close(c.inMessage)
	close(c.outMessage)
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// FIX: Change to a blocking send so no messages are dropped under heavy load.
	// Since we unlocked s.mu early in the server, this will safely pause only
	// this sender thread until the client's receiver catches up.
	c.outMessage <- message
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	msg, ok := <-c.outMessage
	if !ok {
		return "" 
	}
	return msg
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	mu      sync.Mutex
	// FIX: Changed map value type from Client to *Client
	clients map[string]*Client 
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

	_, ok := s.clients[username]
	if ok {
		return nil, ErrUsernameAlreadyTaken
	}

	const bufferSize = 10

	client := &Client{
		username:     username,
		inMessage:    make(chan string, bufferSize),
		outMessage:   make(chan string, bufferSize),
		disconnected: false,
	}

	s.clients[username] = client

	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, client.Username())
	client.CloseChannels()
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// 1. Take a quick snapshot of clients under the lock
	s.mu.Lock()
	snapshot := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		if client != sender {
			snapshot = append(snapshot, client)
		}
	}
	s.mu.Unlock() // Release the lock immediately!

	// 2. Format and send the message outside the critical lock section
	formattedMsg := fmt.Sprintf("[%s]: %s", sender.Username(), message)
	for _, client := range snapshot {
		client.Send(formattedMsg) // Now blocks safely if buffer is full
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	
	s.mu.Lock()
	
	_, senderExists := s.clients[sender.Username()]
	if !senderExists {
		s.mu.Unlock()
		return ErrClientDisconnected // FIX: Returns the expected error for the test
	}
	
	targetClient, exists := s.clients[recipient]
	s.mu.Unlock() // Release the lock immediately!

	if !exists {
		return ErrRecipientNotFound
	}

	// 2. Send the message outside the critical lock section
	formattedMsg := fmt.Sprintf("[Whisper from %s]: %s", sender.Username(), message)
	targetClient.Send(formattedMsg) // Now blocks safely if buffer is full
	
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)
