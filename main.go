package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"recipeApi/api"
	"recipeApi/storage"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newAuthenticator() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		//Log request and context to ensure auth is called for expected routes
		fmt.Println("Authenticating request", input.RequestValidationInput.Request.URL.Path, ctx)

		// Extract the auth header
		authHdr := input.RequestValidationInput.Request.Header.Get("Authorization")
		fmt.Println("Auth header : ", authHdr)

		//Validate

		return nil
	}
}

// func setupDatabase() *pgxpool.Pool {
// 	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer dbpool.Close()

// }

func main() {
	// pool := setupDatabase()

	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.Parse()

	recipeStore := NewService(storage.NewInMemoryUserStore(), storage.NewInMemoryRecipeStore())
	recipeStoreStrictHandler := api.NewStrictHandler(recipeStore, nil)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// Use middleware to check all requests against the
	// OpenAPI schema and authenticate requests
	validator := CreateAuthMiddleware(newAuthenticator())
	r.Use(validator)

	api.HandlerFromMux(recipeStoreStrictHandler, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}

	log.Fatal(s.ListenAndServe())
}
