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
	Username  string
	Connected bool
	ready     chan bool
	Messages  chan string
	server    *ChatServer
	msgMutx   sync.RWMutex
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// TODO: Implement this method
	// Hint: thread-safe, non-blocking send
	if !c.Connected {
		return
	}
	c.msgMutx.Lock()
	defer c.msgMutx.Unlock()
	c.Messages <- message
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// TODO: Implement this method
	// Hint: read from channel, handle closed channel
	c.msgMutx.RLock()
	defer c.msgMutx.RUnlock()
	msg, ok := <-c.Messages
	if !ok {
		return ""
	}
	return msg
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	// Hint: clients map, mutex
	clients   map[string]*Client
	broadcast chan string
	join      chan *Client
	leave     chan *Client
	mu        sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	// TODO: Implement this function
	cs := &ChatServer{
		clients:   make(map[string]*Client),
		broadcast: make(chan string),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		mu:        sync.RWMutex{},
	}
	go cs.run()

	return cs
}

func (s *ChatServer) run() {
	for {
		select {
		case client := <-s.join:
			s.mu.Lock()
			s.clients[client.Username] = client
			client.ready <- true
			s.mu.Unlock()
		case client := <-s.leave:
			s.mu.Lock()
			delete(s.clients, client.Username)
			client.ready <- true
			s.mu.Unlock()
		case message := <-s.broadcast:
			s.mu.Lock()
			for _, c := range s.clients {
				c.Messages <- message
			}
			s.mu.Unlock()
		}
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	if _, ok := s.clients[username]; ok {
		return nil, ErrUsernameAlreadyTaken
	}

	c := &Client{
		Username:  username,
		Connected: true,
		ready:     make(chan bool, 1),
		Messages:  make(chan string),
		server:    s,
		msgMutx:   sync.RWMutex{},
	}

	s.join <- c
	<-c.ready
	return c, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// TODO: Implement this method
	// Hint: remove from map, close channels
	if _, exists := s.clients[client.Username]; !exists {
		return
	}
	client.Connected = false
	s.leave <- client
	<-client.ready
	close(client.Messages)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
	msg := fmt.Sprintf("[%s] %s", sender.Username, message)

	s.broadcast <- msg
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	if !sender.Connected {
		return ErrClientDisconnected
	}
	r, exists := s.clients[recipient]
	if !exists {
		return ErrRecipientNotFound
	}
	r.Messages <- fmt.Sprintf("[%s] %s", sender.Username, message)

	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
