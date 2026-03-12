package usecase

import "context"

type ReplicateClient interface {
	PostTryOn(ctx context.Context, personURL, garmentURL string) (string, error)
}

type FedjazVtonClient interface {
	PostTryOn(ctx context.Context, personPath, garmentPath string) ([]byte, error)
}
