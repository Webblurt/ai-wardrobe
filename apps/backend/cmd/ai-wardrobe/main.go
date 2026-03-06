package main

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/app/deps"
	"ai-wardrobe/internal/app/http/server"
	"ai-wardrobe/internal/platform/logger"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config file", err)
	}

	logger := logger.NewLogger(cfg.LogLevel)

	server, err := server.New(deps.Deps{
		Logger: logger,
		Config: cfg,
	})
	if err != nil {
		logger.Fatal("Failed to create server: ", err)
	}

	logger.Info("Server starting on ", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, server); err != nil {
		logger.Fatal("While listening and serving HTTP: ", err)
	}
}
