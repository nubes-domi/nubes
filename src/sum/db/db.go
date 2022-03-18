package db

import (
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	handle *gorm.DB
}

var DB Database

func databasePath() string {
	store := os.Getenv("NUBES_STORE")
	if os.Getenv("NUBES_DEV") != "" || store == "" {
		store = "sum.sqlite"
	} else {
		store = path.Join(store, "sql", "sum.sqlite")
	}

	return store
}

func InitDatabase() {
	config := gorm.Config{}
	if os.Getenv("NUBES_DEV") != "" {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	handle, err := gorm.Open(sqlite.Open(databasePath()), &config)
	if err != nil {
		log.Panicf("Could not open database")
	}

	DB := Database{handle}
	Migrate(DB.handle)
}
