---
stepsCompleted: ["step-01-validate-prerequisites", "step-02-design-epics", "step-03-create-stories", "step-04-final-validation"]
status: 'complete'
completedAt: '2026-01-23'
inputDocuments:
  - "/app/_bmad-output/planning-artifacts/prd.md"
  - "/app/_bmad-output/planning-artifacts/architecture.md"
---

# dropmail - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for dropmail, decomposing the requirements from the PRD and Architecture into implementable stories.

## Requirements Inventory

### Functional Requirements

**Form Submission & Collection (FR1-FR8):**
- FR1: Visitors can submit their email address via an embedded form
- FR2: Visitors can see immediate client-side validation feedback on email format
- FR3: Visitors can see an on-screen success message when their email is successfully submitted (no email confirmation sent)
- FR4: Visitors can see clear error messages when submission fails
- FR5: The system can validate email format on the server side before storing
- FR6: The system can reject invalid or malformed email addresses
- FR7: The system can store submitted emails with submission timestamp
- FR8: The system does NOT send verification emails or require email confirmation from visitors

**Spam Protection & Security (FR9-FR13):**
- FR9: The system can verify human users via Cloudflare Turnstile before accepting submissions
- FR10: The system can rate-limit submissions per IP address to prevent abuse
- FR11: The system can enforce HTTPS for all communications
- FR12: The system can sanitize inputs to prevent SQL injection and XSS attacks
- FR13: The system can restrict API access via CORS to authorized origins only

**Source Tracking & Analytics (FR14-FR17):**
- FR14: Bastien can embed forms with different source identifiers (e.g., portfolio-site, notion-page-x)
- FR15: The system can capture and store the source parameter with each submission
- FR16: Bastien can query submissions grouped by source to understand traffic sources
- FR17: Bastien can see submission timestamps for all collected emails

**Data Access & Export (FR18-FR22):**
- FR18: Bastien can access the database directly via standard SQLite client tools
- FR19: Bastien can query the total count of submissions
- FR20: Bastien can query submissions filtered by source
- FR21: Bastien can export email data to CSV format via SQL commands
- FR22: Bastien can view all submission data (email, source, timestamp) without additional tools

**Deployment & Configuration (FR23-FR28):**
- FR23: Bastien can deploy the entire system with a single Docker Compose command
- FR24: The system can initialize the database automatically on first run
- FR25: Bastien can configure Cloudflare Turnstile keys via environment variables
- FR26: Bastien can configure allowed CORS origins via environment variables
- FR27: The system can validate required environment variables on startup and fail fast with clear errors
- FR28: Bastien can access deployment health check endpoints to verify system status

**Multi-Property Embedding (FR29-FR33):**
- FR29: Bastien can embed the form as an iframe in Notion pages
- FR30: Bastien can embed the form as an iframe in portfolio websites
- FR31: Bastien can embed the form in multiple properties using the same backend
- FR32: The form can adapt its width to the parent container
- FR33: The form can maintain fixed minimum height to prevent layout shift

**Responsive Design & Device Support (FR34-FR38):**
- FR34: Visitors can use the form on mobile devices (320px-768px screens)
- FR35: Visitors can use the form on tablet devices (769px-1024px screens)
- FR36: Visitors can use the form on desktop devices (1025px+ screens)
- FR37: Visitors can interact with touch-friendly form controls on mobile (minimum 44px tap targets)
- FR38: The system can serve the form to modern browsers (Chrome, Firefox, Safari, Edge - last 2 versions)

**Accessibility & Keyboard Navigation (FR39-FR45):**
- FR39: Visitors can navigate the form using keyboard only (Tab key navigation)
- FR40: Visitors can submit the form by pressing Enter from the email input field
- FR41: Visitors using screen readers can understand form structure through semantic HTML
- FR42: Visitors using screen readers can hear validation errors and success messages
- FR43: Visitors can see clear focus indicators on all interactive elements
- FR44: Visitors can resize text up to 200% without loss of functionality
- FR45: Visitors can distinguish form states (normal, error, success) without relying on color alone

**Performance & Reliability (FR46-FR49):**
- FR46: The form can load and become interactive in under 1 second
- FR47: The system can respond to form submissions in under 2 seconds (including Turnstile verification)
- FR48: The system can maintain 100% uptime after initial deployment
- FR49: The system can operate with zero maintenance for extended periods (months)

