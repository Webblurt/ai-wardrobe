package domain

type CreateJobReq struct {
	Provider string
	Category string
	Fit      string
	Person   Image
	Garment  Image
}

type Image struct {
	Data        []byte
	ContentType string
}
