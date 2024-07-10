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
