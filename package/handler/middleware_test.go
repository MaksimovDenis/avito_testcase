package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_service "github.com/MaksimovDenis/avito_testcase/package/service/mocks"

	"github.com/MaksimovDenis/avito_testcase/package/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                  string
		headerName            string
		headerValue           string
		token                 string
		mockBehaivior         mockBehaivior
		expectedStatusCode    int
		exptextedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaivior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:    200,
			exptextedResponseBody: "1",
		},
		{
			name:                  "No header",
			headerName:            "",
			token:                 "token",
			mockBehaivior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:    401,
			exptextedResponseBody: `{"error":"Пользователь не авторизован"}`,
		},
		{
			name:                  "Invailed Bearer or Empty Token",
			headerName:            "Authorization",
			token:                 "Bearrtoken",
			mockBehaivior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:    401,
			exptextedResponseBody: `{"error":"Пользователь не авторизован"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaivior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, errors.New("Пользователь не авторизован"))
			},
			expectedStatusCode:    401,
			exptextedResponseBody: `{"error":"Пользователь не авторизован"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehaivior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test server
			mux := http.NewServeMux()
			mux.HandleFunc("/", handler.userIdentity(func(w http.ResponseWriter, r *http.Request) {
				userId := r.Context().Value(userCtx).(int)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", userId)))
			}))

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			//Perform Request
			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.exptextedResponseBody)
			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}
