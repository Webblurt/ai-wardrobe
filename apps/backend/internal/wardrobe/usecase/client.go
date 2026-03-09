package usecase

import "context"

type ReplicateClient interface {
	PostTryOn(ctx context.Context, personPath, garmentPath string) (string, error)
}
