package config

import (
	"fmt"
	"go-api-article/models"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	connection := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", ENV.DB_HOST, ENV.DB_USER, ENV.DB_PASSWORD, ENV.DB_DATABASE, ENV.DB_PORT)
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})

	if err != nil {
		panic("Koneksi Ke Database Gagal")
	}

	db.AutoMigrate(&models.Category{}, &models.Article{})

	DB = db
	log.Println("Database Berhasil Dihubungkan")
}
