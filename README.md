# Recipe API

This is a go service for managing recipes. This project uses go-chi for routing and [oapi-codegen](https://github.com/deepmap/oapi-codegen) to generate boilerplate from the OpenAPI spec defined [here](./recipe-api-v1.json).

## Endpoints

- `/recipes`
  - `POST`: Create a new recipe. Expects a `BaseRecipe` schema in the request body. Returns a `Recipe` schema.
- `/recipes/{id}`
  - `GET`: Retrieve a recipe by ID. Returns a `Recipe` schema.
  - `PUT`: Update a recipe by ID. Expects a `Recipe` schema in the request body.
  - `DELETE`: Delete a recipe by ID.

## User Authentication

- `/users`
  - `POST`: Register a new user. Expects a `NewUser` schema in the request body.
- `/login`
  - `POST`: Login a user. Expects a `LoginCredentials` schema in the request body. Returns a `LoginResponse` schema.

## Schemas

- `BaseRecipe`: Used for creating a new recipe. 
  - `title`: A string that represents the title of the recipe.
  - `description`: A string that provides a brief description of the recipe.
  - `ingredients`: An array of strings, each representing an ingredient required for the recipe.
  - `steps`: An array of strings, each representing a step in the recipe.

- `Recipe`: Extends `BaseRecipe` with an `id` field. Used for returning existing recipes.
  - `id`: An integer that represents the unique identifier of the recipe.
  - `creatorId`: An integer that represents the unique identifier of user that created the recipe.

- `NewUser`: Used for registering a new user. 
  - `email`: A string that represents the user's email address.
  - `password`: A string that represents the user's password.

- `LoginCredentials`: Used for logging in a user. 
  - `email`: A string that represents the user's email address.
  - `password`: A string that represents the user's password.

- `LoginResponse`: Returned when a user logs in successfully.
  - `token`: A string that represents the authentication token.