### NonFunctional Requirements

**Performance:**
- NFR1: Time to Interactive (TTI) < 1 second for form loading
- NFR2: Initial page load < 500ms for form rendering
- NFR3: Total page weight < 50KB (HTML + CSS + minimal JavaScript)
- NFR4: Form submission complete within 2 seconds (including Turnstile verification)
- NFR5: API response time < 200ms at 95th percentile
- NFR6: Database write latency < 50ms
- NFR7: Frontend container < 10MB RAM, Backend container < 100MB RAM
- NFR8: Total system footprint < 150MB RAM

**Security:**
- NFR9: HTTPS enforcement for all communications (HTTP redirects to HTTPS)
- NFR10: Input sanitization on server-side to prevent SQL injection and XSS
- NFR11: Strict CORS whitelist of allowed origins via environment configuration
- NFR12: Per-IP rate limiting to prevent abuse and DoS attacks
- NFR13: Cloudflare Turnstile verification on all submissions
- NFR14: Environment variables for all sensitive configuration (never hardcoded)
- NFR15: Minimal data collection (email, source, timestamp only)

**Reliability:**
- NFR16: 100% availability target after initial deployment
- NFR17: Health check endpoints for container monitoring
- NFR18: Automatic database initialization on first run
- NFR19: Graceful degradation if Turnstile unavailable (fail closed with clear error)
- NFR20: Startup validation with fail-fast behavior for missing environment variables

**Accessibility:**
- NFR21: WCAG 2.1 AA compliance for all form components
- NFR22: Full keyboard accessibility (Tab navigation, Enter submission)
- NFR23: Screen reader support with semantic HTML and ARIA attributes
- NFR24: Color contrast ratio ≥ 4.5:1 for text
- NFR25: Focus indicators with ≥ 3:1 contrast
- NFR26: Text resizable up to 200% without functionality loss

### Additional Requirements

**From Architecture - Starter Template:**
- ARCH-1: Frontend initialized with Vite vanilla-ts template (`npm create vite@latest dropmail-frontend -- --template vanilla-ts`)
- ARCH-2: Backend initialized with Standard Go project layout (`go mod init github.com/bastien/dropmail-backend`)

**From Architecture - Container Architecture:**
- ARCH-3: 4-container Docker setup: db, backend, frontend, proxy
- ARCH-4: Internal Docker network with only proxy exposed to host
- ARCH-5: Multi-stage Docker builds for minimal images (< 40MB total stack)

**From Architecture - Technology Stack:**
- ARCH-6: Go 1.25.6 for backend
- ARCH-7: Node.js 24.13.0 LTS for frontend build
- ARCH-8: Vite 7.3.1 for frontend tooling
- ARCH-9: TypeScript 5.x for frontend type safety
- ARCH-10: SQLite with mattn/go-sqlite3 driver
- ARCH-11: golang-migrate v4.19.1 for database migrations
- ARCH-12: rs/cors for CORS middleware
- ARCH-13: golang.org/x/time/rate for rate limiting
- ARCH-14: slog for structured logging

**From Architecture - Implementation Patterns:**
- ARCH-15: Middleware order: CORS → Rate Limit → Logging → Handler
- ARCH-16: snake_case for JSON/database fields
- ARCH-17: Dual validation (frontend TypeScript + backend Go)
- ARCH-18: Turnstile rendering after email validation
- ARCH-19: Simple JSON error format: `{"error": "message"}`
- ARCH-20: Co-located tests for Go (`{file}_test.go`)

### FR Coverage Map

