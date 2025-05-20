package middleware

import (
	"fmt"
	"strings"

	"github.com/KamilSirazhdin/FocusUp/models"
	"github.com/KamilSirazhdin/FocusUp/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// JWTMiddleware проверяет JWT-токен и устанавливает пользователя в контекст
func JWTMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение токена из заголовка Authorization
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Отсутствует токен авторизации"))
		}

		// Извлечение токена
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		// Разбор и проверка токена
		userID, role, err := utils.ParseJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Недействительный токен"))
		}

		// Поиск пользователя в БД
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Пользователь не найден"))
		}

		// Установка пользователя и роли в контекст
		c.Locals("user", user)
		c.Locals("role", role)
		return c.Next()
	}
}

// RoleMiddleware проверяет роль пользователя
func RoleMiddleware(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение роли из контекста
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Роль пользователя не определена"))
		}

		// Проверка роли
		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(models.NewErrorResponse(
				fmt.Sprintf("Недостаточно прав. Требуется роль: %s", requiredRole)))
		}

		return c.Next()
	}
}

// ErrorLogger логирует ошибки запросов
func ErrorLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Вызов следующего обработчика
		err := c.Next()

		// Если возникла ошибка, логируем её
		if err != nil {
			fmt.Printf("Route Error: %s %s - %v\n", c.Method(), c.Path(), err)
		}

		return err
	}
}
