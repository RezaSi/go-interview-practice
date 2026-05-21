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

// errors
var (
	ErrInvalidClient           = errors.New("invalid_client")
	ErrInvalidRequest          = errors.New("invalid_request")
	ErrInvalidScope            = errors.New("invalid_scope")
	ErrInvalidGrant            = errors.New("invalid_grant")
	ErrUnauthorizedClient      = errors.New("unauthorized_client")
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")
	ErrUnsupportedGrantType    = errors.New("unsupported_grant_type")
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
	if client == nil {
		return fmt.Errorf("Nil client info pointer %w", ErrInvalidClient)
	}

	if client.ClientID == "" {
		return fmt.Errorf("Bad client ID %w", ErrInvalidClient)
	}

	_, exists := s.clients[client.ClientID]
	if exists {
		return fmt.Errorf("Client allready registered %w", ErrInvalidClient)
	}

	if client.ClientSecret == "" {
		return fmt.Errorf("Invalid client secret %w", ErrInvalidClient)
	}

	if len(client.RedirectURIs) == 0 {
		return fmt.Errorf("Invalid redirect URIs %w", ErrInvalidClient)
	}

	if len(client.AllowedScopes) == 0 {
		return fmt.Errorf("Invalid scope %w", ErrInvalidScope)
	}
	verifiedClient := &OAuth2ClientInfo{
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		RedirectURIs:  append([]string{}, client.RedirectURIs...),
		AllowedScopes: append([]string{}, client.AllowedScopes...),
	}
	s.clients[client.ClientID] = verifiedClient

	return nil

}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("generate random string %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

type AuthRequestDTO struct {
	userID              string
	clientID            string
	redirectUri         string
	responseType        string
	scope               string
	state               string
	codeChallenge       string
	codeChallengeMethod string
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement authorization endpoint
	// 1. Validate request parameters (client_id, redirect_uri, response_type, scope, state)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	requestDTO := getAuthRequestParams(r)

	if err := s.validAuthRequestParams(requestDTO); err != nil {
		if errors.Is(err, ErrUnsupportedResponseType) {
			authErrorRedirectResponse(w, r, err.Error(), requestDTO)
			return
		}
		authErrorResponse(w, err.Error())
		return
	}

	// 2. Authenticate the user (for this challenge, could be a simple login form)
	// 3. Present a consent screen to the user
	// 4. Generate an authorization code and redirect to the client with the code
	authCode, err := GenerateAuthCode(requestDTO)
	if err != nil {
		authErrorResponse(w, err.Error())
		return
	}
	s.mu.Lock()
	s.authCodes[authCode.Code] = authCode
	s.mu.Unlock()

	authRedirectResponse(w, r, authCode, requestDTO.state)
}

func authErrorResponse(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}

func authErrorRedirectResponse(w http.ResponseWriter, r *http.Request, err string, dto AuthRequestDTO) {
	response := make(map[string]string)
	response["error"] = err
	response["state"] = dto.state
	authResponse(w, r, dto.redirectUri, response)

}
func authRedirectResponse(w http.ResponseWriter, r *http.Request, code *AuthorizationCode, state string) {
	response := make(map[string]string)
	response["code"] = code.Code
	response["state"] = state
	authResponse(w, r, code.RedirectURI, response)
}
func authResponse(w http.ResponseWriter, r *http.Request, redirectUri string, response map[string]string) {
	uri, _ := url.Parse(redirectUri)
	val := url.Values{}
	for k, v := range response {
		val.Add(k, v)
	}
	uri.RawQuery = val.Encode()
	http.Redirect(w, r, uri.String(), http.StatusFound)

}
func getAuthRequestParams(r *http.Request) AuthRequestDTO {
	query := r.URL.Query()

	var userID string
	if user := r.Context().Value("user_id"); user != nil {
		userID = user.(string)
		if userID == "" {
			userID = "user1"
		}
	}

	return AuthRequestDTO{
		userID:              userID,
		clientID:            query.Get("client_id"),
		redirectUri:         query.Get("redirect_uri"),
		responseType:        query.Get("response_type"),
		scope:               query.Get("scope"),
		state:               query.Get("state"),
		codeChallenge:       query.Get("code_challenge"),
		codeChallengeMethod: query.Get("code_challenge_method"),
	}

}
func (s *OAuth2Server) validAuthRequestParams(dto AuthRequestDTO) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[dto.userID]; !exists {
		return fmt.Errorf("user not found %w", ErrInvalidRequest)
	}
	client, exists := s.clients[dto.clientID]
	if !exists {
		return fmt.Errorf("client not found %w", ErrInvalidRequest)
	}
	if !slices.Contains(client.RedirectURIs, dto.redirectUri) {
		return fmt.Errorf("bad redirect uri %w", ErrInvalidRequest)
	}
	if dto.responseType != "code" {
		return fmt.Errorf("%w", ErrUnsupportedResponseType)
	}
	scopes := strings.Fields(dto.scope)
	for _, scope := range scopes {
		if !slices.Contains(client.AllowedScopes, scope) {
			return fmt.Errorf("invalid scope %w", ErrInvalidScope)
		}
	}
	if dto.codeChallenge == "" {
		return fmt.Errorf("code challenge is required %w", ErrInvalidRequest)
	}
	return nil
}

