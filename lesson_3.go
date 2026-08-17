package main

import (
	"fmt"
	"time"
)

type TelemetryEvent struct {
	PlayerID  string
	EventType string
	Timestamp time.Time
}

func main() {
	var events []TelemetryEvent

	events = append(events, TelemetryEvent{
		PlayerID:  "cnel-01",
		EventType: "hack_attempt",
		Timestamp: time.Now(),
	},
	)

	events = append(events, TelemetryEvent{
		PlayerID:  "cnel-02",
		EventType: "hack_success",
		Timestamp: time.Now(),
	},
	)

	events = append(events, TelemetryEvent{
		PlayerID:  "ike-02",
		EventType: "login",
		Timestamp: time.Now(),
	},
	)

	fmt.Println("total events:", len(events))

	for i, e := range events {
		fmt.Printf("%d: %s -> %s\n", i, e.PlayerID, e.EventType)
	}

	countsByPlayer := make(map[string]int)

	for _, e := range events {
		countsByPlayer[e.PlayerID]++
	}

	fmt.Println("\nevent counts:")

	for player, count := range countsByPlayer {
		fmt.Printf("%s: %d events\n", player, count)
	}

	count, exists := countsByPlayer["ike-02"]
	if !exists {
		fmt.Println("\nnagano-99 has no events yet")
	} else {
		fmt.Println(count)
	}
}
