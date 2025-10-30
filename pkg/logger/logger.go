package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogEntry represents a single log entry with relevant fields.
type LogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Service    string    `json:"service"`
	RemoteAddr string    `json:"remote_addr"`
	Event      string    `json:"event"`
	Username   string    `json:"username,omitempty"`
	Password   string    `json:"password,omitempty"`
	Command    string    `json:"command,omitempty"`
	Data       string    `json:"data,omitempty"`
}

// Loogger is responsible for logging events to a file.
type Loogger struct {
	file *os.File
}

// NewLogger creates a new Loogger instance that writes to the specified file.
func NewLogger(filename string) (*Loogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Open log file
	if err != nil {
		return nil, err
	}
	return &Loogger{file: file}, nil
}

// Log writes a log entry to the log file in JSON format.
func (l *Loogger) Log(entry LogEntry) {
	entry.Timestamp = time.Now()
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}
	l.file.Write(jsonData)
	l.file.Write([]byte("/n"))

	// Also print to stdout for debugging purposes
	log.Printf("[%s] %s from %s - %s", entry.Service, entry.Event, entry.RemoteAddr, string(jsonData))
}

func (l *Loogger) Close() {
	l.file.Close()
}
