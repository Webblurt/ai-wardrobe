package domain

type Job struct {
	JobID     string    `json:"job_id"`
	Status    JobStatus `json:"status"`
	ResultURL string    `json:"result_url,omitempty"`
}
