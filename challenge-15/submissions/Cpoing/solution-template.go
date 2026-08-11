package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
)

// OAuth2Config contains configuration for the OAuth2 server
type OAuth2Config struct {
	// AuthorizationEndpoint is the endpoint for authorization requests
	AuthorizationEndpoint string
	// TokenEndpoint is the endpoint for token requests
	TokenEndpoint string
	// ClientID is the OAuth2 client identifier
	ClientID string
	// ClientSecret is the secret for the client
	ClientSecret string
	// RedirectURI is the URI to redirect to after authorization
	RedirectURI string
	// Scopes is a list of requested scopes
	Scopes []string
}

// OAuth2Server implements an OAuth2 authorization server
type OAuth2Server struct {
	// clients stores registered OAuth2 clients
	clients map[string]*OAuth2ClientInfo
	// authCodes stores issued authorization codes
	authCodes map[string]*AuthorizationCode
	// tokens stores issued access tokens
	tokens map[string]*Token
	// refreshTokens stores issued refresh tokens
	refreshTokens map[string]*RefreshToken
	// users stores user credentials for demonstration purposes
	users map[string]*User
	// mutex for concurrent access to data
	mu sync.RWMutex
}

// OAuth2ClientInfo represents a registered OAuth2 client
type OAuth2ClientInfo struct {
	// ClientID is the unique identifier for the client
	ClientID string
	// ClientSecret is the secret for the client
	ClientSecret string
	// RedirectURIs is a list of allowed redirect URIs
	RedirectURIs []string
	// AllowedScopes is a list of scopes the client can request
	AllowedScopes []string
}

// User represents a user in the system
type User struct {
	// ID is the unique identifier for the user
	ID string
	// Username is the username for the user
	Username string
	// Password is the password for the user (in a real system, this would be hashed)
	Password string
}

// AuthorizationCode represents an issued authorization code
type AuthorizationCode struct {
	// Code is the authorization code string
	Code string
	// ClientID is the client that requested the code
	ClientID string
	// UserID is the user that authorized the client
	UserID string
	// RedirectURI is the URI to redirect to
	RedirectURI string
	// Scopes is a list of authorized scopes
	Scopes []string
	// ExpiresAt is when the code expires
	ExpiresAt time.Time
	// CodeChallenge is for PKCE
	CodeChallenge string
	// CodeChallengeMethod is for PKCE
	CodeChallengeMethod string
}

// Token represents an issued access token
type Token struct {
	// AccessToken is the token string
	AccessToken string
	// ClientID is the client that owns the token
	ClientID string
	// UserID is the user that authorized the token
	UserID string
	// Scopes is a list of authorized scopes
	Scopes []string
	// ExpiresAt is when the token expires
	ExpiresAt time.Time
}

// RefreshToken represents an issued refresh token
type RefreshToken struct {
	// RefreshToken is the token string
	RefreshToken string
	// ClientID is the client that owns the token
	ClientID string
	// UserID is the user that authorized the token
	UserID string
	// Scopes is a list of authorized scopes
	Scopes []string
	// ExpiresAt is when the token expires
	ExpiresAt time.Time
}

// NewOAuth2Server creates a new OAuth2Server
func NewOAuth2Server() *OAuth2Server {
	server := &OAuth2Server{
		clients:       make(map[string]*OAuth2ClientInfo),
		authCodes:     make(map[string]*AuthorizationCode),
		tokens:        make(map[string]*Token),
		refreshTokens: make(map[string]*RefreshToken),
		users:         make(map[string]*User),
	}

	// Pre-register some users
	server.users["user1"] = &User{
		ID:       "user1",
		Username: "testuser",
		Password: "password",
	}

	return server
}

