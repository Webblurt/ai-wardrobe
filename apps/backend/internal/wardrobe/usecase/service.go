package usecase

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/platform/logger"
	"ai-wardrobe/internal/wardrobe/domain"
	"context"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Service struct {
	rp     ReplicateClient
	st     Storage
	logger *logger.Logger
	cfg    *config.Config
}

func New(replicate ReplicateClient, storage Storage, logger *logger.Logger, cfg *config.Config) (*Service, error) {
	return &Service{
		rp:     replicate,
		st:     storage,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (s *Service) CreateJob(ctx context.Context, req domain.CreateJobReq) (domain.CreateJobResp, error) {
	jobID := uuid.New().String()

	personPath := filepath.Join(s.cfg.Storage.UploadsDir, jobID+"_person.jpg")
	garmentPath := filepath.Join(s.cfg.Storage.UploadsDir, jobID+"_garment.jpg")

	if err := os.WriteFile(personPath, req.Person.Data, 0644); err != nil {
		return domain.CreateJobResp{}, err
	}

	if err := os.WriteFile(garmentPath, req.Garment.Data, 0644); err != nil {
		return domain.CreateJobResp{}, err
	}

	job := domain.Job{
		JobID:  jobID,
		Status: domain.StatusProcessing,
	}

	if err := s.st.SaveJob(job); err != nil {
		return domain.CreateJobResp{}, err
	}

	go s.runTryOn(context.WithoutCancel(ctx), jobID, personPath, garmentPath)

	return domain.CreateJobResp{
		JobID:  jobID,
		Status: domain.StatusProcessing,
	}, nil
}

func (s *Service) runTryOn(ctx context.Context, jobID, personPath, garmentPath string) {
	resultURL, err := s.rp.PostTryOn(ctx, personPath, garmentPath)

	if err != nil {

		s.st.UpdateJobStatus(jobID, domain.StatusFailed, "")
		s.logger.Error("try-on failed", err)

		return
	}

	s.st.UpdateJobStatus(jobID, domain.StatusCompleted, resultURL)
}

func (s *Service) GetJobByID(ctx context.Context, jobID string) (domain.GetJobByIDResp, error) {
	job, err := s.st.LoadJob(jobID)
	if err != nil {
		return domain.GetJobByIDResp{}, err
	}

	return domain.GetJobByIDResp{
		JobID:     job.JobID,
		Status:    job.Status,
		ResultURL: job.ResultURL,
	}, nil
}
