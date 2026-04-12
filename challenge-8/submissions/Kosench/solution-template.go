package challenge8

import (
	"errors"
	"fmt"
	"sync"
)

// Message represents a message to be delivered
type Message struct {
	Sender    *Client
	Content   string
	Recipient string // empty for broadcast
}

// joinRequest represents a request to join the chat
type joinRequest struct {
	username string
	response chan *Client
	errChan  chan error
}

// leaveRequest represents a request to leave the chat
type leaveRequest struct {
	client *Client
	done   chan struct{}
}

// Client represents a connected chat client
type Client struct {
	username string
	messages chan string
	server   *ChatServer
	mu       sync.RWMutex
	active   bool
}

func newClient(username string, server *ChatServer) *Client {
	return &Client{
		username: username,
		messages: make(chan string, 50),
		server:   server,
		active:   true,
	}
}

// Send sends a message to the client (non-blocking)
func (c *Client) Send(message string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.active {
		return
	}

	// Non-blocking send (protected by RLock to prevent close during send)
	select {
	case c.messages <- message:
	default:
		// Channel full, drop message
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	c.mu.RLock()
	msgChan := c.messages
	c.mu.RUnlock()

	msg, ok := <-msgChan
	if !ok {
		return ""
	}
	return msg
}

// isActive checks if client is still active
func (c *Client) isActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.active
}

// markInactive marks client as inactive and closes channel
func (c *Client) markInactive() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.active {
		c.active = false
		close(c.messages)
	}
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients   map[string]*Client
	mu        sync.RWMutex // Protects clients map (read and write)
	broadcast chan Message
	join      chan joinRequest
	leave     chan leaveRequest
	shutdown  chan struct{}
	wg        sync.WaitGroup
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	server := &ChatServer{
		clients:   make(map[string]*Client),
		broadcast: make(chan Message, 100),
		join:      make(chan joinRequest),
		leave:     make(chan leaveRequest),
		shutdown:  make(chan struct{}),
	}

	server.wg.Add(1)
	go server.run()

	return server
}

// run is the central goroutine that handles all state modifications
func (s *ChatServer) run() {
	defer s.wg.Done()

	for {
		select {
		case req := <-s.join:
			// Check for duplicate username and register client under lock
			s.mu.Lock()
			if _, exists := s.clients[req.username]; exists {
				s.mu.Unlock()
				req.errChan <- ErrUsernameAlreadyTaken
				close(req.response)
				close(req.errChan)
				continue
			}

			client := newClient(req.username, s)
			s.clients[req.username] = client
			s.mu.Unlock()

			// Send response
			req.response <- client
			req.errChan <- nil
			close(req.response)
			close(req.errChan)

		case req := <-s.leave:
			// Remove client if exists
			s.mu.Lock()
			client, exists := s.clients[req.client.username]
			if exists {
				delete(s.clients, req.client.username)
			}
			s.mu.Unlock()

			if exists {
				client.markInactive()
			}
			close(req.done)

		case msg := <-s.broadcast:
			// Deliver message (read clients map under lock)
			s.mu.RLock()
			if msg.Recipient == "" {
				// Broadcast to all except sender
				for _, client := range s.clients {
					if client != msg.Sender && client.isActive() {
						client.Send(msg.Content)
					}
				}
			} else {
				// Private message
				if recipient, exists := s.clients[msg.Recipient]; exists && recipient.isActive() {
					recipient.Send(msg.Content)
				}
			}
			s.mu.RUnlock()

		case <-s.shutdown:
			// Cleanup all clients
			s.mu.RLock()
			clientsCopy := make([]*Client, 0, len(s.clients))
			for _, client := range s.clients {
				clientsCopy = append(clientsCopy, client)
			}
			s.mu.RUnlock()

			// Mark inactive outside the lock to avoid deadlock
			for _, client := range clientsCopy {
				client.markInactive()
			}

			s.mu.Lock()
			s.clients = make(map[string]*Client)
			s.mu.Unlock()
			return
		}
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}

	// Create request channels
	req := joinRequest{
		username: username,
		response: make(chan *Client, 1),
		errChan:  make(chan error, 1),
	}

	// Send join request to central goroutine
	s.join <- req

	// Wait for response (blocking until processed)
	client := <-req.response
	err := <-req.errChan

	if err != nil {
		return nil, err
	}

	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	if client == nil {
		return
	}

	// Create request with done channel
	req := leaveRequest{
		client: client,
		done:   make(chan struct{}),
	}

	// Send leave request
	s.leave <- req

	// Wait for completion (blocking until processed)
	<-req.done
}

// Broadcast sends a message to all connected clients except sender
func (s *ChatServer) Broadcast(sender *Client, message string) {
	if sender == nil || !sender.isActive() {
		return
	}

	formattedMsg := fmt.Sprintf("[%s]: %s", sender.username, message)

	msg := Message{
		Sender:    sender,
		Content:   formattedMsg,
		Recipient: "",
	}

	select {
	case s.broadcast <- msg:
	default:
		// Broadcast channel full, skip
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipientUsername string, message string) error {
	if sender == nil || !sender.isActive() {
		return ErrClientDisconnected
	}

	// Check if recipient exists (under read lock)
	s.mu.RLock()
	_, exists := s.clients[recipientUsername]
	s.mu.RUnlock()

	if !exists {
		return ErrRecipientNotFound
	}

	formattedMsg := fmt.Sprintf("[PM from %s]: %s", sender.username, message)

	msg := Message{
		Sender:    sender,
		Content:   formattedMsg,
		Recipient: recipientUsername,
	}

	select {
	case s.broadcast <- msg:
		return nil
	default:
		return errors.New("message queue full")
	}
}

// Shutdown gracefully shuts down the chat server
func (s *ChatServer) Shutdown() {
	close(s.shutdown)
	s.wg.Wait()
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	ErrEmptyUsername        = errors.New("username cannot be empty")
)
