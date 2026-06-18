package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	errFormatClientSecret = errors.New("client secret must be at least 8 characters long")
	errFormatRedirectURIs = errors.New("at least one redirect URI is required")
)

// OAuth2Config contains configuration for the OAuth2 server
type OAuth2Config struct {
	AuthorizationEndpoint string   // AuthorizationEndpoint is the endpoint for authorization requests
	TokenEndpoint         string   // TokenEndpoint is the endpoint for token requests
	ClientID              string   // ClientID is the OAuth2 client identifier
	ClientSecret          string   // ClientSecret is the secret for the client
	RedirectURI           string   // RedirectURI is the URI to redirect to after authorization
	Scopes                []string // Scopes is a list of requested scopes
}

// OAuth2Server implements an OAuth2 authorization server
type OAuth2Server struct {
	clients       map[string]*OAuth2ClientInfo  // clients stores registered OAuth2 clients
	authCodes     map[string]*AuthorizationCode // authCodes stores issued authorization codes
	tokens        map[string]*Token             // tokens stores issued access tokens
	refreshTokens map[string]*RefreshToken      // refreshTokens stores issued refresh tokens
	users         map[string]*User              // users stores user credentials for demonstration purposes
	mu            sync.RWMutex                  // mutex for concurrent access to data
}

// OAuth2ClientInfo represents a registered OAuth2 client
type OAuth2ClientInfo struct {
	ClientID      string   // ClientID is the unique identifier for the client
	ClientSecret  string   // ClientSecret is the secret for the client
	RedirectURIs  []string // RedirectURIs is a list of allowed redirect URIs
	AllowedScopes []string // AllowedScopes is a list of scopes the client can request
}

// User represents a user in the system
type User struct {
	ID       string // ID is the unique identifier for the user
	Username string // Username is the username for the user
	Password string // Password is the password for the user (in a real system, this would be hashed)
}

// AuthorizationCode represents an issued authorization code
type AuthorizationCode struct {
	Code                string    // Code is the authorization code string
	ClientID            string    // ClientID is the client that requested the code
	UserID              string    // UserID is the user that authorized the client
	RedirectURI         string    // RedirectURI is the URI to redirect to
	Scopes              []string  // Scopes is a list of authorized scopes
	ExpiresAt           time.Time // ExpiresAt is when the code expires
	CodeChallenge       string    // CodeChallenge is for PKCE
	CodeChallengeMethod string    // CodeChallengeMethod is for PKCE
}

// Token represents an issued access token
type Token struct {
	AccessToken string    // AccessToken is the token string
	ClientID    string    // ClientID is the client that owns the token
	UserID      string    // UserID is the user that authorized the token
	Scopes      []string  // Scopes is a list of authorized scopes
	ExpiresAt   time.Time // ExpiresAt is when the token expires
}

// RefreshToken represents an issued refresh token
type RefreshToken struct {
	RefreshToken string    // RefreshToken is the token string
	ClientID     string    // ClientID is the client that owns the token
	UserID       string    // UserID is the user that authorized the token
	Scopes       []string  // Scopes is a list of authorized scopes
	ExpiresAt    time.Time // ExpiresAt is when the token expires
}

// NewOAuth2Server creates a new OAuth2Server
// ===================================================================
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
		Password: "password",
	}

	return server
}

// RegisterClient registers a new OAuth2 client
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateClientInfo(client); err != nil {
		return err
	}

	if _, exists := s.clients[client.ClientID]; exists {
		return fmt.Errorf("client with ID '%s' already registered", client.ClientID)
	}

	// Create a copy to avoid external modifications
	registeredClient := &OAuth2ClientInfo{
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		RedirectURIs:  make([]string, len(client.RedirectURIs)),
		AllowedScopes: make([]string, len(client.AllowedScopes)),
	}

	copy(registeredClient.RedirectURIs, client.RedirectURIs)
	copy(registeredClient.AllowedScopes, client.AllowedScopes)

	s.clients[client.ClientID] = registeredClient

	return nil
}

func (s *OAuth2Server) GetClient(clientID string) (*OAuth2ClientInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[clientID]
	if !exists {
		return nil, errors.New("client not found")
	}

	return &OAuth2ClientInfo{
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		RedirectURIs:  append([]string{}, client.RedirectURIs...),
		AllowedScopes: append([]string{}, client.AllowedScopes...),
	}, nil
}