| FR | Epic | Description |
|----|------|-------------|
| FR1 | Epic 2 | Visitors can submit email via embedded form |
| FR2 | Epic 2 | Client-side validation feedback on email format |
| FR3 | Epic 2 | On-screen success message after submission |
| FR4 | Epic 2 | Clear error messages when submission fails |
| FR5 | Epic 2 | Server-side email format validation |
| FR6 | Epic 2 | Reject invalid/malformed email addresses |
| FR7 | Epic 2 | Store emails with submission timestamp |
| FR8 | Epic 2 | No verification emails required |
| FR9 | Epic 2 | Cloudflare Turnstile verification |
| FR10 | Epic 2 | Rate-limit submissions per IP |
| FR11 | Epic 2 | HTTPS enforcement |
| FR12 | Epic 2 | Input sanitization (SQL injection, XSS) |
| FR13 | Epic 2 | CORS restriction to authorized origins |
| FR14 | Epic 3 | Embed forms with source identifiers |
| FR15 | Epic 3 | Capture and store source parameter |
| FR16 | Epic 3 | Query submissions grouped by source |
| FR17 | Epic 3 | View submission timestamps |
| FR18 | Epic 3 | Direct database access via SQLite client |
| FR19 | Epic 3 | Query total submission count |
| FR20 | Epic 3 | Query submissions filtered by source |
| FR21 | Epic 3 | Export email data to CSV |
| FR22 | Epic 3 | View all submission data without additional tools |
| FR23 | Epic 1 | Deploy with single Docker Compose command |
| FR24 | Epic 1 | Automatic database initialization |
| FR25 | Epic 1 | Configure Turnstile keys via env vars |
| FR26 | Epic 1 | Configure CORS origins via env vars |
| FR27 | Epic 1 | Validate env vars on startup, fail fast |
| FR28 | Epic 1 | Health check endpoints |
| FR29 | Epic 3 | Embed form as iframe in Notion |
| FR30 | Epic 3 | Embed form as iframe in portfolio |
| FR31 | Epic 3 | Embed in multiple properties, same backend |
| FR32 | Epic 3 | Form adapts width to parent container |
| FR33 | Epic 3 | Fixed minimum height prevents layout shift |
| FR34 | Epic 3 | Form works on mobile (320px-768px) |
| FR35 | Epic 3 | Form works on tablet (769px-1024px) |
| FR36 | Epic 3 | Form works on desktop (1025px+) |
| FR37 | Epic 3 | Touch-friendly controls (44px tap targets) |
| FR38 | Epic 3 | Serve to modern browsers |
| FR39 | Epic 4 | Keyboard-only navigation (Tab) |
| FR40 | Epic 4 | Enter key submits form |
| FR41 | Epic 4 | Screen reader support via semantic HTML |
| FR42 | Epic 4 | Screen readers hear validation/success messages |
| FR43 | Epic 4 | Clear focus indicators |
| FR44 | Epic 4 | Text resizable to 200% |
| FR45 | Epic 4 | Form states distinguishable without color |
| FR46 | Epic 4 | Form loads in under 1 second |
| FR47 | Epic 4 | Submission completes in under 2 seconds |
| FR48 | Epic 4 | 100% uptime after deployment |
| FR49 | Epic 4 | Zero maintenance operation |

## Epic List

### Epic 1: Project Foundation & Deployable Infrastructure
Bastien can deploy the complete 4-container infrastructure with a single `docker-compose up` command and verify it's running via health checks.

**FRs covered:** FR23, FR24, FR25, FR26, FR27, FR28
**NFRs addressed:** NFR7, NFR8, NFR16, NFR17, NFR18, NFR20
**ARCH requirements:** ARCH-1 through ARCH-14

---

### Epic 2: Core Email Collection with Security
Visitors can submit their email through a secure form with spam protection (Turnstile), rate limiting, and proper validation. Bastien can see submissions in the database.

**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR7, FR8, FR9, FR10, FR11, FR12, FR13
**NFRs addressed:** NFR9, NFR10, NFR11, NFR12, NFR13, NFR14, NFR15
**ARCH requirements:** ARCH-15, ARCH-16, ARCH-17, ARCH-18, ARCH-19, ARCH-20

---

### Epic 3: Multi-Property Embedding & Source Tracking
Bastien can embed forms across Notion and portfolio sites with source tracking, responsive design across all devices, and data export capabilities.

**FRs covered:** FR14, FR15, FR16, FR17, FR18, FR19, FR20, FR21, FR22, FR29, FR30, FR31, FR32, FR33, FR34, FR35, FR36, FR37, FR38
**NFRs addressed:** NFR1, NFR2, NFR3

---

### Epic 4: Accessibility & Production Polish
Form is fully accessible to all users (WCAG 2.1 AA compliant), performs optimally, and runs reliably with zero maintenance.

**FRs covered:** FR39, FR40, FR41, FR42, FR43, FR44, FR45, FR46, FR47, FR48, FR49
**NFRs addressed:** NFR4, NFR5, NFR6, NFR19, NFR21, NFR22, NFR23, NFR24, NFR25, NFR26

