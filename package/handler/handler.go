package handler

import (
	"avito_testcase/package/service"
	"net/http"

	_ "avito_testcase/docs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.Handler())

	mux.Handle("/metrics", promhttp.Handler())

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

	//GET
	//.../user_banner
	mux.Handle("/user_banner", HTTPMetrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetBannerByTagAndFeature)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	//GET, POST, PATCH, DELETE
	//.../banner
	mux.Handle("/banner", HTTPMetrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.userIdentity(h.handleGetAllBanners)(w, r)
			return
		} else if r.Method == http.MethodPatch {
			h.userIdentity(h.handleUpdateBanner)(w, r)
			return
		} else if r.Method == http.MethodPost {
			h.userIdentity(h.handleCreateBanner)(w, r)
			return
		} else if r.Method == http.MethodDelete {
			h.userIdentity(h.handleDeleteBanner)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	//PATCH, DELETE
	//.../banner/{id}
	mux.Handle("/banner/", HTTPMetrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			h.userIdentity(h.handleUpdateBanner)(w, r)
			return
		} else if r.Method == http.MethodDelete {
			h.userIdentity(h.handleDeleteBanner)(w, r)
			return
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	return mux
}
