package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppEnvConfig struct {
	AppSecret string
	AppPort   string

	Dbdriver   string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

// Извлкает переменные окружения и складывает в DBEnvConfig
func GetEnvConfig() *AppEnvConfig {

	// log.Print("Извлекаем переменные окружения...")
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Не удалось извлечь .env")
	}
	// log.Print("...успешно")

	return &AppEnvConfig{
		AppSecret: os.Getenv("API_SECRET"),
		AppPort:   os.Getenv("APP_PORT"),

		Dbdriver:   os.Getenv("MONGO_DB_DRIVER"),
		DbUser:     os.Getenv("MONGO_ROOT_USERNAME"),
		DbPassword: os.Getenv("MONGO_ROOT_PASSWORD"),
		DbHost:     os.Getenv("MONGO_DB_HOST"),
		DbPort:     os.Getenv("MONGO_DB_PORT"),
		DbName:     os.Getenv("MONGO_DB_NAME"),
	}
}
