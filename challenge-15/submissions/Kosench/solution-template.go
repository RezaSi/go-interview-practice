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
// In production this would be a protected admin panel
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	//Validation of required field
	if client.ClientID == "" {
		return errors.New("client_id is required")
	}
	if client.ClientSecret == "" {
		return errors.New("client_secret is required")
	}
	if len(client.RedirectURIs) == 0 {
		return errors.New("at least one redirect_uri is required")
	}

	// Check that a client with this ID is not already registered
	if _, exists := s.clients[client.ClientID]; exists {
		return fmt.Errorf("client with id %s already exists", client.ClientID)
	}

	s.clients[client.ClientID] = client
	return nil
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	// Fill bytes with random values
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes, %w", err)
	}

	// Convert bytes to characters from the alphabet
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b), nil
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 1. Validate request parameters (client_id, redirect_uri, response_type, scope, state)
	responseType := query.Get("response_type")
	if responseType != "code" {
		// If there is a redirect_uri, redirect with error, otherwise return 400
		redirectURI := query.Get("redirect_uri")
		if redirectURI != "" {
			s.redirectWithError(w, r, redirectURI, "unsupported_response_type",
				"Only 'code' response type is supported", query.Get("state"))
			return
		}
		http.Error(w, "unsupported_response_type", http.StatusBadRequest)
		return
	}

	clientID := query.Get("client_id")
	if clientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()
	if !exists {
		http.Error(w, "invalid_client", http.StatusBadRequest)
		return
	}

	redirectURI := query.Get("redirect_uri")
	if redirectURI == "" {
		http.Error(w, "redirect_uri is required", http.StatusBadRequest)
		return
	}

	if !s.isValidRedirectURI(client, redirectURI) {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}

	state := query.Get("state")

	scopeStr := query.Get("scope")
	requestedScopes := strings.Split(scopeStr, " ")

	for _, scope := range requestedScopes {
		if scope != "" && !s.isScopeAllowed(client, scope) {
			s.redirectWithError(w, r, redirectURI, "invalid_scope",
				fmt.Sprintf("Scope '%s' is not allowed", scope), state)
			return
		}
	}
	// PKCE parameters (optional but recommended)
	codeChallenge := query.Get("code_challenge")
	codeChallengeMethod := query.Get("code_challenge_method")

	// If code_challenge is specified, method is required
	if codeChallenge != "" && codeChallengeMethod == "" {
		s.redirectWithError(w, r, redirectURI, "invalid_request",
			"code_challenge_method is required when code_challenge is present", state)
		return
	}

	// 2. Authenticate the user (for this challenge, could be a simple login form)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		// In production, a login form would be shown here
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	code, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, "failed to generate code", http.StatusInternalServerError)
		return
	}

	authCode := &AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		Scopes:              requestedScopes,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	s.mu.Lock()
	s.authCodes[code] = authCode
	s.mu.Unlock()

	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "invalid redirect_uri", http.StatusInternalServerError)
		return
	}
	q := redirectURL.Query()
	q.Set("code", code)
	if state != "" {
		q.Set("state", state)
	}
	redirectURL.RawQuery = q.Encode()

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

func (s *OAuth2Server) isValidRedirectURI(client *OAuth2ClientInfo, redirectURI string) bool {
	for _, allowed := range client.RedirectURIs {
		if allowed == redirectURI {
			return true
		}
	}
	return false
}

func (s *OAuth2Server) isScopeAllowed(client *OAuth2ClientInfo, scope string) bool {
	for _, allowed := range client.AllowedScopes {
		if allowed == scope {
			return true
		}
	}
	return false
}

