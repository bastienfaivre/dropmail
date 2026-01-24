---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
status: 'complete'
completedAt: '2026-01-23'
inputDocuments:
  - "/app/_bmad-output/planning-artifacts/product-brief-dropmail-2026-01-22.md"
  - "/app/_bmad-output/planning-artifacts/prd.md"
workflowType: 'architecture'
project_name: 'app'
user_name: 'Bastien'
date: '2026-01-22'
---

# Architecture Decision Document - dropmail

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**

dropmail has 49 functional requirements organized into 9 categories, reflecting a focused single-purpose tool designed for simplicity:

1. **Form Submission & Collection (FR1-FR8)**: Core email collection with client-side and server-side validation, success/error messaging, timestamp tracking, and explicitly NO verification emails
2. **Spam Protection & Security (FR9-FR13)**: Cloudflare Turnstile verification (required), IP-based rate limiting, HTTPS enforcement, input sanitization, CORS restrictions
3. **Source Tracking & Analytics (FR14-FR17)**: URL parameter-based source identification for multi-property analytics
4. **Data Access & Export (FR18-FR22)**: Direct SQLite database access via CLI tools, no admin UI required
5. **Deployment & Configuration (FR23-FR28)**: Single-command Docker Compose deployment, automatic database initialization, environment-based configuration with fail-fast validation
6. **Multi-Property Embedding (FR29-FR33)**: iframe embedding across Notion, portfolio sites, and future properties with responsive container adaptation
7. **Responsive Design & Device Support (FR34-FR38)**: Mobile-first design supporting 320px-1024px+ screens with touch-friendly controls
8. **Accessibility & Keyboard Navigation (FR39-FR45)**: Full WCAG 2.1 AA compliance with semantic HTML, screen reader support, keyboard navigation
9. **Performance & Reliability (FR46-FR49)**: < 1s TTI, < 2s submission, 100% uptime, zero maintenance goal

**Non-Functional Requirements:**

Critical NFRs that will drive architectural decisions:

- **Performance**: TTI < 1s, page weight < 50KB, API response < 200ms (95th percentile), submission complete < 2s
- **Security**: HTTPS enforcement, input sanitization, Cloudflare Turnstile required, CORS whitelist, rate limiting, environment-based secrets
- **Reliability**: 100% uptime target, zero maintenance operation, graceful degradation if Turnstile unavailable (fail closed)
- **Accessibility**: WCAG 2.1 AA compliance, keyboard navigation, screen reader support, 4.5:1 color contrast
- **Resource Efficiency**: Total footprint < 150MB RAM, frontend < 10MB, backend < 100MB, database < 1MB for thousands of emails
- **Privacy**: Minimal data collection (email, source, timestamp only), no tracking, no third-party data sharing, complete data sovereignty

**Scale & Complexity:**

- **Primary domain**: Web Application (iframe-embeddable form + REST API backend)
- **Complexity level**: Low (single-purpose utility tool)
- **Estimated architectural components**: 3 core components (frontend static app, backend REST API, SQLite database)
- **Browser support**: Modern browsers only (Chrome, Firefox, Safari, Edge - last 2 versions)
- **Device support**: Full responsive design (mobile, tablet, desktop)

### Technical Constraints & Dependencies

**Required External Dependencies:**
- **Cloudflare Turnstile**: Non-negotiable spam protection requirement (free tier, mature service)
- Graceful degradation strategy: Fail closed if Turnstile unavailable

**Technology Constraints:**
- **SQLite**: Mandated database choice for simplicity (file-based, zero configuration, direct CLI access)
- **Docker Compose**: Required deployment mechanism for single-command setup
- **Vanilla JavaScript**: No framework dependencies - lightweight, minimal build tooling
- **Self-hosted**: Complete stack must run on owned infrastructure

**Design Philosophy Constraints:**
- Purpose-built minimalism: Build exactly what's needed, nothing more
- Portfolio-grade code quality: Clean architecture demonstrating technical capability
- Zero vendor lock-in: No SaaS dependencies beyond Turnstile
- Anti-feature-creep: MVP is the final product

### Cross-Cutting Concerns Identified

**Security (affects all components):**
- HTTPS enforcement across frontend and backend
- Input sanitization at API boundary
- CORS policy management for iframe embedding
- Rate limiting to prevent abuse
- Secret management via environment variables

**Accessibility (affects frontend):**
- WCAG 2.1 AA compliance throughout
- Semantic HTML structure
- Keyboard navigation support
- Screen reader compatibility
- Touch-friendly mobile interactions

**Performance (affects all components):**
- Aggressive optimization for < 50KB page weight
- API response time optimization
- Database write performance
- Minimal resource footprint

**Monitoring & Observability:**
- Health check endpoints for container monitoring
- Startup validation with fail-fast behavior
- Clear error messages for debugging

**Deployment & DevOps:**
- Docker containerization for all components
- Environment-based configuration
- Automatic database initialization
- Single-command deployment experience

## Starter Template Evaluation

### Primary Technology Domain

**Web Application** (iframe-embeddable form + REST API backend) based on project requirements analysis

### Technical Preferences Established

**Language & Runtime Versions:**
- **Frontend**: TypeScript (latest), Node.js 24.13.0 LTS
- **Backend**: Go 1.25.6 (latest stable)
- **Build Tool**: Vite 7.3.1 (latest stable)
- **Styling**: Modern CSS with custom properties
- **Reverse Proxy**: nginx

**Infrastructure Decisions:**
- **Database**: SQLite (PRD-mandated)
- **Deployment**: Docker Compose with 4-container architecture
- **Container Strategy**: Multi-stage builds for minimal images

### Container Architecture

**4-Container Setup:**

1. **db** - SQLite database container
   - Persists data via Docker volume mount
   - Accessible via container name on internal network

2. **backend** - Go API server
   - Distroless/scratch-based image (< 10MB)
   - Connects to `db` container
   - Exposes API endpoints

3. **frontend** - Static TypeScript form
   - nginx alpine serving built assets (< 10MB)
   - Connects to `backend` via internal network

4. **proxy** - nginx reverse proxy
   - HTTPS termination
   - Routes requests to frontend/backend
   - Only container exposed to host network

**Network Architecture:**
- All containers on shared internal Docker network
- Communication via container names (e.g., `http://backend:8080`)
- Only `proxy` container exposes ports to host
- DB persistence via named volume mount

