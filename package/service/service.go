package service

import (
	avito "avito_testcase"
	"avito_testcase/package/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user avito.User) (int, error)
	GetUserStatus(id int) (bool, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Banner interface {
	CreateBanner(banner avito.Banners, tagIDs []int) (int, error)
	GetBannerByTagAndFeature(params avito.BannerQueryParams) (avito.Banners, error)
	GetBannerByTagAndFeatureForAdmin(params avito.BannerQueryParams) (avito.Banners, error)
	GetAllBanners(params avito.BannerQueryParams) ([]avito.AllBanners, error)
	GetAllBannersForAdmin(params avito.BannerQueryParams) ([]avito.AllBanners, error)
	UpdateBanner(bannerID int, input avito.UpdateBanner) error
	DeleteBanner(bannerID int) error
}

type Service struct {
	Authorization
	Banner
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Banner:        NewBannerService(repos.Banner),
	}
}
