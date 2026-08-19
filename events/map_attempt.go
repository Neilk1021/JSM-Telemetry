package events

import "time"

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

func (m *MapAttemptEvent) SetTimestamp(t time.Time) {
	m.Timestamp = t
}

type MapSummaryStat struct {
	MapName       string  `json:"map_name"`
	TotalAttempts int     `json:"total_attempts"`
	PassRate      float64 `json:"pass_rate"`
	AvgDuration   float64 `json:"avg_duration"`
	AvgRounds     float64 `json:"avg_rounds"`
}
