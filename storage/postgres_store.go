package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	// _ "github.com/lib/pq"
)

type PostgresDatatore struct {
	db *pgxpool.Pool
}

func NewPostgresDatastore(db *pgxpool.Pool) (*PostgresDatatore, error) {
	return &PostgresDatatore{db: db}, nil
}

func (s *PostgresDatatore) CreateUser(user User) error {
	_, err := s.db.Exec(context.TODO(),
		"INSERT INTO users (email, password) VALUES ($1, $2)",
		user.Email, user.Password,
	)
	return err
}

func (s *PostgresDatatore) ReadUser(email string) (User, error) {
	var user User
	err := s.db.QueryRow(context.TODO(),
		"SELECT email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return User{}, errors.New("user not found")
	} else if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *PostgresDatatore) CreateRecipe(recipe Recipe) (Recipe, error) {
	err := s.db.QueryRow(context.TODO(),
		"INSERT INTO recipes (title, ingredients, steps) VALUES ($1, $2, $3) RETURNING id",
		recipe.Title, recipe.Ingredients, recipe.Steps,
	).Scan(&recipe.ID)
	if err != nil {
		//print error
		fmt.Println(err)
		return Recipe{}, err
	}
	return recipe, nil
}

func (s *PostgresDatatore) ReadRecipe(id int64) (Recipe, error) {
	var recipe Recipe
	err := s.db.QueryRow(context.TODO(),
		"SELECT id, title, ingredients, steps FROM recipes WHERE id = $1",
		id,
	).Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Steps)
	if err == sql.ErrNoRows {
		return Recipe{}, errors.New("recipe not found")
	} else if err != nil {
		return Recipe{}, err
	}
	return recipe, nil
}

func (s *PostgresDatatore) ListRecipes() ([]Recipe, error) {
	rows, err := s.db.Query(context.TODO(), "SELECT id, title, ingredients, steps FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Steps)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (s *PostgresDatatore) UpdateRecipe(id int64, recipe Recipe) error {
	result, err := s.db.Exec(context.TODO(),
		"UPDATE recipes SET title = $1, ingredients = $2, steps = $3 WHERE id = $4",
		recipe.Title, recipe.Ingredients, recipe.Steps, id,
	)
	rowsAffected := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("recipe not found")
	}
	return nil
}

func (s *PostgresDatatore) DeleteRecipe(id int64) error {
	result, err := s.db.Exec(context.TODO(), "DELETE FROM recipes WHERE id = $1", id)

	rowsAffected := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("recipe not found")
	}
	return nil
}
