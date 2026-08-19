package events

import (
	"encoding/json"
	"time"
)

type Timestamped interface {
	SetTimestamp(t time.Time)
}

type TelemetryEvent struct {
	PlayerID  string                 `json:"player_id"`
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

func (e *TelemetryEvent) UnmarshalJSON(data []byte) error {
	type Alias TelemetryEvent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	delete(raw, "player_id")
	delete(raw, "event_type")
	delete(raw, "timestamp")
	e.Payload = raw
	return nil
}

func (e *TelemetryEvent) SetTimestamp(t time.Time) {
	e.Timestamp = t
}

type DailyEventStat struct {
	Date      string `json:"date"`
	EventType string `json:"event_type"`
	Count     int    `json:"count"`
}
