package handler

/*
func TestHadnler_GetBannerByTagAndFeature(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBanner, tagID int, featureID int, lastVersion bool)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "OK",
			requestURL: "/user_banner/1/1/true",
			mockBehavior: func(s *mock_service.MockBanner, tagID int, featureID int, lastVersion bool) {
				banner := avito.Banners{
					BannerId:  1,
					FeatureId: 1,
					Title:     "some_title",
					Text:      "some_text",
					URL:       "some_url",
					IsActive:  true,
				}
				s.EXPECT().GetBannerByTagAndFeature(1, 1, true).Return(banner, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"title":"some_title","text":"some_text","url":"some_url"}`,
		},
		{
			name:       "Invailid ID parameter",
			requestURL: "/user_banner/1/1",
			mockBehavior: func(s *mock_service.MockBanner, tagID, featureID int, lastVersion bool) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"missing id parameter"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bannerService := mock_service.NewMockBanner(c)
			testCase.mockBehavior(bannerService, 1, 1, true)

			services := &service.Service{Banner: bannerService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/user_banner/", handler.handleGetBannerByTagAndFeature)

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
	type mockBehavior func(s *mock_service.MockBanner, tagID int, featureID int, limit int, offset int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "OK",
			requestURL: "/banner/1/1/10/0",
			mockBehavior: func(s *mock_service.MockBanner, tagID int, featureID int, limit int, offset int) {
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
				s.EXPECT().GetAllBanners(1, 1, 10, 0).Return(banners, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"banner_id":1,"tag_ids":[1],"feature_id":1,"content":{"title":"some_title","text":"some_text","url":"some_url"},"is_active":true,"created_at":"2024-04-12T09:55:56.927Z","updated_at":"2024-04-12T09:55:56.927Z"}]`,
		},
		{
			name:       "Invailid ID parameter",
			requestURL: "/banner/1/1",
			mockBehavior: func(s *mock_service.MockBanner, tagID int, featureID int, limit int, offset int) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"missing id parameter"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bannerService := mock_service.NewMockBanner(c)
			testCase.mockBehavior(bannerService, 1, 1, 10, 0)

			services := &service.Service{Banner: bannerService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/banner/", handler.handleGetAllBanners)

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
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
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
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
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
*/
