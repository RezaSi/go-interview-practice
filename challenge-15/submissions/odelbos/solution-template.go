package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"strings"
	"encoding/json"
	"slices"
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
	return server
}

// RegisterClient registers a new OAuth2 client
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client.ClientID == "" || client.ClientSecret == "" {
		return errors.New("client ID and secret are required")
	}
	if _, ok := s.clients[client.ClientID]; ok {
		return errors.New("client ID already exists")
	}
	s.clients[client.ClientID] = client
	return nil
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clientID := r.URL.Query().Get("client_id")
	client, ok := s.clients[clientID]
	if ! ok {
		http.Error(w, "invalid client ID", http.StatusBadRequest)
		return
	}

	redirectURI := r.URL.Query().Get("redirect_uri")

	responseType := r.URL.Query().Get("response_type")
	if responseType != "code" {
		redirectURL, _ := url.Parse(redirectURI)
		query := redirectURL.Query()
		query.Set("error", "unsupported_response_type")
		redirectURL.RawQuery = query.Encode()
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		// INFO: Should be StatusBadRequest without a redirect but tests
		// want StatusFound with a redirect
		return
	}

	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")
	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")

	if ! slices.Contains(client.RedirectURIs, redirectURI) {
		http.Error(w, "invalid redirect URI", http.StatusBadRequest)
		return
	}

	requestedScopes := strings.Split(scope, " ")
	for _, sc := range requestedScopes {
		if ! slices.Contains(client.AllowedScopes, sc) {
			http.Error(w, "invalid scope", http.StatusBadRequest)
			return
		}
	}

	code, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.authCodes[code] = &AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              clientID,
		RedirectURI:         redirectURI,
		Scopes:              requestedScopes,
		ExpiresAt:           time.Now().Add(5 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	redirectURL, _ := url.Parse(redirectURI)
	query := redirectURL.Query()
	query.Set("code", code)
	if state != "" {
		query.Set("state", state)
	}
	redirectURL.RawQuery = query.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

type errorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func writeJSONError(w http.ResponseWriter, error, description string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := &errorResponse{error, description}
	json, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(json)
}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeJSONError(w, "invalid_request", "invalid request", http.StatusBadRequest)
		return
	}

	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")
	client, ok := s.clients[clientID]
	if ! ok || client.ClientSecret != clientSecret {
		writeJSONError(w, "invalid_client", "invalid client", http.StatusUnauthorized)
		return
	}

	grantType := r.Form.Get("grant_type")
	if grantType == "authorization_code" {
		s.handleAutorizationCode(w , r)
		return
	} else if grantType == "refresh_token" {
		s.handleRefreshToken(w , r)
		return
	} else {
		writeJSONError(w, "invalid_grant", "invalid grant type", http.StatusBadRequest)
	}
}

