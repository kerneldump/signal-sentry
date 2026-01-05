package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"tmobile-stats/internal/models"
)

type JSONLogger struct {
	file *os.File
}

func NewJSONLogger(filename string) (*JSONLogger, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %w", err)
	}
	return &JSONLogger{file: f}, nil
}

func (l *JSONLogger) Log(data *models.CombinedStats) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %w", err)
	}

	_, err = l.file.Write(append(bytes, '\n'))
	if err != nil {
		return fmt.Errorf("could not write to log file: %w", err)
	}

	return nil
}


func (l *JSONLogger) Close() error {
	return l.file.Close()
}
