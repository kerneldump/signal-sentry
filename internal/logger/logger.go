package logger

import "tmobile-stats/internal/models"

// Logger defines the interface for logging gateway statistics.
type Logger interface {
	Log(data *models.CombinedStats) error
	Close() error
}

