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

// ratelimit.go

package aiyou

import (
	"context"
	"sync"
	"time"
)

// RateLimiter contrôle le taux de requêtes
type RateLimiter struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
	mutex      sync.Mutex
	logger     Logger
}

// RateLimiterConfig contient les options de configuration
type RateLimiterConfig struct {
	RequestsPerSecond float64
	BurstSize         int
	WaitTimeout       time.Duration
}

// NewRateLimiter crée un nouveau rate limiter
func NewRateLimiter(config RateLimiterConfig, logger Logger) *RateLimiter {
	return &RateLimiter{
		tokens:     float64(config.BurstSize),
		capacity:   float64(config.BurstSize),
		refillRate: config.RequestsPerSecond,
		lastRefill: time.Now(),
		logger:     logger,
	}
}

// Wait attend qu'un token soit disponible
func (r *RateLimiter) Wait(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.refill()

	if r.tokens >= 1 {
		r.tokens--
		return nil
	}

	// Calculer le temps d'attente
	waitTime := time.Duration((1 - r.tokens) / r.refillRate * float64(time.Second))

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(waitTime):
		r.tokens--
		return nil
	}
}

// refill recharge les tokens
func (r *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens = min(r.capacity, r.tokens+(elapsed*r.refillRate))
	r.lastRefill = now
}

func (r *RateLimiter) GetWaitTime() time.Duration {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.refill()
	if r.tokens >= 1 {
		return 0
	}

	return time.Duration((1 - r.tokens) / r.refillRate * float64(time.Second))
}
