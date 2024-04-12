package handler

import (
	logger "avito_testcase/logs"
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)

		if header == "" {
			logger.Log.Error("empty auth header")
			NewErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			logger.Log.Error("invalid auth header")
			NewErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		token := headerParts[1]

		if headerParts[0] != "Bearer" {
			logger.Log.Error("invalid auth header")
			NewErrorResponse(w, http.StatusUnauthorized, "invailed auth header")
		}

		if token == "" {
			logger.Log.Error("token is empty")
			NewErrorResponse(w, http.StatusUnauthorized, "token is empty")
		}

		userId, err := h.service.Authorization.ParseToken(token)
		if err != nil {
			NewErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userId)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func getUserId(r *http.Request) (int, error) {
	userId := r.Context().Value(userCtx)
	if userId == nil {
		return 0, errors.New("user id not found")
	}

	idInt, ok := userId.(int)
	if !ok {
		return 0, errors.New("user id is of invailid type")
	}

	return idInt, nil
}

func (h *Handler) checkAdminStatus(w http.ResponseWriter, r *http.Request) error {
	userId, err := getUserId(r)
	if err != nil {
		return err
	}

	user, err := h.service.GetUserStatus(userId)
	if err != nil {
		return err
	}

	if !user {
		return errors.New("user is not admin")
	}

	return nil
}
