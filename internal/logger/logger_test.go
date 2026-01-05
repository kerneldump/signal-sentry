package logger

import (
	"testing"
	"tmobile-stats/internal/models"
)

// MockLogger implements Logger for testing
type MockLogger struct {
	LastEntry *models.CombinedStats
}

func (m *MockLogger) Log(data *models.CombinedStats) error {
	m.LastEntry = data
	return nil
}


func (m *MockLogger) Close() error {
	return nil
}

func TestLoggerInterface(t *testing.T) {
	// This test simply verifies that MockLogger satisfies the Logger interface
	// which doesn't exist yet, so it should fail to compile.
	var _ Logger = &MockLogger{}
}
