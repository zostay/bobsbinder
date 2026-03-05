package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/zostay/bobsbinder/internal/config"
	"github.com/zostay/bobsbinder/internal/db"
	"github.com/zostay/bobsbinder/internal/router"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	database, err := db.Connect(cfg, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	if _, err := db.RunMigrations(database, logger); err != nil {
		logger.Warn("failed to run migrations", zap.Error(err))
	}

	r := router.New(database, cfg, logger)

	addr := ":" + cfg.APIPort
	logger.Info("starting server", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