func GenerateAuthCode(dto AuthRequestDTO) (*AuthorizationCode, error) {
	randStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	return &AuthorizationCode{
		Code:                randStr,
		ClientID:            dto.clientID,
		UserID:              dto.userID,
		RedirectURI:         dto.redirectUri,
		Scopes:              strings.Fields(dto.scope),
		ExpiresAt:           time.Now().Add(5 * time.Minute),
		CodeChallenge:       dto.codeChallenge,
		CodeChallengeMethod: dto.codeChallengeMethod,
	}, nil
}

// =======================================
// TOKEN HANDLER
// =======================================

type TokenRequestDTO struct {
	grantType    string
	code         string
	redirectUri  string
	clientID     string
	clientSecret string
	codeVerifier string
	refreshToken string
}

type TokenErrorBody struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"description"`
}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {

	// TODO: Implement token endpoint
	// 1. Validate request parameters (grant_type, code, redirect_uri, client_id, client_secret)
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		tokenErrorResponse(w, "unsupported Content-Type", ErrInvalidRequest)
		return
	}

	if r.Method != http.MethodPost {
		tokenErrorResponse(w, "method not allowed %w", ErrInvalidRequest)
	}

	tokenRequest, err := GetTokenRequestParam(r)
	if err != nil {
		tokenErrorResponse(w, "get token params", err)
	}
	switch tokenRequest.grantType {
	case "authorization_code":
		s.handdleAccessRequest(w, tokenRequest)
	case "refresh_token":
		s.handleRefreshRequest(w, tokenRequest)
	default:
		tokenErrorResponse(w, "handle token %w", ErrUnsupportedGrantType)
	}
	// 2. Verify the authorization code
	// 3. For PKCE, verify the code_verifier
	// 4. Generate access and refresh tokens
	// 5. Return the tokens as a JSON response
}

func tokenErrorResponse(w http.ResponseWriter, msg string, err error) {
	//TODO implement error map to status code
	var status int
	switch {
	case errors.Is(err, ErrInvalidRequest):
		status = http.StatusBadRequest
	case errors.Is(err, ErrUnauthorizedClient):
		status = http.StatusUnauthorized
	case errors.Is(err, ErrInvalidClient):
		status = http.StatusUnauthorized
	default:
		status = http.StatusBadRequest
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             err.Error(),
		"error_description": msg,
	})

}

func tokenResponse(w http.ResponseWriter, token *Token, refreshToken *RefreshToken) {
	// TODO implement
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"access_token":  token.AccessToken,
		"token_type":    "Bearer",
		"expires_in":    int(time.Until(token.ExpiresAt).Seconds()),
		"refresh_token": refreshToken.RefreshToken,
	}

	if len(token.Scopes) > 0 {
		response["scope"] = strings.Join(token.Scopes, " ")
	}

	json.NewEncoder(w).Encode(response)
}

func GetTokenRequestParam(r *http.Request) (*TokenRequestDTO, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parse form %w", err)
	}

	return &TokenRequestDTO{
		grantType:    r.PostFormValue("grant_type"),
		code:         r.PostFormValue("code"),
		redirectUri:  r.PostFormValue("redirect_uri"),
		clientID:     r.PostFormValue("client_id"),
		clientSecret: r.PostFormValue("client_secret"),
		codeVerifier: r.PostFormValue("code_verifier"),
		refreshToken: r.PostFormValue("refresh_token"),
	}, nil

}
func (s *OAuth2Server) handdleAccessRequest(w http.ResponseWriter, dto *TokenRequestDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()

	code, exists := s.authCodes[dto.code]
	if !exists {
		tokenErrorResponse(w, "handl access request", ErrInvalidRequest)
		return
	}
	delete(s.authCodes, code.Code)

	client, exists := s.clients[dto.clientID]
	if !exists {
		tokenErrorResponse(w, "client validate", ErrUnauthorizedClient)
		return
	}

	if dto.redirectUri != code.RedirectURI {
		tokenErrorResponse(w, "redirect validate", ErrInvalidRequest)
		return
	}

	if dto.clientSecret != client.ClientSecret {
		tokenErrorResponse(w, "cliend secret validate", ErrInvalidClient)
		return
	}

	if !VerifyCodeChallenge(dto.codeVerifier, code.CodeChallenge, code.CodeChallengeMethod) {
		tokenErrorResponse(w, "PKCE failed", ErrInvalidGrant)
		return
	}

	accToken, refToken, err := generateTokens(code.ClientID, code.UserID, code.Scopes)
	if err != nil {
		tokenErrorResponse(w, "generate tokens", err)
		return
	}

	s.tokens[accToken.AccessToken] = accToken
	s.refreshTokens[refToken.RefreshToken] = refToken
	tokenResponse(w, accToken, refToken)

}

func (s *OAuth2Server) handleRefreshRequest(w http.ResponseWriter, dto *TokenRequestDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()

	refreshToken, ok := s.refreshTokens[dto.refreshToken]
	if !ok {
		tokenErrorResponse(w, "Bad refresh token", ErrInvalidGrant)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		tokenErrorResponse(w, "Expired refresh token", ErrInvalidGrant)
		return
	}

	if dto.clientSecret != s.clients[refreshToken.ClientID].ClientSecret {
		tokenErrorResponse(w, "Expired refresh token", ErrInvalidClient)
		return
	}

	accToken, refToken, err := generateTokens(refreshToken.ClientID, refreshToken.UserID, refreshToken.Scopes)
	if err != nil {
		tokenErrorResponse(w, "generate tokens", err)
		return
	}

	s.tokens[accToken.AccessToken] = accToken
	s.refreshTokens[refToken.RefreshToken] = refToken
	delete(s.refreshTokens, refreshToken.RefreshToken)
	tokenResponse(w, accToken, refToken)

}

func generateTokens(clientID, userID string, scopes []string) (*Token, *RefreshToken, error) {
	tokenStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, fmt.Errorf("generate token string %w", err)
	}
	reftokenStr, err := GenerateRandomString(32)
	if err != nil {
		return nil, nil, fmt.Errorf("generate token string %w", err)
	}

	return &Token{
			AccessToken: tokenStr,
			ClientID:    clientID,
			UserID:      userID,
			Scopes:      scopes,
			ExpiresAt:   time.Now().Add(time.Hour),
		},
		&RefreshToken{
			RefreshToken: reftokenStr,
			ClientID:     clientID,
			UserID:       userID,
			Scopes:       scopes,
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}, nil

}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	// TODO: Implement token validation
	s.mu.Lock()
	defer s.mu.Unlock()
	accToken, ok := s.tokens[token]
	if !ok {
		return nil, fmt.Errorf("token not found %w", ErrInvalidRequest)
	}

	if accToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("expired token %w", ErrInvalidGrant)
	}

	// client, ok := s.clients[accToken.ClientID]
	// if !ok {
	// 	return nil, fmt.Errorf("client not found %w", ErrInvalidClient)
	// }
	// _, ok = s.users[accToken.UserID]
	// if !ok {
	// 	return nil, fmt.Errorf("user not found %w", ErrInvalidRequest)
	// }
	// for _, scope := range accToken.Scopes {
	// 	if !slices.Contains(client.AllowedScopes, scope) {
	// 		return nil, fmt.Errorf("scope not allowed %w", ErrInvalidScope)
	// 	}
	// }
	return accToken, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	refToken, ok := s.refreshTokens[refreshToken]
	if !ok {
		return nil, nil, fmt.Errorf("Bad refresh token %w", ErrInvalidGrant)
	}
	return generateTokens(refToken.ClientID, refToken.UserID, refToken.Scopes)
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	// TODO: Implement token revocation
	s.mu.Lock()
	defer s.mu.Unlock()
	if isRefreshToken {
		if _, ok := s.refreshTokens[token]; !ok {
			return fmt.Errorf("invalid refresh token %w", ErrInvalidRequest)
		}
		delete(s.refreshTokens, token)
		return nil
	}
	if _, ok := s.tokens[token]; !ok {
		return fmt.Errorf("invalid token %w", ErrInvalidRequest)
	}
	delete(s.tokens, token)
	return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	// TODO: Implement PKCE verification
	switch method {
	case "S256":
		hash := sha256.Sum256([]byte(codeVerifier))
		challenge := base64.RawURLEncoding.EncodeToString(hash[:])
		if codeChallenge != challenge {
			return false
		}
		return true

	case "plain":
		if codeChallenge != codeVerifier {
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

// =================================
// Client code to demonstrate usage
// =================================

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
	parsedPath, err := url.Parse(c.Config.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("Auth URL: %w", err)
	}
	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("client_id", c.Config.ClientID)
	v.Add("redirect_uri", c.Config.RedirectURI)
	v.Add("scope", strings.Join(c.Config.Scopes, " "))
	v.Add("state", state)
	v.Add("codeChallenge", codeChallenge)
	v.Add("codeChallengeMethod", codeChallengeMethod)

	parsedPath.RawQuery = v.Encode()

	return parsedPath.String(), nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	// TODO: Implement token exchange
	params := url.Values{}
	params.Add("grant_type", "")
	params.Add("code", code)
	params.Add("redirect_uri", c.Config.RedirectURI)
	params.Add("client_id", c.Config.ClientID)
	params.Add("client_secret", c.Config.ClientSecret)
	params.Add("codeVerifier", codeVerifier)

	r, err := http.NewRequest(http.MethodPost, c.Config.AuthorizationEndpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("Exchange code for token request: %w", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	response, err := httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("http client error: %w", err)
	}
	type ResponseDTO struct {
		AccessToken  string        `json:"asserr_token"`
		TokenType    string        `json:"token_type"`
		ExpiresIn    time.Duration `json:"expires_in"`
		RefreshToken string        `json:"refresh_token"`
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server response %d", response.StatusCode)
	}

	var responseDTO ResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&responseDTO); err != nil {
		return fmt.Errorf("json decode request body: %w", err)
	}
	r.Body.Close()

	if responseDTO.TokenType != "Bearer" {
		return fmt.Errorf("Wrong token type")
	}

	c.AccessToken = responseDTO.AccessToken
	c.RefreshToken = responseDTO.RefreshToken
	c.TokenExpiry = time.Now().Add(responseDTO.ExpiresIn)

	return errors.New("not implemented")

}

// RefreshToken refreshes the access token using the refresh token
func (c *OAuth2Client) DoRefreshToken() error {
	// TODO: Implement token refresh
	return errors.New("not implemented")
}

// MakeAuthenticatedRequest makes a request with the access token
func (c *OAuth2Client) MakeAuthenticatedRequest(url string, method string) (*http.Response, error) {
	// TODO: Implement authenticated request
	return nil, errors.New("not implemented")
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
