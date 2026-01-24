# Story 1.1: Initialize Project Structure with Starter Templates

Status: review

## Story

As a **developer (Bastien)**,
I want **the frontend and backend projects initialized with the specified starter templates and directory structure**,
So that **I have a consistent foundation that matches the architecture specification**.

## Acceptance Criteria

1. **Frontend Initialization**
   - **Given** the project root directory is empty
   - **When** I run the frontend initialization command `npm create vite@latest dropmail-frontend -- --template vanilla-ts`
   - **Then** the `dropmail-frontend/` directory is created with Vite vanilla-ts template
   - **And** TypeScript configuration is present with strict mode enabled
   - **And** the project uses Node.js 24.13.0 LTS and Vite 7.3.1

2. **Backend Initialization**
   - **Given** the project root directory exists
   - **When** I run the backend initialization commands (`mkdir dropmail-backend && cd dropmail-backend && go mod init github.com/bastien/dropmail-backend`)
   - **Then** the `dropmail-backend/` directory is created with Go module initialized
   - **And** the standard Go project layout is established (cmd/api/, internal/, migrations/)
   - **And** the project uses Go 1.25.6

3. **Directory Structure Verification**
   - **Given** both projects are initialized
   - **When** I review the directory structure
   - **Then** frontend follows: `src/`, `index.html`, `tsconfig.json`, `vite.config.ts`, `package.json`
   - **And** backend follows: `cmd/api/main.go`, `internal/config/`, `internal/handler/`, `internal/middleware/`, `internal/database/`, `internal/model/`

## Tasks / Subtasks

- [x] Task 1: Initialize Frontend Project (AC: #1)
  - [x] 1.1: Run `npm create vite@latest dropmail-frontend -- --template vanilla-ts`
  - [x] 1.2: Verify Vite 7.3.1 and TypeScript 5.x versions in package.json
  - [x] 1.3: Enable strict mode in tsconfig.json
  - [x] 1.4: Verify Node.js 24.13.0 LTS compatibility
  - [x] 1.5: Create initial src/ structure with main.ts, style.css, vite-env.d.ts

- [x] Task 2: Initialize Backend Project (AC: #2)
  - [x] 2.1: Create `dropmail-backend` directory
  - [x] 2.2: Run `go mod init github.com/bastien/dropmail-backend`
  - [x] 2.3: Verify Go 1.25.6 in go.mod
  - [x] 2.4: Create standard Go project layout directories:
    - [x] cmd/api/
    - [x] internal/config/
    - [x] internal/handler/
    - [x] internal/middleware/
    - [x] internal/database/
    - [x] internal/model/
    - [x] migrations/

- [x] Task 3: Create Placeholder Files (AC: #3)
  - [x] 3.1: Create `cmd/api/main.go` with minimal entry point
  - [x] 3.2: Create placeholder .go files in each internal package
  - [x] 3.3: Verify directory structure matches architecture specification

- [x] Task 4: Verify Both Projects Build (AC: #1, #2, #3)
  - [x] 4.1: Run `npm install` in frontend and verify success
  - [x] 4.2: Run `npm run build` in frontend and verify success
  - [x] 4.3: Run `go build ./...` in backend and verify success

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md#Starter-Template-Evaluation]:**

The architecture mandates specific starter templates and version pinning:

| Component | Version | Initialization Command |
|-----------|---------|----------------------|
| Frontend | Vite 7.3.1, TypeScript 5.x, Node 24.13.0 LTS | `npm create vite@latest dropmail-frontend -- --template vanilla-ts` |
| Backend | Go 1.25.6 | `go mod init github.com/bastien/dropmail-backend` |

**Standard Go Project Layout [Source: architecture.md#Go-Project-Organization]:**

```
dropmail-backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point, server setup
├── internal/
│   ├── config/
│   │   └── config.go            # Environment variable loading
│   ├── database/
│   │   └── db.go                # SQLite connection, migrations
│   ├── handler/
│   │   └── submit.go            # HTTP request handlers
│   ├── middleware/
│   │   └── cors.go              # CORS, rate limiting, logging
│   └── model/
│       └── submission.go        # Data models
├── migrations/
│   └── (empty for now)
├── go.mod
└── go.sum
```

**Frontend Project Layout [Source: architecture.md#TypeScript-Project-Organization]:**

```
dropmail-frontend/
├── src/
│   ├── main.ts                  # Entry point
│   ├── style.css                # Global styles
│   └── vite-env.d.ts            # Vite type definitions
├── index.html                   # Entry HTML
├── tsconfig.json                # TypeScript configuration
├── vite.config.ts               # Vite configuration
└── package.json
```

### Project Structure Notes

**Naming Conventions [Source: architecture.md#Naming-Patterns]:**
- Go files: `snake_case.go` (e.g., `submit_handler.go`)
- Go packages: lowercase, single word (e.g., `handler`, `middleware`)
- Go exports: PascalCase (e.g., `SubmitRequest`)
- TypeScript files: `kebab-case.ts` (e.g., `form-handler.ts`)

**Go Module Path:**
- Use `github.com/bastien/dropmail-backend` as the module path
- This matches the architecture specification for import paths

**TypeScript Configuration:**
- Enable strict mode for type safety
- Target ES2020 for modern browser support

### Important Constraints

1. **No Framework Dependencies**: Frontend uses vanilla TypeScript (no React, Vue, etc.)
2. **Minimal Dependencies**: Only add dependencies explicitly listed in architecture
3. **Standard Layouts**: Follow Go and Vite conventions exactly
4. **Version Pinning**: Verify exact versions match architecture spec

### References

- [Source: architecture.md#Starter-Template-Evaluation] - Complete starter template decisions
- [Source: architecture.md#Go-Project-Organization] - Backend directory structure
- [Source: architecture.md#TypeScript-Project-Organization] - Frontend directory structure
- [Source: architecture.md#Naming-Patterns] - All naming conventions
- [Source: project-context.md#Technology-Stack-Versions] - Pinned version requirements
- [Source: epics.md#Story-1.1] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Go 1.25.6 installed manually (not pre-installed in environment)
- CGO_ENABLED=0 required for placeholder Go files (no SQLite driver yet)

### Completion Notes List

- Frontend initialized with Vite 7.3.1 and TypeScript 5.9.3
- Backend initialized with Go 1.25.6 module at github.com/bastien/dropmail-backend
- All placeholder Go files follow architecture naming conventions
- Both projects build successfully (frontend: vite build, backend: go build ./...)
- Directory structures match architecture specification exactly

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Initial story file generation |
| 2026-01-23 | Story implemented | All 4 tasks completed, projects build successfully |

### File List

**Frontend (dropmail-frontend/):**
- dropmail-frontend/package.json (new)
- dropmail-frontend/tsconfig.json (new)
- dropmail-frontend/vite.config.ts (new)
- dropmail-frontend/index.html (new)
- dropmail-frontend/src/main.ts (new)
- dropmail-frontend/src/style.css (new)
- dropmail-frontend/src/vite-env.d.ts (new)
- dropmail-frontend/src/counter.ts (new - Vite template)
- dropmail-frontend/src/typescript.svg (new - Vite template)

**Backend (dropmail-backend/):**
- dropmail-backend/go.mod (new)
- dropmail-backend/cmd/api/main.go (new)
- dropmail-backend/internal/config/config.go (new)
- dropmail-backend/internal/database/db.go (new)
- dropmail-backend/internal/handler/submit.go (new)
- dropmail-backend/internal/middleware/cors.go (new)
- dropmail-backend/internal/model/submission.go (new)
