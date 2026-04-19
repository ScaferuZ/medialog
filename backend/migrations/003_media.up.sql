-- Media table (polymorphic base for films, anime, books, manga, games, doujin)
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('film', 'anime', 'book', 'manga', 'game', 'doujin')),
    title VARCHAR(500) NOT NULL,
    original_title VARCHAR(500),
    description TEXT,
    cover_image TEXT,
    release_date DATE,
    metadata JSONB,
    tmdb_id VARCHAR(50),
    mal_id INTEGER,
    google_books_id VARCHAR(100),
    igdb_id INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Genres table
CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Media-Genre relationship (many-to-many)
CREATE TABLE media_genres (
    media_id UUID REFERENCES media(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (media_id, genre_id)
);

-- Indexes for performance
CREATE INDEX idx_media_type ON media(type);
CREATE INDEX idx_media_title ON media USING gin(to_tsvector('english', title));
CREATE INDEX idx_media_tmdb ON media(tmdb_id);
CREATE INDEX idx_media_mal ON media(mal_id);
CREATE INDEX idx_media_release_date ON media(release_date);
