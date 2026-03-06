package deps

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/platform/logger"
)

type Deps struct {
	Logger *logger.Logger
	Config *config.Config
}
