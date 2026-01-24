# Story 3.1: Source Parameter Capture & Storage

Status: review

## Story

As a **developer (Bastien)**,
I want **to track which property each email submission comes from**,
So that **I can understand which sites drive the most engagement**.

## Acceptance Criteria

1. **URL Parameter Extraction**
   - **Given** the form is loaded with URL parameter `?source=portfolio-site`
   - **When** JavaScript extracts the source parameter
   - **Then** the source value "portfolio-site" is captured
   - **And** the source is included in the submission payload

2. **Default Source Handling**
   - **Given** the form is loaded without a source parameter
   - **When** JavaScript checks for the source
   - **Then** the source defaults to "unknown" or "direct"
   - **And** submissions still work correctly

3. **Database Storage**
   - **Given** a submission includes a source parameter
   - **When** the backend stores the submission
   - **Then** the source is stored in the `source` column
   - **And** I can query `SELECT * FROM submissions WHERE source = 'portfolio-site'`

4. **Multi-Source Querying**
   - **Given** multiple forms are embedded with different sources
   - **When** I query the database
   - **Then** I can see submissions grouped by source
   - **And** each submission has the correct source identifier

## Tasks / Subtasks

- [x] Task 1: Frontend Source Parameter Extraction (AC: #1, #2)
  - [x] 1.1: Create source utility module (src/source.ts)
  - [x] 1.2: Extract source from URL query parameters
  - [x] 1.3: Default to "direct" if source not provided
  - [x] 1.4: Integrate with form submission

- [x] Task 2: Verify Backend Support (AC: #3)
  - [x] 2.1: Confirm source field already in submission model
  - [x] 2.2: Confirm source stored in database
  - [x] 2.3: Update backend default to "direct" for consistency

- [x] Task 3: Test Source Tracking (AC: #4)
  - [x] 3.1: Source module sanitizes input (alphanumeric, hyphens, underscores)
  - [x] 3.2: Max length enforced (50 chars)
  - [x] 3.3: Frontend build succeeds (6 modules, ~3.4KB JS)

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**Source Parameter:**
- Captured from URL query string: `?source=portfolio-site`
- Default to "direct" if not provided
- Passed in JSON payload: `{"email": "...", "source": "portfolio-site", "turnstile_token": "..."}`

**Backend Already Supports:**
- Source field in SubmitRequest model
- Source column in submissions table
- Default "direct" if source empty

### Important Constraints

1. **URL Safety**: Sanitize source parameter input
2. **Max Length**: Consider limiting source length (e.g., 50 chars)
3. **Allowed Characters**: Alphanumeric, hyphens, underscores only

### References

- [Source: architecture.md#Source-Tracking] - Source tracking requirements
- [Source: epics.md#Story-3.1] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Frontend build: 6 modules transformed, ~3.4KB JS
- Source sanitization: alphanumeric, hyphens, underscores only
- Max length: 50 characters

### Completion Notes List

- Created source.ts module with:
  - `getSource()` function to extract source from URL params
  - Default to "direct" when no source parameter
  - Sanitization: alphanumeric, hyphens, underscores only
  - Max length: 50 characters
- Updated main.ts to import source from module
- Updated backend submit handler to default to "direct"
- Updated backend tests to expect "direct" default
- All builds successful

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Source tracking implemented | All tasks complete |

### File List

**New:**
- dropmail-frontend/src/source.ts (source parameter extraction)

**Modified:**
- dropmail-frontend/src/main.ts (import getSource from module)
- dropmail-backend/internal/handler/submit.go (default to "direct")
- dropmail-backend/internal/handler/submit_test.go (test expects "direct")
