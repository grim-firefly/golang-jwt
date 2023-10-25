package migration

import "github.com/grim-firefly/golang-jwt/models"

func MigrateUser() {
	DB.AutoMigrate(&models.User{})
}
