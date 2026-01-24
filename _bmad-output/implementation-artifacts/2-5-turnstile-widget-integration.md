# Story 2.5: Turnstile Widget Integration

Status: review

## Story

As a **visitor**,
I want **invisible spam protection that doesn't interrupt my experience**,
So that **I can submit my email without solving puzzles**.

## Acceptance Criteria

1. **Script Loading**
   - **Given** I load the form page
   - **When** the page initializes
   - **Then** the Cloudflare Turnstile script is loaded asynchronously
   - **And** a container element exists for the widget

2. **Explicit Rendering**
   - **Given** I enter a valid email and click submit
   - **When** the form validates successfully
   - **Then** the Turnstile widget renders (explicit rendering)
   - **And** verification happens (usually invisible)
   - **And** the token is captured via callback

3. **Token Submission**
   - **Given** Turnstile verification completes
   - **When** the callback fires with the token
   - **Then** the form submits to `/api/submit` with email, source, and turnstile_token
   - **And** the token is included in the JSON body

4. **Error Handling**
   - **Given** Turnstile verification fails (e.g., bot detected)
   - **When** the failure callback fires
   - **Then** an error message is displayed to the user
   - **And** the form is not submitted

## Tasks / Subtasks

- [x] Task 1: Add Turnstile Script (AC: #1)
  - [x] 1.1: Add Cloudflare Turnstile script to index.html
  - [x] 1.2: Use async loading with explicit render mode (?render=explicit)
  - [x] 1.3: Add VITE_TURNSTILE_SITE_KEY support (with test key fallback)

- [x] Task 2: Implement Turnstile Integration (AC: #2, #3)
  - [x] 2.1: Create turnstile.ts module with TypeScript types
  - [x] 2.2: Implement explicit rendering on form submit
  - [x] 2.3: Handle success callback with token
  - [x] 2.4: Submit form with token after verification

- [x] Task 3: Handle Errors (AC: #4)
  - [x] 3.1: Implement error callback and expired-callback
  - [x] 3.2: Display user-friendly error message
  - [x] 3.3: Allow retry after failure (resetTurnstile)

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**Turnstile Script:**
```html
<script src="https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit" async defer></script>
```

**Explicit Rendering:**
```typescript
turnstile.render('#turnstile-container', {
  sitekey: 'YOUR_SITE_KEY',
  callback: (token) => { /* handle token */ },
  'error-callback': () => { /* handle error */ },
});
```

**Environment Variable:**
- VITE_TURNSTILE_SITE_KEY for frontend access

### Important Constraints

1. **Explicit Render**: Use explicit mode to control when widget appears
2. **Async Loading**: Load script asynchronously to not block page
3. **Environment Config**: Site key via Vite environment variable
4. **No Puzzle**: Configure for invisible/managed mode

### References

- [Cloudflare Turnstile Client-side Rendering](https://developers.cloudflare.com/turnstile/get-started/client-side-rendering/)
- [Source: architecture.md#Turnstile-Config] - Turnstile configuration
- [Source: epics.md#Story-2.5] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- TypeScript compiles without errors
- Vite build successful: ~6KB total
- Turnstile script loaded asynchronously
- Test site key used as fallback when VITE_TURNSTILE_SITE_KEY not set

### Completion Notes List

- Turnstile script added to index.html with explicit render mode
- turnstile.ts module with full TypeScript typing:
  - Global Window interface extension for turnstile API
  - TurnstileOptions interface
  - waitForTurnstile() - polls for script load with timeout
  - renderTurnstile() - explicit rendering with callbacks
  - resetTurnstile() - reset widget for retry
- Form flow updated in main.ts:
  - Validate email first
  - Render Turnstile widget
  - On success: submit form with token
  - On error: show message, allow retry
- Error handling for:
  - Script load failure
  - Widget render failure
  - Verification failure
  - Token expiration

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Turnstile integration implemented | All tasks complete |

### File List

**New:**
- dropmail-frontend/src/turnstile.ts (new - Turnstile integration module)

**Modified:**
- dropmail-frontend/index.html (updated - Turnstile script tag)
- dropmail-frontend/src/main.ts (updated - integrated Turnstile flow)
