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
	username = strings.TrimSpace(username)
	if username == "" {
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
	stored, exists := s.clients[client.Username]
	if !exists || stored != client {
		return
	}
	delete(s.clients, client.Username)
	stored.mu.Lock()
	if !stored.disconnected {
		stored.disconnected = true
		close(stored.Messages)
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

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	// Create a scanner to read messages line by line
// 	scanner := bufio.NewScanner(conn)

// 	// Read username
// 	var username string
// 	if scanner.Scan() {
// 		username = scanner.Text()
// 	}

// 	// Add client to chat
// 	client := &Client{
// 		Conn:     conn,
// 		Username: username,
// 		Outgoing: make(chan string),
// 	}
// 	chat.Join(client)

// 	// Set up bidirectional communication
// 	go client.ReadMessages(scanner, chat)
// 	client.WriteMessages()
// }
// func (s *ChatServer) run() {
// 	for {
// 		select {
// 		case client := <-s.connect:
// 			// Handle new connection
// 		case client := <-s.disconnect:
// 			// Handle disconnection
// 		case msg := <-s.broadcast:
// 			// Handle broadcast message
// 		}
// 	}
// 	/*
// 	   for {
// 	       select {
// 	       case client := <-c.join:
// 	           c.mu.Lock()
// 	           c.clients[client] = true
// 	           c.mu.Unlock()
// 	           c.broadcast <- fmt.Sprintf("%s has joined the chat", client.Username)

// 	       case client := <-c.leave:
// 	           c.mu.Lock()
// 	           delete(c.clients, client)
// 	           c.mu.Unlock()
// 	           close(client.Outgoing)
// 	           c.broadcast <- fmt.Sprintf("%s has left the chat", client.Username)

// 	       case message := <-c.broadcast:
// 	           c.mu.Lock()
// 	           for client := range c.clients {
// 	               select {
// 	               case client.Outgoing <- message:
// 	                   // Message sent successfully
// 	               default:
// 	                   // Client buffer is full, remove them
// 	                   delete(c.clients, client)
// 	                   close(client.Outgoing)
// 	               }
// 	           }
// 	           c.mu.Unlock()
// 	       }
// 	   }
// 	*/
// }

// // WriteMessages sends messages from the chat to the client
// func (c *Client) WriteMessages() {
// 	for message := range c.Outgoing {
// 		fmt.Fprintln(c.Conn, message)
// 	}
// }

// func processMessages(input <-chan string, workers int) <-chan string {
// 	output := make(chan string)

// 	// Fan out to workers
// 	for i := 0; i < workers; i++ {
// 		go func() {
// 			for message := range input {
// 				// Process message (e.g., check for commands, filter bad words)
// 				processed := processMessage(message)
// 				output <- processed
// 			}
// 		}()
// 	}

// 	return output
// }

// func processMessage(message string) string {
// 	// Apply formatting, filtering, etc.
// 	return message
// }

// // ReadMessages reads messages from the client and sends them to the chat
// func (c *Client) ReadMessages(scanner *bufio.Scanner, chat *Chat) {
// 	defer func() {
// 		chat.leave <- c
// 	}()

// 	for scanner.Scan() {
// 		message := scanner.Text()
// 		if message == "/quit" {
// 			break
// 		}

// 		chat.broadcast <- fmt.Sprintf("%s: %s", c.Username, message)
// 	}
// 	// Check for error
// 	if err := scanner.Err(); err != nil {
// 		log.Printf("Error reading from %s: %v", c.Username, err)
// 	}
// }

// func acceptConnections() {
// 	// Accept new connections
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Println("Error accepting connection:", err)
// 			continue
// 		}

// 		// Handle each connection in a goroutine
// 		go handleConnection(conn)
// 	}
// }
// func main() {
// 	// Set up signal handling
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

// 	// Start chat server
// 	chat := NewChat()
// 	go chat.Run()

// 	// Start TCP server
// 	listener, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Accept connections in a separate goroutine
// 	go acceptConnections(listener, chat)

// 	// Wait for termination signal
// 	<-stop
// 	log.Println("Shutting down server...")

// 	// Notify all clients
// 	chat.broadcast <- "SERVER: Chat server is shutting down"

// 	// Give clients time to disconnect
// 	time.Sleep(2 * time.Second)
// 	listener.Close()
// }

/*

Hint 5
Use a mutex when accessing the clients map:

s.mutex.Lock()
s.clients[client.Username] = client
s.mutex.Unlock()
---------------------------------------------------



Non-blocking Channel Operations with Select
The select statement allows you to wait on multiple channel operations:

select {
case message := <-messageChan:
    // Handle message
case client := <-joinChan:
    // Handle new client
case <-time.After(30 * time.Second):
    // Handle timeout
default:
    // Non-blocking path (only executes if no other case is ready)
}
Timeouts and Deadlines
It's important to handle timeouts to prevent blocked connections:

// Set read deadline
conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

// Set write deadline
conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

// Context timeout for graceful shutdown
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Shutdown(ctx)
Fan-out/Fan-in Pattern for Message Processing
For advanced message handling, you might process messages concurrently:

Buffered vs Unbuffered Channels
Choose channel type based on your requirements:

// Unbuffered channel - blocks until receiver is ready
unbuffered := make(chan string)

// Buffered channel - only blocks when buffer is full
buffered := make(chan string, 10)
Considerations:

Unbuffered channels provide synchronization points
Buffered channels allow sender to continue before receiver is ready
Buffer size should be based on expected load patterns
Handling Client Disconnects
Gracefully handle clients that disconnect unexpectedly:

Rate Limiting
Prevent clients from flooding the chat with messages:

// Create a rate limiter that allows 5 messages per second
limiter := rate.NewLimiter(5, 10)

// Use in a client handler
if !limiter.Allow() {
    // Rate limit exceeded
    fmt.Fprintln(conn, "Rate limit exceeded. Please slow down.")
    continue
}
Logging and Monitoring
Add logging to track server activity:

// Structure logs with context
log.Printf("Client connected: %s (%s)", client.Username, conn.RemoteAddr())
log.Printf("Broadcast: %s", message)
log.Printf("Client disconnected: %s after %v", client.Username, time.Since(client.ConnectedAt))

// Count active connections
log.Printf("Active connections: %d", len(chat.clients))
Graceful Shutdown
Implement proper server shutdown to avoid dropping connections:

Best Practices for Chat Servers
Use context for cancelation: Propagate cancelation through your application
Implement health checks: Monitor server health and connection status
Add authentication: Verify users before allowing them to join
Use heartbeats: Detect disconnected clients that don't close properly
Handle backpressure: Deal with slow clients to prevent memory issues
Add metrics: Track message volume, user counts, and error rates
Further Reading
*/