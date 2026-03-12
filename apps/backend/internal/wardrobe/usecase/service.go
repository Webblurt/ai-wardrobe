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
	fc     FedjazVtonClient
	st     Storage
	logger *logger.Logger
	cfg    *config.Config
}

func New(replicate ReplicateClient, fc FedjazVtonClient, storage Storage, logger *logger.Logger, cfg *config.Config) (*Service, error) {
	return &Service{
		rp:     replicate,
		fc:     fc,
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

	go func() {

		switch req.Provider {

		case "fedjaz":

			s.runFedjazVtonTryOn(
				context.Background(),
				jobID,
				personPath,
				garmentPath,
			)

		case "replicate":

			s.runReplicateTryOn(
				context.Background(),
				jobID,
				personURL,
				garmentURL,
			)

		default:

			s.logger.Error("unknown provider", req.Provider)
			s.st.UpdateJobStatus(jobID, domain.StatusFailed, "")
		}

	}()

	return domain.CreateJobResp{
		JobID:  jobID,
		Status: domain.StatusProcessing,
	}, nil
}

func (s *Service) runFedjazVtonTryOn(ctx context.Context, jobID, personPath, garmentPath string) {

	img, err := s.fc.PostTryOn(ctx, personPath, garmentPath)
	if err != nil {
		s.st.UpdateJobStatus(jobID, domain.StatusFailed, "")
		s.logger.Error("try-on failed", err)
		return
	}

	resultFile := jobID + "_result.png"
	resultPath := filepath.Join(s.cfg.Storage.UploadsDir, resultFile)

	err = os.WriteFile(resultPath, img, 0644)
	if err != nil {
		s.logger.Error("write result failed", err)
		s.st.UpdateJobStatus(jobID, domain.StatusFailed, "")
		return
	}

	resultURL := s.buildImageURL(resultFile)

	s.st.UpdateJobStatus(jobID, domain.StatusCompleted, resultURL)
}

func (s *Service) runReplicateTryOn(ctx context.Context, jobID, personURL, garmentURL string) {
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
