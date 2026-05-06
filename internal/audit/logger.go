package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	User      string    `json:"user,omitempty"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes structured audit entries as JSON lines.
type Logger struct {
	writer  io.Writer
	encoder *json.Encoder
}

// NewLogger creates a Logger that writes to w.
// Pass nil to default to os.Stdout.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &Logger{writer: w, encoder: enc}
}

// Log records an audit entry for the given operation and path.
func (l *Logger) Log(operation, path, user string, success bool, message string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		Path:      path,
		User:      user,
		Success:   success,
		Message:   message,
	}
	return l.encoder.Encode(e)
}

// LogRead is a convenience wrapper for read operations.
func (l *Logger) LogRead(path, user string, success bool) error {
	msg := ""
	if !success {
		msg = "read failed"
	}
	return l.Log("read", path, user, success, msg)
}

// LogList is a convenience wrapper for list operations.
func (l *Logger) LogList(path, user string, success bool) error {
	msg := ""
	if !success {
		msg = "list failed"
	}
	return l.Log("list", path, user, success, msg)
}
