package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"recipeApi/api"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.Parse()

	spec, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	spec.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	recipeStore := api.NewRecipeStore()

	recipeStoreStrictHandler := api.NewStrictHandler(recipeStore, nil)

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec,
		&nethttpmiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: recipeStore.NewAuthenticator(),
			},
		})
	r.Use(validator)

	// r.Use(jwt.AuthMiddleware)

	// We now register our petStore above as the handler for the interface
	api.HandlerFromMux(recipeStoreStrictHandler, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
