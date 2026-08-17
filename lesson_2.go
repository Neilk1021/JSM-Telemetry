package main

import (
	"fmt"
	"time"
)

type TelemetryEvent struct {
	PlayerID  string
	EventType string
	Timestamp time.Time
	Payload   map[string]interface{}
}

func (e TelemetryEvent) Summary() string {
	return fmt.Sprintf(
		"[%s] player= %s type=%s",
		e.Timestamp.Format(time.RFC3339),
		e.PlayerID,
		e.EventType,
	)
}

func (e *TelemetryEvent) MarkedProcessed() {
	e.Payload["processed"] = true
}

func main() {
	event := TelemetryEvent{
		PlayerID:  "cnel-01",
		EventType: "hack_attempt",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"target": "mainframe", "success": true},
	}

	fmt.Println(event.Summary())

	event.MarkedProcessed()
	fmt.Println(event.Payload)
}
