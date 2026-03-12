package api

import (
	"ai-wardrobe/internal/wardrobe/domain"
	"io"
	"net/http"
)

func (h *Handler) postTryOn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(30 << 20); err != nil {
		h.writeError(w, "file too large", http.StatusRequestEntityTooLarge)
		return
	}

	description := r.FormValue("description")
	steps := parseInt(r.FormValue("steps"))
	seed := parseInt(r.FormValue("seed"))
	autocrop := parseBool(r.FormValue("autocrop"))
	upscale := parseInt(r.FormValue("upscale"))
	upscaler := r.FormValue("upscaler")

	provider := r.FormValue("provider")
	switch provider {
	case "fedjaz", "replicate":
		// ok
	case "":
		h.writeError(w, "provider is required", http.StatusBadRequest)
		return
	default:
		h.writeError(w, "unknown provider", http.StatusBadRequest)
		return
	}

	category := r.FormValue("category")
	if category == "" {
		h.writeError(w, "category is required", http.StatusBadRequest)
		return
	}

	fit := r.FormValue("fit")
	if fit == "" {
		h.writeError(w, "fit is required", http.StatusBadRequest)
		return
	}

	personFile, personHeader, err := r.FormFile("person")
	if err != nil {
		h.writeError(w, "person image required", http.StatusBadRequest)
		return
	}
	defer personFile.Close()

	garmentFile, garmentHeader, err := r.FormFile("garment")
	if err != nil {
		h.writeError(w, "garment image required", http.StatusBadRequest)
		return
	}
	defer garmentFile.Close()

	personData, err := io.ReadAll(io.LimitReader(personFile, 10<<20))
	if err != nil {
		h.writeError(w, "cannot read person image", http.StatusBadRequest)
		return
	}

	garmentData, err := io.ReadAll(io.LimitReader(garmentFile, 10<<20))
	if err != nil {
		h.writeError(w, "cannot read garment image", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateJob(ctx, domain.CreateJobReq{
		Provider:    provider,
		Description: description,
		Category:    category,
		Steps:       steps,
		Seed:        seed,
		Autocrop:    autocrop,
		Upscale:     upscale,
		Upscaler:    upscaler,

		Person: domain.Image{
			Data:        personData,
			ContentType: personHeader.Header.Get("Content-Type"),
		},
		Garment: domain.Image{
			Data:        garmentData,
			ContentType: garmentHeader.Header.Get("Content-Type"),
		},
	})

	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, resp, http.StatusCreated)
}

func (h *Handler) getTryOnByJobID(w http.ResponseWriter, r *http.Request, jobID string) {
	ctx := r.Context()

	resp, err := h.service.GetJobByID(ctx, jobID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, resp, http.StatusOK)
}
