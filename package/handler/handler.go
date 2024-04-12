package handler

import (
	"avito_testcase/package/service"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	auth := "/auth"

	mux.HandleFunc(auth+"/sing-up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.handleSingUp(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(auth+"/log-in", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.handleSingIn(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//POST
	//.../banner/{bannerID}
	mux.HandleFunc("/banner", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.userIdentity(h.handleCreateBanner)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET
	//.../user_banner/{tagID}/{featureID}/{lastVersion}
	mux.HandleFunc("/user_banner/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetBannerByTagAndFeature)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	//GET, PATCH, DELETE
	//.../banner/{featureID}/{tagID}/{limit}/{offset} - GET
	//.../banner/{bannerID} - PATCH, DELETE
	mux.HandleFunc("/banner/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetAllBanners)(w, r)
			return
		} else if r.Method == http.MethodPatch {
			h.userIdentity(h.handleUpdateBanner)(w, r)
			return
		} else if r.Method == http.MethodDelete {
			h.userIdentity(h.handleDeleteBanner)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
