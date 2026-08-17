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

func (c *OAuth2ClientInfo) hasScope(scope string) bool {
	hasScope := false

	for _, s := range c.AllowedScopes {
		if s == scope {
			hasScope = true
		}
	}

	return hasScope
}

func (c *OAuth2ClientInfo) hasRedirectURI(URI string) bool {
	hasURI := false

	for _, ru := range c.RedirectURIs {
		if ru == URI {
			hasURI = true
		}
	}

	return hasURI
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

func (u *User) passwordMatches(password string) bool {
	return u.Password == password
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

func (s *OAuth2Server) findUserByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New("not found")
}

// RegisterClient registers a new OAuth2 client
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	// TODO: Implement client registration
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[client.ClientID]; ok {
		return errors.New("duplicate client")
	}

	s.clients[client.ClientID] = client

	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	// TODO: Implement secure random string generation
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range length {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func hashCodeChallenge(codeChallenge, method string) (string, error) {
	switch method {
	case "S256":
		hasher := sha256.New()
		_, err := hasher.Write([]byte(codeChallenge))
		if err != nil {
			return "", err
		}
		hash := hasher.Sum(nil)

		codeChallenge := base64.RawURLEncoding.EncodeToString(hash)
		return codeChallenge, nil

	case "plain":
		return codeChallenge, nil
	default:
		return "", errors.New("unsupported code challenge method")
	}
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement authorization endpoint
	// 1. Validate request parameters (client_id, redirect_uri, response_type, scope, state)
	clientId := r.URL.Query().Get("client_id")
	if clientId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	client, ok := s.clients[clientId]
	if !ok {
		s.mu.RUnlock()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.mu.RUnlock()

	redirectURI := r.URL.Query().Get("redirect_uri")

	if !client.hasRedirectURI(redirectURI) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scope := r.URL.Query().Get("scope")
	if scope == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !client.hasScope(scope) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resoponseType := r.URL.Query().Get("response_type")
	if resoponseType == "" || resoponseType != "code" {
		URL := fmt.Sprintf("%s?error=%s", redirectURI, "unsupported_response_type")
		http.Redirect(w, r, URL, http.StatusFound)
		return
	}

	code, err := GenerateRandomString(32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")

	codeChallengeHash, err := hashCodeChallenge(codeChallenge, codeChallengeMethod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authCode := &AuthorizationCode{
		ClientID:            clientId,
		Code:                code,
		RedirectURI:         redirectURI,
		Scopes:              []string{scope},
		CodeChallenge:       codeChallengeHash,
		CodeChallengeMethod: codeChallengeMethod,
	}

	// 2. Authenticate the user (for this challenge, could be a simple login form)
	if r.Method == http.MethodPost {
		userId := r.Context().Value("user_id")
		if userId == "" {
			fmt.Fprintf(w, `
			<form action="/authorize?client_id=%s&redirect_uri=%s&response_type=%s&scope=%s&state=%s" method="POST">
      <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" name="username" required autocomplete="username">
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" name="password" required autocomplete="current-password">
      </div>

      <button type="submit">Log In</button>
    </form>`, clientId, redirectURI, resoponseType, scope, state)
			return
		}

		s.mu.Lock()
		user, ok := s.users[userId.(string)]
		if !ok {
			s.mu.Unlock()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		s.mu.Unlock()

		authCode.UserID = user.ID

	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := s.findUserByUsername(username)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !user.passwordMatches(password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authCode.UserID = user.ID
	}

	authCode.ExpiresAt = time.Now().Add(10 * time.Minute)

	// 3. Present a consent screen to the user
	// 4. Generate an authorization code and redirect to the client with the code
	s.mu.Lock()
	s.authCodes[authCode.Code] = authCode
	s.mu.Unlock()

	formattedURL := fmt.Sprintf("%s?code=%s&state=%s", authCode.RedirectURI, authCode.Code, state)
	http.Redirect(w, r, formattedURL, http.StatusFound)
}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	grantType := r.FormValue("grant_type")
	if grantType == "" {
		WriteError(w, http.StatusBadRequest, ErrorResponse{
			Error:       "invalid_grant",
			Description: "grant type must be authorization_code",
		})
		return
	}

	switch grantType {
	case "authorization_code":
		s.handleCodeExchange(w, r)

	case "refresh_token":
		s.handleRefreshToken(w, r)
	}

}

func (s *OAuth2Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	grantType := r.FormValue("grant_type")
	refreshToken := r.FormValue("refresh_token")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	if grantType != "refresh_token" {
		WriteError(w, http.StatusBadRequest, ErrorResponse{
			Error:       "invalid_grant",
			Description: "grant type must be refresh_token",
		})
		return
	}

	if refreshToken == "" {
		WriteError(w, http.StatusBadRequest, ErrorResponse{
			Error:       "invalid_grant",
			Description: "refresh token is required",
		})
		return
	}

	if clientID == "" {
		WriteError(w, http.StatusBadRequest, ErrorResponse{
			Error:       "invalid_client",
			Description: "client_id is required",
		})
		return
	}

	if clientSecret == "" {
		WriteError(w, http.StatusBadRequest, ErrorResponse{
			Error:       "invalid_client",
			Description: "client_secret is required",
		})
		return
	}

	s.mu.RLock()
	client, ok := s.clients[clientID]
	if !ok {
		s.mu.RUnlock()
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_client",
			Description: "no client",
		})
		return
	}
	if client.ClientSecret != clientSecret {
		s.mu.RUnlock()
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_client",
			Description: "invalid client secret",
		})
		return
	}
	s.mu.RUnlock()

	refreshTokenObj, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_token",
			Description: err.Error(),
		})
		return
	}

	newRefreshTokenCode, err := GenerateRandomString(32)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, ErrorResponse{
			Error:       "internal_server_error",
			Description: err.Error(),
		})
		return
	}

	newRefreshToken := &RefreshToken{
		RefreshToken: newRefreshTokenCode,
		ClientID:     clientID,
		UserID:       refreshTokenObj.UserID,
		Scopes:       refreshTokenObj.Scopes,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.refreshTokens[newRefreshTokenCode] = newRefreshToken
	s.mu.Unlock()

	accessToken, _, err := s.RefreshAccessToken(newRefreshTokenCode)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, ErrorResponse{
			Error:       "internal_server_error",
			Description: err.Error(),
		})
		return
	}

	if err := s.RevokeToken(refreshToken, true); err != nil {
		WriteError(w, http.StatusInternalServerError, ErrorResponse{
			Error:       "internal_server_error",
			Description: err.Error(),
		})
		return
	}

	WriteResponse(w, http.StatusOK, TokenRespone{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: newRefreshTokenCode,
		TokenType:    "Bearer",
		ExpiresIn:    accessToken.ExpiresAt.Second(),
		Scope:        strings.Join(accessToken.Scopes, " "),
	})

}

