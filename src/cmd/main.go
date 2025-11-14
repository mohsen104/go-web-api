package main

import (
	"github.com/mohsen104/web-api/api"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/data/cache"
	"github.com/mohsen104/web-api/data/db"
	"github.com/mohsen104/web-api/pkg/logging"
)

func main() {
	cfg := config.GetConfig()

	logger := logging.NewLogger(cfg)

	api.InitServer(cfg)

	defer cache.CloseRedis()
	err := cache.InitRedis(cfg)
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	defer db.CloseDb()
	err = db.InitDb(cfg)
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
}
