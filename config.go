package main

import "fmt"

func validateInterval(interval int) error {
	if interval <= 0 {
		return fmt.Errorf("interval must be greater than 0, got %d", interval)
	}
	return nil
}
