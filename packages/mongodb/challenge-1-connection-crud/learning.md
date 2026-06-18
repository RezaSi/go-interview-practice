# Learning: MongoDB Fundamentals with Go

## 🌟 **What is MongoDB?**

MongoDB is a document-oriented NoSQL database that stores data in flexible, JSON-like documents called BSON (Binary JSON). Unlike traditional relational databases, MongoDB doesn't require a predefined schema.

### **Why MongoDB with Go?**
- **Natural fit**: Go structs map perfectly to MongoDB documents
- **Performance**: Efficient BSON serialization/deserialization
- **Flexibility**: Schema-less design adapts to changing requirements
- **Scalability**: Built for modern, distributed applications

## 🏗️ **Core Concepts**

### **1. Documents and Collections**
- **Document**: A record in MongoDB (like a row in SQL)
- **Collection**: A group of documents (like a table in SQL)
- **Database**: A container for collections

```go
// Document structure maps to Go struct
type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name  string             `bson:"name" json:"name"`
    Email string             `bson:"email" json:"email"`
}
```

### **2. BSON and ObjectID**
BSON extends JSON with additional types like ObjectID, dates, and binary data.

```go
// ObjectID is MongoDB's primary key type
id := primitive.NewObjectID()

// Convert string to ObjectID
objectID, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

// ObjectID contains timestamp
timestamp := id.Timestamp()
```

### **3. BSON Tags**
Control how Go structs are marshaled to/from BSON:

```go
type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty"`     // MongoDB's _id field
    Name  string             `bson:"name"`              // Maps to 'name' field
    Email string             `bson:"email_address"`     // Maps to 'email_address'
    Internal string          `bson:"-"`                 // Excluded from BSON
}
```

## 🔌 **MongoDB Connection**

### **Connection Setup**
```go
func ConnectMongoDB(uri string) (*mongo.Client, error) {
    // Create client options
    clientOptions := options.Client().ApplyURI(uri)
    
    // Connect with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }
    
    // Test connection
    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }
    
    return client, nil
}
```

### **Connection Hierarchy**
```
Client (mongo.Client)
├── Database (mongo.Database)
│   ├── Collection (mongo.Collection)
│   │   └── Documents
│   └── Indexes
└── Sessions
```

### **Getting Collections**
```go
client := // ... connected client
db := client.Database("myapp")
collection := db.Collection("users")
```

## 📄 **CRUD Operations**

### **Create (Insert)**
```go
// Insert single document
user := User{
    ID:    primitive.NewObjectID(),
    Name:  "John Doe",
    Email: "john@example.com",
}

result, err := collection.InsertOne(ctx, user)
if err != nil {
    // Handle error
}

// Insert multiple documents
users := []interface{}{user1, user2, user3}
results, err := collection.InsertMany(ctx, users)
```

### **Read (Find)**
```go
// Find one document
var user User
filter := bson.M{"_id": objectID}
err := collection.FindOne(ctx, filter).Decode(&user)
if err != nil {
    if err == mongo.ErrNoDocuments {
        // Document not found
    }
    // Handle other errors
}

// Find multiple documents
cursor, err := collection.Find(ctx, bson.M{})
if err != nil {
    // Handle error
}
defer cursor.Close(ctx)

var users []User
if err = cursor.All(ctx, &users); err != nil {
    // Handle error
}
```

### **Update**
```go
// Update one document
filter := bson.M{"_id": objectID}
update := bson.M{"$set": bson.M{"name": "New Name"}}

result, err := collection.UpdateOne(ctx, filter, update)
if err != nil {
    // Handle error
}

// Check if document was modified
if result.ModifiedCount == 0 {
    // No document was updated
}
```

### **Delete**
```go
// Delete one document
filter := bson.M{"_id": objectID}
result, err := collection.DeleteOne(ctx, filter)
if err != nil {
    // Handle error
}

// Check if document was deleted
if result.DeletedCount == 0 {
    // No document was deleted
}
```

## 🔍 **Query Filters**

### **Basic Filters**
```go
// Exact match
filter := bson.M{"name": "John Doe"}

