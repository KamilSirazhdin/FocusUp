package models

// Response представляет стандартизированный ответ API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewSuccessResponse создает новый успешный ответ
func NewSuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse создает новый ответ с ошибкой
func NewErrorResponse(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// AuthResponse представляет ответ для аутентификации
type AuthResponse struct {
	Token string   `json:"token"`
	User  SafeUser `json:"user"`
}

// StatsResponse представляет ответ со статистикой пользователя
type StatsResponse struct {
	Points int `json:"points"`
	Streak int `json:"streak"`
}

// TaskResponse представляет ответ с задачей
type TaskResponse struct {
	ID       uint   `json:"id"`
	Question string `json:"question"`
	Points   int    `json:"points"`
}

// TaskAnswerResponse представляет ответ на решение задачи
type TaskAnswerResponse struct {
	Correct bool `json:"correct"`
	Points  int  `json:"points"`
	Streak  int  `json:"streak"`
}

// UserStatsResponse представляет статистику пользователя для лидерборда
type UserStatsResponse struct {
	Username string `json:"username"`
	Points   int    `json:"points"`
	Streak   int    `json:"streak"`
}
