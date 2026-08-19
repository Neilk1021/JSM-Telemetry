package storage

import (
	"context"
	"telemetry/events"
	"time"
)

func (s *Store) StoreMapAttempt(ctx context.Context, event events.MapAttemptEvent) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx,
		`INSERT INTO map_attempts 
    			(player_id, map_name, passed, attempts_to_complete, level_rounds, ram_used, cpu_used, duration, first_time, timestamp) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		event.PlayerID, event.MapName, event.Passed, event.Attempts, event.Rounds,
		event.RamUsed, event.CpuUsed, event.Duration, event.FirstTime, time.Now(),
	).Scan(&id)

	if err != nil {
		return err
	}

	for _, e := range event.Towers {
		_, err := tx.Exec(ctx,
			`INSERT INTO map_tower_usage (map_attempt_id, tower_type, count, level) VALUES ($1, $2, $3, $4)`,
			id, e.Tower, e.Count, e.Level,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) GetAllMapAttempts(ctx context.Context) ([]events.MapAttemptEvent, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, player_id, map_name, passed, attempts_to_complete, level_rounds, ram_used, cpu_used, duration, first_time, timestamp
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

	towerRows, err := s.pool.Query(ctx,
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

	result := make([]events.MapAttemptEvent, 0, len(order))
	for _, id := range order {
		result = append(result, *attemptsByID[id])
	}
	return result, nil
}

func (s *Store) GetMapAttemptsPaginated(ctx context.Context, limit, offset int) ([]events.MapAttemptEvent, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, player_id, map_name, passed, attempts_to_complete, level_rounds, ram_used, cpu_used, duration, first_time, timestamp
		 FROM map_attempts ORDER BY timestamp DESC LIMIT $1 OFFSET $2`,
		limit, offset,
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
	if len(order) == 0 {
		return []events.MapAttemptEvent{}, nil
	}

	towerRows, err := s.pool.Query(ctx,
		`SELECT map_attempt_id, tower_type, count, level FROM map_tower_usage WHERE map_attempt_id = ANY($1)`,
		order,
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

	result := make([]events.MapAttemptEvent, 0, len(order))
	for _, id := range order {
		result = append(result, *attemptsByID[id])
	}
	return result, nil
}

func (s *Store) GetMapSummaryStats(ctx context.Context) ([]events.MapSummaryStat, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT map_name, COUNT(*) as total_attempts, SUM(CASE WHEN passed THEN 1 ELSE 0 END)::FLOAT / COUNT(*) as pass_rate,
		        AVG(duration) as avg_duration, AVG(level_rounds) as avg_rounds
		 FROM map_attempts GROUP BY map_name ORDER BY total_attempts DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []events.MapSummaryStat
	for rows.Next() {
		var stat events.MapSummaryStat
		if err := rows.Scan(&stat.MapName, &stat.TotalAttempts, &stat.PassRate, &stat.AvgDuration, &stat.AvgRounds); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}
