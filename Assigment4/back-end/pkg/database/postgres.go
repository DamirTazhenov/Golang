package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"taskmanager/internal/config"
	"taskmanager/models"
	"time"
)

var (
	dbase *gorm.DB
)

func InitDatabase() *gorm.DB {
	var configData = config.NewConfig()

	dsn := configData.DatabaseURL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Task{},
		&models.TeamMember{},
	)

	return db
}

func GetDB() *gorm.DB {
	if dbase == nil {
		dbase = InitDatabase()
		var sleepTime = time.Duration(1)
		for dbase == nil {
			sleepTime = sleepTime * 2
			fmt.Printf("Database unavailable, trying again in %s", sleepTime)
			time.Sleep(sleepTime)
			dbase = InitDatabase()
		}
	}
	return dbase
}
