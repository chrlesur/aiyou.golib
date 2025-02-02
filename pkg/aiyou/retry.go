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
	"time"
)

// retryOperation exécute une opération avec une logique de retry.
// Elle réessaie l'opération jusqu'à ce qu'elle réussisse, que le nombre maximum
// de tentatives soit atteint, ou que le contexte expire.
func retryOperation(ctx context.Context, logger Logger, maxRetries int, initialDelay time.Duration, operation func() error) error {
    var err error
    delay := initialDelay

    for attempt := 0; attempt <= maxRetries; attempt++ {
        logger.Debugf("Attempt %d of %d", attempt+1, maxRetries+1)
        
        err = operation()
        if err == nil {
            logger.Debugf("Operation successful on attempt %d", attempt+1)
            return nil
        }

        if !isRetryableError(err) {
            logger.Errorf("Non-retryable error encountered: %v", err)
            return err
        }

        select {
        case <-ctx.Done():
            logger.Warnf("Context cancelled, stopping retries")
            return ctx.Err()
        default:
            if attempt < maxRetries {
                logger.Infof("Retrying after %v", delay)
                time.Sleep(delay)
                delay *= 2 // exponential backoff
            }
        }
    }

    logger.Errorf("Max retries reached, last error: %v", err)
    return err
}

// isRetryableError détermine si une erreur peut être retentée.
// Actuellement, seules les erreurs de type NetworkError et RateLimitError sont considérées comme retryables.
func isRetryableError(err error) bool {
    switch err.(type) {
    case *NetworkError, *RateLimitError:
        return true
    default:
        return false
    }
}
