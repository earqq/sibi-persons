package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/earqq/gqlgen-easybill/auth"
	"github.com/earqq/gqlgen-easybill/db"
	"github.com/earqq/gqlgen-easybill/graph"
	"github.com/earqq/gqlgen-easybill/graph/generated"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "8087"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := chi.NewRouter()
	db.ConnectDB()
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"htps://chambea.pe"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             3599, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)

	router.Use(auth.Middleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playgro	und", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
