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
	emailRegex    = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	tokenRegex    = regexp.MustCompile(`Bearer\s+[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`)
	passwordRegex = regexp.MustCompile(`("password"\s*:\s*)"[^"]*"`)
)

// MaskSensitiveInfo masks sensitive information in the given string.
// It replaces email addresses, JWT tokens, and password fields with redacted placeholders.
func MaskSensitiveInfo(input string) string {
	maskedInput := emailRegex.ReplaceAllString(input, "[EMAIL REDACTED]")
	maskedInput = tokenRegex.ReplaceAllString(maskedInput, "Bearer [TOKEN REDACTED]")
	maskedInput = passwordRegex.ReplaceAllString(maskedInput, "${1}\"[PASSWORD REDACTED]\"")
	return maskedInput
}

// SafeLog returns a closure that can be used to safely log messages with sensitive info masked.
// It takes a Logger as input and returns a function that can be used to log messages at different levels,
// while automatically masking sensitive information.
func SafeLog(logger Logger) func(level LogLevel, format string, args ...interface{}) {
	return func(level LogLevel, format string, args ...interface{}) {
		safeFormat := MaskSensitiveInfo(format)
		safeArgs := make([]interface{}, len(args))
		for i, arg := range args {
			safeArgs[i] = MaskSensitiveInfo(fmt.Sprint(arg))
		}

		maskedMessage := fmt.Sprintf(safeFormat, safeArgs...)

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
