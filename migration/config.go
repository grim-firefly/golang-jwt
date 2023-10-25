package migration

import (
	"github.com/grim-firefly/golang-jwt/database"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.GetDB()

func Migrate() {
	MigrateUser()
}
