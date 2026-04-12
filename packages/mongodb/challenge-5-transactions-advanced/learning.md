# MongoDB Challenge 5: Learning Guide

## 🎓 Advanced MongoDB Concepts

### 1. ACID Transactions in MongoDB

#### What are ACID Transactions?
```
A - Atomicity: All operations succeed or all fail
C - Consistency: Data remains in valid state
I - Isolation: Concurrent transactions don't interfere
D - Durability: Committed changes persist
```

#### MongoDB Transaction Implementation
```go
// Basic transaction pattern
session, err := client.StartSession()
defer session.EndSession(ctx)

err = mongo.WithTransaction(ctx, session, func(sc mongo.SessionContext) error {
    // All operations must use SessionContext (sc)
    _, err := collection1.UpdateOne(sc, filter1, update1)
    if err != nil {
        return err // Triggers rollback
    }
    
    _, err = collection2.InsertOne(sc, document)
    if err != nil {
        return err // Triggers rollback
    }
    
    return nil // Commits transaction
})
```

#### Why Use Transactions?
- **Financial Operations**: Money transfers must be atomic
- **Data Consistency**: Related documents stay synchronized
- **Error Recovery**: Automatic rollback on failures

### 2. Change Streams - Real-time Data Monitoring

#### What are Change Streams?
Change streams provide real-time notifications when data changes in MongoDB collections.

```go
// Watch for changes in accounts collection
pipeline := mongo.Pipeline{
    {{"$match", bson.M{"operationType": "update"}}},
    {{"$match", bson.M{"fullDocument._id": accountID}}},
}

changeStream, err := collection.Watch(ctx, pipeline)
defer changeStream.Close(ctx)

for changeStream.Next(ctx) {
    var event bson.M
    changeStream.Decode(&event)
    
    // Process the change event
    fmt.Printf("Account %s was updated\n", accountID)
}
```

#### Use Cases
- **Real-time Dashboards**: Update UI when balances change
- **Fraud Detection**: Monitor suspicious transaction patterns
- **Audit Logging**: Track all data modifications
- **Cache Invalidation**: Update caches when data changes

### 3. GridFS - Large File Storage

#### What is GridFS?
GridFS is MongoDB's specification for storing large files that exceed the 16MB document size limit.

```go
// Store a large document
bucket := gridfs.NewBucket(database)

uploadStream, err := bucket.OpenUploadStream("contract.pdf")
defer uploadStream.Close()

// Write file data
_, err = uploadStream.Write(fileData)

// Get the file ID
fileID := uploadStream.FileID
```

#### GridFS Architecture
- **fs.files**: Stores file metadata
- **fs.chunks**: Stores file data in 255KB chunks
- **Automatic Sharding**: Large files split across chunks

#### Use Cases
- **Document Storage**: Legal contracts, reports
- **Media Files**: Images, videos, audio
- **Backup Files**: Database dumps, logs
- **Binary Data**: Any large binary content

### 4. Optimistic Locking

#### The Concurrency Problem
```go
// Two users try to transfer money simultaneously
// User A: Transfer $100 from Account X (balance: $500)
// User B: Transfer $200 from Account X (balance: $500)
// Without locking: Both see $500, both succeed, balance becomes wrong!
```

#### Solution: Version-based Locking
```go
type Account struct {
    ID      primitive.ObjectID `bson:"_id"`
    Balance float64            `bson:"balance"`
    Version int64              `bson:"version"` // Key field!
}

// Update with version check
filter := bson.M{
    "_id": accountID,
    "version": currentVersion, // Must match current version
}

update := bson.M{
    "$inc": bson.M{
        "balance": -amount,
        "version": 1, // Increment version
    },
}

result, err := collection.UpdateOne(ctx, filter, update)
if result.ModifiedCount == 0 {
    return errors.New("concurrent modification detected")
}
```

### 5. Error Handling Strategies

#### Categorize Errors
```go
func handleError(err error) Response {
    if err == nil {
        return Response{Success: true}
    }
    
    // MongoDB specific errors
    if mongo.IsDuplicateKeyError(err) {
        return Response{Success: false, Error: "Duplicate key", Code: 409}
    }
    
    if mongo.IsTimeout(err) {
        return Response{Success: false, Error: "Operation timeout", Code: 408}
    }
    
    // Network errors (retry-able)
    if mongo.IsNetworkError(err) {
        return Response{Success: false, Error: "Network error", Code: 503}
    }
    
    // Generic error
    return Response{Success: false, Error: "Internal error", Code: 500}
}
```

#### Retry Logic with Exponential Backoff
```go
func retryOperation(operation func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        // Don't retry validation errors
        if isValidationError(err) {
            return err
        }
        
        // Exponential backoff: 100ms, 200ms, 400ms, 800ms...
        backoff := time.Duration(100 * (1 << attempt)) * time.Millisecond
        time.Sleep(backoff)
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}
```

### 6. Audit Logging for Compliance

#### Why Audit Logging?
- **Regulatory Compliance**: Financial regulations require audit trails
- **Security**: Track unauthorized access attempts
- **Debugging**: Understand what happened when issues occur
- **Business Intelligence**: Analyze user behavior patterns

