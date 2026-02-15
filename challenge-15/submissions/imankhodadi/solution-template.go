package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
)

type OAuth2Config struct {
	// AuthorizationEndpoint is the endpoint for authorization requests
	AuthorizationEndpoint string
	// TokenEndpoint is the endpoint for token requests
	TokenEndpoint string
	ClientID      string
	ClientSecret  string
	RedirectURI   string
	Scopes        []string
}

// TODO:
// No mechanism to evict expired tokens — memory grows unbounded.
// Expired authorization codes, access tokens, and refresh tokens are never cleaned up
// (except authorization codes which are deleted on use).
// In a long-running process, the in-memory maps will grow indefinitely.
// Consider adding a periodic background goroutine that sweeps expired entries,
// or check-and-evict lazily during lookups.

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
	mu    sync.RWMutex
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
type User struct {
	ID       string
	Username string
	Password string
}

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
	AccessToken  string
	ClientID     string
	UserID       string
	Scopes       []string
	ExpiresAt    time.Time
	RefreshToken string
}
type RefreshToken struct {
	RefreshToken string
	ClientID     string
	UserID       string
	Scopes       []string // Scopes is a list of authorized scopes
	ExpiresAt    time.Time
	AccessToken  string
}

func NewOAuth2Server() *OAuth2Server {
	server := &OAuth2Server{
		clients:       make(map[string]*OAuth2ClientInfo),
		authCodes:     make(map[string]*AuthorizationCode),
		tokens:        make(map[string]*Token),
		refreshTokens: make(map[string]*RefreshToken),
		users:         make(map[string]*User),
	}
	server.users["user1"] = &User{
		ID:       "user1",
		Username: "testuser",
		Password: "password", // just for test
	}
	return server
}

// registers a new OAuth2 client
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[client.ClientID]; ok {
		return fmt.Errorf("client already exists")
	}
	s.clients[client.ClientID] = client
	return nil
}

// TODO: You allocate length random bytes but then base64-encode them (which expands the output) and truncate back to length characters.
// Each base64 character carries 6 bits of entropy, so a 32-character token has ~192 bits of effective randomness—not the ~256 bits
// you'd expect from 32 random bytes. This is still more than sufficient for token security, but if the intent is to control entropy
// rather than string length, consider allocating ceil(length * 3/4) bytes instead (to avoid wasting randomness), or document the effective entropy.
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	//  Implement authorization endpoint
	// 1. Validate request parameters (client_id, redirect_uri, response_type, scope, state)
	// 2. Authenticate the user (for this challenge, could be a simple login form)
	// 3. Present a consent screen to the user
	// 4. Generate an authorization code and redirect to the client with the code
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")
	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")
	if codeChallenge != "" && codeChallengeMethod == "" {
		codeChallengeMethod = "plain"
		// RFC 7636, when code_challenge is present but code_challenge_method is omitted,
		// the default must be "plain"
	}
	s.mu.RLock()
	client, ok := s.clients[clientID]
	s.mu.RUnlock()
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
	// TODO: make a default scope. When scope is an empty string, strings.Fields returns []string{}, the loop body never executes,
	// and the authorization code is issued with no scopes. Depending on your requirements, you may want to require at
	// least one valid scope or assign a default.
	reqScopes := strings.Fields(scope)
	for _, rs := range reqScopes {
		if !slices.Contains(client.AllowedScopes, rs) {
			http.Redirect(w, r, fmt.Sprintf("%s?error=invalid_scope&state=%s", redirectURI, url.QueryEscape(state)), http.StatusFound)
			return
		}
	}
	code, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	authCode := &AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              "user1", // TODO: replace with actual user authentication
		RedirectURI:         redirectURI,
		Scopes:              reqScopes,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	s.mu.Lock()
	s.authCodes[code] = authCode
	s.mu.Unlock()
	v := url.Values{}
	v.Set("code", code)
	if state != "" {
		v.Set("state", state)
	}
	http.Redirect(w, r, redirectURI+"?"+v.Encode(), http.StatusFound)
}

