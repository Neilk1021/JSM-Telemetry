package main

import (
	"context"
	"telemetry/events"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TelemetryStore struct {
	pool *pgxpool.Pool
}

func (t *TelemetryStore) StoreEvent(ctx context.Context, event events.TelemetryEvent) error {
	_, err := t.pool.Exec(ctx,
		"INSERT INTO events (player_id, event_type, timestamp) VALUES ($1, $2, $3)",
		event.PlayerID, event.EventType, event.Timestamp,
	)

	return err
}

func (t *TelemetryStore) GetAllEvents(ctx context.Context) ([]events.TelemetryEvent, error) {
	rows, err := t.pool.Query(ctx,
		"SELECT player_id, event_type, timestamp FROM events",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var telemetryEvents []events.TelemetryEvent

	for rows.Next() {
		var e events.TelemetryEvent

		if err := rows.Scan(&e.PlayerID, &e.EventType, &e.Timestamp); err != nil {
			return nil, err
		}

		telemetryEvents = append(telemetryEvents, e)
	}

	return telemetryEvents, nil
}

func (t *TelemetryStore) StoreMapAttempt(ctx context.Context, event events.MapAttemptEvent) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx,
		`INSERT INTO map_attempts 
    			( 
				 player_id, 
				 map_name, 
				 passed, 
    			 attempts_to_complete, 
    			 level_rounds, 
    			 ram_used, 
    			 cpu_used, 
    			 duration, 
    			 first_time, 
    			 timestamp
    			 ) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		event.PlayerID, event.MapName, event.Passed, event.Attempts, event.Rounds,
		event.RamUsed, event.CpuUsed, event.Duration, event.FirstTime, time.Now(),
	).Scan(&id)

	if err != nil {
		return err
	}

	for _, e := range event.Towers {
		_, err := tx.Exec(ctx,
			`INSERT INTO map_tower_usage 
    				(map_attempt_id, 
    				 tower_type, 
    				 count, 
    				 level
    				 ) 
					VALUES  ($1, $2, $3, $4)`,
			id, e.Tower, e.Count, e.Level,
		)

		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	return err
}

func (t *TelemetryStore) GetAllMapAttempts(ctx context.Context) ([]events.MapAttemptEvent, error) {
	rows, err := t.pool.Query(ctx,
		`SELECT id, player_id, map_name, passed, attempts_to_complete, level_rounds,
		        ram_used, cpu_used, duration, first_time, timestamp
		 FROM map_attempts`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attemptsByID := make(map[int]*events.MapAttemptEvent)
	var order []int

	for rows.Next() {
		var id int
		var e events.MapAttemptEvent

		if err := rows.Scan(&id, &e.PlayerID, &e.MapName, &e.Passed, &e.Attempts,
			&e.Rounds, &e.RamUsed, &e.CpuUsed, &e.Duration, &e.FirstTime, &e.Timestamp); err != nil {
			return nil, err
		}

		e.Towers = []events.MapTowerUsage{}
		attemptsByID[id] = &e
		order = append(order, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	towerRows, err := t.pool.Query(ctx,
		`SELECT map_attempt_id, tower_type, count, level FROM map_tower_usage`,
	)
	if err != nil {
		return nil, err
	}
	defer towerRows.Close()

	for towerRows.Next() {
		var mapAttemptID int
		var tower events.MapTowerUsage

		if err := towerRows.Scan(&mapAttemptID, &tower.Tower, &tower.Count, &tower.Level); err != nil {
			return nil, err
		}

		if attempt, ok := attemptsByID[mapAttemptID]; ok {
			attempt.Towers = append(attempt.Towers, tower)
		}
	}
	if err := towerRows.Err(); err != nil {
		return nil, err
	}

	result := make([]events.MapAttemptEvent, 0, len(order))
	for _, id := range order {
		result = append(result, *attemptsByID[id])
	}

	return result, nil
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
