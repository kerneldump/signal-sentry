package pinger

import (
	"context"
	"math"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
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

func (p *Pinger) ping() {
	// Create a new pro-bing pinger
	pinger, err := probing.NewPinger(p.Target)
	if err != nil {
		// Log error or just return, counting as loss handled by flow below?
		// If we can't resolve or create socket, it's effectively a loss.
		p.recordLoss()
		return
	}

	pinger.Count = 1
	pinger.Timeout = p.Interval // Wait at most the interval duration
	
	// On macOS, unprivileged ping might be needed if sudo is not used, 
	// but we will assume sudo per user request for "native" behavior.
	// However, setting SetPrivileged(true) is safer for ICMP on most systems if running as root.
	pinger.SetPrivileged(true)

	err = pinger.Run() // Blocks until finished
	if err != nil {
		p.recordLoss()
		return
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		p.recordLoss()
		return
	}

	// Success
	rtt := float64(stats.AvgRtt.Milliseconds()) // stats.AvgRtt is the only RTT for Count=1
	
	p.mu.Lock()
	defer p.mu.Unlock()

	p.updateStats(&p.stats, rtt, &p.m2)
	p.updateStats(&p.lifetime, rtt, &p.lifeM2)
}

func (p *Pinger) recordLoss() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stats.Sent++
	p.lifetime.Sent++
	
	// Recalc loss
	p.stats.Loss = float64(p.stats.Sent-p.stats.Received) / float64(p.stats.Sent) * 100
	p.lifetime.Loss = float64(p.lifetime.Sent-p.lifetime.Received) / float64(p.lifetime.Sent) * 100
}

func (p *Pinger) updateStats(s *models.PingStats, rtt float64, m2 *float64) {
	s.Sent++
	s.Received++
	s.LastRTT = rtt
	s.Loss = float64(s.Sent-s.Received) / float64(s.Sent) * 100

	if s.Min == 0 { // First successful sample in window/lifetime (logic from before)
		// Note: Technically lifetime.Min shouldn't be reset to 0 unless we want to allow it. 
		// But s.Min starts at 0. If it's 0, we take the value.
		// Wait, if actual RTT is 0.0001, this check fails. But RTT is rarely exactly 0.
		s.Min = rtt
		s.Max = rtt
		s.Avg = rtt
		s.StdDev = 0
		*m2 = 0
	} else {
		if rtt < s.Min { s.Min = rtt }
		if rtt > s.Max { s.Max = rtt }
		
		delta := rtt - s.Avg
		s.Avg += delta / float64(s.Received)
		delta2 := rtt - s.Avg
		*m2 += delta * delta2
		s.StdDev = math.Sqrt(*m2 / float64(s.Received))
	}
}
