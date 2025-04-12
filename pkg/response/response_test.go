package response_test

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	testCases := []struct {
		name           string
		code           int
		errorMessage   string
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name:           "HTTP 400 Bad Request",
			code:           400,
			errorMessage:   "Bad Request",
			expectedCode:   400,
			expectedErrMsg: "Bad Request",
		},
		{
			name:           "HTTP 404 Not Found",
			code:           404,
			errorMessage:   "Not Found",
			expectedCode:   404,
			expectedErrMsg: "Not Found",
		},
		{
			name:           "HTTP 500 Internal Server Error",
			code:           500,
			errorMessage:   "Internal Server Error",
			expectedCode:   500,
			expectedErrMsg: "Internal Server Error",
		},
		{
			name:           "Custom error code",
			code:           422,
			errorMessage:   "Unprocessable Entity",
			expectedCode:   422,
			expectedErrMsg: "Unprocessable Entity",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewError(tc.code, tc.errorMessage)

			// Type assertion to access the Error struct fields
			respErr, ok := err.(*Error)
			if !ok {
				t.Fatalf("NewError() did not return *response.Error, got %T", err)
			}

			// Check the error code
			if respErr.Code != tc.expectedCode {
				t.Errorf("NewError() code = %v, expectedCode %v", respErr.Code, tc.expectedCode)
			}

			// Check the error message
			if respErr.Error() != tc.expectedErrMsg {
				t.Errorf("NewError() message = %v, expectedErrMsg %v", respErr.Error(), tc.expectedErrMsg)
			}
		})
	}
}

func TestErrorComparison(t *testing.T) {
	// Create some errors for comparison
	err1 := NewError(400, "Bad Request")
	err2 := NewError(400, "Bad Request")
	err3 := NewError(404, "Not Found")
	standardErr := errors.New("Bad Request")

	// Test using errors.Is
	if !errors.Is(err1, err2) {
		t.Error("errors.Is() failed to match identical response errors")
	}

	if errors.Is(err1, err3) {
		t.Error("errors.Is() incorrectly matched different response errors")
	}

	if errors.Is(err1, standardErr) {
		t.Error("errors.Is() incorrectly matched response error with standard error")
	}
}
