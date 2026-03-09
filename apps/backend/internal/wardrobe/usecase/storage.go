package usecase

import "ai-wardrobe/internal/wardrobe/domain"

type Storage interface {
	SaveJob(job domain.Job) error
	UpdateJobStatus(jobID string, status domain.JobStatus, resultURL string) error
	LoadJob(jobID string) (domain.Job, error)
}
