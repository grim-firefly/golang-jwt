package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login
func Login(w http.ResponseWriter, r *http.Request) {

	var credential Credentials
	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		data := make(map[string]string)
		data["message"] = "Error"
		helpers.ResponseJson(w, http.StatusOK, map[string]string{
			"message": "Error",
		})
		return
	}
	var exUser models.User
	if DB.Where("email = ? ", credential.Email).Select("id", "email", "password").First(&exUser).Error != nil {
		helpers.ResponseJson(w, http.StatusOK, map[string]string{
			"message": "Invalid Email or Password",
		})
		return
	}
	if exUser.Password != credential.Password {
		helpers.ResponseJson(w, http.StatusOK, map[string]string{
			"message": "Password is not matched",
		})
		return
	}

	//jwt  token creating
	claims := jwt.MapClaims{
		"id":     exUser.ID,
		"email":  exUser.Email,
		"expire": time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	exUser.Token = tokenString
	// jwt token ends

	helpers.ResponseJson(w, http.StatusOK, exUser)

	// helpers.ResponseJson(w, http.StatusOK, Credential) ]

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

// for validating token
func Validation(w http.ResponseWriter, r *http.Request) {
	if r.Header["Token"] != nil {
		tokenString := r.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				helpers.ResponseJson(w, http.StatusUnauthorized, map[string]string{
					"message": "Unauthorized",
				})
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("secret"), nil
		})
		if err != nil {
			helpers.ResponseJson(w, http.StatusUnauthorized, map[string]string{
				"message": "Unauthorized",
			})
		}
		if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			helpers.ResponseJson(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Authorization Successfull",
				"expire":  claim["expire"],
				"email":   claim["email"],
			})
		} else {
			helpers.ResponseJson(w, http.StatusUnauthorized, map[string]string{
				"message": "Unauthorized",
			})
		}

		// jwt.ParseWithClaims(token)
	}
}

// for refreshing token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header["Token"][0]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if err != nil {
		helpers.ResponseJson(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims["expire"] = time.Now().Add(time.Minute * 5).Unix()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("secret"))

		helpers.ResponseJson(w, http.StatusUnauthorized, map[string]string{
			"message": "Authorization Successfull",
			"Token":   tokenString,
		})
	} else {
		helpers.ResponseJson(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
	}
}
