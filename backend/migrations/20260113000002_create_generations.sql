-- Migration: Create generations table
-- Requirements: 5.1, 5.2

CREATE TABLE IF NOT EXISTS generations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_idea TEXT NOT NULL,
    experience_level VARCHAR(20) NOT NULL,
    hook_preset VARCHAR(20) NOT NULL,
    files JSONB NOT NULL,
    category_id INTEGER REFERENCES categories(id) DEFAULT 5,
    avg_rating DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for filtering by category
CREATE INDEX IF NOT EXISTS idx_generations_category ON generations(category_id);

-- Index for sorting by newest (created_at descending)
CREATE INDEX IF NOT EXISTS idx_generations_created_at ON generations(created_at DESC);

-- Index for sorting by highest rated
CREATE INDEX IF NOT EXISTS idx_generations_avg_rating ON generations(avg_rating DESC);

-- Index for sorting by most viewed
CREATE INDEX IF NOT EXISTS idx_generations_view_count ON generations(view_count DESC);
