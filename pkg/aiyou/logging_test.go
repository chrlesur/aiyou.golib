// File: pkg/aiyou/logging_test.go

package aiyou

import (
    "bytes"
    "strings"
    "testing"
)

func TestDefaultLogger(t *testing.T) {
    var buf bytes.Buffer
    logger := NewDefaultLogger(&buf)

    testCases := []struct {
        level   LogLevel
        message string
    }{
        {DEBUG, "Debug message"},
        {INFO, "Info message"},
        {WARN, "Warning message"},
        {ERROR, "Error message"},
    }

    for _, tc := range testCases {
        buf.Reset()
        logger.SetLevel(tc.level)

        logger.Debugf("Debug message")
        logger.Infof("Info message")
        logger.Warnf("Warning message")
        logger.Errorf("Error message")

        output := buf.String()
        if !strings.Contains(output, tc.message) {
            t.Errorf("Expected log to contain %q at level %v, but got: %s", tc.message, tc.level, output)
        }

        for _, otherMsg := range testCases {
            if otherMsg.level < tc.level && strings.Contains(output, otherMsg.message) {
                t.Errorf("Log should not contain %q at level %v, but got: %s", otherMsg.message, tc.level, output)
            }
        }
    }
}

func TestMaskSensitiveInfo(t *testing.T) {
    testCases := []struct {
        input    string
        expected string
    }{
        {
            "Email: user@example.com, Token: Bearer abc123.def456.ghi789",
            "Email: [EMAIL REDACTED], Token: Bearer [TOKEN REDACTED]",
        },
        {
            `{"email": "user@example.com", "password": "secret123"}`,
            `{"email": "[EMAIL REDACTED]", "password": "[PASSWORD REDACTED]"}`,
        },
        {
            "No sensitive info here",
            "No sensitive info here",
        },
    }

    for _, tc := range testCases {
        result := MaskSensitiveInfo(tc.input)
        if result != tc.expected {
            t.Errorf("MaskSensitiveInfo(%q) = %q, want %q", tc.input, result, tc.expected)
        }
    }
}

func TestSafeLog(t *testing.T) {
    var buf bytes.Buffer
    logger := NewDefaultLogger(&buf)
    safeLog := SafeLog(logger)

    testCases := []struct {
        level    LogLevel
        format   string
        args     []interface{}
        expected string
    }{
        {
            INFO,
            "User %s logged in with token %s",
            []interface{}{"user@example.com", "Bearer abc123.def456.ghi789"},
            "User [EMAIL REDACTED] logged in with token Bearer [TOKEN REDACTED]",
        },
        {
            ERROR,
            "Failed to authenticate: %v",
            []interface{}{&AuthenticationError{Message: "Invalid password for user@example.com"}},
            "Failed to authenticate: Authentication error: Invalid password for [EMAIL REDACTED]",
        },
    }

    for _, tc := range testCases {
        buf.Reset()
        safeLog(tc.level, tc.format, tc.args...)
        output := buf.String()
        if !strings.Contains(output, tc.expected) {
            t.Errorf("SafeLog output should contain %q, but got: %s", tc.expected, output)
        }
    }
}