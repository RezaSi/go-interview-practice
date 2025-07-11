# Hints for HTTP Authentication Middleware

## Hint 1: Middleware Function Signature
HTTP middleware in Go is a function that takes an `http.Handler` and returns an `http.Handler`:
```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // middleware logic here
    })
}
```

## Hint 2: Getting Header Values
Extract the authentication token from the request header:
```go
token := r.Header.Get("X-Auth-Token")
```

## Hint 3: Token Validation
Check if the token equals the expected secret value. If not, return 401:
```go
if token != "secret" {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}
```

## Hint 4: Calling Next Handler
If the token is valid, call the next handler in the chain:
```go
next.ServeHTTP(w, r)
```

## Hint 5: Setting Up Routes
Create a router with the required endpoints:
```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", helloHandler)
mux.HandleFunc("/secure", secureHandler)
```

## Hint 6: Applying Middleware
Wrap your router with the authentication middleware:
```go
server := &http.Server{
    Addr:    ":8080",
    Handler: authMiddleware(mux),
}
```

## Hint 7: Handler Functions
Create simple handler functions:
```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello!"))
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("You are authorized!"))
}
```

## Hint 8: Missing Header Handling
If the header is missing (empty string), treat it as invalid and return 401. The `Header.Get()` method returns an empty string if the header doesn't exist. 