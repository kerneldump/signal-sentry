package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"tmobile-stats/internal/gateway"
)

type CSVLogger struct {
	file   *os.File
	writer *csv.Writer
}

func NewCSVLogger(filename string) (*CSVLogger, error) {
	fileInfo, err := os.Stat(filename)
	isNew := os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %w", err)
	}

	writer := csv.NewWriter(f)
	l := &CSVLogger{file: f, writer: writer}

	if isNew {
		header := []string{
			"Timestamp",
			"5G_Band", "5G_RSRP", "5G_SINR", "5G_Bars",
			"4G_Band", "4G_RSRP", "4G_SINR", "4G_Bars",
		}
		if err := writer.Write(header); err != nil {
			return nil, fmt.Errorf("could not write CSV header: %w", err)
		}
		writer.Flush()
	}

	return l, nil
}

func (l *CSVLogger) Log(data *gateway.GatewayResponse) error {
	row := []string{
		time.Now().Format(time.RFC3339),
		strings.Join(data.Signal.FiveG.Bands, ","),
		strconv.Itoa(data.Signal.FiveG.RSRP),
		strconv.Itoa(data.Signal.FiveG.SINR),
		strconv.FormatFloat(data.Signal.FiveG.Bars, 'f', 1, 64),
		strings.Join(data.Signal.FourG.Bands, ","),
		strconv.Itoa(data.Signal.FourG.RSRP),
		strconv.Itoa(data.Signal.FourG.SINR),
		strconv.FormatFloat(data.Signal.FourG.Bars, 'f', 1, 64),
	}

	if err := l.writer.Write(row); err != nil {
		return fmt.Errorf("could not write CSV row: %w", err)
	}
	l.writer.Flush()
	return nil
}

func (l *CSVLogger) Close() error {
	l.writer.Flush()
	return l.file.Close()
}
