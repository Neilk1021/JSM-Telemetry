package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type TelemetryEvent struct {
	PlayerID  string    `json:"player_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {

	ctx := context.Background()

	connStr := os.Getenv("DATABASE_URL")

	store, err := NewTelemetryStore(ctx, connStr)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	apiKey := os.Getenv("TELEMETRY_API_KEY")
	if apiKey == "" {
		panic("TELEMETRY_API_KEY must be set")
	}

	eventCh := make(chan TelemetryEvent, 101)

	eventHandler := NewEventHandler(eventCh, apiKey)
	statsHandler := NewStatsHandler(store)

	http.HandleFunc("/event", eventHandler.PostEvent)
	http.HandleFunc("/stats", statsHandler.GetStats)

	go func() {
		for event := range eventCh {
			err := store.StoreEvent(ctx, event)
			if err != nil {
				fmt.Println("failed to store event:", err)
			}
		}
	}()

	fmt.Println("listening on :9082")
	http.ListenAndServe(":9082", nil)
}
