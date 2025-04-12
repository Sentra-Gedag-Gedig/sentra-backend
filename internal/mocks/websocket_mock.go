package mocks

import (
	"ProjectGolang/internal/api/detection"
	"ProjectGolang/internal/entity"
	"encoding/json"
	"errors"
	"sync"
)

// MockWebSocket implements the websocketPkg.IWebsocket interface for testing
type MockWebSocket struct {
	mutex               sync.RWMutex
	connected           map[detection.DetectionType]bool
	faceFrameResponse   *entity.DetectionResult
	ktpFrameResponse    *entity.KTPDetectionResult
	qrisFrameResponse   *entity.QRISDetectionResult
	faceError           error
	ktpError            error
	qrisError           error
	reconnectError      error
	processedFaceFrames [][]byte
	processedKTPFrames  [][]byte
	processedQRISFrames [][]byte
}

// NewMockWebSocket creates a new instance of MockWebSocket
func NewMockWebSocket() *MockWebSocket {
	return &MockWebSocket{
		connected:           make(map[detection.DetectionType]bool),
		processedFaceFrames: make([][]byte, 0),
		processedKTPFrames:  make([][]byte, 0),
		processedQRISFrames: make([][]byte, 0),
		// Set default connected status for all types
		faceFrameResponse: &entity.DetectionResult{
			Status:   "PERFECT_POSITION",
			FaceSize: ptr(0.8),
			FacePosition: &entity.Position{
				X: 320,
				Y: 240,
			},
			FrameCenter: entity.Position{
				X: 320,
				Y: 240,
			},
			Deviations: map[string]float64{
				"x": 0.0,
				"y": 0.0,
			},
		},
		ktpFrameResponse: &entity.KTPDetectionResult{
			Message:    "KTP detected",
			Confidence: 0.95,
			BBox:       []float64{10, 20, 200, 300},
			KTPPosition: &entity.KTPPosition{
				X1: 10,
				Y1: 20,
				X2: 200,
				Y2: 300,
			},
		},
		qrisFrameResponse: &entity.QRISDetectionResult{
			Message:    "QRIS detected",
			Confidence: 0.98,
			BBox:       []float64{50, 60, 250, 350},
			QRISPosition: &entity.QRISPosition{
				X1: 50,
				Y1: 60,
				X2: 250,
				Y2: 350,
			},
		},
	}
}

// ProcessFaceFrame implements the ProcessFaceFrame method of websocketPkg.IWebsocket
func (m *MockWebSocket) ProcessFaceFrame(frame []byte) (*entity.DetectionResult, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.faceError != nil {
		return nil, m.faceError
	}

	if !m.IsConnected(detection.FaceDetection) {
		return nil, errors.New("not connected to face detection service")
	}

	// Store the frame for later verification
	m.processedFaceFrames = append(m.processedFaceFrames, frame)

	return m.faceFrameResponse, nil
}

// ProcessKTPFrame implements the ProcessKTPFrame method of websocketPkg.IWebsocket
func (m *MockWebSocket) ProcessKTPFrame(frame []byte) (*entity.KTPDetectionResult, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.ktpError != nil {
		return nil, m.ktpError
	}

	if !m.IsConnected(detection.KTPDetection) {
		return nil, errors.New("not connected to KTP detection service")
	}

	// Store the frame for later verification
	m.processedKTPFrames = append(m.processedKTPFrames, frame)

	return m.ktpFrameResponse, nil
}

// ProcessQRISFrame implements the ProcessQRISFrame method of websocketPkg.IWebsocket
func (m *MockWebSocket) ProcessQRISFrame(frame []byte) (*entity.QRISDetectionResult, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.qrisError != nil {
		return nil, m.qrisError
	}

	if !m.IsConnected(detection.QRISDetection) {
		return nil, errors.New("not connected to QRIS detection service")
	}

	// Store the frame for later verification
	m.processedQRISFrames = append(m.processedQRISFrames, frame)

	return m.qrisFrameResponse, nil
}

// IsConnected implements the IsConnected method of websocketPkg.IWebsocket
func (m *MockWebSocket) IsConnected(detectionType detection.DetectionType) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	connected, exists := m.connected[detectionType]
	return exists && connected
}

// Reconnect implements the Reconnect method of websocketPkg.IWebsocket
func (m *MockWebSocket) Reconnect(detectionType detection.DetectionType) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.reconnectError != nil {
		return m.reconnectError
	}

	m.connected[detectionType] = true
	return nil
}

// CloseConnections implements the CloseConnections method of websocketPkg.IWebsocket
func (m *MockWebSocket) CloseConnections() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key := range m.connected {
		m.connected[key] = false
	}
}

