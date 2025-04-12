package mocks

import (
	"context"
	"errors"
	"sync"
)

// MockWhatsapp is a mock implementation of the whatsapp.IWhatsappSender interface for testing
type MockWhatsapp struct {
	mutex           sync.RWMutex
	sentMessages    map[string][]string
	isConnected     bool
	sendError       error
	disconnectError error
}

// NewMockWhatsapp creates a new instance of MockWhatsapp
func NewMockWhatsapp() *MockWhatsapp {
	return &MockWhatsapp{
		sentMessages: make(map[string][]string),
		isConnected:  true,
	}
}

// SendMessage implements the SendMessage method of whatsapp.IWhatsappSender
func (m *MockWhatsapp) SendMessage(ctx context.Context, phoneNumber, message string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.sendError != nil {
		return m.sendError
	}

	if !m.isConnected {
		return errors.New("not connected to WhatsApp")
	}

	// Store the message
	if _, exists := m.sentMessages[phoneNumber]; !exists {
		m.sentMessages[phoneNumber] = make([]string, 0)
	}
	m.sentMessages[phoneNumber] = append(m.sentMessages[phoneNumber], message)

	return nil
}

// Disconnect implements the Disconnect method of whatsapp.IWhatsappSender
func (m *MockWhatsapp) Disconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.disconnectError != nil {
		return m.disconnectError
	}

	m.isConnected = false
	return nil
}

// IsConnected implements the IsConnected method of whatsapp.IWhatsappSender
func (m *MockWhatsapp) IsConnected() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.isConnected
}

// GetSentMessages returns all messages sent to a specific phone number
func (m *MockWhatsapp) GetSentMessages(phoneNumber string) []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if messages, exists := m.sentMessages[phoneNumber]; exists {
		return messages
	}
	return []string{}
}

// GetLastSentMessage returns the last message sent to a specific phone number
func (m *MockWhatsapp) GetLastSentMessage(phoneNumber string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	messages, exists := m.sentMessages[phoneNumber]
	if !exists || len(messages) == 0 {
		return "", errors.New("no messages sent to this phone number")
	}

	return messages[len(messages)-1], nil
}

// GetTotalSentMessagesCount returns the total number of messages sent
func (m *MockWhatsapp) GetTotalSentMessagesCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := 0
	for _, messages := range m.sentMessages {
		total += len(messages)
	}
	return total
}

// GetUniqueRecipientsCount returns the number of unique phone numbers that received messages
func (m *MockWhatsapp) GetUniqueRecipientsCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.sentMessages)
}

// SetSendError sets an error to be returned by SendMessage
func (m *MockWhatsapp) SetSendError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.sendError = err
}

// SetDisconnectError sets an error to be returned by Disconnect
func (m *MockWhatsapp) SetDisconnectError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.disconnectError = err
}

// SetConnectedStatus sets the connected status
func (m *MockWhatsapp) SetConnectedStatus(connected bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.isConnected = connected
}

// Reset clears all sent messages and errors
func (m *MockWhatsapp) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.sentMessages = make(map[string][]string)
	m.isConnected = true
	m.sendError = nil
	m.disconnectError = nil
}
