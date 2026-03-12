package usecase

import (
	"ai-wardrobe/internal/wardrobe/domain"
	"context"
)

type ReplicateClient interface {
	PostTryOn(ctx context.Context, params domain.TryOnParams, personURL, garmentURL string) (string, error)
}

type FedjazVtonClient interface {
	PostTryOn(ctx context.Context, params domain.TryOnParams, personPath, garmentPath string) ([]byte, error)
}
