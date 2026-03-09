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

	personData, err := io.ReadAll(personFile)
	if err != nil {
		h.writeError(w, "cannot read person image", http.StatusInternalServerError)
		return
	}

	garmentData, err := io.ReadAll(garmentFile)
	if err != nil {
		h.writeError(w, "cannot read garment image", http.StatusInternalServerError)
		return
	}

	resp, err := h.service.CreateJob(ctx, domain.CreateJobReq{
		Category: category,
		Fit:      fit,
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
