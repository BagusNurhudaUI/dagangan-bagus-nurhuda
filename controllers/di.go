package controllers

import (
	"gorm.io/gorm"
)

// struct for database injection
type InDB struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *InDB {
	return &InDB{
		DB: db,
	}
}
