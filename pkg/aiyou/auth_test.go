package aiyou

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestJWTAuthentication(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	t.Run("Successful Authentication", func(t *testing.T) {
		auth := NewJWTAuthenticator(
			testConfig.Email,
			testConfig.Password,
			testConfig.BaseURL,
			&http.Client{Timeout: 10 * time.Second},
			logger,
		)

		err := auth.Authenticate(context.Background())
		if err != nil {
			t.Fatalf("Authentication failed: %v", err)
		}

		token := auth.Token()
		if token == "" {
			t.Error("Expected non-empty token after successful authentication")
		}

		t.Logf("Successfully authenticated with token length: %d", len(token))
	})

	t.Run("Failed Authentication", func(t *testing.T) {
		auth := NewJWTAuthenticator(
			"invalid@example.com",
			"wrong_password",
			testConfig.BaseURL,
			&http.Client{Timeout: 10 * time.Second},
			logger,
		)

		err := auth.Authenticate(context.Background())
		if err == nil {
			t.Error("Expected error with invalid credentials")
		} else {
			t.Logf("Got expected authentication error: %v", err)
		}
	})
}

func TestTokenExpiration(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	auth := NewJWTAuthenticator(
		testConfig.Email,
		testConfig.Password,
		testConfig.BaseURL,
		&http.Client{Timeout: 10 * time.Second},
		logger,
	)

	// Première authentification
	err := auth.Authenticate(context.Background())
	if err != nil {
		t.Fatalf("Initial authentication failed: %v", err)
	}

	initialToken := auth.Token()
	t.Logf("Initial token obtained, length: %d", len(initialToken))

	// Vérifier que le token est valide
	if auth.tokenExpired() {
		t.Error("Token should not be expired immediately after authentication")
	}

	// Forcer l'expiration
	auth.expiry = time.Now().Add(-1 * time.Hour)
	t.Log("Forced token expiration")

	// Vérifier que le token est maintenant expiré
	if !auth.tokenExpired() {
		t.Error("Token should be expired after setting past expiry time")
	}

	// Réauthentification après expiration
	err = auth.Authenticate(context.Background())
	if err != nil {
		t.Fatalf("Re-authentication failed: %v", err)
	}

	newToken := auth.Token()
	t.Logf("New token obtained after expiration, length: %d", len(newToken))

	if newToken == "" {
		t.Error("Expected non-empty token after re-authentication")
	}
}

func TestBearerAuthentication(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	t.Run("Valid Bearer Token", func(t *testing.T) {
		// D'abord obtenir un token valide via JWT
		jwtAuth := NewJWTAuthenticator(
			testConfig.Email,
			testConfig.Password,
			testConfig.BaseURL,
			&http.Client{Timeout: 10 * time.Second},
			logger,
		)

		err := jwtAuth.Authenticate(context.Background())
		if err != nil {
			t.Fatalf("Failed to get initial token: %v", err)
		}

		validToken := jwtAuth.Token()

		// Créer un auth Bearer avec le token valide
		bearerAuth := NewBearerAuthenticator(validToken, logger)

		// Tester l'authentification Bearer
		err = bearerAuth.Authenticate(context.Background())
		if err != nil {
			t.Errorf("Bearer authentication failed: %v", err)
		}

		if bearerAuth.Token() != validToken {
			t.Error("Bearer token mismatch")
		}

		t.Log("Successfully tested bearer authentication")
	})

	t.Run("Invalid Bearer Token", func(t *testing.T) {
		bearerAuth := NewBearerAuthenticator("invalid_token", logger)

		err := bearerAuth.Authenticate(context.Background())
		if err == nil {
			t.Error("Expected error with invalid bearer token")
		} else {
			t.Logf("Got expected error with invalid bearer token: %v", err)
		}
	})

	t.Run("Empty Bearer Token", func(t *testing.T) {
		bearerAuth := NewBearerAuthenticator("", logger)

		err := bearerAuth.Authenticate(context.Background())
		if err == nil {
			t.Error("Expected error with empty bearer token")
		} else {
			t.Logf("Got expected error with empty bearer token: %v", err)
		}
	})
}

func TestAuthenticatorInterface(t *testing.T) {
	// Vérifier que les deux types implémentent bien l'interface Authenticator
	var _ Authenticator = (*JWTAuthenticator)(nil)
	var _ Authenticator = (*BearerAuthenticator)(nil)

	t.Log("Verified that both authenticators implement the Authenticator interface")
}
