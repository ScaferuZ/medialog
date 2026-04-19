-- Reviews table (separate from logs)
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    log_id UUID REFERENCES logs(id) ON DELETE SET NULL,

    title VARCHAR(255),
    content TEXT NOT NULL,
    rating DECIMAL(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 10),
    contains_spoilers BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Likes table (for logs and reviews)
CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_id UUID REFERENCES logs(id) ON DELETE CASCADE,
    review_id UUID REFERENCES reviews(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CHECK (
        (log_id IS NOT NULL AND review_id IS NULL) OR
        (log_id IS NULL AND review_id IS NOT NULL)
    ),
    UNIQUE(user_id, log_id),
    UNIQUE(user_id, review_id)
);

-- Indexes
CREATE INDEX idx_reviews_user ON reviews(user_id, created_at DESC);
CREATE INDEX idx_reviews_media ON reviews(media_id, created_at DESC);
CREATE INDEX idx_likes_user ON likes(user_id);
CREATE INDEX idx_likes_log ON likes(log_id);
CREATE INDEX idx_likes_review ON likes(review_id);
