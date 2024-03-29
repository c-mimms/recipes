package api

import (
	"context"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64
	Email    string
	Password string
}

type RecipeStore struct {
	//Temporary in memory stores
	Users        map[string]User
	Recipes      map[int64]Recipe
	NextUserId   int64
	NextRecipeId int64
	Lock         sync.Mutex
}

var _ StrictServerInterface = (*RecipeStore)(nil)

func NewRecipeStore() *RecipeStore {
	return &RecipeStore{
		Users:        make(map[string]User),
		Recipes:      make(map[int64]Recipe),
		NextUserId:   1000,
		NextRecipeId: 1,
	}
}

func (p *RecipeStore) CreateUser(_ context.Context, request CreateUserRequestObject) (CreateUserResponseObject, error) {
	//Secure hash password before storing
	hashed, _ := bcrypt.GenerateFromPassword([]byte(request.Body.Password), 8)

	var user User
	user.Email = string(request.Body.Email)
	user.Password = string(hashed)

	p.Lock.Lock()
	defer p.Lock.Unlock()

	user.Id = p.NextUserId
	p.NextUserId++

	// Insert into map
	p.Users[string(user.Email)] = user

	return CreateUser201Response{}, nil
}
func (p *RecipeStore) LoginUser(_ context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error) {
	//Check if user exists and password matches after hashing
	user, found := p.Users[string(request.Body.Email)]
	if !found {
		return LoginUser401Response{}, nil
	}

	//Generate new auth token
	//TODO implement this
	token := "authToken"

	if nil == bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Body.Password)) {
		return LoginUser200JSONResponse(LoginResponse{Token: &token}), nil
	}

	return LoginUser401Response{}, nil
}

func (p *RecipeStore) CreateRecipe(_ context.Context, request CreateRecipeRequestObject) (CreateRecipeResponseObject, error) {
	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	var recipe Recipe
	recipe.Title = request.Body.Title
	recipe.Description = request.Body.Description
	recipe.Ingredients = request.Body.Ingredients
	recipe.Steps = request.Body.Steps
	recipe.Id = p.NextRecipeId
	p.NextRecipeId++

	// Insert into map
	p.Recipes[recipe.Id] = recipe

	return CreateRecipe201JSONResponse(recipe), nil
}

func (p *RecipeStore) ListRecipes(context.Context, ListRecipesRequestObject) (ListRecipesResponseObject, error) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	var result []Recipe

	for _, recipe := range p.Recipes {
		result = append(result, recipe)
	}
	return ListRecipes200JSONResponse(result), nil
}

func (p *RecipeStore) GetRecipeById(_ context.Context, request GetRecipeByIdRequestObject) (GetRecipeByIdResponseObject, error) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	recipe, found := p.Recipes[request.Id]
	if !found {
		return GetRecipeById404Response{}, nil
	}

	return GetRecipeById200JSONResponse(recipe), nil
}

func (p *RecipeStore) UpdateRecipe(context.Context, UpdateRecipeRequestObject) (UpdateRecipeResponseObject, error) {
	//TODO implement
	return UpdateRecipe404Response{}, nil
}

func (p *RecipeStore) DeleteRecipe(_ context.Context, request DeleteRecipeRequestObject) (DeleteRecipeResponseObject, error) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, found := p.Recipes[request.Id]
	if !found {
		return DeleteRecipe404Response{}, nil
	}
	delete(p.Recipes, request.Id)

	return DeleteRecipe200Response{}, nil
}