### Development Environment

**.devcontainer Configuration:**

Ubuntu-based development container with:
- Go 1.25.6
- Node.js 24.13.0 LTS
- Docker and Docker Compose
- Essential development tools (git, curl, etc.)
- Go development tools (gopls, delve)
- Node development tools (npm, typescript)

**Purpose**: Consistent development environment for testing all components locally

### Testing Strategy

**Integration Tests:**
- Build each Docker container separately
- Spin up complete 4-container stack
- Test full request flow: proxy → frontend → backend → db
- Verify container networking and communication
- Test HTTPS termination and CORS
- Validate Turnstile integration
- Tear down stack after tests

**Unit Tests:**
- Frontend: TypeScript validation, form logic tests
- Backend: Handler tests, database tests, middleware tests

### Starter Options Considered

**Frontend Options:**
1. **Vite vanilla-ts** - Official Vite template for TypeScript static pages ✓
2. **VITAM** - Front-end template with Vite for static websites
3. **Custom from scratch** - Manual TypeScript setup

**Backend Options:**
1. **GRAB** - Go REST API Boilerplate (2026, AI-focused)
2. **qiangxue/go-rest-api** - Clean Architecture approach
3. **benc-uk/go-rest-api** - Lightweight net/http microservice
4. **Standard Go project layout** - Clean structure from scratch ✓

**Decision Rationale:**

Given dropmail's intentional simplicity (single POST endpoint, no auth, SQLite only), traditional REST API boilerplates introduce unnecessary complexity (PostgreSQL setup, JWT authentication, full CRUD patterns, complex middleware chains). The portfolio-grade quality goal is better served by clean, purposeful code rather than adapting heavy scaffolding.

### Selected Approach

**Frontend: Vite vanilla-ts (minimal starter)**

**Initialization Command:**
```bash
npm create vite@latest dropmail-frontend -- --template vanilla-ts
```

**Backend: Standard Go project layout (from scratch)**

**Initialization Command:**
```bash
mkdir dropmail-backend && cd dropmail-backend
go mod init github.com/bastien/dropmail-backend
```

### Architectural Decisions Provided by Starters

**Frontend (Vite vanilla-ts):**

**Language & Runtime:**
- TypeScript with strict type checking enabled
- Target: ES2020 for modern browser support
- Node.js 24.13.0 LTS required for development
- Vite 7.3.1 for build tooling

**Build Tooling:**
- Vite 7.x for development server and production builds
- esbuild for fast transpilation
- Rollup for optimized production bundles
- Hot Module Replacement (HMR) for development

**Styling Solution:**
- No pre-configured CSS framework (as requested)
- Modern CSS with custom properties support built-in
- CSS modules available if needed
- PostCSS integration available

**Project Structure:**
```
dropmail-frontend/
├── src/
│   ├── main.ts          # Entry point
│   ├── style.css        # Global styles
│   └── vite-env.d.ts    # Vite type definitions
├── index.html           # HTML entry
├── tsconfig.json        # TypeScript configuration
├── vite.config.ts       # Vite configuration
├── Dockerfile           # nginx alpine multi-stage build
└── package.json
```

**Development Experience:**
- Development server with instant HMR
- TypeScript error checking in IDE and build
- Production build with minification and tree-shaking
- Asset optimization (< 50KB target easily achievable)

**Backend (Standard Go Layout):**

**Language & Runtime:**
- Go 1.25.6 (latest stable, January 2026)
- Standard library net/http package
- Minimal external dependencies (SQLite driver, rate limiting)

**Project Structure (Standard Go Layout):**
```
dropmail-backend/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── handlers/             # HTTP request handlers
│   ├── middleware/           # CORS, rate limiting, logging
│   ├── models/               # Data models
│   ├── database/             # SQLite connection & queries
│   └── config/               # Environment configuration
├── migrations/               # Database schema migrations
├── Dockerfile                # Multi-stage build (distroless/scratch)
├── go.mod
└── go.sum
```

**Code Organization:**
- Clean separation of concerns (handlers, middleware, database)
- Internal package for application code (non-importable)
- cmd package for executables
- Standard Go project layout conventions

**Development Experience:**
- Go standard tooling (go build, go test, go fmt)
- Air for hot-reload during development (optional)
- Delve for debugging
- Native Go testing framework

### Docker Compose Architecture

**4-Container Setup:**

```yaml
version: '3.8'

services:
  db:
    image: alpine:latest
    container_name: dropmail-db
    volumes:
      - db-data:/data
    networks:
      - dropmail-internal
    # SQLite file stored at /data/dropmail.db

  backend:
    build:
      context: ./dropmail-backend
      dockerfile: Dockerfile
    container_name: dropmail-backend
    depends_on:
      - db
    environment:
      - DB_PATH=/data/dropmail.db
      - TURNSTILE_SECRET=${TURNSTILE_SECRET}
      - CORS_ORIGINS=${CORS_ORIGINS}
    volumes:
      - db-data:/data
    networks:
      - dropmail-internal
    # Communicates with db via container name

  frontend:
    build:
      context: ./dropmail-frontend
      dockerfile: Dockerfile
    container_name: dropmail-frontend
    networks:
      - dropmail-internal
    # Serves static assets via nginx

  proxy:
    image: nginx:alpine
    container_name: dropmail-proxy
    depends_on:
      - frontend
      - backend
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    networks:
      - dropmail-internal
    # Routes to frontend and backend containers

networks:
  dropmail-internal:
    driver: bridge

volumes:
  db-data:
    driver: local
```

**Network Communication:**
- `proxy` → `frontend` via `http://dropmail-frontend:80`
- `proxy` → `backend` via `http://dropmail-backend:8080`
- `backend` → `db` via shared volume at `/data/dropmail.db`

### Deployment Strategy

**Multi-stage Docker builds:**

**Frontend Dockerfile:**
```dockerfile
# Build stage
FROM node:24.13.0-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
```

**Backend Dockerfile:**
```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o api cmd/api/main.go

# Production stage
FROM gcr.io/distroless/static
COPY --from=builder /app/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
```

**Expected Image Sizes:**
- `db`: < 5MB (alpine)
- `backend`: < 15MB (distroless with Go binary + SQLite)
- `frontend`: < 10MB (nginx alpine + static assets)
- `proxy`: < 10MB (nginx alpine)
- **Total**: < 40MB for entire stack

