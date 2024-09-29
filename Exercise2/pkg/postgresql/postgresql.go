package postgresql

import (
	"database/sql"
	"exercise2/internal/config"
	"exercise2/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var dbase *gorm.DB

func InitDatabase(configData *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configData.Host, configData.Port, configData.User, configData.Password, configData.DbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		models.User{},
		models.AdvancedGormUser{},
		models.AdvancedGormProfile{},
		models.RestUser{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}
	fmt.Println("Tables migrated successfully!")
	dbase = db
	return dbase
}

func GetDB(configData *config.Config) (*sql.DB, error) {
	if dbase == nil {
		dbase = InitDatabase(configData)
		var sleepTime = time.Duration(1)
		for dbase == nil {
			sleepTime = sleepTime * 2
			fmt.Printf("Database unavailable, trying again in %s", sleepTime)
			time.Sleep(sleepTime)
			dbase = InitDatabase(configData)
		}
	}
	sqlDB, err := dbase.DB()
	// Configure connection pooling
	sqlDB.SetMaxOpenConns(25)                 // Max number of open connections
	sqlDB.SetMaxIdleConns(25)                 // Max number of idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Max lifetime of a connection
	return sqlDB, err
}
