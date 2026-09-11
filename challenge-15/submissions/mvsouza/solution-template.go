package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
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

// Token Server reponse
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
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

// OAuth2 error payload
type OAuth2Error struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
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
	if s == nil || s.clients == nil {
		return errors.New("Server not properly instatiated")
	}
	if client == nil {
		return errors.New("Client cannot be nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[client.ClientID]; !ok {
		s.clients[client.ClientID] = client.Clone()
		return nil
	} else {
		return errors.New("Client id already in use")
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

func (c *OAuth2ClientInfo) HasCallbackUrl(urlParam *url.URL) bool {
	for _, urlString := range c.RedirectURIs {
		urlItem, err := url.Parse(urlString)
		if err != nil {
			continue
		}
		if urlItem.Scheme == urlParam.Scheme && urlItem.Host == urlParam.Host && urlItem.Path == urlParam.Path && urlItem.RawQuery == urlParam.RawQuery {
			return true
		}
	}
	return false
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	clientId := query.Get("client_id")

	s.mu.RLock()
	client, ok := s.clients[clientId]
	if !ok {
		s.mu.RUnlock()
		http.Error(w, fmt.Sprintf("Invalid client_id: %q", clientId), http.StatusBadRequest)
		return
	}
	client = client.Clone()
	s.mu.RUnlock()
	redirectUri := query.Get("redirect_uri")
	urlParsed, err := url.Parse(redirectUri)
	if err != nil {
		http.Error(w, "Missing redirect_uri", http.StatusBadRequest)
		return
	}
	if !client.HasCallbackUrl(urlParsed) {
		http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
		return
	}
	responseType := query.Get("response_type")
	newQuery := urlParsed.Query()
	if responseType != "code" {
		newQuery.Add("error", "unsupported_response_type")
		urlParsed.RawQuery = newQuery.Encode()
		w.Header().Add("Location", urlParsed.String())
		http.Redirect(w, r, urlParsed.String(), http.StatusFound)
		return
	}

	state := query.Get("state")
	if state != "" {
		newQuery.Set("state", state)
	}
	codeChallenge := query.Get("code_challenge")
	codeChallengeMethod := query.Get("code_challenge_method")
	scope := query.Get("scope")
	requestedScopes := strings.Fields(scope)
	allowed := make(map[string]struct{}, len(client.AllowedScopes))
	for _, s := range client.AllowedScopes {
		allowed[s] = struct{}{}
	}
	for _, s := range requestedScopes {
		if _, ok := allowed[s]; !ok {
			newQuery.Set("error", "invalid_scope")
			urlParsed.RawQuery = newQuery.Encode()
			http.Redirect(w, r, urlParsed.String(), http.StatusFound)
			return
		}
	}
	if len(requestedScopes) > 0 {
		newQuery.Set("scope", strings.Join(requestedScopes, " "))
	}
	code, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate authorization code", http.StatusInternalServerError)
		return
	}
	newQuery.Set("code", code)

	s.mu.Lock()
	s.authCodes[code] = &AuthorizationCode{
		Code:                code,
		ClientID:            clientId,
		RedirectURI:         redirectUri,
		Scopes:              requestedScopes,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	s.mu.Unlock()

	urlParsed.RawQuery = newQuery.Encode()
	http.Redirect(w, r, urlParsed.String(), http.StatusFound)
}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}

	client := s.GetClient(r)
	if client == nil {
		errorPayload := OAuth2Error{
			Error:       "invalid_client",
			Description: "Client details dont match",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorPayload)
		return
	}
	grantType := r.FormValue("grant_type")
	var token *TokenResponse
	switch grantType {
	case "authorization_code":
		code, ok := s.ConsumeVerifiedCode(r, client)
		if !ok {
			errorPayload := OAuth2Error{
				Error:       "invalid_grant",
				Description: "Incorect code",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorPayload)
			return
		}
		token = s.GenareteTokenFromAuthCode(code)
	case "refresh_token":
		refreshTokenString := r.FormValue("refresh_token")
		token = s.RefreshToken(client.ClientID, refreshTokenString)
		if token == nil {
			errorPayload := OAuth2Error{
				Error:       "invalid_grant",
				Description: "Invalid refresh token",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorPayload)
			return
		}
	default:
		errorPayload := OAuth2Error{
			Error:       "unsupported_grant_type",
			Description: "Unsupported grant_type",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorPayload)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK)
}

func (s *OAuth2Server) GetClient(r *http.Request) *OAuth2ClientInfo {
	clientId := r.FormValue("client_id")
	secret := r.FormValue("client_secret")
	return s.GetClientById(clientId, secret)
}

func (s *OAuth2Server) GetClientById(clientId, secret string) *OAuth2ClientInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	client, ok := s.clients[clientId]
	if !ok {
		return nil
	}
	if secret != client.ClientSecret {
		return nil
	}
	return client.Clone()
}

