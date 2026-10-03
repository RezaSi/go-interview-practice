// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Client represents a connected chat client
type Client struct {
	disconnected    bool
	username        string
	messageCh       chan string
	
	mu              sync.Mutex
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.disconnected {
	    return
	}
	
	select {
	    case c.messageCh <- message:
	    default:
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	msg, ok := <-c.messageCh
    if !ok {
        return ""
    }
    return msg
}

func (c *Client) isDisconnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.disconnected
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	mu          sync.RWMutex
	clients     map[string]*Client
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
	    clients: make(map[string]*Client, 100),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.clients[username]; exists {
	    return nil, ErrUsernameAlreadyTaken
	}
	
	newClient := &Client{
	    username: username,
	    messageCh: make(chan string, 100),
	}
	
	s.clients[username] = newClient
	
	return newClient, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	if _, exists := s.clients[client.username]; !exists {
	    s.mu.Unlock()
	    return
	}
	delete(s.clients, client.username)
	s.mu.Unlock()
	
	client.mu.Lock()
	if !client.disconnected {
	    client.disconnected = true
	    close(client.messageCh)
	}
	client.mu.Unlock()
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for name, client := range s.clients {
	    if sender.username == name {
	        continue
	    }
	    
	    client.Send(sender.username + ": " + message)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, exists := s.clients[recipient]
	if !exists {
	    return ErrRecipientNotFound
	}
	
	if sender.isDisconnected() || user.isDisconnected() {
	    return ErrClientDisconnected
	}
	
	if sender.username == user.username {
	    user.Send("[Note]: " + message)
	} else {
	    user.Send("[Private] " + sender.username + ": " + message)
	}

	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)

// TestChatServer_RaceConditions tests chat concurrency
func TestChatServer_RaceConditions(t *testing.T) {
    server := NewChatServer()
    
    const clientCount = 100
    
    var wg sync.WaitGroup
    
    clients := make([]*Client, clientCount)
    for i := 0; i < clientCount; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            username := fmt.Sprintf("user_%d", id)
            client, err := server.Connect(username)
            if err != nil {
                t.Errorf("failed to connect %s: %v", username, err)
                return
            }
            clients[id] = client
        }(i)
    }
    wg.Wait()
    
    for i := 0; i < clientCount; i++ {
        go func(c *Client) {
            for {
                msg := c.Receive()
                if msg == "" {
                    return
                }
            }
        }(clients[i])
    }
    
    for i := 0; i < clientCount; i++ {
        wg.Add(3)
        go func(c *Client) {
            defer wg.Done()
            for j := 0; j < 5; j++ {
                server.Broadcast(c, fmt.Sprintf("hello everyone %d", j))
                time.Sleep(time.Millisecond)
            }
        }(clients[i])
        
        go func(c *Client, id int) {
            defer wg.Done()
            nextUser := fmt.Sprintf("user_%d", (id+1)%clientCount)
            for j := 0; j < 5; j++ {
                _ = server.PrivateMessage(c, nextUser, "secret msg")
                time.Sleep(time.Millisecond)
            }
        }(clients[i], i)
        
        go func(c *Client) {
            defer wg.Done()
            time.Sleep(2 * time.Millisecond)
            server.Disconnect(c)
        }(clients[i])
    }
    
    wg.Wait()
}
