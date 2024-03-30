package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"recipeApi/api"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newAuthenticator() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		//Log request and context to ensure auth is called for expected routes
		fmt.Println("Authenticating request", input.RequestValidationInput.Request.URL.Path, ctx)
		fmt.Println("Auth scopes : ", ctx.Value(api.BearerAuthScopes))

		// Extract the auth header
		authHdr := input.RequestValidationInput.Request.Header.Get("Authorization")
		fmt.Println("Auth header : ", authHdr)

		return nil
	}
}

func main() {
	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.Parse()

	// Create an instance of our handler which satisfies the generated interface
	recipeStore := api.NewRecipeStore()
	recipeStoreStrictHandler := api.NewStrictHandler(recipeStore, nil)

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	validator := api.CreateMiddleware(newAuthenticator())
	r.Use(validator)

	// We now register our petStore above as the handler for the interface
	api.HandlerFromMux(recipeStoreStrictHandler, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
