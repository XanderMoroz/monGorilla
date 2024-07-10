package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents user model
type UserModel struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Password    string             `json:"password,omitempty"`
	FirstName   string             `json:"first_name,omitempty"`
	LastName    string             `json:"last_name,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty"`
	Email       string             `json:"email,omitempty"`
	// BirthDate   time.Time `json:"birth_date,omitempty"`
}

// ------------ USER REGISTER ARGS & RESULT -------------

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

// ------------ USER LOGIN ARGS & RESULT -------------

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

// ------------ CURRENT USER ARGS & RESULT -------------

type CurrentUserModel struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	FirstName   string             `json:"first_name,omitempty"`
	LastName    string             `json:"last_name,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty"`
	Email       string             `json:"email,omitempty"`
}

type CurrentUserResult struct {
	CurrentUser         CurrentUserModel `json:"current_user"`
	AuthenticationToken string           `json:"authentication_token"`
	Result              Result           `json:"result"`
}
