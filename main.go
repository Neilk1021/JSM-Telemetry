package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"telemetry/handlers"
	"telemetry/storage"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()

	connStr := os.Getenv("DATABASE_URL")
	apiKey := os.Getenv("TELEMETRY_API_KEY")
	if apiKey == "" {
		panic("TELEMETRY_API_KEY must be set")
	}

	m, err := migrate.New("file://db/migrations", connStr)
	if err != nil {
		fmt.Printf("Warning: failed to initialize migrate: %v\n", err)
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Warning: failed to run migrations: %v\n", err)
		} else {
			fmt.Println("Database migrations ran successfully")
		}
	}

	store, err := storage.NewStore(ctx, connStr)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	telemetryHandler := handlers.NewTelemetryHandler(store, apiKey, 101)
	telemetryHandler.StartWorker(ctx)

	mapAttemptHandler := handlers.NewMapAttemptHandler(store, apiKey, 101)
	mapAttemptHandler.StartWorker(ctx)

	http.HandleFunc("/event", telemetryHandler.PostEvent)
	http.HandleFunc("/stats/events", telemetryHandler.GetEventStats)
	http.HandleFunc("/stats/events/paginated", telemetryHandler.GetEventsPaginated)
	http.HandleFunc("/stats/events/daily", telemetryHandler.GetEventDailyStats)

	http.HandleFunc("/map_attempt", mapAttemptHandler.PostEvent)
	http.HandleFunc("/stats/map_attempts", mapAttemptHandler.GetMapStats)
	http.HandleFunc("/stats/map_attempts/paginated", mapAttemptHandler.GetMapAttemptsPaginated)
	http.HandleFunc("/stats/map_attempts/summary", mapAttemptHandler.GetMapSummaryStats)

	fmt.Println("listening on :9082")
	http.ListenAndServe(":9082", nil)
}
