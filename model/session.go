package model

type Session struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}
