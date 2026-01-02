package main

import (
	"fmt"
	"net/http"
	"os"
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
	// Secure route: /secure. Wrap with AuthMiddleware
	secureRoute := http.HandlerFunc(secureHandler)
	mux.Handle("/secure", AuthMiddleware(secureRoute))
	return mux
}

func main() {
	if err := http.ListenAndServe(":8080", SetupServer()); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
	}
}

// this value is fixed for this assignment
const validToken = "secret"

func validateToken(token string) bool {
	return token == validToken
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		isValid := validateToken(token)
		if !isValid {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}