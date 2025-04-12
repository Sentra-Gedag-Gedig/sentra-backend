package redis

import (
	"ProjectGolang/internal/mocks"
	"context"
	"errors"
	"testing"
	"time"
)

func TestSetAndGetOTP(t *testing.T) {
	mockRedis := mocks.NewMockRedis()

	testCases := []struct {
		name           string
		key            string
		otp            string
		expiration     time.Duration
		setError       error
		getError       error
		expectedOTP    string
		expectSetError bool
		expectGetError bool
	}{
		{
			name:           "Valid OTP",
			key:            "user:123:otp",
			otp:            "12345",
			expiration:     5 * time.Minute,
			expectedOTP:    "12345",
			expectSetError: false,
			expectGetError: false,
		},
		{
			name:           "Set error",
			key:            "user:456:otp",
			otp:            "67890",
			expiration:     5 * time.Minute,
			setError:       errors.New("connection error"),
			expectSetError: true,
			expectGetError: false,
		},
		{
			name:           "Get error",
			key:            "user:789:otp",
			otp:            "54321",
			expiration:     5 * time.Minute,
			getError:       errors.New("connection error"),
			expectSetError: false,
			expectGetError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock state
			mockRedis.Reset()

			// Set up mock to return specified errors
			if tc.setError != nil {
				mockRedis.SetSetError(tc.setError)
			}

			if tc.getError != nil {
				mockRedis.SetGetError(tc.getError)
			}

			// Test SetOTP
			err := mockRedis.SetOTP(context.Background(), tc.key, tc.otp, tc.expiration)

			if tc.expectSetError {
				if err == nil {
					t.Errorf("SetOTP() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("SetOTP() unexpected error: %v", err)
				return
			}

			// Test GetOTP
			retrievedOTP, err := mockRedis.GetOTP(context.Background(), tc.key)

			if tc.expectGetError {
				if err == nil {
					t.Errorf("GetOTP() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("GetOTP() unexpected error: %v", err)
				return
			}

			// Check retrieved OTP matches expected
			if retrievedOTP != tc.expectedOTP {
				t.Errorf("GetOTP() = %v, expected %v", retrievedOTP, tc.expectedOTP)
			}
		})
	}
}

func TestOTPExpiration(t *testing.T) {
	mockRedis := mocks.NewMockRedis()
	ctx := context.Background()

	key := "user:123:otp"
	otp := "12345"
	shortExpiration := 10 * time.Millisecond // Very short expiration for testing

	// Set OTP with short expiration
	err := mockRedis.SetOTP(ctx, key, otp, shortExpiration)
	if err != nil {
		t.Fatalf("SetOTP() unexpected error: %v", err)
	}

	// Verify OTP is retrievable immediately
	retrievedOTP, err := mockRedis.GetOTP(ctx, key)
	if err != nil {
		t.Fatalf("GetOTP() unexpected error immediately after set: %v", err)
	}
	if retrievedOTP != otp {
		t.Errorf("GetOTP() = %v, expected %v", retrievedOTP, otp)
	}

	// Wait for expiration
	time.Sleep(50 * time.Millisecond)

	// Verify OTP is no longer retrievable
	_, err = mockRedis.GetOTP(ctx, key)
	if err == nil {
		t.Error("GetOTP() expected error after expiration, got nil")
	}
}
