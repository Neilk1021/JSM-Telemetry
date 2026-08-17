package main

import (
	"fmt"
	"time"
)

type TelemetryEvent struct {
	PlayerID  string
	EventType string
}

func simulateClient(playerID string, eventCh chan<- TelemetryEvent, done chan<- bool) {
	events := []string{"login", "hack_attempt", "hack_success", "logout"}
	for _, evt := range events {
		time.Sleep(50 * time.Millisecond)
		eventCh <- TelemetryEvent{PlayerID: playerID, EventType: evt}
	}
	done <- true
}

func main() {
	eventCh := make(chan TelemetryEvent)
	done := make(chan bool)

	players := []string{"cnel-01", "ike-02", "nagano-99"}

	for _, p := range players {
		go simulateClient(p, eventCh, done)
	}

	received := 0
	finished := 0
	total := len(players) * 4

	for received < total {
		select {
		case event := <-eventCh:
			received++
			fmt.Printf("[received %d/%d] %s -> %s\n", received, total, event.PlayerID, event.EventType)
		case <-done:
			finished++
		}
	}

	fmt.Println("\nall events processed")
}
