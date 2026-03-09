package domain

type JobStatus string

const (
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type CreateJobResp struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
}

type GetJobByIDResp struct {
	JobID     string    `json:"job_id"`
	Status    JobStatus `json:"status"`
	ResultURL string    `json:"result_url,omitempty"`
}
