package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/XanderMoroz/mongoMovies/database"
	"github.com/XanderMoroz/mongoMovies/internal/models"
)

const connectRemote = "mongodb://127.0.0.1:27017/authSource=admin"
const dbName = "netflix"
const colName = "watchedlist"

// Important
var collection *mongo.Collection

func checkNilError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitMongoDB() {
	// Client option
	clientOptions := options.Client().ApplyURI(connectRemote)

	//Connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection Done")

	collection = client.Database(dbName).Collection(colName)
	//Collection ready
	fmt.Println("Collection is ready")
}

//MongoDB helpers

//insert one record

func insertOneMovie(movie models.Netflix) {

	inserted, err := collection.InsertOne(context.Background(), movie)
	checkNilError(err)
	fmt.Println("Inserted one movie with ID:", inserted.InsertedID)
}

// update one record

func updateOneMovie(movieID string) {
	id, err := primitive.ObjectIDFromHex(movieID)
	checkNilError(err)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	checkNilError(err)
	fmt.Println("Modified Count:", result.ModifiedCount)
}

// delete one record
func deleteOneMovie(movieID string) {
	id, err := primitive.ObjectIDFromHex(movieID)
	checkNilError(err)
	filter := bson.M{"_id": id}
	delCount, err := collection.DeleteOne(context.Background(), filter)
	checkNilError(err)
	fmt.Println("Deleted Movie Count:", delCount)
}

//delete all record

func deleteAllMovie() int64 {
	delCount, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	checkNilError(err)
	fmt.Println("No of movies deleted:", delCount.DeletedCount)
	return delCount.DeletedCount

}

//get all movies from DB

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	checkNilError(err)

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		checkNilError(err)
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies

}

//Actual Controllers

func GetAlIMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

// @Summary        create new movie
// @Description    Creating Movie in DB with given request body
// @Tags           Movies
// @Accept         json
// @Produce        json
// @Param          request         	body        models.AddMovieBody    true    "Введите фильм"
// @Success        201              {string}    string
// @Failure        400              {string}    string    "Bad Request"
// @Router         /api/movie 			[post]
func CreateMovie(w http.ResponseWriter, r *http.Request) {

	log.Println("Поступил запрос на создание новой записи в БД...")
	var movie models.Netflix
	log.Println("Извлекаю тело запроса...")
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Printf("При извлечении тела запроса - Произошла ошибка: <%v>\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		log.Println("...успешно")
		log.Printf("Тело запроса: %+v", movie)
	}
	// moviesCollection := database.MongoDB.Collection("movie_collection")
	// moviesCollection.Drop(database.MongoCtx)

	insertedDocument := bson.M{
		"name":       "Царица",
		"content":    "test content",
		"bank_money": 1000,
		"create_at":  time.Now(),
	}
	insertedResult, err := database.MongoCollection.InsertOne(
		database.MongoCtx,
		insertedDocument,
	)
	// inserted, err := database.MongoCollection.InsertOne(context.Background(), movie)
	// checkNilError(err)
	// fmt.Println("Inserted one movie with ID:", inserted.InsertedID)

	if err != nil {
		log.Fatalf("inserted error : %v", err)
		return
	}
	fmt.Println("======= inserted id ================")
	log.Printf("inserted ID is : %v", insertedResult.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params)

}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")
	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)

}
