package domain

type ReplicateUploadResp struct {
	URL string `json:"url"`
}

type ReplicatePredictionResp struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Output []string `json:"output"`
}
