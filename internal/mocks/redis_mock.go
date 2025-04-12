package mocks

import (
	"context"
	"errors"
	"sync"
	"time"
)

// MockRedis is a mock implementation of the redis.IRedis interface for testing
type MockRedis struct {
	data     map[string]string
	expiry   map[string]time.Time
	mutex    sync.RWMutex
	setError error
	getError error
}

// NewMockRedis creates a new instance of MockRedis
func NewMockRedis() *MockRedis {
	return &MockRedis{
		data:   make(map[string]string),
		expiry: make(map[string]time.Time),
	}
}

// SetOTP mocks the SetOTP method of redis.IRedis
func (m *MockRedis) SetOTP(ctx context.Context, key string, code string, expiration time.Duration) error {
	if m.setError != nil {
		return m.setError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = code
	m.expiry[key] = time.Now().Add(expiration)

	return nil
}

// GetOTP mocks the GetOTP method of redis.IRedis
func (m *MockRedis) GetOTP(ctx context.Context, key string) (string, error) {
	if m.getError != nil {
		return "", m.getError
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	code, exists := m.data[key]
	if !exists {
		return "", errors.New("key not found")
	}

	expiryTime, exists := m.expiry[key]
	if !exists || time.Now().After(expiryTime) {
		delete(m.data, key)
		delete(m.expiry, key)
		return "", errors.New("key expired")
	}

	return code, nil
}

// SetGetError sets an error to be returned by GetOTP
func (m *MockRedis) SetGetError(err error) {
	m.getError = err
}

// SetSetError sets an error to be returned by SetOTP
func (m *MockRedis) SetSetError(err error) {
	m.setError = err
}

// ClearErrors clears any set errors
func (m *MockRedis) ClearErrors() {
	m.getError = nil
	m.setError = nil
}

// Reset clears all stored data
func (m *MockRedis) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = make(map[string]string)
	m.expiry = make(map[string]time.Time)
	m.getError = nil
	m.setError = nil
}
