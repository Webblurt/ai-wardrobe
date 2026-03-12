package api

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func parseInt(v string) int {
	if v == "" {
		return 0
	}

	i, _ := strconv.Atoi(v)
	return i
}

func parseBool(v string) bool {
	if v == "" {
		return false
	}

	b, _ := strconv.ParseBool(v)
	return b
}
