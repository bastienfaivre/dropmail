-- Drop indexes first
DROP INDEX IF EXISTS idx_submissions_source;
DROP INDEX IF EXISTS idx_submissions_created_at;

-- Drop submissions table
DROP TABLE IF EXISTS submissions;
