CREATE TYPE activity_event_type AS ENUM (
    'item_added',
    'item_updated',
    'item_deleted',
    'status_changed',
    'rating_updated'
);

CREATE TABLE IF NOT EXISTS activity_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    media_item_id UUID REFERENCES media_items(id) ON DELETE SET NULL,
    event_type activity_event_type NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_activity_user_id ON activity_events (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_media_item_id ON activity_events (media_item_id);
