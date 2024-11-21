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
package aiyou

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

func TestRetryOperation(t *testing.T) {
	ctx := context.Background()
	maxRetries := 3
	initialDelay := 10 * time.Millisecond
	logger := NewDefaultLogger(io.Discard) // Use a silent logger for tests

	t.Run("Successful operation", func(t *testing.T) {
		calls := 0
		err := retryOperation(ctx, logger, maxRetries, initialDelay, func() error {
			calls++
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if calls != 1 {
			t.Errorf("Expected 1 call, got %d", calls)
		}
	})

	t.Run("Retry on NetworkError", func(t *testing.T) {
		calls := 0
		err := retryOperation(ctx, logger, maxRetries, initialDelay, func() error {
			calls++
			if calls < 3 {
				return &NetworkError{Err: errors.New("connection failed")}
			}
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if calls != 3 {
			t.Errorf("Expected 3 calls, got %d", calls)
		}
	})

	t.Run("Max retries exceeded", func(t *testing.T) {
		calls := 0
		err := retryOperation(ctx, logger, maxRetries, initialDelay, func() error {
			calls++
			return &NetworkError{Err: errors.New("persistent error")}
		})

		if err == nil {
			t.Error("Expected an error, got nil")
		}
		if calls != maxRetries+1 {
			t.Errorf("Expected %d calls, got %d", maxRetries+1, calls)
		}
	})

	t.Run("Non-retryable error", func(t *testing.T) {
		calls := 0
		err := retryOperation(ctx, logger, maxRetries, initialDelay, func() error {
			calls++
			return errors.New("non-retryable error")
		})

		if err == nil {
			t.Error("Expected an error, got nil")
		}
		if calls != 1 {
			t.Errorf("Expected 1 call, got %d", calls)
		}
	})
}
