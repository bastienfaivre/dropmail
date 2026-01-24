# Data Access & Export Guide

This guide documents how to access and export email submission data from the dropmail system.

## Database Location

The SQLite database is stored at `/data/dropmail.db` in the backend container's volume mount.

## Accessing the Database

### Option 1: Copy Database to Host

```bash
# Copy database file from the volume
docker cp dropmail-backend:/data/dropmail.db ./dropmail.db

# Open with sqlite3 on your host machine
sqlite3 ./dropmail.db
```

### Option 2: Run SQLite Container

```bash
# Run a temporary container with sqlite3
docker run --rm -it \
  -v dropmail_db-data:/data \
  keinos/sqlite3:latest \
  /data/dropmail.db
```

## Common SQL Queries

### Count Total Submissions

```sql
SELECT COUNT(*) FROM submissions;
```

### View All Submissions

```sql
SELECT email, source, created_at
FROM submissions
ORDER BY created_at DESC;
```

### Count Submissions by Source

```sql
SELECT source, COUNT(*) as count
FROM submissions
GROUP BY source
ORDER BY count DESC;
```

### Filter by Source

```sql
SELECT email, created_at
FROM submissions
WHERE source = 'portfolio-site'
ORDER BY created_at DESC;
```

### Recent Submissions (Last 7 Days)

```sql
SELECT email, source, created_at
FROM submissions
WHERE created_at >= datetime('now', '-7 days')
ORDER BY created_at DESC;
```

## Exporting to CSV

### Export All Data

```bash
# Using sqlite3 on host after copying database
sqlite3 -header -csv ./dropmail.db \
  "SELECT email, source, created_at FROM submissions ORDER BY created_at DESC;" \
  > submissions.csv
```

### Export Filtered by Source

```bash
sqlite3 -header -csv ./dropmail.db \
  "SELECT email, source, created_at FROM submissions WHERE source = 'portfolio-site';" \
  > portfolio-submissions.csv
```

### Export Using Temporary Container

```bash
docker run --rm \
  -v dropmail_db-data:/data \
  -v $(pwd):/output \
  keinos/sqlite3:latest \
  -header -csv /data/dropmail.db \
  "SELECT email, source, created_at FROM submissions;" \
  > /output/submissions.csv
```

## Database Schema

```sql
CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'direct',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Notes

- The database file is persisted in a Docker volume named `dropmail_db-data`
- All timestamps are stored in UTC
- The `source` field defaults to "direct" if not provided via URL parameter
