// File: pkg/aiyou/ratelimit_test.go

package aiyou

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimiter_Basic(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         5,
		WaitTimeout:       time.Second,
	}

	limiter := NewRateLimiter(config, NewDefaultLogger(io.Discard))
	ctx := context.Background()

	t.Run("Burst Capacity", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < config.BurstSize; i++ {
			err := limiter.Wait(ctx)
			if err != nil {
				t.Errorf("Expected no error for request %d in burst, got %v", i, err)
			}
		}
		duration := time.Since(start)
		if duration > time.Millisecond*100 {
			t.Errorf("Burst requests took too long: %v", duration)
		}
	})

	t.Run("Rate Limiting", func(t *testing.T) {
		start := time.Now()
		err := limiter.Wait(ctx)
		duration := time.Since(start)

		expectedWait := time.Second / time.Duration(config.RequestsPerSecond)
		if duration < expectedWait {
			t.Errorf("Rate limiting not working, request completed too quickly: %v < %v",
				duration, expectedWait)
		}
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestRateLimiter_Concurrent(t *testing.T) {
	config := RateLimiterConfig{
		RequestsPerSecond: 5,
		BurstSize:         3,
		WaitTimeout:       time.Second,
	}

	limiter := NewRateLimiter(config, NewDefaultLogger(io.Discard))
	ctx := context.Background()

	t.Run("Concurrent Requests", func(t *testing.T) {
		var wg sync.WaitGroup
		requestCount := 10
		errorCount := int32(0)

		start := time.Now()

		for i := 0; i < requestCount; i++ {
			wg.Add(1)
			go func(reqNum int) {
				defer wg.Done()
				err := limiter.Wait(ctx)
				if err != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)

		expectedMinDuration := time.Duration(float64(requestCount-config.BurstSize) * float64(time.Second) / float64(config.RequestsPerSecond))
		if duration < expectedMinDuration {
			t.Errorf("Concurrent requests completed too quickly: %v < %v",
				duration, expectedMinDuration)
		}
	})
}

func TestRateLimiter_Integration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"test_token","expires_at":"2099-01-01T00:00:00Z"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	client, err := NewClient(
		WithEmailPassword("test@example.com", "password"),
		WithBaseURL(server.URL),
		WithRateLimiter(RateLimiterConfig{
			RequestsPerSecond: 2,
			BurstSize:         2,
			WaitTimeout:       time.Second,
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("Integration Test", func(t *testing.T) {
		ctx := context.Background()
		start := time.Now()

		for i := 0; i < 5; i++ {
			_, err := client.AuthenticatedRequest(ctx, "GET", "/test", nil)
			if err != nil {
				t.Errorf("Request %d failed: %v", i, err)
			}
		}

		duration := time.Since(start)
		// 1.5 secondes converties en Duration
		expectedMinDuration := time.Duration(1500 * time.Millisecond)
		if duration < expectedMinDuration {
			t.Errorf("Requests completed too quickly: %v < %v", duration, expectedMinDuration)
		}
	})
}
