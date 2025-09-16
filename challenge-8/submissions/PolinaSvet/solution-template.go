// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
// package main
package challenge8

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// *** Client ***
// ==================================================================

// Client represents a connected chat client
type Client struct {
	Username string
	Messages chan string
	mutex    sync.RWMutex
	server   *ChatServer
	conn     net.Conn
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	select {
	case c.Messages <- message:
	default:
		log.Println(c.Username, ErrClientSendDef)
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	return <-c.Messages
}

// SetConnection устанавливает соединение для клиента
func (c *Client) SetConnection(conn net.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.conn = conn
}

// GetConnection возвращает соединение клиента
func (c *Client) GetConnection() net.Conn {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.conn
}

// Common errors that can be returned by the Chat Client
var (
	ErrClientSendDef = errors.New("channel is full")
)

// *** Server ***
// ==================================================================

// Define message structures for different operations
type BroadcastMessage struct {
	Sender  *Client
	Content string
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients    map[string]*Client
	broadcast  chan BroadcastMessage
	connect    chan *Client
	disconnect chan *Client
	mutex      sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:    make(map[string]*Client),
		broadcast:  make(chan BroadcastMessage, 100),
		connect:    make(chan *Client, 100),
		disconnect: make(chan *Client, 100),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.clients[username]
	if exists {
		return nil, ErrUsernameAlreadyTaken
	}

	client := &Client{
		Username: username,
		Messages: make(chan string, 100), // buffered channel
		server:   s,
	}

	s.clients[username] = client
	s.connect <- client
	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.clients[client.Username]; exists {
		delete(s.clients, client.Username)
		close(client.Messages)

		if conn := client.GetConnection(); conn != nil {
			conn.Close()
		}

		log.Printf("Client disconnected: %s", client.Username)
	}
}

// disconnectClient безопасно отключает клиента (функция для defer)
func (s *ChatServer) disconnectClient(client *Client) {
	s.disconnect <- client
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	select {
	case s.broadcast <- BroadcastMessage{
		Sender:  sender,
		Content: message,
	}:
	default:
		log.Println("Broadcast channel full, message dropped")
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	s.mutex.RLock()
	client, existsRecipient := s.clients[recipient]
	_, existsSender := s.clients[sender.Username]
	s.mutex.RUnlock()

	if !existsRecipient {
		return ErrRecipientNotFound
	}

	if !existsSender {
		return ErrClientDisconnected
	}

	select {
	case client.Messages <- fmt.Sprintf("[PM from %s] %s", sender.Username, message):
		return nil
	default:
		return ErrRecipientMessFull
	}
}

// ListUsers returns list of connected users
func (s *ChatServer) ListUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]string, 0, len(s.clients))
	for username := range s.clients {
		users = append(users, username)
	}
	return users
}

// Run starts the server's main loop
func (s *ChatServer) Run() {
	for {
		select {
		case client := <-s.connect:
			log.Printf("Client connected: %s", client.Username)
			s.Broadcast(nil, fmt.Sprintf("*** %s joined the chat ***", client.Username))

		case client := <-s.disconnect:
			s.Disconnect(client)
			s.Broadcast(nil, fmt.Sprintf("*** %s left the chat ***", client.Username))

		case message := <-s.broadcast:
			s.mutex.RLock()
			for username, client := range s.clients {
				// Don't send to sender if it's a broadcast message
				if message.Sender == nil || username != message.Sender.Username {
					var formattedMessage string
					if message.Sender == nil {
						formattedMessage = fmt.Sprintf("*** %s ***", message.Content)
					} else {
						formattedMessage = fmt.Sprintf("[%s] %s", message.Sender.Username, message.Content)
					}

					select {
					case client.Messages <- formattedMessage:
					default:
						log.Printf("Failed to send to %s: channel full", username)
					}
				}
			}
			s.mutex.RUnlock()
		}
	}
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	ErrRecipientMessFull    = errors.New("recipient's message queue is full")
)

// handleClientConnection handles individual client connections
func handleClientConnection(conn net.Conn, server *ChatServer) {
	defer conn.Close()

	// Get username
	conn.Write([]byte("Enter your username: "))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	username := strings.TrimSpace(scanner.Text())

	if username == "" {
		conn.Write([]byte("Username cannot be empty\n"))
		return
	}

	// Connect client
	client, err := server.Connect(username)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}

	client.SetConnection(conn)
	defer server.disconnectClient(client)

	// Welcome message
	client.Send("Welcome to the chat! Commands: /users, /pm <user> <message>, /quit, /help")
	client.Send(fmt.Sprintf("Users online: %v", server.ListUsers()))

	// Message reader goroutine
	go func() {
		for message := range client.Messages {
			currentConn := client.GetConnection()
			if currentConn != nil {
				currentConn.Write([]byte(message + "\n"))
			}
		}
	}()

	// Command processor
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "/quit" {
			break
		}

		if strings.HasPrefix(text, "/") {
			// Process commands
			parts := strings.SplitN(text, " ", 3)
			switch parts[0] {
			case "/users":
				client.Send(fmt.Sprintf("Online users: %v", server.ListUsers()))
			case "/pm":
				if len(parts) < 3 {
					client.Send("Usage: /pm <username> <message>")
				} else {
					err := server.PrivateMessage(client, parts[1], parts[2])
					if err != nil {
						client.Send("Error: " + err.Error())
					} else {
						client.Send(fmt.Sprintf("PM sent to %s", parts[1]))
					}
				}
			case "/help":
				client.Send("Commands: /users, /pm <user> <message>, /quit, /help")
			default:
				client.Send("Unknown command. Type /help for help.")
			}
		} else if text != "" {
			// Broadcast message
			server.Broadcast(client, text)
		}
	}
}

// telnet localhost 8085
func main() {
	fmt.Println("Starting chat server on :8085...")

	// Create and start server
	server := NewChatServer()
	go server.Run()

	// Start TCP listener
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Server started. Waiting for connections...")
	fmt.Println("Connect using: telnet localhost 8085")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		log.Printf("New connection from: %s", conn.RemoteAddr())
		go handleClientConnection(conn, server)
	}
}
