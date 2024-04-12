package handler

import (
	avito "avito_testcase"
	"avito_testcase/package/service"
	mock_service "avito_testcase/package/service/mocks"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHadnler_singUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user avito.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           avito.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"test", "is_admin":true}`,
			inputUser: avito.User{
				Username: "test",
				Password: "test",
				Is_admin: true,
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user avito.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"username":"test"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, user avito.User) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Username and password are required"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"test", "password":"test","is_admin":true}`,
			inputUser: avito.User{
				Username: "test",
				Password: "test",
				Is_admin: true,
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user avito.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"service failure"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/auth/sing-up", handler.handleSingUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sing-up",
				bytes.NewBufferString(testCase.inputBody))

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestHandler_handleSingIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, username, password string)

	testTable := []struct {
		name                string
		inputBody           string
		username            string
		password            string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"test"}`,
			username:  "test",
			password:  "test",
			mockBehavior: func(s *mock_service.MockAuthorization, username, password string) {
				s.EXPECT().GenerateToken(username, password).Return("testoken", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"testoken"}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"password":"test"}`,
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Username and password are required"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"test", "password":"test"}`,
			username:  "test",
			password:  "test",
			mockBehavior: func(s *mock_service.MockAuthorization, username, password string) {
				s.EXPECT().GenerateToken(username, password).Return("", errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"service failure"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			if testCase.mockBehavior != nil {
				testCase.mockBehavior(auth, testCase.username, testCase.password)
			}

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			mux := http.NewServeMux()
			mux.HandleFunc("/auth/log-in", handler.handleSingIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/log-in",
				bytes.NewBufferString(testCase.inputBody))

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)
		})
	}
}
