package database

import (
	"log"

	"github.com/biskitsx/Refresh-Token-Go/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.User{}, &model.Session{})
	return db
}
