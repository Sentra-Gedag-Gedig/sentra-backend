package s3

import (
	"ProjectGolang/internal/mocks"
	"errors"
	"mime/multipart"
	"strings"
	"testing"
)

// MockFileHeader creates a mock multipart.FileHeader for testing
func MockFileHeader(filename string, size int64) *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: filename,
		Size:     size,
		Header:   make(map[string][]string),
	}
}

func TestUploadFile(t *testing.T) {
	mockS3 := mocks.NewMockS3()

	testCases := []struct {
		name      string
		file      *multipart.FileHeader
		setError  error
		expectErr bool
	}{
		{
			name:      "Successful upload",
			file:      MockFileHeader("test.jpg", 1024),
			expectErr: false,
		},
		{
			name:      "Nil file",
			file:      nil,
			expectErr: true,
		},
		{
			name:      "Upload error",
			file:      MockFileHeader("error.jpg", 1024),
			setError:  errors.New("upload failed"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockS3.Reset()

			// Set error if specified
			if tc.setError != nil {
				mockS3.SetUploadError(tc.setError)
			}

			// Test file upload
			url, err := mockS3.UploadFile(tc.file)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("UploadFile() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("UploadFile() unexpected error: %v", err)
				return
			}

			// Verify URL format
			if !strings.HasPrefix(url, "https://") {
				t.Errorf("UploadFile() returned invalid URL: %s", url)
			}

			// Verify file count increased
			count := mockS3.GetUploadedFilesCount()
			if count != 1 {
				t.Errorf("Expected 1 file to be uploaded, got %d", count)
			}
		})
	}
}

func TestPresignUrl(t *testing.T) {
	mockS3 := mocks.NewMockS3()

	// First upload a file to get a valid URL
	file := MockFileHeader("test.jpg", 1024)
	uploadedUrl, err := mockS3.UploadFile(file)
	if err != nil {
		t.Fatalf("Failed to set up test: %v", err)
	}

	testCases := []struct {
		name      string
		fileUrl   string
		setError  error
		expectErr bool
	}{
		{
			name:      "Successful presign",
			fileUrl:   uploadedUrl,
			expectErr: false,
		},
		{
			name:      "Presign non-existent file",
			fileUrl:   "https://mock-bucket.s3.amazonaws.com/non-existent.jpg",
			expectErr: true,
		},
		{
			name:      "Presign error",
			fileUrl:   uploadedUrl,
			setError:  errors.New("presign failed"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset errors
			mockS3.SetPresignError(nil)

			// Set error if specified
			if tc.setError != nil {
				mockS3.SetPresignError(tc.setError)
			}

			// Test URL presigning
			presignedUrl, err := mockS3.PresignUrl(tc.fileUrl)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("PresignUrl() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("PresignUrl() unexpected error: %v", err)
				return
			}

			// Verify presigned URL format
			if !strings.Contains(presignedUrl, "?presigned=true") {
				t.Errorf("PresignUrl() returned invalid presigned URL: %s", presignedUrl)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	mockS3 := mocks.NewMockS3()

	// First upload a file to get a valid URL
	file := MockFileHeader("test.jpg", 1024)
	uploadedUrl, err := mockS3.UploadFile(file)
	if err != nil {
		t.Fatalf("Failed to set up test: %v", err)
	}

	// Extract the key from the URL
	parts := strings.Split(uploadedUrl, "https://mock-bucket.s3.amazonaws.com/")
	fileName := parts[1]

	testCases := []struct {
		name      string
		fileName  string
		setError  error
		expectErr bool
	}{
		{
			name:      "Successful delete with filename",
			fileName:  fileName,
			expectErr: false,
		},
		{
			name:      "Successful delete with full URL",
			fileName:  uploadedUrl,
			expectErr: false,
		},
		{
			name:      "Delete non-existent file",
			fileName:  "non-existent.jpg",
			expectErr: true,
		},
		{
			name:      "Delete error",
			fileName:  fileName,
			setError:  errors.New("delete failed"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock for each test
			mockS3.Reset()

			// Upload a file again for this test
			uploadedUrl, _ := mockS3.UploadFile(file)
			parts := strings.Split(uploadedUrl, "https://mock-bucket.s3.amazonaws.com/")
			fileName := parts[1]

			// Override fileName if the test is using the full URL
			if tc.fileName == uploadedUrl {
				tc.fileName = uploadedUrl
			} else if tc.fileName == "non-existent.jpg" {
				// Keep it as is
			} else {
				tc.fileName = fileName
			}

			// Set error if specified
			if tc.setError != nil {
				mockS3.SetDeleteError(tc.setError)
			}

			// Test file deletion
			err := mockS3.DeleteFile(tc.fileName)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("DeleteFile() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("DeleteFile() unexpected error: %v", err)
				return
			}

			// Verify file was deleted
			if !mockS3.IsFileDeleted(fileName) && !tc.expectErr {
				t.Errorf("DeleteFile() file was not marked as deleted")
			}

			// Verify file count decreased
			count := mockS3.GetUploadedFilesCount()
			if count != 0 && !tc.expectErr {
				t.Errorf("Expected 0 files after deletion, got %d", count)
			}
		})
	}
}
