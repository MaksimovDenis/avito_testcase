package service

import (
	avito "avito_testcase"
	"avito_testcase/package/repository"
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

func (b *BannerService) GetBannerByTagAndFeature(tagID int, featureID int, lastVersion bool) (avito.Banners, error) {
	return b.repo.GetBannerByTagAndFeature(tagID, featureID, lastVersion)
}

func (b *BannerService) GetBannerByTagAndFeatureForAdmin(tagID int, featureID int, lastVersion bool) (avito.Banners, error) {
	return b.repo.GetBannerByTagAndFeatureForAdmin(tagID, featureID, lastVersion)
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