### Development Container Configuration

**.devcontainer/devcontainer.json:**
```json
{
  "name": "dropmail Development",
  "image": "ubuntu:24.04",
  "features": {
    "ghcr.io/devcontainers/features/go:1": {
      "version": "1.25"
    },
    "ghcr.io/devcontainers/features/node:1": {
      "version": "24"
    },
    "ghcr.io/devcontainers/features/docker-in-docker:2": {}
  },
  "customizations": {
    "vscode": {
      "extensions": [
        "golang.go",
        "dbaeumer.vscode-eslint",
        "esbenp.prettier-vscode"
      ]
    }
  },
  "postCreateCommand": "go version && node --version && docker --version"
}
```

### Build and Development Workflow

**Frontend Development:**
```bash
cd dropmail-frontend
npm install
npm run dev          # Start dev server (http://localhost:5173)
npm run build        # Production build to dist/
npm run preview      # Preview production build
```

**Backend Development:**
```bash
cd dropmail-backend
go mod download
go run cmd/api/main.go              # Start API server
go build -o bin/api cmd/api/main.go # Build binary
go test ./...                        # Run tests
```

**Docker Development:**
```bash
# Build individual containers
docker build -t dropmail-frontend ./dropmail-frontend
docker build -t dropmail-backend ./dropmail-backend

# Start complete stack
docker-compose up --build

# Run integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### Testing Infrastructure

**Unit Tests:**

**Frontend:**
- Vitest for unit testing (can be added post-initialization)
- Playwright for E2E testing (separate installation)
- TypeScript provides compile-time safety

**Backend:**
- Go standard testing package (testing, httptest)
- Table-driven tests for handlers
- Mock database for unit tests
- Integration tests with test SQLite database

**Integration Tests:**

**Test Strategy:**
1. Build all 4 containers separately
2. Spin up complete stack with `docker-compose.test.yml`
3. Run test suite against running stack:
   - Health check all containers
   - Test proxy → frontend routing
   - Test proxy → backend routing
   - Test backend → db communication
   - Test form submission end-to-end
   - Test CORS configuration
   - Test Turnstile integration (with test keys)
   - Test rate limiting
4. Capture logs from all containers
5. Tear down stack and cleanup volumes

**Test Docker Compose:**
```yaml
# docker-compose.test.yml
version: '3.8'

services:
  # Same 4 services as production
  # Plus test-runner service
  test-runner:
    build:
      context: ./tests
      dockerfile: Dockerfile
    depends_on:
      - proxy
      - frontend
      - backend
      - db
    networks:
      - dropmail-internal
    environment:
      - PROXY_URL=http://dropmail-proxy
    command: ["npm", "run", "test:integration"]
