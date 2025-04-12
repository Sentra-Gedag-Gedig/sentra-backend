package mocks

import (
	"errors"
	"sync"
)

// MockSMTP is a mock implementation of the smtp.ItfSmtp interface for testing
type MockSMTP struct {
	sentEmails map[string]string // Map of email address to OTP
	mutex      sync.RWMutex
	error      error
}

// NewMockSMTP creates a new instance of MockSMTP
func NewMockSMTP() *MockSMTP {
	return &MockSMTP{
		sentEmails: make(map[string]string),
	}
}

// CreateSmtp mocks the CreateSmtp method of smtp.ItfSmtp
func (m *MockSMTP) CreateSmtp(userEmail string, otp string) error {
	if m.error != nil {
		return m.error
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.sentEmails[userEmail] = otp
	return nil
}

// SetError sets an error to be returned by CreateSmtp
func (m *MockSMTP) SetError(err error) {
	m.error = err
}

// ClearError clears any set error
func (m *MockSMTP) ClearError() {
	m.error = nil
}

// GetSentOTP returns the OTP sent to a specific email
func (m *MockSMTP) GetSentOTP(email string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	otp, exists := m.sentEmails[email]
	if !exists {
		return "", errors.New("no OTP sent to this email")
	}

	return otp, nil
}

// Reset clears all sent emails
func (m *MockSMTP) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.sentEmails = make(map[string]string)
	m.error = nil
}

// GetSentEmailsCount returns the number of emails sent
func (m *MockSMTP) GetSentEmailsCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.sentEmails)
}
