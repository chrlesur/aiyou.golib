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

import "fmt"

// APIError représente une erreur retournée par l'API AI.YOU.
// Il contient le code de statut HTTP et le message d'erreur.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d - %s", e.StatusCode, e.Message)
}

// AuthenticationError représente une erreur d'authentification
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error: %s", e.Message)
}

// RateLimitError indique que la limite de taux a été atteinte.
// RetryAfter indique le nombre de secondes à attendre avant de réessayer.
type RateLimitError struct {
	RetryAfter   int  // en secondes
	IsClientSide bool // pour distinguer entre le rate limiting côté client et serveur
}

func (e *RateLimitError) Error() string {
	source := "server"
	if e.IsClientSide {
		source = "client"
	}
	return fmt.Sprintf("%s-side rate limit exceeded. Retry after %d seconds", source, e.RetryAfter)
}

// NetworkError représente une erreur de réseau survenue lors d'une requête.
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network error: %v", e.Err)
}
