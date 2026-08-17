package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TelemetryEvent struct {
	PlayerID  string    `json:"player_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {

	ctx := context.Background()

	connStr := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	eventCh := make(chan TelemetryEvent, 100)

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var event TelemetryEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "bad_request", http.StatusBadRequest)
			return
		}

		event.Timestamp = time.Now()

		eventCh <- event
		w.WriteHeader(http.StatusAccepted)
	})

	go func() {
		for event := range eventCh {
			_, err := pool.Exec(ctx,
				"INSERT INTO events (player_id, event_type, timestamp) VALUES ($1, $2, $3)",
				event.PlayerID, event.EventType, event.Timestamp,
			)

			if err != nil {
				fmt.Println("failed to store event:", err)
			}
		}
	}()

	fmt.Println("listening on :9081")
	http.ListenAndServe(":9081", nil)
}