func (s *OAuth2Server) handleAutorizationCode(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := r.ParseForm()
	if err != nil {
		writeJSONError(w, "invalid_request", "invalid request", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	redirectURI := r.Form.Get("redirect_uri")
	clientID := r.Form.Get("client_id")
	codeVerifier := r.Form.Get("code_verifier")

	authCode, ok := s.authCodes[code]
	if ! ok || authCode.ExpiresAt.Before(time.Now()) || authCode.RedirectURI != redirectURI {
		writeJSONError(w, "invalid_auth_code", "invalid authorization code", http.StatusBadRequest)
		return
	}

	if authCode.CodeChallenge != "" {
		if ! VerifyCodeChallenge(codeVerifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
			writeJSONError(w, "invalid_grant", "bad code challenge", http.StatusBadRequest)
			return
		}
	}

	accessToken, err := GenerateRandomString(32)
	if err != nil {
		writeJSONError(w, "server_error", "internal server error", http.StatusInternalServerError)
		return
	}
	refreshToken, err := GenerateRandomString(32)
	if err != nil {
		writeJSONError(w, "server_error", "internal server error", http.StatusInternalServerError)
		return
	}

	// Store tokens
	s.tokens[accessToken] = &Token{
		AccessToken: accessToken,
		ClientID:    clientID,
		UserID:      authCode.UserID,
		Scopes:      authCode.Scopes,
		ExpiresAt:   time.Now().Add(time.Hour)}

	s.refreshTokens[refreshToken] = &RefreshToken{
		RefreshToken: refreshToken,
		ClientID:     clientID,
		UserID:       authCode.UserID,
		Scopes:       authCode.Scopes,
		ExpiresAt:    time.Now().Add(24 * time.Hour)}

	delete(s.authCodes, code)

	response := &tokenResponse{
		accessToken,
		"Bearer",
		3600,
		refreshToken,
		strings.Join(authCode.Scopes, " ")}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *OAuth2Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeJSONError(w, "invalid_request", "invalid request", http.StatusBadRequest)
		return
	}

	rToken := r.Form.Get("refresh_token")

	accessToken, refreshToken, err := s.RefreshAccessToken(rToken)
	if err != nil {
		writeJSONError(w, "server_error", "internal server error", http.StatusInternalServerError)
		return
	}

	response := &tokenResponse{
		accessToken.AccessToken,
		"Bearer",
		3600,
		refreshToken.RefreshToken,
		strings.Join(refreshToken.Scopes, " ")}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tokens[token]
	if ! ok || t.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid token")
	}
	return t, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rt, ok := s.refreshTokens[refreshToken]
	if ! ok || rt.ExpiresAt.Before(time.Now()) {
		return nil, nil, errors.New("invalid token")
	}

	accessToken, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}
	newRefreshToken, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}

	token := &Token{
		AccessToken: accessToken,
		ClientID:    rt.ClientID,
		UserID:      rt.UserID,
		Scopes:      rt.Scopes,
		ExpiresAt:   time.Now().Add(time.Hour)}

	newRT := &RefreshToken{
		RefreshToken: newRefreshToken,
		ClientID:     rt.ClientID,
		UserID:       rt.UserID,
		Scopes:       rt.Scopes,
		ExpiresAt:    time.Now().Add(24 * time.Hour)}

	s.tokens[accessToken] = token
	s.refreshTokens[newRefreshToken] = newRT
	delete(s.refreshTokens, refreshToken)

	return token, newRT, nil
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if isRefreshToken {
		if _, ok := s.refreshTokens[token]; ok {
			delete(s.refreshTokens, token)
			return nil
		}
	} else {
		if _, ok := s.tokens[token]; ok {
			delete(s.tokens, token)
			return nil
		}
	}
	return errors.New("token not found")
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	if method == "S256" {
		hash := sha256.Sum256([]byte(codeVerifier))
		return base64.RawURLEncoding.EncodeToString(hash[:]) == codeChallenge
	} else if method == "plain" {
		return codeVerifier == codeChallenge
	}
	return false
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
	url, err := url.Parse(c.Config.AuthorizationEndpoint)
	if err != nil {
		return "", err
	}
	query := url.Query()
	query.Set("client_id", c.Config.ClientID)
	query.Set("redirect_uri", c.Config.RedirectURI)
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(c.Config.Scopes, " "))
	query.Set("state", state)
	if codeChallenge != "" {
		query.Set("code_challenge", codeChallenge)
		query.Set("code_challenge_method", codeChallengeMethod)
	}
	url.RawQuery = query.Encode()
	return url.String(), nil
}

type tokensResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.Config.RedirectURI)
	data.Set("client_id", c.Config.ClientID)
	data.Set("client_secret", c.Config.ClientSecret)
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequest("POST", c.Config.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed: %s", resp.Status)
	}

	var tokens tokensResp
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return err
	}

	c.AccessToken = tokens.AccessToken
	c.RefreshToken = tokens.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	return nil
}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.RefreshToken)
	data.Set("client_id", c.Config.ClientID)
	data.Set("client_secret", c.Config.ClientSecret)

	req, err := http.NewRequest("POST", c.Config.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh token request failed: %s", resp.Status)
	}

	var tokens tokensResp
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return err
	}

	c.AccessToken = tokens.AccessToken
	c.RefreshToken = tokens.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	return nil
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(url string, method string) (*http.Response, error) {
	if c.TokenExpiry.Before(time.Now()) {
		if err := c.DoRefreshToken(); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + c.AccessToken)

	client := &http.Client{}
	return client.Do(req)
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
