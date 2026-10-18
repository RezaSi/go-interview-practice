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

const (
	codeTimeout         = 5 * time.Minute
	maxBodySize         = 1 * 1024 * 1024
	randomStringSize    = 32
	accessTokenTimeout  = 24 * time.Hour
	refreshTokenTimeout = 7 * 24 * time.Hour
)

var (
	ErrClientNotFound       = errors.New("client not found")
	ErrExpiredCode          = errors.New("expired code")
	ErrIDNotUnique          = errors.New("provided ClientID already exists")
	ErrInvalidGrant         = errors.New("invalid grant")
	ErrInvalidPKCEMethod    = errors.New("invalid challenge method")
	ErrInvalidRedirectURI   = errors.New("invalid redirect URI")
	ErrNoAuthCode           = errors.New("no auth code found")
	ErrATGen                = errors.New("access token generation error")
	ErrATNotFound           = errors.New("unregistered access token")
	ErrRTGen                = errors.New("refresh token generation error")
	ErrRTNotFound           = errors.New("unregistered refresh token")
	ErrSecretMismatch       = errors.New("secret mismatch")
	ErrTokenTimeout         = errors.New("token timeout")
	ErrUnprocessablePayload = errors.New("unprocessable payload")
	ErrUnsupportedGrant     = errors.New("grant type error")
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

type TokenResponce struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// RegisterClient registers a new OAuth2 client
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[client.ClientID]; exists {
		return ErrIDNotUnique
	}

	s.clients[client.ClientID] = client

	return nil
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = alphabet[b%byte(len(alphabet))]
	}

	return string(bytes), nil
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	clientID := q.Get("client_id")
	codeChallenge := q.Get("code_challenge")
	codeChallengeMethod := q.Get("code_challenge_method")
	redirectURI := q.Get("redirect_uri")
	responseType := q.Get("response_type")
	scope := q.Get("scope")
	state := q.Get("state")

	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, ErrClientNotFound.Error(), http.StatusBadRequest)
		return
	}

	if !slices.Contains(client.RedirectURIs, redirectURI) {
		http.Error(w, ErrInvalidRedirectURI.Error(), http.StatusBadRequest)
		return
	}

	u, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "redirectURI parsing error", http.StatusInternalServerError)
		return
	}

	resQ := u.Query()

	if responseType != "code" {
		resQ.Set("error", "unsupported_response_type")
		u.RawQuery = resQ.Encode()

		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}

	if codeChallengeMethod == "" {
		codeChallengeMethod = "plain"
	}

	if codeChallenge != "" && codeChallengeMethod != "S256" && codeChallengeMethod != "plain" {
		http.Error(w, ErrInvalidPKCEMethod.Error(), http.StatusBadRequest)
		return
	}

	randomString, err := GenerateRandomString(randomStringSize)
	if err != nil {
		http.Error(w, "string generation error", http.StatusInternalServerError)
		return
	}

	allowedMap := make(map[string]struct{}, len(client.AllowedScopes))
	for _, v := range client.AllowedScopes {
		allowedMap[v] = struct{}{}
	}

	checkedScopes := make([]string, 0, len(client.AllowedScopes))
	for _, s := range strings.Split(scope, " ") {
		if s == "" {
			continue
		}
		if _, exists := allowedMap[s]; exists {
			checkedScopes = append(checkedScopes, s)
		}
	}

	authCode := AuthorizationCode{
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		Scopes:              checkedScopes,
		Code:                randomString,
		ExpiresAt:           time.Now().Add(codeTimeout),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	s.mu.Lock()
	s.authCodes[randomString] = &authCode
	s.mu.Unlock()

	resQ.Set("code", randomString)
	resQ.Set("state", state)

	u.RawQuery = resQ.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	if err := r.ParseForm(); err != nil {
		http.Error(w, ErrUnprocessablePayload.Error(), http.StatusUnprocessableEntity)
		return
	}

	grantType := r.PostForm.Get("grant_type")
	if grantType != "authorization_code" && grantType != "refresh_token" {
		http.Error(w, ErrUnsupportedGrant.Error(), http.StatusBadRequest)
		return
	}

	clientID := r.PostForm.Get("client_id")

	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, ErrClientNotFound.Error(), http.StatusNotFound)
		return
	}

	secret := r.PostForm.Get("client_secret")
	if client.ClientSecret != secret {
		sendJSONError(w, http.StatusUnauthorized, "invalid_client", ErrSecretMismatch.Error())
		return
	}

	switch grantType {
	case "authorization_code":
		code := r.PostForm.Get("code")

		s.mu.Lock()
		authCode, exists := s.authCodes[code]
		if !exists {
			s.mu.Unlock()
			sendJSONError(w, http.StatusUnauthorized, "invalid_grant", ErrNoAuthCode.Error())
			return
		}

		if authCode.ClientID != clientID {
			s.mu.Unlock()
			sendJSONError(w, http.StatusBadRequest, "invalid_grant", "client mismatch")
			return
		}

		delete(s.authCodes, code)
		s.mu.Unlock()

		redirectURI := r.PostForm.Get("redirect_uri")
		if authCode.RedirectURI != redirectURI {
			http.Error(w, ErrInvalidRedirectURI.Error(), http.StatusBadRequest)
			return
		}

		if time.Now().After(authCode.ExpiresAt) {
			http.Error(w, ErrExpiredCode.Error(), http.StatusUnauthorized)
			return
		}

		if authCode.CodeChallenge != "" {
			verifier := r.PostForm.Get("code_verifier")
			if !VerifyCodeChallenge(verifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
				sendJSONError(w, http.StatusBadRequest, "invalid_grant", ErrInvalidGrant.Error())
				return
			}
		}

		at, err := GenerateRandomString(randomStringSize)
		if err != nil {
			http.Error(w, "at token generation error", http.StatusInternalServerError)
			return
		}

		rt, err := GenerateRandomString(randomStringSize)
		if err != nil {
			http.Error(w, "rt token generation error", http.StatusInternalServerError)
			return
		}

		accToken := &Token{
			ClientID:    clientID,
			AccessToken: at,
			Scopes:      authCode.Scopes,
			ExpiresAt:   time.Now().Add(accessTokenTimeout),
		}

		refToken := &RefreshToken{
			ClientID:     clientID,
			RefreshToken: rt,
			Scopes:       authCode.Scopes,
			ExpiresAt:    time.Now().Add(refreshTokenTimeout),
		}

		s.mu.Lock()
		s.tokens[at] = accToken
		s.refreshTokens[rt] = refToken
		s.mu.Unlock()

		payload := &TokenResponce{
			AccessToken:  at,
			TokenType:    "Bearer",
			ExpiresIn:    int64(time.Until(accToken.ExpiresAt).Seconds()),
			RefreshToken: rt,
			Scope:        strings.Join(accToken.Scopes, " "),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	case "refresh_token":
		refreshToken := r.PostForm.Get("refresh_token")
		newAt, newRt, err := s.RefreshAccessToken(refreshToken, clientID)

		if err != nil {
			if errors.Is(err, ErrInvalidGrant) {
				sendJSONError(w, http.StatusBadRequest, "invalid_grant", ErrInvalidGrant.Error())
				return
			}

			if errors.Is(err, ErrRTNotFound) ||
				errors.Is(err, ErrClientNotFound) ||
				errors.Is(err, ErrTokenTimeout) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		payload := &TokenResponce{
			AccessToken:  newAt.AccessToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(time.Until(newAt.ExpiresAt).Seconds()),
			RefreshToken: newRt.RefreshToken,
			Scope:        strings.Join(newAt.Scopes, " "),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(payload)
		if err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
			return
		}
	}
}

func sendJSONError(w http.ResponseWriter, status int, errorType, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	payload := struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}{
		Error:       errorType,
		Description: desc,
	}

	_ = json.NewEncoder(w).Encode(payload)
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	at, exists := s.tokens[token]
	s.mu.RUnlock()

	if !exists {
		return nil, ErrATNotFound
	}

	if time.Now().After(at.ExpiresAt) {
		return nil, ErrTokenTimeout
	}

	return at, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string, clientID string) (*Token, *RefreshToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	refToken, tokenExists := s.refreshTokens[refreshToken]
	if !tokenExists {
		return nil, nil, ErrRTNotFound
	}

	if refToken.ClientID != clientID {
		return nil, nil, ErrInvalidGrant
	}

	delete(s.refreshTokens, refreshToken)

	if time.Now().After(refToken.ExpiresAt) {
		return nil, nil, ErrTokenTimeout
	}

	client, clientExists := s.clients[refToken.ClientID]
	if !clientExists {
		return nil, nil, ErrClientNotFound
	}

	newAtStr, err := GenerateRandomString(randomStringSize)
	if err != nil {
		return nil, nil, ErrATGen
	}

	newRtStr, err := GenerateRandomString(randomStringSize)
	if err != nil {
		return nil, nil, ErrRTGen
	}

	newAccessToken := &Token{
		ClientID:    client.ClientID,
		AccessToken: newAtStr,
		Scopes:      refToken.Scopes,
		ExpiresAt:   time.Now().Add(accessTokenTimeout),
	}

	newRefreshToken := &RefreshToken{
		ClientID:     client.ClientID,
		RefreshToken: newRtStr,
		Scopes:       refToken.Scopes,
		ExpiresAt:    time.Now().Add(refreshTokenTimeout),
	}

	s.tokens[newAtStr] = newAccessToken
	s.refreshTokens[newRtStr] = newRefreshToken

	return newAccessToken, newRefreshToken, nil
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch isRefreshToken {
	case true:
		if _, exists := s.refreshTokens[token]; !exists {
			return ErrRTNotFound
		}
		delete(s.refreshTokens, token)

	case false:
		if _, exists := s.tokens[token]; !exists {
			return ErrATNotFound
		}
		delete(s.tokens, token)
	}

	return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	switch method {
	case "S256":
		checksum := sha256.Sum256([]byte(codeVerifier))
		strChecksum := base64.RawURLEncoding.EncodeToString(checksum[:])
		if strChecksum != codeChallenge {
			return false
		}
		return true

	case "plain":
		if codeVerifier != codeChallenge {
			return false
		}
		return true

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
	authPath := c.Config.AuthorizationEndpoint
	u, err := url.Parse(authPath)
	if err != nil {
		return "", err
	}

	v := u.Query()
	v.Set("response_type", "code")
	v.Set("client_id", c.Config.ClientID)
	v.Set("redirect_uri", c.Config.RedirectURI)
	v.Set("scope", strings.Join(c.Config.Scopes, " "))
	v.Set("state", state)
	v.Set("code_challenge", codeChallenge)
	v.Set("code_challenge_method", codeChallengeMethod)
	u.RawQuery = v.Encode()

	return u.String(), nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("redirect_uri", c.Config.RedirectURI)
	v.Set("client_id", c.Config.ClientID)
	v.Set("client_secret", c.Config.ClientSecret)
	if codeVerifier != "" {
		v.Set("code_verifier", codeVerifier)
	}

	resp, err := http.PostForm(c.Config.TokenEndpoint, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %v", resp.Status)
	}

	var tokenResp TokenResponce
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return err
	}

	c.AccessToken = tokenResp.AccessToken
	c.RefreshToken = tokenResp.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", c.RefreshToken)
	v.Set("client_id", c.Config.ClientID)
	v.Set("client_secret", c.Config.ClientSecret)

	resp, err := http.PostForm(c.Config.TokenEndpoint, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %v", resp.Status)
	}

	var respToken TokenResponce
	err = json.NewDecoder(resp.Body).Decode(&respToken)
	if err != nil {
		return err
	}

	c.AccessToken = respToken.AccessToken
	c.RefreshToken = respToken.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(respToken.ExpiresIn) * time.Second)

	return nil
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(url string, method string) (*http.Response, error) {
	// TODO: Implement authenticated request
	if time.Now().After(c.TokenExpiry) {
		err := c.DoRefreshToken()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
