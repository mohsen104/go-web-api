package main

import (
	"log"

	"github.com/mohsen104/web-api/api"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/data/cache"
	"github.com/mohsen104/web-api/data/db"
)

func main() {
	cfg := config.GetConfig()

	api.InitServer(cfg)

	defer cache.CloseRedis()
	err := cache.InitRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.CloseDb()
	err = db.InitDb(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