func (c *OAuth2ClientInfo) Clone() *OAuth2ClientInfo {
	return &OAuth2ClientInfo{
		ClientID:      c.ClientID,
		ClientSecret:  c.ClientSecret,
		RedirectURIs:  append([]string(nil), c.RedirectURIs...),
		AllowedScopes: append([]string(nil), c.AllowedScopes...),
	}
}

func (s *OAuth2Server) RefreshToken(clientid, refreshTokenString string) *TokenResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	refreshToken, ok := s.refreshTokens[refreshTokenString]
	if !ok {
		return nil
	}
	if clientid != refreshToken.ClientID || refreshToken.Expired() {
		return nil
	}
	accessTokem, err := GenerateRandomString(256)
	if err != nil {
		return nil
	}
	newRefreshToken, err := GenerateRandomString(256)
	if err != nil {
		return nil
	}
	iat := time.Now().Add(10 * time.Minute)
	s.tokens[accessTokem] = &Token{
		AccessToken: accessTokem,
		ClientID:    refreshToken.ClientID,
		UserID:      refreshToken.UserID,
		Scopes:      refreshToken.Scopes,
		ExpiresAt:   iat,
	}
	s.refreshTokens[newRefreshToken] = &RefreshToken{
		RefreshToken: newRefreshToken,
		ClientID:     refreshToken.ClientID,
		UserID:       refreshToken.UserID,
		Scopes:       refreshToken.Scopes,
		ExpiresAt:    iat,
	}
	delete(s.refreshTokens, refreshToken.RefreshToken)
	return &TokenResponse{
		AccessToken:  accessTokem,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(iat).Seconds()),
		Scope:        strings.Join(refreshToken.Scopes, " "),
	}
}

func (s *OAuth2Server) GenareteTokenFromAuthCode(authCode *AuthorizationCode) *TokenResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	accessTokem, _ := GenerateRandomString(256)
	refreshToken, _ := GenerateRandomString(256)
	iat := time.Now().Add(10 * time.Minute)
	s.tokens[accessTokem] = &Token{
		AccessToken: accessTokem,
		ClientID:    authCode.ClientID,
		UserID:      authCode.UserID,
		Scopes:      authCode.Scopes,
		ExpiresAt:   iat,
	}
	s.refreshTokens[refreshToken] = &RefreshToken{
		RefreshToken: refreshToken,
		ClientID:     authCode.ClientID,
		UserID:       authCode.UserID,
		Scopes:       authCode.Scopes,
		ExpiresAt:    iat,
	}
	delete(s.authCodes, authCode.Code)
	return &TokenResponse{
		AccessToken:  accessTokem,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(iat).Seconds()),
		Scope:        strings.Join(authCode.Scopes, " "),
	}
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(tokenString string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	token, ok := s.tokens[tokenString]
	if !ok {
		return nil, errors.New("Token were not found")
	} else if token.Expired() {
		return nil, errors.New("Token expired")
	}
	return token.Clone(), nil
}

func (token *Token) Clone() *Token {
	return &Token{
		AccessToken: token.AccessToken,
		ClientID:    token.ClientID,
		UserID:      token.UserID,
		Scopes:      append([]string(nil), token.Scopes...),
		ExpiresAt:   token.ExpiresAt,
	}
}

func (t *Token) Expired() bool {
	return t.ExpiresAt.Before(time.Now())
}

func (t *RefreshToken) Expired() bool {
	return t.ExpiresAt.Before(time.Now())
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if isRefreshToken {
		_, ok := s.refreshTokens[token]
		delete(s.refreshTokens, token)
		if ok {
			return nil
		}
	} else {
		_, ok := s.tokens[token]
		delete(s.tokens, token)
		if ok {
			return nil
		}
	}
	return errors.New("Token not found")
}

func (s *OAuth2Server) ConsumeVerifiedCode(r *http.Request, client *OAuth2ClientInfo) (*AuthorizationCode, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	code := r.FormValue("code")
	authCode, ok := s.authCodes[code]
	if !ok {
		return nil, false
	}
	if authCode.ClientID != client.ClientID || authCode.RedirectURI != r.FormValue("redirect_uri") {
		return nil, false
	}
	if time.Now().After(authCode.ExpiresAt) {
		return nil, false
	}
	codeVerifier := r.FormValue("code_verifier")
	if !VerifyCodeChallenge(codeVerifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
		return nil, false
	}
	delete(s.authCodes, code)
	return authCode, true
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	var hashedCode string
	switch method {
	case "S256":
		hash32Bytes := sha256.Sum256([]byte(codeVerifier))
		hashBytes := hash32Bytes[:]
		hashedCode = base64.RawURLEncoding.EncodeToString(hashBytes)
	case "plain":
		hashedCode = codeVerifier
	}

	return hashedCode == codeChallenge
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
