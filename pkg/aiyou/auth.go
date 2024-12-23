/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package aiyou provides authentication functionalities for the AI.YOU API.
package aiyou

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// JWTAuthenticator implements the Authenticator interface for JWT-based authentication
// using email and password credentials.
type JWTAuthenticator struct {
	email    string
	password string
	token    string
	expiry   time.Time
	client   *http.Client
	baseURL  string
	logger   Logger
}

// BearerAuthenticator implements the Authenticator interface for direct bearer token authentication
// without requiring email/password credentials.
type BearerAuthenticator struct {
	token  string
	logger Logger
}

// NewJWTAuthenticator creates a new instance of JWTAuthenticator for email/password authentication.
func NewJWTAuthenticator(email, password, baseURL string, client *http.Client, logger Logger) *JWTAuthenticator {
	if logger == nil {
		logger = NewDefaultLogger(io.Discard) // Default silent logger
	}
	return &JWTAuthenticator{
		email:    email,
		password: password,
		baseURL:  baseURL,
		client:   client,
		logger:   logger,
	}
}

// NewBearerAuthenticator creates a new instance of BearerAuthenticator for direct token authentication.
func NewBearerAuthenticator(token string, logger Logger) *BearerAuthenticator {
	if logger == nil {
		logger = NewDefaultLogger(io.Discard) // Default silent logger
	}
	logger.Debugf("Initializing Bearer authenticator with token: %s", maskSensitiveInfo(token))
	return &BearerAuthenticator{
		token:  token,
		logger: logger,
	}
}

// SetLogger sets a custom logger for the JWT authenticator
func (a *JWTAuthenticator) SetLogger(logger Logger) {
	a.logger = logger
}

// SetLogger sets a custom logger for the Bearer authenticator
func (a *BearerAuthenticator) SetLogger(logger Logger) {
	a.logger = logger
}

// Authenticate performs the authentication process and obtains a JWT token
// for email/password authentication.
func (a *JWTAuthenticator) Authenticate(ctx context.Context) error {
	if !a.tokenExpired() {
		a.logger.Debugf("JWT token is still valid, skipping authentication")
		return nil
	}

	a.logger.Debugf("Authenticating user: %s", maskSensitiveInfo(a.email))
	loginReq := LoginRequest{
		Email:    a.email,
		Password: a.password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		a.logger.Errorf("Failed to marshal login request: %v", err)
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/api/login", bytes.NewBuffer(jsonData))
	if err != nil {
		a.logger.Errorf("Failed to create login request: %v", err)
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	a.logger.Debugf("Sending login request")
	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Errorf("Failed to send login request: %v", err)
		return fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.logger.Warnf("Authentication failed with status code: %d", resp.StatusCode)
		return fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		a.logger.Errorf("Failed to decode login response: %v", err)
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	a.token = loginResp.Token
	a.expiry = loginResp.ExpiresAt

	a.logger.Debugf("Authentication successful, token expires at %v", a.expiry)
	return nil
}

// Authenticate for BearerAuthenticator validates the token existence
// and returns immediately as no API call is needed.
func (a *BearerAuthenticator) Authenticate(ctx context.Context) error {
	if a.token == "" {
		a.logger.Errorf("Bearer token authentication failed: token is empty")
		return &AuthenticationError{Message: "bearer token is empty"}
	}
	a.logger.Debugf("Using provided bearer token for authentication")
	return nil
}

// Token returns the current JWT token
func (a *JWTAuthenticator) Token() string {
	return a.token
}

// Token returns the bearer token
func (a *BearerAuthenticator) Token() string {
	return a.token
}

// SetToken updates the bearer token
func (a *BearerAuthenticator) SetToken(token string) {
	a.logger.Infof("Updating bearer token")
	a.token = token
	a.logger.Debugf("Bearer token has been successfully updated: %s", maskSensitiveInfo(token))
}

// tokenExpired checks if the current JWT token has expired
func (a *JWTAuthenticator) tokenExpired() bool {
	return a.token == "" || time.Now().After(a.expiry)
}

// maskSensitiveInfo masque les informations sensibles dans une chaîne
func maskSensitiveInfo(input string) string {
	// Masquer les adresses e-mail
	emailPattern := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	maskedInput := emailPattern.ReplaceAllString(input, "[EMAIL REDACTED]")

	// Masquer les tokens JWT
	tokenPattern := regexp.MustCompile(`Bearer\s+[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`)
	maskedInput = tokenPattern.ReplaceAllString(maskedInput, "Bearer [TOKEN REDACTED]")

	return maskedInput
}
