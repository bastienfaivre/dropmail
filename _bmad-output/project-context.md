---
project_name: 'dropmail'
user_name: 'Bastien'
date: '2026-01-23'
sections_completed: []
source_document: 'architecture.md'
sections_completed: ['technology_stack', 'language_rules', 'middleware_architecture', 'testing_antipatterns']
status: 'complete'
rule_count: 35
optimized_for_llm: true
---

# Project Context for AI Agents

_This file contains critical rules and patterns that AI agents must follow when implementing code for dropmail. Focus on unobvious details that agents might otherwise miss._

---

## Technology Stack & Versions

### Backend (Go)
- **Go**: 1.25.6 (required - do not use older versions)
- **SQLite Driver**: mattn/go-sqlite3 (CGO-enabled, latest)
- **Migrations**: golang-migrate/migrate v4.19.1
- **CORS**: rs/cors (latest)
- **Rate Limiting**: golang.org/x/time/rate (latest)
- **Logging**: slog (Go 1.21+ standard library)

### Frontend (TypeScript)
- **Node.js**: 24.13.0 LTS (required for development)
- **TypeScript**: 5.x (latest)
- **Vite**: 7.3.1 (build tool)
- **Template**: vanilla-ts (no framework)

### Infrastructure
- **Database**: SQLite (file-based, single file)
- **Containers**: Docker Compose v3.8 (4-container architecture)
- **Proxy**: nginx alpine (HTTPS termination)
- **Base Images**: distroless/static (Go), nginx:alpine (frontend/proxy)

### External Dependencies
- **Cloudflare Turnstile**: Required for spam protection (CDN loaded)

## Critical Implementation Rules

### Go Backend Rules

**Project Structure (MANDATORY):**
- Use standard Go project layout: `cmd/api/`, `internal/`, `migrations/`
- All application code in `internal/` (non-importable from outside)
- Entry point at `cmd/api/main.go`
- Co-locate tests with implementation: `handler.go` → `handler_test.go`

**Naming Conventions:**
- Exported types/functions: `PascalCase` (e.g., `SubmitRequest`, `HandleSubmit`)
- Unexported: `camelCase` (e.g., `rateLimiter`, `validateEmail`)
- Package names: lowercase, single word (e.g., `handler`, `middleware`)
- Files: `snake_case.go` (e.g., `submit_handler.go`, `rate_limiter.go`)

**Logging (CRITICAL):**
- Use `slog` ONLY - no `fmt.Println` or `log.Println` in production code
- Log successful submissions at `info` level with: email, source, IP, timestamp
- Log validation failures at `warn` level with: reason, IP
- NEVER log Turnstile tokens (security sensitive)

**Error Handling:**
- Always return simple JSON: `{"error": "message"}` - no error codes or nested structures
- Use `net/mail.ParseAddress()` for email validation (standard library)
- SQL queries use parameter binding only - never string concatenation

### TypeScript Frontend Rules

**File Naming:**
- Use `kebab-case.ts` for files (e.g., `form-handler.ts`, `api-client.ts`)
- Use `PascalCase` for interfaces/types (e.g., `SubmitRequest`, `ApiResponse`)
- Use `camelCase` for functions/variables (e.g., `handleSubmit`, `turnstileToken`)

**Validation Flow (CRITICAL ORDER):**
1. HTML5 validation first (`<input type="email" required>`)
2. TypeScript regex validation before Turnstile
3. Render Turnstile widget ONLY after email validation passes
4. Submit to API with email + source + token

**API Communication:**
- All JSON fields use `snake_case` (matches backend/database)
- Request: `{ "email": "", "source": "", "turnstile_token": "" }`
- Success: `{ "success": true, "message": "" }`
- Error: `{ "error": "" }`

### Middleware Order (CRITICAL - Never Reorder)

Go middleware MUST be applied in this exact order:
1. **CORS** - First, handles preflight requests
2. **Rate Limiting** - Second, rejects excessive requests early
3. **Logging** - Third, logs all requests that pass rate limiting
4. **Handler** - Last, executes business logic

```go
// Correct wrapping (reverse order in code)
var h http.Handler = mux
h = middleware.CORS(h)      // Applied first
h = middleware.RateLimit(h) // Applied second
h = middleware.Logging(h)   // Applied third
```

### Database Rules

**Schema Naming:**
- Table names: `snake_case`, lowercase, plural (e.g., `submissions`)
- Column names: `snake_case` (e.g., `created_at`, not `createdAt`)
- Primary keys: `id` (simple, no table prefix)
- Indexes: `idx_{table}_{column}` (e.g., `idx_submissions_created_at`)

**Migrations:**
- Use golang-migrate with embedded SQL files
- File pattern: `000001_init.up.sql`, `000001_init.down.sql`
- Run automatically on startup (idempotent)
- SQLite: No explicit BEGIN/COMMIT needed in migration files

### Container Architecture (4 Containers)

| Container | Role | Network Access |
|-----------|------|----------------|
| `db` | SQLite volume mount | Internal only |
| `backend` | Go API on :8080 | Internal only |
| `frontend` | nginx serving static | Internal only |
| `proxy` | nginx HTTPS termination | **Only exposed container** |

**Communication:**
- `proxy` → `frontend`: `http://dropmail-frontend:80`
- `proxy` → `backend`: `http://dropmail-backend:8080`
- `backend` → `db`: Shared volume at `/data/dropmail.db`

### Testing Rules

**Go Tests:**
- Co-locate tests: `submit_handler.go` → `submit_handler_test.go` (same directory)
- Use Go standard `testing` package and `httptest` for handlers
- Table-driven tests for validation logic
- Integration tests with test SQLite database (separate file)

**Integration Tests:**
- Build all 4 containers separately
- Spin up complete stack with `docker-compose.test.yml`
- Test full flow: proxy → frontend → backend → db
- Tear down and cleanup volumes after tests

### Critical Anti-Patterns (NEVER DO)

**Go Backend:**
- NEVER use `fmt.Println` for logging (use `slog`)
- NEVER concatenate SQL strings (use parameter binding)
- NEVER log Turnstile tokens
- NEVER reorder middleware (CORS must be first)
- NEVER use complex error response formats

**TypeScript Frontend:**
- NEVER render Turnstile before email validation
- NEVER use `camelCase` in JSON fields (use `snake_case`)
- NEVER trust client-side validation alone (backend validates too)

**Infrastructure:**
- NEVER expose containers other than `proxy` to host network
- NEVER hardcode secrets (use environment variables)
- NEVER skip HTTPS in production

### Performance Constraints

- Page weight: < 50KB total
- TTI: < 1 second
- API response: < 200ms (95th percentile)
- Form submission: < 2 seconds complete

---

## Usage Guidelines

**For AI Agents:**
- Read this file before implementing any code
- Follow ALL rules exactly as documented
- When in doubt, prefer the more restrictive option
- Refer to architecture.md for detailed examples and rationale

**For Humans:**
- Keep this file lean and focused on agent needs
- Update when technology stack changes
- Review periodically for outdated rules

---

_Generated from architecture.md on 2026-01-23_

