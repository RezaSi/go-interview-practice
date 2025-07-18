package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestClientRegistration(t *testing.T) {
	server := NewOAuth2Server()

	client := &OAuth2ClientInfo{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		RedirectURIs:  []string{"https://client.example.com/callback"},
		AllowedScopes: []string{"read", "write", "profile"},
	}

	err := server.RegisterClient(client)
	if err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	// Try to register a client with the same ID
	duplicateClient := &OAuth2ClientInfo{
		ClientID:      "test-client",
		ClientSecret:  "different-secret",
		RedirectURIs:  []string{"https://different.example.com/callback"},
		AllowedScopes: []string{"read"},
	}

	err = server.RegisterClient(duplicateClient)
	if err == nil {
		t.Fatalf("Expected error when registering client with duplicate ID, got nil")
	}
}

func TestRandomStringGeneration(t *testing.T) {
	// Test generating random strings of different lengths
	lengths := []int{16, 32, 64, 128}

	for _, length := range lengths {
		s1, err := GenerateRandomString(length)
		if err != nil {
			t.Fatalf("Failed to generate random string: %v", err)
		}

		if len(s1) != length {
			t.Errorf("Expected random string of length %d, got %d", length, len(s1))
		}

		// Generate another string of the same length and ensure they're different
		s2, err := GenerateRandomString(length)
		if err != nil {
			t.Fatalf("Failed to generate second random string: %v", err)
		}

		if s1 == s2 {
			t.Errorf("Expected different random strings, got the same string: %s", s1)
		}
	}
}

