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