package migrations

import (
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/data/db"
	"github.com/mohsen104/web-api/data/models"
	"github.com/mohsen104/web-api/pkg/logging"
)

var logger = logging.NewLogger(config.GetConfig())

func Up1() {
	database := db.GetDb()

	tables := []interface{}{}

	country := models.Country{}
	city := models.City{}

	if !database.Migrator().HasTable(country) {
		tables = append(tables, country)
	}

	if !database.Migrator().HasTable(city) {
		tables = append(tables, city)
	}

	database.Migrator().CreateTable(tables...)
	logger.Info(logging.Postgres, logging.Migration, "tables created", nil)
}

func Down1() {

}
