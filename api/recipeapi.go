package api

import (
	"context"
	"sync"
)

type User struct {
	Id       int64
	Email    string
	Password string
}

type RecipeStore struct {
	//Temporary in memory stores
	Users        map[int64]User
	Recipes      map[int64]Recipe
	NextUserId   int64
	NextRecipeId int64
	Lock         sync.Mutex
}

var _ StrictServerInterface = (*RecipeStore)(nil)

func NewRecipeStore() *RecipeStore {
	return &RecipeStore{
		Users:        make(map[int64]User),
		Recipes:      make(map[int64]Recipe),
		NextUserId:   1000,
		NextRecipeId: 1,
	}
}

func (p *RecipeStore) CreateUser(context.Context, CreateUserRequestObject) (CreateUserResponseObject, error) {
	return CreateUser201Response{}, nil
}
func (p *RecipeStore) LoginUser(context.Context, LoginUserRequestObject) (LoginUserResponseObject, error) {
	return LoginUser401Response{}, nil
}

func (p *RecipeStore) CreateRecipe(context.Context, CreateRecipeRequestObject) (CreateRecipeResponseObject, error) {
	return CreateRecipe201JSONResponse{}, nil
}

func (p *RecipeStore) ListRecipes(context.Context, ListRecipesRequestObject) (ListRecipesResponseObject, error) {
	return ListRecipes200JSONResponse{}, nil
}

func (p *RecipeStore) GetRecipeById(_ context.Context, request GetRecipeByIdRequestObject) (GetRecipeByIdResponseObject, error) {
	return GetRecipeById404Response{}, nil
}

func (p *RecipeStore) UpdateRecipe(context.Context, UpdateRecipeRequestObject) (UpdateRecipeResponseObject, error) {
	return UpdateRecipe404Response{}, nil
}

func (p *RecipeStore) DeleteRecipe(_ context.Context, request DeleteRecipeRequestObject) (DeleteRecipeResponseObject, error) {
	return DeleteRecipe404Response{}, nil
}
