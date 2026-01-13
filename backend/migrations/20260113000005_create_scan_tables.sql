-- Migration: Create scan tables for security scanning feature
-- Requirements: 11.1, 11.3

-- Scan jobs table to track scan requests and their status
CREATE TABLE IF NOT EXISTS scan_jobs (
    id VARCHAR(36) PRIMARY KEY,
    repo_url TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    languages TEXT[] DEFAULT '{}',
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '7 days'
);

-- Index for filtering by status (pending, cloning, scanning, reviewing, completed, failed)
CREATE INDEX IF NOT EXISTS idx_scan_jobs_status ON scan_jobs(status);

-- Index for cleanup job to find expired scans
CREATE INDEX IF NOT EXISTS idx_scan_jobs_expires_at ON scan_jobs(expires_at);

-- Index for sorting by newest
CREATE INDEX IF NOT EXISTS idx_scan_jobs_created_at ON scan_jobs(created_at DESC);

-- Scan findings table to store individual security findings
CREATE TABLE IF NOT EXISTS scan_findings (
    id VARCHAR(36) PRIMARY KEY,
    scan_job_id VARCHAR(36) NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    severity VARCHAR(10) NOT NULL,
    tool VARCHAR(50) NOT NULL,
    file_path TEXT NOT NULL,
    line_number INTEGER,
    description TEXT NOT NULL,
    remediation TEXT,
    code_example TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for looking up findings by scan job
CREATE INDEX IF NOT EXISTS idx_scan_findings_job_id ON scan_findings(scan_job_id);

-- Index for filtering findings by severity
CREATE INDEX IF NOT EXISTS idx_scan_findings_severity ON scan_findings(severity);
