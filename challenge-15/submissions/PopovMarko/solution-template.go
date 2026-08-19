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

func oauthErrorCode(err error) string {
	switch {
	case errors.Is(err, ErrInvalidClient):
		return ErrInvalidClient.Error()
	case errors.Is(err, ErrInvalidRequest):
		return ErrInvalidRequest.Error()
	case errors.Is(err, ErrInvalidScope):
		return ErrInvalidScope.Error()
	case errors.Is(err, ErrInvalidGrant):
		return ErrInvalidGrant.Error()
	case errors.Is(err, ErrUnauthorizedClient):
		return ErrUnauthorizedClient.Error()
	case errors.Is(err, ErrUnsupportedResponseType):
		return ErrUnsupportedResponseType.Error()
	case errors.Is(err, ErrUnsupportedGrantType):
		return ErrUnsupportedGrantType.Error()
	default:
		return ErrInvalidRequest.Error()
	}
}

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

// RegisterClient registers a new OAuth2 client. The caller-supplied slices
// are copied so subsequent caller mutations cannot affect the stored value.
// Returns a wrapped ErrInvalidClient / ErrInvalidScope sentinel for any
// validation failure.
func (s *OAuth2Server) RegisterClient(client *OAuth2ClientInfo) error {
	if client == nil {
		return fmt.Errorf("nil client info: %w", ErrInvalidClient)
	}
	if client.ClientID == "" {
		return fmt.Errorf("empty client ID: %w", ErrInvalidClient)
	}
	if client.ClientSecret == "" {
		return fmt.Errorf("empty client secret: %w", ErrInvalidClient)
	}
	if len(client.RedirectURIs) == 0 {
		return fmt.Errorf("no redirect URIs: %w", ErrInvalidClient)
	}
	if len(client.AllowedScopes) == 0 {
		return fmt.Errorf("no allowed scopes: %w", ErrInvalidScope)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.clients[client.ClientID]; exists {
		return fmt.Errorf("client already registered: %w", ErrInvalidClient)
	}
	s.clients[client.ClientID] = &OAuth2ClientInfo{
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		RedirectURIs:  append([]string{}, client.RedirectURIs...),
		AllowedScopes: append([]string{}, client.AllowedScopes...),
	}
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

// HandleAuthorize handles GET requests to the OAuth2 authorization endpoint.
// On success it issues a short-lived authorization code and 302-redirects the
// user-agent back to the client's redirect_uri with code and state appended.
// On a recoverable error (e.g. unsupported response_type) it redirects with
// an OAuth2 error code; on a fatal error (unknown client / bad redirect URI)
// it returns 400 directly without redirecting to an untrusted target.
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	requestDTO := getAuthRequestParams(r)

	if err := s.validAuthRequestParams(requestDTO); err != nil {
		if errors.Is(err, ErrUnsupportedResponseType) {
			authErrorRedirectResponse(w, r, oauthErrorCode(err), requestDTO)
			return
		}
		authErrorResponse(w, err.Error())
		return
	}

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
	val := uri.Query()
	for k, v := range response {
		val.Set(k, v)
	}
	uri.RawQuery = val.Encode()
	http.Redirect(w, r, uri.String(), http.StatusFound)

}

// getAuthRequestParams reads OAuth2 authorization parameters from the request's
// URL query and pulls the logged-in user id off the request context. A missing
// context value falls back to "user1" so the demo flow works end-to-end without
// a real session middleware.
func getAuthRequestParams(r *http.Request) AuthRequestDTO {
	query := r.URL.Query()

	userID, _ := r.Context().Value("user_id").(string)
	if userID == "" {
		userID = "user1"
	}

	method := query.Get("code_challenge_method")
	if method == "" {
		method = "plain"
	}
	return AuthRequestDTO{
		userID:              userID,
		clientID:            query.Get("client_id"),
		redirectUri:         query.Get("redirect_uri"),
		responseType:        query.Get("response_type"),
		scope:               query.Get("scope"),
		state:               query.Get("state"),
		codeChallenge:       query.Get("code_challenge"),
		codeChallengeMethod: method,
	}
}

// validAuthRequestParams validates an authorize-endpoint DTO against the
// configured users and clients. Returns a wrapped sentinel error so callers
// can branch on errors.Is.
func (s *OAuth2Server) validAuthRequestParams(dto AuthRequestDTO) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if dto.userID == "" {
		return fmt.Errorf("user not authenticated: %w", ErrInvalidRequest)
	}
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
	if dto.codeChallengeMethod != "plain" && dto.codeChallengeMethod != "S256" {
		return fmt.Errorf("unsupported code challenge method: %w", ErrInvalidRequest)
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

// HandleToken handles POST requests to the OAuth2 token endpoint.
// It dispatches by grant_type to either the authorization_code or
// refresh_token flow. The request body must be application/x-www-form-urlencoded
// per RFC 6749 §3.2.
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		tokenErrorResponse(w, "method not allowed", ErrInvalidRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		tokenErrorResponse(w, "unsupported Content-Type", ErrInvalidRequest)
		return
	}

	tokenRequest, err := GetTokenRequestParam(r)
	if err != nil {
		tokenErrorResponse(w, "get token params", err)
		return
	}

	switch tokenRequest.grantType {
	case "authorization_code":
		s.handleAccessRequest(w, tokenRequest)
	case "refresh_token":
		s.handleRefreshRequest(w, tokenRequest)
	default:
		tokenErrorResponse(w, "handle token", ErrUnsupportedGrantType)
	}
}

// tokenErrorResponse writes an OAuth2-style JSON error body and maps the
// supplied error sentinel to an appropriate HTTP status code per RFC 6749 §5.2.
func tokenErrorResponse(w http.ResponseWriter, msg string, err error) {
	var status int
	switch {
	case errors.Is(err, ErrUnauthorizedClient),
		errors.Is(err, ErrInvalidClient):
		status = http.StatusUnauthorized
	default:
		status = http.StatusBadRequest
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":             oauthErrorCode(err),
		"error_description": msg,
	})
}