func (s *OAuth2Server) handleCodeExchange(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")
	codeVerifier := r.FormValue("code_verifier")

	s.mu.RLock()
	client, ok := s.clients[clientID]
	if !ok {
		s.mu.RUnlock()
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_client",
			Description: "no client",
		})
		return
	}
	s.mu.RUnlock()

	if client.ClientSecret != clientSecret {
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_client",
			Description: "invalid client secret",
		})
	}

	if !client.hasRedirectURI(redirectURI) {
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_grant",
			Description: "redirect URI not allowed",
		})
		return
	}

	s.mu.RLock()
	authCode, ok := s.authCodes[code]
	if !ok {
		s.mu.RUnlock()
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_grant",
			Description: "code not found",
		})
		return
	}
	s.mu.RUnlock()

	if !time.Now().Before(authCode.ExpiresAt) {
		WriteError(w, http.StatusUnauthorized, ErrorResponse{
			Error:       "invalid_grant",
			Description: "code expired",
		})
		return
	}

	if !VerifyCodeChallenge(codeVerifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
		WriteError(w, http.StatusBadRequest, ErrorResponse{
			Error:       "invalid_grant",
			Description: "code verifier mismatch",
		})
		return
	}

	accessTokenCode, err := GenerateRandomString(32)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, ErrorResponse{
			Error:       "internal server error",
			Description: "cannot generate token",
		})
		return
	}
	token := &Token{
		AccessToken: accessTokenCode,
		ClientID:    clientID,
		UserID:      authCode.UserID,
		Scopes:      authCode.Scopes,
		ExpiresAt:   time.Now().Add(15 * time.Minute),
	}

	refreshTokenCode, err := GenerateRandomString(32)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, ErrorResponse{
			Error:       "internal server error",
			Description: "cannot generate token",
		})
		return
	}
	refreshToken := &RefreshToken{
		RefreshToken: refreshTokenCode,
		ClientID:     clientID,
		UserID:       authCode.UserID,
		Scopes:       authCode.Scopes,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	s.mu.Lock()
	s.refreshTokens[refreshTokenCode] = refreshToken
	s.tokens[accessTokenCode] = token
	delete(s.authCodes, authCode.Code)
	s.mu.Unlock()

	WriteResponse(w, http.StatusOK, TokenRespone{
		AccessToken:  accessTokenCode,
		TokenType:    "Bearer",
		ExpiresIn:    token.ExpiresAt.Second(),
		RefreshToken: refreshTokenCode,
		Scope:        strings.Join(authCode.Scopes, " "),
	})
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	accessToken, ok := s.tokens[token]
	if !ok {
		return nil, errors.New("token not found")
	}

	if accessToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	return accessToken, nil
}