```

### Configuration Management

**Frontend:**
- Environment variables via Vite (.env files)
- Build-time variable injection for API URLs
- Runtime configuration via window object if needed

**Backend:**
- Environment variables for all configuration
- Graceful degradation with sensible defaults
- Fail-fast validation on startup

**Docker Compose Environment:**
```env
# .env file
TURNSTILE_SECRET=your_secret_key
CORS_ORIGINS=https://example.com,https://notion.so
DB_PATH=/data/dropmail.db
API_PORT=8080
```

### Initial Setup Notes

**First Implementation Story:**
1. Initialize frontend with Vite vanilla-ts template
2. Initialize Go module and create standard project structure
3. Create 4 Dockerfiles (db, frontend, backend, proxy)
4. Configure docker-compose.yml with internal network
5. Set up .devcontainer for development environment
6. Configure nginx reverse proxy routing
7. Verify hot-reload works for frontend and backend
8. Set up integration test infrastructure

**Version Tracking:**
- Go: 1.25.6
- Node.js: 24.13.0 LTS
- Vite: 7.3.1
- TypeScript: latest (5.x)

**No starter-imposed constraints** - Clean foundation for implementing exactly what dropmail needs, nothing more.

## Core Architectural Decisions

### Decision Summary

All critical architectural decisions have been made collaboratively, focusing on simplicity, maintainability, and portfolio-grade quality.

### Data Architecture

**SQLite Driver:**
- **Decision**: mattn/go-sqlite3
- **Version**: Latest (CGO-enabled)
- **Rationale**: Battle-tested, full SQLite feature support, works well with multi-stage Docker builds
- **Affects**: Backend database layer, Dockerfile build configuration

**Database Migrations:**
- **Decision**: golang-migrate/migrate v4
- **Version**: v4.19.1 (November 2025)
- **Migration Files**: `migrations/000001_init.up.sql` and `000001_init.down.sql`
- **Execution Strategy**: Embedded in Go binary via `embed` package, automatic execution on startup
- **SQLite-Specific**: Implicit transaction wrapping (no explicit BEGIN/COMMIT needed in migration files)
- **Rationale**: Industry standard for Go migrations, versioned schema changes, rollback support
- **Affects**: Backend initialization, database schema management

**Schema Initialization:**
- **Approach**: Automatic on first run
- **Implementation**: Migrations run on application startup, idempotent execution
- **Rationale**: Aligns with "zero maintenance" and "automatic database initialization" requirements

### Security & Middleware

**Rate Limiting:**
- **Decision**: golang.org/x/time/rate (standard library extension)
- **Algorithm**: Token bucket
- **Storage**: In-memory (resets on restart)
- **Rationale**: Lightweight, no external dependencies, sufficient for single-instance deployment
- **Configuration**: Per-IP limits via environment variables
- **Affects**: Backend middleware layer

**CORS Middleware:**
- **Decision**: rs/cors
- **Configuration**: Environment-based origin whitelist
- **Rationale**: Clean origin validation, lightweight, well-maintained
- **Environment Variable**: CORS_ORIGINS (comma-separated list)
- **Affects**: Backend middleware, iframe embedding capability

**Input Validation & Sanitization:**
- **Email Validation**: net/mail.ParseAddress (standard library)
- **Sanitization**: Manual SQL parameter binding (no raw SQL concatenation)
- **Rationale**: Single email field doesn't warrant heavy validation framework, standard library sufficient
- **Affects**: Backend handler layer, database queries

**Logging:**
- **Decision**: slog (Go 1.21+ standard library)
- **Format**: Structured logging (JSON in production)
- **Rationale**: Official Go standard for structured logging, no external dependency
- **Affects**: All backend components (handlers, middleware, database)

### API & Communication Patterns

**API Design:**
- **Pattern**: REST
- **Endpoint**: Single POST /api/submit
- **Request Format**: JSON `{"email": "user@example.com", "source": "portfolio-site", "turnstile_token": "..."}`
- **Response Format**: JSON
  - Success: `{"success": true, "message": "Email submitted successfully"}`
  - Error: `{"error": "Invalid email format"}` (simple JSON, no error codes)
- **Rationale**: Simplicity over complexity, single endpoint doesn't need complex API patterns
- **Affects**: Frontend/backend contract, API tests

**Error Handling:**
- **Format**: Simple JSON error messages
- **Example**: `{"error": "Invalid email format"}`
- **HTTP Status Codes**: Standard (200 success, 400 bad request, 429 rate limit, 500 server error)
- **Rationale**: Lightweight, easy to parse, sufficient for form feedback
- **Affects**: Frontend error display, backend error responses

**Request Validation:**
- **Server-side**: Required (email format, Turnstile token, rate limits)
- **Client-side**: Optional (HTML5 validation, TypeScript checks for UX)
- **Rationale**: Defense in depth, never trust client validation alone
- **Affects**: Both frontend and backend

### Frontend Architecture

**Cloudflare Turnstile Integration:**
- **Approach**: Explicit rendering (programmatic control via JS API)
- **Script Loading**: `<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>`
- **Token Extraction**: Via callback function before form submission
- **Rationale**: More control over widget lifecycle than automatic rendering, better for programmatic submission
- **Affects**: Frontend form submission flow

**Frontend Validation:**
- **HTML5 Native**: `<input type="email" required>` for basic validation
- **TypeScript**: Regex check before API call for additional UX feedback
- **No Validation Library**: Keeps frontend lightweight
- **Rationale**: Single email field, native browser validation sufficient
- **Affects**: Frontend form component, user experience

**Form Submission Flow:**
1. User enters email
2. HTML5 validates on blur/submit
3. TypeScript validates before Turnstile render
4. Turnstile widget renders and completes
5. Token extracted via callback
6. POST to /api/submit with email + source + token
7. Display success/error message based on response

### Infrastructure & Deployment

**Container Architecture** (established in Step 3):
- 4 containers: db (volume), backend (Go API), frontend (nginx), proxy (nginx + HTTPS)
- Internal Docker network, only proxy exposed
- Multi-stage builds for minimal images

**Configuration Management:**
- **Backend**: Environment variables for all configuration
- **Frontend**: Vite .env files for build-time injection
- **Secrets**: TURNSTILE_SECRET, CORS_ORIGINS via Docker Compose .env
- **Validation**: Fail-fast on startup if required env vars missing

**Monitoring & Observability:**
- **Health Checks**: Simple HTTP endpoint on backend (/health)
- **Logging**: Structured logs via slog (stdout for container capture)
- **Startup Validation**: Environment check, database connection check, migration check

### Cross-Component Dependencies

**Implementation Sequence:**
1. **Backend foundation**: Go project structure, config, slog setup
2. **Database layer**: SQLite driver, golang-migrate integration, schema
3. **Middleware**: Rate limiting, CORS, logging middleware
4. **Handlers**: POST /api/submit handler with validation
5. **Frontend**: TypeScript form, Turnstile integration, API client
6. **Docker**: 4 Dockerfiles, docker-compose.yml, nginx config
7. **Integration tests**: Full stack testing

**Critical Dependencies:**
- Frontend depends on backend API contract (request/response format)
- Backend depends on Turnstile token verification
- CORS configuration must allow iframe embedding origins
- Rate limiter must be initialized before handler registration
- Migrations must run before database queries

### Technology Version Matrix

**Backend:**
- Go: 1.25.6
- mattn/go-sqlite3: Latest
- golang-migrate/migrate: v4.19.1
- rs/cors: Latest
- golang.org/x/time/rate: Latest (go get)

**Frontend:**
- Node.js: 24.13.0 LTS
- TypeScript: 5.x (latest)
- Vite: 7.3.1
- Cloudflare Turnstile: Latest (CDN)

**Infrastructure:**
- Docker Compose: v3.8 format
- nginx: alpine (latest)
- Alpine Linux: latest
- Distroless: gcr.io/distroless/static

### Deferred Decisions

No decisions deferred - all critical architecture finalized for implementation.

**Post-MVP Considerations** (not planned during development):
- Redis for distributed rate limiting (if multi-instance needed)
- PostgreSQL migration (if SQLite scale insufficient - highly unlikely)
- Admin dashboard (explicitly out of scope per PRD)

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Identified:** 18 areas where AI agents could make different implementation choices, now standardized for consistency.

### Naming Patterns

**Database Naming Conventions:**

- **Table names**: `snake_case`, lowercase, plural (e.g., `submissions`)
- **Column names**: `snake_case`, lowercase (e.g., `email`, `created_at`)
- **Primary keys**: `id` (simple, no table prefix)
- **Foreign keys**: `{table}_id` (e.g., `user_id`)
- **Timestamps**: `created_at`, `updated_at` (explicit naming, not `createdAt`)
- **Indexes**: `idx_{table}_{column}` (e.g., `idx_submissions_email`)

**Example Schema:**
```sql
CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    source TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_submissions_created_at ON submissions(created_at);
```

**API Naming Conventions:**

- **REST endpoints**: Lowercase with hyphens, descriptive verbs for actions (e.g., `/api/submit`)
- **JSON field names**: `snake_case` for consistency with database (e.g., `"email"`, `"created_at"`)
- **HTTP methods**: POST for mutations, GET for queries
- **Headers**: Standard casing (e.g., `Content-Type`, `X-Turnstile-Token` for custom headers)

**Go Code Naming Conventions:**

- **Exported types**: PascalCase (e.g., `SubmitRequest`, `SubmitResponse`)
- **Unexported types**: camelCase (e.g., `submitHandler`, `rateLimiter`)
- **Package names**: Lowercase, single word (e.g., `handler`, `middleware`, `database`)
- **Interface names**: PascalCase with "-er" suffix (e.g., `Validator`, `Submitter`)
- **Constants**: PascalCase for exported, camelCase for unexported (e.g., `DefaultRateLimit`)
- **Files**: `snake_case.go` (e.g., `submit_handler.go`, `rate_limiter.go`)
- **Tests**: `{file}_test.go` (e.g., `submit_handler_test.go`), co-located with implementation

**TypeScript Code Naming Conventions:**

- **Files**: `kebab-case.ts` (e.g., `form-handler.ts`, `api-client.ts`)
- **Interfaces/Types**: PascalCase (e.g., `SubmitRequest`, `ApiResponse`)
- **Functions**: camelCase (e.g., `handleSubmit`, `validateEmail`)
- **Variables/Constants**: camelCase (e.g., `turnstileToken`, `apiEndpoint`)
- **CSS classes**: `kebab-case` (e.g., `.submit-button`, `.error-message`)

### Structure Patterns

**Go Project Organization:**

```
dropmail-backend/
├── cmd/
│   └── api/
│       ├── main.go              # Entry point, server setup
│       └── main_test.go         # Integration tests (if needed)
├── internal/
│   ├── config/
│   │   ├── config.go            # Environment variable loading
│   │   └── config_test.go       # Co-located tests
│   ├── database/
│   │   ├── db.go                # SQLite connection, migrations
│   │   ├── queries.go           # Database operations
│   │   └── db_test.go           # Co-located tests
│   ├── handler/
│   │   ├── submit.go            # POST /api/submit handler
│   │   ├── health.go            # GET /health handler
│   │   └── submit_test.go       # Co-located handler tests
│   ├── middleware/
│   │   ├── cors.go              # CORS middleware
│   │   ├── ratelimit.go         # Rate limiting middleware
│   │   ├── logging.go           # Logging middleware
│   │   └── middleware_test.go   # Co-located middleware tests
│   └── model/
│       ├── submission.go        # Data models
│       └── validation.go        # Validation logic
├── migrations/
│   ├── 000001_init.up.sql       # Initial schema
│   └── 000001_init.down.sql     # Rollback
├── Dockerfile
├── go.mod
└── go.sum
```

**TypeScript Project Organization:**

```
dropmail-frontend/
├── src/
│   ├── main.ts                  # Entry point
│   ├── api-client.ts            # Backend API communication
│   ├── form-handler.ts          # Form submission logic
│   ├── turnstile.ts             # Turnstile integration
│   ├── validation.ts            # Client-side validation
│   └── style.css                # Global styles
├── index.html                   # Entry HTML
├── vite.config.ts
├── tsconfig.json
├── Dockerfile
└── package.json
```

**Test Location Rules:**

- **Go**: Tests co-located with implementation (`{file}_test.go`)
- **TypeScript**: Tests co-located with implementation (`{file}.test.ts`) or in `__tests__/` directory
- **Integration tests**: Separate `tests/` directory at project root (outside frontend/backend)

### Format Patterns

**API Response Formats:**

**Success Response:**
```json
{
  "success": true,
  "message": "Email submitted successfully"
}
```

**Error Response:**
```json
{
  "error": "Invalid email format"
}
```

**Field Naming in JSON:**
- Always `snake_case` (matches database, simpler transformation)
- Boolean fields: descriptive names (`is_valid`, not `valid`)
- Timestamps: ISO 8601 format when serialized (`"created_at": "2026-01-22T10:30:00Z"`)

**HTTP Status Codes:**
- `200 OK`: Successful submission
- `400 Bad Request`: Invalid email, missing fields, validation failure
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Database error, server failure

**Data Exchange Formats:**

**Request Body (POST /api/submit):**
```json
{
  "email": "user@example.com",
  "source": "portfolio-site",
  "turnstile_token": "0.abc123..."
}
```

**Go Struct Mapping:**
```go
type SubmitRequest struct {
    Email          string `json:"email"`
    Source         string `json:"source"`
    TurnstileToken string `json:"turnstile_token"`
}

type SubmitResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```

### Communication Patterns

**Middleware Order (CRITICAL - must be consistent):**

1. **CORS** - First, to handle preflight requests and set headers
2. **Rate Limiting** - Second, to reject excessive requests early
3. **Logging** - Third, to log all requests that pass rate limiting
4. **Handler** - Last, business logic execution

**Go Implementation:**
```go
// In cmd/api/main.go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/submit", handler.Submit)
    mux.HandleFunc("/health", handler.Health)

    // Wrap in reverse order (logging wraps rate limiting wraps CORS wraps mux)
    var h http.Handler = mux
    h = middleware.CORS(h)          // Applied first
    h = middleware.RateLimit(h)     // Applied second
    h = middleware.Logging(h)       // Applied third

    http.ListenAndServe(":8080", h)
}
```

**Logging Patterns:**

**What to Log:**
- **Successful submissions**: Log at `info` level with full request details (email, source, IP)
- **Validation failures**: Log at `warn` level with error reason
- **Rate limit exceeded**: Log at `warn` level with IP address
- **Server errors**: Log at `error` level with stack trace

**Log Format (slog structured):**
```go
// Successful submission
slog.Info("email submitted",
    "email", req.Email,
    "source", req.Source,
    "ip", r.RemoteAddr,
    "timestamp", time.Now())

// Validation failure
slog.Warn("validation failed",
    "email", req.Email,
    "reason", "invalid format",
    "ip", r.RemoteAddr)

// Rate limit
slog.Warn("rate limit exceeded",
    "ip", r.RemoteAddr)
```

**Sensitive Data Policy:**
- Emails CAN be logged (user explicitly provided for collection)
- IP addresses CAN be logged (needed for rate limiting analysis)
- Turnstile tokens should NOT be logged (security tokens)

### Process Patterns

**Validation Pattern (Dual Validation):**

**Frontend Validation (TypeScript):**
```typescript
// HTML5 validation first
<input type="email" required pattern="[^@]+@[^@]+\.[^@]+">

// TypeScript validation before Turnstile
function validateEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Validation flow
if (!validateEmail(emailInput.value)) {
    showError("Please enter a valid email address");
    return;
}
// Proceed to Turnstile
```

**Backend Validation (Go):**
```go
import "net/mail"

func validateEmail(email string) error {
    _, err := mail.ParseAddress(email)
    if err != nil {
        return fmt.Errorf("invalid email format")
    }
    return nil
}

