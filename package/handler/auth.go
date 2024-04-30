package handler

import (
	"encoding/json"
	"net/http"

	avito "github.com/MaksimovDenis/avito_testcase"

	"github.com/sirupsen/logrus"
)

// @Summary SingUp
// @Description  Create account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body avito.User true "Account info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router       /auth/sing-up [post]
func (h *Handler) handleSingUp(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Handling Sign Up")

	var input avito.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logrus.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if input.Username == "" || input.Password == "" {
		errMsg := "Username and password are required"
		logrus.Error(errMsg)
		NewErrorResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		logrus.Error("Failed to create new user:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Error("Failed to encode response", err.Error())
	}
}

type logInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary LogIn
// @Description  LogIn
// @Tags auth
// @Accept json
// @Produce json
// @Param input body logInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router       /auth/log-in [post]
func (h *Handler) handleSingIn(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Handling Log In")

	var input logInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logrus.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if input.Username == "" || input.Password == "" {
		errMsg := "Username and password are required"
		logrus.Error(errMsg)
		NewErrorResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	token, err := h.service.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		logrus.Error("Failed to generate JWT Token:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Error("Failed to encode response", err.Error())
	}

}
