// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
	// Add any other necessary imports
)

// Client represents a connected chat client
type Client struct {
	// TODO: Implement this struct
	// Hint: username, message channel, mutex, disconnected flag
	username string
	message  chan string
	mutex    sync.Locker
	isOnline bool
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
	message, ok := <-c.message
	if !ok {
		return ""
	}
	return message
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	// Hint: clients map, mutex
	onlineClient     map[*Client]chan string
	onlineClientList map[string]struct{}
	l                *sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	// TODO: Implement this function
	return &ChatServer{
		onlineClient:     make(map[*Client]chan string),
		onlineClientList: map[string]struct{}{},
		l:                &sync.Mutex{},
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	if _, ok := s.onlineClientList[username]; ok {
		return nil, ErrUsernameAlreadyTaken
	}
	Client := &Client{
		username: username,
		message:  make(chan string),
		isOnline: true,
	}
	s.l.Lock()
	defer s.l.Unlock()
	s.onlineClient[Client] = Client.message
	s.onlineClientList[username] = struct{}{}
	return Client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// TODO: Implement this method
	// Hint: remove from map, close channels

	if client == nil || !client.isOnline {
		return
	}
	s.l.Lock()
	// client.mutex.Lock()
	defer s.l.Unlock()
	// defer client.mutex.Unlock()
	delete(s.onlineClient, client)
	delete(s.onlineClientList, client.username)

	client.isOnline = false
	close(client.message)

}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
	if _, ok := s.onlineClient[sender]; !ok {
		return
	}
	for _, ch := range s.onlineClient {
		ch <- message
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	if _, ok := s.onlineClient[sender]; !ok {
		return ErrClientDisconnected
	}
	if _, ok := s.onlineClientList[recipient]; !ok {
		return ErrRecipientNotFound
	}
	for k, v := range s.onlineClient {
		if k.username == recipient {
			v <- message
			return nil
		}
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
