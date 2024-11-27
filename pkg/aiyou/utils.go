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

// File: pkg/aiyou/utils.go

package aiyou

import (
	"fmt"
	"regexp"
)

var (
	// Patterns for sensitive information
	emailPattern    = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	tokenPattern    = regexp.MustCompile(`(Bearer\s+)([A-Za-z0-9-._~+/]+=*)`)
	passwordPattern = regexp.MustCompile(`("password"\s*:\s*)"[^"]*"`)
)

// MaskSensitiveInfo masks sensitive information in the given string.
// It replaces email addresses, JWT tokens, and password fields with redacted placeholders.
func MaskSensitiveInfo(input string) string {
	// Mask email addresses
	emailPattern := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	maskedInput := emailPattern.ReplaceAllString(input, "[EMAIL REDACTED]")

	// Mask JWT tokens
	tokenPattern := regexp.MustCompile(`Bearer\s+[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`)
	maskedInput = tokenPattern.ReplaceAllString(maskedInput, "Bearer [TOKEN REDACTED]")

	// Mask passwords
	passwordPattern := regexp.MustCompile(`("password"\s*:\s*)"[^"]*"`)
	maskedInput = passwordPattern.ReplaceAllString(maskedInput, `${1}"[PASSWORD REDACTED]"`)

	return maskedInput
}

// SafeLog returns a closure that can be used to safely log messages with sensitive info masked.
// It takes a Logger as input and returns a function that can be used to log messages at different levels,
// while automatically masking sensitive information.
func SafeLog(logger Logger) func(level LogLevel, format string, args ...interface{}) {
	return func(level LogLevel, format string, args ...interface{}) {
		// Mask sensitive info in format and args
		safeFormat := MaskSensitiveInfo(format)
		safeArgs := make([]interface{}, len(args))
		for i, arg := range args {
			safeArgs[i] = MaskSensitiveInfo(fmt.Sprint(arg))
		}

		// Create the masked message
		maskedMessage := fmt.Sprintf(safeFormat, safeArgs...)

		// Log the masked message
		switch level {
		case DEBUG:
			logger.Debugf("%s", maskedMessage)
		case INFO:
			logger.Infof("%s", maskedMessage)
		case WARN:
			logger.Warnf("%s", maskedMessage)
		case ERROR:
			logger.Errorf("%s", maskedMessage)
		}
	}
}
