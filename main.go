package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"recipeApi/api"
	"recipeApi/storage"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

func newAuthenticator() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		//Log request and context to ensure auth is called for expected routes
		// fmt.Println("Authenticating request", input.RequestValidationInput.Request.URL.Path, ctx)

		// Extract the auth header
		authHdr := input.RequestValidationInput.Request.Header.Get("Authorization")
		// fmt.Println("Auth header : ", authHdr)

		//Parse the auth header and verify the token
		//TODO add token cache, check for real, and set user in context
		authHdr = authHdr[len("Bearer "):]
		if authHdr == "test" {
			return nil
		}

		return errors.New("invalid token")
	}
}

func setupDatabase(db_url string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), db_url)
	println(os.Getenv("DATABASE_URL"))
	println(db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		db_url = "postgresql://postgres:postgres@localhost/postgres"
	}

	pool := setupDatabase(db_url)
	defer pool.Close()

	postgresStore, _ := storage.NewPostgresDatastore(pool)
	recipeStore := NewService(postgresStore, postgresStore)

	// recipeStore := NewService(storage.NewInMemoryUserStore(), storage.NewInMemoryRecipeStore())
	recipeStoreStrictHandler := api.NewStrictHandler(recipeStore, nil)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Default().Handler)
	// Use middleware to check all requests against the
	// OpenAPI schema and authenticate requests
	validator := CreateAuthMiddleware(newAuthenticator())
	r.Use(validator)

	api.HandlerFromMux(recipeStoreStrictHandler, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}

	log.Fatal(s.ListenAndServe())
}
