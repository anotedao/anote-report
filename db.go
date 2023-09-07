package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDb() *gorm.DB {
	var db *gorm.DB
	var err error
	dbconf := gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	dbconf.Logger = logger.Default.LogMode(logger.Error)

	db, err = gorm.Open(postgres.Open(conf.DSN), &dbconf)

	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if err := db.AutoMigrate(&Address{}, &KeyValue{}); err != nil {
		panic(err.Error())
	}

	return db
}
