package controllers

import (
	"net/http"

	"github.com/grim-firefly/golang-jwt/database"
	"github.com/grim-firefly/golang-jwt/helpers"
	"github.com/grim-firefly/golang-jwt/models"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.GetDB()

func SignUp(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {

}
func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []models.User
	DB.Find(&users)
	helpers.ResponseJson(w, 200, users)

}
func GetUser(w http.ResponseWriter, r *http.Request) {

}
