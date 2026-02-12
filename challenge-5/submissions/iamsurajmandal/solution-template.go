package main

import (
    "context"
	"fmt"
	"net/http"
)

const validToken = "secret"

// AuthMiddleware checks the "X-Auth-Token" header.
// If it's "secret", call the next handler.
// Otherwise, respond with 401 Unauthorized.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    auth := r.Header.Get("X-Auth-Token")
	    fmt.Println("hey this token",auth)
	    if (auth == validToken) {
	        ctx := context.WithValue(r.Context(), validToken, auth)
		    next.ServeHTTP(w, r.WithContext(ctx)) // calling next handleer
	    } else {
	       http.Error(w, "", http.StatusUnauthorized)
	       return;
	    }
		// TODO: Implement the logic:
		//  1) Grab the "X-Auth-Token" header
		//  2) Compare against validToken
		//  3) If mismatch or missing, respond with 401
		//  4) Otherwise pass to next handler
	})
}

// helloHandler returns "Hello!" on GET /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

// secureHandler returns "You are authorized!" on GET /secure
func secureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You are authorized!")
}

// SetupServer configures the HTTP routes with the authentication middleware.
func SetupServer() http.Handler {
	mux := http.NewServeMux()

	// Public route: /hello (no auth required)
	mux.HandleFunc("/hello", helloHandler)

	// Secure route: /secure
	// Wrap with AuthMiddleware
	secureRoute := http.HandlerFunc(secureHandler)
	mux.Handle("/secure", AuthMiddleware(secureRoute))

	return mux
}

func main() {
	// Optional: you can run a real server for local testing
	// http.ListenAndServe(":8080", SetupServer())
}
