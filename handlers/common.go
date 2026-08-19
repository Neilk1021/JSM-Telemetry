package handlers

import (
	"encoding/json"
	"net/http"
	"telemetry/events"
	"time"
)

type PostHandler interface {
	GetApiKey() string
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
