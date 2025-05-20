package controllers

import (
	"errors"
	"strings"
	"time"

	"github.com/KamilSirazhdin/FocusUp/models"
	"github.com/KamilSirazhdin/FocusUp/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterInput структура для регистрации пользователя
type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role"`
}

// Register регистрирует нового пользователя
func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input RegisterInput

		// Парсинг и валидация входных данных
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		if err := utils.ValidateStruct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
		}

		// Нормализация email
		input.Email = strings.ToLower(strings.TrimSpace(input.Email))

		// Проверка существования пользователя
		var existingUser models.User
		if result := db.Where("email = ?", input.Email).First(&existingUser); !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusConflict).JSON(models.NewErrorResponse("Пользователь с таким email уже существует"))
		}

		// Хеширование пароля
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при хешировании пароля"))
		}

		// Создание пользователя
		user := models.User{
			Username:       input.Username,
			Email:          input.Email,
			PasswordHash:   string(hash),
			Role:           input.Role,
			LastActiveDate: time.Now(),
		}

		// Установка роли по умолчанию
		if user.Role == "" {
			user.Role = "user"
		}

		// Сохранение пользователя в БД
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при создании пользователя"))
		}

		return c.Status(fiber.StatusCreated).JSON(models.NewSuccessResponse("Регистрация прошла успешно", user.ToSafeUser()))
	}
}

// LoginInput структура для входа пользователя
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login аутентифицирует пользователя и возвращает JWT-токен
func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input LoginInput

		// Парсинг и валидация входных данных
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		if err := utils.ValidateStruct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
		}

		// Нормализация email
		input.Email = strings.ToLower(strings.TrimSpace(input.Email))

		// Поиск пользователя
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Неверный email или пароль"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при поиске пользователя"))
		}

		// Проверка пароля
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Неверный email или пароль"))
		}

		// Генерация JWT токена
		token, err := utils.GenerateJWT(user.ID, user.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при создании токена"))
		}

		// Создание ответа
		authResponse := models.AuthResponse{
			Token: token,
			User:  user.ToSafeUser(),
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Вход выполнен успешно", authResponse))
	}
}

// Profile возвращает информацию о текущем пользователе
func Profile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение пользователя из контекста (установлен в middleware)
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Пользователь не аутентифицирован"))
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Профиль пользователя", user.ToSafeUser()))
	}
}

// Logout выполняет выход пользователя
func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Очистка cookie с токеном
		c.ClearCookie("token")
		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Выход выполнен успешно", nil))
	}
}

// UpdateProfileInput структура для обновления профиля
type UpdateProfileInput struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

// UpdateProfile обновляет профиль пользователя
func UpdateProfile(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получение пользователя из контекста
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse("Пользователь не аутентифицирован"))
		}

		var input UpdateProfileInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		if err := utils.ValidateStruct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
		}

		// Обновление данных пользователя
		updates := false

		if input.Username != "" && input.Username != user.Username {
			// Проверка уникальности username
			var count int64
			db.Model(&models.User{}).Where("username = ? AND id != ?", input.Username, user.ID).Count(&count)
			if count > 0 {
				return c.Status(fiber.StatusConflict).JSON(models.NewErrorResponse("Пользователь с таким именем уже существует"))
			}

			user.Username = input.Username
			updates = true
		}

		if input.Email != "" && input.Email != user.Email {
			// Нормализация email
			email := strings.ToLower(strings.TrimSpace(input.Email))

			// Проверка уникальности email
			var count int64
			db.Model(&models.User{}).Where("email = ? AND id != ?", email, user.ID).Count(&count)
			if count > 0 {
				return c.Status(fiber.StatusConflict).JSON(models.NewErrorResponse("Пользователь с таким email уже существует"))
			}

			user.Email = email
			updates = true
		}

		if input.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при хешировании пароля"))
			}
			user.PasswordHash = string(hash)
			updates = true
		}

		if updates {
			if err := db.Save(&user).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при обновлении профиля"))
			}
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Профиль обновлен", user.ToSafeUser()))
	}
}
