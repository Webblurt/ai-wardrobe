package wardrobe

import (
	"ai-wardrobe/internal/app/deps"
	"ai-wardrobe/internal/wardrobe/api"
	"ai-wardrobe/internal/wardrobe/clients/replicate"
	"ai-wardrobe/internal/wardrobe/storage"
	"ai-wardrobe/internal/wardrobe/usecase"
	"fmt"
	"net/http"
)

func Register(mux *http.ServeMux, d deps.Deps) error {
	replicateCli, err := replicate.New(d.Config, d.Logger)
	if err != nil {
		return fmt.Errorf("init replicate client: %w", err)
	}
	storage, err := storage.New(&d.Config.Storage, d.Logger)
	if err != nil {
		return fmt.Errorf("init storage: %w", err)
	}
	service, err := usecase.New(replicateCli, storage, d.Logger, d.Config)
	if err != nil {
		return fmt.Errorf("init service: %w", err)
	}
	handler, err := api.New(service, d.Logger, d.Config)
	if err != nil {
		return fmt.Errorf("init handler: %w", err)
	}

	mux.HandleFunc("/api/v1/wardrobe/try-on", handler.TryOn)
	mux.HandleFunc("/api/v1/wardrobe/try-on/", handler.TryOnSubroutes)

	return nil
}