// RegisterClient registers a new OAuth2 client
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	// TODO: Implement client registration
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[client.ClientID]; ok {
		return fmt.Errorf("client already exists")
	}

	s.clients[client.ClientID] = client

	return nil
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	// TODO: Implement secure random string generation

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement authorization endpoint
	// 1. Validate request parameters (client_id, redirect_uri, response_type, scope, state)
	// 2. Authenticate the user (for this challenge, could be a simple login form)
	// 3. Present a consent screen to the user
	// 4. Generate an authorization code and redirect to the client with the code
	s.mu.Lock()
	defer s.mu.Unlock()

	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")
	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")

	client, ok := s.clients[clientID]
	if !ok {
		http.Error(w, "client id does not exist", http.StatusBadRequest)
		return

	}

	if !slices.Contains(client.RedirectURIs, redirectURI) {
		http.Error(w, "invalid_redirect_uri", http.StatusBadRequest)
		return
	}

	if responseType != "code" {
		http.Redirect(w, r, fmt.Sprintf("%s?error=unsupported_response_type&state=%s",
			redirectURI, url.QueryEscape(state)), http.StatusFound)
		return
	}

	reqScopes := strings.Fields(scope)
	for _, rs := range reqScopes {
		if !slices.Contains(client.AllowedScopes, rs) {
			http.Error(w, "scope not allowed", http.StatusBadRequest)
			return
		}
	}

	code, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	s.authCodes[code] = &AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              clientID,
		RedirectURI:         redirectURI,
		Scopes:              reqScopes,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	v := url.Values{}
	v.Set("code", code)
	if state != "" {
		v.Set("state", state)
	}
	http.Redirect(w, r, redirectURI+"?"+v.Encode(), http.StatusFound)
}

