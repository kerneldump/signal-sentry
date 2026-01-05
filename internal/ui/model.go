package ui

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tmobile-stats/internal/config"
	"tmobile-stats/internal/gateway"
	"tmobile-stats/internal/logger"
	"tmobile-stats/internal/models"
	"tmobile-stats/internal/pinger"
)

// Msg types
type tickMsg time.Time
type dataMsg struct {
	Stats        *models.CombinedStats
	LifetimePing models.PingStats
	Err          error
}

// Model represents the state of the TUI.
type Model struct {
	cfg          *config.Config
	client       *http.Client
	pinger       *pinger.Pinger
	loggers      []logger.Logger
	buffer       []*models.CombinedStats
	lifetimePing models.PingStats
	interval     time.Duration
	width        int
	height       int
	showHelp     bool
	err          error
}

func NewModel(cfg *config.Config, client *http.Client, pg *pinger.Pinger, loggers []logger.Logger) *Model {
	return &Model{
		cfg:      cfg,
		client:   client,
		pinger:   pg,
		loggers:  loggers,
		interval: time.Duration(cfg.RefreshInterval) * time.Second,
		buffer:   make([]*models.CombinedStats, 0, 30),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchData(),
		m.tick(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "i":
			m.showHelp = !m.showHelp
		case "+", "=":
			m.interval += time.Second
			if m.interval > 60*time.Second {
				m.interval = 60 * time.Second
			}
		case "-":
			m.interval -= time.Second
			if m.interval < time.Second {
				m.interval = time.Second
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		return m, tea.Batch(
			m.fetchData(),
			m.tick(),
		)

	case dataMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.err = nil
			m.lifetimePing = msg.LifetimePing // Update lifetime stats
			
			// 1. Log data
			for _, l := range m.loggers {
				_ = l.Log(msg.Stats) // Best effort logging
			}

			// 2. Prepend to buffer
			m.buffer = append([]*models.CombinedStats{msg.Stats}, m.buffer...)
			if len(m.buffer) > 30 {
				m.buffer = m.buffer[:30]
			}
		}
	}

	return m, nil
}

func (m *Model) tick() tea.Cmd {
	return tea.Tick(m.interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) fetchData() tea.Cmd {
	return func() tea.Msg {
		gatewayData, err := gateway.FetchStats(m.client, m.cfg.RouterURL)
		if err != nil {
			return dataMsg{Err: err}
		}

		pingData := m.pinger.GetStatsAndReset()
		lifetimePing := m.pinger.GetLifetimeStats()

		return dataMsg{
			Stats: &models.CombinedStats{
				Gateway: *gatewayData,
				Ping:    pingData,
			},
			LifetimePing: lifetimePing,
		}
	}
}

// Rendering Logic (reusing the logic from main.go but adapted for Bubble Tea)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

const helpText = `SIGNAL METRICS GUIDE:
---------------------
* BAND:  The frequency band in use.
         n41: High speed, shorter range (Ultra Capacity).
         n25: Balanced speed and range.
         n71: Long range, slower speeds.

* RSRP:  (Reference Signal Received Power) Your main signal strength.
         Excellent > -80  | Good -80 to -95
         Fair -95 to -110 | Poor < -110 (Risk of drops).

* SINR:  (Signal-to-Interference-plus-Noise Ratio) Signal quality.
         Higher is better. > 20 is excellent.
         < 0 means high noise (your speed will suffer).

* RSRQ:  (Reference Signal Received Quality) The congestion indicator.
         If SINR is Good (high) but RSRQ is Bad (low), the tower is
         likely congested with heavy traffic.

* CID & gNBID:
         gNBID identifies the physical TOWER.
         CID identifies the SECTOR (which side of the tower you are facing).`

func (m *Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var s strings.Builder

	// 1. Device Info
	if len(m.buffer) > 0 {
		d := m.buffer[0].Gateway.Device
		s.WriteString(fmt.Sprintf("DEVICE: %s | FW: %s | Serial: %s\n", d.Model, d.SoftwareVersion, d.Serial))
	} else {
		s.WriteString("Waiting for data...\n")
	}

	// 2. Metrics Guide (Small version)
	s.WriteString("RSRP: Exc >-80, Good -95, Fair -110, Poor <-110 | SINR: Exc >20, Poor <0\n")
	
	// 3. Lifetime Ping Stats
	// PING: 531 packets transmitted, 531 packets received, 0.0% packet loss
	// round-trip min/avg/max/stddev = 20.986/49.955/855.485/53.432 ms
	lp := m.lifetimePing
	s.WriteString(fmt.Sprintf("PING: %d packets transmitted, %d packets received, %.1f%% packet loss\n", 
		lp.Sent, lp.Received, lp.Loss))
	s.WriteString(fmt.Sprintf("round-trip min/avg/max/stddev = %.3f/%.3f/%.3f/%.3f ms\n", 
		lp.Min, lp.Avg, lp.Max, lp.StdDev))

	s.WriteString(fmt.Sprintf("Interval: %v (Press +/- to adjust, i for info, q to quit)\n\n", m.interval))

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\n")
	}

	if m.showHelp {
		s.WriteString(helpText + "\n")
	} else {
		// 3. Header
		s.WriteString(headerStyle.Render(" TYPE  | BANDS      | BARS | RSRP  | SINR  | RSRQ | RSSI | CID   | TWR gNBID/PCIDE | MIN AVG MAX STD LOSS") + "\n")
		s.WriteString("-------+------------+------+-------+-------+------+------+-------+-----------------+-------------------------\n")

		// 4. Buffer
		// guideLines: Device(1), Metrics(1), PingStats(2), Interval(1), Empty(1), Header(1), Separator(1) = 8
		guideLines := 9 // Adjusted for 2 extra ping lines + safety
		linesUsed := 0
		maxLines := m.height - guideLines
		if maxLines < 0 {
			maxLines = 0
		}

		for _, data := range m.buffer {
			// 5G row
			row5g := m.renderRow("5G", data.Gateway.Signal.FiveG, data.Ping)
			if linesUsed < maxLines {
				s.WriteString(row5g)
				linesUsed++
			} else {
				break
			}

			// 4G row
			if len(data.Gateway.Signal.FourG.Bands) > 0 || data.Gateway.Signal.FourG.Bars > 0 {
				row4g := m.renderRow("4G", data.Gateway.Signal.FourG, data.Ping)
				if linesUsed < maxLines {
					s.WriteString(row4g)
					linesUsed++
				} else {
					break
				}
			}
		}
	}

	return s.String()
}

func (m *Model) renderRow(connType string, stats models.ConnectionStats, ping models.PingStats) string {
	bands := strings.Join(stats.Bands, ",")
	if bands == "" {
		bands = "---"
	}

	towerID := stats.GNBID
	if towerID == 0 {
		towerID = stats.PCID
	}

	// Colorize Connection Type
	typeStr := connType
	if connType == "5G" {
		typeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(fmt.Sprintf("%-2s", connType))
	} else if connType == "4G" {
		typeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render(fmt.Sprintf("%-2s", connType))
	}

	// Format loss string
	lossValStr := fmt.Sprintf("%.1f%%", ping.Loss)
	lossStr := lossValStr
	if ping.Loss > 0 {
		lossStr = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(lossValStr)
	}

	// Row Printf with explicit spaces to match header
	return fmt.Sprintf("  %s   | %-10s | %s  | %s  | %s  | %-4d | %-4d | %-5d | %-15d | %.1f %.1f %.1f %.1f %s \n",
		typeStr,
		bands,
		m.colorizeBars(stats.Bars),
		m.colorizeRSRP(stats.RSRP),
		m.colorizeSINR(stats.SINR),
		stats.RSRQ,
		stats.RSSI,
		stats.CID,
		towerID,
		ping.Min, ping.Avg, ping.Max, ping.StdDev, lossStr,
	)
}

// Helpers (Stripped down version of main.go colorize)
func (m *Model) colorizeRSRP(val int) string {
	s := fmt.Sprintf("%4d", val)
	if val > -80 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(s)
	} else if val >= -110 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render(s)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(s)
}

func (m *Model) colorizeSINR(val int) string {
	s := fmt.Sprintf("%4d", val)
	if val > 20 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(s)
	} else if val >= 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render(s)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(s)
}

func (m *Model) colorizeBars(val float64) string {
	s := fmt.Sprintf("%3.1f", val)
	if val >= 4.0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(s)
	} else if val >= 2.0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render(s)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(s)
}
