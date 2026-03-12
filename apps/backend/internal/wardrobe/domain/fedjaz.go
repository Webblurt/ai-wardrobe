package domain

type TryOnMode string

const (
	TryOnModeDefault TryOnMode = "fedjaz"
	TryOnModeFashn   TryOnMode = "fedjaz_fashn_v1.5"
)

type TryOnParams struct {
	Mode TryOnMode

	Description string
	Category    string
	Steps       int
	Seed        int
	Autocrop    bool
	Upscale     int
	Upscaler    string

	// fashn params
	GarmentPhotoType string
	NumSamples       int
	NumTimesteps     int
	GuidanceScale    float32
	SegmentationFree bool
}
