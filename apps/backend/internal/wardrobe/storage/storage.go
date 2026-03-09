package storage

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/platform/logger"
	"ai-wardrobe/internal/wardrobe/domain"
	"encoding/json"
	"os"
	"path/filepath"
)

type Storage struct {
	cfg    *config.Storage
	logger *logger.Logger
}

func New(cfg *config.Storage, logger *logger.Logger) (*Storage, error) {
	return &Storage{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (s *Storage) SaveJob(job domain.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	path := filepath.Join(s.cfg.JobsDir, job.JobID+".json")

	return os.WriteFile(path, data, 0644)
}

func (s *Storage) UpdateJobStatus(jobID string, status domain.JobStatus, resultURL string) error {
	job, err := s.LoadJob(jobID)
	if err != nil {
		s.logger.Error("load job failed", err)
		return err
	}

	job.Status = status
	job.ResultURL = resultURL

	if err := s.SaveJob(job); err != nil {
		s.logger.Error("save job failed", err)
		return err
	}

	return nil
}

func (s *Storage) LoadJob(jobID string) (domain.Job, error) {
	path := filepath.Join(s.cfg.JobsDir, jobID+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Job{}, err
	}

	var job domain.Job

	if err := json.Unmarshal(data, &job); err != nil {
		return domain.Job{}, err
	}

	return job, nil
}