// ValidateRefreshToken validates an refrech token
func (s *OAuth2Server) ValidateRefreshToken(token string) (*RefreshToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	refreshToken, ok := s.refreshTokens[token]
	if !ok {
		return nil, errors.New("token not found")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	return refreshToken, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	refreshTokenObj, ok := s.refreshTokens[refreshToken]
	if !ok {
		return nil, nil, errors.New("invalid refresh token")
	}

	if refreshTokenObj.ExpiresAt.Before(time.Now()) {
		return nil, nil, errors.New("expired refresh token")
	}

	accessTokenCode, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}

	newAccessToken := &Token{
		AccessToken: accessTokenCode,
		ClientID:    refreshTokenObj.ClientID,
		UserID:      refreshTokenObj.UserID,
		Scopes:      refreshTokenObj.Scopes,
		ExpiresAt:   time.Now().Add(15 * time.Minute),
	}

	return newAccessToken, refreshTokenObj, nil
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
		return errors.New("refresh token not found")
	}

	if _, ok := s.tokens[token]; ok {
		delete(s.tokens, token)
		return nil
	}
	return errors.New("token not found")
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	matches := false

	hash, err := hashCodeChallenge(codeVerifier, method)
	if err == nil && hash == codeChallenge {
		matches = true
	}

	return matches
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

type TokenRespone struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func WriteResponse(w http.ResponseWriter, statusCode int, response TokenRespone) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, statusCode int, errResponse ErrorResponse) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResponse)
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
	if c.Config.AuthorizationEndpoint == "" {
		return "", errors.New("authorization endpoint is required")
	}
	if c.Config.ClientID == "" {
		return "", errors.New("client ID is required")
	}

	u, err := url.Parse(c.Config.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("invalid authorization endpoint: %w", err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.Config.ClientID)
	q.Set("redirect_uri", c.Config.RedirectURI)
	if len(c.Config.Scopes) > 0 {
		q.Set("scope", strings.Join(c.Config.Scopes, " "))
	}
	if state != "" {
		q.Set("state", state)
	}
	if codeChallenge != "" {
		q.Set("code_challenge", codeChallenge)
	}
	if codeChallengeMethod != "" {
		q.Set("code_challenge_method", codeChallengeMethod)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	if c.Config.TokenEndpoint == "" {
		return errors.New("token endpoint is required")
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.Config.RedirectURI)
	data.Set("client_id", c.Config.ClientID)
	data.Set("client_secret", c.Config.ClientSecret)
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	resp, err := http.PostForm(c.Config.TokenEndpoint, data)
	if err != nil {
		return fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("token request failed (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.Description)
		}
		return fmt.Errorf("token request failed with status code: %d", resp.StatusCode)
	}

	var tokenResp TokenRespone
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	c.RefreshToken = tokenResp.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	if c.Config.TokenEndpoint == "" {
		return errors.New("token endpoint is required")
	}
	if c.RefreshToken == "" {
		return errors.New("no refresh token available")
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.RefreshToken)
	data.Set("client_id", c.Config.ClientID)
	data.Set("client_secret", c.Config.ClientSecret)

	resp, err := http.PostForm(c.Config.TokenEndpoint, data)
	if err != nil {
		return fmt.Errorf("failed to make refresh token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("refresh token request failed (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.Description)
		}
		return fmt.Errorf("refresh token request failed with status code: %d", resp.StatusCode)
	}

	var tokenResp TokenRespone
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		c.RefreshToken = tokenResp.RefreshToken
	}
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(reqURL string, method string) (*http.Response, error) {
	if c.AccessToken == "" {
		return nil, errors.New("no access token available")
	}

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
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
		fmt.Spt
				// After authorization, exchange the code for tokens
				client.ExchangeCodeForToken("returned-code", codeVerifier)

				// Make an authenticated request
				resp, _ := client.MakeAuthenticatedRequest("http://api.example.com/resource", "GET")
				fmt.Printf("Response: %v\n", resp)
	*/
}
