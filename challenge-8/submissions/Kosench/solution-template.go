// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"fmt"
	"sync"
	"time"
	// Add any other necessary imports
)

type Message struct {
	Sender    *Client
	Content   string
	Recipient string // пустая строка для broadcast
}

// Client represents a connected chat client
type Client struct {
	username     string
	messages     chan string
	server       *ChatServer
	mu           sync.Mutex
	disconnected bool
}

func newClient(username string, server *ChatServer) *Client {
	return &Client{
		username:     username,
		messages:     make(chan string, 50),
		server:       server,
		disconnected: false,
	}
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.mu.Lock()
	disconnected := c.disconnected
	c.mu.Unlock()

	if disconnected {
		return
	}

	// Non-blocking send используя select с default [web:49][web:53]
	select {
	case c.messages <- message:
		// Сообщение отправлено успешно
	default:
		// Канал переполнен, пропускаем сообщение
		// В production здесь можно логировать
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// Blocking read с проверкой закрытого канала [web:7][web:56]
	message, ok := <-c.messages
	if !ok {
		// Канал закрыт
		return ""
	}
	return message
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients   map[string]*Client
	broadcast chan Message
	join      chan *Client
	leave     chan *Client
	mu        sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	server := &ChatServer{
		clients:   make(map[string]*Client),
		broadcast: make(chan Message, 100),
		join:      make(chan *Client, 10),
		leave:     make(chan *Client, 10),
	}

	go server.run()

	return server
}

func (s *ChatServer) run() {
	for {
		select {
		case client := <-s.join:
			s.mu.Lock()
			s.clients[client.username] = client
			s.mu.Unlock()

		case client := <-s.leave:
			s.mu.Lock()
			if _, exists := s.clients[client.username]; exists {
				delete(s.clients, client.username)
				close(client.messages)
			}
			s.mu.Unlock()

		case msg := <-s.broadcast:
			s.mu.RLock()
			if msg.Recipient == "" {
				for _, client := range s.clients {
					if client != msg.Sender {
						client.Send(msg.Content)
					}
				}
			} else {
				if recipient, exists := s.clients[msg.Recipient]; exists {
					recipient.Send(msg.Content)
				}
			}
			s.mu.RUnlock()
		}
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}

	s.mu.RLock()
	_, exist := s.clients[username]
	s.mu.RUnlock()

	if exist {
		return nil, ErrUsernameAlreadyTaken
	}

	client := newClient(username, s)
	s.join <- client

	// Даём время центральной горутине обработать join
	time.Sleep(10 * time.Millisecond)

	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	if client == nil {
		return
	}

	client.mu.Lock()
	client.disconnected = true
	client.mu.Unlock()

	s.leave <- client

	time.Sleep(10 * time.Millisecond)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	if sender == nil {
		return
	}

	sender.mu.Lock()
	disconnected := sender.disconnected
	sender.mu.Unlock()

	if disconnected {
		return
	}

	formattedMsg := fmt.Sprintf("[%s]: %s", sender.username, message)

	msg := Message{
		Sender:    sender,
		Content:   formattedMsg,
		Recipient: "",
	}

	s.broadcast <- msg

}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	if sender == nil {
		return ErrClientDisconnected
	}

	// Проверяем, что отправитель не отключён [web:7]
	sender.mu.Lock()
	disconnected := sender.disconnected
	sender.mu.Unlock()

	if disconnected {
		return ErrClientDisconnected
	}

	// Проверяем существование получателя
	s.mu.RLock()
	_, exists := s.clients[recipient]
	s.mu.RUnlock()

	if !exists {
		return ErrRecipientNotFound
	}

	// Форматируем сообщение (важно: содержит оригинальный message для теста strings.Contains)
	formattedMsg := fmt.Sprintf("[PM from %s]: %s", sender.username, message)

	// Отправляем в канал broadcast с указанием получателя
	msg := Message{
		Sender:    sender,
		Content:   formattedMsg,
		Recipient: recipient,
	}

	s.broadcast <- msg

	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	ErrEmptyUsername        = errors.New("username cannot be empty")
)
