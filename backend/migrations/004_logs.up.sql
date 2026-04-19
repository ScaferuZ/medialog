-- Activity logs table (core feature - tracks media consumption)
CREATE TABLE logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    
    -- Log status
    status VARCHAR(20) NOT NULL CHECK (status IN ('planned', 'in_progress', 'completed', 'dropped')),
    
    -- Rating (0-10 with 0.5 increments, stored as DECIMAL(3,1))
    rating DECIMAL(3,1) CHECK (rating >= 0 AND rating <= 10),
    
    -- Dates
    started_at DATE,
    completed_at DATE,
    
    -- Rewatch/reread/replay count
    rewatch_count INTEGER DEFAULT 0,
    
    -- Progress tracking
    progress INTEGER, -- Pages read, episodes watched, hours played
    total INTEGER,    -- Total pages/episodes/hours
    
    -- User content
    note TEXT,
    contains_spoilers BOOLEAN DEFAULT false,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- One log per user per media
    UNIQUE(user_id, media_id)
);

-- Indexes for common queries
CREATE INDEX idx_logs_user_created ON logs(user_id, created_at DESC);
CREATE INDEX idx_logs_media ON logs(media_id);
CREATE INDEX idx_logs_status ON logs(status);
CREATE INDEX idx_logs_user_status ON logs(user_id, status);
