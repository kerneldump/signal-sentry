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
	stats    models.PingStats // Reset every interval
	lifetime models.PingStats // Cumulative for session
	m2       float64          // for interval variance
	lifeM2   float64          // for lifetime variance
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

// GetLifetimeStats returns the cumulative statistics for the session.
func (p *Pinger) GetLifetimeStats() models.PingStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lifetime
}

// GetStatsAndReset returns the current statistics and resets the internal counters.
// This is useful for interval-based reporting (e.g. "stats for the last 10 seconds").
func (p *Pinger) GetStatsAndReset() models.PingStats {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Capture current state
	currentStats := p.stats

	// Reset state for next window
	p.stats = models.PingStats{
		Min:     0,
		Max:     0,
		Avg:     0,
		StdDev:  0,
		Loss:    0,
		Sent:    0,
		Received: 0,
		LastRTT: currentStats.LastRTT, // Preserve LastRTT for continuity context if needed
	}
	p.m2 = 0

	return currentStats
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
	p.lifetime.Sent++

	if err != nil {
		// Update Loss %
		p.stats.Loss = float64(p.stats.Sent-p.stats.Received) / float64(p.stats.Sent) * 100
		p.lifetime.Loss = float64(p.lifetime.Sent-p.lifetime.Received) / float64(p.lifetime.Sent) * 100
		return
	}

	p.stats.Received++
	p.lifetime.Received++
	p.stats.Loss = float64(p.stats.Sent-p.stats.Received) / float64(p.stats.Sent) * 100
	p.lifetime.Loss = float64(p.lifetime.Sent-p.lifetime.Received) / float64(p.lifetime.Sent) * 100

	matches := macStatsRegex.FindStringSubmatch(string(out))
	if len(matches) == 5 {
		rtt, _ := strconv.ParseFloat(matches[2], 64) // Use avg from ping output as the RTT
		p.stats.LastRTT = rtt
		p.lifetime.LastRTT = rtt

		// Update Interval Stats
		if p.stats.Min == 0 { // First successful sample in window
			p.stats.Min = rtt
			p.stats.Max = rtt
			p.stats.Avg = rtt
			p.stats.StdDev = 0
			p.m2 = 0
		} else {
			if rtt < p.stats.Min { p.stats.Min = rtt }
			if rtt > p.stats.Max { p.stats.Max = rtt }
			delta := rtt - p.stats.Avg
			p.stats.Avg += delta / float64(p.stats.Received)
			delta2 := rtt - p.stats.Avg
			p.m2 += delta * delta2
			p.stats.StdDev = math.Sqrt(p.m2 / float64(p.stats.Received))
		}

		// Update Lifetime Stats
		if p.lifetime.Received == 1 {
			p.lifetime.Min = rtt
			p.lifetime.Max = rtt
			p.lifetime.Avg = rtt
			p.lifetime.StdDev = 0
			p.lifeM2 = 0
		} else {
			if rtt < p.lifetime.Min { p.lifetime.Min = rtt }
			if rtt > p.lifetime.Max { p.lifetime.Max = rtt }
			delta := rtt - p.lifetime.Avg
			p.lifetime.Avg += delta / float64(p.lifetime.Received)
			delta2 := rtt - p.lifetime.Avg
			p.lifeM2 += delta * delta2
			p.lifetime.StdDev = math.Sqrt(p.lifeM2 / float64(p.lifetime.Received))
		}
	}
}
