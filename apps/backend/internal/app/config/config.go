package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App        `yaml:"app"`
	Server     `yaml:"server"`
	Storage    `yaml:"storage"`
	Replicate  `yaml:"replicate"`
	FedjazVton `yaml:"fedjaz_vton"`
	LogLevel   string `yaml:"log_level"` // possible variants: trace, debug, info, warn, error, fatal, success
}

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	URL     string `yaml:"url"`
}

type Server struct {
	Port string `yaml:"port"`
	Cors `yaml:"cors"`
	// Handlers `yaml:"handlers"`
}

type Cors struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AlloweMethods  []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type Storage struct {
	BaseDir    string `yaml:"base_dir"`
	UploadsDir string `yaml:"uploads_dir"`
	ResultsDir string `yaml:"results_dir"`
	JobsDir    string `yaml:"jobs_dir"`
}

type Replicate struct {
	Token        string `yaml:"token"`
	BaseURL      string `yaml:"base_url"`
	ModelVersion string `yaml:"model_version"`
	PollInterval int    `yaml:"poll_interval"`
	MaxWaitSec   int    `yaml:"max_wait_sec"`
}

type FedjazVton struct {
	Token        string `yaml:"token"`
	BaseURL      string `yaml:"base_url"`
	ModelVersion string `yaml:"model_version"`
	PollInterval int    `yaml:"poll_interval"`
	MaxWaitSec   int    `yaml:"max_wait_sec"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, errors.New("config path is empty")
	}

	if _, err := os.Stat(path); err != nil {
		return nil, errors.New("config file does not exist: " + path)
	}

	var cfg Config

	// Load YAML file
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	// Override with environment variables
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	// err := validate(&cfg)
	// if err != nil {
	// 	return nil, err
	// }

	return &cfg, nil
}
