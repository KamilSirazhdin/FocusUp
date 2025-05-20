// controllers/password_reset_controller.go
package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"time"

	"github.com/KamilSirazhdin/FocusUp/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ForgotPasswordInput структура для запроса сброса пароля
type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyCodeInput структура для верификации кода
type VerifyCodeInput struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

// ResetPasswordInput структура для изменения пароля
type ResetPasswordInput struct {
	Email       string `json:"email" validate:"required,email"`
	Token       string `json:"token" validate:"required"`
	Code        string `json:"code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func ForgotPassword(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("Получен запрос на сброс пароля")

		var input ForgotPasswordInput
		if err := c.BodyParser(&input); err != nil {
			fmt.Printf("Ошибка парсинга запроса: %v\n", err)
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		fmt.Printf("Запрос на сброс пароля для email: %s\n", input.Email)

		// Проверка: есть ли пользователь с таким email
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			// Чтобы не палить, что пользователь не существует — возвращаем успех
			return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(
				"Если email зарегистрирован, на него был отправлен код сброса пароля", nil))
		}

		// Генерация кода и токена
		code, err := generateRandomCode(6)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка генерации кода"))
		}

		token, err := generateRandomToken(16)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка генерации токена"))
		}

		// Создание записи в таблице password_resets
		reset := models.PasswordReset{
			Email:     input.Email,
			Code:      code,
			Token:     token,
			Used:      false,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		if err := db.Create(&reset).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка сохранения запроса"))
		}

		// Отправка письма
		if err := sendResetEmail(input.Email, code); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Не удалось отправить email"))
		}

		// Возвращаем ответ
		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(
			"Код для сброса пароля отправлен на указанный email", nil))
	}
}

// VerifyResetCode проверяет код, введенный пользователем
func VerifyResetCode(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input VerifyCodeInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		// Проверяем наличие активного запроса на сброс пароля
		var resetRequest models.PasswordReset
		if err := db.Where("email = ? AND code = ? AND used = ? AND expires_at > ?",
			input.Email, input.Code, false, time.Now()).
			Order("created_at DESC").
			First(&resetRequest).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный или устаревший код"))
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(
			"Код подтвержден",
			fiber.Map{
				"token": resetRequest.Token,
				"email": resetRequest.Email,
			}))
	}
}

// ResetPassword меняет пароль пользователя
func ResetPassword(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input ResetPasswordInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный формат запроса"))
		}

		// Проверяем наличие активного запроса на сброс пароля
		var resetRequest models.PasswordReset
		if err := db.Where("email = ? AND token = ? AND code = ? AND used = ? AND expires_at > ?",
			input.Email, input.Token, input.Code, false, time.Now()).
			Order("created_at DESC").
			First(&resetRequest).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Неверный или устаревший запрос"))
		}

		// Находим пользователя
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("Пользователь не найден"))
		}

		// Хешируем новый пароль
		hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при хешировании пароля"))
		}

		// Обновляем пароль пользователя
		user.PasswordHash = string(hash)
		if err := db.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("Ошибка при обновлении пароля"))
		}

		// Помечаем запрос как использованный
		resetRequest.Used = true
		db.Save(&resetRequest)

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse("Пароль успешно изменен", nil))
	}
}

// Генерация случайного токена
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Генерация случайного цифрового кода
func generateRandomCode(length int) (string, error) {
	code := ""
	for i := 0; i < length; i++ {
		// Генерация случайной цифры от 0 до 9
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", n)
	}
	return code, nil
}

// Отправка email с кодом для сброса пароля
func sendResetEmail(email, code string) error {
	// Данные для SMTP сервера
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	// Проверка наличия настроек SMTP
	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" {
		// Для тестирования без отправки email можно вернуть успех
		fmt.Printf("SMTP not configured. Would send reset code %s to %s\n", code, email)
		return nil
	}

	// Содержимое письма
	subject := "Сброс пароля в FocusUp"
	body := fmt.Sprintf(`
Здравствуйте!

Вы запросили сброс пароля в системе FocusUp.
Ваш код для сброса пароля: %s

Код действителен в течение 1 часа.
Если вы не запрашивали сброс пароля, проигнорируйте это сообщение.

С уважением,
Команда FocusUp
`, code)

	message := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body)

	// Аутентификация
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// Отправка письма
	addr := smtpHost + ":" + smtpPort
	err := smtp.SendMail(addr, auth, from, []string{email}, message)
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return err
	}

	return nil
}
