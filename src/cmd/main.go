package main

import (
	"github.com/mohsen104/web-api/api"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/data/cache"
	"github.com/mohsen104/web-api/data/db"
	"github.com/mohsen104/web-api/data/db/migrations"
	"github.com/mohsen104/web-api/pkg/logging"
)

func main() {
	cfg := config.GetConfig()

	logger := logging.NewLogger(cfg)

	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	err = db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	migrations.Up1()

	api.InitServer(cfg)
}