---

## Epic 1: Project Foundation & Deployable Infrastructure

Bastien can deploy the complete 4-container infrastructure with a single `docker-compose up` command and verify it's running via health checks.

### Story 1.1: Initialize Project Structure with Starter Templates

As a **developer (Bastien)**,
I want **the frontend and backend projects initialized with the specified starter templates and directory structure**,
So that **I have a consistent foundation that matches the architecture specification**.

**Acceptance Criteria:**

**Given** the project root directory is empty
**When** I run the frontend initialization command `npm create vite@latest dropmail-frontend -- --template vanilla-ts`
**Then** the `dropmail-frontend/` directory is created with Vite vanilla-ts template
**And** TypeScript configuration is present with strict mode enabled
**And** the project uses Node.js 24.13.0 LTS and Vite 7.3.1

**Given** the project root directory exists
**When** I run the backend initialization commands (`mkdir dropmail-backend && cd dropmail-backend && go mod init github.com/bastien/dropmail-backend`)
**Then** the `dropmail-backend/` directory is created with Go module initialized
**And** the standard Go project layout is established (cmd/api/, internal/, migrations/)
**And** the project uses Go 1.25.6

**Given** both projects are initialized
**When** I review the directory structure
**Then** frontend follows: `src/`, `index.html`, `tsconfig.json`, `vite.config.ts`, `package.json`
**And** backend follows: `cmd/api/main.go`, `internal/config/`, `internal/handler/`, `internal/middleware/`, `internal/database/`, `internal/model/`

---

### Story 1.2: Docker Container Configuration

As a **developer (Bastien)**,
I want **all four Docker containers configured with multi-stage builds and Docker Compose orchestration**,
So that **I can deploy the entire system with a single command**.

**Acceptance Criteria:**

**Given** the frontend project exists with buildable TypeScript code
**When** I build the frontend Dockerfile
**Then** a multi-stage build produces an nginx:alpine image serving static assets
**And** the final image is < 10MB

**Given** the backend project exists with compilable Go code
**When** I build the backend Dockerfile
**Then** a multi-stage build produces a distroless/static image with the Go binary
**And** CGO is enabled for SQLite support
**And** the final image is < 15MB

**Given** all Dockerfiles are created
**When** I create docker-compose.yml
**Then** four services are defined: db, backend, frontend, proxy
**And** an internal Docker network `dropmail-internal` connects all containers
**And** only the proxy container exposes ports (80, 443) to the host
**And** a named volume `db-data` persists the SQLite database

**Given** docker-compose.yml exists
**When** I run `docker-compose build`
**Then** all four images build successfully
**And** total stack size is < 40MB

---

### Story 1.3: Environment Configuration & Startup Validation

As a **developer (Bastien)**,
I want **environment variables for Turnstile keys and CORS origins with fail-fast validation on startup**,
So that **misconfiguration is caught immediately rather than at runtime**.

**Acceptance Criteria:**

**Given** the backend application starts
**When** required environment variables are missing (TURNSTILE_SECRET, CORS_ORIGINS)
**Then** the application exits immediately with exit code 1
**And** a clear error message indicates which variable is missing
**And** the error message includes expected format/example

**Given** a `.env.example` file exists in the project root
**When** I copy it to `.env` and fill in values
**Then** docker-compose loads the environment variables correctly
**And** TURNSTILE_SECRET is available to the backend container
**And** CORS_ORIGINS is available to the backend container

**Given** valid environment variables are provided
**When** the backend application starts
**Then** configuration is logged at startup (without exposing secrets)
**And** the application proceeds to initialize

---

### Story 1.4: Health Checks & Database Initialization

As a **developer (Bastien)**,
I want **a health check endpoint and automatic database initialization**,
So that **I can verify the deployment is working and the database is ready**.

**Acceptance Criteria:**

**Given** the backend application is running
**When** I send a GET request to `/health`
**Then** the response status is 200 OK
**And** the response body is `{"status": "healthy"}`
**And** response time is < 50ms

**Given** the database file does not exist
**When** the backend application starts
**Then** SQLite database is created at the configured DB_PATH
**And** golang-migrate runs all pending migrations automatically
**And** the submissions table is created with columns: id, email, source, created_at

**Given** the database already exists with current migrations
**When** the backend application starts
**Then** golang-migrate detects no pending migrations
**And** the application starts normally without errors

