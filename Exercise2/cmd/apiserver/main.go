package main

import (
	"database/sql"
	"exercise2/internal/config"
	"exercise2/pkg/postgresql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB
var sqlDB *sql.DB

func main() {
	var configData = config.NewConfig()

	//init dbase
	db = postgresql.InitDatabase(configData)

	sqlDB, _ = postgresql.GetDB(configData)

	//connect dbase

	if err := configData.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	//test(configData)
	//test2(db)
	//test3()
	//advancedTest(configData)
	//advancedTest2(db)
	advancedTest3()
}
