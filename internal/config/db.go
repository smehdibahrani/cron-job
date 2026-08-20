package config

import (
	"cron_job/internal/entity"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := GetString("DB_HOST", "localhost")
	port := GetInt("DB_PORT", 5432)
	name := GetString("DB_NAME", "cronjob")
	user := GetString("DB_USER", "postgres")
	pass := GetString("DB_PASS", "postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable", host, user, pass, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Job{},
		&entity.Group{},
		&entity.Notification{},
		&entity.RequestHttp{},
		&entity.Log{})
	if err != nil {
		fmt.Println("failed to migrate database", err.Error())
	}
}
