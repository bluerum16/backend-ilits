package config

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"backend-ilits/models"
)

func ConnectDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("akta.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal konek ke database: ", err)
	}

	if err := db.AutoMigrate(&models.AktaKelahiran{}); err != nil {
		log.Fatal("gagal migrate: ", err)
	}

	return db
}