// validate
func validateClientInfo(client *OAuth2ClientInfo) error {

	if err := validClientID(client.ClientID); err != nil {
		return err
	}

	if len(client.ClientSecret) < 8 {
		return errFormatClientSecret
	}

	if len(client.RedirectURIs) == 0 {
		return errFormatRedirectURIs
	}

	for _, uri := range client.RedirectURIs {
		if err := validateRedirectURI(uri); err != nil {
			return fmt.Errorf("invalid redirect URI '%s': %w", uri, err)
		}
	}

	for _, scope := range client.AllowedScopes {
		if err := validateScope(scope); err != nil {
			return fmt.Errorf("invalid scope '%s': %w", scope, err)
		}
	}

	return nil
}

func validateResponseType(responseType string) error {
	if responseType == "" {
		return fmt.Errorf("response_type is required")
	}

	supportedTypes := map[string]bool{"code": true}

	if _, ok := supportedTypes[responseType]; !ok {
		return fmt.Errorf("unsupported_response_type")
	}

	return nil
}

func validClientID(clientID string) error {
	if clientID == "" {
		return fmt.Errorf("client_id is required")
	}

	if len(clientID) < 8 || len(clientID) > 32 {
		return fmt.Errorf("invalid client_id length")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, clientID)
	if !matched {
		return fmt.Errorf("client_id contains invalid characters")
	}

	return nil
}

