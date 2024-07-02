package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
	// "github.com/XanderMoroz/goBlog/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type AppEnvConfig struct {
	AppSecret  string
	Dbdriver   string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

// Указатель на БД
// (он будет осуществлять запросы)
// var DB *gorm.DB

// Извлкает переменные окружения и складывает в DBEnvConfig
func GetEnvConfig() *AppEnvConfig {

	log.Println("Извлекаем переменные окружения...")
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Не удалось извлечь .env")
	}
	log.Println("...успешно")

	return &AppEnvConfig{
		AppSecret: os.Getenv("API_SECRET"),

		Dbdriver:   os.Getenv("MONGO_DB_DRIVER"),
		DbUser:     os.Getenv("MONGO_ROOT_USERNAME"),
		DbPassword: os.Getenv("MONGO_ROOT_PASSWORD"),
		DbHost:     os.Getenv("MONGO_DB_HOST"),
		DbPort:     os.Getenv("MONGO_DB_PORT"),
		DbName:     os.Getenv("MONGO_DB_NAME"),
	}
}

// Подключается к БД
func ConnectToMongo() {

	envConfig := GetEnvConfig()

	DBURL := fmt.Sprintf(
		"%s://%s:%s@%s:%s/",
		envConfig.Dbdriver,
		envConfig.DbUser,
		envConfig.DbPassword,
		envConfig.DbHost,
		envConfig.DbPort,
	)
	log.Println("Подключаемся к БД:")
	log.Printf("Адрес БД: <%s> ...", DBURL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(DBURL),
	)

	defer func() {
		cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("mongodb disconnect error : %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("connection error :%v", err)
		return
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ping mongodb error :%v", err)
		return
	}
	fmt.Println("ping success")
}
