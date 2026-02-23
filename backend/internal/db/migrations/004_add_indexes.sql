-- Additional composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_media_user_type ON media_items (user_id, media_type);
CREATE INDEX IF NOT EXISTS idx_media_user_status ON media_items (user_id, status);
CREATE INDEX IF NOT EXISTS idx_media_user_created ON media_items (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_media_release_year ON media_items (release_year) WHERE release_year IS NOT NULL;
