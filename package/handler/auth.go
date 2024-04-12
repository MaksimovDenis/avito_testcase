package handler

import (
	avito "avito_testcase"
	logger "avito_testcase/logs"
	"encoding/json"
	"net/http"
)

func (h *Handler) handleSingUp(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Handling Sign Up")

	var input avito.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Log.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if input.Username == "" || input.Password == "" {
		errMsg := "Username and password are required"
		logger.Log.Error(errMsg)
		NewErrorResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		logger.Log.Error("Failed to create new user:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

type logInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) handleSingIn(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Log In")

	var input logInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Log.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if input.Username == "" || input.Password == "" {
		errMsg := "Username and password are required"
		logger.Log.Error(errMsg)
		NewErrorResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	token, err := h.service.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		logger.Log.Error("Failed to generate JWT Token:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}

}