**Given** all containers are configured
**When** I run `docker-compose up` from the project root
**Then** all four containers start successfully
**And** the proxy routes `/health` to the backend
**And** the health check returns 200 OK
**And** logs show successful database initialization

## Epic 2: Core Email Collection with Security

Visitors can submit their email through a secure form with spam protection (Turnstile), rate limiting, and proper validation. Bastien can see submissions in the database.

### Story 2.1: Backend Submit Endpoint with Validation

As a **visitor**,
I want **to submit my email address to the backend API**,
So that **my email is validated and stored securely for future updates**.

**Acceptance Criteria:**

**Given** the backend API is running
**When** I send a POST request to `/api/submit` with valid JSON `{"email": "user@example.com", "source": "test", "turnstile_token": "valid-token"}`
**Then** the response status is 200 OK
**And** the response body is `{"success": true, "message": "Email submitted successfully"}`
**And** the email is stored in the database with current timestamp

**Given** I submit an invalid email format (e.g., "not-an-email")
**When** the backend processes the request
**Then** the response status is 400 Bad Request
**And** the response body is `{"error": "Invalid email format"}`
**And** no record is created in the database

**Given** I submit a request with missing required fields
**When** the backend processes the request
**Then** the response status is 400 Bad Request
**And** the response body indicates which field is missing

**Given** a valid submission is processed
**When** the email is stored
**Then** SQL parameter binding is used (no string concatenation)
**And** the email is stored exactly as submitted (no XSS vectors executed)
**And** slog logs the submission at info level with email, source, IP, timestamp

---

### Story 2.2: CORS & Rate Limiting Middleware

As a **developer (Bastien)**,
I want **CORS restrictions and rate limiting on the API**,
So that **only authorized origins can embed the form and abuse is prevented**.

**Acceptance Criteria:**

**Given** CORS_ORIGINS is set to "https://example.com,https://notion.so"
**When** a request arrives from an origin not in the whitelist
**Then** the CORS preflight fails
**And** the request is rejected

**Given** a request arrives from an authorized origin
**When** the preflight OPTIONS request is sent
**Then** the response includes `Access-Control-Allow-Origin` header
**And** the response includes `Access-Control-Allow-Methods: POST`
**And** the actual request proceeds

**Given** rate limiting is configured (e.g., 10 requests per minute per IP)
**When** an IP exceeds the rate limit
**Then** the response status is 429 Too Many Requests
**And** the response body is `{"error": "Rate limit exceeded"}`
**And** slog logs the rate limit event at warn level with IP

**Given** middleware is configured
**When** a request is processed
**Then** middleware executes in order: CORS → Rate Limit → Logging → Handler

---

### Story 2.3: Cloudflare Turnstile Backend Verification

As a **developer (Bastien)**,
I want **the backend to verify Turnstile tokens before accepting submissions**,
So that **bot submissions are rejected**.

**Acceptance Criteria:**

**Given** a submission includes a valid Turnstile token
**When** the backend verifies the token with Cloudflare's siteverify API
**Then** verification succeeds
**And** the submission proceeds to storage

**Given** a submission includes an invalid or expired Turnstile token
**When** the backend verifies the token
**Then** verification fails
**And** the response status is 400 Bad Request
**And** the response body is `{"error": "Turnstile verification failed"}`
**And** no record is created in the database

**Given** a submission is missing the turnstile_token field
**When** the backend processes the request
**Then** the response status is 400 Bad Request
**And** the response body is `{"error": "Turnstile token required"}`

**Given** Cloudflare's siteverify API is unavailable
**When** the backend attempts verification
**Then** the submission fails closed (rejected)
**And** the response body is `{"error": "Verification service unavailable"}`
**And** slog logs the error at error level

---

### Story 2.4: Frontend Form with Client-Side Validation

As a **visitor**,
I want **a clean email form with immediate validation feedback**,
So that **I know if my email is valid before submitting**.

**Acceptance Criteria:**

**Given** I load the form page
**When** the page renders
**Then** I see an email input field with a label
**And** I see a submit button
**And** the form uses semantic HTML (`<form>`, `<label>`, `<input type="email">`, `<button>`)

**Given** I enter an invalid email format
**When** I blur the email field or attempt to submit
**Then** HTML5 validation displays an error message
**And** TypeScript validation provides additional feedback if needed

