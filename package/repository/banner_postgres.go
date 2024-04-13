package repository

import (
	avito "avito_testcase"
	logger "avito_testcase/logs"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type BannerPostgres struct {
	db *sqlx.DB
}

func NewBannerPostgres(db *sqlx.DB) *BannerPostgres {
	return &BannerPostgres{db: db}
}

func (b *BannerPostgres) CreateBanner(banner avito.Banners, tagIDs []int) (int, error) {

	//CHECK BANNER
	tagIDStr := ""
	for i, id := range tagIDs {
		if i > 0 {
			tagIDStr += ","
		}
		tagIDStr += fmt.Sprintf("%d", id)
	}

	query := fmt.Sprintf(`
		SELECT bt.banner_id FROM %s bt
		JOIN %s b ON bt.banner_id = b.banner_id
		WHERE b.feature_id=$1
		AND bt.tag_id IN (%s)
	`, bannerTagsTable, bannersTable, tagIDStr)

	var existingID int
	row := b.db.QueryRow(query, banner.FeatureId)
	if err := row.Scan(&existingID); err == nil {
		return existingID, errors.New("Banner with the same feature_id and tag_ids already exists")
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	//START TRANSACTION
	tx, err := b.db.Begin()
	if err != nil {
		logger.Log.Errorf("failed to start transaction: %v", err)
		return 0, err
	}
	defer func() {
		if err != nil {
			logger.Log.Errorf("rolling back transaction due to error: %v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Log.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	var id int
	currentTime := time.Now()

	//INSERT INTO bannersTable
	query = fmt.Sprintf(`INSERT INTO %s (feature_id, title, text, url, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING banner_id`, bannersTable)
	row = tx.QueryRow(query, banner.FeatureId, banner.Title, banner.Text, banner.URL, banner.IsActive, currentTime, currentTime)
	if err := row.Scan(&id); err != nil {
		logger.Log.Errorf("error occured during CreateBanner to DB (bannersTable): %v", err)
		return 0, err
	}

	//INSERT INTO featureTable
	query = fmt.Sprintf(`INSERT INTO %s (feature_id, feature_name) VALUES ($1, $2) ON CONFLICT DO NOTHING`, featureTable)
	_, err = tx.Exec(query, banner.FeatureId, strconv.Itoa(banner.FeatureId))
	if err != nil {
		logger.Log.Errorf("error occurred during CreateBanner to DB (featureTable): %v", err)
		return 0, err
	}

	//INSERT INTO bannerFeaturesTable
	query = fmt.Sprintf(`INSERT INTO %s (banner_id, feature_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, bannerFeaturesTable)
	_, err = tx.Exec(query, id, banner.FeatureId)
	if err != nil {
		logger.Log.Errorf("error occurred during CreateBanner to DB (bannerFeaturesTable): %v", err)
		return 0, err
	}

	//INSERT INTO tagsTable AND bannerTagsTable
	for _, tagID := range tagIDs {
		query = fmt.Sprintf(`INSERT INTO %s (tag_id, tag_name) VALUES ($1, $2) ON CONFLICT DO NOTHING`, tagsTable)
		if _, err := tx.Exec(query, tagID, strconv.Itoa(tagID)); err != nil {
			logger.Log.Errorf("error occurred during CreateBanner to DB (tagsTable): %v", err)
			return 0, err
		}

		query = fmt.Sprintf(`INSERT INTO %s (banner_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, bannerTagsTable)
		if _, err := tx.Exec(query, id, tagID); err != nil {
			logger.Log.Errorf("error occurred during CreateBanner to DB (bannerTagsTable): %v", err)
			return 0, err
		}
	}

	//END TRANSACTION
	if err := tx.Commit(); err != nil {
		logger.Log.Errorf("failed to commit transaction: %v", err)
		return 0, err
	}

	return id, nil
}

func (b *BannerPostgres) GetBannerByTagAndFeature(params avito.BannerQueryParams) (avito.Banners, error) {
	ctx := context.Background()
	key := "bannerByTagAndFeature:" + strconv.Itoa(params.TagID) + ":" + strconv.Itoa(params.FeatureID) + ":" + strconv.FormatBool(params.LastVersion)

	if !params.LastVersion {
		//CHECK KEY EXISTENCE IN CACHE
		exists, err := ClientRedis.Exists(ctx, key).Result()
		if err != nil {
			logger.Log.Error("Failed to check key existence in Redis (GetBannerByTagAndFeature)", err.Error())
			return avito.Banners{}, err
		}

		if exists > 0 {
			value, err := ClientRedis.Get(ctx, key).Result()
			if err != nil {
				logger.Log.Error("Failed to get data from Redis (GetBannerByTagAndFeature)", err.Error())
				return avito.Banners{}, err
			}

			var banner avito.Banners
			err = json.Unmarshal([]byte(value), &banner)
			if err != nil {
				logger.Log.Error("Failed to unmarshal JSON (GetBannerByTagAndFeature)", err.Error())
				return avito.Banners{}, err
			}

			return banner, nil
		}
	}

	var banner avito.Banners
	isActive := true
	query := fmt.Sprintf(`
	SELECT 
		b.title, 
		b.text, 
		b.url 
	FROM 
		%s b 
	LEFT JOIN 
		%s t ON b.banner_id=t.banner_id 
	WHERE 
		t.tag_id=$1 
	AND 
		b.feature_id=$2
	AND	
		b.is_active=$3`, bannersTable, bannerTagsTable)
	err := b.db.Get(&banner, query, params.TagID, params.FeatureID, isActive)
	if err != nil {
		logger.Log.Error("Failed to get banner by tag and feature", err.Error())
		return banner, err
	}

	value, err := json.Marshal(banner)
	if err != nil {
		logger.Log.Error("Failed to marshal JSON (GetBannerByTagAndFeature)", err.Error())
		return banner, err
	}

	err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
	if err != nil {
		logger.Log.Error("Failed to set data to redis", err.Error())
		return banner, err
	}

	return banner, nil
}

func (b *BannerPostgres) GetBannerByTagAndFeatureForAdmin(params avito.BannerQueryParams) (avito.Banners, error) {
	ctx := context.Background()
	key := "bannerByTagAndFeature:" + strconv.Itoa(params.TagID) + ":" + strconv.Itoa(params.FeatureID) + ":" + strconv.FormatBool(params.LastVersion)

	if !params.LastVersion {
		//CHECK KEY EXISTENCE IN CACHE
		exists, err := ClientRedis.Exists(ctx, key).Result()
		if err != nil {
			logger.Log.Error("Failed to check key existence in Redis", err.Error())
			return avito.Banners{}, err
		}

		if exists > 0 {
			value, err := ClientRedis.Get(ctx, key).Result()
			if err != nil {
				logger.Log.Error("Failed to get data from Redis", err.Error())
				return avito.Banners{}, err
			}

			var banner avito.Banners
			err = json.Unmarshal([]byte(value), &banner)
			if err != nil {
				logger.Log.Error("Failed to unmarshal JSON (GetBannerByTagAndFeature)", err.Error())
				return avito.Banners{}, err
			}

			return banner, nil
		}
	}

	var banner avito.Banners
	query := fmt.Sprintf(`
	SELECT 
		b.title, 
		b.text, 
		b.url 
	FROM 
		%s b 
	LEFT JOIN 
		%s t ON b.banner_id=t.banner_id 
	WHERE 
		t.tag_id=$1 
	AND 
		b.feature_id=$2`, bannersTable, bannerTagsTable)
	err := b.db.Get(&banner, query, params.TagID, params.FeatureID)
	if err != nil {
		logger.Log.Error("Failed to get banner by tag and feature", err.Error())
		return banner, err
	}

	value, err := json.Marshal(banner)
	if err != nil {
		logger.Log.Error("Failed to marshal JSON (GetBannerByTagAndFeature)", err.Error())
		return banner, err
	}

	err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
	if err != nil {
		logger.Log.Error("Failed to set data to redis", err.Error())
		return banner, err
	}

	return banner, nil
}

/*func (b *BannerPostgres) GetAllBanners(featureID int, tagID int, limit int, offset int) ([]avito.AllBanners, error) {
	ctx := context.Background()
	key := "getAllBanners:" + strconv.Itoa(tagID) + ":" + strconv.Itoa(featureID) + ":" + strconv.Itoa(limit) + ":" + strconv.Itoa(offset)

	//CHECK KEY EXISTENCE IN CACHE
	exists, err := ClientRedis.Exists(ctx, key).Result()
	if err != nil {
		logger.Log.Error("Failed ot check key existence in Redis (GetAllBanners)")
		return []avito.AllBanners{}, err
	}

	if exists > 0 {
		value, err := ClientRedis.Get(ctx, key).Result()
		if err != nil {
			logger.Log.Error("Failed to get data from Redis (GetAllBanners)", err.Error())
			return []avito.AllBanners{}, err
		}

		var banners []avito.AllBanners
		err = json.Unmarshal([]byte(value), &banners)
		if err != nil {
			logger.Log.Error("Failed to unmarshal JSON (GetBannerByTagAndFeature)", err.Error())
			return []avito.AllBanners{}, err
		}
		return banners, nil
	}

	var banners []avito.AllBanners
	isActive := true

	bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title,
		b.text,
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    WHERE
        b.feature_id=$1
	AND
		b.is_active=$2
    OFFSET $3
    LIMIT $4
    `, bannersTable)

	err = b.db.Select(&banners, bannersQuery, featureID, isActive, offset, limit)
	if err != nil {
		return nil, err
	}

	for i := range banners {
		tagsQuery := fmt.Sprintf(`
        SELECT
            t.tag_id
        FROM
            %s t
        WHERE
            t.banner_id=$1
        `, bannerTagsTable)

		var tagIDs []int
		err := b.db.Select(&tagIDs, tagsQuery, banners[i].BannerID)
		if err != nil {
			return nil, err
		}

		banners[i].TagIDs = tagIDs
	}

	value, err := json.Marshal(banners)
	if err != nil {
		logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
		return banners, err
	}

	err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
	if err != nil {
		logger.Log.Error("Failed to set data to redis", err.Error())
		return banners, err
	}

	return banners, nil
}*/

func (b *BannerPostgres) GetAllBanners(params avito.BannerQueryParams) ([]avito.AllBanners, error) {
	ctx := context.Background()
	key := "getAllBanners:" + strconv.Itoa(params.TagID) + ":" + strconv.Itoa(params.FeatureID) + ":" + strconv.Itoa(params.Limit) + ":" + strconv.Itoa(params.Offset)

	//CHECK KEY EXISTENCE IN CACHE
	exists, err := ClientRedis.Exists(ctx, key).Result()
	if err != nil {
		logger.Log.Error("Failed ot check key existence in Redis (GetAllBanners)")
		return []avito.AllBanners{}, err
	}

	if exists > 0 {
		value, err := ClientRedis.Get(ctx, key).Result()
		if err != nil {
			logger.Log.Error("Failed to get data from Redis (GetAllBanners)", err.Error())
			return []avito.AllBanners{}, err
		}

		var banners []avito.AllBanners
		err = json.Unmarshal([]byte(value), &banners)
		if err != nil {
			logger.Log.Error("Failed to unmarshal JSON (GetBannerByTagAndFeature)", err.Error())
			return []avito.AllBanners{}, err
		}
		return banners, nil
	}

	var banners []avito.AllBanners
	isActive := true

	if params.TagID != 0 && params.FeatureID != 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    WHERE
        b.feature_id=$1
	AND
		b.is_active=$2
    OFFSET $3
    LIMIT CASE WHEN $4=0 THEN 10000 ELSE $4 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, params.FeatureID, isActive, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		for i := range banners {
			tagsQuery := fmt.Sprintf(`
        SELECT
            t.tag_id
        FROM
            %s t
        WHERE
            t.banner_id=$1
        `, bannerTagsTable)

			var tagIDs []int
			err := b.db.Select(&tagIDs, tagsQuery, banners[i].BannerID)
			if err != nil {
				return nil, err
			}

			banners[i].TagIDs = tagIDs
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	} else if params.TagID != 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    WHERE
		b.is_active=$1
    OFFSET $2
    LIMIT CASE WHEN $3=0 THEN 10000 ELSE $3 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, isActive, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		for i := range banners {
			tagsQuery := fmt.Sprintf(`
        SELECT
            t.tag_id
        FROM
            %s t
        WHERE
            t.banner_id=$1
        `, bannerTagsTable)

			var tagIDs []int
			err := b.db.Select(&tagIDs, tagsQuery, banners[i].BannerID)
			if err != nil {
				return nil, err
			}

			banners[i].TagIDs = tagIDs
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	} else if params.FeatureID != 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    WHERE
        b.feature_id=$1
	AND
		b.is_active=$2
    OFFSET $3
    LIMIT CASE WHEN $4=0 THEN 10000 ELSE $4 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, params.FeatureID, isActive, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	} else if params.FeatureID == 0 && params.TagID == 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    OFFSET $1
    LIMIT CASE WHEN $2=0 THEN 10000 ELSE $2 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	}

	return banners, nil
}

func (b *BannerPostgres) GetAllBannersForAdmin(params avito.BannerQueryParams) ([]avito.AllBanners, error) {
	ctx := context.Background()
	key := "getAllBanners:" + strconv.Itoa(params.TagID) + ":" + strconv.Itoa(params.FeatureID) + ":" + strconv.Itoa(params.Limit) + ":" + strconv.Itoa(params.Offset)

	//CHECK KEY EXISTENCE IN CACHE
	exists, err := ClientRedis.Exists(ctx, key).Result()
	if err != nil {
		logger.Log.Error("Failed ot check key existence in Redis (GetAllBanners)")
		return []avito.AllBanners{}, err
	}

	if exists > 0 {
		value, err := ClientRedis.Get(ctx, key).Result()
		if err != nil {
			logger.Log.Error("Failed to get data from Redis (GetAllBanners)", err.Error())
			return []avito.AllBanners{}, err
		}

		var banners []avito.AllBanners
		err = json.Unmarshal([]byte(value), &banners)
		if err != nil {
			logger.Log.Error("Failed to unmarshal JSON (GetBannerByTagAndFeature)", err.Error())
			return []avito.AllBanners{}, err
		}
		return banners, nil
	}

	var banners []avito.AllBanners

	if params.TagID != 0 && params.FeatureID != 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    WHERE
        b.feature_id=$1
    OFFSET $2
    LIMIT CASE WHEN $3=0 THEN 10000 ELSE $3 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, params.FeatureID, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		for i := range banners {
			tagsQuery := fmt.Sprintf(`
        SELECT
            t.tag_id
        FROM
            %s t
        WHERE
            t.banner_id=$1
        `, bannerTagsTable)

			var tagIDs []int
			err := b.db.Select(&tagIDs, tagsQuery, banners[i].BannerID)
			if err != nil {
				return nil, err
			}

			banners[i].TagIDs = tagIDs
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	} else if params.TagID != 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    OFFSET $1
    LIMIT CASE WHEN $2=0 THEN 10000 ELSE $2 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		for i := range banners {
			tagsQuery := fmt.Sprintf(`
        SELECT
            t.tag_id
        FROM
            %s t
        WHERE
            t.banner_id=$1
        `, bannerTagsTable)

			var tagIDs []int
			err := b.db.Select(&tagIDs, tagsQuery, banners[i].BannerID)
			if err != nil {
				return nil, err
			}

			banners[i].TagIDs = tagIDs
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	} else if params.FeatureID != 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    WHERE
        b.feature_id=$1
    OFFSET $2
    LIMIT CASE WHEN $3=0 THEN 10000 ELSE $3 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, params.FeatureID, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	} else if params.FeatureID == 0 && params.TagID == 0 {
		bannersQuery := fmt.Sprintf(`
    SELECT
        b.banner_id,
        b.feature_id,
        b.title, 
		b.text, 
		b.url,
        b.is_active,
        b.created_at,
        b.updated_at
    FROM
        %s b
    OFFSET $1
    LIMIT CASE WHEN $2=0 THEN 10000 ELSE $2 END
    `, bannersTable)

		err = b.db.Select(&banners, bannersQuery, params.Offset, params.Limit)
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(banners)
		if err != nil {
			logger.Log.Error("Failed to marshal JSON (GetAllBanners)", err.Error())
			return banners, err
		}

		err = ClientRedis.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
		if err != nil {
			logger.Log.Error("Failed to set data to redis", err.Error())
			return banners, err
		}
	}

	return banners, nil
}

func (b *BannerPostgres) UpdateBanner(bannerID int, input avito.UpdateBanner) error {

	tagIDStr := ""
	for i, id := range *input.TagIDs {
		if i > 0 {
			tagIDStr += ","
		}
		tagIDStr += fmt.Sprintf("%d", id)
	}

	query := fmt.Sprintf(`
	SELECT bt.banner_id FROM %s bt
	JOIN %s b ON bt.banner_id = b.banner_id
	WHERE b.feature_id=$1
	AND bt.tag_id IN (%s)
`, bannerTagsTable, bannersTable, tagIDStr)

	var existingID int
	row := b.db.QueryRow(query, input.FeatureId)
	if err := row.Scan(&existingID); err == nil {
		return errors.New("Banner with the same feature_id and tag_ids already exists")
	} else if err != sql.ErrNoRows {
		return err
	}

	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.FeatureId != nil {
		setValue = append(setValue, fmt.Sprintf("feature_id=$%d", argID))
		args = append(args, *input.FeatureId)
		argID++

		deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE banner_id=$1`, bannerFeaturesTable)
		_, err := b.db.Exec(deleteQuery, bannerID)
		if err != nil {
			logger.Log.Errorf("Error occured during delete string from bannerFeaturesTable (UpdateBanner): %v", err)
			return err
		}

		insertQuery := fmt.Sprintf(`INSERT INTO %s (banner_id, feature_id) VALUES ($1, $2)`, bannerFeaturesTable)
		_, err = b.db.Exec(insertQuery, bannerID, *input.FeatureId)
		if err != nil {
			logger.Log.Errorf("Error occured during insert updtated info to bannerFeaturesTable(UpdateBanner): %v", err)
			return err
		}
	}

	if input.Content.Title != nil {
		setValue = append(setValue, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Content.Title)
		argID++
	}

	if input.Content.Text != nil {
		setValue = append(setValue, fmt.Sprintf("text=$%d", argID))
		args = append(args, *input.Content.Text)
		argID++
	}

	if input.Content.URL != nil {
		setValue = append(setValue, fmt.Sprintf("url=$%d", argID))
		args = append(args, *input.Content.URL)
		argID++
	}

	if input.IsActive != nil {
		setValue = append(setValue, fmt.Sprintf("is_active=$%d", argID))
		args = append(args, *input.IsActive)
		argID++
	}

	if input.TagIDs != nil {
		deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE banner_id=$1`, bannerTagsTable)
		_, err := b.db.Exec(deleteQuery, bannerID)
		if err != nil {
			logger.Log.Errorf("Error occured during delete string from bannerTagsTable (UpdateBanner): %v", err)
			return err
		}

		for _, tagID := range *input.TagIDs {
			insertQuery := fmt.Sprintf(`INSERT INTO %s (banner_id, tag_id) VALUES ($1, $2)`, bannerTagsTable)
			_, err := b.db.Exec(insertQuery, bannerID, tagID)
			if err != nil {
				logger.Log.Errorf("Error occured during insert updtated info to bannerTagsTable(UpdateBanner): %v", err)
				return err
			}
		}
	}

	setQuery := strings.Join(setValue, ", ")
	setQuery += fmt.Sprintf(", updated_at=$%d", argID)
	args = append(args, time.Now())
	argID++

	query = fmt.Sprintf("UPDATE %s SET %s WHERE banner_id=$%d", bannersTable, setQuery, argID)

	args = append(args, bannerID)

	logger.Log.Printf("updateQuerry: %s", query)
	logger.Log.Printf("args :%v", args)

	_, err := b.db.Exec(query, args...)
	if err != nil {
		logger.Log.Errorf("failed to execute update query: %v", err)
		return fmt.Errorf("failed to execute update query: %v", err)
	}

	return nil
}

func (b *BannerPostgres) DeleteBanner(bannerID int) error {

	//CHECK BANNER
	var count int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE banner_id=$1`, bannersTable)
	err := b.db.QueryRow(countQuery, bannerID).Scan(&count)
	if err != nil {
		logger.Log.Errorf("Failed to query banner count: %v", err)
		return err
	}

	if count == 0 {
		return fmt.Errorf("banner with ID %d not found", bannerID)
	}

	//START TRANSACTION
	tx, err := b.db.Begin()
	if err != nil {
		logger.Log.Errorf("failed to start transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			logger.Log.Errorf("rolling back transaction due to error: %v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Log.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	query := fmt.Sprintf(`DELETE FROM %s WHERE banner_id=$1`, bannerFeaturesTable)
	if _, err := b.db.Exec(query, bannerID); err != nil {
		return err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE banner_id=$1`, bannerTagsTable)
	if _, err = b.db.Exec(query, bannerID); err != nil {
		return err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE banner_id=$1`, bannersTable)
	if _, err = b.db.Exec(query, bannerID); err != nil {
		return err
	}

	//END TRANSACTION
	if err := tx.Commit(); err != nil {
		logger.Log.Errorf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}
