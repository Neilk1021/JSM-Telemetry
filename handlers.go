package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type EventHandler struct {
	eventCh chan TelemetryEvent
	apiKey  string
}

type StatsHandler struct {
	store *TelemetryStore
}

func (s *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events, err := s.store.GetAllEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(&events); err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (e *EventHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providedKey := r.Header.Get("X-API-Key")
	if providedKey != e.apiKey {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var event TelemetryEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}

	event.Timestamp = time.Now()

	e.eventCh <- event
	w.WriteHeader(http.StatusAccepted)
}

func NewEventHandler(eventCh chan TelemetryEvent, apiKey string) *EventHandler {
	return &EventHandler{
		eventCh: eventCh,
		apiKey:  apiKey,
	}
}

func NewStatsHandler(store *TelemetryStore) *StatsHandler {
	return &StatsHandler{
		store: store,
	}
}
