package models

import "gorm.io/gorm"

// Task представляет задачу в системе
type Task struct {
	gorm.Model
	Question    string `gorm:"not null" json:"question"`
	Answer      string `gorm:"not null" json:"answer"`
	Points      int    `gorm:"default:10" json:"points"`
	CreatedByID uint   `json:"created_by_id"`
}
