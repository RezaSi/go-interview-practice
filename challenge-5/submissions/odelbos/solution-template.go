package main

import (
	"fmt"
	"net/http"
	"crypto/subtle"
)

const validToken = "secret"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("X-Auth-Token")
        // Crypto constant time compare to avoid timing attack.
        if subtle.ConstantTimeCompare([]byte(token), []byte(validToken)) == 1 {
            next.ServeHTTP(w, r)
            return
        }
        http.Error(w, "", http.StatusUnauthorized)
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
