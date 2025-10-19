# MongoDB Challenge 5: Hints & Tips

## 🎯 Progressive Implementation Strategy

### Phase 1: Foundation (Start Here)
1. **Basic Account Operations**
   ```go
   // Start with simple CRUD operations
   func (bs *BankingService) CreateAccount(ctx context.Context, userID string, initialBalance float64) Response {
       // Validate inputs first
       if userID == "" || strings.TrimSpace(userID) == "" {
           return Response{Success: false, Error: "User ID cannot be empty", Code: 400}
       }
       // Then implement database operations
   }
   ```

2. **Input Validation Pattern**
   ```go
   // Use consistent validation across all methods
   if accountID == primitive.NilObjectID {
       return Response{Success: false, Error: "Account ID cannot be empty", Code: 400}
   }
   ```

### Phase 2: Transactions (Core Feature)
3. **MongoDB Sessions**
   ```go
   session, err := bs.Client.StartSession()
   if err != nil {
       return Response{Success: false, Error: "Failed to start session", Code: 500}
   }
   defer session.EndSession(ctx)
   ```

4. **Transaction Pattern**
   ```go
   err = mongo.WithTransaction(ctx, session, func(sc mongo.SessionContext) error {
       // All database operations must use 'sc' (SessionContext)
       _, err := bs.AccountsCollection.UpdateOne(sc, filter, update)
       return err
   })
   ```

### Phase 3: Advanced Features
5. **Change Streams Setup**
   ```go
   pipeline := mongo.Pipeline{
       {{"$match", bson.M{"fullDocument._id": accountID}}},
   }
   changeStream, err := bs.AccountsCollection.Watch(ctx, pipeline)
   ```

6. **GridFS Operations**
   ```go
   bucket := gridfs.NewBucket(bs.Database)
   uploadStream, err := bucket.OpenUploadStream(filename)
   ```

## 🔧 Implementation Tips

### Transaction Best Practices
```go
// ✅ Good: Use session context in transactions
err = mongo.WithTransaction(ctx, session, func(sc mongo.SessionContext) error {
    result, err := bs.AccountsCollection.UpdateOne(sc, filter, update)
    return err
})

// ❌ Bad: Using regular context in transactions
err = mongo.WithTransaction(ctx, session, func(sc mongo.SessionContext) error {
    result, err := bs.AccountsCollection.UpdateOne(ctx, filter, update) // Wrong!
    return err
})
```

### Optimistic Locking Pattern
```go
// Include version in filter for atomic updates
filter := bson.M{
    "_id": accountID,
    "version": currentVersion,
}
update := bson.M{
    "$inc": bson.M{"balance": -amount, "version": 1},
    "$set": bson.M{"updated_at": time.Now()},
}
```

### Error Handling Strategy
```go
// Distinguish between validation and database errors
if amount <= 0 {
    return Response{Success: false, Error: "Amount must be positive", Code: 400}
}

if err != nil {
    return Response{Success: false, Error: "Database operation failed", Code: 500}
}
```

## 🚨 Common Pitfalls & Solutions

### 1. **Session Context Confusion**
```go
// ❌ Problem: Mixing contexts
func (bs *BankingService) TransferMoney(ctx context.Context, ...) Response {
    err = mongo.WithTransaction(ctx, session, func(sc mongo.SessionContext) error {
        // Using 'ctx' instead of 'sc' breaks transactions!
        _, err := bs.AccountsCollection.UpdateOne(ctx, filter, update)
        return err
    })
}

// ✅ Solution: Always use session context
func (bs *BankingService) TransferMoney(ctx context.Context, ...) Response {
    err = mongo.WithTransaction(ctx, session, func(sc mongo.SessionContext) error {
        // Use 'sc' for all database operations
        _, err := bs.AccountsCollection.UpdateOne(sc, filter, update)
        return err
    })
}
```

### 2. **Validation Order**
```go
// ✅ Validate inputs before database operations
func (bs *BankingService) TransferMoney(ctx context.Context, fromAccountID, toAccountID primitive.ObjectID, amount float64, description string) Response {
    // 1. Validate inputs first
    if amount <= 0 {
        return Response{Success: false, Error: "Amount must be positive", Code: 400}
    }
    
    if fromAccountID == toAccountID {
        return Response{Success: false, Error: "Cannot transfer to same account", Code: 400}
    }
    
    // 2. Then check database connectivity
    if bs.Client == nil {
        return Response{Success: false, Error: "Database not connected", Code: 500}
    }
    
    // 3. Finally perform operations
}
```

