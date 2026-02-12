package challenge8

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Client struct {
	Username string
	Messages chan string
	server   *ChatServer
}

type ChatServer struct {
	clients    map[string]*Client
	broadcast  chan BroadcastMessage
	connect    chan *Client
	disconnect chan *Client
	mutex      sync.RWMutex
}

type BroadcastMessage struct {
	Sender  *Client
	Content string
}

var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)

func NewChatServer() *ChatServer {
	s := &ChatServer{
		clients:    make(map[string]*Client),
		broadcast:  make(chan BroadcastMessage, 1000),
		connect:    make(chan *Client, 1000),
		disconnect: make(chan *Client, 1000),
	}
	go s.run()
	return s
}

func (s *ChatServer) run() {
	for {
		select {
		case client := <-s.connect:
			s.mutex.Lock()
			s.clients[client.Username] = client
			s.mutex.Unlock()
		case client := <-s.disconnect:
			s.mutex.Lock()
			s.disconnectClient(client)
			s.mutex.Unlock()
		case msg := <-s.broadcast:
			s.mutex.RLock()
			for _, client := range s.clients {
				if client.Username != msg.Sender.Username {
					client.Send(fmt.Sprintf("%s: %s", msg.Sender.Username, msg.Content))
				}
			}
			s.mutex.RUnlock()
		}
	}
}

func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mutex.RLock()
	_, exists := s.clients[username]
	s.mutex.RUnlock()
	
	if exists {
		return nil, ErrUsernameAlreadyTaken
	}
	
	client := &Client{
		Username: username,
		Messages: make(chan string, 1000),
		server:   s,
	}
	
	// Send connect request and wait a bit for it to be processed
	s.connect <- client
	time.Sleep(10 * time.Millisecond)
	
	return client, nil
}

func (s *ChatServer) disconnectClient(client *Client) {
	if _, exists := s.clients[client.Username]; exists {
		close(client.Messages)
		delete(s.clients, client.Username)
	}
}

func (s *ChatServer) Disconnect(client *Client) {
	s.disconnect <- client
	time.Sleep(10 * time.Millisecond) // Give time for disconnect to process
}

func (s *ChatServer) Broadcast(sender *Client, message string) {
	// Check if sender is still connected
	s.mutex.RLock()
	_, exists := s.clients[sender.Username]
	s.mutex.RUnlock()
	
	if exists {
		s.broadcast <- BroadcastMessage{Sender: sender, Content: message}
	}
}

func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// Check if sender is still connected
	s.mutex.RLock()
	_, senderExists := s.clients[sender.Username]
	recipientClient, recipientExists := s.clients[recipient]
	s.mutex.RUnlock()
	
	if !senderExists {
		return ErrClientDisconnected
	}
	
	if !recipientExists {
		return ErrRecipientNotFound
	}
	
	formattedMessage := fmt.Sprintf("[Private from %s]: %s", sender.Username, message)
	
	select {
	case recipientClient.Messages <- formattedMessage:
		return nil
	default:
		return errors.New("recipient's message queue is full")
	}
}

func (c *Client) Send(message string) {
	select {
	case c.Messages <- message:
	default:
		// Message dropped - queue full
	}
}

func (c *Client) Receive() string {
	msg, ok := <-c.Messages
	if !ok {
		// Channel closed
		return ""
	}
	return msg
}