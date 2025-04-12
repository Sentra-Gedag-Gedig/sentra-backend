package bcrypt

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	bcryptService := New()

	testCases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "secure_password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "Long password",
			password: "very_long_password_that_is_more_than_72_bytes_which_is_the_limit_for_bcrypt_very_long_password_that_is_more",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := bcryptService.HashPassword(tc.password)

			// Check error status matches expectations
			if (err != nil) != tc.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// If we don't expect an error, do further validation
			if !tc.wantErr {
				if hashedPassword == "" {
					t.Errorf("HashPassword() returned empty string, expected non-empty hash")
				}

				if hashedPassword == tc.password {
					t.Errorf("HashPassword() returned original password, expected hashed value")
				}
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	// Initialize the bcrypt service
	bcryptService := New()

	// Hash a known password
	password := "test_password"
	hashedPassword, err := bcryptService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	testCases := []struct {
		name        string
		hashedPwd   string
		providedPwd string
		expectMatch bool
	}{
		{
			name:        "Matching password",
			hashedPwd:   hashedPassword,
			providedPwd: password,
			expectMatch: true,
		},
		{
			name:        "Non-matching password",
			hashedPwd:   hashedPassword,
			providedPwd: "wrong_password",
			expectMatch: false,
		},
		{
			name:        "Empty provided password",
			hashedPwd:   hashedPassword,
			providedPwd: "",
			expectMatch: false,
		},
		{
			name:        "Invalid hash format",
			hashedPwd:   "not_a_valid_bcrypt_hash",
			providedPwd: password,
			expectMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bcryptService.ComparePassword(tc.hashedPwd, tc.providedPwd)

			// Check if result matches expectation
			if (err == nil) != tc.expectMatch {
				t.Errorf("ComparePassword() error = %v, expectMatch %v", err, tc.expectMatch)
			}
		})
	}
}