#### Implementation Pattern
```go
type AuditLog struct {
    ID        primitive.ObjectID     `bson:"_id,omitempty"`
    Action    string                 `bson:"action"`
    UserID    string                 `bson:"user_id"`
    AccountID primitive.ObjectID     `bson:"account_id,omitempty"`
    Details   map[string]interface{} `bson:"details"`
    Timestamp time.Time              `bson:"timestamp"`
    IPAddress string                 `bson:"ip_address,omitempty"`
}

func (bs *BankingService) logAuditEvent(action, userID string, accountID primitive.ObjectID, details map[string]interface{}) {
    auditLog := AuditLog{
        Action:    action,
        UserID:    userID,
        AccountID: accountID,
        Details:   details,
        Timestamp: time.Now(),
    }
    
    // Log asynchronously to avoid blocking main operations
    go func() {
        bs.AuditCollection.InsertOne(context.Background(), auditLog)
    }()
}

// Usage
bs.logAuditEvent("transfer_money", userID, fromAccountID, map[string]interface{}{
    "to_account": toAccountID,
    "amount":     amount,
    "description": description,
})
```

### 7. Production Considerations

#### Connection Configuration
```go
// Production-ready connection options
clientOptions := options.Client().
    ApplyURI(uri).
    SetMaxPoolSize(100).                    // Connection pool size
    SetMinPoolSize(10).                     // Minimum connections
    SetMaxConnIdleTime(30 * time.Second).   // Idle connection timeout
    SetServerSelectionTimeout(5 * time.Second). // Server selection timeout
    SetSocketTimeout(10 * time.Second).     // Socket timeout
    SetRetryWrites(true).                   // Enable retryable writes
    SetRetryReads(true)                     // Enable retryable reads

client, err := mongo.Connect(ctx, clientOptions)
```

#### Read/Write Concerns
```go
// Ensure data durability
writeOptions := options.Collection().
    SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
    SetReadConcern(readconcern.Majority())

collection := database.Collection("accounts", writeOptions)
```

#### Monitoring and Metrics
```go
// Track operation performance
start := time.Now()
result, err := collection.UpdateOne(ctx, filter, update)
duration := time.Since(start)

// Log slow operations
if duration > 100*time.Millisecond {
    log.Printf("Slow operation: %v took %v", operation, duration)
}

// Track success/failure rates
if err != nil {
    metrics.IncrementCounter("mongodb.operations.failed")
} else {
    metrics.IncrementCounter("mongodb.operations.succeeded")
}
```

## 🏗️ Architecture Patterns

### 1. Repository Pattern
```go
type AccountRepository interface {
    Create(ctx context.Context, account Account) error
    GetByID(ctx context.Context, id primitive.ObjectID) (*Account, error)
    UpdateBalance(ctx context.Context, id primitive.ObjectID, amount float64) error
}

type MongoAccountRepository struct {
    collection *mongo.Collection
}

func (r *MongoAccountRepository) UpdateBalance(ctx context.Context, id primitive.ObjectID, amount float64) error {
    filter := bson.M{"_id": id}
    update := bson.M{"$inc": bson.M{"balance": amount}}
    
    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}
```

### 2. Service Layer Pattern
```go
type BankingService struct {
    accountRepo     AccountRepository
    transactionRepo TransactionRepository
    auditRepo       AuditRepository
}

func (s *BankingService) TransferMoney(ctx context.Context, req TransferRequest) error {
    // 1. Validate request
    if err := s.validateTransferRequest(req); err != nil {
        return err
    }
    
    // 2. Start transaction
    return s.executeTransfer(ctx, req)
}
```

## 🎯 Real-World Applications

### Banking Systems
- **Core Banking**: Account management, transfers, loans
- **Payment Processing**: Credit card transactions, digital wallets
- **Risk Management**: Fraud detection, compliance monitoring

### Enterprise Applications
- **ERP Systems**: Inventory management, order processing
- **CRM Systems**: Customer data, interaction tracking
- **Content Management**: Document storage, version control

### E-commerce Platforms
- **Order Processing**: Multi-step checkout workflows
- **Inventory Management**: Stock level tracking
- **Payment Systems**: Secure transaction processing

## 📊 Performance Optimization

### Indexing Strategy
```go
// Compound indexes for common queries
db.accounts.createIndex({"user_id": 1, "status": 1})
db.transactions.createIndex({"from_account": 1, "created_at": -1})
db.transactions.createIndex({"to_account": 1, "created_at": -1})

// Text indexes for search
db.transactions.createIndex({"description": "text"})
```

### Query Optimization
```go
// Use projections to limit data transfer
opts := options.FindOne().SetProjection(bson.M{
    "balance": 1,
    "status": 1,
    "version": 1,
})

// Use limits for large result sets
opts := options.Find().SetLimit(100).SetSort(bson.M{"created_at": -1})
```

---

**Master these concepts and you'll be ready to build enterprise-grade applications with MongoDB!** 🚀
