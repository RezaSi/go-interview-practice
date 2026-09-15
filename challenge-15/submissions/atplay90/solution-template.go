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
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.clients[client.ClientID]
	if exists {
		return errors.New("client already exists")
	}

	s.clients[client.ClientID] = client

	return nil
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	str := base64.RawURLEncoding.EncodeToString(bytes)
	return str[:length], nil
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement authorization endpoint
	// 1. Validate request parameters (client_id, redirect_uri, response_type, scope, state)
	urlValues := r.URL.Query()
	responseType := urlValues.Get("response_type")
	clientID := urlValues.Get("client_id")
	redirectUri := urlValues.Get("redirect_uri")
	scope := urlValues.Get("scope")
	state := urlValues.Get("state")
	codeChallenge := urlValues.Get("code_challenge")
	codeChallengeMethod := urlValues.Get("code_challenge_method")
	
	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()
	if !exists {
		http.Error(w, "client is missing", http.StatusBadRequest)
		return
	}
	rUriFound := false
	for _, rUri := range client.RedirectURIs {
		if rUri == redirectUri {
			rUriFound = true
			break
		}
	}
	if !rUriFound {
		http.Error(w, "redirect URI not found", http.StatusBadRequest)
		return
	}
	if responseType != "code" {
		errURL := redirectUri + "?error=unsupported_response_type"
		if state != "" {
			errURL += "&state=" + url.QueryEscape(state)
		}
		http.Redirect(w, r, errURL, http.StatusFound)
		return
	}
	// 2. Authenticate the user (for this challenge, could be a simple login form)
	userID, ok := (r.Context().Value("user_id")).(string)
	if !ok {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	// 3. Present a consent screen to the user
	// 4. Generate an authorization code and redirect to the client with the code
	code, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, "failed to generate authorization code", http.StatusInternalServerError)
		return
	}
	authCode := &AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectUri,
		Scopes:              strings.Split(scope, " "),
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	s.mu.Lock()
	s.authCodes[code] = authCode
	s.mu.Unlock()

	q := url.Values{}
	q.Set("code", code)
	if state != "" {
		q.Set("state", state)
	}
	http.Redirect(w, r, redirectUri+"?"+q.Encode(), http.StatusFound)
}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "failed to parse form")
		return
	}

	grantType := r.FormValue("grant_type")

	switch grantType {
	case "authorization_code":
		s.handleAuthorizationCodeGrant(w, r)
	case "refresh_token":
		s.handleRefreshTokenGrant(w, r)
	default:
		writeOAuthError(w, http.StatusBadRequest, "unsupported_grant_type", "grant_type must be authorization_code or refresh_token")
	}
}

func (s *OAuth2Server) handleAuthorizationCodeGrant(w http.ResponseWriter, r *http.Request) {
	codeStr := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")
	codeVerifier := r.FormValue("code_verifier")

	s.mu.RLock()
	authCode, exists := s.authCodes[codeStr]
	client, clientExists := s.clients[clientID]
	s.mu.RUnlock()

	if !clientExists || client.ClientSecret != clientSecret {
		writeOAuthError(w, http.StatusUnauthorized, "invalid_client", "client authentication failed")
		return
	}

	if !exists || time.Now().After(authCode.ExpiresAt) || authCode.ClientID != clientID {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "authorization code is invalid or expired")
		return
	}

	if authCode.RedirectURI != redirectURI {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "redirect_uri does not match")
		return
	}

	if authCode.CodeChallenge != "" {
		if codeVerifier == "" || !VerifyCodeChallenge(codeVerifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
			writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "code_verifier is invalid")
			return
		}
	}

	s.mu.Lock()
	delete(s.authCodes, codeStr)
	s.mu.Unlock()

	accessToken, refreshToken, err := s.issueTokens(clientID, authCode.UserID, authCode.Scopes)
	if err != nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to issue tokens")
		return
	}

	writeTokenResponse(w, accessToken, refreshToken)
}

