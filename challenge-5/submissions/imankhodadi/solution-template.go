package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You are authorized!")
}

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
	http.ListenAndServe(":8080", SetupServer())
}

const validToken = "secret"

func validateToken(token string) (bool, error) {
	if token == validToken {
		return true, nil
	}
	return false, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		token := r.Header.Get("X-Auth-Token") // Directly get the value of X-Auth-Token

		if token == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Validate token using your existing function
		isValid, _ := validateToken(token) // We only care about the boolean result for this challenge

		if !isValid {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// If valid, pass the request to the next handler
		// For this challenge, you don't strictly need to pass user info via context
		// as the secureHandler doesn't consume it.
		next.ServeHTTP(w, r)
	})
}