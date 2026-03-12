package domain

type CreateJobReq struct {
	Provider    string
	Description string
	Category    string
	Steps       int
	Seed        int
	Autocrop    bool
	Upscale     int
	Upscaler    string
	Person      Image
	Garment     Image
}

type Image struct {
	Data        []byte
	ContentType string
}
