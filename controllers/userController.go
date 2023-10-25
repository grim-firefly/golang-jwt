package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/grim-firefly/golang-jwt/database"
	"github.com/grim-firefly/golang-jwt/helpers"
	"github.com/grim-firefly/golang-jwt/models"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.GetDB()

func SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		panic(err)
		return
	}
	newUser.User_type = "User"
	Err := DB.Create(&newUser)
	if Err != nil {
		helpers.ResponseJson(w, 502, struct {
			message string
		}{
			message: "Error Creating User",
		})

	}
	helpers.ResponseJson(w, 200, newUser)

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
