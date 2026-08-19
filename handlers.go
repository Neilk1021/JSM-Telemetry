package main

import (
	"encoding/json"
	"net/http"
	"telemetry/events"
	"time"
)

type PostHandler interface {
	GetApiKey() string
}

type EventHandler struct {
	eventCh chan events.TelemetryEvent
	apiKey  string
}

type StatsHandler struct {
	store *TelemetryStore
}

type MapHandler struct {
	apiKey  string
	eventCh chan events.MapAttemptEvent
}

func (e *EventHandler) GetApiKey() string {
	return e.apiKey
}

func (m *MapHandler) GetApiKey() string {
	return m.apiKey
}

func IsPostMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func ValidApiKey(w http.ResponseWriter, r *http.Request, p PostHandler) bool {
	providedKey := r.Header.Get("X-API-Key")
	if providedKey != p.GetApiKey() {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func ValidatePostRequest(w http.ResponseWriter, r *http.Request, p PostHandler) bool {
	if !IsPostMethod(w, r) {
		return false
	}

	if !ValidApiKey(w, r, p) {
		return false
	}
	return true
}

func (s *StatsHandler) GetEventStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	telemetryEvents, err := s.store.GetAllEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(&telemetryEvents); err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *StatsHandler) GetMapStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mapAttempts, err := s.store.GetAllMapAttempts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(&mapAttempts); err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func HandlePostEvent[T any, PT interface {
	*T
	events.Timestamped
}](w http.ResponseWriter, r *http.Request, ch chan T, apiKey string) {
	var event T
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}

	PT(&event).SetTimestamp(time.Now())

	ch <- event
	w.WriteHeader(http.StatusAccepted)
}

func (e *EventHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	if !ValidatePostRequest(w, r, e) {
		return
	}

	HandlePostEvent(w, r, e.eventCh, e.apiKey)
}

func (m *MapHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	if !ValidatePostRequest(w, r, m) {
		return
	}

	HandlePostEvent(w, r, m.eventCh, m.apiKey)
}

//Handlers

func NewEventHandler(eventCh chan events.TelemetryEvent, apiKey string) *EventHandler {
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

func NewMapHandler(mapCh chan events.MapAttemptEvent, apiKey string) *MapHandler {
	return &MapHandler{
		apiKey:  apiKey,
		eventCh: mapCh,
	}
}
