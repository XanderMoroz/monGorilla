package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents user Recipe
type RecipeModel struct {
	Id     primitive.ObjectID `json:"id,omitempty"`
	Title  string             `json:"title,omitempty"`
	Stages []StageModel       `json:"stages,omitempty"`
}

// Represents user model
type StageModel struct {
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Ingredients []IngredientModel `json:"ingredients,omitempty"`
}

// Represents user model
type IngredientModel struct {
	Subject   string `json:"password,omitempty"`
	Condition string `json:"first_name,omitempty"`
}