### 3. **Change Stream Lifecycle**
```go
// ✅ Proper change stream handling
func (bs *BankingService) StartChangeStream(ctx context.Context, accountID primitive.ObjectID) (<-chan ChangeStreamEvent, error) {
    if accountID == primitive.NilObjectID {
        return nil, fmt.Errorf("account ID cannot be empty")
    }
    
    if bs.AccountsCollection == nil {
        return nil, fmt.Errorf("accounts collection not initialized")
    }
    
    // Create pipeline and start stream
    pipeline := mongo.Pipeline{...}
    changeStream, err := bs.AccountsCollection.Watch(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    
    // Create channel and start goroutine
    eventChan := make(chan ChangeStreamEvent)
    go func() {
        defer close(eventChan)
        defer changeStream.Close(ctx)
        
        for changeStream.Next(ctx) {
            var event ChangeStreamEvent
            if err := changeStream.Decode(&event); err == nil {
                eventChan <- event
            }
        }
    }()
    
    return eventChan, nil
}
```

## 🎯 Testing Strategy

### Test Your Validation Logic
```go
// Test with nil collections to focus on validation
bankingService := &BankingService{
    Client: nil,
    AccountsCollection: nil,
    // ... other nil fields
}

// This should fail validation, not panic
response := bankingService.TransferMoney(ctx, primitive.NilObjectID, toAccountID, 100, "test")
// Should return error about invalid account ID
```

### Test Edge Cases
```go
// Test boundary conditions
response := bankingService.GetTransactionHistory(ctx, accountID, 0)     // Should fail
response := bankingService.GetTransactionHistory(ctx, accountID, 101)   // Should fail
response := bankingService.GetTransactionHistory(ctx, accountID, 100)   // Should pass validation
```

## 🚀 Performance Tips

### 1. **Index Your Collections**
```go
// Create indexes for better query performance
accountsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{"user_id", 1}},
    Options: options.Index().SetUnique(true),
})

transactionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{"from_account", 1}, {"created_at", -1}},
})
```

### 2. **Use Projections**
```go
// Only fetch needed fields
opts := options.FindOne().SetProjection(bson.M{"balance": 1, "status": 1, "version": 1})
result := bs.AccountsCollection.FindOne(ctx, filter, opts)
```

### 3. **Batch Operations**
```go
// Use bulk operations for multiple updates
models := []mongo.WriteModel{
    mongo.NewUpdateOneModel().SetFilter(filter1).SetUpdate(update1),
    mongo.NewUpdateOneModel().SetFilter(filter2).SetUpdate(update2),
}
bs.AccountsCollection.BulkWrite(ctx, models)
```

## 🔍 Debugging Tips

### 1. **Enable MongoDB Logging**
```go
// Add logging to understand what's happening
log.Printf("Starting transfer: from=%s, to=%s, amount=%.2f", fromAccountID.Hex(), toAccountID.Hex(), amount)
```

### 2. **Check Transaction Status**
```go
err = mongo.WithTransaction(ctx, session, func(sc mongo.SessionContext) error {
    log.Printf("Transaction started")
    
    result, err := bs.AccountsCollection.UpdateOne(sc, filter, update)
    if err != nil {
        log.Printf("Update failed: %v", err)
        return err
    }
    
    log.Printf("Updated %d documents", result.ModifiedCount)
    return nil
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

### 3. **Validate Your Pipeline**
```go
// Test change stream pipeline separately
pipeline := mongo.Pipeline{
    {{"$match", bson.M{"fullDocument._id": accountID}}},
}

// Test with a simple find operation first
cursor, err := bs.AccountsCollection.Aggregate(ctx, pipeline)
```

## 📚 MongoDB Driver Documentation

- **Transactions**: https://docs.mongodb.com/drivers/go/current/fundamentals/transactions/
- **Change Streams**: https://docs.mongodb.com/drivers/go/current/fundamentals/change-streams/
- **GridFS**: https://docs.mongodb.com/drivers/go/current/fundamentals/gridfs/

---

**Remember**: Start simple, validate everything, and build complexity gradually! 🎯
