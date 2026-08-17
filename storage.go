package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventSink interface {
	Store(ctx context.Context, event TelemetryEvent) error
	Get(ctx context.Context) ([]TelemetryEvent, error)
}

type TelemetryStore struct {
	pool *pgxpool.Pool
}

func (t *TelemetryStore) StoreEvent(ctx context.Context, event TelemetryEvent) error {
	_, err := t.pool.Exec(ctx,
		"INSERT INTO events (player_id, event_type, timestamp) VALUES ($1, $2, $3)",
		event.PlayerID, event.EventType, event.Timestamp,
	)

	return err
}

func (t *TelemetryStore) GetAllEvents(ctx context.Context) ([]TelemetryEvent, error) {
	rows, err := t.pool.Query(ctx,
		"SELECT player_id, event_type, timestamp FROM events",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var events []TelemetryEvent

	for rows.Next() {
		var e TelemetryEvent

		if err := rows.Scan(&e.PlayerID, &e.EventType, &e.Timestamp); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func NewTelemetryStore(ctx context.Context, connStr string) (*TelemetryStore, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &TelemetryStore{pool: pool}, nil
}

func (t *TelemetryStore) Close() {
	t.pool.Close()
}
