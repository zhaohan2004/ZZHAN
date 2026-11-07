package main

import (
	"ZZHAN/internal/service"
	"ZZHAN/pkg/config"
	"ZZHAN/pkg/database"
	"ZZHAN/pkg/logger"
	"ZZHAN/pkg/storage"
	"context"
	"os"
	"os/signal"
)

func runCoverMigration(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	if err := logger.Init(&cfg.Log); err != nil {
		return err
	}
	db, err := database.InitMySQL(&cfg.Database.MySQL)
	if err != nil {
		return err
	}
	store, err := storage.New(cfg.Storage)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	return service.MigrateArticleCovers(ctx, db, store)
}
