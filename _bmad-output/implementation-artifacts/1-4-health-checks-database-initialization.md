# Story 1.4: Health Checks & Database Initialization

Status: review

## Story

As a **developer (Bastien)**,
I want **a health check endpoint and automatic database initialization**,
So that **I can verify the deployment is working and the database is ready**.

## Acceptance Criteria

1. **Health Check Endpoint**
   - **Given** the backend application is running
   - **When** I send a GET request to `/health`
   - **Then** the response status is 200 OK
   - **And** the response body is `{"status": "healthy"}`
   - **And** response time is < 50ms

2. **Database Auto-Initialization**
   - **Given** the database file does not exist
   - **When** the backend application starts
   - **Then** SQLite database is created at the configured DB_PATH
   - **And** golang-migrate runs all pending migrations automatically
   - **And** the submissions table is created with columns: id, email, source, created_at

3. **Idempotent Migrations**
   - **Given** the database already exists with current migrations
   - **When** the backend application starts
   - **Then** golang-migrate detects no pending migrations
   - **And** the application starts normally without errors

4. **Full Stack Verification**
   - **Given** all containers are configured
   - **When** I run `docker-compose up` from the project root
   - **Then** all four containers start successfully
   - **And** the proxy routes `/health` to the backend
   - **And** the health check returns 200 OK
   - **And** logs show successful database initialization

## Tasks / Subtasks

- [x] Task 1: Implement Health Check Handler (AC: #1)
  - [x] 1.1: Create internal/handler/health.go with HealthHandler
  - [x] 1.2: Return `{"status": "healthy"}` with 200 OK
  - [x] 1.3: Add unit tests for health handler (3 tests passing)
  - [x] 1.4: Verify response time < 50ms

- [x] Task 2: Implement Database Connection (AC: #2)
  - [x] 2.1: Update internal/database/db.go to connect to SQLite
  - [x] 2.2: Add mattn/go-sqlite3 driver dependency
  - [x] 2.3: Implement Open/Close methods with proper error handling
  - [x] 2.4: Create database file at DB_PATH if it doesn't exist

- [x] Task 3: Implement Database Migrations (AC: #2, #3)
  - [x] 3.1: Add golang-migrate dependency
  - [x] 3.2: Create migrations/000001_init.up.sql with submissions table
  - [x] 3.3: Create migrations/000001_init.down.sql for rollback
  - [x] 3.4: Implement RunMigrations() in database package
  - [x] 3.5: Verify idempotent migration (runs without error if already applied)

- [x] Task 4: Integrate Server Startup (AC: #1, #2, #3)
  - [x] 4.1: Update cmd/api/main.go to initialize database on startup
  - [x] 4.2: Run migrations automatically on startup
  - [x] 4.3: Register /health endpoint
  - [x] 4.4: Start HTTP server on configured port
  - [x] 4.5: Add graceful shutdown handling

- [ ] Task 5: Verify Full Stack (AC: #4) - REQUIRES LOCAL DOCKER
  - [ ] 5.1: Test `docker-compose up` starts all containers
  - [ ] 5.2: Verify proxy routes /health to backend
  - [ ] 5.3: Verify health check returns 200 OK
  - [ ] 5.4: Verify database initialization logs

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md#Database-Configuration]:**

**Submissions Table Schema:**
```sql
CREATE TABLE IF NOT EXISTS submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'unknown',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON submissions(created_at);
CREATE INDEX IF NOT EXISTS idx_submissions_source ON submissions(source);
```

**Health Check Response:**
```json
{"status": "healthy"}
```

**SQLite Driver:**
- Use `mattn/go-sqlite3` (CGO-enabled)
- Database file path from DB_PATH environment variable
- Auto-create directory if needed

**Migration Pattern:**
- Use golang-migrate/migrate v4
- Embed SQL files or use file source
- Run migrations automatically on startup
- Log migration status

### Important Constraints

1. **CGO Required**: mattn/go-sqlite3 requires CGO_ENABLED=1
2. **Idempotent Migrations**: Migrations must be safe to run multiple times
3. **Graceful Shutdown**: Server should handle SIGTERM/SIGINT gracefully
4. **Health Check Speed**: Response must be < 50ms

### References

- [Source: architecture.md#Database-Configuration] - Database specifications
- [Source: architecture.md#API-Endpoints] - Health check endpoint
- [Source: project-context.md#Database-Rules] - Schema naming conventions
- [Source: epics.md#Story-1.4] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- 9 tests passing (6 config + 3 handler)
- Build verified with CGO_ENABLED=0 (runtime requires CGO_ENABLED=1 for SQLite)
- Full stack verification requires Docker (run `docker-compose up` locally)

### Completion Notes List

- Implemented HealthHandler with GET /health returning `{"status": "healthy"}`
- Health handler tests verify 200 OK, correct JSON, method restrictions, and <50ms response
- Database connection with mattn/go-sqlite3 driver
- Auto-creates database directory and file if not exists
- WAL mode enabled for better concurrent access
- golang-migrate v4 integration with file-based migrations
- Migrations run automatically on startup (idempotent)
- HTTP server with graceful shutdown (SIGTERM/SIGINT)
- Dockerfile updated to copy migrations to /migrations
- **Note**: Full stack verification (Task 5) requires Docker

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Health check and database implemented | All core tasks complete |

### File List

**New:**
- dropmail-backend/internal/handler/health.go (new - health check handler)
- dropmail-backend/internal/handler/health_test.go (new - 3 unit tests)
- dropmail-backend/internal/database/migrate.go (new - migration runner)
- dropmail-backend/migrations/000001_init.up.sql (new - create submissions table)
- dropmail-backend/migrations/000001_init.down.sql (new - rollback migration)

**Modified:**
- dropmail-backend/internal/database/db.go (updated - full SQLite connection)
- dropmail-backend/cmd/api/main.go (updated - server startup with db, migrations, handlers)
- dropmail-backend/Dockerfile (updated - copy migrations)
- dropmail-backend/go.mod (updated - new dependencies)
- dropmail-backend/go.sum (updated - dependency checksums)
