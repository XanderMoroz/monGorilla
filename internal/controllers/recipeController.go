package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/XanderMoroz/mongoMovies/configs"
	"github.com/XanderMoroz/mongoMovies/internal/models"
	"github.com/XanderMoroz/mongoMovies/internal/utils"
	"github.com/gorilla/mux"
	// "github.com/XanderMoroz/mongoMovies/internal/utils"
)

var recipeCollection *mongo.Collection = configs.GetCollection(configs.DB, "recipes")

// @Summary     create new recipe
// @Description Creating Recipe in DB with given request body
// @Tags        Recipes
// @ID			create-new-recipe
// @Accept      json
// @Produce     json
// @Param       request         	body        models.RecipeCreateBody    true    "Enter recipe data"
// @Success     201              	{string}    string
// @Failure     400              	{string}    string    "Bad Request"
// @Security  	Bearer
// @Router      /api/recipes 		[post]
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	CurrentUserResult := new(models.CurrentUserResult)

	log.Println("Валидируем токен из заголовка...")
	jwtFromHeader := r.Header.Get("Authorization")
	userId, err := utils.ParseUserIDFromJWTToken(jwtFromHeader)
	if err != nil {
		log.Println("При извлечении userId произошла ошибка")
		CurrentUserResult.Result.Success = false
		CurrentUserResult.Result.ErrorCode = utils.ERR0304
		CurrentUserResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		CurrentUserResult.Result.ErrorException = utils.ExceptionToString(err)
		json.NewEncoder(w).Encode(CurrentUserResult)
	} else {
		log.Println("...успешно")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.CurrentUserModel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err = userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		log.Printf("При извлечении записи -произошла ошибка: <%v>\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		log.Printf("Запись успешно извлечена: <%+v>\n", user)

		log.Println("Поступил запрос на создание новой записи в БД...")
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var recipeBody models.RecipeCreateBody
		defer cancel()

		err := json.NewDecoder(r.Body).Decode(&recipeBody)
		if err != nil {
			log.Printf("При извлечении тела запроса - Произошла ошибка: <%v>\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("...успешно")
			// log.Printf("Тело запроса: %+v", user)
		}

		newRecipe := models.RecipeModel{
			Id:          primitive.NewObjectID(),
			Title:       recipeBody.Title,
			Stages:      recipeBody.Stages,
			AuthorEmail: user.Email,
		}

		result, err := recipeCollection.InsertOne(ctx, newRecipe)
		if err != nil {
			log.Printf("При добавлении новой записи - Произошла ошибка: <%v>\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("Новая запись успешно добавлена:")
			log.Printf("ID новой записи: %v", result.InsertedID)
		}

		json.NewEncoder(w).Encode(newRecipe)
	}
}

// @Summary		get all my recipes
// @Description Get all recipes of authenticated user
// @Tags 		Recipes
// @ID			get-all-recipes-of-current-user
// @Produce		json
// @Success		200		{object}	[]models.RecipeModel
// @Security  	Bearer
// @Router		/api/recipes [get]
func GetAllMyRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	CurrentUserResult := new(models.CurrentUserResult)

	log.Println("Валидируем токен из заголовка...")
	jwtFromHeader := r.Header.Get("Authorization")
	userId, err := utils.ParseUserIDFromJWTToken(jwtFromHeader)
	if err != nil {
		log.Println("При извлечении userId произошла ошибка")
		CurrentUserResult.Result.Success = false
		CurrentUserResult.Result.ErrorCode = utils.ERR0304
		CurrentUserResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		CurrentUserResult.Result.ErrorException = utils.ExceptionToString(err)
		json.NewEncoder(w).Encode(CurrentUserResult)
	} else {
		log.Println("...успешно")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.CurrentUserModel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err = userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		log.Printf("При извлечении записи -произошла ошибка: <%v>\n", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		log.Printf("Запись успешно извлечена: <%+v>\n", user)

		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var userRecipes []models.RecipeModel
		// defer cancel()

		results, err := recipeCollection.Find(ctx, bson.M{"authoremail": user.Email})

		if err != nil {
			log.Printf("При извлечении списка записей - произошла ошибка: <%v>\n", err.Error())
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleRecipe models.RecipeModel
			if err = results.Decode(&singleRecipe); err != nil {
				log.Printf("При обработке списка записей -произошла ошибка: <%v>\n", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			userRecipes = append(userRecipes, singleRecipe)
		}

		json.NewEncoder(w).Encode(userRecipes)
	}
}

// @Summary		get a recipe by ID
// @Description Get a recipe by ID
// @Tags 		Recipes
// @ID			get-recipe-by-id
// @Produce		json
// @Param		id					path		string			true	"RecipeID"
// @Success		200					{object}	models.RecipeModel
// @Failure		404					{object}	[]string
// @Router		/api/recipes/{id} 	[get]
func GetRecipeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	log.Printf("Поступил запрос на извлечение записи по ID: <%s>\n", params["id"])

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := params["id"]
	var recipe models.RecipeModel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := recipeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&recipe)
	if err != nil {
		log.Printf("При извлечении записи -произошла ошибка: <%v>\n", err.Error())
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	} else {
		log.Printf("Рецепт успешно извлечен: <%+v>\n", recipe)
	}
	json.NewEncoder(w).Encode(recipe)
}

// // @Summary			update user by ID
// // @Description 	Update user by ID
// // @ID				update-user-by-id
// // @Tags 			Users
// // @Produce			json
// // @Param			id					path		string									true	"UserID"
// // @Param           request         	body        models.CreateUserBody    true    	"Введите новые данные пользователя"
// // @Success			200	{object}		[]string
// // @Failure			404	{object}		[]string
// // @Router			/api/users/{id} 	[put]
// func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

// 	params := mux.Vars(r)
// 	log.Printf("Поступил запрос на обновление записи по ID: <%s>\n", params["id"])
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := params["id"]
// 	var user models.User
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	//validate the request body
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		log.Printf("При извлечении тела запроса - Произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	} else {
// 		log.Println("...успешно")
// 	}

// 	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

// 	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
// 	if err != nil {
// 		log.Printf("При обновлении записи -произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	//get updated user details
// 	var updatedUser models.User
// 	if result.MatchedCount == 1 {
// 		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
// 		if err != nil {
// 			log.Printf("При извлечении записи -произошла ошибка: <%v>\n", err.Error())
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}
// 	}

// 	json.NewEncoder(w).Encode(updatedUser)
// }

// // @Summary		delete a user by ID
// // @Description Delete a user by ID
// // @ID			delete-user-by-id
// // @Tags 		Users
// // @Produce		json
// // @Param		id					path		string		true	"UserID"
// // @Success		200					{object}	[]string
// // @Failure		404					{object}	[]string
// // @Router		/api/users/{id} 	[delete]
// func DeleteBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

// 	params := mux.Vars(r)
// 	log.Printf("Поступил запрос на удаление записи по ID: <%s>\n", params["id"])
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := params["id"]
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
// 	if err != nil {
// 		log.Printf("При удалении записи - произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if result.DeletedCount < 1 {
// 		log.Printf("При извлечении тела запроса - Произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	json.NewEncoder(w).Encode("User successfully deleted!")
// }

// //MongoDB helpers
// // func checkNilError(err error) {
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // }

// // func deleteAllMovie() int64 {
// // 	delCount, err := userCollection.DeleteMany(context.Background(), bson.D{{}}, nil)
// // 	checkNilError(err)
// // 	fmt.Println("No of movies deleted:", delCount.DeletedCount)
// // 	return delCount.DeletedCount
// // }
