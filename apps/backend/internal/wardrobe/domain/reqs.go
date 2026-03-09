package domain

type CreateJobReq struct {
	Category string
	Fit      string
	Person   Image
	Garment  Image
}

type Image struct {
	Data        []byte
	ContentType string
}
