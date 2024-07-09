package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	Title    string             `json:"title,omitempty" validate:"required"`
}

type CreateUserBody struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Location string `json:"location,omitempty" validate:"required"`
	Title    string `json:"title,omitempty" validate:"required"`
}

type UserModel struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Password    string             `json:"password,omitempty"`
	FirstName   string             `json:"first_name,omitempty"`
	LastName    string             `json:"last_name,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty"`
	Email       string             `json:"email,omitempty"`
	// BirthDate   time.Time `json:"birth_date,omitempty"`
}

// USER LOGIN ARGS & RESULT
type UserLoginArgs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResult struct {
	Id                  primitive.ObjectID `json:"id"`
	UserInfos           map[string]string  `json:"user_infos"`
	AuthenticationToken string             `json:"authentication_token"`
	Result              Result             `json:"result"`
}

// -------------------------

type UserRegisterArgs struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	PhoneNumber      string `json:"phone_number"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	ValidatePassword string `json:"validate_password"`
	// BirthDate        time.Time `json:"birth_date"`
}

type UserRegisterResult struct {
	Result Result `json:"result"`
}
