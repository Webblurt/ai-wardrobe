package domain

type TryOnParams struct {
	Description string
	Category    string
	Steps       int
	Seed        int
	Autocrop    bool
	Upscale     int
	Upscaler    string
}
