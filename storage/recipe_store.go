package storage

import (
	"errors"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Recipe struct {
	ID          int64    `json:"id"`
	CreatorID   int64    `json:"creatorId"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
	Title       string   `json:"title"`
}

var (
	ErrRecipeNotFound = errors.New("recipe not found")
)

type RecipeDatastore interface {
	CreateRecipe(recipe Recipe) (Recipe, error)
	ReadRecipe(id int64) (Recipe, error)
	ListRecipes() ([]Recipe, error)
	UpdateRecipe(id int64, recipe Recipe) error
	DeleteRecipe(id int64) error
}

type UserDatastore interface {
	CreateUser(user User) error
	ReadUser(email string) (User, error)
}
