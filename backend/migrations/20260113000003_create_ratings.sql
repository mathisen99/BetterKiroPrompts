-- Migration: Create ratings table
-- Requirements: 7.2

CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    generation_id UUID NOT NULL REFERENCES generations(id) ON DELETE CASCADE,
    score SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5),
    voter_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(generation_id, voter_hash)
);

-- Index for efficient rating lookups by generation
CREATE INDEX IF NOT EXISTS idx_ratings_generation ON ratings(generation_id);
