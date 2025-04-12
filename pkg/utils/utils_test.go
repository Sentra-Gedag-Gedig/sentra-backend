package utils

import (
	"bytes"
	"mime/multipart"
	"testing"
	"time"
)

func TestNewULIDFromTimestamp(t *testing.T) {
	utils := New()

	testCases := []struct {
		name      string
		timestamp time.Time
		wantErr   bool
	}{
		{
			name:      "Current time",
			timestamp: time.Now(),
			wantErr:   false,
		},
		{
			name:      "Past time",
			timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "Future time",
			timestamp: time.Now().Add(24 * time.Hour),
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ulid, err := utils.NewULIDFromTimestamp(tc.timestamp)

			// Check error status matches expectations
			if (err != nil) != tc.wantErr {
				t.Errorf("NewULIDFromTimestamp() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// If we don't expect an error, do further validation
			if !tc.wantErr {
				if len(ulid) != 26 {
					t.Errorf("NewULIDFromTimestamp() returned ULID of length %d, expected 26", len(ulid))
				}
			}
		})
	}
}

func TestValidateImageFile(t *testing.T) {
	utils := New()

	// Create a mock file header for testing
	createMockFileHeader := func(contentType string, size int64) *multipart.FileHeader {
		return &multipart.FileHeader{
			Filename: "test.jpg",
			Header: map[string][]string{
				"Content-Type": {contentType},
			},
			Size: size,
		}
	}

	testCases := []struct {
		name        string
		fileHeader  *multipart.FileHeader
		wantErr     bool
		errContains string
	}{
		{
			name:       "Valid image file",
			fileHeader: createMockFileHeader("image/jpeg", 1024*1024), // 1MB
			wantErr:    false,
		},
		{
			name:        "Nil file header",
			fileHeader:  nil,
			wantErr:     true,
			errContains: "no file uploaded",
		},
		{
			name:        "File too large",
			fileHeader:  createMockFileHeader("image/jpeg", 10*1024*1024), // 10MB
			wantErr:     true,
			errContains: "file size exceeds limit",
		},
		{
			name:        "Not an image file",
			fileHeader:  createMockFileHeader("application/pdf", 1024*1024),
			wantErr:     true,
			errContains: "not an image",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.ValidateImageFile(tc.fileHeader)

			// Check error status matches expectations
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateImageFile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// If we expect an error, check the error message
			if tc.wantErr && err != nil && tc.errContains != "" {
				if !bytes.Contains([]byte(err.Error()), []byte(tc.errContains)) {
					t.Errorf("ValidateImageFile() error = %v, expected to contain %v", err, tc.errContains)
				}
			}
		})
	}
}

// Helper function to create a mock multipart.File
type mockFile struct {
	*bytes.Reader
}

func (m *mockFile) Close() error {
	return nil
}

func TestConvertFileToBase64(t *testing.T) {
	utils := New()

	testCases := []struct {
		name        string
		fileContent []byte
		wantErr     bool
	}{
		{
			name:        "Valid file content",
			fileContent: []byte("test file content"),
			wantErr:     false,
		},
		{
			name:        "Empty file content",
			fileContent: []byte{},
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockFile := &mockFile{bytes.NewReader(tc.fileContent)}

			base64Str, err := utils.ConvertFileToBase64(mockFile)

			// Check error status matches expectations
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertFileToBase64() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// If we don't expect an error, do further validation
			if !tc.wantErr {
				if len(base64Str) == 0 && len(tc.fileContent) > 0 {
					t.Errorf("ConvertFileToBase64() returned empty string for non-empty file")
				}
			}
		})
	}
}
