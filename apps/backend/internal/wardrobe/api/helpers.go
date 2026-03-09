package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, msg string, status int) {
	h.writeJSON(w, map[string]string{
		"error": msg,
	}, status)
}
