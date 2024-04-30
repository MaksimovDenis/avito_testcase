package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/MaksimovDenis/avito_testcase/package/metrics"

	"github.com/MaksimovDenis/avito_testcase/package/helpers"

	"strings"

	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)

		if header == "" {
			logrus.Error("empty auth header")
			NewErrorResponse(w, http.StatusUnauthorized, "Пользователь не авторизован")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			logrus.Error("invalid auth header")
			NewErrorResponse(w, http.StatusUnauthorized, "Пользователь не авторизован")
			return
		}

		token := headerParts[1]

		if headerParts[0] != "Bearer" {
			logrus.Error("invalid auth header")
			NewErrorResponse(w, http.StatusUnauthorized, "Пользователь не авторизован")
		}

		if token == "" {
			logrus.Error("token is empty")
			NewErrorResponse(w, http.StatusUnauthorized, "token is empty")
		}

		userId, err := h.service.Authorization.ParseToken(token)
		if err != nil {
			NewErrorResponse(w, http.StatusUnauthorized, "Пользователь не авторизован")
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
		return 0, errors.New("Пользователь не найден")
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
		return errors.New("Пользователь не имеет доступа")
	}

	return nil
}

func HTTPMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := helpers.NewStatusResponseWriter(w)
		now := time.Now()

		next.ServeHTTP(srw, r)

		elapsedSeconds := time.Since(now).Seconds()
		status := srw.GetStatusString()
		pattern := r.URL.Path
		method := r.Method

		metrics.HttpRequestsTotal.WithLabelValues(pattern, method, status).Inc()
		metrics.HttpRequestsDurationHistogram.WithLabelValues(pattern, method).Observe(elapsedSeconds)
		metrics.HttpRequestsDurationSummary.WithLabelValues(pattern, method).Observe(elapsedSeconds)

	})
}
