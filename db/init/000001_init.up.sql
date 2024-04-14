CREATE TABLE Banners (
    banner_id SERIAL PRIMARY KEY, 
    feature_id INTEGER NOT NULL, 
    title VARCHAR(255) NOT NULL, 
    text TEXT NOT NULL,
    url VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Tags (
    tag_id SERIAL PRIMARY KEY,
    tag_name VARCHAR(100) NOT NULL
);

CREATE TABLE Features (
    feature_id SERIAL PRIMARY KEY,
    feature_name VARCHAR(100) NOT NULL
);

CREATE TABLE BannerTags (
    banner_id INTEGER REFERENCES Banners(banner_id),
    tag_id INTEGER REFERENCES Tags(tag_id),
    PRIMARY KEY (banner_id, tag_id)
);

CREATE TABLE BannerFeatures (
    banner_id INTEGER REFERENCES Banners(banner_id),
    feature_id INTEGER REFERENCES Features(feature_id),
    PRIMARY KEY (banner_id, feature_id)
);

CREATE TABLE Users(
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL, 
    is_admin BOOLEAN NOT NULL
);


DO $$
BEGIN
    FOR i IN 1..5000 LOOP
        -- Вставка данных в таблицу Banners
        INSERT INTO Banners (feature_id, title, text, url)
        VALUES (i, 'Title ' || i, 'Text ' || i, 'http://example.com/banner' || i);

        -- Вставка данных в таблицу Tags
        INSERT INTO Tags (tag_name)
        VALUES ('Tag ' || i);

        -- Вставка данных в таблицу Features
        INSERT INTO Features (feature_name)
        VALUES ('Feature ' || i);

        -- Вставка данных в таблицу BannerTags (связь между Banners и Tags)
        INSERT INTO BannerTags (banner_id, tag_id)
        VALUES (i, i); -- Примерное предположение о связи между Banners и Tags

        -- Вставка данных в таблицу BannerFeatures (связь между Banners и Features)
        INSERT INTO BannerFeatures (banner_id, feature_id)
        VALUES (i, i); -- Примерное предположение о связи между Banners и Features
    END LOOP;
END$$;