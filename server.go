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
)

const defaultPort = "8085"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db.ConnectDB()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router := chi.NewRouter()
	router.Use(auth.Middleware())
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
