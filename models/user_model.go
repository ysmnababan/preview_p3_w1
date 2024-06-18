package models

type User struct {
	UserID   uint    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username string  `json:"username" gorm:"size:255"`
	Email    string  `json:"email" gorm:"unique;size:255;not null"`
	Password string  `json:"password" gorm:"size:255"`
	Limit    float64 `json:"limit" gorm:"type:decimal(10,2)"`
	Balance  float64 `json:"balance" gorm:"type:decimal(10,2)"`
}