**Given** I submit a valid email successfully
**When** the API returns success
**Then** I see a success message: "Thanks! You're subscribed."
**And** the email input is cleared
**And** the form returns to initial state

**Given** the API returns an error
**When** the error response is received
**Then** I see the error message displayed clearly
**And** I can correct and resubmit

**Given** a submission is in progress
**When** the form is waiting for API response
**Then** the submit button is disabled
**And** the button text changes to "Submitting..."

---

### Story 2.5: Turnstile Widget Integration

As a **visitor**,
I want **invisible spam protection that doesn't interrupt my experience**,
So that **I can submit my email without solving puzzles**.

**Acceptance Criteria:**

**Given** I load the form page
**When** the page initializes
**Then** the Cloudflare Turnstile script is loaded asynchronously
**And** a container element exists for the widget

**Given** I enter a valid email and click submit
**When** the form validates successfully
**Then** the Turnstile widget renders (explicit rendering)
**And** verification happens (usually invisible)
**And** the token is captured via callback

**Given** Turnstile verification completes
**When** the callback fires with the token
**Then** the form submits to `/api/submit` with email, source, and turnstile_token
**And** the token is included in the JSON body

**Given** Turnstile verification fails (e.g., bot detected)
**When** the failure callback fires
**Then** an error message is displayed to the user
**And** the form is not submitted

---

### Story 2.6: HTTPS Enforcement via Proxy

As a **developer (Bastien)**,
I want **all traffic encrypted via HTTPS**,
So that **email submissions are secure in transit**.

**Acceptance Criteria:**

**Given** the nginx proxy is configured
**When** a request arrives on port 80 (HTTP)
**Then** a 301 redirect is returned to the HTTPS version
**And** no content is served over HTTP

**Given** HTTPS is configured with valid certificates
**When** a request arrives on port 443
**Then** TLS termination happens at the proxy
**And** internal traffic to backend/frontend is over HTTP (within Docker network)

**Given** the proxy configuration is complete
**When** I access the form via HTTPS
**Then** the browser shows a secure connection
**And** mixed content warnings do not appear

**Given** security headers are configured
**When** any response is returned
**Then** appropriate headers are set (Content-Security-Policy, X-Frame-Options allowing iframe, etc.)

## Epic 3: Multi-Property Embedding & Source Tracking

Bastien can embed forms across Notion and portfolio sites with source tracking, responsive design across all devices, and data export capabilities.

### Story 3.1: Source Parameter Capture & Storage

As a **developer (Bastien)**,
I want **to track which property each email submission comes from**,
So that **I can understand which sites drive the most engagement**.

**Acceptance Criteria:**

**Given** the form is loaded with URL parameter `?source=portfolio-site`
**When** JavaScript extracts the source parameter
**Then** the source value "portfolio-site" is captured
**And** the source is included in the submission payload

**Given** the form is loaded without a source parameter
**When** JavaScript checks for the source
**Then** the source defaults to "unknown" or "direct"
**And** submissions still work correctly

**Given** a submission includes a source parameter
**When** the backend stores the submission
**Then** the source is stored in the `source` column
**And** I can query `SELECT * FROM submissions WHERE source = 'portfolio-site'`

**Given** multiple forms are embedded with different sources
**When** I query the database
**Then** I can see submissions grouped by source
**And** each submission has the correct source identifier

---

### Story 3.2: iframe Embedding Support

As a **developer (Bastien)**,
I want **the form to work correctly when embedded as an iframe**,
So that **I can embed it in Notion pages and portfolio sites**.

**Acceptance Criteria:**

**Given** the nginx proxy serves the form
**When** the response headers are configured
**Then** `X-Frame-Options` is set to allow embedding from authorized origins
**And** `Content-Security-Policy frame-ancestors` permits the configured CORS origins

**Given** I create an iframe embed in Notion
**When** I paste the form URL into a Notion embed block
**Then** the form renders correctly within Notion
**And** form submission works from within the iframe

**Given** I create an iframe embed in my portfolio site
**When** I add `<iframe src="https://dropmail.example.com/?source=portfolio">` to my HTML
**Then** the form renders correctly
**And** the source parameter is captured

**Given** the same backend serves multiple embedded forms
**When** forms are embedded in different properties
**Then** each form works independently
**And** source tracking correctly identifies the origin

---

