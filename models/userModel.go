package models

import (
	"time"
)

type User struct {
	ID            uint64    `json:"id"`
	First_name    string    `gorm:"not null" json:"first_name" validate:"required, min=2, max=100"`
	Last_name     string    `gorm:"not null" json:"last_name" validate:"required, min=2, max=100"`
	Password      string    `gorm:"not null" json:"password" validate:"required, min=6"`
	Email         string    `gorm:"not null" json:"email" validate:"email,required"`
	Phone         string    `gorm:"not null" json:"phone" validate:"required"`
	Token         string    `json:"token"`
	User_type     string    `gorm:"not null" json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Refresh_token string    `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	User_id       string    `gorm:"unique;not null;type:varchar(100)"json:"user_id"`
}
