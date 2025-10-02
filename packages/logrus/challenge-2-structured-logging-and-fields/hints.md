# Hints for Challenge 2: Structured Logging & Fields

## Hint 1: Generating the Request ID

To generate a unique ID for each request, you'll need the `github.com/google/uuid` package. The `uuid.New()` function creates a new UUID object, and its `.String()` method returns it in the standard format

```go
import "github.com/google/uuid"

// In the loggingMiddleware...
requestID := uuid.New().String()
```

---

## Hint 2: Creating the Context-Aware Logger

The goal is to create a `logrus.Entry` that is pre-loaded with fields. Use `logrus.WithFields()` and pass it a `logrus.Fields` map. You can get the required information directly from the `*http.Request` object (`r`)

```go
// In the loggingMiddleware...
logger := logrus.WithFields(logrus.Fields{
    "request_id":  requestID,
    "http_method": r.Method,
    "uri":         r.RequestURI,
    "user_agent":  r.UserAgent(),
})

// Now, any log call using this `logger` variable will include these fields
logger.Info("Request received")
```

---

## Hint 3: Putting the Logger into the Context

To pass the logger to the handler, you need to add it to the request's context. This is a two-step process:

1. Create a new context from the old one, adding your value
2. Create a new request object that uses this new context

```go
// In the loggingMiddleware...

// 1. Create a new context with the logger stored under our custom `key`.\
ctx := context.WithValue(r.Context(), key, logger)

// 2. Call the next handler, but replace the request `r` with a copy
// that has the new context
next.ServeHTTP(w, r.WithContext(ctx))
```

---

## Hint 4: Retrieving the Logger in the Handler

In the `helloHandler`, you need to get the logger back out of the context. The `Value()` method returns an `interface{}`, so you must use a type assertion to convert it back to a `*logrus.Entry`. Always check if the assertion was successful with the `ok` variable

```go
// In the helloHandler...
var logger *logrus.Entry

// Try to get the logger from context.
loggerFromCtx, ok := r.Context().Value(key).(*logrus.Entry)
if ok {
    // Success! Use the logger from the context
    logger = loggerFromCtx
} else {
    // Fallback: If no logger is in the context, use the default global one
    logger = logrus.NewEntry(logrus.StandardLogger())
}
```

---

## Hint 5: Adding More Fields

Once you have a `*logrus.Entry`, you can chain more `.WithField()` or `.WithFields()` calls to it. This creates a new entry with the combined fields.

```go
// In the helloHandler, after retrieving the logger...
logger = logger.WithField("user_id", "user-99")

// Now, this log will have the middleware fields AND the user_id field
logger.Info("Processing hello request")
```

---

## Hint 6: Setting Up the Server in `main`

Don't forget to configure the global logger in `main`. Then, chain your handlers together. The request will flow through the middleware first, and then to the final handler.

```go
func main() {
    // Configure the global logger
    logrus.SetFormatter(&logrus.JSONFormatter{})

    // Wrap the final handler with the middleware
    finalHandler := loggingMiddleware(http.HandlerFunc(helloHandler))

    // Start the server with the wrapped handler
    http.ListenAndServe(":8080", finalHandler)
}
```

---