package api

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/platform/logger"
	"ai-wardrobe/internal/wardrobe/usecase"
	"net/http"
	"strings"
)

type Handler struct {
	service *usecase.Service
	logger  *logger.Logger
	cfg     *config.Config
}

func New(service *usecase.Service, logger *logger.Logger, cfg *config.Config) (*Handler, error) {
	return &Handler{
		service: service,
		logger:  logger,
		cfg:     cfg,
	}, nil
}

func (h *Handler) TryOn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.postTryOn(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TryOnSubroutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/try-on/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	jobID := parts[0]

	switch len(parts) {

	case 1:
		switch r.Method {
		case http.MethodGet:
			h.getTryOnByJobID(w, r, jobID)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	http.NotFound(w, r)
}
