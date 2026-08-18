package events

import "time"

type TelemetryEvent struct {
	PlayerID  string    `json:"player_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

type MapAttemptEvent struct {
	PlayerID  string          `json:"player_id"`
	MapName   string          `json:"map_name"`
	Passed    bool            `json:"passed"`
	Attempts  int             `json:"attempts"`
	Rounds    int             `json:"rounds"`
	RamUsed   int             `json:"ram_used"`
	CpuUsed   int             `json:"cpu_used"`
	Duration  time.Duration   `json:"duration"`
	FirstTime bool            `json:"first_time"`
	Towers    []MapTowerUsage `json:"towers"`
	Timestamp time.Time       `json:"timestamp"`
}

type MapTowerUsage struct {
	Tower string `json:"tower"`
	Level int    `json:"level"`
	Count int    `json:"count"`
}

type Timestamped interface {
	SetTimestamp(t time.Time)
}

func (e *TelemetryEvent) SetTimestamp(t time.Time) {
	e.Timestamp = t
}

func (m *MapAttemptEvent) SetTimestamp(t time.Time) {
	m.Timestamp = t
}