func validateState(state string) error {
	if state == "" {
		return fmt.Errorf("state parameter missing - CSRF protection weakened")
	}

	if len(state) < 6 {
		return fmt.Errorf("state too short (min 8 chars)")
	}

	if len(state) > 1024 {
		return fmt.Errorf("state too long (max 1024 chars)")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._~-]+$`, state)
	if !matched {
		return fmt.Errorf("state contains invalid characters")
	}

	return nil
}

func validateRedirectURI(uri string) error {
	if uri == "" {
		return errors.New("URI cannot be empty")
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return errors.New("only http and https schemes are allowed")
	}

	if parsed.Host == "" {
		return errors.New("host is required")
	}

	if strings.Contains(parsed.Host, "@") {
		return errors.New("invalid host format")
	}

	return nil
}

func validateScope(scope string) error {
	if scope == "" {
		return errors.New("scope cannot be empty")
	}

	requestedScopes := strings.Fields(scope)

	scopeMap := make(map[string]bool)
	for _, requestedScope := range requestedScopes {
		if err := validateScopeFormat(requestedScope); err != nil {
			return fmt.Errorf("invalid scope format: %s", requestedScope)
		}

		if !isScopeSupported(requestedScope) {
			return fmt.Errorf("unsupported scope: %s", requestedScope)
		}

		if _, ok := scopeMap[requestedScope]; ok {
			return fmt.Errorf("duplicate scopes found")
		}

		scopeMap[requestedScope] = true
	}

	return nil
}

func validateScopeFormat(scope string) error {
	if scope == "" {
		return fmt.Errorf("scope cannot be empty")
	}

	if len(scope) > 100 {
		return fmt.Errorf("scope too long")
	}

	matched, _ := regexp.MatchString(`^[\x21\x23-\x5B\x5D-\x7E]+$`, scope)
	if !matched {
		return fmt.Errorf("scope contains invalid characters")
	}

	return nil
}

func isScopeSupported(scope string) bool {
	supportedScopes := map[string]bool{
		"read": true, "write": true, "email": true, "profile": true,
		"offline_access": true, "openid": true}

	if _, ok := supportedScopes[scope]; !ok {
		return false
	}

	return true
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}

	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charSet[num.Int64()]
	}

	return string(result), nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeErrorToken(w http.ResponseWriter, statusCode int, message, redirectURI string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             message,
		"error_description": fmt.Sprintf("%v", statusCode),
	})

	log.Println("writeErrorToken:", statusCode, message)
}

func writeErrorAuthorize(w http.ResponseWriter, statusCode int, message, redirectURI string) {
	parsedURL, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "Invalid redirect URI: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Add error parameters
	query := parsedURL.Query()
	query.Add("error_code", fmt.Sprintf("%v", statusCode))
	query.Add("error", message)
	parsedURL.RawQuery = query.Encode()

	w.Header().Set("Location", parsedURL.String())
	w.WriteHeader(statusCode)

	log.Println("writeErrorAuthorize:", statusCode, message)
}

// HandleAuthorize
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client_id := r.URL.Query().Get("client_id")
	redirect_uri := r.URL.Query().Get("redirect_uri")
	response_type := r.URL.Query().Get("response_type")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")
	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")

	// 1. Validate request parameters
	if err := validClientID(client_id); err != nil {
		writeErrorAuthorize(w, http.StatusBadRequest, err.Error(), redirect_uri)
		return
	}

	if err := validateRedirectURI(redirect_uri); err != nil {
		writeErrorAuthorize(w, http.StatusBadRequest, err.Error(), redirect_uri)
		return
	}

	if err := validateScope(scope); err != nil {
		writeErrorAuthorize(w, http.StatusBadRequest, err.Error(), redirect_uri)
		return
	}

	if err := validateState(state); err != nil {
		writeErrorAuthorize(w, http.StatusBadRequest, err.Error(), redirect_uri)
		return
	}

	if err := validateResponseType(response_type); err != nil {
		writeErrorAuthorize(w, http.StatusFound, err.Error(), redirect_uri)
		return
	}

	client, err := s.GetClient(client_id)
	if err != nil {
		writeErrorAuthorize(w, http.StatusBadRequest, err.Error(), redirect_uri)
		return
	}

	if !slices.Contains(client.RedirectURIs, redirect_uri) {
		writeErrorAuthorize(w, http.StatusBadRequest, "invalid_redirect_uri", redirect_uri)
		return
	}

	// 2. Check if user is already authenticated
	userID := s.getAuthenticatedUserID(r)
	if userID == "" {
		userID = "user1"
	}

	// 3. Generate authorization code
	authCode, err := s.generateAuthorizationCode(client_id, userID, redirect_uri, scope, codeChallenge, codeChallengeMethod)
	if err != nil {
		writeErrorAuthorize(w, http.StatusBadRequest, "failed to generate authorization code", redirect_uri)
		return
	}

	// 4. Redirect with code
	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirect_uri, authCode.Code, state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *OAuth2Server) getAuthenticatedUserID(r *http.Request) string {
	if user := r.Context().Value("user_id"); user != nil {
		if userID, ok := user.(string); ok {
			return userID
		}
	}

	if cookie, err := r.Cookie("user_id"); err == nil {
		return cookie.Value
	}

	return ""
}

func (s *OAuth2Server) generateAuthorizationCode(clientID, userID, redirectURI, scopes string, codeChallenge, challengeMethod string) (*AuthorizationCode, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	code, err := GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	authCode := &AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		Scopes:              strings.Fields(scopes),
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: challengeMethod,
	}

	s.authCodes[code] = authCode
	return authCode, nil
}

// HandleToken
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")
	codeVerifier := r.FormValue("code_verifier")
	refreshToken := r.FormValue("refresh_token")

	switch grantType {
	case "authorization_code":
		s.handleAuthorizationCodeGrant(w, code, redirectURI, clientID, clientSecret, codeVerifier)
	case "refresh_token":
		s.handleRefreshTokenGrant(w, refreshToken, clientID, clientSecret, redirectURI)
	default:
		writeErrorToken(w, http.StatusBadRequest, "Unsupported grant type", redirectURI)
	}
}

func (s *OAuth2Server) handleAuthorizationCodeGrant(w http.ResponseWriter, code, redirectURI, clientID, clientSecret, codeVerifier string) {
	// Validate client credentials
	if !s.validateClientCredentials(clientID, clientSecret) {
		writeErrorToken(w, http.StatusUnauthorized, "invalid_client", redirectURI)
		return
	}

	authCode, err := s.getAuthorizationCode(code)
	if err != nil {
		writeErrorToken(w, http.StatusBadRequest, "Invalid authorization code", redirectURI)
		return
	}

	if authCode.ClientID != clientID {
		writeErrorToken(w, http.StatusBadRequest, "Client mismatch", redirectURI)
		return
	}

	if authCode.RedirectURI != redirectURI {
		writeErrorToken(w, http.StatusBadRequest, "Redirect URI mismatch", redirectURI)
		return
	}

	// Verify PKCE if code challenge was provided
	if authCode.CodeChallenge != "" {
		if codeVerifier == "" {
			writeErrorToken(w, http.StatusBadRequest, "Code verifier required", redirectURI)
			return
		}
		if !VerifyCodeChallenge(codeVerifier, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
			writeErrorToken(w, http.StatusBadRequest, "invalid_grant", redirectURI)
			return
		}
	}

	accessToken, refreshToken, err := s.generateTokens(authCode.ClientID, authCode.UserID, authCode.Scopes)
	if err != nil {
		writeErrorToken(w, http.StatusBadRequest, "Failed to generate tokens", redirectURI)
		return
	}

	s.deleteAuthorizationCode(code)

	s.sendTokenResponse(w, accessToken, refreshToken)
}

func (s *OAuth2Server) handleRefreshTokenGrant(w http.ResponseWriter, refreshToken, clientID, clientSecret, redirectURI string) {
	if !s.validateClientCredentials(clientID, clientSecret) {
		writeErrorToken(w, http.StatusUnauthorized, "Invalid client credentials", redirectURI)
		return
	}

	refreshTokenData, err := s.getRefreshToken(refreshToken)
	if err != nil {
		writeErrorToken(w, http.StatusBadRequest, "Invalid refresh token", redirectURI)
		return
	}

	if refreshTokenData.ClientID != clientID {
		writeErrorToken(w, http.StatusBadRequest, "Client mismatch", redirectURI)
		return
	}

	accessToken, newRefreshToken, err := s.generateTokens(refreshTokenData.ClientID, refreshTokenData.UserID, refreshTokenData.Scopes)
	if err != nil {
		writeErrorToken(w, http.StatusBadRequest, "Failed to generate tokens", redirectURI)
		return
	}

	s.deleteRefreshToken(refreshToken)

	s.sendTokenResponse(w, accessToken, newRefreshToken)
}

func (s *OAuth2Server) validateClientCredentials(clientID, clientSecret string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[clientID]
	return exists && client.ClientSecret == clientSecret
}

func (s *OAuth2Server) getAuthorizationCode(code string) (*AuthorizationCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	authCode, exists := s.authCodes[code]
	if !exists {
		return nil, errors.New("code not found")
	}
	if time.Now().After(authCode.ExpiresAt) {
		return nil, errors.New("code expired")
	}
	return authCode, nil
}

func (s *OAuth2Server) getRefreshToken(token string) (*RefreshToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	refreshToken, exists := s.refreshTokens[token]
	if !exists {
		return nil, errors.New("refresh token not found")
	}
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}
	return refreshToken, nil
}

func (s *OAuth2Server) deleteAuthorizationCode(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.authCodes, code)
}

func (s *OAuth2Server) deleteRefreshToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.refreshTokens, token)
}

func (s *OAuth2Server) generateTokens(clientID, userID string, scopes []string) (*Token, *RefreshToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate access token
	accessTokenStr, err := GenerateRandomString(32)
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

	// Generate refresh token
	refreshTokenStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}

	refreshToken := &RefreshToken{
		RefreshToken: refreshTokenStr,
		ClientID:     clientID,
		UserID:       userID,
		Scopes:       scopes,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30), // 30 days
	}

	s.tokens[accessTokenStr] = accessToken
	s.refreshTokens[refreshTokenStr] = refreshToken

	return accessToken, refreshToken, nil
}

func (s *OAuth2Server) sendTokenResponse(w http.ResponseWriter, accessToken *Token, refreshToken *RefreshToken) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"access_token":  accessToken.AccessToken,
		"token_type":    "Bearer",
		"expires_in":    int(time.Until(accessToken.ExpiresAt).Seconds()),
		"refresh_token": refreshToken.RefreshToken,
	}

	if len(accessToken.Scopes) > 0 {
		response["scope"] = strings.Join(accessToken.Scopes, " ")
	}

	json.NewEncoder(w).Encode(response)
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tokens[token]
	if !ok {
		return nil, errors.New("non existent token")
	}
	if t.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expired token")
	}
	return t, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rt, ok := s.refreshTokens[refreshToken]
	if !ok || rt.ExpiresAt.Before(time.Now()) {
		return nil, nil, errors.New("invalid token")
	}

	accessTokenStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}
	refreshTokenStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, err
	}

	aToken := &Token{
		AccessToken: accessTokenStr,
		ClientID:    rt.ClientID,
		UserID:      rt.UserID,
		Scopes:      rt.Scopes,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	rToken := &RefreshToken{
		RefreshToken: refreshTokenStr,
		ClientID:     rt.ClientID,
		UserID:       rt.UserID,
		Scopes:       rt.Scopes,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
	}

	s.tokens[accessTokenStr] = aToken
	s.refreshTokens[refreshTokenStr] = rToken
	delete(s.refreshTokens, refreshToken)

	return aToken, rToken, nil
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
	switch method {
	case "S256":
		hash := sha256.Sum256([]byte(codeVerifier))
		expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
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
// ===================================================================

// OAuth2Client represents a client application using OAuth2
type OAuth2Client struct {
	Config       OAuth2Config // Config is the OAuth2 configuration
	AccessToken  string       // Token is the current access token
	RefreshToken string       // RefreshToken is the current refresh token
	TokenExpiry  time.Time    // TokenExpiry is when the access token expires
	HTTPClient   *http.Client // HTTPClient for making requests
}

// NewOAuth2Client creates a new OAuth2 client
func NewOAuth2Client(config OAuth2Config) *OAuth2Client {
	return &OAuth2Client{
		Config:     config,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetAuthorizationURL returns the URL to redirect the user for authorization
func (c *OAuth2Client) GetAuthorizationURL(state string, codeChallenge string, codeChallengeMethod string) (string, error) {
	authURL, err := url.Parse(c.Config.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("invalid auth URL: %w", err)
	}

	params := url.Values{}
	params.Add("client_id", c.Config.ClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", c.Config.RedirectURI)
	params.Add("state", state)

	if len(c.Config.Scopes) > 0 {
		params.Add("scope", strings.Join(c.Config.Scopes, " "))
	}

	// PKCE parameters
	if codeChallenge != "" {
		params.Add("code_challenge", codeChallenge)
		params.Add("code_challenge_method", codeChallengeMethod)
	}

	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.Config.RedirectURI)
	data.Set("client_id", c.Config.ClientID)
	data.Set("client_secret", c.Config.ClientSecret)

	// PKCE code verifier
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequest("POST", c.Config.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	c.RefreshToken = tokenResp.RefreshToken
	if tokenResp.ExpiresIn > 0 {
		c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return nil
}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	if c.RefreshToken == "" {
		return errors.New("no refresh token available")
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.RefreshToken)
	data.Set("client_id", c.Config.ClientID)
	data.Set("client_secret", c.Config.ClientSecret)

	req, err := http.NewRequest("POST", c.Config.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("refresh request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode refresh response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		c.RefreshToken = tokenResp.RefreshToken
	}
	if tokenResp.ExpiresIn > 0 {
		c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return nil
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(urlStr string, method string) (*http.Response, error) {
	if c.AccessToken == "" {
		return nil, errors.New("no access token available")
	}

	// Check if token needs refresh
	if !c.TokenExpiry.IsZero() && time.Until(c.TokenExpiry) < 30*time.Second {
		if err := c.DoRefreshToken(); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
	}

	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Accept", "application/json")

	return c.HTTPClient.Do(req)
}

func (c *OAuth2Client) IsTokenValid() bool {
	if c.AccessToken == "" {
		return false
	}
	if c.TokenExpiry.IsZero() {
		return true
	}
	return time.Until(c.TokenExpiry) > 30*time.Second
}

func GeneratePKCECodeVerifier() (string, error) {
	return GenerateRandomString(32)
}

func GeneratePKCECodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// ===================================================================
func main() {

	server := NewOAuth2Server()

	testClient := &OAuth2ClientInfo{
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret-123",
		RedirectURIs:  []string{"http://localhost:8085/callback"},
		AllowedScopes: []string{"read", "write"},
	}

	if err := server.RegisterClient(testClient); err != nil {
		log.Fatal("Failed to register client:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/authorize", server.HandleAuthorize)
	mux.HandleFunc("/token", server.HandleToken)

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		error := r.URL.Query().Get("error")

		if error != "" {
			fmt.Fprintf(w, "<h1>Authorization Error</h1>")
			fmt.Fprintf(w, "<p>Error: %s</p>", error)
			fmt.Fprintf(w, "<p>Error description: %s</p>", r.URL.Query().Get("error_description"))
			return
		}

		if code == "" {
			fmt.Fprintf(w, "<h1>No Authorization Code</h1>")
			fmt.Fprintf(w, "<p>No authorization code received in callback</p>")
			fmt.Fprintf(w, "<p>Query parameters: %+v</p>", r.URL.Query())
			return
		}

		fmt.Fprintf(w, "<h1>Authorization Successful!</h1>")
		fmt.Fprintf(w, "<p>Authorization code: <strong>%s</strong></p>", code)
		fmt.Fprintf(w, "<p>State: %s</p>", state)
		fmt.Fprintf(w, "<p>Copy this code and paste it in the terminal</p>")
		fmt.Fprintf(w, `<script>
			console.log('Authorization code:', '%s');
			console.log('State:', '%s');
		</script>`, code, state)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "<h1>OAuth2 Server is running!</h1>")
		fmt.Fprintf(w, "<p>Available endpoints:</p>")
		fmt.Fprintf(w, "<ul>")
		fmt.Fprintf(w, "<li><a href='/authorize?client_id=test-client-id&response_type=code&redirect_uri=http://localhost:8085/callback&scope=read write&state=test123'>/authorize</a> - Authorization endpoint</li>")
		fmt.Fprintf(w, "</ul>")
	})

	go func() {
		fmt.Printf("Starting OAuth2 server\n")
		fmt.Printf("Server URL: http://localhost:8085\n")
		if err := http.ListenAndServe(":8085", mux); err != nil {
			log.Fatal("Server failed:", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	clientConfig := OAuth2Config{
		AuthorizationEndpoint: "http://localhost:8085/authorize",
		TokenEndpoint:         "http://localhost:8085/token",
		ClientID:              "test-client-id",
		ClientSecret:          "test-client-secret-123",
		RedirectURI:           "http://localhost:8085/callback",
		Scopes:                []string{"read", "write"},
	}

	client := NewOAuth2Client(clientConfig)

	codeVerifier, err := GeneratePKCECodeVerifier()
	if err != nil {
		log.Fatal("Failed to generate code verifier:", err)
	}
	codeChallenge := GeneratePKCECodeChallenge(codeVerifier)

	state := "test-state-123"
	authURL, err := client.GetAuthorizationURL(state, codeChallenge, "S256")
	if err != nil {
		log.Fatal("Failed to get authorization URL:", err)
	}

	fmt.Println("===================================================================")
	fmt.Println("OAUTH2 - PORT 8085")
	fmt.Println("===================================================================")
	fmt.Printf("1. Open this URL in your browser:\n")
	fmt.Printf("   %s\n", authURL)
	fmt.Println("")

	fmt.Print("2. After authorization, copy the code from browser and paste it here: ")
	var authCode string
	fmt.Scanln(&authCode)

	authCode = strings.TrimSpace(authCode)

	fmt.Println("")
	fmt.Println("3. Exchanging authorization code for tokens...")

	if err := client.ExchangeCodeForToken(authCode, codeVerifier); err != nil {
		log.Fatal("Failed to exchange code for tokens:", err)
	}

	fmt.Println("✅ Token exchange successful!")
	fmt.Printf("Access Token: %s\n", client.AccessToken)
	fmt.Printf("Refresh Token: %s\n", client.RefreshToken)
	fmt.Printf("Token Expiry: %s\n", client.TokenExpiry.Format("2006-01-02 15:04:05"))
	fmt.Println("")

	fmt.Printf("4. Token valid: %v\n", client.IsTokenValid())
	fmt.Println("")

	fmt.Println("5. Testing authenticated request to protected endpoint...")

	resp, err := client.MakeAuthenticatedRequest("http://localhost:8085/protected", "GET")
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("✅ Request successful! Status: %d\n", resp.StatusCode)
		fmt.Printf("Response: %s\n", string(body))
	}
	fmt.Println("")

	fmt.Println("6. Testing token refresh...")
	if err := client.DoRefreshToken(); err != nil {
		fmt.Printf("❌ Refresh failed: %v\n", err)
	} else {
		fmt.Println("✅ Refresh successful!")
		fmt.Printf("New Access Token: %s\n", client.AccessToken)
		fmt.Printf("New Refresh Token: %s\n", client.RefreshToken)
		fmt.Printf("New Token Expiry: %s\n", client.TokenExpiry.Format("2006-01-02 15:04:05"))
	}

	fmt.Println("")
	fmt.Println("===================================================================")
	fmt.Println("OAUTH2 COMPLETED SUCCESSFULLY!")
	fmt.Println("===================================================================")
	log.Fatal("")

	select {}
}