func writeTokenResponse(w http.ResponseWriter, resp *TokenResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding token response: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, errCode, desc string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
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
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	// 1. Validate request parameters (grant_type, code, redirect_uri, client_id, client_secret)
	// 2. Verify the authorization code
	// 3. For PKCE, verify the code_verifier
	// 4. Generate access and refresh tokens
	// 5. Return the tokens as a JSON response
	if r.Method != http.MethodPost {
		writeJSONError(w, "invalid_request", "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		writeJSONError(w, "invalid_request", "error parsing form", http.StatusBadRequest)
		return
	}
	grantType := r.Form.Get("grant_type")
	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")
	s.mu.RLock()
	client, ok := s.clients[clientID]
	s.mu.RUnlock()
	if !ok {
		writeJSONError(w, "invalid_client", "invalid client", http.StatusUnauthorized)
		return
	}
	// instead of "client.ClientSecret != clientSecret", Use crypto/subtle.ConstantTimeCompare for secret comparisons
	if subtle.ConstantTimeCompare([]byte(client.ClientSecret), []byte(clientSecret)) != 1 {
		writeJSONError(w, "invalid_client", "invalid client", http.StatusUnauthorized)
		return
	}
	switch grantType {
	case "refresh_token":
		refToken := r.Form.Get("refresh_token")

		accessToken, refreshToken, err := s.RefreshAccessToken(refToken, clientID)

		if err != nil {
			writeJSONError(w, "invalid_grant", "invalid refresh token", http.StatusBadRequest)
			return
		}
		tokenType := "Bearer"
		expiresIn := 3600
		scope := strings.Join(refreshToken.Scopes, " ")
		resp := &TokenResponse{
			AccessToken:  accessToken.AccessToken,
			TokenType:    tokenType,
			ExpiresIn:    expiresIn,
			RefreshToken: refreshToken.RefreshToken,
			Scope:        scope,
		}
		writeTokenResponse(w, resp)
	case "authorization_code":
		code := r.Form.Get("code")
		redirectURI := r.Form.Get("redirect_uri")
		codeVerifier := r.Form.Get("code_verifier")
		s.mu.Lock()
		authReq, ok := s.authCodes[code]
		// Always delete the code on first use attempt to enforce single-use
		delete(s.authCodes, code)
		s.mu.Unlock()
		if !ok || authReq.ExpiresAt.Before(time.Now()) || authReq.RedirectURI != redirectURI || authReq.ClientID != clientID {
			writeJSONError(w, "invalid_grant", "invalid_grant", http.StatusBadRequest)
			return
		}
		if authReq.CodeChallenge != "" {
			if !VerifyCodeChallenge(codeVerifier, authReq.CodeChallenge, authReq.CodeChallengeMethod) {
				writeJSONError(w, "invalid_grant", "invalid_grant", http.StatusBadRequest)
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
			AccessToken:  accessToken,
		}
		s.mu.Lock()
		s.tokens[accessToken] = token
		s.refreshTokens[refreshToken] = rt
		s.mu.Unlock()
		tokenType := "Bearer"
		expiresIn := 3600
		scope := strings.Join(authReq.Scopes, " ")
		resp := &TokenResponse{
			AccessToken:  accessToken,
			TokenType:    tokenType,
			ExpiresIn:    expiresIn,
			RefreshToken: refreshToken,
			Scope:        scope,
		}
		writeTokenResponse(w, resp)
	default:
		writeJSONError(w, "unsupported_grant_type", "unsupported grant type", http.StatusBadRequest)
	}
}

func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	tokenInfo, exists := s.tokens[token]
	s.mu.RUnlock()
	if !exists {
		return nil, errors.New("token not found")
	}
	if tokenInfo.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}
	return tokenInfo, nil
}

func (s *OAuth2Server) RefreshAccessToken(refreshToken string, clientID string) (*Token, *RefreshToken, error) {
	// Refresh token rotation: old refresh token and its associated access token
	// are both revoked when a new pair is issued.
	s.mu.Lock()
	defer s.mu.Unlock()
	rt, exists := s.refreshTokens[refreshToken]
	if !exists {
		return nil, nil, errors.New("invalid refresh token")
	}
	if rt.ExpiresAt.Before(time.Now()) {
		return nil, nil, errors.New("token expired")
	}
	if rt.ClientID != clientID {
		return nil, nil, errors.New("refresh token does not belong to this client")
	}
	//Note: the random generation inside the mutex, which is fine for an in-memory server
	// but worth considering if lock contention becomes a concern.
	// This is called twice while holding the write lock, blocking all other server operations
	// (token validation, authorization, etc.).
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
		AccessToken:  rToken,
	}
	delete(s.refreshTokens, refreshToken)
	delete(s.tokens, rt.AccessToken) // revoke old access token
	s.tokens[rToken] = t
	s.refreshTokens[rRefreshToken] = r
	return t, r, nil
}

// revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {

	if isRefreshToken {
		// Revocation of unknown tokens should be a no-op per RFC 7009.
		// RFC 7009 specifies that the server should respond with 200 even for invalid or already-revoked tokens.
		// Returning an error here forces callers to special-case "token doesn't exist,"
		// which complicates any future /revoke HTTP endpoint built on top of this method.
		// delete flaged lines for future use (required for this assignment)
		s.mu.Lock()
		defer s.mu.Unlock()
		rt, exists := s.refreshTokens[token]
		if !exists {
			return errors.New("token does not exist") // cannot delete this, part of assignment unittests
		}
		if rt.AccessToken != "" {
			delete(s.tokens, rt.AccessToken) // Also revoke the associated access token
		}
		delete(s.refreshTokens, token)
	} else {
		s.mu.Lock()
		defer s.mu.Unlock()
		tokenTemp, exists := s.tokens[token]
		if !exists {
			return errors.New("token does not exist") // cannot delete this, part of assignment unittests
		}
		if tokenTemp.RefreshToken != "" {
			delete(s.refreshTokens, tokenTemp.RefreshToken)
		}
		delete(s.tokens, token)
	}
	return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	//  Implement PKCE verification
	switch method {
	case "S256":
		hashedVerifier := sha256.Sum256([]byte(codeVerifier))
		expectedChallenge := base64.RawURLEncoding.EncodeToString(hashedVerifier[:])
		return subtle.ConstantTimeCompare([]byte(expectedChallenge), []byte(codeChallenge)) == 1
	case "plain":
		return subtle.ConstantTimeCompare([]byte(codeVerifier), []byte(codeChallenge)) == 1
	default:
		return false
	}
}

func (s *OAuth2Server) StartServer(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/authorize", s.HandleAuthorize)
	mux.HandleFunc("/token", s.HandleToken)
	// server with timeout
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv.ListenAndServe()
}

// Client code to demonstrate usage////////////////////////////////////////////////////////////////////////

// OAuth2Client represents a client application using OAuth2
type OAuth2Client struct {
	// Config is the OAuth2 configuration
	Config       OAuth2Config
	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time
	mu           sync.RWMutex
}

func NewOAuth2Client(config OAuth2Config) *OAuth2Client {
	return &OAuth2Client{Config: config}
}

// returns the URL to redirect the user for authorization
func (c *OAuth2Client) GetAuthorizationURL(state string, codeChallenge string, codeChallengeMethod string) (string, error) {
	//  Implement building the authorization URL
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

// exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	//  Implement token exchange
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

	var httpClient = &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&e)
		desc, _ := e["error_description"].(string)
		return fmt.Errorf("token exchange failed: %s", desc)
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AccessToken = tr.AccessToken
	c.RefreshToken = tr.RefreshToken
	if tr.ExpiresIn > 0 {
		c.TokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	} else {
		c.TokenExpiry = time.Time{}
	}
	return nil
}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	c.mu.RLock()
	refreshToken := c.RefreshToken
	cfg := c.Config
	c.mu.RUnlock()
	if refreshToken == "" {
		return errors.New("no refresh token")
	}
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)
	v.Set("client_id", cfg.ClientID)
	v.Set("client_secret", cfg.ClientSecret)
	req, err := http.NewRequest("POST", cfg.TokenEndpoint, strings.NewReader(v.Encode()))

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var httpClient = &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
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
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AccessToken = tr.AccessToken
	if tr.RefreshToken != "" {
		c.RefreshToken = tr.RefreshToken
	}
	if tr.ExpiresIn > 0 {
		c.TokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	} else {
		c.TokenExpiry = time.Time{}
	}
	return nil
}

// makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(targetURL string, method string) (*http.Response, error) {
	//  Implement authenticated request
	c.mu.RLock()
	needsRefresh := !c.TokenExpiry.IsZero() && c.TokenExpiry.Before(time.Now())
	c.mu.RUnlock()

	if needsRefresh {
		if err := c.DoRefreshToken(); err != nil {
			return nil, fmt.Errorf("token expired and refresh failed: %w", err)
		}
	}
	c.mu.RLock()
	accessToken := c.AccessToken
	c.mu.RUnlock()
	if accessToken == "" {
		return nil, errors.New("no access token available; call ExchangeCodeForToken first")
	}
	req, err := http.NewRequest(method, targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}

func main() {
	server := NewOAuth2Server()
	// Register a client
	client := &OAuth2ClientInfo{
		ClientID:      "example-client",
		ClientSecret:  "example-secret",
		RedirectURIs:  []string{"http://localhost:8080/callback"},
		AllowedScopes: []string{"read", "write"},
	}
	if err := server.RegisterClient(client); err != nil {
		fmt.Printf("Error registering client: %v\n", err)
		return
	}
	fmt.Println("OAuth2 server is running on port 9000")
	// TODO The *http.Server is created locally and never returned or stored, so callers (and tests) cannot call srv.Shutdown(ctx).
	// For a challenge solution this is fine, but it limits testability.
	err := server.StartServer(9000)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
