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

	provider := r.FormValue("provider")

	switch provider {
	case "fedjaz", "fedjaz_fashn_v1.5", "replicate":
	default:
		h.writeError(w, "unknown provider", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	category := r.FormValue("category")

	if category == "" {
		h.writeError(w, "category is required", http.StatusBadRequest)
		return
	}

	steps := parseInt(r.FormValue("steps"))
	seed := parseInt(r.FormValue("seed"))
	autocrop := parseBool(r.FormValue("autocrop"))
	upscale := parseInt(r.FormValue("upscale"))
	upscaler := r.FormValue("upscaler")

	// ---- fashn params ----

	garmentPhotoType := r.FormValue("garmentPhotoType")
	numSamples := parseInt(r.FormValue("numSamples"))
	numTimesteps := parseInt(r.FormValue("numTimesteps"))
	guidanceScale := parseFloat(r.FormValue("guidanceScale"))
	segmentationFree := parseBool(r.FormValue("segmentationFree"))

	// ---- files ----

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
		Provider: provider,

		Description: description,
		Category:    category,
		Steps:       steps,
		Seed:        seed,
		Autocrop:    autocrop,
		Upscale:     upscale,
		Upscaler:    upscaler,

		GarmentPhotoType: garmentPhotoType,
		NumSamples:       numSamples,
		NumTimesteps:     numTimesteps,
		GuidanceScale:    guidanceScale,
		SegmentationFree: segmentationFree,

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
