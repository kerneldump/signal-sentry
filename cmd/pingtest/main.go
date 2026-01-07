package main

import (
	"fmt"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// This script compares "One-Shot" pinging (Current Implementation)
// vs "Continuous" pinging (Proposed Implementation).

const (
	target = "8.8.8.8"
	count  = 10
)

func main() {
	fmt.Println("--- Starting Ping Reliability Test ---")
	fmt.Printf("Target: %s\n", target)
	fmt.Println("Note: Run with sudo if packet loss is observed due to permission errors.")

	fmt.Println("\n[Test 1] One-Shot Mode (Current Implementation)")
	testOneShot()

	fmt.Println("\n[Test 2] Continuous Mode (Proposed Implementation)")
	testContinuous()
}

func testOneShot() {

lost := 0
totalRtt := time.Duration(0)

	start := time.Now()
	for i := 0; i < count; i++ {
		pinger, err := probing.NewPinger(target)
		if err != nil {
			fmt.Printf("Error init: %v\n", err)
		
lost++
			continue
		}
		pinger.Count = 1
		pinger.Timeout = 1 * time.Second
		pinger.SetPrivileged(true)

		err = pinger.Run()
		if err != nil {
			fmt.Printf("Error run: %v\n", err)
		
lost++
			continue
		}

		stats := pinger.Statistics()
		if stats.PacketsRecv == 0 {
			fmt.Print("X") // Lost
		
lost++
		} else {
			fmt.Print(".") // Received
			totalRtt += stats.AvgRtt
		}
		time.Sleep(1 * time.Second) // Wait for next tick
	}
	duration := time.Since(start)
	if count-lost == 0 {
		fmt.Printf("\nOne-Shot Results: 0/%d received, %d lost. Total Time: %v\n", count, lost, duration)
	} else {
		fmt.Printf("\nOne-Shot Results: %d/%d received, %d lost. Avg RTT: %v. Total Time: %v\n", 
			count-lost, count, lost, totalRtt/time.Duration(count-lost), duration)
	}
}

func testContinuous() {
	pinger, err := probing.NewPinger(target)
	if err != nil {
		fmt.Printf("Error init: %v\n", err)
		return
	}

	pinger.Count = count
	pinger.Interval = 1 * time.Second
	pinger.Timeout = time.Duration(count+2) * time.Second // Allow extra time for final packet
	pinger.SetPrivileged(true)

	// Callback for real-time feedback
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Print(".")
	}
	
	// We don't need manual loop, pinger handles it
	start := time.Now()
	err = pinger.Run()
	if err != nil {
		fmt.Printf("Error run: %v\n", err)
		return
	}
	duration := time.Since(start)

	stats := pinger.Statistics()
	fmt.Printf("\nContinuous Results: %d/%d received, %.1f%% loss. Avg RTT: %v. Total Time: %v\n", 
		stats.PacketsRecv, stats.PacketsSent, stats.PacketLoss, stats.AvgRtt, duration)
}
