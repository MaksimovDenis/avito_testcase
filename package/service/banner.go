package service

import (
	"github.com/MaksimovDenis/avito_testcase/package/repository"

	avito "github.com/MaksimovDenis/avito_testcase"
)

type BannerService struct {
	repo repository.Banner
}

func NewBannerService(repo repository.Banner) *BannerService {
	return &BannerService{repo: repo}
}

func NewAllBannerService(repo repository.Banner) *BannerService {
	return &BannerService{repo: repo}
}

func (b *BannerService) CreateBanner(banner avito.Banners, tagIDs []int) (int, error) {
	return b.repo.CreateBanner(banner, tagIDs)
}

func (b *BannerService) GetBannerByTagAndFeature(params avito.BannerQueryParams) (avito.Banners, error) {
	return b.repo.GetBannerByTagAndFeature(params)
}

func (b *BannerService) GetBannerByTagAndFeatureForAdmin(params avito.BannerQueryParams) (avito.Banners, error) {
	return b.repo.GetBannerByTagAndFeatureForAdmin(params)
}

func (b *BannerService) GetAllBanners(params avito.BannerQueryParams) ([]avito.AllBanners, error) {
	return b.repo.GetAllBanners(params)
}

func (b *BannerService) GetAllBannersForAdmin(params avito.BannerQueryParams) ([]avito.AllBanners, error) {
	return b.repo.GetAllBannersForAdmin(params)
}

func (b *BannerService) UpdateBanner(bannerID int, input avito.UpdateBanner) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return b.repo.UpdateBanner(bannerID, input)
}

func (b *BannerService) DeleteBanner(bannerID int) error {
	return b.repo.DeleteBanner(bannerID)
}
