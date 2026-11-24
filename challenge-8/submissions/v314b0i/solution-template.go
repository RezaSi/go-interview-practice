// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
	"log"
	// Add any other necessary imports
)

// Client represents a connected chat client
type Client struct {
	userName    string
	inboxChan   chan string
	isConnected bool
	// Hint: username, message channel, mutex, disconnected flag
}

// Send sends a message to the client
func (c *Client) Send(message string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to send message to client. Client, message:", c.userName, message)
		}
	}()
	if !c.isConnected {
	    return ErrClientDisconnected
	}
	c.inboxChan <- message
	return nil
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// Hint: read from channel, handle closed channel
	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to recieve message for user ", c.userName, err)
		}
	}()
	return <-c.inboxChan
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	mu      sync.Mutex
	clients map[string]*Client
	// Hint: clients map, mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	newChatServer := ChatServer{}
	newChatServer.clients = map[string]*Client{}
	return &newChatServer
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (retClient *Client, err error) {
	defer func() {
	    if err != nil {
	        return
	    }
	    if r := recover(); r!=nil{
		    err = r.(error)
	    }
	}()

	if _, exists := s.clients[username]; !exists {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.clients[username] = &Client{}
		s.clients[username].userName = username
		s.clients[username].inboxChan = make(chan string)
		s.clients[username].isConnected = true
	} else {
	    err = ErrUsernameAlreadyTaken
	    return
	}
	// Hint: check username, create client, add to map
	retClient = s.clients[username]
	return
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error in disconnect", err)
		}
	}()
	if _, exists := s.clients[client.userName]; !exists {
		return
	}
	close(s.clients[client.userName].inboxChan)
	s.clients[client.userName].isConnected = false
	delete(s.clients, client.userName)
	return
	// Hint: remove from map, close channels
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// Hint: format message, send to all clients
	formattedMsg := "[BROADCAST] " + message
	for _, client := range s.clients {
		if client.userName != sender.userName {
			err := s.PrivateMessage(sender, client.userName, formattedMsg)
			if err != nil {
				log.Println("Error Broadcasting message. sender, message, error", sender.userName, message, err)
			}
		}
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) (err error) {
	defer func() {
	    if err != nil {
	        return // Return the original error
	    }
		if r := recover(); r != nil {
			err = ErrFailedSend
		}
	}()
	if _, exists := s.clients[recipient]; !exists {
		err = ErrRecipientNotFound
		return
	}
	if !s.clients[recipient].isConnected || !sender.isConnected {
		err = ErrClientDisconnected
		return
	}
	s.clients[recipient].Send(message)
	// Hint: find recipient, check errors, send message
	return
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
	ErrFailedSend = errors.New("failed message send to channel")
)