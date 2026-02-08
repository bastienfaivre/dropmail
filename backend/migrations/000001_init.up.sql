-- Create submissions table for email collection
CREATE TABLE IF NOT EXISTS submissions (
    email TEXT PRIMARY KEY NOT NULL,
    source TEXT NOT NULL DEFAULT 'unknown',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for querying by creation time (most recent first)
CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON submissions(created_at);

-- Index for querying by source (analytics)
CREATE INDEX IF NOT EXISTS idx_submissions_source ON submissions(source);