### Story 3.3: Responsive Form Layout

As a **visitor**,
I want **the form to look good and work well on any device**,
So that **I can subscribe from my phone, tablet, or desktop**.

**Acceptance Criteria:**

**Given** I view the form on a mobile device (320px-768px width)
**When** the page renders
**Then** the form displays in a single column layout
**And** the email input is full width
**And** the submit button is full width
**And** all tap targets are at least 44px

**Given** I view the form on a tablet (769px-1024px width)
**When** the page renders
**Then** the form is centered with constrained width
**And** elements have appropriate spacing for the viewport

**Given** I view the form on a desktop (1025px+ width)
**When** the page renders
**Then** the form is centered with optimal reading width
**And** the layout remains clean and professional

**Given** any device viewport
**When** I interact with the form
**Then** touch targets (input, button) are at least 44px in height
**And** font size is at least 16px (prevents iOS zoom on focus)

---

### Story 3.4: Browser Compatibility & Container Adaptation

As a **visitor**,
I want **the form to work in my modern browser**,
So that **I can subscribe regardless of which browser I use**.

**Acceptance Criteria:**

**Given** I use Chrome (last 2 versions)
**When** I load and submit the form
**Then** all functionality works correctly

**Given** I use Firefox (last 2 versions)
**When** I load and submit the form
**Then** all functionality works correctly

**Given** I use Safari (last 2 versions)
**When** I load and submit the form
**Then** all functionality works correctly

**Given** I use Edge (last 2 versions)
**When** I load and submit the form
**Then** all functionality works correctly

**Given** the form is embedded in a container of any width
**When** the iframe loads
**Then** the form adapts its width to 100% of the parent container
**And** the form maintains a fixed minimum height (e.g., 150px) to prevent layout shift

---

### Story 3.5: Data Access & Export Documentation

As a **developer (Bastien)**,
I want **clear documentation for accessing and exporting email data**,
So that **I can query submissions and export data when needed**.

**Acceptance Criteria:**

**Given** the SQLite database is running in the db container
**When** I connect via `docker exec -it dropmail-db sqlite3 /data/dropmail.db`
**Then** I can access the database directly
**And** I can run SQL queries interactively

**Given** I want to count total submissions
**When** I run `SELECT COUNT(*) FROM submissions;`
**Then** I get the total number of email submissions

**Given** I want to see submissions by source
**When** I run `SELECT source, COUNT(*) FROM submissions GROUP BY source;`
**Then** I see a breakdown of submissions per source

**Given** I want to view all submission data
**When** I run `SELECT email, source, created_at FROM submissions ORDER BY created_at DESC;`
**Then** I see all submissions with timestamps

**Given** I want to export to CSV
**When** I run `sqlite3 /data/dropmail.db ".mode csv" ".headers on" ".output /data/export.csv" "SELECT * FROM submissions;" ".exit"`
**Then** a CSV file is created with all submission data
**And** I can copy it from the container or access via the mounted volume

**Given** documentation needs to be provided
**When** the project is complete
**Then** a README section documents these access patterns
**And** example queries are provided for common operations

## Epic 4: Accessibility & Production Polish

Form is fully accessible to all users (WCAG 2.1 AA compliant), performs optimally, and runs reliably with zero maintenance.

### Story 4.1: Keyboard Navigation & Form Submission

As a **visitor using keyboard navigation**,
I want **to navigate and submit the form using only my keyboard**,
So that **I can subscribe without needing a mouse**.

**Acceptance Criteria:**

**Given** I navigate to the form using Tab key
**When** I press Tab from the page
**Then** focus moves to the email input field first
**And** pressing Tab again moves focus to the submit button
**And** pressing Tab again moves focus out of the form (no keyboard trap)

**Given** I have entered my email in the input field
**When** I press Enter
**Then** the form submits (same as clicking the submit button)
**And** Turnstile verification initiates

**Given** I am navigating with keyboard
**When** I interact with any form element
**Then** I never get trapped in a keyboard loop
**And** I can always Tab out of the form

**Given** the Turnstile widget is rendered
**When** I interact via keyboard
**Then** the widget is accessible via keyboard (handled by Cloudflare)
**And** verification can complete without mouse interaction

---

### Story 4.2: Screen Reader Support & ARIA

As a **visitor using a screen reader**,
I want **the form to be fully announced and understandable**,
So that **I can subscribe independently**.

