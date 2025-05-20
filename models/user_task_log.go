package models

import (
	"gorm.io/gorm"
	"time"
)

// UserTaskLog представляет запись о попытке пользователя решить задачу
type UserTaskLog struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	TaskID     uint      `json:"task_id"`
	AnsweredAt time.Time `json:"answered_at"`
	Correct    bool      `json:"correct"`
}
