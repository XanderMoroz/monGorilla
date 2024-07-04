package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {

	envConfig := GetEnvConfig()

	DBURL := fmt.Sprintf(
		"%s://%s:%s@%s:%s/",
		envConfig.Dbdriver,
		envConfig.DbUser,
		envConfig.DbPassword,
		envConfig.DbHost,
		envConfig.DbPort,
	)
	// log.Println("Подключаемся к БД...")
	// log.Printf("Адрес БД: <%s>:", DBURL)
	// log.Printf("Название БД: <%s>:", envConfig.DbName)

	client, err := mongo.NewClient(options.Client().ApplyURI(DBURL))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golangAPI").Collection(collectionName)
	return collection
}
