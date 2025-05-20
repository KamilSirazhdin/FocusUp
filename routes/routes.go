package routes

import (
	"github.com/KamilSirazhdin/FocusUp/controllers"
	"github.com/KamilSirazhdin/FocusUp/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes настраивает маршруты API
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Публичные эндпоинты
	authRoutes := app.Group("/api/auth")
	authRoutes.Post("/register", controllers.Register(db))
	authRoutes.Post("/login", controllers.Login(db))

	// Добавляем новые эндпоинты для сброса пароля
	authRoutes.Post("/forgot-password", controllers.ForgotPassword(db))
	authRoutes.Post("/verify-reset-code", controllers.VerifyResetCode(db))
	authRoutes.Post("/reset-password", controllers.ResetPassword(db))

	// Публичная статистика
	app.Get("/api/stats", controllers.GetStats(db))

	// Маршруты пользователя (требуют аутентификации)
	userRoutes := app.Group("/api/user", middleware.JWTMiddleware(db))
	userRoutes.Get("/profile", controllers.Profile())
	userRoutes.Get("/streak", controllers.StreakStats())
	userRoutes.Get("/history", controllers.GetUserTaskHistory(db))
	userRoutes.Get("/task", controllers.GetRandomTask(db))
	userRoutes.Post("/task/:id/answer", controllers.AnswerTask(db))
	userRoutes.Post("/logout", controllers.Logout())
	userRoutes.Put("/profile", controllers.UpdateProfile(db))

	// Маршруты администратора
	adminRoutes := app.Group("/api/admin",
		middleware.JWTMiddleware(db),
		middleware.RoleMiddleware("admin"))

	adminRoutes.Get("/tasks", controllers.GetAllTasks(db))
	adminRoutes.Post("/task", controllers.CreateTask(db))
	adminRoutes.Put("/task/:id", controllers.UpdateTask(db))
	adminRoutes.Delete("/task/:id", controllers.DeleteTask(db))

	// Управление заданиями
	adminRoutes.Post("/task", controllers.CreateTask(db))
	adminRoutes.Put("/task/:id", controllers.UpdateTask(db))
	adminRoutes.Delete("/task/:id", controllers.DeleteTask(db))

	// Управление пользователями
	adminRoutes.Get("/users", controllers.GetAllUsers(db))
	adminRoutes.Delete("/users/:id", controllers.DeleteUser(db))
}
