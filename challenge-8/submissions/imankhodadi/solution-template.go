// Challenge 8: Chat Server with Channels
package challenge8

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type ChatServer struct {
	clients map[string]*Client
	mu      sync.RWMutex
}
type Client struct {
	Username     string
	Messages     chan string
	mu           sync.Mutex
	disconnected bool
}

var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)

func (c *Client) Send(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.disconnected {
		return
	}
	select {
	case c.Messages <- message:
	default:
		// non-blocking
	}
}
func (c *Client) Receive() string {
	msg, ok := <-c.Messages
	if !ok {
		return ""
	}
	return msg
}
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]*Client),
	}
}
func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.TrimSpace(username) == "" {
		return nil, errors.New("username cannot be empty")
	}
	if _, exists := s.clients[username]; exists {
		return nil, ErrUsernameAlreadyTaken
	}
	client := &Client{
		Username: username,
		Messages: make(chan string, 100),
	}
	s.clients[username] = client
	return client, nil
}

func (s *ChatServer) Disconnect(client *Client) {
	if client == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.clients[client.Username]; !exists {
		return
	}
	delete(s.clients, client.Username)
	client.mu.Lock()
	if !client.disconnected {
		client.disconnected = true
		close(client.Messages)
	}
	client.mu.Unlock()
}
func (s *ChatServer) Broadcast(sender *Client, message string) {
	if sender != nil {
		sender.mu.Lock()
		disconnected := sender.disconnected
		sender.mu.Unlock()
		if disconnected {
			return
		}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	formatted := message
	if sender != nil {
		formatted = fmt.Sprintf("%s: %s", sender.Username, message)
	}
	for _, client := range s.clients {
		client.Send(formatted)
	}
}
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	if sender != nil {
		sender.mu.Lock()
		disconnected := sender.disconnected
		sender.mu.Unlock()
		if disconnected {
			return ErrClientDisconnected
		}
	}
	s.mu.RLock()
	client, exists := s.clients[recipient]
	s.mu.RUnlock()
	if !exists {
		return ErrRecipientNotFound
	}
	formatted := message
	if sender != nil {
		formatted = fmt.Sprintf("[PM from %s] %s", sender.Username, message)
	}
	client.Send(formatted)
	return nil
}