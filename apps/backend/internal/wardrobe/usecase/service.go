package usecase

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/platform/logger"
	"ai-wardrobe/internal/wardrobe/domain"
	"context"
	"os"
	"path/filepath"
	"strings"

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

	personFile := jobID + "_person.jpg"
	garmentFile := jobID + "_garment.jpg"

	personPath := filepath.Join(s.cfg.Storage.UploadsDir, personFile)
	garmentPath := filepath.Join(s.cfg.Storage.UploadsDir, garmentFile)

	if err := os.WriteFile(personPath, req.Person.Data, 0644); err != nil {
		s.logger.Error("write person file failed", err)
		return domain.CreateJobResp{}, err
	}

	if err := os.WriteFile(garmentPath, req.Garment.Data, 0644); err != nil {
		s.logger.Error("write garment file failed", err)
		return domain.CreateJobResp{}, err
	}

	personURL := s.buildImageURL(personFile)
	garmentURL := s.buildImageURL(garmentFile)

	job := domain.Job{
		JobID:  jobID,
		Status: domain.StatusProcessing,
	}

	if err := s.st.SaveJob(job); err != nil {
		s.logger.Error("save job failed", err)
		return domain.CreateJobResp{}, err
	}

	go s.runTryOn(context.WithoutCancel(ctx), jobID, personURL, garmentURL)

	return domain.CreateJobResp{
		JobID:  jobID,
		Status: domain.StatusProcessing,
	}, nil
}

func (s *Service) runTryOn(ctx context.Context, jobID, personURL, garmentURL string) {
	resultURL, err := s.rp.PostTryOn(ctx, personURL, garmentURL)
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

func (s *Service) buildImageURL(filename string) string {
	return strings.TrimRight(s.cfg.App.URL, "/") + "/images/" + filename
}
