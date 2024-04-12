package handler

import (
	avito "avito_testcase"
	logger "avito_testcase/logs"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type bannerRequest struct {
	TagIDs    []int `json:"tag_ids"`
	FeatureID int   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
	IsActive bool `json:"is_active"`
}

func (h *Handler) handleCreateBanner(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Create Banner request")

	/*if err := h.checkAdminStatus(w, r); err != nil {
		logger.Log.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}*/

	var request bannerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Log.Error("Failed to decode request bodt: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	banner := avito.Banners{
		FeatureId: request.FeatureID,
		Title:     request.Content.Title,
		Text:      request.Content.Text,
		URL:       request.Content.URL,
		IsActive:  request.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	bannerID, err := h.service.Banner.CreateBanner(banner, request.TagIDs)
	if err != nil {
		logger.Log.Error("Failed to create banner:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, "Failed to create banner")
		return
	}

	response := map[string]interface{}{
		"id": bannerID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

type getBannerByTagAndFeatureResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

func (h *Handler) handleGetBannerByTagAndFeature(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get Banner By Tag And Feature")

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		logger.Log.Info("Missing ID parameters")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	tagIdStr := parts[2]
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		logger.Log.Error("Invalid Tag ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
		return
	}

	featureIdStr := parts[3]
	featureId, err := strconv.Atoi(featureIdStr)
	if err != nil {
		logger.Log.Error("Invalid Feature ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Feature ID parameter")
		return
	}

	lastVersionStr := parts[4]
	lastVersion, err := strconv.ParseBool(lastVersionStr)
	if err != nil {
		logger.Log.Error("Invalid Last Version parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Last Version parameter")
		return
	}

	if err := h.checkAdminStatus(w, r); err != nil {
		banner, err := h.service.Banner.GetBannerByTagAndFeature(tagId, featureId, lastVersion)
		if err != nil {
			logger.Log.Error("Failed to Get Banner By TagID and FeatureID", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := getBannerByTagAndFeatureResponse{
			Title: banner.Title,
			Text:  banner.Text,
			URL:   banner.URL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Log.Error("Failed to encode response", err.Error())
		}

	} else {
		banner, err := h.service.Banner.GetBannerByTagAndFeatureForAdmin(tagId, featureId, lastVersion)
		if err != nil {
			logger.Log.Error("Failed to Get Banner By TagID and FeatureID", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := getBannerByTagAndFeatureResponse{
			Title: banner.Title,
			Text:  banner.Text,
			URL:   banner.URL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Log.Error("Failed to encode response", err.Error())
		}
	}

}

type getAllBannersResponse struct {
	BannerID  int   `json:"banner_id" db:"banner_id"`
	TagIDs    []int `json:"tag_ids" db:"tag_ids"`
	FeatureId int   `json:"feature_id" db:"feature_id"`
	Content   struct {
		Title string `json:"title" db:"title"`
		Text  string `json:"text" db:"text"`
		URL   string `json:"url" db:"url"`
	} `json:"content"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (h *Handler) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Get All Banners")

	var tagID, featureID, limit, offset int
	var err error

	tagIDstr := r.URL.Query().Get("tagID")
	featureIDstr := r.URL.Query().Get("featureID")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if tagIDstr != "" {
		tagID, err = strconv.Atoi(tagIDstr)
		if err != nil {
			logger.Log.Error("Invalid Tag ID parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
			return
		}
	} else {
		tagID = 0
	}

	if featureIDstr != "" {
		featureID, err = strconv.Atoi(featureIDstr)
		if err != nil {
			logger.Log.Error("Invalid Feature ID parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Feature ID parameter")
			return
		}
	} else {
		featureID = 0
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			logger.Log.Error("Invalid Limit parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Limit parameter")
			return
		}
	} else {
		limit = 0
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			logger.Log.Error("Invalid Offset parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Offset parameter")
			return
		}
	} else {
		offset = 0
	}

	params := avito.BannerQueryParams{
		TagID:     tagID,
		FeatureID: featureID,
		Limit:     limit,
		Offset:    offset,
	}

	var responses []getAllBannersResponse

	if err := h.checkAdminStatus(w, r); err != nil {
		banners, err := h.service.Banner.GetAllBanners(params)
		if err != nil {
			logger.Log.Error("Failed to Get All Banners", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(banners) == 0 {
			logger.Log.Info("The list of banners is empty")
			NewErrorResponse(w, http.StatusOK, "The list of banners is empty")
			return
		}

		for _, banner := range banners {
			response := getAllBannersResponse{
				BannerID:  banner.BannerID,
				TagIDs:    banner.TagIDs,
				FeatureId: banner.FeatureId,
				Content: struct {
					Title string `json:"title" db:"title"`
					Text  string `json:"text" db:"text"`
					URL   string `json:"url" db:"url"`
				}{
					Title: banner.Title,
					Text:  banner.Text,
					URL:   banner.URL,
				},
				IsActive:  banner.IsActive,
				CreatedAt: banner.CreatedAt,
				UpdatedAt: banner.UpdatedAt,
			}
			responses = append(responses, response)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			logger.Log.Error("Failed to encode response", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		banners, err := h.service.Banner.GetAllBannersForAdmin(params)
		if err != nil {
			logger.Log.Error("Failed to Get All Banners", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(banners) == 0 {
			logger.Log.Info("The list of banners is empty")
			NewErrorResponse(w, http.StatusOK, "The list of banners is empty")
			return
		}

		for _, banner := range banners {
			response := getAllBannersResponse{
				BannerID:  banner.BannerID,
				TagIDs:    banner.TagIDs,
				FeatureId: banner.FeatureId,
				Content: struct {
					Title string `json:"title" db:"title"`
					Text  string `json:"text" db:"text"`
					URL   string `json:"url" db:"url"`
				}{
					Title: banner.Title,
					Text:  banner.Text,
					URL:   banner.URL,
				},
				IsActive:  banner.IsActive,
				CreatedAt: banner.CreatedAt,
				UpdatedAt: banner.UpdatedAt,
			}
			responses = append(responses, response)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			logger.Log.Error("Failed to encode response", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

/*func (h *Handler) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Handling Get All Banners")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		logger.Log.Info("Missing ID parameters")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	featureID, err := strconv.Atoi(parts[2])
	if err != nil {
		logger.Log.Error("Invalid Feature ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Feature ID parameter")
		return
	}

	tagID, err := strconv.Atoi(parts[3])
	if err != nil {
		logger.Log.Error("Invalid Tag ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
		return
	}

	limit, err := strconv.Atoi(parts[4])
	if err != nil {
		logger.Log.Error("Invalid Limit parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Limit parameter")
		return
	}

	offset, err := strconv.Atoi(parts[5])
	if err != nil {
		logger.Log.Error("Invalid Offset parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Offset parameter")
		return
	}

	var responses []getAllBannersResponse

	if err := h.checkAdminStatus(w, r); err != nil {
		banners, err := h.service.Banner.GetAllBanners(featureID, tagID, limit, offset)
		if err != nil {
			logger.Log.Error("Failed to Get All Banners", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(banners) == 0 {
			logger.Log.Info("The list of banners is empty")
			NewErrorResponse(w, http.StatusOK, "The list of banners is empty")
			return
		}

		for _, banner := range banners {
			response := getAllBannersResponse{
				BannerID:  banner.BannerID,
				TagIDs:    banner.TagIDs,
				FeatureId: banner.FeatureId,
				Content: struct {
					Title string `json:"title" db:"title"`
					Text  string `json:"text" db:"text"`
					URL   string `json:"url" db:"url"`
				}{
					Title: banner.Title,
					Text:  banner.Text,
					URL:   banner.URL,
				},
				IsActive:  banner.IsActive,
				CreatedAt: banner.CreatedAt,
				UpdatedAt: banner.UpdatedAt,
			}
			responses = append(responses, response)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			logger.Log.Error("Failed to encode response", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		banners, err := h.service.Banner.GetAllBannersForAdmin(featureID, tagID, limit, offset)
		if err != nil {
			logger.Log.Error("Failed to Get All Banners", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(banners) == 0 {
			logger.Log.Info("The list of banners is empty")
			NewErrorResponse(w, http.StatusOK, "The list of banners is empty")
			return
		}

		for _, banner := range banners {
			response := getAllBannersResponse{
				BannerID:  banner.BannerID,
				TagIDs:    banner.TagIDs,
				FeatureId: banner.FeatureId,
				Content: struct {
					Title string `json:"title" db:"title"`
					Text  string `json:"text" db:"text"`
					URL   string `json:"url" db:"url"`
				}{
					Title: banner.Title,
					Text:  banner.Text,
					URL:   banner.URL,
				},
				IsActive:  banner.IsActive,
				CreatedAt: banner.CreatedAt,
				UpdatedAt: banner.UpdatedAt,
			}
			responses = append(responses, response)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			logger.Log.Error("Failed to encode response", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}*/

func (h *Handler) handleUpdateBanner(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Update Banner By BannerID")

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		logger.Log.Info("Missing ID parameters")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	if err := h.checkAdminStatus(w, r); err != nil {
		logger.Log.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	bannerIDstr := parts[2]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		logger.Log.Error("Invalid Banner ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
		return
	}

	var input avito.UpdateBanner
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Log.Error("Failed to decode request body: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Banner.UpdateBanner(bannerID, input); err != nil {
		logger.Log.Error("Failed to update movie: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := StatusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}

}

func (h *Handler) handleDeleteBanner(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Delete Banner by ID")

	if err := h.checkAdminStatus(w, r); err != nil {
		logger.Log.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "This function is only available to the administrator")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		logger.Log.Info("Missing ID parameters")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	bannerIDstr := parts[2]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		logger.Log.Error("Invalid Banner ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
		return
	}

	err = h.service.Banner.DeleteBanner(bannerID)
	if err != nil {
		logger.Log.Error("Failed to delete banner", err.Error())
		NewErrorResponse(w, http.StatusNoContent, err.Error())
		return
	}

	response := StatusResponse{Status: "Баннер успешно удален"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}

}
