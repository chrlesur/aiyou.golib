/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

// File: pkg/aiyou/logging.go

package aiyou

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// DEBUG level for detailed debugging information
	DEBUG LogLevel = iota
	// INFO level for general operational information
	INFO
	// WARN level for warning messages
	WARN
	// ERROR level for error messages
	ERROR
)

// Pool de buffers pour le logging
var logBufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Logger interface extends the standard log.Logger interface with additional methods
// for different log levels and the ability to set the logging level.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	SetLevel(level LogLevel)
}

// defaultLogger implements the Logger interface
type defaultLogger struct {
	level  LogLevel
	writer io.Writer
}

// NewDefaultLogger creates a new instance of defaultLogger with the specified
// output writer, prefix, and flags. It sets the initial log level to INFO.
func NewDefaultLogger(w io.Writer) *defaultLogger {
	return &defaultLogger{
		level:  INFO,
		writer: w,
	}
}

// SetLevel sets the logging level for the logger. Only messages with a severity
// level equal to or higher than the set level will be logged.
func (l *defaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

// log logs a message at the specified level
func (l *defaultLogger) log(level LogLevel, format string, args ...interface{}) {
	if level >= l.level {
		// Récupérer un buffer du pool
		buf := logBufferPool.Get().(*bytes.Buffer)
		buf.Reset() // Réinitialiser le buffer pour réutilisation
		defer func() {
			// Remettre le buffer dans le pool après utilisation
			logBufferPool.Put(buf)
		}()

		// Get file and line information
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "unknown"
			line = 0
		}

		// Extract just the filename from the full path
		filename := filepath.Base(file)

		// Construire le message dans le buffer
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(buf, "[%s] %s %s:%d: ", timestamp, level.String(), filename, line)
		fmt.Fprintf(buf, format, args...)
		buf.WriteByte('\n')

		// Écrire le contenu du buffer en une seule opération
		l.writer.Write(buf.Bytes())
	}
}

// Debugf logs a debug message
func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Infof logs an info message
func (l *defaultLogger) Infof(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warnf logs a warning message
func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Errorf logs an error message
func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
