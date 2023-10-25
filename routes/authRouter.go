package routes

import (
	"github.com/go-chi/chi/v5"
	controller "github.com/grim-firefly/golang-jwt/controllers"
)

func AuthRoutes(router *chi.Mux) {
	router.Post("/signup", controller.SignUp)
	router.Post("/login", controller.Login)
}
