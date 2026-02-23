CREATE TYPE media_type AS ENUM ('movie', 'music', 'game');
CREATE TYPE media_status AS ENUM ('owned', 'wishlist', 'currently_using', 'completed');

CREATE TABLE IF NOT EXISTS media_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    media_type media_type NOT NULL,
    status media_status NOT NULL DEFAULT 'owned',
    creator TEXT NOT NULL DEFAULT '',
    genre TEXT[] NOT NULL DEFAULT '{}',
    release_year INT,
    cover_url TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    rating NUMERIC(3,1) CHECK (rating >= 0 AND rating <= 10),
    tmdb_id TEXT,
    musicbrainz_id TEXT,
    igdb_id TEXT,
    metadata JSONB NOT NULL DEFAULT '{}',
    search_vector TSVECTOR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- GIN indexes for full-text and trigram search
CREATE INDEX IF NOT EXISTS idx_media_search ON media_items USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS idx_media_title_trgm ON media_items USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_media_creator_trgm ON media_items USING GIN (creator gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_media_genre ON media_items USING GIN (genre);
CREATE INDEX IF NOT EXISTS idx_media_user_id ON media_items (user_id);
CREATE INDEX IF NOT EXISTS idx_media_type ON media_items (media_type);
CREATE INDEX IF NOT EXISTS idx_media_status ON media_items (status);

-- Auto-update trigger for search_vector
-- Weighted: title A, creator B, genre C, notes D
CREATE OR REPLACE FUNCTION media_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.creator, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(array_to_string(NEW.genre, ' '), '')), 'C') ||
        setweight(to_tsvector('english', coalesce(NEW.notes, '')), 'D');
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER media_search_vector_trigger
    BEFORE INSERT OR UPDATE ON media_items
    FOR EACH ROW EXECUTE FUNCTION media_search_vector_update();

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER media_items_updated_at
    BEFORE UPDATE ON media_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
