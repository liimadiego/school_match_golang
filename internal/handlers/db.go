package handlers

import (
	"gorm.io/gorm"

	"github.com/liimadiego/schoolmatch/internal/database"
)

var (
	testDB *gorm.DB
)

func GetDB() *gorm.DB {
	if testDB != nil {
		return testDB
	}
	return database.DB
}

func SetDB(db *gorm.DB) {
	testDB = db
}
