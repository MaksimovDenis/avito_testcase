package handler

import (
	avito "avito_testcase"
	"avito_testcase/package/service"
	mock_service "avito_testcase/package/service/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHadnler_GetBannerByTagAndFeature(t *testing.T) {
	testTable := []struct {
		name                string
		requestURL          string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "OK",
			requestURL:          "/user_banner?tag_id=1&feature_id=1&use_last_revision=false",
			expectedStatusCode:  200,
			expectedRequestBody: `{"title":"some_title","text":"some_text","url":"some_url"}`,
		},
		{
			name:                "Invalid Feature ID parameter",
			requestURL:          "/user_banner?tag_id=1&feature_id=D&use_last_revision=false",
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Invalid Feature ID parameter"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bannerService := mock_service.NewMockBanner(c)
			if testCase.name == "OK" {
				// Ожидаемый вызов метода GetBannerByTagAndFeature для случая "OK"
				banner := avito.Banners{
					BannerId:  1,
					FeatureId: 1,
					Title:     "some_title",
					Text:      "some_text",
					URL:       "some_url",
					IsActive:  true,
				}
				expectedParams := avito.BannerQueryParams{
					TagID:       1,
					FeatureID:   1,
					LastVersion: false,
				}
				bannerService.EXPECT().GetBannerByTagAndFeature(expectedParams).Return(banner, nil)
			}

			services := &service.Service{Banner: bannerService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/user_banner", handler.handleGetBannerByTagAndFeature)

			req := httptest.NewRequest("GET", testCase.requestURL, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestHadnler_GetAllBanners(t *testing.T) {
	testTable := []struct {
		name                string
		requestURL          string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "OK",
			requestURL:          "/banner?tag_id=1&feature_id=1&limit=10&offset0",
			expectedStatusCode:  200,
			expectedRequestBody: `[{"banner_id":1,"tag_ids":[1],"feature_id":1,"content":{"title":"some_title","text":"some_text","url":"some_url"},"is_active":true,"created_at":"2024-04-12T09:55:56.927Z","updated_at":"2024-04-12T09:55:56.927Z"}]`,
		},
		{
			name:                "Positive parameters",
			requestURL:          "/banner?tag_id=-1&feature_id=0&limit=0&offset=0",
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Parameter must be positive"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bannerService := mock_service.NewMockBanner(c)
			if testCase.name == "OK" {
				createdAt, _ := time.Parse(time.RFC3339, "2024-04-12T09:55:56.927Z")
				updatedAt, _ := time.Parse(time.RFC3339, "2024-04-12T09:55:56.927Z")

				banners := []avito.AllBanners{
					{
						BannerID:  1,
						TagIDs:    []int{1},
						FeatureId: 1,
						Title:     "some_title",
						Text:      "some_text",
						URL:       "some_url",
						IsActive:  true,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
				}
				expectedParams := avito.BannerQueryParams{
					TagID:     1,
					FeatureID: 1,
					Limit:     10,
					Offset:    0,
				}
				bannerService.EXPECT().GetAllBanners(expectedParams).Return(banners, nil)
			}
			services := &service.Service{Banner: bannerService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/banner", handler.handleGetAllBanners)

			req := httptest.NewRequest("GET", testCase.requestURL, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestHandler_handleUpdateBanner(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBanner, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			requestURL:          "/banner/1",
			mockBehavior:        func(s *mock_service.MockBanner, id int) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"Пользователь не имеет доступа"}`,
		},
		{
			name:                "Missing ID Parameter",
			requestURL:          "/banner",
			mockBehavior:        func(s *mock_service.MockBanner, id int) {},
			expectedStatusCode:  301,
			expectedRequestBody: ``,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			bannerService := mock_service.NewMockBanner(c)
			testCase.mockBehavior(bannerService, 1)

			services := &service.Service{Banner: bannerService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/banner/", handler.handleUpdateBanner)

			req := httptest.NewRequest("PATCH", testCase.requestURL, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleDeleteBanner(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBanner, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			requestURL:          "/banner/1",
			mockBehavior:        func(s *mock_service.MockBanner, id int) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"Пользователь не имеет доступа"}`,
		},
		{
			name:                "Missing ID Parameter",
			requestURL:          "/banner",
			mockBehavior:        func(s *mock_service.MockBanner, id int) {},
			expectedStatusCode:  301,
			expectedRequestBody: ``,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			bannerService := mock_service.NewMockBanner(c)
			testCase.mockBehavior(bannerService, 1)

			services := &service.Service{Banner: bannerService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/banner/", handler.handleDeleteBanner)

			req := httptest.NewRequest("DELETE", testCase.requestURL, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}