func TestAuthorizationEndpoint(t *testing.T) {
	server := NewOAuth2Server()

	// Register a client
	client := &OAuth2ClientInfo{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		RedirectURIs:  []string{"https://client.example.com/callback"},
		AllowedScopes: []string{"read", "write", "profile"},
	}
	server.RegisterClient(client)

	t.Run("ValidRequest", func(t *testing.T) {
		// Create a request to the authorization endpoint
		req := httptest.NewRequest("GET", "/authorize?response_type=code&client_id=test-client&redirect_uri=https://client.example.com/callback&scope=read&state=xyz123&code_challenge=abc123&code_challenge_method=S256", nil)

		// Add a "logged in" user
		req = req.WithContext(createMockUserContext(req, "user1"))

		// Record the response
		w := httptest.NewRecorder()
		server.HandleAuthorize(w, req)

		// Check that we got a redirect
		resp := w.Result()
		if resp.StatusCode != http.StatusFound {
			t.Errorf("Expected status Found (302), got %v", resp.StatusCode)
		}

		// Check the location header
		location := resp.Header.Get("Location")
		if location == "" {
			t.Fatalf("Expected Location header, got none")
		}

		// Parse the redirect URL
		redirectURL, err := url.Parse(location)
		if err != nil {
			t.Fatalf("Failed to parse redirect URL: %v", err)
		}

		// Check that it's redirecting to the correct URI
		if !strings.HasPrefix(redirectURL.String(), "https://client.example.com/callback") {
			t.Errorf("Expected redirect to client callback URL, got %s", redirectURL.String())
		}

		// Check for authorization code in the query
		code := redirectURL.Query().Get("code")
		if code == "" {
			t.Errorf("Expected authorization code in redirect URL, got none")
		}

		// Check that the state parameter is preserved
		state := redirectURL.Query().Get("state")
		if state != "xyz123" {
			t.Errorf("Expected state 'xyz123', got '%s'", state)
		}
	})

	t.Run("InvalidClientID", func(t *testing.T) {
		// Create a request with an invalid client ID
		req := httptest.NewRequest("GET", "/authorize?response_type=code&client_id=invalid-client&redirect_uri=https://client.example.com/callback&scope=read&state=xyz123", nil)
		req = req.WithContext(createMockUserContext(req, "user1"))

		w := httptest.NewRecorder()
		server.HandleAuthorize(w, req)

		// Check that we get an error
		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest (400), got %v", resp.StatusCode)
		}
	})

	t.Run("InvalidRedirectURI", func(t *testing.T) {
		// Create a request with an invalid redirect URI
		req := httptest.NewRequest("GET", "/authorize?response_type=code&client_id=test-client&redirect_uri=https://attacker.example.com/callback&scope=read&state=xyz123", nil)
		req = req.WithContext(createMockUserContext(req, "user1"))

		w := httptest.NewRecorder()
		server.HandleAuthorize(w, req)

		// Check that we get an error
		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest (400), got %v", resp.StatusCode)
		}
	})

	t.Run("InvalidResponseType", func(t *testing.T) {
		// Create a request with an invalid response type
		req := httptest.NewRequest("GET", "/authorize?response_type=token&client_id=test-client&redirect_uri=https://client.example.com/callback&scope=read&state=xyz123", nil)
		req = req.WithContext(createMockUserContext(req, "user1"))

		w := httptest.NewRecorder()
		server.HandleAuthorize(w, req)

		// Check that we get a redirect with an error
		resp := w.Result()
		if resp.StatusCode != http.StatusFound {
			t.Errorf("Expected status Found (302), got %v", resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location == "" {
			t.Fatalf("Expected Location header, got none")
		}

		redirectURL, err := url.Parse(location)
		if err != nil {
			t.Fatalf("Failed to parse redirect URL: %v", err)
		}

		// Check for error in the query
		errorParam := redirectURL.Query().Get("error")
		if errorParam != "unsupported_response_type" {
			t.Errorf("Expected error 'unsupported_response_type', got '%s'", errorParam)
		}
	})
}

func TestTokenEndpoint(t *testing.T) {
	server := NewOAuth2Server()

	// Register a client
	client := &OAuth2ClientInfo{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		RedirectURIs:  []string{"https://client.example.com/callback"},
		AllowedScopes: []string{"read", "write", "profile"},
	}
	server.RegisterClient(client)

	// Create a valid authorization code
	codeStr, _ := GenerateRandomString(32)
	code := &AuthorizationCode{
		Code:                codeStr,
		ClientID:            "test-client",
		UserID:              "user1",
		RedirectURI:         "https://client.example.com/callback",
		Scopes:              []string{"read", "profile"},
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM", // SHA256 hash of "test-verifier"
		CodeChallengeMethod: "S256",
	}
	server.authCodes[codeStr] = code

	t.Run("ValidAuthorizationCode", func(t *testing.T) {
		// Create a token request
		form := url.Values{}
		form.Add("grant_type", "authorization_code")
		form.Add("code", codeStr)
		form.Add("redirect_uri", "https://client.example.com/callback")
		form.Add("client_id", "test-client")
		form.Add("client_secret", "test-secret")
		form.Add("code_verifier", "test-verifier")

		req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		server.HandleToken(w, req)

		// Check that we get a successful response
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK (200), got %v", resp.StatusCode)
		}

		// Check the response body
		var tokenResponse struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
			Scope        string `json:"scope"`
		}

		err := json.NewDecoder(resp.Body).Decode(&tokenResponse)
		if err != nil {
			t.Fatalf("Failed to decode token response: %v", err)
		}

		if tokenResponse.AccessToken == "" {
			t.Errorf("Expected access token, got none")
		}
		if tokenResponse.TokenType != "Bearer" {
			t.Errorf("Expected token type 'Bearer', got '%s'", tokenResponse.TokenType)
		}
		if tokenResponse.ExpiresIn <= 0 {
			t.Errorf("Expected positive expires_in value, got %d", tokenResponse.ExpiresIn)
		}
		if tokenResponse.RefreshToken == "" {
			t.Errorf("Expected refresh token, got none")
		}
		if tokenResponse.Scope != "read profile" {
			t.Errorf("Expected scope 'read profile', got '%s'", tokenResponse.Scope)
		}

		// Check that the code was consumed
		if _, exists := server.authCodes[codeStr]; exists {
			t.Errorf("Authorization code should be consumed after use")
		}
	})

	t.Run("InvalidClientCredentials", func(t *testing.T) {
		// Create a new authorization code
		newCodeStr, _ := GenerateRandomString(32)
		newCode := &AuthorizationCode{
			Code:        newCodeStr,
			ClientID:    "test-client",
			UserID:      "user1",
			RedirectURI: "https://client.example.com/callback",
			Scopes:      []string{"read"},
			ExpiresAt:   time.Now().Add(10 * time.Minute),
		}
		server.authCodes[newCodeStr] = newCode

		// Create a token request with invalid client credentials
		form := url.Values{}
		form.Add("grant_type", "authorization_code")
		form.Add("code", newCodeStr)
		form.Add("redirect_uri", "https://client.example.com/callback")
		form.Add("client_id", "test-client")
		form.Add("client_secret", "wrong-secret")

		req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		server.HandleToken(w, req)

		// Check that we get an error
		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized (401), got %v", resp.StatusCode)
		}

		// Check error response
		var errorResponse struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}

		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if errorResponse.Error != "invalid_client" {
			t.Errorf("Expected error 'invalid_client', got '%s'", errorResponse.Error)
		}
	})

	t.Run("InvalidCodeVerifier", func(t *testing.T) {
		// Create a new authorization code with PKCE
		newCodeStr, _ := GenerateRandomString(32)
		newCode := &AuthorizationCode{
			Code:                newCodeStr,
			ClientID:            "test-client",
			UserID:              "user1",
			RedirectURI:         "https://client.example.com/callback",
			Scopes:              []string{"read"},
			ExpiresAt:           time.Now().Add(10 * time.Minute),
			CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM", // SHA256 hash of "test-verifier"
			CodeChallengeMethod: "S256",
		}
		server.authCodes[newCodeStr] = newCode

		// Create a token request with invalid code_verifier
		form := url.Values{}
		form.Add("grant_type", "authorization_code")
		form.Add("code", newCodeStr)
		form.Add("redirect_uri", "https://client.example.com/callback")
		form.Add("client_id", "test-client")
		form.Add("client_secret", "test-secret")
		form.Add("code_verifier", "wrong-verifier")

		req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		server.HandleToken(w, req)

		// Check that we get an error
		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest (400), got %v", resp.StatusCode)
		}

		// Check error response
		var errorResponse struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}

		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if errorResponse.Error != "invalid_grant" {
			t.Errorf("Expected error 'invalid_grant', got '%s'", errorResponse.Error)
		}
	})
}

