package controllers

import (
	"github.com/KamilSirazhdin/FocusUp/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

// GetStats возвращает статистику топ-10 пользователей по очкам
func GetStats(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var topUsers []models.User

		// Получение топ-10 пользователей по очкам
		if err := db.Order("points desc").Limit(10).Find(&topUsers).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при получении статистики"))
		}

		// Формирование ответа
		var response []models.UserStatsResponse
		for _, user := range topUsers {
			response = append(response, models.UserStatsResponse{
				Username: user.Username,
				Points:   user.Points,
				Streak:   user.Streak,
			})
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Статистика получена успешно", response))
	}
}

// StreakStats возвращает статистику стрика текущего пользователя
func StreakStats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение пользователя из контекста
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Пользователь не аутентифицирован"))
		}

		// Формирование ответа
		response := models.StatsResponse{
			Points: user.Points,
			Streak: user.Streak,
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Статистика стрика получена успешно", response))
	}
}

// GetUserTaskHistory возвращает историю ответов пользователя на задания
func GetUserTaskHistory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение пользователя из контекста
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Пользователь не аутентифицирован"))
		}

		// Получение истории ответов
		var logs []models.UserTaskLog
		if err := db.Where("user_id = ?", user.ID).
			Order("answered_at desc").
			Limit(50).
			Find(&logs).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при получении истории ответов"))
		}

		// Получение информации о заданиях
		var taskIDs []uint
		for _, log := range logs {
			taskIDs = append(taskIDs, log.TaskID)
		}

		var tasks []models.Task
		if err := db.Where("id IN ?", taskIDs).Find(&tasks).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при получении информации о заданиях"))
		}

		// Создание карты заданий для быстрого доступа
		taskMap := make(map[uint]models.Task)
		for _, task := range tasks {
			taskMap[task.ID] = task
		}

		// Формирование ответа
		type HistoryItem struct {
			ID         uint      `json:"id"`
			TaskID     uint      `json:"task_id"`
			Question   string    `json:"question"`
			Answer     string    `json:"answer"`
			UserAnswer string    `json:"user_answer"`
			Correct    bool      `json:"correct"`
			Points     int       `json:"points"`
			AnsweredAt time.Time `json:"answered_at"`
		}

		var history []HistoryItem
		for _, log := range logs {
			task, exists := taskMap[log.TaskID]
			if !exists {
				continue
			}

			history = append(history, HistoryItem{
				ID:         log.ID,
				TaskID:     log.TaskID,
				Question:   task.Question,
				Answer:     task.Answer,
				Correct:    log.Correct,
				Points:     task.Points,
				AnsweredAt: log.AnsweredAt,
			})
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("История ответов получена успешно", history))
	}
}
