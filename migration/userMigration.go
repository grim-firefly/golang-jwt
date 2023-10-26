package migration

import "github.com/grim-firefly/golang-jwt/models"

func MigrateUser() {
	DB.AutoMigrate(&models.User{})
	if DB.Migrator().HasColumn(&models.User{}, "user_id") {
		DB.Migrator().DropColumn(&models.User{}, "user_id")
	}
}
