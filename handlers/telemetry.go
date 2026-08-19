package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"telemetry/events"
	"telemetry/storage"
)

type TelemetryHandler struct {
	eventCh chan events.TelemetryEvent
	apiKey  string
	store   *storage.Store
}

func NewTelemetryHandler(store *storage.Store, apiKey string, bufferSize int) *TelemetryHandler {
	return &TelemetryHandler{
		eventCh: make(chan events.TelemetryEvent, bufferSize),
		apiKey:  apiKey,
		store:   store,
	}
}

func (h *TelemetryHandler) GetApiKey() string {
	return h.apiKey
}

func (h *TelemetryHandler) StartWorker(ctx context.Context) {
	go func() {
		for event := range h.eventCh {
			err := h.store.StoreEvent(ctx, event)
			if err != nil {
				fmt.Println("failed to store event:", err)
			}
		}
	}()
}

func (h *TelemetryHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	if !ValidatePostRequest(w, r, h) {
		return
	}
	HandlePostEvent(w, r, h.eventCh, h.apiKey)
}

func (h *TelemetryHandler) GetEventStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	telemetryEvents, err := h.store.GetAllEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&telemetryEvents); err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}
}

func (h *TelemetryHandler) GetEventsPaginated(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	limit, offset := 50, 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	eventsList, err := h.store.GetEventsPaginated(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&eventsList); err != nil {
		http.Error(w, "internal_error", http.StatusInternalServerError)
	}
}

func (h *TelemetryHandler) GetEventDailyStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stats, err := h.store.GetEventDailyStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&stats); err != nil {
		http.Error(w, "internal_error", http.StatusInternalServerError)
	}
}