func writeJSONError(w http.ResponseWriter, errCode, desc string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":             errCode,
		"error_description": desc,
	})
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Refreshtoken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement token endpoint
	// 1. Validate request parameters (grant_type, code, redirect_uri, client_id, client_secret)
	// 2. Verify the authorization code
	// 3. For PKCE, verify the code_verifier
	// 4. Generate access and refresh tokens
	// 5. Return the tokens as a JSON response

	if err := r.ParseForm(); err != nil {
		writeJSONError(w, "invalid request", "error parsing form", http.StatusBadRequest)
		return
	}

	grantType := r.Form.Get("grant_type")
	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	client, ok := s.clients[clientID]
	if !ok || client.ClientSecret != clientSecret {
		writeJSONError(w, "invalid_client", "invalid client", http.StatusUnauthorized)
		return
	}

	switch grantType {
	case "refresh_token":
		refToken := r.Form.Get("refresh_token")

		accessToken, refreshToken, err := s.RefreshAccessToken(refToken)
		if err != nil {
			writeJSONError(w, "server error", "internal server error", http.StatusInternalServerError)
			return
		}

		tokenType := "Bearer"
		expiresIn := 3600
		scope := strings.Join(refreshToken.Scopes, " ")

		resp := &TokenResponse{
			AccessToken:  accessToken.AccessToken,
			TokenType:    tokenType,
			ExpiresIn:    expiresIn,
			Refreshtoken: refreshToken.RefreshToken,
			Scope:        scope,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case "authorization_code":
		code := r.Form.Get("code")
		redirectURI := r.Form.Get("redirect_uri")
		codeVerifier := r.Form.Get("code_verifier")

		authReq, ok := s.authCodes[code]
		if !ok || authReq.ExpiresAt.Before(time.Now()) || authReq.RedirectURI != redirectURI {
			writeJSONError(w, "invalid_grant", "invalid grant type", http.StatusBadRequest)
			return
		}

		if authReq.CodeChallenge != "" {
			if !VerifyCodeChallenge(codeVerifier, authReq.CodeChallenge, authReq.CodeChallengeMethod) {
				writeJSONError(w, "invalid_grant", "invalid grant type", http.StatusBadRequest)
				return
			}
		}

		delete(s.authCodes, code)

		accessToken, err := GenerateRandomString(32)
		if err != nil {
			writeJSONError(w, "server error", "internal server error", http.StatusInternalServerError)
			return
		}
		refreshToken, err := GenerateRandomString(32)
		if err != nil {
			writeJSONError(w, "server error", "internal server error", http.StatusInternalServerError)
			return
		}

		token := &Token{
			AccessToken: accessToken,
			ClientID:    clientID,
			UserID:      authReq.UserID,
			Scopes:      authReq.Scopes,
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		}
		rt := &RefreshToken{
			RefreshToken: refreshToken,
			ClientID:     clientID,
			UserID:       authReq.UserID,
			Scopes:       authReq.Scopes,
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}

		s.tokens[accessToken] = token
		s.refreshTokens[refreshToken] = rt

		tokenType := "Bearer"
		expiresIn := 3600
		scope := strings.Join(authReq.Scopes, " ")

		resp := &TokenResponse{
			AccessToken:  accessToken,
			TokenType:    tokenType,
			ExpiresIn:    expiresIn,
			Refreshtoken: refreshToken,
			Scope:        scope,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		writeJSONError(w, "invalid_grant", "invalid grant type", http.StatusBadRequest)
	}
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	// TODO: Implement token validation
	s.mu.Lock()
	defer s.mu.Unlock()

	tokenInfo, exists := s.tokens[token]
	if !exists {
		return nil, errors.New("token not found")
	}

	if tokenInfo.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return tokenInfo, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	// TODO: Implement token refresh
	s.mu.Lock()
	defer s.mu.Unlock()

	rt, exists := s.refreshTokens[refreshToken]
	if !exists || rt.ExpiresAt.Before(time.Now()) {
		return nil, nil, errors.New("token expired")
	}

	rToken, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}

	rRefreshToken, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}

	t := &Token{
		AccessToken: rToken,
		ClientID:    rt.ClientID,
		UserID:      rt.UserID,
		Scopes:      rt.Scopes,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	r := &RefreshToken{
		RefreshToken: rRefreshToken,
		ClientID:     rt.ClientID,
		UserID:       rt.UserID,
		Scopes:       rt.Scopes,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	delete(s.refreshTokens, refreshToken)

	s.tokens[rToken] = t
	s.refreshTokens[rRefreshToken] = r

	return t, r, nil
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	// TODO: Implement token revocation
	s.mu.Lock()
	defer s.mu.Unlock()

	if isRefreshToken {
		_, exists := s.refreshTokens[token]
		if !exists {
			return errors.New("token does not exist")
		}
		delete(s.refreshTokens, token)
	} else {
		_, exists := s.tokens[token]
		if !exists {
			return errors.New("token does not exist")
		}
		delete(s.tokens, token)
	}

	return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	// TODO: Implement PKCE verification
	switch method {
	case "S256":
		hashedVerifier := sha256.Sum256([]byte(codeVerifier))
		expectedChallenge := base64.RawURLEncoding.EncodeToString(hashedVerifier[:])
		return expectedChallenge == codeChallenge
	case "plain":
		return codeVerifier == codeChallenge
	default:
		return false
	}
}

// StartServer starts the OAuth2 server
func (s *OAuth2Server) StartServer(port int) error {
	// Register HTTP handlers
	http.HandleFunc("/authorize", s.HandleAuthorize)
	http.HandleFunc("/token", s.HandleToken)

	// Start the server
	fmt.Printf("Starting OAuth2 server on port %d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// Client code to demonstrate usage

// OAuth2Client represents a client application using OAuth2
type OAuth2Client struct {
	// Config is the OAuth2 configuration
	Config OAuth2Config
	// Token is the current access token
	AccessToken string
	// RefreshToken is the current refresh token
	RefreshToken string
	// TokenExpiry is when the access token expires
	TokenExpiry time.Time
}

// NewOAuth2Client creates a new OAuth2 client
func NewOAuth2Client(config OAuth2Config) *OAuth2Client {
	return &OAuth2Client{Config: config}
}

// GetAuthorizationURL returns the URL to redirect the user for authorization
func (c *OAuth2Client) GetAuthorizationURL(state string, codeChallenge string, codeChallengeMethod string) (string, error) {
	// TODO: Implement building the authorization URL
	authURL, err := url.Parse(c.Config.AuthorizationEndpoint)
	if err != nil {
		return "", err
	}

	q := authURL.Query()
	q.Set("client_id", c.Config.ClientID)
	q.Set("redirect_uri", c.Config.RedirectURI)
	q.Set("response_type", "code")
	q.Set("scope", strings.Join(c.Config.Scopes, " "))
	q.Set("state", state)

	if codeChallenge != "" {
		q.Set("code_challenge", codeChallenge)
		q.Set("code_challenge_method", codeChallengeMethod)
	}

	authURL.RawQuery = q.Encode()
	return authURL.String(), nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	// TODO: Implement token exchange
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("redirect_uri", c.Config.RedirectURI)
	if codeVerifier != "" {
		v.Set("code_verifier", codeVerifier)
	}
	v.Set("client_id", c.Config.ClientID)
	v.Set("client_secret", c.Config.ClientSecret)

	req, err := http.NewRequest("POST", c.Config.TokenEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&e)
		return errors.New("token exchange failed")
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return err
	}

	c.AccessToken = tr.AccessToken
	c.RefreshToken = tr.Refreshtoken
	if tr.ExpiresIn > 0 {
		c.TokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	} else {
		c.TokenExpiry = time.Time{}
	}

	return nil
}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	// TODO: Implement token refresh
	if c.RefreshToken == "" {
		return errors.New("no refresh token")
	}

	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", c.RefreshToken)
	v.Set("client_id", c.Config.ClientID)
	v.Set("client_secret", c.Config.ClientSecret)

	req, err := http.NewRequest("POST", c.Config.TokenEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&e)
		return errors.New("refresh token request failed")
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return err
	}

	c.AccessToken = tr.AccessToken
	if tr.Refreshtoken != "" {
		c.RefreshToken = tr.Refreshtoken
	}
	if tr.ExpiresIn > 0 {
		c.TokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	} else {
		c.TokenExpiry = time.Time{}
	}

	return nil
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(url string, method string) (*http.Response, error) {
	// TODO: Implement authenticated request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}
	return http.DefaultClient.Do(req)
}

func main() {
	// Example of starting the OAuth2 server
	server := NewOAuth2Server()

	// Register a client
	client := &OAuth2ClientInfo{
		ClientID:      "example-client",
		ClientSecret:  "example-secret",
		RedirectURIs:  []string{"http://localhost:8080/callback"},
		AllowedScopes: []string{"read", "write"},
	}
	server.RegisterClient(client)

	// Start the server in a goroutine
	go func() {
		err := server.StartServer(9000)
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	fmt.Println("OAuth2 server is running on port 9000")

	// Example of using the client (this wouldn't actually work in main, just for demonstration)
	/*
		client := NewOAuth2Client(OAuth2Config{
			AuthorizationEndpoint: "http://localhost:9000/authorize",
			TokenEndpoint:         "http://localhost:9000/token",
			ClientID:              "example-client",
			ClientSecret:          "example-secret",
			RedirectURI:           "http://localhost:8080/callback",
			Scopes:                []string{"read", "write"},
		})

		// Generate a code verifier and challenge for PKCE
		codeVerifier, _ := GenerateRandomString(64)
		codeChallenge := GenerateCodeChallenge(codeVerifier, "S256")

		// Get the authorization URL and redirect the user
		authURL, _ := client.GetAuthorizationURL("random-state", codeChallenge, "S256")
		fmt.Printf("Please visit: %s\n", authURL)

		// After authorization, exchange the code for tokens
		client.ExchangeCodeForToken("returned-code", codeVerifier)

		// Make an authenticated request
		resp, _ := client.MakeAuthenticatedRequest("http://api.example.com/resource", "GET")
		fmt.Printf("Response: %v\n", resp)
	*/
}
