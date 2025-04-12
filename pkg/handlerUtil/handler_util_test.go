package handlerUtil

import (
	"ProjectGolang/pkg/response"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"testing"
)

// Mock Fiber context for testing
type mockFiberCtx struct {
	statusCode int
	body       []byte
	headers    map[string]string
	path       string
}

func (m *mockFiberCtx) Status(status int) *fiber.Ctx {
	m.statusCode = status
	return &fiber.Ctx{}
}

func (m *mockFiberCtx) JSON(data interface{}) error {
	// In a real test we'd serialize the JSON, but for this test
	// we're simplifying by just checking it's not nil
	return nil
}

func (m *mockFiberCtx) SendStatus(status int) error {
	m.statusCode = status
	return nil
}

func (m *mockFiberCtx) Path() string {
	return m.path
}

// Setup test function that creates a mock context and error handler
func setupTest() (*mockFiberCtx, *ErrorHandler) {
	mockCtx := &mockFiberCtx{
		headers: make(map[string]string),
		path:    "/test/path",
	}

	logger := logrus.New()
	// Disable logger output during tests
	logger.Out = &strings.Builder{}

	errorHandler := New(logger)

	return mockCtx, errorHandler
}

func TestHandleSuccess(t *testing.T) {
	mockCtx, errHandler := setupTest()

	testCases := []struct {
		name       string
		statusCode int
		data       interface{}
		expectData bool
	}{
		{
			name:       "Success with data",
			statusCode: http.StatusOK,
			data:       map[string]string{"key": "value"},
			expectData: true,
		},
		{
			name:       "Success without data",
			statusCode: http.StatusNoContent,
			data:       nil,
			expectData: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := errHandler.HandleSuccess(mockCtx, tc.statusCode, tc.data)

			if err != nil {
				t.Errorf("HandleSuccess() returned unexpected error: %v", err)
			}

			if mockCtx.statusCode != tc.statusCode {
				t.Errorf("HandleSuccess() set status code %v, expected %v", mockCtx.statusCode, tc.statusCode)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	mockCtx, errHandler := setupTest()

	testCases := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "Response error",
			err:            response.NewError(http.StatusBadRequest, "Bad Request"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Standard error",
			err:            errors.New("Standard error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestID := "test-request-id"
			operation := "test-operation"

			err := errHandler.Handle(mockCtx, requestID, tc.err, mockCtx.Path(), operation)

			if err != nil {
				t.Errorf("Handle() returned unexpected error: %v", err)
			}

			if mockCtx.statusCode != tc.expectedStatus {
				t.Errorf("Handle() set status code %v, expected %v", mockCtx.statusCode, tc.expectedStatus)
			}
		})
	}
}

func TestHandleValidationError(t *testing.T) {
	mockCtx, errHandler := setupTest()

	requestID := "test-request-id"
	validationErr := errors.New("Validation failed")

	err := errHandler.HandleValidationError(mockCtx, requestID, validationErr, mockCtx.Path())

	if err != nil {
		t.Errorf("HandleValidationError() returned unexpected error: %v", err)
	}

	if mockCtx.statusCode != http.StatusBadRequest {
		t.Errorf("HandleValidationError() set status code %v, expected %v", mockCtx.statusCode, http.StatusBadRequest)
	}
}

func TestHandleRequestTimeout(t *testing.T) {
	mockCtx, errHandler := setupTest()

	err := errHandler.HandleRequestTimeout(mockCtx)

	if err != nil {
		t.Errorf("HandleRequestTimeout() returned unexpected error: %v", err)
	}

	if mockCtx.statusCode != http.StatusRequestTimeout {
		t.Errorf("HandleRequestTimeout() set status code %v, expected %v", mockCtx.statusCode, http.StatusRequestTimeout)
	}
}

func TestHandleUnauthorized(t *testing.T) {
	mockCtx, errHandler := setupTest()

	requestID := "test-request-id"
	message := "Unauthorized access"

	err := errHandler.HandleUnauthorized(mockCtx, requestID, message)

	if err != nil {
		t.Errorf("HandleUnauthorized() returned unexpected error: %v", err)
	}

	if mockCtx.statusCode != http.StatusUnauthorized {
		t.Errorf("HandleUnauthorized() set status code %v, expected %v", mockCtx.statusCode, http.StatusUnauthorized)
	}
}
