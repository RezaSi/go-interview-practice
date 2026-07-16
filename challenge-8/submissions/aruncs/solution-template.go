// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
	// Add any other necessary imports
)

type Message struct {
	Sender   string
	Receiver string
	Message  string
}

// Client represents a connected chat client
type Client struct {
	// TODO: Implement this struct
	// Hint: username, message channel, mutex, disconnected flag
	username   string
	outChannel chan Message
	inChannel  chan Message
	connected  bool
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// TODO: Implement this method
	// Hint: thread-safe, non-blocking send

}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// TODO: Implement this method
	// Hint: read from channel, handle closed channel
	message := <-c.inChannel
	return message.Message
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	// Hint: clients map, mutex
	mu      sync.RWMutex
	clients map[string]*Client
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	// TODO: Implement this function
	return &ChatServer{
		clients: make(map[string]*Client, 10),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isUsernameExists(username) {
		return nil, ErrUsernameAlreadyTaken
	}

	client := &Client{
		username:   username,
		inChannel:  make(chan Message, 10),
		outChannel: make(chan Message, 10),
		connected:  true,
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

	delete(s.clients, client.username)
	close(client.inChannel)
	close(client.outChannel)

	client.connected = false
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, receiver := range s.clients {
		if receiver.username == sender.username {
			continue
		}
		receiver.inChannel <- Message{
			Sender:   sender.username,
			Receiver: receiver.username,
			Message:  message,
		}
	}

}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !sender.connected {
		return ErrClientDisconnected
	}

	receiver, exists := s.clients[recipient]

	if !exists {
		return ErrRecipientNotFound
	}

	receiver.inChannel <- Message{
		Sender:   sender.username,
		Receiver: recipient,
		Message:  message,
	}
	return nil
}

func (s *ChatServer) isUsernameExists(username string) bool {
	_, exists := s.clients[username]
	return exists
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
