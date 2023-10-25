package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grim-firefly/golang-jwt/migration"
	"github.com/grim-firefly/golang-jwt/routes"
)

func main() {
	router := chi.NewRouter()
	router.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// migration

	migration.Migrate()
	// migration ends
	http.ListenAndServe(":8080", router)
}
