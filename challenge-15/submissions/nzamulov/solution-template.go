package main

import (
    "crypto/rand"
	"fmt"
	"net/http"
	"sync"
	"time"
	"strings"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

// OAuth2Config contains configuration for the OAuth2 server
type OAuth2Config struct {
	AuthorizationEndpoint string
	TokenEndpoint string
	ClientID string
	ClientSecret string
	RedirectURI string
	Scopes []string
}

// OAuth2Server implements an OAuth2 authorization server
type OAuth2Server struct {
	clients map[string]*OAuth2ClientInfo
	authCodes map[string]*AuthorizationCode
	tokens map[string]*Token
	refreshTokens map[string]*RefreshToken
	users map[string]*User
	mu sync.RWMutex
}

// OAuth2ClientInfo represents a registered OAuth2 client
type OAuth2ClientInfo struct {
	ClientID string
	ClientSecret string
	RedirectURIs []string
	AllowedScopes []string
}

// User represents a user in the system
type User struct {
	ID string
	Username string
	Password string
}

// AuthorizationCode represents an issued authorization code
type AuthorizationCode struct {
	Code string
	ClientID string
	UserID string
	RedirectURI string
	Scopes []string
	ExpiresAt time.Time
	CodeChallenge string
	CodeChallengeMethod string
}

// Token represents an issued access token
type Token struct {
	AccessToken string
	ClientID string
	UserID string
	Scopes []string
	ExpiresAt time.Time
}

// RefreshToken represents an issued refresh token
type RefreshToken struct {
	RefreshToken string
	ClientID string
	UserID string
	Scopes []string
	ExpiresAt time.Time
}

type TokenOrCode struct {
    ClientID string
    UserID string
    Scopes []string
    ExpiresAt time.Time
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type errorResponse struct {
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if client.ClientID == "" {
	    return fmt.Errorf("invalid client id: %s", client.ClientID)
	}
	if _, found := s.clients[client.ClientID]; found {
	    return fmt.Errorf("client with id %s already exists", client.ClientID)
	}
	s.clients[client.ClientID] = client
	return nil
}

// GenerateRandomString returns a URL‑safe random string of exact length using crypto/rand.
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("cannot generate random string with length %d", length)
	}
	// Over-generate bytes to ensure encoded output >= length, then trim.
	byteLen := length // good enough with base64 expansion
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(buf)
	if len(s) < length {
		// Extremely unlikely; top up once.
		extra := make([]byte, length)
		if _, err := rand.Read(extra); err != nil {
			return "", err
		}
		s += base64.RawURLEncoding.EncodeToString(extra)
	}
	return s[:length], nil
}

var authorizeParams = []string{"client_id", "redirect_uri", "response_type", "scope", "state"}

// HandleAuthorize handles the authorization endpoint
func (s *OAuth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	for _, param := range authorizeParams {
	    if val := q.Get(param); len(val) == 0 {
	        w.WriteHeader(http.StatusBadRequest)
	        return
	    }
	}
	
	clientID := q.Get("client_id")
	s.mu.RLock()
	client, found := s.clients[clientID]
	if !found {
	    s.mu.RUnlock()
	    w.WriteHeader(http.StatusBadRequest)
	    return
	}
	s.mu.RUnlock()

    redirectURI := q.Get("redirect_uri")
    allowedURI := false
    for _, allowedRedirectURI := range client.RedirectURIs {
        if redirectURI == allowedRedirectURI {
            allowedURI = true
            break
        }
    }
    if !allowedURI {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    
    scopes := strings.Fields(q.Get("scope"))
    anyNotAllowedScope := false
    for _, gotScope := range scopes {
        found := false
        for _, allowedScope := range client.AllowedScopes {
            found = found || gotScope == allowedScope
        }
        anyNotAllowedScope = anyNotAllowedScope || !found
    }
    if anyNotAllowedScope {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if q.Get("response_type") != "code" {
        w.Header().Set("Location", fmt.Sprintf("%s?error=%s&state=%s", redirectURI, "unsupported_response_type", q.Get("state")))
	    w.WriteHeader(http.StatusFound)
	    return
    }

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
	    w.WriteHeader(http.StatusBadRequest)
		return
	}
	
    codeStr, err := GenerateRandomString(32)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, "server_error", "failed to generate authorization code")
        return
    }

	code := &AuthorizationCode{
	    Code: codeStr,
	    ClientID: clientID,
	    UserID: userID,
	    RedirectURI: redirectURI,
	    Scopes: scopes,
	    ExpiresAt: time.Now().Add(10 * time.Minute),
	    CodeChallenge: q.Get("code_challenge"),
	    CodeChallengeMethod: q.Get("code_challenge_method"),
	}

    s.mu.Lock()
    s.authCodes[codeStr] = code
    s.mu.Unlock()

	w.Header().Set("Location", fmt.Sprintf("%s?code=%s&state=%s", redirectURI, codeStr, q.Get("state")))
	w.WriteHeader(http.StatusFound)
}

