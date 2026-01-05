package pinger

import (
	"context"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"tmobile-stats/internal/models"
)

// Pinger manages the background ping process.
type Pinger struct {
	Target   string
	Interval time.Duration
	stats    models.PingStats
	m2       float64 // for running variance calculation
	mu       sync.RWMutex
}

// NewPinger creates a new Pinger instance.
func NewPinger(target string, interval time.Duration) *Pinger {
	return &Pinger{
		Target:   target,
		Interval: interval,
	}
}

// Run starts the ping loop.
func (p *Pinger) Run(ctx context.Context) {
	// First ping immediately
	p.ping()

	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.ping()
		}
	}
}

// GetStats returns a copy of the current statistics.
func (p *Pinger) GetStats() models.PingStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stats
}


var (
	// macOS ping output example for stats:
	// round-trip min/avg/max/stddev = 14.545/14.545/14.545/0.000 ms
	macStatsRegex = regexp.MustCompile(`min/avg/max/stddev = ([\d.]+)/([\d.]+)/([\d.]+)/([\d.]+) ms`)
)

func (p *Pinger) ping() {
	ctx, cancel := context.WithTimeout(context.Background(), p.Interval)
	defer cancel()

	// Use -t 1 to set timeout to 1 second on macOS (optional, but good)
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-t", "1", p.Target)
	out, err := cmd.Output()

	p.mu.Lock()
	defer p.mu.Unlock()

	p.stats.Sent++

	if err != nil {
		p.stats.Loss = float64(p.stats.Sent-p.stats.Received) / float64(p.stats.Sent) * 100
		return
	}

	p.stats.Received++
	p.stats.Loss = float64(p.stats.Sent-p.stats.Received) / float64(p.stats.Sent) * 100

	matches := macStatsRegex.FindStringSubmatch(string(out))
	if len(matches) == 5 {
		rtt, _ := strconv.ParseFloat(matches[2], 64) // Use avg from ping output as the RTT
		p.stats.LastRTT = rtt

		if p.stats.Received == 1 {
			p.stats.Min = rtt
			p.stats.Max = rtt
			p.stats.Avg = rtt
			p.stats.StdDev = 0
			p.m2 = 0
		} else {
			if rtt < p.stats.Min {
				p.stats.Min = rtt
			}
			if rtt > p.stats.Max {
				p.stats.Max = rtt
			}

			// Welford's algorithm for running mean and variance
			delta := rtt - p.stats.Avg
			p.stats.Avg += delta / float64(p.stats.Received)
			delta2 := rtt - p.stats.Avg
			p.m2 += delta * delta2

			variance := p.m2 / float64(p.stats.Received)
			p.stats.StdDev = math.Sqrt(variance)
		}
	}
}
