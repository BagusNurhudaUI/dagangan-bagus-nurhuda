package config

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// initialize database connection return gorm.DB
func DBInit() *gorm.DB {
	var err error
	portdb, _ := strconv.Atoi(GetEnv("portdb"))
	var (
		host     = GetEnv("POSTGRES_HOST")
		port     = portdb
		user     = GetEnv("POSTGRES_USER")
		password = GetEnv("POSTGRES_PASSWORD")
		dbname   = GetEnv("POSTGRES_DATABASE")
		db       *gorm.DB
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	fmt.Println("Connected to database")
	// err = db.AutoMigrate(models.Product{})
	// if err != nil {
	// 	log.Println("Failed to migrate, error: ", err)
	// }
	return db
}
