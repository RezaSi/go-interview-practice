// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
	"fmt"
)

// Client represents a connected chat client
type Client struct {
	// TODO: Implement this struct
	// Hint: username, message channel, mutex, disconnected flag
	username string
	incomingMessages chan string
	connected bool
	
	mu sync.RWMutex
}

var (
    clientBufferSize int = 5
)

func (c *Client) isConnected() bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.connected
}

func (c *Client) getName() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.username
}

func (c *Client) onConnection() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.incomingMessages = make(chan string, clientBufferSize)
    c.connected = true
}

func (c *Client) onDisconnection() {
    c.mu.Lock()
    defer c.mu.Unlock()
    close(c.incomingMessages)
    c.connected = false
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// TODO: Implement this method
	// Hint: thread-safe, non-blocking send
	select {
	case c.incomingMessages <- message:
	    return
	default:
	    return
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// TODO: Implement this method
	// Hint: read from channel, handle closed channel
	msg, ok := <- c.incomingMessages
    if !ok || msg == "" {
        return ""
    }
    return fmt.Sprintln(msg)
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	// Hint: clients map, mutex
	connectedClients map[string]*Client
	
	mu sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	// TODO: Implement this function
	return &ChatServer{
	    connectedClients: map[string]*Client{},
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if client, ok := s.connectedClients[username]; ok {
	    if client.isConnected(){
	        return nil, ErrUsernameAlreadyTaken
	    }
	    client.onConnection()
	    return client, nil
	}
	
	newClient := Client{
	    username: username,
	}
	newClient.onConnection()
	s.connectedClients[username] = &newClient
	return &newClient, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// TODO: Implement this method
	// Hint: remove from map, close channels
	s.mu.Lock()
	defer s.mu.Unlock()
	if _,ok := s.connectedClients[client.getName()]; !ok {
	    return //ErrClientNotFound
	}
	if !client.isConnected() {
	    return //ErrClientDisconnected
	}
	client.onDisconnection()
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
	if !sender.isConnected() {
	    return 
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for _, client := range s.connectedClients {
	    sender.mu.RLock()
	    formattedMessage := fmt.Sprintf("%s: %s", sender.username, message)
	    sender.mu.RUnlock()
	    client.Send(formattedMessage)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	if !sender.isConnected() {
	    return ErrClientDisconnected
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for name, client := range s.connectedClients{
	    if recipient == name{
	        client.mu.RLock()
	        if !client.connected{ //was in server but left
	            client.mu.RUnlock()
	            return ErrClientDisconnected
	        }
	        client.mu.RUnlock()
	        
	        sender.mu.RLock()
	        formattedMessage := fmt.Sprintf("<private>%s: %s", sender.username, message)
	        sender.mu.RUnlock()
	        client.Send(formattedMessage)
	        return nil
	    }
	}
	return ErrRecipientNotFound //never in server
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	//ErrClientNotFound = errors.New("client not found in server")
)
