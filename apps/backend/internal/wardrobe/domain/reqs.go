package domain

type CreateJobReq struct {
	Provider         string
	Description      string
	Category         string
	Steps            int
	Seed             int
	Autocrop         bool
	Upscale          int
	Upscaler         string
	Person           Image
	Garment          Image
	GarmentPhotoType string
	NumSamples       int
	NumTimesteps     int
	GuidanceScale    float32
	SegmentationFree bool
}

type Image struct {
	Data        []byte
	ContentType string
}
