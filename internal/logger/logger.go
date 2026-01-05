package logger

import "tmobile-stats/internal/gateway"

// Logger defines the interface for logging gateway statistics.
type Logger interface {
	Log(data *gateway.GatewayResponse) error
	Close() error
}
