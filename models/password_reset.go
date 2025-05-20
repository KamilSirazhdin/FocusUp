// models/password_reset.go
package models

import "time"

type PasswordReset struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"index"`
	Code      string
	Token     string
	Used      bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
