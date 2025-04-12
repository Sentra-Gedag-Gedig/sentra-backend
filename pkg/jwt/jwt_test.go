package jwtPkg

import (
	"os"
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	// Set up environment variable for testing
	originalValue := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	defer os.Setenv("JWT_ACCESS_TOKEN_SECRET", originalValue)

	os.Setenv("JWT_ACCESS_TOKEN_SECRET", "test_secret_key")

	testCases := []struct {
		name    string
		data    map[string]interface{}
		expiry  time.Duration
		wantErr bool
	}{
		{
			name: "Valid data",
			data: map[string]interface{}{
				"id":       "user123",
				"email":    "test@example.com",
				"username": "testuser",
			},
			expiry:  time.Hour,
			wantErr: false,
		},
		{
			name:    "Empty data",
			data:    map[string]interface{}{},
			expiry:  time.Hour,
			wantErr: false,
		},
		{
			name: "Zero expiry",
			data: map[string]interface{}{
				"id": "user123",
			},
			expiry:  0,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, exp, err := Sign(tc.data, tc.expiry)

			// Check error status matches expectations
			if (err != nil) != tc.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// If we don't expect an error, do further validation
			if !tc.wantErr {
				if token == "" {
					t.Errorf("Sign() returned empty token")
				}

				expectedExpiry := time.Now().Add(tc.expiry).Unix()
				// Allow for a small time difference (5 seconds) in test execution
				if exp < expectedExpiry-5 || exp > expectedExpiry+5 {
					t.Errorf("Sign() expiry time = %v, expected around %v", exp, expectedExpiry)
				}
			}
		})
	}
}

func TestSignWithMissingSecret(t *testing.T) {
	// Save original value and restore it after the test
	originalValue := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	defer os.Setenv("JWT_ACCESS_TOKEN_SECRET", originalValue)

	// Unset the environment variable for this test
	os.Unsetenv("JWT_ACCESS_TOKEN_SECRET")

	data := map[string]interface{}{
		"id": "user123",
	}

	_, _, err := Sign(data, time.Hour)

	if err == nil {
		t.Errorf("Sign() with missing secret should return error, got nil")
	}
}

// Note: Testing VerifyTokenHeader and GetUserLoginData would require mocking fiber.Ctx
// which is beyond the scope of this basic test suite. Those tests would be better
// placed in an integration test suite.
