// Package logs provides debug logging
package log

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
	// DefaultSize Default buffer size if any
	DefaultSize = 256
	// DefaultFormat Default formatter
	DefaultFormat = TextFormat
)

// Log is debug logs interface for reading and writing logs
type Log interface {
	// Read reads logs entries from the logger
	Read(...ReadOption) ([]Record, error)
	// Write writes records to logs
	Write(Record) error
	// Stream logs records
	Stream() (Stream, error)
}

// Record is logs record entry
type Record struct {
	// Timestamp of logged event
	Timestamp time.Time `json:"timestamp"`
	// Metadata to enrich logs record
	Metadata map[string]string `json:"metadata"`
	// Value contains logs entry
	Message interface{} `json:"message"`
}

// Stream returns a logs stream
type Stream interface {
	Chan() <-chan Record
	Stop() error
}

// FormatFunc is a function which formats the output
type FormatFunc func(Record) string

// TextFormat returns text format
func TextFormat(r Record) string {
	t := r.Timestamp.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s %v ", t, r.Message)
}

// JSONFormat is a json Format func
func JSONFormat(r Record) string {
	b, _ := json.Marshal(r)
	return string(b) + " "
}
