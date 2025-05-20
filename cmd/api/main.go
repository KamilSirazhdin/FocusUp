package main

import (
	"database/sql"
	"log"

	"github.com/KamilSirazhdin/FocusUp/config"
	"github.com/KamilSirazhdin/FocusUp/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pressly/goose/v3"
)

func main() {
	// Загрузка конфигурации
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Подключение к БД
	db, sqlDB, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Настройка и запуск миграций
	if err := runMigrations(sqlDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация приложения
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Настройка маршрутов
	routes.SetupRoutes(app, db)

	// Запуск сервера
	port := ":" + config.GetPort()
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// runMigrations запускает миграции базы данных с помощью Goose
func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

// customErrorHandler обрабатывает ошибки Fiber и возвращает стандартизированный ответ
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Определение кода ошибки в зависимости от типа ошибки
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