func TestRefreshToken(t *testing.T) {
	server := NewOAuth2Server()

	// Register a client
	client := &OAuth2ClientInfo{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		RedirectURIs:  []string{"https://client.example.com/callback"},
		AllowedScopes: []string{"read", "write", "profile"},
	}
	server.RegisterClient(client)

	// Create a valid refresh token
	refreshTokenStr, _ := GenerateRandomString(64)
	refreshToken := &RefreshToken{
		RefreshToken: refreshTokenStr,
		ClientID:     "test-client",
		UserID:       "user1",
		Scopes:       []string{"read", "profile"},
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	server.refreshTokens[refreshTokenStr] = refreshToken

	t.Run("ValidRefreshToken", func(t *testing.T) {
		// Create a refresh token request
		form := url.Values{}
		form.Add("grant_type", "refresh_token")
		form.Add("refresh_token", refreshTokenStr)
		form.Add("client_id", "test-client")
		form.Add("client_secret", "test-secret")

		req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		server.HandleToken(w, req)

		// Check that we get a successful response
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK (200), got %v", resp.StatusCode)
		}

		// Check the response body
		var tokenResponse struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
			Scope        string `json:"scope"`
		}

		err := json.NewDecoder(resp.Body).Decode(&tokenResponse)
		if err != nil {
			t.Fatalf("Failed to decode token response: %v", err)
		}

		if tokenResponse.AccessToken == "" {
			t.Errorf("Expected access token, got none")
		}
		if tokenResponse.TokenType != "Bearer" {
			t.Errorf("Expected token type 'Bearer', got '%s'", tokenResponse.TokenType)
		}
		if tokenResponse.ExpiresIn <= 0 {
			t.Errorf("Expected positive expires_in value, got %d", tokenResponse.ExpiresIn)
		}
		if tokenResponse.RefreshToken == "" {
			t.Errorf("Expected refresh token, got none")
		}
		if tokenResponse.RefreshToken == refreshTokenStr {
			t.Errorf("Expected new refresh token, got the same one")
		}
		if tokenResponse.Scope != "read profile" {
			t.Errorf("Expected scope 'read profile', got '%s'", tokenResponse.Scope)
		}

		// Check that the old refresh token was invalidated
		if _, exists := server.refreshTokens[refreshTokenStr]; exists {
			t.Errorf("Old refresh token should be invalidated")
		}
	})

	// Create another valid refresh token for the next test
	anotherRefreshTokenStr, _ := GenerateRandomString(64)
	anotherRefreshToken := &RefreshToken{
		RefreshToken: anotherRefreshTokenStr,
		ClientID:     "test-client",
		UserID:       "user1",
		Scopes:       []string{"read", "profile"},
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	server.refreshTokens[anotherRefreshTokenStr] = anotherRefreshToken

	t.Run("InvalidClientForRefreshToken", func(t *testing.T) {
		// Create a refresh token request with wrong client credentials
		form := url.Values{}
		form.Add("grant_type", "refresh_token")
		form.Add("refresh_token", anotherRefreshTokenStr)
		form.Add("client_id", "test-client")
		form.Add("client_secret", "wrong-secret")

		req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		server.HandleToken(w, req)

		// Check that we get an error
		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized (401), got %v", resp.StatusCode)
		}

		// Check that the refresh token was not invalidated
		if _, exists := server.refreshTokens[anotherRefreshTokenStr]; !exists {
			t.Errorf("Refresh token should not be invalidated on client authentication failure")
		}
	})
}

