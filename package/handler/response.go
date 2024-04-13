package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

var ErrBadRequest = Err{Message: "Некорректные данные"}

var ErrUnauthorized = Err{Message: "Пользователь не авторизован"}

var ErrForbidden = Err{Message: "Пользователь не имеет доступа"}

var ErrNotFound = Err{Message: "Баннер не найден"}

var ErrInternalServerError = Err{Message: "Внутренняя ошибка сервера"}

type Err struct {
	Message string `json:"message"`
}

func (e Err) Error() string {
	return e.Message
}

type StatusResponse struct {
	Status string `json:"status"`
}

func NewErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"error": message,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
