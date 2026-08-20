package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsBool, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return valAsBool
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}
