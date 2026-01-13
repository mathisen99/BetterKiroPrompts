-- Migration: Add review_stats column to scan_jobs
-- Stores AI review statistics as JSON

ALTER TABLE scan_jobs ADD COLUMN IF NOT EXISTS review_stats JSONB;