var accessTokenParams = []string{"grant_type", "code", "redirect_uri", "client_id", "client_secret"}
var refreshTokenParams = []string{"grant_type", "refresh_token", "client_id", "client_secret"}

const authorizationCodeGrantType = "authorization_code"
const refreshTokenGrantType = "refresh_token"

// HandleToken handles the token endpoint
func (s *OAuth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	grantType := r.FormValue("grant_type")
	if grantType != authorizationCodeGrantType && grantType != refreshTokenGrantType {
	    w.WriteHeader(http.StatusBadRequest)
	    return
	}
	
	var params = accessTokenParams
	if grantType == refreshTokenGrantType {
	    params = refreshTokenParams
	}
	
	for _, param := range params {
	    if val := r.FormValue(param); len(val) == 0 {
	        w.WriteHeader(http.StatusBadRequest)
	        return
	    }
	}
	
	clientID := r.FormValue("client_id")
	s.mu.Lock()
	defer s.mu.Unlock()

	client, found := s.clients[clientID]
	if !found {
	    w.WriteHeader(http.StatusBadRequest)
	    return
	}
	
	clientSecret := r.FormValue("client_secret")
	if client.ClientSecret != clientSecret {
	    WriteError(w, http.StatusUnauthorized, "invalid_client", "client secret is invalid")
		return
	}
	
	var tokenOrCode = &TokenOrCode{}

	if grantType == refreshTokenGrantType {
	    refreshTokenStr := r.FormValue("refresh_token")
    	refreshToken, found := s.refreshTokens[refreshTokenStr]
    	if !found {
    	    w.WriteHeader(http.StatusBadRequest)
    	    return
    	}

        tokenOrCode.ClientID = refreshToken.ClientID
	    tokenOrCode.UserID = refreshToken.UserID
	    tokenOrCode.Scopes = refreshToken.Scopes
	    tokenOrCode.ExpiresAt = refreshToken.ExpiresAt
	    
	    delete(s.refreshTokens, refreshTokenStr)
	}
	
	code := r.FormValue("code")
	if grantType == authorizationCodeGrantType {
    	authCode, found := s.authCodes[code]
    	if !found {
    	    w.WriteHeader(http.StatusBadRequest)
    	    return
    	}
    	
    	if authCode.Code != code {
    	    w.WriteHeader(http.StatusBadRequest)
    	    return
    	}

    	if r.FormValue("redirect_uri") != authCode.RedirectURI {
    	    WriteError(w, http.StatusBadRequest, "invalid_grant", "redirect_uri mismatch")
    	    return
    	}
    	
    	codeVerifierStr := r.FormValue("code_verifier")
        if authCode.CodeChallenge != "" {
            if codeVerifierStr == "" {
                WriteError(w, http.StatusBadRequest, "invalid_grant", "code_verifier is required for PKCE")
                return
            }
            if !VerifyCodeChallenge(codeVerifierStr, authCode.CodeChallenge, authCode.CodeChallengeMethod) {
                WriteError(w, http.StatusBadRequest, "invalid_grant", "code_verifier does not match code_challenge")
                return
            }
        }
        
        tokenOrCode.ClientID = authCode.ClientID
        tokenOrCode.UserID = authCode.UserID
    	tokenOrCode.Scopes = authCode.Scopes
    	tokenOrCode.ExpiresAt = authCode.ExpiresAt
    	
    	delete(s.authCodes, code)
	}
	
	if tokenOrCode.ClientID != clientID {
        WriteError(w, http.StatusUnauthorized, "invalid_client", "refresh token or code not issued to this client")
        return
    }

	if tokenOrCode.ExpiresAt.Before(time.Now()) {
	    w.WriteHeader(http.StatusUnauthorized)
	    return
	}
	
	accessToken, err := GenerateRandomString(32)
	if err != nil {
	    WriteError(w, http.StatusInternalServerError, "server_error", "failed to generate access token")
	    return
	}
	aToken := &Token{
	    AccessToken: accessToken,
	    ClientID: client.ClientID,
	    UserID: tokenOrCode.UserID,
	    Scopes: tokenOrCode.Scopes,
	    ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	
	refreshToken, err := GenerateRandomString(32)
	if err != nil {
	    WriteError(w, http.StatusInternalServerError, "server_error", "failed to generate refresh token")
	    return
	}
	rToken := &RefreshToken{
	    RefreshToken: refreshToken,
	    ClientID: client.ClientID,
	    UserID: tokenOrCode.UserID,
	    Scopes: tokenOrCode.Scopes,
	    ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.tokens[accessToken] = aToken
	s.refreshTokens[refreshToken] = rToken

	response := &tokenResponse{
	    AccessToken: accessToken,
	    TokenType: "Bearer",
	    ExpiresIn: int(time.Until(aToken.ExpiresAt).Seconds()),
	    RefreshToken: refreshToken,
	    Scope: strings.Join(tokenOrCode.Scopes, " "),
	} 
	
	jsonData, err := json.Marshal(response)
	if err != nil {
	    w.WriteHeader(http.StatusInternalServerError)
	    return
	}
	
	_ = WriteResponse(w, http.StatusOK, jsonData)
}

// ValidateToken validates an access token
func (s *OAuth2Server) ValidateToken(token string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	fetchedToken, found := s.tokens[token]
	if !found {
	    return nil, fmt.Errorf("token not found")
	}
	
	if fetchedToken.ExpiresAt.Before(time.Now()) {
	    return nil, fmt.Errorf("token has expired")
	}
	
	return fetchedToken, nil
}

// RevokeToken revokes an access or refresh token
func (s *OAuth2Server) RevokeToken(token string, isRefreshToken bool) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if isRefreshToken {
        if _, found := s.refreshTokens[token]; !found {
            return fmt.Errorf("refresh token %s not found", token)
        }
        
        delete(s.refreshTokens, token)
    } else {
        if _, found := s.tokens[token]; !found {
            return fmt.Errorf("access token %s not found", token)
        }
        
        delete(s.tokens, token)
    }
    
    return nil
}

// VerifyCodeChallenge verifies a PKCE code challenge
func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
    switch method {
        case "S256": {
            hashBytes := sha256.Sum256([]byte(codeVerifier))
            hash := base64.RawURLEncoding.EncodeToString(hashBytes[:])
            return codeChallenge == hash
        }
        case "plain": {
            return codeVerifier == codeChallenge
        }
        default: {
            return false
        }
    }
}

func WriteError(w http.ResponseWriter, statusCode int, errorStr, description string) {
    var errResp = errorResponse{
	    Error: errorStr,
	    Description: description,
	}
    
    jsonData, err := json.Marshal(errResp)
	if err != nil {
	    w.WriteHeader(http.StatusInternalServerError)
	    return
	}
	
    _ = WriteResponse(w, statusCode, jsonData)
}

func WriteResponse(w http.ResponseWriter, statusCode int, jsonData []byte) error {
    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
    w.WriteHeader(statusCode)
    _, err := w.Write(jsonData)
    return err
}