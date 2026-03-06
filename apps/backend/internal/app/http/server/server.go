package server

import (
	"ai-wardrobe/internal/app/deps"
	"fmt"
	"net/http"
)

func New(d deps.Deps) (http.Handler, error) {
	wardrobeMux := http.NewServeMux()
	if err := wardrobe.Register(wardrobeMux, d); err != nil {
		return nil, fmt.Errorf("register wardrobe module: %w", err)
	}

	// root mux
	mux := http.NewServeMux()

	mux.Handle("/api/v1/wardrobe/", wardrobeMux)

	return mux, nil
}
