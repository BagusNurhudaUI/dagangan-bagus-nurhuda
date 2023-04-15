package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/BagusNurhudaUI/dagangan-bagus-nurhuda/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func DBInit() {
	var err error
	portdb, _ := strconv.Atoi(GetEnv("portdb"))
	var (
		host     = GetEnv("host")
		port     = portdb
		user     = GetEnv("user")
		password = GetEnv("password")
		dbname   = GetEnv("dbname")
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	fmt.Println("Connected to database")
	err = db.AutoMigrate(models.Product{})
	if err != nil {
		return
	}
}

func GetDB() *gorm.DB {
	return db
}
