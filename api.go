package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"recipeApi/api"
	"recipeApi/storage"

	"github.com/getkin/kin-openapi/openapi3filter"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"golang.org/x/crypto/bcrypt"
)

func convertRecipe(recipe storage.Recipe) api.Recipe {
	//Convert storage recipe to API recipe
	apiRecipe := api.Recipe{
		Id:          recipe.ID,
		CreatorId:   recipe.CreatorID,
		Title:       recipe.Title,
		Description: recipe.Description,
		Ingredients: recipe.Ingredients,
		Steps:       recipe.Steps,
	}

	return apiRecipe
}

type Service struct {
	UserDB   storage.UserDatastore
	RecipeDB storage.RecipeDatastore
}

var _ api.StrictServerInterface = (*Service)(nil)

func NewService(userStore storage.UserDatastore, recipeStore storage.RecipeDatastore) *Service {
	return &Service{
		UserDB:   userStore,
		RecipeDB: recipeStore,
	}
}

// Creaates middleware to validate requests against the OpenAPI schema and authenticate requests
func CreateAuthMiddleware(authenticator openapi3filter.AuthenticationFunc) func(next http.Handler) http.Handler {
	spec, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	spec.Servers = nil

	return nethttpmiddleware.OapiRequestValidatorWithOptions(spec,
		&nethttpmiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: authenticator,
			},
		})
}

func (p *Service) CreateUser(_ context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	//Secure hash password before storing
	hashed, _ := bcrypt.GenerateFromPassword([]byte(request.Body.Password), 8)

	var user storage.User
	user.Email = string(request.Body.Email)
	user.Password = string(hashed)

	p.UserDB.CreateUser(user)

	return api.CreateUser201Response{}, nil
}
func (p *Service) LoginUser(_ context.Context, request api.LoginUserRequestObject) (api.LoginUserResponseObject, error) {
	//Check if user exists and password matches after hashing
	user, err := p.UserDB.ReadUser(string(request.Body.Email))
	if err != nil {
		return api.LoginUser401Response{}, nil
	}

	if nil == bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Body.Password)) {
		//Generate new auth token
		//TODO implement JWT token generation
		token := string(user.Email)

		return api.LoginUser200JSONResponse(api.LoginResponse{Token: &token}), nil
	}

	return api.LoginUser401Response{}, nil
}

func (p *Service) CreateRecipe(ctx context.Context, request api.CreateRecipeRequestObject) (api.CreateRecipeResponseObject, error) {
	var recipe storage.Recipe
	recipe.Title = request.Body.Title
	recipe.Description = request.Body.Description
	recipe.Ingredients = request.Body.Ingredients
	recipe.Steps = request.Body.Steps

	createdRecipe, _ := p.RecipeDB.CreateRecipe(recipe)

	return api.CreateRecipe201JSONResponse(convertRecipe(createdRecipe)), nil
}

func (p *Service) ListRecipes(context.Context, api.ListRecipesRequestObject) (api.ListRecipesResponseObject, error) {

	var result []api.Recipe
	recipes, _ := p.RecipeDB.ListRecipes()

	for _, recipe := range recipes {
		result = append(result, convertRecipe(recipe))
	}
	return api.ListRecipes200JSONResponse(result), nil
}

func (p *Service) GetRecipeById(_ context.Context, request api.GetRecipeByIdRequestObject) (api.GetRecipeByIdResponseObject, error) {

	recipe, _ := p.RecipeDB.ReadRecipe(request.Id)

	return api.GetRecipeById200JSONResponse(convertRecipe(recipe)), nil
}

func (p *Service) UpdateRecipe(context.Context, api.UpdateRecipeRequestObject) (api.UpdateRecipeResponseObject, error) {
	//TODO implement
	return api.UpdateRecipe404Response{}, nil
}

func (p *Service) DeleteRecipe(_ context.Context, request api.DeleteRecipeRequestObject) (api.DeleteRecipeResponseObject, error) {

	p.RecipeDB.DeleteRecipe(request.Id)

	return api.DeleteRecipe200Response{}, nil
}
