package controllers

import (
	"errors"

	"github.com/KamilSirazhdin/FocusUp/models"
	"github.com/KamilSirazhdin/FocusUp/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// TaskInput структура для создания/update задания
type TaskInput struct {
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
	Points   int    `json:"points"`
}

// CreateTask создает новое задание
func CreateTask(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение пользователя из контекста
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Пользователь не аутентифицирован"))
		}

		var input TaskInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		if err := utils.ValidateStruct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
		}

		// Проверка на существование похожего задания
		var existingTask models.Task
		result := db.Where("question = ? OR answer = ?", input.Question, input.Answer).First(&existingTask)
		if result.Error == nil {
			// Задание уже существует
			return c.Status(fiber.StatusConflict).JSON(models.NewErrorResponse("Задание с похожим вопросом или ответом уже существует"))
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Другая ошибка БД
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при проверке существования задания"))
		}

		// Установка значения по умолчанию для очков
		if input.Points <= 0 {
			input.Points = 10
		}

		// Создание задания
		task := models.Task{
			Question:    input.Question,
			Answer:      input.Answer,
			Points:      input.Points,
			CreatedByID: user.ID,
		}

		if err := db.Create(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при создании задания"))
		}

		return c.Status(fiber.StatusCreated).JSON(models.NewSuccessResponse("Задание создано успешно", task))
	}
}

// UpdateTask обновляет существующее задание
func UpdateTask(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Проверка существования задания
		var task models.Task
		if err := db.First(&task, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse("Задание не найдено"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при поиске задания"))
		}

		var input TaskInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		if err := utils.ValidateStruct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
		}

		// Обновление задания
		task.Question = input.Question
		task.Answer = input.Answer

		// Проверка на положительные очки
		if input.Points > 0 {
			task.Points = input.Points
		}

		if err := db.Save(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при обновлении задания"))
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Задание обновлено успешно", task))
	}
}

// DeleteTask удаляет задание
func DeleteTask(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Проверка существования задания
		var task models.Task
		if err := db.First(&task, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse("Задание не найдено"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при поиске задания"))
		}

		// Удаление задания (soft delete благодаря gorm.Model)
		if err := db.Delete(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при удалении задания"))
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Задание удалено успешно", nil))
	}
}

// GetAllUsers возвращает список всех пользователей
func GetAllUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при получении списка пользователей"))
		}

		// Преобразование пользователей в безопасный формат
		safeUsers := make([]models.SafeUser, len(users))
		for i, user := range users {
			safeUsers[i] = user.ToSafeUser()
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Список пользователей получен успешно", safeUsers))
	}
}

// DeleteUser удаляет пользователя
func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Проверка существования пользователя
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse("Пользователь не найден"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при поиске пользователя"))
		}

		// Удаление пользователя (soft delete благодаря gorm.Model)
		if err := db.Delete(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при удалении пользователя"))
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Пользователь удален успешно", nil))
	}
}
