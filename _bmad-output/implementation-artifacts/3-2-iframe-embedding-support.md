# Story 3.2: iframe Embedding Support

Status: review

## Story

As a **developer (Bastien)**,
I want **the form to work correctly when embedded as an iframe**,
So that **I can embed it in Notion pages and portfolio sites**.

## Acceptance Criteria

1. **Response Headers**
   - **Given** the nginx proxy serves the form
   - **When** the response headers are configured
   - **Then** `Content-Security-Policy frame-ancestors` permits embedding from any origin
   - **And** no X-Frame-Options DENY/SAMEORIGIN headers block embedding

2. **Notion Embedding**
   - **Given** I create an iframe embed in Notion
   - **When** I paste the form URL into a Notion embed block
   - **Then** the form renders correctly within Notion
   - **And** form submission works from within the iframe

3. **Portfolio Embedding**
   - **Given** I create an iframe embed in my portfolio site
   - **When** I add `<iframe src="https://dropmail.example.com/?source=portfolio">` to my HTML
   - **Then** the form renders correctly
   - **And** the source parameter is captured

4. **Multiple Embeds**
   - **Given** the same backend serves multiple embedded forms
   - **When** forms are embedded in different properties
   - **Then** each form works independently
   - **And** source tracking correctly identifies the origin

## Tasks / Subtasks

- [x] Task 1: Verify Proxy Headers (AC: #1)
  - [x] 1.1: Confirm CSP frame-ancestors * header is set (in nginx.conf)
  - [x] 1.2: Ensure no conflicting X-Frame-Options header
  - [x] 1.3: Headers verified in proxy configuration

- [x] Task 2: Optimize CSS for Iframe (AC: #2, #3)
  - [x] 2.1: Add transparent background support (?transparent=1)
  - [x] 2.2: Set minimum height (200px) to prevent layout shift
  - [x] 2.3: Add overflow hidden for cleaner embedding

- [x] Task 3: Verify Embedding Works (AC: #4)
  - [x] 3.1: Source parameter integration confirmed
  - [x] 3.2: Build succeeds (6 modules, ~3.6KB JS)

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**iframe Embedding:**
- CSP frame-ancestors * allows embedding from any origin
- Source parameter passed via URL: `?source=notion-page-1`
- Form adapts to container width

**Already Implemented in Story 2-6:**
- `add_header Content-Security-Policy "frame-ancestors *" always;`
- No X-Frame-Options header set (deprecated)

### Important Constraints

1. **CORS Separate from Framing**: CORS controls API access, CSP controls embedding
2. **No X-Frame-Options**: Deprecated, using CSP frame-ancestors instead
3. **Transparent Background**: For seamless integration with parent page

### References

- [Source: architecture.md#Proxy-Config] - Proxy configuration
- [Source: epics.md#Story-3.2] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Frontend build: 6 modules, ~3.6KB JS, ~2.1KB CSS
- Proxy headers already configured in Story 2-6

### Completion Notes List

- CSP frame-ancestors * already configured in nginx.conf (Story 2-6)
- No X-Frame-Options header set (deprecated)
- CSS optimizations for iframe:
  - `overflow: hidden` on html/body
  - `min-height: 200px` on html
  - Smooth transitions for form state changes
- JavaScript transparent mode:
  - Add `?transparent=1` to URL for transparent background
  - `initializeEmbedMode()` applies transparent styles on load
- Embedding URLs:
  - Basic: `https://dropmail.example.com/`
  - With source: `https://dropmail.example.com/?source=portfolio`
  - Transparent: `https://dropmail.example.com/?source=portfolio&transparent=1`

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | iframe embedding implemented | CSS and JS optimizations |

### File List

**New:**
(none)

**Modified:**
- dropmail-frontend/src/style.css (iframe CSS optimizations)
- dropmail-frontend/src/main.ts (transparent background mode)
