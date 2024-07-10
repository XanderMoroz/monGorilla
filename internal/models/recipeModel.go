package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents user Recipe
type RecipeModel struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Title       string             `json:"title,omitempty"`
	Stages      []StageModel       `json:"stages,omitempty"`
	AuthorEmail string             `json:"author_email,omitempty"`
}

// Represents user model
type StageModel struct {
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Ingredients []IngredientModel `json:"ingredients,omitempty"`
}

// Represents user model
type IngredientModel struct {
	Subject   string `json:"subject,omitempty"`
	Condition string `json:"condition,omitempty"`
}

type RecipeCreateBody struct {
	Title  string       `json:"title,omitempty"`
	Stages []StageModel `json:"stages,omitempty"`
}

type RecipeUpdateBody struct {
	Title       string       `json:"title,omitempty"`
	Stages      []StageModel `json:"stages,omitempty"`
	AuthorEmail string       `json:"author_email,omitempty"`
}

// Represents Result of creating recipe response Body
type RecipeCreateResult struct {
	Result Result `json:"result"`
}
