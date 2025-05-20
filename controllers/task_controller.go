package controllers

import (
	"errors"
	"time"

	"github.com/KamilSirazhdin/FocusUp/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetRandomTask возвращает случайное задание
func GetRandomTask(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var task models.Task

		// Получение случайного задания
		if err := db.Order("RANDOM()").
			Limit(1).
			Select("id", "question", "points").
			First(&task).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse("Задания не найдены"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при получении задания"))
		}

		// Формирование ответа
		taskResponse := models.TaskResponse{
			ID:       task.ID,
			Question: task.Question,
			Points:   task.Points,
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Задание получено успешно", taskResponse))
	}
}

// AnswerInput структура для ввода ответа на задание
type AnswerInput struct {
	Answer string `json:"answer" validate:"required"`
}

// AnswerTask обрабатывает ответ пользователя на задание
func AnswerTask(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение пользователя из контекста
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Пользователь не аутентифицирован"))
		}

		// Получение ID задания из URL
		taskID := c.Params("id")

		// Поиск задания
		var task models.Task
		if err := db.First(&task, taskID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse("Задание не найдено"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при поиске задания"))
		}

		// Получение ответа пользователя
		var input AnswerInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		// Проверка ответа
		isCorrect := input.Answer == task.Answer

		// Создание лога ответа
		log := models.UserTaskLog{
			UserID:     user.ID,
			TaskID:     task.ID,
			AnsweredAt: time.Now(),
			Correct:    isCorrect,
		}

		// Сохранение лога в БД
		if err := db.Create(&log).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при сохранении ответа"))
		}

		// Если ответ правильный, обновляем очки и стрик
		if isCorrect {
			// Добавление очков
			user.AddPoints(task.Points)

			// Обновление стрика
			user.UpdateStreak()

			// Сохранение изменений пользователя
			if err := db.Save(&user).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при обновлении данных пользователя"))
			}
		}

		// Формирование ответа
		response := models.TaskAnswerResponse{
			Correct: isCorrect,
			Points:  user.Points,
			Streak:  user.Streak,
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Ответ обработан", response))
	}
}
