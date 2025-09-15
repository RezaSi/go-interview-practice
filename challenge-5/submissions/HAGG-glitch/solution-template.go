package main

import (
	"net/http"
)

// AuthMiddleware checks for "X-Auth-Token"
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Public handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}

// Secure handler
func secureHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are authorized!"))
}

// SetupServer builds the mux (needed for tests)
func SetupServer() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/hello", http.HandlerFunc(helloHandler))
	mux.Handle("/secure", AuthMiddleware(http.HandlerFunc(secureHandler)))
	return mux
}

func main() {
	http.ListenAndServe(":8080", SetupServer())
}
