package domain

type ReplicateUploadResp struct {
	URL string `json:"url"`
}

type ReplicatePredictionResp struct {
	ID     string `json:"id"`
	Status any    `json:"status"`
	Output any    `json:"output"`
}
