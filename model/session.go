package model

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID uint `json:"user_id"`
}
