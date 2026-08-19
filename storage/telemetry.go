package storage

import (
	"context"
	"encoding/json"
	"telemetry/events"
	"time"
)

func (s *Store) StoreEvent(ctx context.Context, event events.TelemetryEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx,
		"INSERT INTO events (player_id, event_type, timestamp, payload) VALUES ($1, $2, $3, $4)",
		event.PlayerID, event.EventType, event.Timestamp, payloadJSON,
	)
	return err
}

func (s *Store) GetAllEvents(ctx context.Context) ([]events.TelemetryEvent, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT player_id, event_type, timestamp, payload FROM events",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var telemetryEvents []events.TelemetryEvent
	for rows.Next() {
		var e events.TelemetryEvent
		var payloadJSON []byte
		if err := rows.Scan(&e.PlayerID, &e.EventType, &e.Timestamp, &payloadJSON); err != nil {
			return nil, err
		}
		if len(payloadJSON) > 0 {
			json.Unmarshal(payloadJSON, &e.Payload)
		}
		telemetryEvents = append(telemetryEvents, e)
	}
	return telemetryEvents, nil
}

func (s *Store) GetEventsPaginated(ctx context.Context, limit, offset int) ([]events.TelemetryEvent, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT player_id, event_type, timestamp, payload FROM events ORDER BY timestamp DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var telemetryEvents []events.TelemetryEvent
	for rows.Next() {
		var e events.TelemetryEvent
		var payloadJSON []byte
		if err := rows.Scan(&e.PlayerID, &e.EventType, &e.Timestamp, &payloadJSON); err != nil {
			return nil, err
		}
		if len(payloadJSON) > 0 {
			json.Unmarshal(payloadJSON, &e.Payload)
		}
		telemetryEvents = append(telemetryEvents, e)
	}
	return telemetryEvents, nil
}

func (s *Store) GetEventDailyStats(ctx context.Context) ([]events.DailyEventStat, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT DATE(timestamp) as date, event_type, COUNT(*) as count 
		 FROM events 
		 GROUP BY DATE(timestamp), event_type 
		 ORDER BY date DESC, count DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []events.DailyEventStat
	for rows.Next() {
		var stat events.DailyEventStat
		var t time.Time
		if err := rows.Scan(&t, &stat.EventType, &stat.Count); err != nil {
			return nil, err
		}
		stat.Date = t.Format("2006-01-02")
		stats = append(stats, stat)
	}
	return stats, nil
}