func (s *OAuth2Server) redirectWithError(w http.ResponseWriter, r *http.Request,
	redirectURI, errorCode, errorDescription, state string) {
	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	q := redirectURL.Query()
	q.Set("error", errorCode)
	if errorDescription != "" {
		q.Set("error_description", errorDescription)
	}
	if state != "" {
		q.Set("state", state)
	}
	redirectURL.RawQuery = q.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

// HandleToken handles the token endpoint
// This is where the authorization code is exchanged for tokens
// Or the access token is refreshed using a refresh token
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	// Token endpoint only accepts POST requests
	if r.Method != http.MethodPost {
		s.tokenError(w, "invalid_request", "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		s.tokenError(w, "invalid_request", "Failed to parse form", http.StatusBadRequest)
		return
	}

	// grant_type determines the request type
	grantType := r.FormValue("grant_type")

	switch grantType {
	case "authorization_code":
		// Exchange authorization code for tokens
		s.handleAuthorizationCodeGrant(w, r)
	case "refresh_token":
		// Refresh access token using refresh token
		s.handleRefreshTokenGrant(w, r)
	default:
		s.tokenError(w, "unsupported_grant_type",
			fmt.Sprintf("Grant type '%s' is not supported", grantType), http.StatusBadRequest)
	}
}

func (s *OAuth2Server) handleAuthorizationCodeGrant(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		s.tokenError(w, "invalid_request", "code is required", http.StatusBadRequest)
		return
	}

	redirectURI := r.FormValue("redirect_uri")
	if redirectURI == "" {
		s.tokenError(w, "invalid_request", "redirect_uri is required", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	if clientID == "" || clientSecret == "" {
		s.tokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
		return
	}

	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()

	if !exists {
		s.tokenError(w, "invalid_client", "Client not found", http.StatusUnauthorized)
		return
	}

	// Check client_secret (in production this would be a bcrypt hash)
	if client.ClientSecret != clientSecret {
		s.tokenError(w, "invalid_client", "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	// 3. VALIDATE AUTHORIZATION CODE
	s.mu.Lock()
	authCode, exists := s.authCodes[code]
	if !exists {
		s.mu.Unlock()
		s.tokenError(w, "invalid_grant", "Authorization code is invalid or expired", http.StatusBadRequest)
		return
	}

	// One-time use code - delete immediately (protection against replay attacks)
	delete(s.authCodes, code)
	s.mu.Unlock()

	// Check expiration
	if time.Now().After(authCode.ExpiresAt) {
		s.tokenError(w, "invalid_grant", "Authorization code has expired", http.StatusBadRequest)
		return
	}

	// Check that the code was issued to this client
	if authCode.ClientID != clientID {
		s.tokenError(w, "invalid_grant", "Authorization code was issued to another client", http.StatusBadRequest)
		return
	}

	// Check redirect_uri (should match the one from when the code was obtained)
	if authCode.RedirectURI != redirectURI {
		s.tokenError(w, "invalid_grant", "redirect_uri mismatch", http.StatusBadRequest)
		return
	}

	// 4. CHECK PKCE (if used)

	if authCode.CodeChallenge != "" {
		codeVerifier := r.FormValue("code_verifier")
		if codeVerifier == "" {
			s.tokenError(w, "invalid_grant", "code_verifier is required", http.StatusBadRequest)
			return
		}

		// Verify PKCE challenge
		if !VerifyCodeChallenge(codeVerifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
			s.tokenError(w, "invalid_grant", "Invalid code_verifier", http.StatusBadRequest)
			return
		}
	}

	// 5. GENERATE TOKENS

	// Generate access token
	accessTokenStr, err := GenerateRandomString(32)
	if err != nil {
		s.tokenError(w, "server_error", "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Access token lasts 1 hour (typical for OAuth2)
	accessToken := &Token{
		AccessToken: accessTokenStr,
		ClientID:    clientID,
		UserID:      authCode.UserID,
		Scopes:      authCode.Scopes,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	// Generate refresh token
	refreshTokenStr, err := GenerateRandomString(64)
	if err != nil {
		s.tokenError(w, "server_error", "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Refresh token lasts 30 days
	refreshToken := &RefreshToken{
		RefreshToken: refreshTokenStr,
		ClientID:     clientID,
		UserID:       authCode.UserID,
		Scopes:       authCode.Scopes,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}

	// Save tokens
	s.mu.Lock()
	s.tokens[accessTokenStr] = accessToken
	s.refreshTokens[refreshTokenStr] = refreshToken
	s.mu.Unlock()

	// 6. RETURN TOKENS TO CLIENT

	response := map[string]interface{}{
		"access_token":  accessTokenStr,
		"token_type":    "Bearer", // Token type - always Bearer for OAuth2
		"expires_in":    3600,     // Expiration time in seconds (1 hour)
		"refresh_token": refreshTokenStr,
		"scope":         strings.Join(authCode.Scopes, " "),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store") // Tokens should not be cached!
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error - response headers already sent, cannot change status
		fmt.Printf("failed to encode token response: %v\n", err)
	}
}

// handleRefreshTokenGrant handles access token refresh
// This allows getting a new access token without re-authorizing the user
func (s *OAuth2Server) handleRefreshTokenGrant(w http.ResponseWriter, r *http.Request) {
	// 1. EXTRACT PARAMETERS

	refreshTokenStr := r.FormValue("refresh_token")
	if refreshTokenStr == "" {
		s.tokenError(w, "invalid_request", "refresh_token is required", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	// 2. CLIENT AUTHENTICATION

	if clientID == "" || clientSecret == "" {
		s.tokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
		return
	}

	s.mu.RLock()
	client, exists := s.clients[clientID]
	s.mu.RUnlock()

	if !exists || client.ClientSecret != clientSecret {
		s.tokenError(w, "invalid_client", "Invalid client credentials", http.StatusUnauthorized)
		return
	}

	// 3. REFRESH TOKEN VALIDATION

	s.mu.Lock()
	refreshToken, exists := s.refreshTokens[refreshTokenStr]
	if !exists {
		s.mu.Unlock()
		s.tokenError(w, "invalid_grant", "Refresh token is invalid", http.StatusBadRequest)
		return
	}

	// Remove old refresh token (refresh token rotation - best practice)
	// This protects against replay attacks
	delete(s.refreshTokens, refreshTokenStr)
	s.mu.Unlock()

	// Check expiration
	if time.Now().After(refreshToken.ExpiresAt) {
		s.tokenError(w, "invalid_grant", "Refresh token has expired", http.StatusBadRequest)
		return
	}

	// Check that the token belongs to this client
	if refreshToken.ClientID != clientID {
		s.tokenError(w, "invalid_grant", "Refresh token was issued to another client", http.StatusBadRequest)
		return
	}

	// 4. GENERATE NEW TOKENS

	// Generate new access token
	accessTokenStr, err := GenerateRandomString(32)
	if err != nil {
		s.tokenError(w, "server_error", "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	accessToken := &Token{
		AccessToken: accessTokenStr,
		ClientID:    clientID,
		UserID:      refreshToken.UserID,
		Scopes:      refreshToken.Scopes,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	// Generate new refresh token (refresh token rotation)
	newRefreshTokenStr, err := GenerateRandomString(64)
	if err != nil {
		s.tokenError(w, "server_error", "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	newRefreshToken := &RefreshToken{
		RefreshToken: newRefreshTokenStr,
		ClientID:     clientID,
		UserID:       refreshToken.UserID,
		Scopes:       refreshToken.Scopes,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}

	// Save new tokens
	s.mu.Lock()
	s.tokens[accessTokenStr] = accessToken
	s.refreshTokens[newRefreshTokenStr] = newRefreshToken
	s.mu.Unlock()

	// 5. RETURN NEW TOKENS

	response := map[string]interface{}{
		"access_token":  accessTokenStr,
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": newRefreshTokenStr, // Return NEW refresh token
		"scope":         strings.Join(refreshToken.Scopes, " "),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	json.NewEncoder(w).Encode(response)
}

func (s *OAuth2Server) tokenError(w http.ResponseWriter, errorCode, description string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             errorCode,
		"error_description": description,
	})
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(tokenStr string) (*Token, error) {
	s.mu.RLock()
	token, exists := s.tokens[tokenStr]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.New("token not found")
	}

	// Check expiration
	if time.Now().After(token.ExpiresAt) {
		// Remove expired token
		s.mu.Lock()
		delete(s.tokens, tokenStr)
		s.mu.Unlock()
		return nil, errors.New("token has expired")
	}

	return token, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshTokenStr string) (*Token, *RefreshToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check existence of refresh token
	refreshToken, exists := s.refreshTokens[refreshTokenStr]
	if !exists {
		return nil, nil, errors.New("refresh token not found")
	}

	// Check expiration
	if time.Now().After(refreshToken.ExpiresAt) {
		delete(s.refreshTokens, refreshTokenStr)
		return nil, nil, errors.New("refresh token has expired")
	}

	// Remove old refresh token (rotation for security)
	delete(s.refreshTokens, refreshTokenStr)

	// Generate new access token
	accessTokenStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newAccessToken := &Token{
		AccessToken: accessTokenStr,
		ClientID:    refreshToken.ClientID,
		UserID:      refreshToken.UserID,
		Scopes:      refreshToken.Scopes,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	// Generate new refresh token
	newRefreshTokenStr, err := GenerateRandomString(64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	newRefreshToken := &RefreshToken{
		RefreshToken: newRefreshTokenStr,
		ClientID:     refreshToken.ClientID,
		UserID:       refreshToken.UserID,
		Scopes:       refreshToken.Scopes,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}

	// Save new tokens
	s.tokens[accessTokenStr] = newAccessToken
	s.refreshTokens[newRefreshTokenStr] = newRefreshToken

	return newAccessToken, newRefreshToken, nil
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(tokenStr string, isRefreshToken bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if isRefreshToken {
		if _, exists := s.refreshTokens[tokenStr]; !exists {
			return errors.New("refresh token not found")
		}
		delete(s.refreshTokens, tokenStr)
	} else {
		if _, exists := s.tokens[tokenStr]; !exists {
			return errors.New("access token not found")
		}
		delete(s.tokens, tokenStr)
	}

	return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	switch method {
	case "S256":
		h := sha256.New()
		h.Write([]byte(codeVerifier))
		computed := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
		return computed == codeChallenge

	case "plain":
		// Plain method - just string comparison
		// Less secure, but supported for backward compatibility
		return codeVerifier == codeChallenge

	default:
		// Unsupported method
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
// This is the first step - user clicks "Sign in with OAuth" and gets redirected here
func (c *OAuth2Client) GetAuthorizationURL(state string, codeChallenge string, codeChallengeMethod string) (string, error) {
	if c.Config.AuthorizationEndpoint == "" {
		return "", errors.New("authorization endpoint is not configured")
	}

	// Parse base URL
	authURL, err := url.Parse(c.Config.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("invalid authorization endpoint: %w", err)
	}

	// Build query parameters
	query := authURL.Query()
	query.Set("response_type", "code")                     // Always "code" for authorization code flow
	query.Set("client_id", c.Config.ClientID)              // ID of our application
	query.Set("redirect_uri", c.Config.RedirectURI)        // Where to return after authorization
	query.Set("scope", strings.Join(c.Config.Scopes, " ")) // Which permissions to request
	query.Set("state", state)                              // CSRF protection - random string

	// Add PKCE parameters if specified (for mobile/SPA applications)
	if codeChallenge != "" {
		query.Set("code_challenge", codeChallenge)
		query.Set("code_challenge_method", codeChallengeMethod)
	}

	authURL.RawQuery = query.Encode()
	return authURL.String(), nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
// This is called on the callback endpoint after the user has authorized
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	if c.Config.TokenEndpoint == "" {
		return errors.New("token endpoint is not configured")
	}

	// Build form data for POST request
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", c.Config.RedirectURI)
	form.Set("client_id", c.Config.ClientID)
	form.Set("client_secret", c.Config.ClientSecret)

	// Add code_verifier if PKCE was used
	if codeVerifier != "" {
		form.Set("code_verifier", codeVerifier)
	}

	// Make POST request to token endpoint
	resp, err := http.PostForm(c.Config.TokenEndpoint, form)
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf("token exchange failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
	}

	// Parse successful response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Save tokens in the client
	c.AccessToken = tokenResp.AccessToken
	c.RefreshToken = tokenResp.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// RefreshToken refreshes the access token using the refresh token
// Called automatically when the access token expires
func (c *OAuth2Client) DoRefreshToken() error {
	if c.RefreshToken == "" {
		return errors.New("no refresh token available")
	}

	if c.Config.TokenEndpoint == "" {
		return errors.New("token endpoint is not configured")
	}

	// Build form data
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", c.RefreshToken)
	form.Set("client_id", c.Config.ClientID)
	form.Set("client_secret", c.Config.ClientSecret)

	// Make POST request
	resp, err := http.PostForm(c.Config.TokenEndpoint, form)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf("token refresh failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
	}

	// Parse new tokens
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Update tokens (including the new refresh token!)
	c.AccessToken = tokenResp.AccessToken
	c.RefreshToken = tokenResp.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(urlStr string, method string) (*http.Response, error) {
	// Check if the token has expired (with a 30 second buffer)
	if time.Now().Add(30 * time.Second).After(c.TokenExpiry) {
		// Token has expired or will expire soon - refresh it
		if err := c.DoRefreshToken(); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
	}

	// Create HTTP request
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add access token to Authorization header
	// Format: "Bearer <access_token>" - OAuth2 standard
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// If we get 401 Unauthorized - the token may have been invalidated
	// Try to refresh and repeat the request once
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()

		if err := c.DoRefreshToken(); err != nil {
			return nil, fmt.Errorf("token refresh after 401 failed: %w", err)
		}

		// Create a new request with the new token
		retryReq, err := http.NewRequest(method, urlStr, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create retry request: %w", err)
		}
		retryReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
		resp, err = client.Do(retryReq)
		if err != nil {
			return nil, fmt.Errorf("retry request failed: %w", err)
		}
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
	if err := server.RegisterClient(client); err != nil {
		fmt.Printf("Failed to register client: %v\n", err)
		return
	}

	// Start the server (blocks)
	fmt.Println("Starting OAuth2 server on port 9000")
	if err := server.StartServer(9000); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
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
