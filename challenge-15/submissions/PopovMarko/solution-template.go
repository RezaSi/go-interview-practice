package main

import (
	"crypto/rand"
	"encoding/hex"
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
	ErrInvalidRequest       = errors.New("invalid request")
	ErrInvaldClient         = errors.New("invalid client")
	ErrInvalidGrant         = errors.New("invalid or extired grant")
	ErrUnauthorizedClient   = errors.New("client not authorise for grant type")
	ErrUnsupportedGrantType = errors.New("unsupported grant type")
	ErrInvalidScope         = errors.New("invalid scope")
	ErrInvaldResponseType   = errors.New("invalid response type")
	ErrInternalServerError  = errors.New("internal server error")
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
	// State of the request
	State string
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
	if client.ClientID == "" {
		return fmt.Errorf("Client ID can't be empty: %w", ErrInvaldClient)
	}
	if _, exists := s.clients[client.ClientID]; exists {
		return fmt.Errorf("Duplicate client ID: %w", ErrInvaldClient)
	}
	if client.ClientSecret == "" {
		return fmt.Errorf("Client secret can't be empty: %w", ErrInvaldClient)
	}
	if len(client.RedirectURIs) == 0 {
		return fmt.Errorf("RedixectURIs can't be empty: %w", ErrInvaldClient)
	}
	if len(client.AllowedScopes) == 0 {
		return fmt.Errorf("AllowedScopes can't be empty: %w", ErrInvaldClient)
	}
	s.clients[client.ClientID] = client
	return nil
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	// TODO: Implement secure random string generation
	bytes := make([]byte, length)
	rand.Read(bytes) // Fills the slice with secure random bytes
	return hex.EncodeToString(bytes)[:length], nil
}

type AuthRequestDTO struct {
	ResponseType        string
	ClientID            string
	UserID              string
	RedirectURI         string
	Scopes              []string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement authorization endpoint
	// 1. Validate request parameters (client_id, redirect_uri, response_type, scope, state)
	authRequest, err := s.GetAndValidateAuthRequest(r)
	if err != nil {
		if errors.Is(err, ErrInvaldResponseType) {
			InvalidRTResponse(w, r, authRequest)
		}
		ErrorResponse(w, "Query param validation", err, http.StatusBadRequest)
		return
	}
	// 2. Authenticate the user (for this challenge, could be a simple login form)
	// 3. Present a consent screen to the user
	// 4. Generate an authorization code and redirect to the client with the code
	authCode, err := GenerateAuthCode(authRequest)
	if err != nil {
		ErrorResponse(w, "Auth code generation", err, http.StatusInternalServerError)
		return
	}

	AuthResponse(w, r, authCode)

}

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement token endpoint
	// 1. Validate request parameters (grant_type, code, redirect_uri, client_id, client_secret)
	// 2. Verify the authorization code
	// 3. For PKCE, verify the code_verifier
	// 4. Generate access and refresh tokens
	// 5. Return the tokens as a JSON response
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	// TODO: Implement token validation
	return nil, errors.New("not implemented")
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	// TODO: Implement token refresh
	return nil, nil, errors.New("not implemented")
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	// TODO: Implement token revocation
	return errors.New("not implemented")
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	// TODO: Implement PKCE verification
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
	// TODO: Implement building the authorization URL
	return "", errors.New("not implemented")
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	// TODO: Implement token exchange
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

func (s *OAuth2Server) GetAndValidateAuthRequest(r *http.Request) (AuthRequestDTO, error) {
	if r.Method != http.MethodGet {
		return AuthRequestDTO{}, fmt.Errorf("Method not allowed")
	}
	user := r.Context().Value("user_id")
	userID, ok := user.(string)
	if !ok {
		return AuthRequestDTO{}, fmt.Errorf("Bad request wrong user id: %v: %w", userID, ErrInvalidRequest)
	}
	query := r.URL.Query()

	clientID := query.Get("client_id")
	client, ok := s.clients[clientID]
	if !ok {
		return AuthRequestDTO{}, fmt.Errorf("client validation: %w", ErrInvaldClient)
	}

	redirectUri := query.Get("redirect_uri")
	if redirectUri == "" || !slices.Contains(client.RedirectURIs, redirectUri) {
		return AuthRequestDTO{}, fmt.Errorf("invalid redirect URI: %w", ErrInvalidRequest)
	}

	scope := query.Get("scope")
	if scope == "" {
		return AuthRequestDTO{}, fmt.Errorf("scope validation: %w", ErrInvalidScope)
	}
	scopes := strings.Split(scope, " ")
	for _, scope := range scopes {
		if slices.Contains(client.AllowedScopes, scope) {
			continue
		}
		return AuthRequestDTO{}, fmt.Errorf("scope not alowed: %w", ErrInvalidScope)
	}

	state := query.Get("state")
	if state == "" {
		return AuthRequestDTO{}, fmt.Errorf("invalid state %s: %w", state, ErrInvalidRequest)
	}

	responseType := query.Get("response_type")
	if responseType != "code" {
		return AuthRequestDTO{
				ResponseType:        responseType,
				ClientID:            clientID,
				UserID:              userID,
				RedirectURI:         redirectUri,
				Scopes:              scopes,
				State:               state,
				CodeChallenge:       query.Get("code_challenge"),
				CodeChallengeMethod: query.Get("code_challenge_method"),
			},
			fmt.Errorf("Bad request, wrong response type: %s: %w", responseType, ErrInvaldResponseType)
	}

	return AuthRequestDTO{
		ResponseType:        responseType,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectUri,
		Scopes:              scopes,
		State:               state,
		CodeChallenge:       query.Get("code_challenge"),
		CodeChallengeMethod: query.Get("code_challenge_method"),
	}, nil
}

func GenerateAuthCode(dto AuthRequestDTO) (AuthorizationCode, error) {
	codeStr, err := GenerateRandomString(32)
	if err != nil {
		return AuthorizationCode{}, err
	}
	return AuthorizationCode{
		Code:                codeStr,
		ClientID:            dto.ClientID,
		UserID:              dto.UserID,
		RedirectURI:         dto.RedirectURI,
		Scopes:              dto.Scopes,
		ExpiresAt:           time.Now().Add(5 * time.Minute),
		State:               dto.State,
		CodeChallenge:       dto.CodeChallenge,
		CodeChallengeMethod: dto.CodeChallengeMethod,
	}, nil
}

type ErrorBody struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func ErrorResponse(w http.ResponseWriter, msg string, err error, status int) {
	response := ErrorBody{
		Error:            err.Error(),
		ErrorDescription: msg,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&response)
}

func AuthResponse(w http.ResponseWriter, r *http.Request, authCode AuthorizationCode) {
	params := make(map[string]string)
	params["code"] = authCode.Code
	params["state"] = authCode.State
	redirectRespons(w, r, authCode.RedirectURI, params)
}

func InvalidRTResponse(w http.ResponseWriter, r *http.Request, authRequest AuthRequestDTO) {
	params := make(map[string]string)
	params["error"] = "unsupported_response_type"
	params["state"] = authRequest.State
	redirectRespons(w, r, authRequest.RedirectURI, params)
}

func redirectRespons(w http.ResponseWriter, r *http.Request, uri string, params map[string]string) {
	u, err := url.Parse(uri)
	if err != nil {
		ErrorResponse(w, "parsing redirect URI", err, http.StatusBadRequest)
	}
	v := u.Query()
	for key, val := range params {
		v.Set(key, val)
	}
	v.Set("error", "unsupported_response_type")
	u.RawQuery = v.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}