func TestTokenValidation(t *testing.T) {
	server := NewOAuth2Server()

	// Create a valid access token
	tokenStr, _ := GenerateRandomString(32)
	token := &Token{
		AccessToken: tokenStr,
		ClientID:    "test-client",
		UserID:      "user1",
		Scopes:      []string{"read", "profile"},
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	server.tokens[tokenStr] = token

	// Create an expired token
	expiredTokenStr, _ := GenerateRandomString(32)
	expiredToken := &Token{
		AccessToken: expiredTokenStr,
		ClientID:    "test-client",
		UserID:      "user1",
		Scopes:      []string{"read"},
		ExpiresAt:   time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}
	server.tokens[expiredTokenStr] = expiredToken

	t.Run("ValidToken", func(t *testing.T) {
		// Validate the token
		validatedToken, err := server.ValidateToken(tokenStr)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if validatedToken == nil {
			t.Fatalf("Expected token, got nil")
		}

		if validatedToken.AccessToken != tokenStr {
			t.Errorf("Expected token string %s, got %s", tokenStr, validatedToken.AccessToken)
		}
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		// Validate an expired token
		_, err := server.ValidateToken(expiredTokenStr)
		if err == nil {
			t.Fatalf("Expected error for expired token, got nil")
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		// Validate a non-existent token
		_, err := server.ValidateToken("non-existent-token")
		if err == nil {
			t.Fatalf("Expected error for invalid token, got nil")
		}
	})
}

func TestTokenRevocation(t *testing.T) {
	server := NewOAuth2Server()

	// Create a valid access token
	tokenStr, _ := GenerateRandomString(32)
	token := &Token{
		AccessToken: tokenStr,
		ClientID:    "test-client",
		UserID:      "user1",
		Scopes:      []string{"read", "profile"},
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	server.tokens[tokenStr] = token

	// Create a valid refresh token
	refreshTokenStr, _ := GenerateRandomString(64)
	refreshToken := &RefreshToken{
		RefreshToken: refreshTokenStr,
		ClientID:     "test-client",
		UserID:       "user1",
		Scopes:       []string{"read", "profile"},
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	server.refreshTokens[refreshTokenStr] = refreshToken

	t.Run("RevokeAccessToken", func(t *testing.T) {
		// Revoke the access token
		err := server.RevokeToken(tokenStr, false)
		if err != nil {
			t.Fatalf("Failed to revoke access token: %v", err)
		}

		// Check that the token was revoked
		if _, exists := server.tokens[tokenStr]; exists {
			t.Errorf("Access token should be revoked")
		}
	})

	t.Run("RevokeRefreshToken", func(t *testing.T) {
		// Revoke the refresh token
		err := server.RevokeToken(refreshTokenStr, true)
		if err != nil {
			t.Fatalf("Failed to revoke refresh token: %v", err)
		}

		// Check that the token was revoked
		if _, exists := server.refreshTokens[refreshTokenStr]; exists {
			t.Errorf("Refresh token should be revoked")
		}
	})

	t.Run("RevokeNonExistentToken", func(t *testing.T) {
		// Try to revoke a non-existent token
		err := server.RevokeToken("non-existent-token", false)
		if err == nil {
			t.Fatalf("Expected error when revoking non-existent token, got nil")
		}
	})
}

func TestPKCE(t *testing.T) {
	t.Run("VerifyValidCodeChallenge", func(t *testing.T) {
		// Test verifying a valid code challenge with S256 method
		codeVerifier := "test-verifier"
		codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM" // SHA256 hash of "test-verifier" in base64url encoding

		valid := VerifyCodeChallenge(codeVerifier, codeChallenge, "S256")
		if !valid {
			t.Errorf("Expected valid code challenge verification")
		}
	})

	t.Run("VerifyInvalidCodeChallenge", func(t *testing.T) {
		// Test verifying an invalid code challenge
		codeVerifier := "test-verifier"
		codeChallenge := "invalid-challenge"

		valid := VerifyCodeChallenge(codeVerifier, codeChallenge, "S256")
		if valid {
			t.Errorf("Expected invalid code challenge verification")
		}
	})

	t.Run("UnsupportedMethod", func(t *testing.T) {
		// Test verifying with an unsupported method
		codeVerifier := "test-verifier"
		codeChallenge := "test-verifier" // plain method would just compare directly

		valid := VerifyCodeChallenge(codeVerifier, codeChallenge, "unsupported")
		if valid {
			t.Errorf("Expected failure for unsupported method")
		}
	})

	t.Run("PlainMethod", func(t *testing.T) {
		// Test verifying with plain method
		codeVerifier := "test-verifier"
		codeChallenge := "test-verifier" // plain method just compares directly

		valid := VerifyCodeChallenge(codeVerifier, codeChallenge, "plain")
		if !valid {
			t.Errorf("Expected valid code challenge verification with plain method")
		}
	})
}

// Helper function to create a mock user context for testing
func createMockUserContext(r *http.Request, userID string) context.Context {
	// In a real application, you would have middleware that sets the user in the context
	// For testing, we'll just create a mock context with a user
	return context.WithValue(r.Context(), "user_id", userID)
}