// Multiple conditions (AND)
filter := bson.M{
    "name": "John Doe",
    "age":  30,
}

// OR conditions
filter := bson.M{
    "$or": []bson.M{
        {"name": "John"},
        {"name": "Jane"},
    },
}
```

### **Comparison Operators**
```go
// Greater than, less than
filter := bson.M{
    "age": bson.M{
        "$gte": 18,  // age >= 18
        "$lt":  65,  // age < 65
    },
}

// In array
filter := bson.M{
    "status": bson.M{"$in": []string{"active", "pending"}},
}
```

## ⚠️ **Error Handling**

### **Common MongoDB Errors**
```go
func handleMongoError(err error) {
    switch {
    case err == mongo.ErrNoDocuments:
        // Document not found
        
    case mongo.IsDuplicateKeyError(err):
        // Unique constraint violation
        
    case mongo.IsTimeout(err):
        // Operation timed out
        
    case mongo.IsNetworkError(err):
        // Network connection issue
        
    default:
        // Other database error
    }
}
```

### **Context and Timeouts**
```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Use context in operations
result, err := collection.InsertOne(ctx, document)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        // Operation timed out
    }
}
```

## 🏗️ **Best Practices**

### **Connection Management**
- Use connection pooling (automatic in Go driver)
- Always close client when done: `defer client.Disconnect(ctx)`
- Test connections with `Ping()`

### **Error Handling**
- Always check for `mongo.ErrNoDocuments`
- Handle network and timeout errors gracefully
- Use appropriate HTTP status codes in responses

### **Data Validation**
```go
func validateUser(user CreateUserRequest) error {
    if user.Name == "" {
        return errors.New("name is required")
    }
    if user.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email format")
    }
    if user.Age <= 0 {
        return errors.New("age must be positive")
    }
    return nil
}
```

### **Response Structure**
```go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
    Error   string      `json:"error,omitempty"`
    Code    int         `json:"code,omitempty"`
}

// Success response
return Response{
    Success: true,
    Data:    user,
    Message: "User created successfully",
    Code:    201,
}

// Error response
return Response{
    Success: false,
    Error:   "User not found",
    Code:    404,
}
```

## 🧪 **Testing MongoDB Applications**

### **Test Database Setup**
```go
func setupTestDB(t *testing.T) *mongo.Collection {
    client, err := mongo.Connect(context.Background(), 
        options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Skip("MongoDB not available")
    }
    
    // Use test database
    db := client.Database("test_" + t.Name())
    collection := db.Collection("users")
    
    // Cleanup after test
    t.Cleanup(func() {
        db.Drop(context.Background())
        client.Disconnect(context.Background())
    })
    
    return collection
}
```

## 🌍 **Real-World Applications**

### **When to Use MongoDB**
- **Content Management**: Blogs, news sites, documentation
- **E-commerce**: Product catalogs, user profiles, orders
- **Analytics**: Event tracking, user behavior data
- **IoT Applications**: Sensor data, device management
- **Social Networks**: User profiles, posts, relationships

### **Production Considerations**
- **Indexing**: Create indexes for frequently queried fields
- **Connection Pooling**: Configure appropriate pool sizes
- **Error Handling**: Implement retry logic for transient failures
- **Monitoring**: Track query performance and errors
- **Security**: Use authentication and validate all inputs

## 📚 **Next Steps**

After mastering basic CRUD operations, explore:
1. **Advanced Queries**: Complex filters, sorting, pagination
2. **Aggregation Pipeline**: Data processing and analytics
3. **Indexing**: Performance optimization
4. **Transactions**: Multi-document operations
5. **Change Streams**: Real-time data monitoring

## 🔗 **Additional Resources**

- [MongoDB Go Driver Documentation](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [MongoDB Manual](https://docs.mongodb.com/manual/)
- [BSON Specification](http://bsonspec.org/)
- [MongoDB University](https://university.mongodb.com/) - Free courses