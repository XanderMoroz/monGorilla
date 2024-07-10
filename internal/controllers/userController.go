package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	// "reflect"

	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/XanderMoroz/mongoMovies/configs"
	"github.com/XanderMoroz/mongoMovies/internal/models"
	"github.com/XanderMoroz/mongoMovies/internal/utils"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// Register		 RegisterAccount godoc
// @Summary      Create a account
// @Description  Register and create account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        userModelArgs 			body 			models.UserRegisterArgs true "UserRegister"
// @Success      200  					{object}  		models.UserRegisterResult
// @Router       /api/users/register 	[post]
func Register(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	return fmt.Errorf("method not allowed %s", r.Method)
	// }

	log.Println("Поступил запрос на регистрацию пользователя")

	userRegisterArgs := new(models.UserRegisterArgs)
	userRegisterResult := new(models.UserRegisterResult)

	log.Println("Извлекаем тело запроса...")
	if err := json.NewDecoder(r.Body).Decode(userRegisterArgs); err != nil {
		userRegisterResult.Result.Success = false
		userRegisterResult.Result.ErrorCode = utils.ERR0303
		userRegisterResult.Result.ErrorDescription = utils.ERR0303.ToDescription()
		userRegisterResult.Result.ErrorException = utils.ExceptionToString(err)
		log.Println("При извлечении тела запроса - произошла ошибка:", userRegisterResult)
		json.NewEncoder(w).Encode(userRegisterResult)
		return
	} else {
		log.Println("...успешно")
		log.Println("Тело запроса", userRegisterArgs)
	}

	log.Println("Валидируем тело запроса...")
	if ok := utils.ValidateCheckSpaceCharacter(
		userRegisterArgs.FirstName,
		userRegisterArgs.LastName,
		userRegisterArgs.Email,
		userRegisterArgs.Password,
		userRegisterArgs.ValidatePassword,
		userRegisterArgs.PhoneNumber,
	); !ok {
		userRegisterResult.Result.Success = false
		userRegisterResult.Result.ErrorCode = utils.ERR0304
		userRegisterResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		log.Println("При валидации тела запроса - произошла ошибка:", userRegisterResult)
		json.NewEncoder(w).Encode(userRegisterResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Println("Сравниваем пароли...")
	if ok := utils.IsStringEqual(userRegisterArgs.Password, userRegisterArgs.ValidatePassword); !ok {
		userRegisterResult.Result.Success = false
		userRegisterResult.Result.ErrorCode = utils.ERR0406
		userRegisterResult.Result.ErrorDescription = utils.ERR0406.ToDescription()
		json.NewEncoder(w).Encode(userRegisterResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Println("Валидируем Email...")
	if ok := utils.ValidateEmail(userRegisterArgs.Email); !ok {
		userRegisterResult.Result.Success = false
		userRegisterResult.Result.ErrorCode = utils.ERR0304
		userRegisterResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		json.NewEncoder(w).Encode(userRegisterResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Println("Валидируем пароль...")
	if ok := utils.ValidatePassword(userRegisterArgs.Password); !ok {
		userRegisterResult.Result.Success = false
		userRegisterResult.Result.ErrorCode = utils.ERR0304
		userRegisterResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		json.NewEncoder(w).Encode(userRegisterResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Printf("Поступил запрос на извлечение записи по e-mail: <%s>\n", userRegisterArgs.Email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// userId := params["id"]
	var user models.UserModel
	defer cancel()
	userEmail := userRegisterArgs.Email
	// objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		log.Printf("Результат: <%v>\n", err.Error())
		// http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		log.Printf("Такой email уже есть в БД: <%+v>\n", user)
		return
	}

	hashed_password, err := utils.HashPassword(userRegisterArgs.Password)
	if err != nil {
		userRegisterResult.Result.Success = false
		userRegisterResult.Result.ErrorCode = utils.ERR0405
		userRegisterResult.Result.ErrorDescription = utils.ERR0405.ToDescription()
		userRegisterResult.Result.ErrorException = utils.ExceptionToString(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	to_register := new(models.UserModel)
	to_register.Id = primitive.NewObjectID()
	to_register.FirstName = userRegisterArgs.FirstName
	to_register.LastName = userRegisterArgs.LastName
	to_register.Email = userRegisterArgs.Email
	to_register.PhoneNumber = userRegisterArgs.PhoneNumber
	to_register.Password = hashed_password

	result, err := userCollection.InsertOne(ctx, to_register)
	if err != nil {
		log.Printf("При добавлении новой записи - Произошла ошибка: <%v>\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		log.Println("Новая запись успешно добавлена:")
		log.Printf("ID новой записи: %v", result.InsertedID)
	}

	json.NewEncoder(w).Encode(to_register)
}

// Login   LoginAccount godoc
// @Summary      Login to your account
// @Description  Login with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        userModelArgs 	body 		models.UserLoginArgs 	true 	"UserLogin"
// @Success      200  			{object}  	models.UserLoginResult
// @Router       /api/users/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	// // Валидируем метод запроса
	// if r.Method != http.MethodPost {
	// 	return fmt.Errorf("method not allowed %s", r.Method)
	// }

	userLoginResult := new(models.UserLoginResult)
	userLoginArgs := new(models.UserLoginArgs)

	log.Println("Извлекаем тело запроса...")
	if err := json.NewDecoder(r.Body).Decode(userLoginArgs); err != nil {
		userLoginResult.Result.Success = false
		userLoginResult.Result.ErrorCode = utils.ERR0303
		userLoginResult.Result.ErrorDescription = utils.ERR0303.ToDescription()
		userLoginResult.Result.ErrorException = utils.ExceptionToString(err)
		json.NewEncoder(w).Encode(userLoginResult)
		return
	} else {
		log.Println("...успешно")
		log.Println("Тело запроса", userLoginArgs)
	}

	log.Println("Валидируем тело запроса...")
	if ok := utils.ValidateCheckSpaceCharacter(userLoginArgs.Email, userLoginArgs.Password); !ok {
		userLoginResult.Result.Success = false
		userLoginResult.Result.ErrorCode = utils.ERR0304
		userLoginResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		json.NewEncoder(w).Encode(userLoginResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Println("Валидируем Email...")
	if ok := utils.ValidateEmail(userLoginArgs.Email); !ok {
		userLoginResult.Result.Success = false
		userLoginResult.Result.ErrorCode = utils.ERR0304
		userLoginResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		json.NewEncoder(w).Encode(userLoginResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Println("Валидируем пароль...")
	if ok := utils.ValidatePassword(userLoginArgs.Password); !ok {
		userLoginResult.Result.Success = false
		userLoginResult.Result.ErrorCode = utils.ERR0304
		userLoginResult.Result.ErrorDescription = utils.ERR0304.ToDescription()
		json.NewEncoder(w).Encode(userLoginResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Printf("Поступил запрос на извлечение записи по e-mail: <%s>\n", userLoginArgs.Email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.UserModel
	defer cancel()
	userEmail := userLoginArgs.Email

	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		log.Printf("Результат: <%v>\n", err.Error())
	} else {
		log.Printf("Такой email уже есть в БД: <%+v>\n", user)
	}

	log.Println("Сравниваем пароль и хэш...")
	if ok := utils.CompareHashAndPassword(user.Password, userLoginArgs.Password); !ok {
		userLoginResult.Result.Success = false
		userLoginResult.Result.ErrorCode = utils.ERR0403
		userLoginResult.Result.ErrorDescription = utils.ERR0403.ToDescription()
		userLoginResult.Result.ErrorException = utils.ExceptionToString(err)
		json.NewEncoder(w).Encode(userLoginResult)
		return
	} else {
		log.Println("...успешно")
	}

	log.Println("Создаем токен...")
	token, err := utils.CreateJSONWebToken(user.Id)
	if err != nil {
		userLoginResult.Result.Success = false
		userLoginResult.Result.ErrorCode = utils.ERR0405
		userLoginResult.Result.ErrorDescription = utils.ERR0405.ToDescription()
		userLoginResult.Result.ErrorException = utils.ExceptionToString(err)
		json.NewEncoder(w).Encode(userLoginResult)
		return
	} else {
		log.Println("...успешно")
	}

	userLoginResult.Id = user.Id
	userLoginResult.AuthenticationToken = token
	userLoginResult.UserInfos = map[string]string{
		"full_name":    user.FirstName + " " + user.LastName,
		"phone_number": user.PhoneNumber,
		"email":        user.Email,
	}

	userLoginResult.Result.Success = true
	userLoginResult.Result.ErrorCode = ""
	userLoginResult.Result.ErrorDescription = ""
	userLoginResult.Result.ErrorException = ""

	json.NewEncoder(w).Encode(userLoginResult)
}

// CurrentUser  	GetCurrentUser godoc
// @Summary      	Check validity of token
// @Description  	Token check method for authentication
// @Tags         	Auth
// @Accept       	json
// @Produce      	json
// @Success      	200  				{object}  	models.TokenCheckResult
// @Security  		Bearer
// @Router       	/api/users/current_user [get]
func CurrentUser(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	//  return fmt.Errorf("method not allowed %s", r.Method)
	// }

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

	}
	CurrentUserResult.AuthenticationToken = jwtFromHeader
	CurrentUserResult.CurrentUser = user
	CurrentUserResult.Result.Success = true
	CurrentUserResult.Result.ErrorCode = ""
	CurrentUserResult.Result.ErrorDescription = ""
	CurrentUserResult.Result.ErrorException = ""
	json.NewEncoder(w).Encode(user)
}

// // @Summary        create new user
// // @Description    Creating User in DB with given request body
// // @Tags           Users
// // @ID				create-new-user
// // @Accept         json
// // @Produce        json
// // @Param          request         		body        models.CreateUserBody    true    "Enter user data"
// // @Success        201              	{string}    string
// // @Failure        400              	{string}    string    "Bad Request"
// // @Router         /api/users 			[post]
// func CreateUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Allow-Control-Allow-Methods", "POST")

// 	log.Println("Поступил запрос на создание новой записи в БД...")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var user models.User
// 	defer cancel()

// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		log.Printf("При извлечении тела запроса - Произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	} else {
// 		log.Println("...успешно")
// 		// log.Printf("Тело запроса: %+v", user)
// 	}

// 	newUser := models.User{
// 		Id:       primitive.NewObjectID(),
// 		Name:     user.Name,
// 		Location: user.Location,
// 		Title:    user.Title,
// 	}

// 	result, err := userCollection.InsertOne(ctx, newUser)
// 	if err != nil {
// 		log.Printf("При добавлении новой записи - Произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	} else {
// 		log.Println("Новая запись успешно добавлена:")
// 		log.Printf("ID новой записи: %v", result.InsertedID)
// 	}

// 	json.NewEncoder(w).Encode(newUser)
// }

// // @Summary		get a user by ID
// // @Description Get a user by ID
// // @Tags 		Users
// // @ID			get-user-by-id
// // @Produce		json
// // @Param		id					path		string			true	"UserID"
// // @Success		200					{object}	models.User
// // @Failure		404					{object}	[]string
// // @Router		/api/users/{id} 	[get]
// func GetUserByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Allow-Control-Allow-Methods", "GET")

// 	params := mux.Vars(r)
// 	log.Printf("Поступил запрос на извлечение записи по ID: <%s>\n", params["id"])

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := params["id"]
// 	var user models.User
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
// 	if err != nil {
// 		log.Printf("При извлечении записи -произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	} else {
// 		log.Printf("Запись успешно извлечена: <%+v>\n", user)
// 	}
// 	json.NewEncoder(w).Encode(user)
// }

// // @Summary		get all users
// // @Description Get all users from db
// // @Tags 		Users
// // @ID			get-all-users
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @Produce		json
// // @Success		200		{object}	[]models.User
// // @Router		/api/users [get]
// func GetAllUsers(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Allow-Control-Allow-Methods", "GET")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var users []models.User
// 	defer cancel()

// 	results, err := userCollection.Find(ctx, bson.M{})

// 	if err != nil {
// 		log.Printf("При извлечении списка записей - произошла ошибка: <%v>\n", err.Error())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	//reading from the db in an optimal way
// 	defer results.Close(ctx)
// 	for results.Next(ctx) {
// 		var singleUser models.User
// 		if err = results.Decode(&singleUser); err != nil {
// 			log.Printf("При обработке списка записей -произошла ошибка: <%v>\n", err.Error())
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		users = append(users, singleUser)
// 	}

// 	json.NewEncoder(w).Encode(users)
// }

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
