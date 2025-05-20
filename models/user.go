package models

import (
	"gorm.io/gorm"
	"time"
)

// User представляет пользователя в системе
type User struct {
	gorm.Model
	Username       string `gorm:"unique;not null"`
	Email          string `gorm:"unique;not null"`
	PasswordHash   string `gorm:"not null"`
	Role           string `gorm:"default:user"` // user | admin
	Points         int    `gorm:"default:0"`
	Streak         int    `gorm:"default:0"`
	LastActiveDate time.Time
}

// SafeUser представляет структуру для безопасного отображения пользователя в API
type SafeUser struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Points    int       `json:"points"`
	Streak    int       `json:"streak"`
	CreatedAt time.Time `json:"created_at"`
}

// ToSafeUser конвертирует User в SafeUser, скрывая конфиденциальную информацию
func (u *User) ToSafeUser() SafeUser {
	return SafeUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Points:    u.Points,
		Streak:    u.Streak,
		CreatedAt: u.CreatedAt,
	}
}

// UpdateStreak обновляет стрик пользователя на основе времени последней активности
func (u *User) UpdateStreak() bool {
	today := time.Now().Truncate(24 * time.Hour)
	lastActive := u.LastActiveDate.Truncate(24 * time.Hour)

	// Если уже был активен сегодня, стрик не меняется
	if today.Equal(lastActive) {
		return false
	}

	// Если был активен вчера, увеличиваем стрик
	yesterday := today.Add(-24 * time.Hour)
	if yesterday.Equal(lastActive) {
		u.Streak++
		return true
	}

	// Если пропустил день или больше, сбрасываем стрик до 1
	u.Streak = 1
	return true
}

// AddPoints добавляет очки пользователю и обновляет дату последней активности
func (u *User) AddPoints(points int) {
	u.Points += points
	u.LastActiveDate = time.Now()
}
