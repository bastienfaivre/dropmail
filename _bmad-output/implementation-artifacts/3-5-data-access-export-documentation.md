# Story 3.5: Data Access & Export Documentation

Status: review

## Story

As a **developer (Bastien)**,
I want **clear documentation for accessing and exporting email data**,
So that **I can query submissions and export data when needed**.

## Acceptance Criteria

1. **Database Access**
   - **Given** the SQLite database is running in the db container
   - **When** I connect via `docker exec -it dropmail-db sqlite3 /data/dropmail.db`
   - **Then** I can access the database directly
   - **And** I can run SQL queries interactively

2. **Count Submissions**
   - **Given** I want to count total submissions
   - **When** I run `SELECT COUNT(*) FROM submissions;`
   - **Then** I get the total number of email submissions

3. **Query by Source**
   - **Given** I want to see submissions by source
   - **When** I run `SELECT source, COUNT(*) FROM submissions GROUP BY source;`
   - **Then** I see a breakdown of submissions per source

4. **View All Submissions**
   - **Given** I want to view all submission data
   - **When** I run `SELECT email, source, created_at FROM submissions ORDER BY created_at DESC;`
   - **Then** I see all submissions with timestamps

5. **CSV Export**
   - **Given** I want to export to CSV
   - **When** I run the sqlite3 CSV export command
   - **Then** a CSV file is created with all submission data
   - **And** I can copy it from the container or access via the mounted volume

6. **Documentation**
   - **Given** documentation needs to be provided
   - **When** the project is complete
   - **Then** a README section documents these access patterns
   - **And** example queries are provided for common operations

## Tasks / Subtasks

- [x] Task 1: Verify Database Access (AC: #1)
  - [x] 1.1: Confirmed db container uses Docker volume
  - [x] 1.2: Documented docker cp and sqlite container methods

- [x] Task 2: Document SQL Queries (AC: #2, #3, #4)
  - [x] 2.1: Count total submissions query
  - [x] 2.2: Group by source query
  - [x] 2.3: View all submissions query
  - [x] 2.4: Filter by source query
  - [x] 2.5: Recent submissions query (last 7 days)

- [x] Task 3: Document CSV Export (AC: #5)
  - [x] 3.1: Document sqlite3 CSV export command
  - [x] 3.2: Document using temporary container
  - [x] 3.3: Document export with filters

- [x] Task 4: Create Documentation (AC: #6)
  - [x] 4.1: Created DATA_ACCESS.md with all examples
  - [x] 4.2: Included database schema documentation

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**Database Access:**
- SQLite database at /data/dropmail.db in backend volume
- Direct access via docker cp or sqlite container
- Mounted volume allows data persistence

**Data Schema:**
```sql
CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'direct',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Important Constraints

1. **No Admin UI**: Access via SQLite CLI only
2. **Volume Mount**: Database persisted at /data/dropmail.db
3. **Export Options**: CSV via sqlite3 commands

### References

- [Source: architecture.md#Database] - Database schema
- [Source: epics.md#Story-3.5] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- DATA_ACCESS.md created with comprehensive examples
- Database schema documented

### Completion Notes List

- Created DATA_ACCESS.md with:
  - Database location documentation
  - Two access methods (docker cp, sqlite container)
  - Common SQL queries:
    - Count total submissions
    - View all submissions
    - Count by source
    - Filter by source
    - Recent submissions (last 7 days)
  - CSV export methods:
    - Direct export to file
    - Filtered export
    - Using temporary container
  - Database schema documentation
- Note: The db container (busybox) doesn't have sqlite3 installed
- Recommended approach: Copy database or use sqlite container

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Documentation created | DATA_ACCESS.md with all examples |

### File List

**New:**
- DATA_ACCESS.md (data access and export documentation)

**Modified:**
(none)
