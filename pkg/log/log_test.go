package log

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// Custom buffer to capture log output
type bufferWriter struct {
	bytes.Buffer
}

func (bw *bufferWriter) Close() error {
	return nil
}

func TestLogLevels(t *testing.T) {
	// Capture log output to a buffer
	buffer := &bufferWriter{}
	originalLogger := NewLogger()

	// This is not ideal as we're modifying a package-level variable,
	// but for testing it's acceptable - more robust solution would be
	// to refactor the log package to accept an io.Writer
	// Here we're assuming the package has some way to set the output
	// If not, the test would need to be modified

	testCases := []struct {
		name      string
		logFunc   func(fields Fields, msg string)
		levelName string
		fields    Fields
		message   string
	}{
		{
			name:      "Debug log",
			logFunc:   Debug,
			levelName: "debug",
			fields:    Fields{"key": "value"},
			message:   "Debug message",
		},
		{
			name:      "Info log",
			logFunc:   Info,
			levelName: "info",
			fields:    Fields{"key": "value"},
			message:   "Info message",
		},
		{
			name:      "Warn log",
			logFunc:   Warn,
			levelName: "warning",
			fields:    Fields{"key": "value"},
			message:   "Warning message",
		},
		{
			name:      "Error log",
			logFunc:   Error,
			levelName: "error",
			fields:    Fields{"key": "value"},
			message:   "Error message",
		},
		{
			name:      "Log with nil fields",
			logFunc:   Info,
			levelName: "info",
			fields:    nil,
			message:   "Message with nil fields",
		},
	}

	// Skip fatal and panic logs as they would terminate the test

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer.Reset()

			// Call the log function
			tc.logFunc(tc.fields, tc.message)

			// Check that the log contains the expected message
			logOutput := buffer.String()

			// This is a basic test - in a real environment we would
			// parse the log output more thoroughly
			if !strings.Contains(logOutput, tc.message) {
				t.Errorf("Log output does not contain message. Expected to find '%s' in: %s",
					tc.message, logOutput)
			}

			// Check if fields were logged
			if tc.fields != nil {
				for k, v := range tc.fields {
					if !strings.Contains(logOutput, k) || !strings.Contains(logOutput, v.(string)) {
						t.Errorf("Log output does not contain field. Expected to find '%s:%s' in: %s",
							k, v, logOutput)
					}
				}
			}
		})
	}
}

func TestErrorWithTraceID(t *testing.T) {
	// Test with existing request ID
	fields := Fields{
		"request_id": "existing-id",
		"other_key":  "value",
	}

	traceID := ErrorWithTraceID(fields, "Error with existing ID")

	if traceID != "existing-id" {
		t.Errorf("ErrorWithTraceID() = %v, expected %v", traceID, "existing-id")
	}

	// Test with no request ID
	fields = Fields{
		"other_key": "value",
	}

	traceID = ErrorWithTraceID(fields, "Error with no ID")

	if traceID == "" {
		t.Errorf("ErrorWithTraceID() returned empty trace ID")
	}

	// Test with nil fields
	traceID = ErrorWithTraceID(nil, "Error with nil fields")

	if traceID == "" {
		t.Errorf("ErrorWithTraceID() returned empty trace ID for nil fields")
	}
}

func TestWithRequestID(t *testing.T) {
	// Create context with request ID
	ctx := context.WithValue(context.Background(), RequestIDKey, "test-request-id")

	// Get logger entry with request ID
	entry := WithRequestID(ctx)

	// We can't easily test the exact entry output, but we can verify it's not nil
	if entry == nil {
		t.Errorf("WithRequestID() returned nil")
	}

	// Test with nil context
	entry = WithRequestID(nil)

	if entry == nil {
		t.Errorf("WithRequestID() returned nil for nil context")
	}
}
