package main

import (
	"fmt"
	"net/http"
)

const validToken = "secret"

// AuthMiddleware checks the "X-Auth-Token" header.
// If it's "secret", call the next handler.
// Otherwise, respond with 401 Unauthorized.
func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, re *http.Request) { // Here noted it down , to satisfy the Handler we use this http.HandlerFunc
		// function. because it implements the ServerHTTP method.
		s := re.Header.Get("X-Auth-Token") // key Authorization value "Bearer 123.456.789"
		//ne := strings.Trim(s, "Bearer")    // to remove the Bearer string
		if s == "" || s != validToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(""))
		return 
		}
		// 		new, err := jwt.Parse(ne, func(k *jwt.Token) (any, error) {
		// 			//if varifietoken == "" {
		// 			//	return nil, http.Error(w, "Empty token", http.StatusInternalServerError)
		// 			//} this is wrong because http.Error writes it on writer.
		// 			return validToken, nil
		// 		})
		// 		if err != nil || !new.Valid { // !new.Valid here it means that Token sahi hai ,Signature match hai, Expire nahi hua
		//             w.WriteHeader(http.StatusUnauthorized)
		// 			http.Error(w, "Error while decoding", http.StatusUnauthorized)
		// 		}
			w.WriteHeader(http.StatusOK)
		h.ServeHTTP(w, re)
	})
} // syntax is very importan please mug up this.


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
