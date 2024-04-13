package avito

import (
	"errors"
	"time"
)

type Banners struct {
	BannerId  int       `json:"banner_id" db:"banner_id"`
	FeatureId int       `json:"feature_id" db:"feature_id"`
	Title     string    `json:"title" db:"title"`
	Text      string    `json:"text" db:"text"`
	URL       string    `json:"url" db:"url"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Tags struct {
	TagId   int    `json:"tag_id" db:"tag_id"`
	TagName string `json:"tag_name" db:"tag_name"`
}

type Features struct {
	FeatureId   int    `json:"feature_id" db:"feature_id"`
	FeatureName string `json:"feature_name" db:"feature_name"`
}

type BannerTags struct {
	BannerId int `json:"banner_id" db:"banner_id"`
	TagId    int `json:"tag_id" db:"tag_id"`
}

type BannerFeatures struct {
	BannerId  int `json:"banner_id" db:"banner_id"`
	FeatureId int `json:"feature_id" db:"feature_id"`
}

type AllBanners struct {
	BannerID  int       `json:"banner_id" db:"banner_id"`
	TagIDs    []int     `json:"tag_ids" db:"tag_ids"`
	FeatureId int       `json:"feature_id" db:"feature_id"`
	Title     string    `json:"title" db:"title"`
	Text      string    `json:"text" db:"text"`
	URL       string    `json:"url" db:"url"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

/*type BannerTest struct {
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
}*/

type UpdateBanner struct {
	BannerID  *int   `json:"banner_id" db:"banner_id"`
	TagIDs    *[]int `json:"tag_ids" db:"tag_ids"`
	FeatureId *int   `json:"feature_id" db:"feature_id"`
	Content   struct {
		Title *string `json:"title" db:"title"`
		Text  *string `json:"text" db:"text"`
		URL   *string `json:"url" db:"url"`
	} `json:"content"`
	IsActive  *bool      `json:"is_active" db:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

func (u UpdateBanner) Validate() error {
	if u.TagIDs == nil && u.FeatureId == nil && u.Content.Title == nil && u.Content.Text == nil && u.Content.URL == nil && u.IsActive == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type BannerQueryParams struct {
	TagID       int
	FeatureID   int
	Limit       int
	Offset      int
	LastVersion bool
}
