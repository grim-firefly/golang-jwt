package controllers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/grim-firefly/golang-jwt/database"
	"github.com/grim-firefly/golang-jwt/helpers"
	"github.com/grim-firefly/golang-jwt/models"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.GetDB()

// sign up with validation
func SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		data := make(map[string]string)
		data["message"] = "Error"
		helpers.ResponseJson(w, http.StatusOK, map[string]string{
			"message": "Error",
		})
		return
	}
	errors := make(map[string]string)

	// validation
	if strings.TrimSpace(newUser.First_name) == "" {
		errors["first_name"] = "First name can't be empty"
	}
	if strings.TrimSpace(newUser.Last_name) == "" {
		errors["last_name"] = "Last name can't be empty"
	}
	if strings.TrimSpace(newUser.Password) == "" {
		errors["password"] = "Password can't be empty"
	}
	if utf8.RuneCountInString(newUser.Password) < 6 {
		errors["password"] = "Password Length Must be atleast 6"
	}
	if utf8.RuneCountInString(newUser.Password) > 20 {
		errors["password"] = "Password Length Must be less than 20"
	}
	if _, err := mail.ParseAddress(newUser.Email); err != nil {
		errors["email"] = "Invalid Email"
	}
	if utf8.RuneCountInString(newUser.Phone) != 11 {
		errors["phone"] = "Phone Number must have 11 digit "
	}

	if DB.Where("email = ?", newUser.Email).First(&models.User{}).Error == nil {
		errors["email"] = "Email is Already Exist"
	}

	if len(errors) > 0 {
		helpers.ResponseJson(w, http.StatusOK, errors)
		return
	}
	// validation ends here

	newUser.User_type = "User"

	if DB.Create(&newUser).Error != nil {
		errors["message"] = "Failed to create New User"
		helpers.ResponseJson(w, http.StatusBadGateway, errors)
		return
	}

	helpers.ResponseJson(w, http.StatusOK, newUser)

}

// login
func Login(w http.ResponseWriter, r *http.Request) {

}

// get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []models.User
	DB.Find(&users)
	helpers.ResponseJson(w, 200, users)

}

// get specific user using id
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var user models.User

	if DB.First(&user, id).Error != nil {

		helpers.ResponseJson(w, http.StatusOK, map[string]string{
			"message": "User not Found",
		})
		return
	}
	helpers.ResponseJson(w, http.StatusOK, user)
}
