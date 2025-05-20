// controllers/admin_task_controller.go
package controllers

import (
	"github.com/KamilSirazhdin/FocusUp/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllTasks возвращает список всех заданий для администратора
func GetAllTasks(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tasks []models.Task

		// Получаем все задания, включая удаленные с мягким удалением
		if err := db.Find(&tasks).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при получении списка заданий"))
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Список заданий получен успешно", tasks))
	}
}
