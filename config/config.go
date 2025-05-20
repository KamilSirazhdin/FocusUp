package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB глобальный экземпляр подключения к базе данных
var DB *gorm.DB

// LoadEnv загружает переменные окружения из .env.example файла
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("ошибка загрузки .env.example файла: %w", err)
	}
	return nil
}

// ConnectDB устанавливает соединение с базой данных
func ConnectDB() (*gorm.DB, *sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Настройка логгера GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Открытие соединения с БД
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Получение базового SQL-соединения для миграций
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка получения SQL DB: %w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	fmt.Println("Успешное подключение к БД!")
	return DB, sqlDB, nil
}

// GetDB возвращает экземпляр соединения с базой данных
func GetDB() *gorm.DB {
	return DB
}