**Acceptance Criteria:**

**Given** a screen reader user loads the form
**When** they navigate to the email input
**Then** the label "Email address" (or similar) is announced
**And** the input type (email) is announced
**And** any placeholder text is supplementary, not replacing the label

**Given** I submit an invalid email
**When** validation fails
**Then** the error message is announced via `aria-live="polite"` region
**And** the error is associated with the input via `aria-describedby`
**And** focus moves to or remains on the email input

**Given** I submit successfully
**When** the success message appears
**Then** the message "Thanks! You're subscribed." is announced via `aria-live="polite"`
**And** the announcement happens without requiring user action

**Given** the form is in a loading state
**When** submission is in progress
**Then** the "Submitting..." state is announced
**And** users understand the form is processing

**Given** any form element
**When** inspected by assistive technology
**Then** semantic HTML is used (`<form>`, `<label>`, `<input>`, `<button>`)
**And** ARIA attributes supplement (not replace) native semantics

---

### Story 4.3: Visual Accessibility

As a **visitor with visual impairments**,
I want **the form to have sufficient contrast and clear visual indicators**,
So that **I can see and interact with all elements**.

**Acceptance Criteria:**

**Given** any text element in the form
**When** I measure the color contrast
**Then** the contrast ratio is at least 4.5:1 against the background (WCAG AA)

**Given** any interactive element (input, button)
**When** the element receives keyboard focus
**Then** a visible focus indicator appears
**And** the focus indicator has at least 3:1 contrast ratio

**Given** I resize the browser text to 200%
**When** the page reflows
**Then** all form elements remain visible and functional
**And** no content is cut off or overlaps
**And** the form remains usable

**Given** the form displays an error state
**When** I view the error
**Then** the error is indicated by text AND an icon (not color alone)
**And** red color is supplemented with a warning icon or text prefix

**Given** the form displays a success state
**When** I view the success message
**Then** the success is indicated by text AND an icon (not color alone)
**And** green color is supplemented with a checkmark icon or text prefix

---

### Story 4.4: Performance Optimization

As a **visitor**,
I want **the form to load and respond quickly**,
So that **I don't wait or experience delays**.

**Acceptance Criteria:**

**Given** I navigate to the form page
**When** the page loads
**Then** Time to Interactive (TTI) is under 1 second
**And** the form is usable within 1 second of navigation

**Given** the form page loads
**When** I measure the total page weight
**Then** HTML + CSS + JavaScript totals less than 50KB
**And** Turnstile script loads asynchronously (not blocking)

**Given** I submit a valid email
**When** I click submit
**Then** the complete submission (Turnstile + API + response) takes under 2 seconds
**And** I see feedback (loading state → success/error) promptly

**Given** I load the form via iframe
**When** the iframe renders
**Then** First Contentful Paint (FCP) occurs within 500ms
**And** the form appears responsive immediately

**Given** the production build is created
**When** assets are served
**Then** CSS and JavaScript are minified
**And** appropriate cache headers are set for static assets

---

### Story 4.5: Production Reliability & Zero Maintenance

As a **developer (Bastien)**,
I want **the system to run reliably with zero maintenance**,
So that **I can deploy once and forget about it**.

**Acceptance Criteria:**

**Given** the system is deployed
**When** it runs for extended periods (weeks/months)
**Then** uptime remains at 100% (no crashes or restarts needed)
**And** no manual intervention is required

**Given** Docker containers are running
**When** a container encounters an error
**Then** Docker restart policy automatically restarts the container
**And** the system self-heals without manual intervention

**Given** the SQLite database grows over time
**When** thousands of submissions accumulate
**Then** performance remains stable
**And** database file size remains manageable (< 10MB for 100k emails)

**Given** Cloudflare Turnstile service is temporarily unavailable
**When** a submission is attempted
**Then** the system fails closed (rejects submission)
**And** a clear error message is shown
**And** the system recovers automatically when Turnstile returns

**Given** the system has minimal dependencies
**When** I review the architecture
**Then** external dependencies are limited to Cloudflare Turnstile only
**And** all other components are self-contained
**And** no scheduled maintenance tasks are required

**Given** logs are generated
**When** I need to debug an issue
**Then** structured logs (slog/JSON) are available via Docker logs
**And** logs include sufficient context (timestamps, request IDs, etc.)
