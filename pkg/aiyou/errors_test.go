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
	"errors"
	"testing"
)

func TestAPIError(t *testing.T) {
	err := &APIError{StatusCode: 400, Message: "Bad Request"}
	if err.Error() != "API error: 400 - Bad Request" {
		t.Errorf("APIError.Error() = %v, want %v", err.Error(), "API error: 400 - Bad Request")
	}
}

func TestAuthenticationError(t *testing.T) {
	err := &AuthenticationError{Message: "Invalid credentials"}
	if err.Error() != "Authentication error: Invalid credentials" {
		t.Errorf("AuthenticationError.Error() = %v, want %v", err.Error(), "Authentication error: Invalid credentials")
	}
}

func TestRateLimitError(t *testing.T) {
	err := &RateLimitError{
		RetryAfter:   30,
		IsClientSide: false,
	}
	expected := "server-side rate limit exceeded. Retry after 30 seconds"
	if err.Error() != expected {
		t.Errorf("RateLimitError.Error() = %v, want %v", err.Error(), expected)
	}

	// Test aussi le cas client-side
	errClient := &RateLimitError{
		RetryAfter:   30,
		IsClientSide: true,
	}
	expectedClient := "client-side rate limit exceeded. Retry after 30 seconds"
	if errClient.Error() != expectedClient {
		t.Errorf("RateLimitError.Error() = %v, want %v", errClient.Error(), expectedClient)
	}
}

func TestNetworkError(t *testing.T) {
	underlying := errors.New("connection reset")
	err := &NetworkError{Err: underlying}
	if err.Error() != "Network error: connection reset" {
		t.Errorf("NetworkError.Error() = %v, want %v", err.Error(), "Network error: connection reset")
	}
}
