package repository

import (
	avito "github.com/MaksimovDenis/avito_testcase"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user avito.User) (int, error)
	GetUser(username, password string) (avito.User, error)
	GetUserStatus(id int) (bool, error)
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

type Repository struct {
	Authorization
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Banner:        NewBannerPostgres(db),
	}
}
