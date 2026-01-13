-- Migration: Create views table for IP-based view deduplication
-- Requirements: 5.1, 5.3

CREATE TABLE IF NOT EXISTS views (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    generation_id UUID NOT NULL REFERENCES generations(id) ON DELETE CASCADE,
    ip_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(generation_id, ip_hash)
);

-- Index for efficient view lookups by generation
CREATE INDEX IF NOT EXISTS idx_views_generation_id ON views(generation_id);
