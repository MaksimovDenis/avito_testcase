package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	avito "github.com/MaksimovDenis/avito_testcase"

	"github.com/sirupsen/logrus"
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

type getBannerByTagAndFeatureResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

// @Summary Получение баннера для пользователя
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param tag_id query int true "tag_id"
// @Param feature_id query int true "feature_id"
// @Param use_last_revision query bool false "use_last_revision: false"
// @Success 200 {object} getBannerByTagAndFeatureResponse
// @Failure 400 {object} Err
// @Failure 401 {object} Err
// @Failure 403 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /user_banner [get]
func (h *Handler) handleGetBannerByTagAndFeature(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Handling Get Banner By Tag And Feature")

	var tag_id, feature_id int
	var use_last_revision bool
	var err error

	tagIDstr := r.URL.Query().Get("tag_id")
	featureIDstr := r.URL.Query().Get("feature_id")
	lastrevisionStr := r.URL.Query().Get("use_last_revision")
	if strings.Contains(tagIDstr, "-") || strings.Contains(featureIDstr, "-") || strings.Contains(lastrevisionStr, "-") {
		logrus.Error("Parameter must be positive")
		NewErrorResponse(w, http.StatusBadRequest, "Parameter must be positive")
		return
	}

	tag_id, err = strconv.Atoi(tagIDstr)
	if err != nil || tag_id < 1 {
		logrus.Error("Invalid Tag ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
		return
	}

	feature_id, err = strconv.Atoi(featureIDstr)
	if err != nil || feature_id < 1 {
		logrus.Error("Invalid Feature ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Feature ID parameter")
		return
	}

	use_last_revision, err = strconv.ParseBool(lastrevisionStr)
	if err != nil {
		logrus.Error("Ivailid use_last_revision parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Ivailid use_last_revision parameter")
		return
	}

	if tag_id == 0 && feature_id == 0 {
		logrus.Error("Invalid Tag ID and Feature ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID and Feature ID parameter")
		return
	}

	params := avito.BannerQueryParams{
		TagID:       tag_id,
		FeatureID:   feature_id,
		LastVersion: use_last_revision,
	}

	if err := h.checkAdminStatus(w, r); err != nil {
		banner, err := h.service.Banner.GetBannerByTagAndFeature(params)
		if err != nil {
			logrus.Error("Failed to Get Banner By TagID and FeatureID", err.Error())
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
			logrus.Error("Failed to encode response", err.Error())
		}

	} else {
		banner, err := h.service.Banner.GetBannerByTagAndFeatureForAdmin(params)
		if err != nil {
			logrus.Error("Failed to Get Banner By TagID and FeatureID", err.Error())
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
			logrus.Error("Failed to encode response", err.Error())
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

// @Summary Получение всех баннеров c фильтрацией по фиче и/или тегу
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param feature_id query int false "feature_id"
// @Param tag_id query int false "tag_id"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} []getAllBannersResponse
// @Failure 400 {object} Err
// @Failure 401 {object} Err
// @Failure 403 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /banner [get]
func (h *Handler) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Handling Get All Banners")

	var tag_id, feature_id, limit, offset int
	var err error

	tagIDstr := r.URL.Query().Get("tag_id")
	featureIDstr := r.URL.Query().Get("feature_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	if strings.Contains(tagIDstr, "-") || strings.Contains(featureIDstr, "-") || strings.Contains(limitStr, "-") || strings.Contains(offsetStr, "-") {
		logrus.Error("Parameter must be positive")
		NewErrorResponse(w, http.StatusBadRequest, "Parameter must be positive")
		return
	}

	if tagIDstr != "" {
		tag_id, err = strconv.Atoi(tagIDstr)
		if err != nil || tag_id < 0 {
			logrus.Error("Invalid Tag ID parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
			fmt.Println(tag_id)
			return
		}
	} else {
		tag_id = 0
	}

	if featureIDstr != "" {
		feature_id, err = strconv.Atoi(featureIDstr)
		if err != nil || feature_id < 0 {
			logrus.Error("Invalid Feature ID parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Feature ID parameter")
			fmt.Println(feature_id)
			return
		}
	} else {
		feature_id = 0
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			logrus.Error("Invalid Limit parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Limit parameter")
			fmt.Println(limit)
			return
		}
	} else {
		limit = 0
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			logrus.Error("Invalid Offset parameter: ", err.Error())
			NewErrorResponse(w, http.StatusBadRequest, "Invalid Offset parameter")
			fmt.Println(offset)
			return
		}
	} else {
		offset = 0
	}

	params := avito.BannerQueryParams{
		TagID:     tag_id,
		FeatureID: feature_id,
		Limit:     limit,
		Offset:    offset,
	}

	var responses []getAllBannersResponse

	if err := h.checkAdminStatus(w, r); err != nil {
		banners, err := h.service.Banner.GetAllBanners(params)
		if err != nil {
			logrus.Error("Failed to Get All Banners", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(banners) == 0 {
			logrus.Info("The list of banners is empty")
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
			logrus.Error("Failed to encode response", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		banners, err := h.service.Banner.GetAllBannersForAdmin(params)
		if err != nil {
			logrus.Error("Failed to Get All Banners", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(banners) == 0 {
			logrus.Info("The list of banners is empty")
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
			logrus.Error("Failed to encode response", err.Error())
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// @Summary Создание нового баннера
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body bannerRequest true "Содержимое баннера"
// @Success 200 {string} string "id"
// @Failure 400 {object} Err
// @Failure 401 {object} Err
// @Failure 403 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /banner [POST]
func (h *Handler) handleCreateBanner(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Handling Create Banner request")

	if err := h.checkAdminStatus(w, r); err != nil {
		logrus.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	var request bannerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logrus.Error("Failed to decode request bodt: ", err.Error())
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
		logrus.Error("Failed to create banner:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, "Failed to create banner")
		return
	}

	response := map[string]interface{}{
		"id": bannerID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Error("Failed to encode response", err.Error())
	}
}

// @Summary Обновление содержимого баннера
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param input body avito.UpdateBanner true "Содержимое баннера"
// @Success 204 {object} StatusResponse
// @Failure 400 {object} Err
// @Failure 401 {object} Err
// @Failure 403 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /banner/{id} [patch]
func (h *Handler) handleUpdateBanner(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Handling Update Banner By BannerID")

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		logrus.Info("Missing ID parameters")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	if err := h.checkAdminStatus(w, r); err != nil {
		logrus.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	bannerIDstr := parts[2]
	id, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		logrus.Error("Invalid Banner ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
		return
	}

	var input avito.UpdateBanner
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logrus.Error("Failed to decode request body: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Banner.UpdateBanner(id, input); err != nil {
		logrus.Error("Failed to update movie: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := StatusResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Error("Failed to encode response", err.Error())
	}

}

// @Summary Удаление баннера по идентификатору
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 204 {object} StatusResponse
// @Failure 400 {object} Err
// @Failure 401 {object} Err
// @Failure 403 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router /banner/{id} [delete]
func (h *Handler) handleDeleteBanner(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Handling Delete Banner by ID")

	if err := h.checkAdminStatus(w, r); err != nil {
		logrus.Error("Admin status is not available:", err.Error())
		NewErrorResponse(w, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		logrus.Info("Missing ID parameters")
		NewErrorResponse(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	bannerIDstr := parts[2]
	id, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		logrus.Error("Invalid Banner ID parameter: ", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, "Invalid Tag ID parameter")
		return
	}

	err = h.service.Banner.DeleteBanner(id)
	if err != nil {
		logrus.Error("Failed to delete banner", err.Error())
		NewErrorResponse(w, http.StatusNoContent, "Баннер не найден")
		return
	}

	response := StatusResponse{Status: "Баннер успешно удален"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Error("Failed to encode response", err.Error())
	}

}
