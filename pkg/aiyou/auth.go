/*
Copyright (C) 2023 Cloud Temple

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
    "net/http"
    "time"
)

// JWTAuthenticator implements the Authenticator interface for JWT-based authentication.
type JWTAuthenticator struct {
    email    string
    password string
    token    string
    expiry   time.Time
    client   *http.Client
    baseURL  string
}

// NewJWTAuthenticator creates a new instance of JWTAuthenticator.
func NewJWTAuthenticator(email, password, baseURL string, client *http.Client) *JWTAuthenticator {
    return &JWTAuthenticator{
        email:    email,
        password: password,
        baseURL:  baseURL,
        client:   client,
    }
}

// Authenticate performs the authentication process and obtains a JWT token.
func (a *JWTAuthenticator) Authenticate(ctx context.Context) error {
    if !a.tokenExpired() {
        return nil
    }

    loginReq := LoginRequest{
        Email:    a.email,
        Password: a.password,
    }

    jsonData, err := json.Marshal(loginReq)
    if err != nil {
        return fmt.Errorf("failed to marshal login request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/api/login", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create login request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := a.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send login request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
    }

    var loginResp LoginResponse
    if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
        return fmt.Errorf("failed to decode login response: %w", err)
    }

    a.token = loginResp.Token
    a.expiry = loginResp.ExpiresAt

    return nil
}

// Token returns the current JWT token.
func (a *JWTAuthenticator) Token() string {
    return a.token
}

// tokenExpired checks if the current token has expired.
func (a *JWTAuthenticator) tokenExpired() bool {
    return a.token == "" || time.Now().After(a.expiry)
}
