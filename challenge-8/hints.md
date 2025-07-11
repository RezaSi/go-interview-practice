# Hints for Chat Server with Channels

## Hint 1: ChatServer Structure
Design your ChatServer with channels for coordination and a map to track clients:
```go
type ChatServer struct {
    clients     map[string]*Client
    broadcast   chan BroadcastMessage
    connect     chan *Client
    disconnect  chan *Client
    mutex       sync.RWMutex
}
```

## Hint 2: Client Structure
Each client needs channels for communication and identifying information:
```go
type Client struct {
    Username string
    Messages chan string
    server   *ChatServer
}
```

## Hint 3: Message Types
Define message structures for different operations:
```go
type BroadcastMessage struct {
    Sender  *Client
    Content string
}
```

## Hint 4: Server Event Loop
The server should run a goroutine that handles all operations through channels:
```go
func (s *ChatServer) run() {
    for {
        select {
        case client := <-s.connect:
            // Handle new connection
        case client := <-s.disconnect:
            // Handle disconnection
        case msg := <-s.broadcast:
            // Handle broadcast message
        }
    }
}
```

## Hint 5: Thread-Safe Client Management
Use a mutex when accessing the clients map:
```go
s.mutex.Lock()
s.clients[client.Username] = client
s.mutex.Unlock()
```

## Hint 6: Connect Method Implementation
Create a new client and send it through the connect channel:
```go
func (s *ChatServer) Connect(username string) (*Client, error) {
    if /* username already exists */ {
        return nil, errors.New("username already taken")
    }
    
    client := &Client{
        Username: username,
        Messages: make(chan string, 100), // buffered channel
        server:   s,
    }
    
    s.connect <- client
    return client, nil
}
```

## Hint 7: Broadcast Implementation
Send the message through the broadcast channel:
```go
func (s *ChatServer) Broadcast(sender *Client, message string) {
    s.broadcast <- BroadcastMessage{
        Sender:  sender,
        Content: message,
    }
}
```

## Hint 8: Private Message Implementation
Find the recipient and send the message directly to their channel:
```go
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
    s.mutex.RLock()
    client, exists := s.clients[recipient]
    s.mutex.RUnlock()
    
    if !exists {
        return errors.New("recipient not found")
    }
    
    select {
    case client.Messages <- message:
        return nil
    default:
        return errors.New("recipient's message queue is full")
    }
}
```

## Hint 9: Client Send and Receive Methods
```go
func (c *Client) Send(message string) {
    select {
    case c.Messages <- message:
    default:
        // Channel is full, handle gracefully
    }
}

func (c *Client) Receive() string {
    return <-c.Messages
}
```

## Hint 10: Graceful Shutdown
When disconnecting, clean up resources and close the client's message channel:
```go
close(client.Messages)
delete(s.clients, client.Username)
``` 