// tokenResponse writes a successful OAuth2 token response per RFC 6749 §5.1.
// The expires_in field is reported in seconds.
func tokenResponse(w http.ResponseWriter, token *Token, refreshToken *RefreshToken) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]any{
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

// GetTokenRequestParam parses the token request's
// application/x-www-form-urlencoded body and returns the parameters as a DTO.
// Per RFC 6749 §3.2 the token endpoint only reads credentials from the body
// — query-string values are ignored.
func GetTokenRequestParam(r *http.Request) (*TokenRequestDTO, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parse form: %w", err)
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

// handleAccessRequest implements the authorization_code grant: it exchanges
// a previously issued authorization code for an access + refresh token pair.
// The auth code is consumed only after all validation succeeds, so a failed
// exchange does not invalidate a usable code.
func (s *OAuth2Server) handleAccessRequest(w http.ResponseWriter, dto *TokenRequestDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()

	code, codeExists := s.authCodes[dto.code]

	client, clientExists := s.clients[dto.clientID]

	if !codeExists {
		tokenErrorResponse(w, "auth code not found", ErrInvalidGrant)
		return
	}
	if !clientExists {
		tokenErrorResponse(w, "client not registered", ErrUnauthorizedClient)
		return
	}
	if code.ExpiresAt.Before(time.Now()) {
		tokenErrorResponse(w, "auth code expired", ErrInvalidGrant)
		return
	}
	if code.ClientID != dto.clientID {
		tokenErrorResponse(w, "client mismatch", ErrUnauthorizedClient)
		return
	}
	if dto.redirectUri != code.RedirectURI {
		tokenErrorResponse(w, "redirect URI mismatch", ErrInvalidRequest)
		return
	}
	if dto.clientSecret != client.ClientSecret {
		tokenErrorResponse(w, "bad client secret", ErrInvalidClient)
		return
	}
	if !VerifyCodeChallenge(dto.codeVerifier, code.CodeChallenge, code.CodeChallengeMethod) {
		tokenErrorResponse(w, "PKCE verification failed", ErrInvalidGrant)
		return
	}

	delete(s.authCodes, code.Code)

	accToken, refToken, err := generateTokens(code.ClientID, code.UserID, code.Scopes)
	if err != nil {
		tokenErrorResponse(w, "generate tokens", err)
		return
	}

	s.tokens[accToken.AccessToken] = accToken
	s.refreshTokens[refToken.RefreshToken] = refToken

	tokenResponse(w, accToken, refToken)
}

// handleRefreshRequest implements the refresh_token grant: it rotates the
// supplied refresh token for a fresh access + refresh token pair. The old
// refresh token is revoked atomically with the issuance of the new one.
func (s *OAuth2Server) handleRefreshRequest(w http.ResponseWriter, dto *TokenRequestDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()

	refreshToken, tokenExists := s.refreshTokens[dto.refreshToken]

	var client *OAuth2ClientInfo

	if tokenExists {
		client = s.clients[refreshToken.ClientID]
	}

	if !tokenExists {
		tokenErrorResponse(w, "refresh token not found", ErrInvalidGrant)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		tokenErrorResponse(w, "refresh token expired", ErrInvalidGrant)
		return
	}
	if client == nil {
		tokenErrorResponse(w, "client no longer registered", ErrInvalidClient)
		return
	}
	if dto.clientID != client.ClientID {
		tokenErrorResponse(w, "client not found", ErrInvalidClient)
		return
	}
	if dto.clientSecret != client.ClientSecret {
		tokenErrorResponse(w, "bad client secret", ErrInvalidClient)
		return
	}

	accToken, refToken, err := generateTokens(refreshToken.ClientID, refreshToken.UserID, refreshToken.Scopes)
	if err != nil {
		tokenErrorResponse(w, "generate tokens", err)
		return
	}

	delete(s.refreshTokens, refreshToken.RefreshToken)
	s.tokens[accToken.AccessToken] = accToken
	s.refreshTokens[refToken.RefreshToken] = refToken

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

// ValidateToken looks up an access token and returns it if present and
// not expired. Returns ErrInvalidRequest when unknown and ErrInvalidGrant
// when expired.
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	accToken, ok := s.tokens[token]
	if !ok {
		s.mu.RUnlock()
		return nil, fmt.Errorf("token not found: %w", ErrInvalidRequest)
	}
	if accToken.ExpiresAt.Before(time.Now()) {
		s.mu.RUnlock()
		return nil, fmt.Errorf("token expired: %w", ErrInvalidGrant)
	}
	out := cloneToken(accToken)
	s.mu.RUnlock()
	return out, nil
}

// RefreshAccessToken rotates the supplied refresh token for a new
// access + refresh token pair. The previous refresh token is revoked.
// Returns ErrInvalidGrant if the token is unknown or expired.
func (s *OAuth2Server) RefreshAccessToken(refreshToken string) (*Token, *RefreshToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	refToken, ok := s.refreshTokens[refreshToken]
	if !ok {
		return nil, nil, fmt.Errorf("refresh token not found: %w", ErrInvalidGrant)
	}
	if refToken.ExpiresAt.Before(time.Now()) {
		return nil, nil, fmt.Errorf("refresh token expired: %w", ErrInvalidGrant)
	}

	accToken, newRefToken, err := generateTokens(refToken.ClientID, refToken.UserID, refToken.Scopes)
	if err != nil {
		return nil, nil, err
	}

	delete(s.refreshTokens, refreshToken)
	s.tokens[accToken.AccessToken] = accToken
	s.refreshTokens[newRefToken.RefreshToken] = newRefToken

	return cloneToken(accToken), cloneRefreshToken(newRefToken), nil
}

// RevokeToken revokes an access or refresh token. The isRefreshToken flag
// selects which store to consult. Returns ErrInvalidRequest if the token
// is unknown.
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if isRefreshToken {
		if _, ok := s.refreshTokens[token]; !ok {
			return fmt.Errorf("unknown refresh token: %w", ErrInvalidRequest)
		}
		delete(s.refreshTokens, token)
		return nil
	}
	if _, ok := s.tokens[token]; !ok {
		return fmt.Errorf("unknown access token: %w", ErrInvalidRequest)
	}
	delete(s.tokens, token)
	return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
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

func cloneToken(t *Token) *Token {
	if t == nil {
		return nil
	}
	cp := *t
	cp.Scopes = append([]string(nil), t.Scopes...)
	return &cp
}

func cloneRefreshToken(t *RefreshToken) *RefreshToken {
	if t == nil {
		return nil
	}
	cp := *t
	cp.Scopes = append([]string(nil), t.Scopes...)
	return &cp
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

// GetAuthorizationURL builds the authorization endpoint URL the user-agent
// should be redirected to. The state value is the caller's CSRF token; the
// server echoes it back unchanged. codeChallenge / codeChallengeMethod are
// the PKCE parameters as defined in RFC 7636.
func (c *OAuth2Client) GetAuthorizationURL(state string, codeChallenge string, codeChallengeMethod string) (string, error) {
	parsedPath, err := url.Parse(c.Config.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("parse authorization endpoint: %w", err)
	}
	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("client_id", c.Config.ClientID)
	v.Add("redirect_uri", c.Config.RedirectURI)
	v.Add("scope", strings.Join(c.Config.Scopes, " "))
	v.Add("state", state)
	v.Add("code_challenge", codeChallenge)
	v.Add("code_challenge_method", codeChallengeMethod)

	parsedPath.RawQuery = v.Encode()

	return parsedPath.String(), nil
}

// ExchangeCodeForToken exchanges a freshly obtained authorization code for
// an access + refresh token pair via the token endpoint. On success the
// client's AccessToken, RefreshToken, and TokenExpiry fields are updated.
func (c *OAuth2Client) ExchangeCodeForToken(code string, codeVerifier string) error {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", code)
	params.Add("redirect_uri", c.Config.RedirectURI)
	params.Add("client_id", c.Config.ClientID)
	params.Add("client_secret", c.Config.ClientSecret)
	params.Add("code_verifier", codeVerifier)

	r, err := http.NewRequest(http.MethodPost, c.Config.TokenEndpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("build token request: %w", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := c.doTokenRequest(r); err != nil {
		return fmt.Errorf("token exchange: %w", err)
	}

	return nil
}

// DoRefreshToken rotates the stored refresh token for a fresh access +
// refresh token pair. On success the client's AccessToken, RefreshToken,
// and TokenExpiry fields are updated in place.
func (c *OAuth2Client) DoRefreshToken() error {
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", c.RefreshToken)
	params.Add("client_id", c.Config.ClientID)
	params.Add("client_secret", c.Config.ClientSecret)

	r, err := http.NewRequest(http.MethodPost, c.Config.TokenEndpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("build refresh request: %w", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := c.doTokenRequest(r); err != nil {
		return fmt.Errorf("refresh token: %w", err)
	}

	return nil
}

// MakeAuthenticatedRequest issues an HTTP request to the given URL with the
// client's current access token attached as a Bearer credential.
func (c *OAuth2Client) MakeAuthenticatedRequest(url string, method string) (*http.Response, error) {
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build authenticated request: %w", err)
	}
	r.Header.Set("Authorization", "Bearer "+c.AccessToken)
	r.Header.Set("Accept", "application/json")
	return c.httpClient().Do(r)
}

// doTokenRequest sends a prepared token request, decodes the JSON response,
// and on success updates the client's AccessToken, RefreshToken, and
// TokenExpiry. ExpiresIn is interpreted as seconds per RFC 6749 §5.1.
func (c *OAuth2Client) doTokenRequest(r *http.Request) error {
	type responseDTO struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	var body responseDTO

	response, err := c.httpClient().Do(r)
	if err != nil {
		return fmt.Errorf("http client: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server response %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return fmt.Errorf("decode token response: %w", err)
	}

	if body.TokenType != "Bearer" {
		return fmt.Errorf("unexpected token_type %q", body.TokenType)
	}

	c.AccessToken = body.AccessToken
	c.RefreshToken = body.RefreshToken
	c.TokenExpiry = time.Now().Add(time.Duration(body.ExpiresIn) * time.Second)

	return nil
}

// httpClient returns a default HTTP client used for token and resource calls.
func (c *OAuth2Client) httpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
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
