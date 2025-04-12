package whatsapp_test

import (
	"ProjectGolang/internal/mocks"
	"context"
	"errors"
	"testing"
)

func TestSendMessage(t *testing.T) {
	mockWA := mocks.NewMockWhatsapp()
	ctx := context.Background()

	testCases := []struct {
		name        string
		phoneNumber string
		message     string
		connected   bool
		setError    error
		expectErr   bool
	}{
		{
			name:        "Successfully send message",
			phoneNumber: "+6281234567890",
			message:     "Hello, this is a test message",
			connected:   true,
			expectErr:   false,
		},
		{
			name:        "Not connected",
			phoneNumber: "+6281234567890",
			message:     "This message should fail",
			connected:   false,
			expectErr:   true,
		},
		{
			name:        "Error during sending",
			phoneNumber: "+6281234567890",
			message:     "This message should fail",
			connected:   true,
			setError:    errors.New("sending error"),
			expectErr:   true,
		},
		{
			name:        "Empty phone number",
			phoneNumber: "",
			message:     "Message with empty phone number",
			connected:   true,
			expectErr:   false, // The mock doesn't validate phone number format
		},
		{
			name:        "Empty message",
			phoneNumber: "+6281234567890",
			message:     "",
			connected:   true,
			expectErr:   false, // The mock doesn't validate message content
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockWA.Reset()

			// Set connected status
			mockWA.SetConnectedStatus(tc.connected)

			// Set error if specified
			if tc.setError != nil {
				mockWA.SetSendError(tc.setError)
			}

			// Test sending message
			err := mockWA.SendMessage(ctx, tc.phoneNumber, tc.message)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("SendMessage() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("SendMessage() unexpected error: %v", err)
				return
			}

			// Verify message was stored
			lastMsg, err := mockWA.GetLastSentMessage(tc.phoneNumber)
			if err != nil {
				t.Errorf("GetLastSentMessage() unexpected error: %v", err)
				return
			}

			if lastMsg != tc.message {
				t.Errorf("GetLastSentMessage() = %v, expected %v", lastMsg, tc.message)
			}

			// Verify message count increased
			count := mockWA.GetTotalSentMessagesCount()
			if count != 1 {
				t.Errorf("Expected 1 message to be sent, got %d", count)
			}
		})
	}
}

func TestSendMultipleMessages(t *testing.T) {
	mockWA := mocks.NewMockWhatsapp()
	ctx := context.Background()

	// Define test messages
	messages := []struct {
		phoneNumber string
		content     string
	}{
		{"+6281234567890", "Message 1"},
		{"+6281234567890", "Message 2"},
		{"+6289876543210", "Message 3"},
	}

	// Send multiple messages
	for _, msg := range messages {
		err := mockWA.SendMessage(ctx, msg.phoneNumber, msg.content)
		if err != nil {
			t.Fatalf("SendMessage() unexpected error: %v", err)
		}
	}

	// Check total message count
	totalCount := mockWA.GetTotalSentMessagesCount()
	if totalCount != len(messages) {
		t.Errorf("Expected %d total messages, got %d", len(messages), totalCount)
	}

	// Check unique recipients count
	uniqueRecipients := mockWA.GetUniqueRecipientsCount()
	expectedUniqueRecipients := 2 // We used 2 different phone numbers
	if uniqueRecipients != expectedUniqueRecipients {
		t.Errorf("Expected %d unique recipients, got %d", expectedUniqueRecipients, uniqueRecipients)
	}

	// Check messages for first phone number
	firstPhoneNumber := "+6281234567890"
	messagesForFirst := mockWA.GetSentMessages(firstPhoneNumber)
	if len(messagesForFirst) != 2 {
		t.Errorf("Expected 2 messages for %s, got %d", firstPhoneNumber, len(messagesForFirst))
	}

	// Check last message for first phone number
	lastMsgForFirst, err := mockWA.GetLastSentMessage(firstPhoneNumber)
	if err != nil {
		t.Errorf("GetLastSentMessage() unexpected error: %v", err)
	} else if lastMsgForFirst != "Message 2" {
		t.Errorf("Last message for %s = %s, expected %s",
			firstPhoneNumber, lastMsgForFirst, "Message 2")
	}
}

func TestDisconnect(t *testing.T) {
	mockWA := mocks.NewMockWhatsapp()

	testCases := []struct {
		name      string
		setError  error
		expectErr bool
	}{
		{
			name:      "Successfully disconnect",
			expectErr: false,
		},
		{
			name:      "Error during disconnection",
			setError:  errors.New("disconnection error"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockWA.Reset()

			// Set error if specified
			if tc.setError != nil {
				mockWA.SetDisconnectError(tc.setError)
			}

			// Test disconnection
			err := mockWA.Disconnect()

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("Disconnect() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("Disconnect() unexpected error: %v", err)
				return
			}

			// Verify connected status
			if mockWA.IsConnected() {
				t.Error("Disconnect() did not set connected status to false")
			}

			// Try to send a message after disconnection
			ctx := context.Background()
			err = mockWA.SendMessage(ctx, "+6281234567890", "This should fail")
			if err == nil {
				t.Error("SendMessage() after disconnect should return error, got nil")
			}
		})
	}
}

func TestIsConnected(t *testing.T) {
	mockWA := mocks.NewMockWhatsapp()

	// By default, the mock is connected
	if !mockWA.IsConnected() {
		t.Error("New mock should be connected by default")
	}

	// Disconnect
	err := mockWA.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect() unexpected error: %v", err)
	}

	// Check connected status
	if mockWA.IsConnected() {
		t.Error("IsConnected() should return false after disconnection")
	}

	// Manually set connected status back to true
	mockWA.SetConnectedStatus(true)

	// Check connected status again
	if !mockWA.IsConnected() {
		t.Error("IsConnected() should return true after setting connected status")
	}
}

func TestReset(t *testing.T) {
	mockWA := mocks.NewMockWhatsapp()
	ctx := context.Background()

	// Send a message
	err := mockWA.SendMessage(ctx, "+6281234567890", "Test message")
	if err != nil {
		t.Fatalf("SendMessage() unexpected error: %v", err)
	}

	// Set an error
	mockWA.SetSendError(errors.New("sending error"))

	// Disconnect
	mockWA.SetConnectedStatus(false)

	// Reset the mock
	mockWA.Reset()

	// Check that the mock is reset
	if !mockWA.IsConnected() {
		t.Error("Reset() did not restore connected status")
	}

	// Check that message count is reset
	count := mockWA.GetTotalSentMessagesCount()
	if count != 0 {
		t.Errorf("Reset() did not clear messages, got count %d", count)
	}

	// Check that error is cleared
	err = mockWA.SendMessage(ctx, "+6281234567890", "Test message after reset")
	if err != nil {
		t.Errorf("SendMessage() after reset unexpected error: %v", err)
	}
}
