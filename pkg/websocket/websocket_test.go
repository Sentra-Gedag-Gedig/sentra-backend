package websocketPkg_test

import (
	"ProjectGolang/internal/api/detection"
	"ProjectGolang/internal/entity"
	"ProjectGolang/internal/mocks"
	"errors"
	"fmt"
	"testing"
)

func TestProcessFaceFrame(t *testing.T) {
	mockWS := mocks.NewMockWebSocket()

	// Define test frame data
	testFrame := []byte("mock face frame data")

	testCases := []struct {
		name       string
		connected  bool
		setError   error
		expectErr  bool
		customResp *entity.DetectionResult
	}{
		{
			name:      "Successfully process face frame",
			connected: true,
			expectErr: false,
		},
		{
			name:      "Not connected to face detection service",
			connected: false,
			expectErr: true,
		},
		{
			name:      "Error during processing",
			connected: true,
			setError:  errors.New("processing error"),
			expectErr: true,
		},
		{
			name:      "Custom response",
			connected: true,
			expectErr: false,
			customResp: &entity.DetectionResult{
				Status: "NO_FACE_DETECTED",
				FrameCenter: entity.Position{
					X: 400,
					Y: 300,
				},
				Instructions: []string{"Please position your face in the frame"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockWS.Reset()

			// Set connection status
			mockWS.SetConnected(detection.FaceDetection, tc.connected)

			// Set error if specified
			if tc.setError != nil {
				mockWS.SetFaceError(tc.setError)
			}

			// Set custom response if specified
			if tc.customResp != nil {
				mockWS.SetFaceFrameResponse(tc.customResp)
			}

			// Test processing face frame
			result, err := mockWS.ProcessFaceFrame(testFrame)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("ProcessFaceFrame() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("ProcessFaceFrame() unexpected error: %v", err)
				return
			}

			// Verify result is not nil
			if result == nil {
				t.Error("ProcessFaceFrame() returned nil result")
				return
			}

			// Check custom response if specified
			if tc.customResp != nil {
				if result.Status != tc.customResp.Status {
					t.Errorf("ProcessFaceFrame() status = %v, expected %v",
						result.Status, tc.customResp.Status)
				}

				// Compare frames center
				if result.FrameCenter.X != tc.customResp.FrameCenter.X ||
					result.FrameCenter.Y != tc.customResp.FrameCenter.Y {
					t.Errorf("ProcessFaceFrame() frame center = %+v, expected %+v",
						result.FrameCenter, tc.customResp.FrameCenter)
				}

				// Compare instructions if present
				if len(tc.customResp.Instructions) > 0 {
					if len(result.Instructions) != len(tc.customResp.Instructions) {
						t.Errorf("ProcessFaceFrame() instructions count = %d, expected %d",
							len(result.Instructions), len(tc.customResp.Instructions))
					} else {
						for i, instruction := range tc.customResp.Instructions {
							if result.Instructions[i] != instruction {
								t.Errorf("ProcessFaceFrame() instruction[%d] = %s, expected %s",
									i, result.Instructions[i], instruction)
							}
						}
					}
				}
			}

			// Verify frame count increased
			count := mockWS.GetProcessedFaceFramesCount()
			if count != 1 {
				t.Errorf("Expected 1 face frame to be processed, got %d", count)
			}
		})
	}
}

func TestProcessKTPFrame(t *testing.T) {
	mockWS := mocks.NewMockWebSocket()

	// Define test frame data
	testFrame := []byte("mock KTP frame data")

	testCases := []struct {
		name       string
		connected  bool
		setError   error
		expectErr  bool
		customResp *entity.KTPDetectionResult
	}{
		{
			name:      "Successfully process KTP frame",
			connected: true,
			expectErr: false,
		},
		{
			name:      "Not connected to KTP detection service",
			connected: false,
			expectErr: true,
		},
		{
			name:      "Error during processing",
			connected: true,
			setError:  errors.New("processing error"),
			expectErr: true,
		},
		{
			name:      "Custom response",
			connected: true,
			expectErr: false,
			customResp: &entity.KTPDetectionResult{
				Message:    "KTP not detected",
				Confidence: 0.3,
				BBox:       []float64{0, 0, 0, 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockWS.Reset()

			// Set connection status
			mockWS.SetConnected(detection.KTPDetection, tc.connected)

			// Set error if specified
			if tc.setError != nil {
				mockWS.SetKTPError(tc.setError)
			}

			// Set custom response if specified
			if tc.customResp != nil {
				mockWS.SetKTPFrameResponse(tc.customResp)
			}

			// Test processing KTP frame
			result, err := mockWS.ProcessKTPFrame(testFrame)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("ProcessKTPFrame() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("ProcessKTPFrame() unexpected error: %v", err)
				return
			}

			// Verify result is not nil
			if result == nil {
				t.Error("ProcessKTPFrame() returned nil result")
				return
			}

			// Check custom response if specified
			if tc.customResp != nil {
				if result.Message != tc.customResp.Message {
					t.Errorf("ProcessKTPFrame() message = %v, expected %v",
						result.Message, tc.customResp.Message)
				}

				if result.Confidence != tc.customResp.Confidence {
					t.Errorf("ProcessKTPFrame() confidence = %v, expected %v",
						result.Confidence, tc.customResp.Confidence)
				}

				// Compare BBox if present
				if len(tc.customResp.BBox) > 0 {
					if len(result.BBox) != len(tc.customResp.BBox) {
						t.Errorf("ProcessKTPFrame() BBox size = %d, expected %d",
							len(result.BBox), len(tc.customResp.BBox))
					} else {
						for i, val := range tc.customResp.BBox {
							if result.BBox[i] != val {
								t.Errorf("ProcessKTPFrame() BBox[%d] = %f, expected %f",
									i, result.BBox[i], val)
							}
						}
					}
				}
			}

			// Verify frame count increased
			count := mockWS.GetProcessedKTPFramesCount()
			if count != 1 {
				t.Errorf("Expected 1 KTP frame to be processed, got %d", count)
			}
		})
	}
}

func TestProcessQRISFrame(t *testing.T) {
	mockWS := mocks.NewMockWebSocket()

	// Define test frame data
	testFrame := []byte("mock QRIS frame data")

	testCases := []struct {
		name       string
		connected  bool
		setError   error
		expectErr  bool
		customResp *entity.QRISDetectionResult
	}{
		{
			name:      "Successfully process QRIS frame",
			connected: true,
			expectErr: false,
		},
		{
			name:      "Not connected to QRIS detection service",
			connected: false,
			expectErr: true,
		},
		{
			name:      "Error during processing",
			connected: true,
			setError:  errors.New("processing error"),
			expectErr: true,
		},
		{
			name:      "Custom response",
			connected: true,
			expectErr: false,
			customResp: &entity.QRISDetectionResult{
				Message:    "QRIS not detected",
				Confidence: 0.2,
				BBox:       []float64{0, 0, 0, 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockWS.Reset()

			// Set connection status
			mockWS.SetConnected(detection.QRISDetection, tc.connected)

			// Set error if specified
			if tc.setError != nil {
				mockWS.SetQRISError(tc.setError)
			}

			// Set custom response if specified
			if tc.customResp != nil {
				mockWS.SetQRISFrameResponse(tc.customResp)
			}

			// Test processing QRIS frame
			result, err := mockWS.ProcessQRISFrame(testFrame)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("ProcessQRISFrame() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("ProcessQRISFrame() unexpected error: %v", err)
				return
			}

			// Verify result is not nil
			if result == nil {
				t.Error("ProcessQRISFrame() returned nil result")
				return
			}

			// Check custom response if specified
			if tc.customResp != nil {
				if result.Message != tc.customResp.Message {
					t.Errorf("ProcessQRISFrame() message = %v, expected %v",
						result.Message, tc.customResp.Message)
				}

				if result.Confidence != tc.customResp.Confidence {
					t.Errorf("ProcessQRISFrame() confidence = %v, expected %v",
						result.Confidence, tc.customResp.Confidence)
				}

				// Compare BBox if present
				if len(tc.customResp.BBox) > 0 {
					if len(result.BBox) != len(tc.customResp.BBox) {
						t.Errorf("ProcessQRISFrame() BBox size = %d, expected %d",
							len(result.BBox), len(tc.customResp.BBox))
					} else {
						for i, val := range tc.customResp.BBox {
							if result.BBox[i] != val {
								t.Errorf("ProcessQRISFrame() BBox[%d] = %f, expected %f",
									i, result.BBox[i], val)
							}
						}
					}
				}
			}

			// Verify frame count increased
			count := mockWS.GetProcessedQRISFramesCount()
			if count != 1 {
				t.Errorf("Expected 1 QRIS frame to be processed, got %d", count)
			}
		})
	}
}

func TestReconnect(t *testing.T) {
	mockWS := mocks.NewMockWebSocket()

	testCases := []struct {
		name          string
		detectionType detection.DetectionType
		setError      error
		expectErr     bool
	}{
		{
			name:          "Successfully reconnect to Face detection",
			detectionType: detection.FaceDetection,
			expectErr:     false,
		},
		{
			name:          "Successfully reconnect to KTP detection",
			detectionType: detection.KTPDetection,
			expectErr:     false,
		},
		{
			name:          "Successfully reconnect to QRIS detection",
			detectionType: detection.QRISDetection,
			expectErr:     false,
		},
		{
			name:          "Error during reconnection",
			detectionType: detection.FaceDetection,
			setError:      errors.New("reconnection error"),
			expectErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockWS.Reset()

			// Set error if specified
			if tc.setError != nil {
				mockWS.SetReconnectError(tc.setError)
			}

			// Test reconnection
			err := mockWS.Reconnect(tc.detectionType)

			// Check error status
			if tc.expectErr {
				if err == nil {
					t.Errorf("Reconnect() expected error, got nil")
				}
				return
			} else if err != nil {
				t.Errorf("Reconnect() unexpected error: %v", err)
				return
			}

			// Verify connected status
			if !mockWS.IsConnected(tc.detectionType) {
				t.Errorf("Reconnect() did not set connected status for %v", tc.detectionType)
			}
		})
	}
}

func TestIsConnected(t *testing.T) {
	mockWS := mocks.NewMockWebSocket()

	testCases := []struct {
		name          string
		detectionType detection.DetectionType
		connected     bool
		expected      bool
	}{
		{
			name:          "Connected to Face detection",
			detectionType: detection.FaceDetection,
			connected:     true,
			expected:      true,
		},
		{
			name:          "Not connected to KTP detection",
			detectionType: detection.KTPDetection,
			connected:     false,
			expected:      false,
		},
		{
			name:          "Connection status not set",
			detectionType: detection.QRISDetection,
			// Don't set connected status
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mockWS.Reset()

			// Set connection status if it should be set
			if tc.name != "Connection status not set" {
				mockWS.SetConnected(tc.detectionType, tc.connected)
			}

			// Test IsConnected
			connected := mockWS.IsConnected(tc.detectionType)

			// Check result
			if connected != tc.expected {
				t.Errorf("IsConnected(%v) = %v, expected %v",
					tc.detectionType, connected, tc.expected)
			}
		})
	}
}

func TestCloseConnections(t *testing.T) {
	mockWS := mocks.NewMockWebSocket()

	// Set up connected status for all detection types
	mockWS.SetConnected(detection.FaceDetection, true)
	mockWS.SetConnected(detection.KTPDetection, true)
	mockWS.SetConnected(detection.QRISDetection, true)

	// Verify they are all connected
	if !mockWS.IsConnected(detection.FaceDetection) ||
		!mockWS.IsConnected(detection.KTPDetection) ||
		!mockWS.IsConnected(detection.QRISDetection) {
		t.Error("Failed to set up test - not all detection types are connected")
		return
	}

	// Test closing connections
	mockWS.CloseConnections()

	// Verify all connections are closed
	if mockWS.IsConnected(detection.FaceDetection) {
		t.Error("CloseConnections() did not close Face detection connection")
	}

	if mockWS.IsConnected(detection.KTPDetection) {
		t.Error("CloseConnections() did not close KTP detection connection")
	}

	if mockWS.IsConnected(detection.QRISDetection) {
		t.Error("CloseConnections() did not close QRIS detection connection")
	}
}

func TestMultipleFrameProcessing(t *testing.T) {
	mockWS := mocks.NewMockWebSocket()

	// Set connection status for all types
	mockWS.SetConnected(detection.FaceDetection, true)
	mockWS.SetConnected(detection.KTPDetection, true)
	mockWS.SetConnected(detection.QRISDetection, true)

	// Process multiple frames of each type
	frameCount := 5
	for i := 0; i < frameCount; i++ {
		frameData := []byte(fmt.Sprintf("frame_%d", i))

		_, err := mockWS.ProcessFaceFrame(frameData)
		if err != nil {
			t.Errorf("ProcessFaceFrame() unexpected error: %v", err)
		}

		_, err = mockWS.ProcessKTPFrame(frameData)
		if err != nil {
			t.Errorf("ProcessKTPFrame() unexpected error: %v", err)
		}

		_, err = mockWS.ProcessQRISFrame(frameData)
		if err != nil {
			t.Errorf("ProcessQRISFrame() unexpected error: %v", err)
		}
	}

	// Verify frame counts
	faceCount := mockWS.GetProcessedFaceFramesCount()
	if faceCount != frameCount {
		t.Errorf("Expected %d face frames to be processed, got %d", frameCount, faceCount)
	}

	ktpCount := mockWS.GetProcessedKTPFramesCount()
	if ktpCount != frameCount {
		t.Errorf("Expected %d KTP frames to be processed, got %d", frameCount, ktpCount)
	}

	qrisCount := mockWS.GetProcessedQRISFramesCount()
	if qrisCount != frameCount {
		t.Errorf("Expected %d QRIS frames to be processed, got %d", frameCount, qrisCount)
	}
}