// SetConnected sets the connection status for a detection type
func (m *MockWebSocket) SetConnected(detectionType detection.DetectionType, connected bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.connected[detectionType] = connected
}

// SetFaceError sets an error to be returned by ProcessFaceFrame
func (m *MockWebSocket) SetFaceError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.faceError = err
}

// SetKTPError sets an error to be returned by ProcessKTPFrame
func (m *MockWebSocket) SetKTPError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.ktpError = err
}

// SetQRISError sets an error to be returned by ProcessQRISFrame
func (m *MockWebSocket) SetQRISError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.qrisError = err
}

// SetReconnectError sets an error to be returned by Reconnect
func (m *MockWebSocket) SetReconnectError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.reconnectError = err
}

// SetFaceFrameResponse sets the response to be returned by ProcessFaceFrame
func (m *MockWebSocket) SetFaceFrameResponse(response *entity.DetectionResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.faceFrameResponse = response
}

// SetKTPFrameResponse sets the response to be returned by ProcessKTPFrame
func (m *MockWebSocket) SetKTPFrameResponse(response *entity.KTPDetectionResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.ktpFrameResponse = response
}

// SetQRISFrameResponse sets the response to be returned by ProcessQRISFrame
func (m *MockWebSocket) SetQRISFrameResponse(response *entity.QRISDetectionResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.qrisFrameResponse = response
}

// GetProcessedFaceFramesCount returns the number of face frames processed
func (m *MockWebSocket) GetProcessedFaceFramesCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.processedFaceFrames)
}

// GetProcessedKTPFramesCount returns the number of KTP frames processed
func (m *MockWebSocket) GetProcessedKTPFramesCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.processedKTPFrames)
}

// GetProcessedQRISFramesCount returns the number of QRIS frames processed
func (m *MockWebSocket) GetProcessedQRISFramesCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.processedQRISFrames)
}

// Reset clears all stored data and errors
func (m *MockWebSocket) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.processedFaceFrames = make([][]byte, 0)
	m.processedKTPFrames = make([][]byte, 0)
	m.processedQRISFrames = make([][]byte, 0)
	m.faceError = nil
	m.ktpError = nil
	m.qrisError = nil
	m.reconnectError = nil

	// Set default connected status
	m.connected = make(map[detection.DetectionType]bool)

	// Reset to default responses
	m.faceFrameResponse = &entity.DetectionResult{
		Status:   "PERFECT_POSITION",
		FaceSize: ptr(0.8),
		FacePosition: &entity.Position{
			X: 320,
			Y: 240,
		},
		FrameCenter: entity.Position{
			X: 320,
			Y: 240,
		},
		Deviations: map[string]float64{
			"x": 0.0,
			"y": 0.0,
		},
	}

	m.ktpFrameResponse = &entity.KTPDetectionResult{
		Message:    "KTP detected",
		Confidence: 0.95,
		BBox:       []float64{10, 20, 200, 300},
		KTPPosition: &entity.KTPPosition{
			X1: 10,
			Y1: 20,
			X2: 200,
			Y2: 300,
		},
	}

	m.qrisFrameResponse = &entity.QRISDetectionResult{
		Message:    "QRIS detected",
		Confidence: 0.98,
		BBox:       []float64{50, 60, 250, 350},
		QRISPosition: &entity.QRISPosition{
			X1: 50,
			Y1: 60,
			X2: 250,
			Y2: 350,
		},
	}
}

// SerializeLastFaceFrame returns the JSON representation of the last processed face frame
func (m *MockWebSocket) SerializeLastFaceFrame() (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.processedFaceFrames) == 0 {
		return "", errors.New("no face frames processed")
	}

	// In a real test you might want to deserialize the frame data or perform other checks
	return string(m.processedFaceFrames[len(m.processedFaceFrames)-1]), nil
}

// SerializeLastKTPFrame returns the JSON representation of the last processed KTP frame
func (m *MockWebSocket) SerializeLastKTPFrame() (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.processedKTPFrames) == 0 {
		return "", errors.New("no KTP frames processed")
	}

	return string(m.processedKTPFrames[len(m.processedKTPFrames)-1]), nil
}

// SerializeLastQRISFrame returns the JSON representation of the last processed QRIS frame
func (m *MockWebSocket) SerializeLastQRISFrame() (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.processedQRISFrames) == 0 {
		return "", errors.New("no QRIS frames processed")
	}

	return string(m.processedQRISFrames[len(m.processedQRISFrames)-1]), nil
}

// Helper function to create float64 pointer
func ptr(v float64) *float64 {
	return &v
}

// ToJSON converts a struct to JSON string (for testing)
func ToJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}
