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

type MapAttemptHandler struct {
	eventCh chan events.MapAttemptEvent
	apiKey  string
	store   *storage.Store
}

func NewMapAttemptHandler(store *storage.Store, apiKey string, bufferSize int) *MapAttemptHandler {
	return &MapAttemptHandler{
		eventCh: make(chan events.MapAttemptEvent, bufferSize),
		apiKey:  apiKey,
		store:   store,
	}
}

func (h *MapAttemptHandler) GetApiKey() string {
	return h.apiKey
}

func (h *MapAttemptHandler) StartWorker(ctx context.Context) {
	go func() {
		for event := range h.eventCh {
			err := h.store.StoreMapAttempt(ctx, event)
			if err != nil {
				fmt.Println("failed to store map_attempt:", err)
			}
		}
	}()
}

func (h *MapAttemptHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
	if !ValidatePostRequest(w, r, h) {
		return
	}
	HandlePostEvent(w, r, h.eventCh, h.apiKey)
}

func (h *MapAttemptHandler) GetMapStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	mapAttempts, err := h.store.GetAllMapAttempts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&mapAttempts); err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}
}

func (h *MapAttemptHandler) GetMapAttemptsPaginated(w http.ResponseWriter, r *http.Request) {
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
	mapAttempts, err := h.store.GetMapAttemptsPaginated(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&mapAttempts); err != nil {
		http.Error(w, "internal_error", http.StatusInternalServerError)
	}
}

func (h *MapAttemptHandler) GetMapSummaryStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stats, err := h.store.GetMapSummaryStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&stats); err != nil {
		http.Error(w, "internal_error", http.StatusInternalServerError)
	}
}
