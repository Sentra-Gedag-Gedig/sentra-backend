package smtp

import (
	"ProjectGolang/internal/mocks"
	"errors"
	"testing"
)

func TestCreateSmtp(t *testing.T) {
	mockSMTP := mocks.NewMockSMTP()

	testCases := []struct {
		name      string
		email     string
		otp       string
		setError  error
		expectErr bool
	}{
		{
			name:      "Successful OTP email",
			email:     "test@example.com",
			otp:       "12345",
			expectErr: false,
		},
		{
			name:      "Email sending error",
			email:     "test@example.com",
			otp:       "67890",
			setError:  errors.New("SMTP connection failed"),
			expectErr: true,
		},
		{
			name:      "Empty email address",
			email:     "",
			otp:       "54321",
			expectErr: false, // The mock doesn't validate email format
		},
		{
			name:      "Empty OTP",
			email:     "test@example.com",
			otp:       "",
			expectErr: false, // The mock doesn't validate OTP format
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockSMTP.Reset()

			// Set error if specified
			if tc.setError != nil {
				mockSMTP.SetError(tc.setError)
			}

			// Test sending OTP via email
			err := mockSMTP.CreateSmtp(tc.email, tc.otp)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("CreateSmtp() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("CreateSmtp() unexpected error: %v", err)
				return
			}

			// Verify the OTP was stored for the email
			storedOTP, err := mockSMTP.GetSentOTP(tc.email)
			if err != nil {
				t.Errorf("GetSentOTP() unexpected error: %v", err)
				return
			}

			if storedOTP != tc.otp {
				t.Errorf("GetSentOTP() = %v, expected %v", storedOTP, tc.otp)
			}

			// Verify email count increased
			count := mockSMTP.GetSentEmailsCount()
			if count != 1 {
				t.Errorf("Expected 1 email to be sent, got %d", count)
			}
		})
	}
}

func TestMultipleEmails(t *testing.T) {
	mockSMTP := mocks.NewMockSMTP()

	emails := []struct {
		address string
		otp     string
	}{
		{"user1@example.com", "12345"},
		{"user2@example.com", "67890"},
		{"user3@example.com", "54321"},
	}

	// Send multiple emails
	for _, email := range emails {
		err := mockSMTP.CreateSmtp(email.address, email.otp)
		if err != nil {
			t.Fatalf("CreateSmtp() unexpected error: %v", err)
		}
	}

	// Check email count
	count := mockSMTP.GetSentEmailsCount()
	if count != len(emails) {
		t.Errorf("Expected %d emails to be sent, got %d", len(emails), count)
	}

	// Verify each email received the correct OTP
	for _, email := range emails {
		storedOTP, err := mockSMTP.GetSentOTP(email.address)
		if err != nil {
			t.Errorf("GetSentOTP() unexpected error for %s: %v", email.address, err)
			continue
		}

		if storedOTP != email.otp {
			t.Errorf("GetSentOTP() for %s = %v, expected %v",
				email.address, storedOTP, email.otp)
		}
	}
}

func TestNonExistentEmail(t *testing.T) {
	mockSMTP := mocks.NewMockSMTP()

	// Try to get OTP for email that hasn't been sent
	_, err := mockSMTP.GetSentOTP("nonexistent@example.com")

	// Should return an error
	if err == nil {
		t.Error("GetSentOTP() expected error for non-existent email, got nil")
	}
}

func TestReset(t *testing.T) {
	mockSMTP := mocks.NewMockSMTP()

	// Send an email
	err := mockSMTP.CreateSmtp("test@example.com", "12345")
	if err != nil {
		t.Fatalf("CreateSmtp() unexpected error: %v", err)
	}

	// Verify email count is 1
	if count := mockSMTP.GetSentEmailsCount(); count != 1 {
		t.Errorf("Expected 1 email to be sent, got %d", count)
	}

	// Reset the mock
	mockSMTP.Reset()

	// Verify email count is now 0
	if count := mockSMTP.GetSentEmailsCount(); count != 0 {
		t.Errorf("Expected 0 emails after reset, got %d", count)
	}

	// Try to get the previously sent OTP
	_, err = mockSMTP.GetSentOTP("test@example.com")
	if err == nil {
		t.Error("GetSentOTP() expected error after reset, got nil")
	}
}
