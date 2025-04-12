package mocks

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
	"time"
)

// MockS3 is a mock implementation of the s3.ItfS3 interface for testing
type MockS3 struct {
	uploadedFiles map[string][]byte
	deletedFiles  map[string]bool
	baseURL       string
	mutex         sync.RWMutex
	uploadError   error
	presignError  error
	deleteError   error
}

// NewMockS3 creates a new instance of MockS3
func NewMockS3() *MockS3 {
	return &MockS3{
		uploadedFiles: make(map[string][]byte),
		deletedFiles:  make(map[string]bool),
		baseURL:       "https://mock-bucket.s3.amazonaws.com/",
	}
}

// UploadFile mocks the UploadFile method of s3.ItfS3
func (m *MockS3) UploadFile(file *multipart.FileHeader) (string, error) {
	if m.uploadError != nil {
		return "", m.uploadError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if file == nil {
		return "", errors.New("nil file header")
	}

	// Generate a mock unique filename
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)

	// In a real implementation, we would read the file content
	// Here we just store the filename
	m.uploadedFiles[uniqueFileName] = []byte(fmt.Sprintf("mock content for %s", file.Filename))

	// Return the full URL to the file
	return m.baseURL + uniqueFileName, nil
}

// PresignUrl mocks the PresignUrl method of s3.ItfS3
func (m *MockS3) PresignUrl(fileUrl string) (string, error) {
	if m.presignError != nil {
		return "", m.presignError
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Extract the key from the URL
	key := m.extractKeyFromS3Url(fileUrl)

	// Check if file exists
	if _, exists := m.uploadedFiles[key]; !exists && !strings.HasSuffix(fileUrl, key) {
		return "", errors.New("file not found for presigning")
	}

	// Return a mock presigned URL
	return fileUrl + "?presigned=true&expires=900", nil
}

// DeleteFile mocks the DeleteFile method of s3.ItfS3
func (m *MockS3) DeleteFile(fileName string) error {
	if m.deleteError != nil {
		return m.deleteError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if file exists
	_, exists := m.uploadedFiles[fileName]
	if !exists {
		// Check if it's a full URL
		key := m.extractKeyFromS3Url(fileName)
		_, exists = m.uploadedFiles[key]
		if exists {
			fileName = key
		} else {
			return errors.New("file not found for deletion")
		}
	}

	// Mark as deleted
	m.deletedFiles[fileName] = true
	// Actually remove from uploaded files
	delete(m.uploadedFiles, fileName)

	return nil
}

// Helper method to extract key from S3 URL
func (m *MockS3) extractKeyFromS3Url(fileUrl string) string {
	parts := strings.Split(fileUrl, m.baseURL)
	if len(parts) > 1 {
		return parts[1]
	}
	return fileUrl
}

// SetUploadError sets an error to be returned by UploadFile
func (m *MockS3) SetUploadError(err error) {
	m.uploadError = err
}

// SetPresignError sets an error to be returned by PresignUrl
func (m *MockS3) SetPresignError(err error) {
	m.presignError = err
}

// SetDeleteError sets an error to be returned by DeleteFile
func (m *MockS3) SetDeleteError(err error) {
	m.deleteError = err
}

// IsFileUploaded checks if a file with the given name has been uploaded
func (m *MockS3) IsFileUploaded(fileName string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.uploadedFiles[fileName]
	return exists
}

// IsFileDeleted checks if a file with the given name has been deleted
func (m *MockS3) IsFileDeleted(fileName string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	deleted, exists := m.deletedFiles[fileName]
	return exists && deleted
}

// GetUploadedFilesCount returns the number of files currently uploaded
func (m *MockS3) GetUploadedFilesCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.uploadedFiles)
}

// Reset clears all stored files and errors
func (m *MockS3) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.uploadedFiles = make(map[string][]byte)
	m.deletedFiles = make(map[string]bool)
	m.uploadError = nil
	m.presignError = nil
	m.deleteError = nil
}