func (s *OAuth2Server) handleRefreshTokenGrant(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr := r.FormValue("refresh_token")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	s.mu.RLock()
	client, clientExists := s.clients[clientID]
	existingRefreshToken, exists := s.refreshTokens[refreshTokenStr]
	s.mu.RUnlock()

	if !clientExists || client.ClientSecret != clientSecret {
		writeOAuthError(w, http.StatusUnauthorized, "invalid_client", "client authentication failed")
		return
	}

	if !exists || time.Now().After(existingRefreshToken.ExpiresAt) || existingRefreshToken.ClientID != clientID {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "refresh token is invalid or expired")
		return
	}

	s.mu.Lock()
	delete(s.refreshTokens, refreshTokenStr)
	s.mu.Unlock()

	accessToken, newRefreshToken, err := s.issueTokens(clientID, existingRefreshToken.UserID, existingRefreshToken.Scopes)
	if err != nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to issue tokens")
		return
	}

	writeTokenResponse(w, accessToken, newRefreshToken)
}

// issueTokens generates and stores a new access token and refresh token pair
func (s *OAuth2Server) issueTokens(clientID, userID string, scopes []string) (*Token, *RefreshToken, error) {
	accessTokenStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}
	refreshTokenStr, err := GenerateRandomString(64)
	if err != nil {
		return nil, nil, err
	}

	accessToken := &Token{
		AccessToken: accessTokenStr,
		ClientID:    clientID,
		UserID:      userID,
		Scopes:      scopes,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	refreshToken := &RefreshToken{
		RefreshToken: refreshTokenStr,
		ClientID:     clientID,
		UserID:       userID,
		Scopes:       scopes,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.tokens[accessTokenStr] = accessToken
	s.refreshTokens[refreshTokenStr] = refreshToken
	s.mu.Unlock()

	return accessToken, refreshToken, nil
}

// writeOAuthError writes a standard OAuth2 JSON error response
func writeOAuthError(w http.ResponseWriter, status int, errCode, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             errCode,
		"error_description": description,
	})
}

// writeTokenResponse writes a successful OAuth2 token JSON response
func writeTokenResponse(w http.ResponseWriter, token *Token, refreshToken *RefreshToken) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  token.AccessToken,
		"token_type":    "Bearer",
		"expires_in":    int(time.Until(token.ExpiresAt).Seconds()),
		"refresh_token": refreshToken.RefreshToken,
		"scope":         strings.Join(token.Scopes, " "),
	})
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	t, exists := s.tokens[token]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.New("token not found")
	}
	if time.Now().After(t.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return t, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	s.mu.RLock()
	rt, exists := s.refreshTokens[refreshToken]
	s.mu.RUnlock()

	if !exists {
		return nil, nil, errors.New("refresh token not found")
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, nil, errors.New("refresh token expired")
	}

	s.mu.Lock()
	delete(s.refreshTokens, refreshToken)
	s.mu.Unlock()

	return s.issueTokens(rt.ClientID, rt.UserID, rt.Scopes)
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if isRefreshToken {
		if _, exists := s.refreshTokens[token]; !exists {
			return errors.New("refresh token not found")
		}
		delete(s.refreshTokens, token)
		return nil
	}

	if _, exists := s.tokens[token]; !exists {
		return errors.New("token not found")
	}
	delete(s.tokens, token)
	return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	switch method {
	case "S256":
		hash := sha256.Sum256([]byte(codeVerifier))
		toCompare := base64.RawURLEncoding.EncodeToString(hash[:])
		return toCompare == codeChallenge
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
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", c.Config.ClientID)
	q.Set("redirect_uri", c.Config.RedirectURI)
	q.Set("scope", strings.Join(c.Config.Scopes, " "))
	q.Set("state", state)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", codeChallengeMethod)

	return c.Config.AuthorizationEndpoint + "?" + q.Encode(), nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", c.Config.RedirectURI)
	form.Set("client_id", c.Config.ClientID)
	form.Set("client_secret", c.Config.ClientSecret)
	form.Set("code_verifier", codeVerifier)

	return c.doTokenRequest(form)
}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", c.RefreshToken)
	form.Set("client_id", c.Config.ClientID)
	form.Set("client_secret", c.Config.ClientSecret)

	return c.doTokenRequest(form)
}

// doTokenRequest posts a token request to the token endpoint and stores the result on the client
func (c *OAuth2Client) doTokenRequest(form url.Values) error {
	resp, err := http.PostForm(c.Config.TokenEndpoint, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	c.AccessToken = tokenResponse.AccessToken
	c.RefreshToken = tokenResponse.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	return nil
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(url string, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

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