// Always validate on backend regardless of frontend validation
```

**Turnstile Integration Pattern:**

**When to Render:** After successful email validation, on form submit button click

**Frontend Flow:**
```typescript
async function handleSubmit(event: Event) {
    event.preventDefault();

    // 1. HTML5 validation (automatic)
    // 2. TypeScript validation
    if (!validateEmail(emailInput.value)) {
        showError("Invalid email");
        return;
    }

    // 3. Render Turnstile widget (spawns challenge)
    const token = await turnstile.render(turnstileContainer, {
        sitekey: TURNSTILE_SITEKEY,
        callback: onTurnstileSuccess,
    });

    // 4. On callback, submit to API
}

function onTurnstileSuccess(token: string) {
    submitToAPI(emailInput.value, sourceParam, token);
}
```

**Rate Limiting Pattern:**

**Implementation:** Global rate limiter (not per-endpoint, since there's only one endpoint)

**Go Implementation:**
```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    r        rate.Limit  // requests per second
    b        int         // burst size
}

func (rl *RateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if _, exists := rl.limiters[ip]; !exists {
        rl.limiters[ip] = rate.NewLimiter(rl.r, rl.b)
    }

    return rl.limiters[ip].Allow()
}

// Middleware
func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        if !rateLimiter.Allow(ip) {
            slog.Warn("rate limit exceeded", "ip", ip)
            http.Error(w, `{"error":"Rate limit exceeded"}`, http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**Configuration:** Via environment variables (`RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`)

**Error Handling Pattern:**

**Backend Error Response:**
```go
func writeError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Usage
if err := validateEmail(req.Email); err != nil {
    writeError(w, "Invalid email format", http.StatusBadRequest)
    return
}
```

**Frontend Error Display:**
```typescript
async function submitToAPI(email: string, source: string, token: string) {
    try {
        const response = await fetch('/api/submit', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, source, turnstile_token: token })
        });

        const data = await response.json();

        if (data.success) {
            showSuccess(data.message);
        } else if (data.error) {
            showError(data.error);
        }
    } catch (err) {
        showError("Network error. Please try again.");
    }
}
```

**Loading State Pattern:**

**Frontend Loading UX:**
```typescript
function showLoading() {
    submitButton.disabled = true;
    submitButton.textContent = "Submitting...";
}

function hideLoading() {
    submitButton.disabled = false;
    submitButton.textContent = "Submit";
}

// Usage
showLoading();
await submitToAPI(email, source, token);
hideLoading();
```

**No global loading state needed** - Single form, local state sufficient

### Enforcement Guidelines

**All AI Agents MUST:**

1. **Follow middleware order**: CORS → Rate Limit → Logging → Handler (never reorder)
2. **Co-locate Go tests**: `{file}_test.go` in same directory as `{file}.go`
3. **Use snake_case for JSON/database**: All API fields and database columns must use snake_case
4. **Validate on both frontend and backend**: Never trust client-side validation alone
5. **Render Turnstile after validation**: Email validation must complete before Turnstile spawns
6. **Use slog for logging**: No fmt.Println or log.Println in production code
7. **Log successful submissions at info level**: Include email, source, IP, timestamp
8. **Use simple JSON error format**: `{"error": "message"}`, no error codes or nested structures
9. **Apply global rate limiting**: Single rate limiter instance per IP, not per endpoint
10. **Use PascalCase for Go exports**: Public types, functions, and constants must use PascalCase
11. **Use camelCase for TypeScript**: Variables, functions, and unexported members must use camelCase
12. **Follow standard Go project layout**: cmd/, internal/, migrations/ structure is mandatory

**Pattern Enforcement:**

**How to Verify Patterns are Followed:**
- Code review checklist comparing implementation to this document
- Automated linting for naming conventions (golangci-lint, ESLint)
- Integration tests verifying middleware order
- Unit tests confirming validation logic
- API contract tests ensuring response format matches specification

**Where to Document Pattern Violations:**
- Code review comments during PR process
- GitHub issues tagged with `architecture-deviation`
- Architecture decision log (ADL) if pattern needs revision

**Process for Updating Patterns:**
- Propose change in architecture.md via PR
- Document rationale for change
- Update affected code to match new pattern
- Update enforcement tooling (linters, tests)

### Pattern Examples

**Good Examples:**

**Go Handler Following All Patterns:**
```go
// internal/handler/submit.go
package handler

import (
    "encoding/json"
    "net/http"
    "net/mail"
    "time"
    "log/slog"

    "github.com/bastien/dropmail-backend/internal/database"
    "github.com/bastien/dropmail-backend/internal/model"
)

type SubmitRequest struct {
    Email          string `json:"email"`
    Source         string `json:"source"`
    TurnstileToken string `json:"turnstile_token"`
}

type SubmitResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func Submit(db *database.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse request
        var req SubmitRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validate email
        if _, err := mail.ParseAddress(req.Email); err != nil {
            slog.Warn("validation failed",
                "email", req.Email,
                "reason", "invalid format",
                "ip", r.RemoteAddr)
            writeError(w, "Invalid email format", http.StatusBadRequest)
            return
        }

        // Verify Turnstile (implementation omitted for brevity)

        // Save to database
        submission := model.Submission{
            Email:     req.Email,
            Source:    req.Source,
            CreatedAt: time.Now(),
        }

        if err := db.SaveSubmission(submission); err != nil {
            slog.Error("database error", "error", err)
            writeError(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        // Log success
        slog.Info("email submitted",
            "email", req.Email,
            "source", req.Source,
            "ip", r.RemoteAddr,
            "timestamp", time.Now())

        // Return success response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(SubmitResponse{
            Success: true,
            Message: "Email submitted successfully",
        })
    }
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}
```

**TypeScript Form Handler Following All Patterns:**
```typescript
// src/form-handler.ts
import { validateEmail } from './validation';
import { renderTurnstile } from './turnstile';
import { submitToAPI } from './api-client';

export class FormHandler {
    private emailInput: HTMLInputElement;
    private submitButton: HTMLButtonElement;
    private messageContainer: HTMLDivElement;

    constructor(formElement: HTMLFormElement) {
        this.emailInput = formElement.querySelector('#email')!;
        this.submitButton = formElement.querySelector('#submit')!;
        this.messageContainer = formElement.querySelector('#message')!;

        formElement.addEventListener('submit', this.handleSubmit.bind(this));
    }

    private async handleSubmit(event: Event): Promise<void> {
        event.preventDefault();

        const email = this.emailInput.value.trim();
        const source = new URLSearchParams(window.location.search).get('source') || 'unknown';

        // Frontend validation
        if (!validateEmail(email)) {
            this.showError('Please enter a valid email address');
            return;
        }

        // Show loading state
        this.setLoading(true);

        try {
            // Render Turnstile after validation
            const token = await renderTurnstile();

            // Submit to API
            const response = await submitToAPI(email, source, token);

            if (response.success) {
                this.showSuccess(response.message);
                this.emailInput.value = '';
            } else {
                this.showError(response.error || 'Submission failed');
            }
        } catch (error) {
            this.showError('Network error. Please try again.');
        } finally {
            this.setLoading(false);
        }
    }

    private setLoading(loading: boolean): void {
        this.submitButton.disabled = loading;
        this.submitButton.textContent = loading ? 'Submitting...' : 'Submit';
    }

    private showSuccess(message: string): void {
        this.messageContainer.textContent = message;
        this.messageContainer.className = 'message success';
    }

    private showError(message: string): void {
        this.messageContainer.textContent = message;
        this.messageContainer.className = 'message error';
    }
}
```

**Anti-Patterns (AVOID):**

**❌ Wrong Middleware Order:**
```go
// INCORRECT - Logging before rate limiting means we log rejected requests
var h http.Handler = mux
h = middleware.RateLimit(h)
h = middleware.Logging(h)  // Wrong order!
h = middleware.CORS(h)
```

**❌ Incorrect JSON Field Naming:**
```go
// INCORRECT - Using camelCase in JSON tags
type SubmitRequest struct {
    Email  string `json:"email"`      // ✓ Correct
    Source string `json:"sourceId"`   // ✗ Wrong - should be "source_id"
}
```

**❌ Missing Backend Validation:**
```go
// INCORRECT - Trusting client-side validation
func Submit(w http.ResponseWriter, r *http.Request) {
    var req SubmitRequest
    json.NewDecoder(r.Body).Decode(&req)

    // ✗ No validation - assumes frontend validated correctly
    db.SaveSubmission(req.Email, req.Source)
}
```

**❌ Turnstile Before Validation:**
```typescript
// INCORRECT - Rendering Turnstile before validating email
async function handleSubmit(event: Event) {
    event.preventDefault();

    // ✗ Wrong order - should validate email first
    const token = await renderTurnstile();

    if (!validateEmail(emailInput.value)) {
        showError("Invalid email");
        return;
    }
}
```

**❌ Inconsistent Test Location:**
```
dropmail-backend/
├── internal/
│   └── handler/
│       └── submit.go
└── tests/               # ✗ Wrong - Go tests should be co-located
    └── handler_test.go
```

**❌ Using fmt.Println for Logging:**
```go
// INCORRECT - Using fmt.Println instead of slog
func Submit(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received submission:", req.Email)  // ✗ Wrong

    // ✓ Correct
    slog.Info("email submitted", "email", req.Email)
}
```

**❌ Complex Error Response Format:**
```go
// INCORRECT - Overengineered error structure
type ErrorResponse struct {
    Error struct {
        Code    int    `json:"code"`     // ✗ Unnecessary
        Type    string `json:"type"`     // ✗ Unnecessary
        Message string `json:"message"`
        Details string `json:"details"`  // ✗ Unnecessary
    } `json:"error"`
}

// ✓ Correct - Simple error format
{"error": "Invalid email format"}
```

### Pattern Summary

**Consistency Checklist for AI Agents:**

- [ ] Middleware order: CORS → Rate Limit → Logging
- [ ] JSON/database fields: snake_case
- [ ] Go exports: PascalCase
- [ ] TypeScript: camelCase
- [ ] Tests co-located in Go projects
- [ ] Dual validation (frontend + backend)
- [ ] Turnstile after email validation
- [ ] Global rate limiting per IP
- [ ] Info-level logging for successful submissions with full request data
- [ ] Simple JSON error format: `{"error": "message"}`
- [ ] slog for structured logging
- [ ] Standard Go project layout (cmd/, internal/, migrations/)

**All patterns documented here are MANDATORY for ensuring consistent, maintainable implementation across all AI agents working on dropmail.**

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:**
All technology choices work together without conflicts:
- Go 1.25.6 with mattn/go-sqlite3 (CGO-enabled) and golang-migrate v4.19.1
- rs/cors and golang.org/x/time/rate integrate cleanly with net/http
- TypeScript 5.x with Vite 7.3.1 produces optimized static assets
- nginx alpine serves both static frontend and proxies to Go backend
- 4-container Docker architecture provides clean separation

**Pattern Consistency:**
All implementation patterns support the architectural decisions:
- Naming conventions (snake_case JSON/DB, PascalCase Go, camelCase TypeScript) are internally consistent
- Middleware order (CORS → Rate Limit → Logging → Handler) is explicitly defined and enforced
- Structure patterns align with Go standard layout and Vite conventions
- Communication patterns match the technology stack capabilities

**Structure Alignment:**
Project structure fully supports all architectural decisions:
- Backend cmd/api/ and internal/ packages follow Go conventions
- Frontend src/ organization supports TypeScript module imports
- Docker Compose orchestrates all 4 containers with proper networking
- Migrations directory enables versioned schema management

### Requirements Coverage Validation ✅

**Functional Requirements Coverage:**
All 49 functional requirements across 9 categories are architecturally supported:
- Form Submission (FR1-FR8): POST /api/submit with TypeScript + Go dual validation
- Spam Protection (FR9-FR13): Turnstile, rate limiting, HTTPS, CORS whitelist
- Source Tracking (FR14-FR17): URL parameter extraction and database storage
- Data Access (FR18-FR22): Direct SQLite CLI access, standard SQL queries
- Deployment (FR23-FR28): Docker Compose single-command, auto-migrations, env vars
- Multi-Property (FR29-FR33): iframe embedding with responsive container
- Responsive Design (FR34-FR38): Mobile-first CSS, modern browser support
- Accessibility (FR39-FR45): WCAG 2.1 AA semantic HTML patterns
- Performance (FR46-FR49): <1s TTI, <2s submission, <50KB page weight

**Non-Functional Requirements Coverage:**
All NFRs are addressed by architectural decisions:
- Performance: Vite optimized builds, Go net/http efficiency, SQLite local writes
- Security: nginx HTTPS termination, SQL parameter binding, per-IP rate limiting
- Reliability: Docker health checks, auto-migrations, minimal external dependencies
- Accessibility: Semantic HTML patterns, ARIA support, keyboard navigation

### Implementation Readiness Validation ✅

**Decision Completeness:**
- All technology versions pinned (Go 1.25.6, Node 24.13.0 LTS, Vite 7.3.1, golang-migrate v4.19.1)
- Library choices documented with rationale
- 18 potential conflict points identified and standardized
- Complete code examples provided for handlers, middleware, and TypeScript components

**Structure Completeness:**
- Backend directory structure fully specified (cmd/, internal/, migrations/)
- Frontend directory structure fully specified (src/ with modular files)
- Docker infrastructure complete (4 Dockerfiles, docker-compose.yml, nginx.conf)
- Test location rules defined (co-located for Go, flexible for TypeScript)

**Pattern Completeness:**
- Request/response formats defined with Go struct tags
- Error handling standardized (`{"error": "message"}` format)
- Logging patterns documented with slog examples
- Validation flow specified (frontend → Turnstile → backend)
- Anti-patterns documented with explanations

### Gap Analysis Results

**Critical Gaps:** None identified

**Technology Decisions (Intentional Improvements):**
- PRD suggested vanilla JavaScript → Architecture chose TypeScript for type safety
- PRD suggested Node.js backend → Architecture chose Go for smaller containers and better performance

These are valid architectural decisions that improve on initial suggestions while maintaining all requirements.

**Areas for Future Enhancement (Non-blocking):**
- nginx configuration patterns could be expanded
- Integration test examples could be more comprehensive
- CI/CD pipeline intentionally deferred to post-MVP

### Architecture Completeness Checklist

**✅ Requirements Analysis**
- [x] Project context thoroughly analyzed (49 FRs, comprehensive NFRs)
- [x] Scale and complexity assessed (Low complexity, single-purpose tool)
- [x] Technical constraints identified (Turnstile, SQLite, Docker Compose)
- [x] Cross-cutting concerns mapped (security, logging, validation, accessibility)

**✅ Architectural Decisions**
- [x] Critical decisions documented with versions
- [x] Technology stack fully specified
- [x] Integration patterns defined (4-container Docker architecture)
- [x] Performance considerations addressed (<50KB, <200ms API)

**✅ Implementation Patterns**
- [x] Naming conventions established (18 rules)
- [x] Structure patterns defined (Go standard layout, Vite conventions)
- [x] Communication patterns specified (middleware order, API contract)
- [x] Process patterns documented (validation, error handling, logging)

**✅ Project Structure**
- [x] Complete directory structure defined for frontend and backend
- [x] Component boundaries established (handlers, middleware, database)
- [x] Integration points mapped (container networking, API endpoints)
- [x] Requirements to structure mapping complete

### Architecture Readiness Assessment

**Overall Status:** READY FOR IMPLEMENTATION

**Confidence Level:** HIGH

**Key Strengths:**
- Comprehensive decision coverage eliminates architectural ambiguity
- Explicit version pinning prevents dependency drift
- Anti-patterns documented help AI agents avoid common mistakes
- Full code examples reduce interpretation errors
- 4-container architecture provides clean separation of concerns

**Areas for Future Enhancement:**
- CI/CD pipeline (post-MVP consideration)
- Distributed rate limiting if multi-instance needed (unlikely)
- Admin dashboard (explicitly out of scope per PRD)

### Implementation Handoff

**AI Agent Guidelines:**
- Follow all architectural decisions exactly as documented
- Use implementation patterns consistently across all components
- Respect project structure and boundaries
- Refer to this document for all architectural questions
- Check the Pattern Summary checklist before completing any implementation

**First Implementation Priority:**
```bash
# Frontend initialization
npm create vite@latest dropmail-frontend -- --template vanilla-ts

# Backend initialization
mkdir dropmail-backend && cd dropmail-backend
go mod init github.com/bastien/dropmail-backend
```

**Implementation Sequence:**
1. Backend foundation: Go project structure, config, slog setup
2. Database layer: SQLite driver, golang-migrate integration, schema
3. Middleware: Rate limiting, CORS, logging middleware
4. Handlers: POST /api/submit handler with validation
5. Frontend: TypeScript form, Turnstile integration, API client
6. Docker: 4 Dockerfiles, docker-compose.yml, nginx config
7. Integration tests: Full stack testing

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED ✅
**Total Steps Completed:** 8
**Date Completed:** 2026-01-23
**Document Location:** _bmad-output/planning-artifacts/architecture.md

### Final Architecture Deliverables

**Complete Architecture Document**
- All architectural decisions documented with specific versions
- Implementation patterns ensuring AI agent consistency
- Complete project structure with all files and directories
- Requirements to architecture mapping
- Validation confirming coherence and completeness

**Implementation Ready Foundation**
- 25+ architectural decisions made (technology stack, libraries, patterns)
- 18 implementation consistency rules defined
- 4 architectural components specified (db, backend, frontend, proxy)
- 49 functional requirements fully supported

**AI Agent Implementation Guide**
- Technology stack with verified versions (Go 1.25.6, Node 24.13.0, Vite 7.3.1)
- Consistency rules that prevent implementation conflicts
- Project structure with clear boundaries
- Integration patterns and communication standards

### Quality Assurance Checklist

**✅ Architecture Coherence**
- [x] All decisions work together without conflicts
- [x] Technology choices are compatible
- [x] Patterns support the architectural decisions
- [x] Structure aligns with all choices

**✅ Requirements Coverage**
- [x] All 49 functional requirements are supported
- [x] All non-functional requirements are addressed
- [x] Cross-cutting concerns are handled
- [x] Integration points are defined

**✅ Implementation Readiness**
- [x] Decisions are specific and actionable
- [x] Patterns prevent agent conflicts
- [x] Structure is complete and unambiguous
- [x] Examples are provided for clarity

---

**Architecture Status:** READY FOR IMPLEMENTATION ✅

**Next Phase:** Begin implementation using the architectural decisions and patterns documented herein.

**Document Maintenance:** Update this architecture when major technical decisions are made during implementation.

