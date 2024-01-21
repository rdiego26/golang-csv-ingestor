package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Logger represents a simple logger with emojis and mutex.
type Logger struct {
	mu sync.Mutex
}

// LogLevel represents the severity of a log entry.
type LogLevel int

const (
	// Info level indicates informational messages.
	Info LogLevel = iota

	// Warning level indicates warnings.
	Warning

	// Fatal level indicates fatal error.
	Fatal

	// Debug level indicates debug.
	Debug

	// Error level indicates errors.
	Error
)

var emojiMap = map[LogLevel]string{
	Info:    "✨",
	Warning: "⚠️",
	Fatal:   "💣",
	Error:   "❌",
	Debug:   "🔎",
}

// Log writes a log entry to the standard output with emoji based on log level.
func (l *Logger) Log(level LogLevel, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get the current time in a readable format.
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Get the emoji based on log level.
	emoji, exists := emojiMap[level]
	if !exists {
		emoji = "ℹ️" // Default to information emoji if level is unknown.
	}

	// Create the log entry.
	logEntry := fmt.Sprintf("[%s] %s: %s\n", currentTime, emoji, message)

	// Write the log entry to the standard output.
	if level != Fatal {
		fmt.Print(logEntry)
	} else {
		log.Fatalf(logEntry)
	}

}
