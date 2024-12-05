package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"shop/internal/config"
	"shop/models"
	"time"
)

var (
	dbase *gorm.DB
)

func InitRoles(db *gorm.DB) {
	var roles = []models.Role{
		{Name: "admin"},
		{Name: "manager"},
		{Name: "user"},
	}

	for _, role := range roles {
		if err := db.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
			log.Fatalf("Ошибка создания роли %s: %v", role.Name, err)
		}
	}
}

func InitDatabase() *gorm.DB {
	var configData = config.NewConfig()

	dsn := configData.DatabaseURL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Role{},
	)

	InitRoles(db)

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
