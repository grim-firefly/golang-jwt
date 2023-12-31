package routes

import (
	"github.com/go-chi/chi/v5"
	controller "github.com/grim-firefly/golang-jwt/controllers"
)

func UserRoutes(router *chi.Mux) {
	router.Get("/users", controller.GetUsers)
	router.Get("/users/{id}", controller.GetUser)
	router.Get("/validation", controller.Validation)
	router.Get("/refresh_token", controller.RefreshToken)

}
