package config

import (
	"fmt"
	"log"
	"os"

	"github.com/mca93/qrcode_service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// 👉 AutoMigrate para criar tabelas se não existirem
	err = db.AutoMigrate(&models.ClientApp{}, &models.ApiKey{})
	// &models.Template{}, &models.QRCode{})
	if err != nil {
		log.Fatal("failed to migrate tables: ", err)
	}
	DB = db
}
