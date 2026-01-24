# Story 2.4: Frontend Form with Client-Side Validation

Status: review

## Story

As a **visitor**,
I want **a clean email form with immediate validation feedback**,
So that **I know if my email is valid before submitting**.

## Acceptance Criteria

1. **Form Rendering**
   - **Given** I load the form page
   - **When** the page renders
   - **Then** I see an email input field with a label
   - **And** I see a submit button
   - **And** the form uses semantic HTML (`<form>`, `<label>`, `<input type="email">`, `<button>`)

2. **Email Validation**
   - **Given** I enter an invalid email format
   - **When** I blur the email field or attempt to submit
   - **Then** HTML5 validation displays an error message
   - **And** TypeScript validation provides additional feedback if needed

3. **Success Handling**
   - **Given** I submit a valid email successfully
   - **When** the API returns success
   - **Then** I see a success message: "Thanks! You're subscribed."
   - **And** the email input is cleared
   - **And** the form returns to initial state

4. **Error Handling**
   - **Given** the API returns an error
   - **When** the error response is received
   - **Then** I see the error message displayed clearly
   - **And** I can correct and resubmit

5. **Loading State**
   - **Given** a submission is in progress
   - **When** the form is waiting for API response
   - **Then** the submit button is disabled
   - **And** the button text changes to "Submitting..."

## Tasks / Subtasks

- [x] Task 1: Create Form HTML Structure (AC: #1)
  - [x] 1.1: Update index.html with semantic form structure
  - [x] 1.2: Add email input with proper label and aria-describedby
  - [x] 1.3: Add submit button with 44px min-height tap target
  - [x] 1.4: Set up form for accessibility (aria-live, role attributes)

- [x] Task 2: Implement Form Styles (AC: #1)
  - [x] 2.1: Create/update CSS for form layout
  - [x] 2.2: Style input and button with dark mode support
  - [x] 2.3: Add focus states with visible focus rings

- [x] Task 3: Implement TypeScript Form Handler (AC: #2, #3, #4, #5)
  - [x] 3.1: Create form submission handler with fetch API
  - [x] 3.2: Add email validation (regex + blur validation)
  - [x] 3.3: Implement API call to /api/submit
  - [x] 3.4: Handle success response (show success message, hide form)
  - [x] 3.5: Handle error response (display API error message)
  - [x] 3.6: Implement loading state (disabled button, "Submitting..." text)

- [x] Task 4: Add Message Display (AC: #3, #4)
  - [x] 4.1: Create success message element with aria-live
  - [x] 4.2: Create error message element with role="alert"
  - [x] 4.3: Show/hide messages appropriately

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**API Endpoint:**
```
POST /api/submit
Content-Type: application/json

{"email": "user@example.com", "source": "test", "turnstile_token": "xxx"}
```

**Success Response:**
```json
{"success": true, "message": "Email submitted successfully"}
```

**Error Response:**
```json
{"error": "Invalid email format"}
```

### Important Constraints

1. **Semantic HTML**: Use proper form elements for accessibility
2. **TypeScript**: All JavaScript must be TypeScript
3. **Turnstile Token**: Form needs placeholder for Turnstile (Story 2.5)
4. **Source Parameter**: Extract from URL query string

### References

- [Source: architecture.md#Frontend] - Frontend structure
- [Source: epics.md#Story-2.4] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- TypeScript compiles without errors
- Vite build successful: ~5KB total (1KB HTML, 2KB CSS, 2KB JS)
- Well under 50KB page weight requirement

### Completion Notes List

- Semantic HTML form structure with proper labels
- Email validation on blur and submit
- Source parameter extracted from URL query string
- API call to /api/submit with JSON payload
- Success message: "Thanks! You're subscribed."
- Error message display with aria-live for screen readers
- Loading state: disabled button + "Submitting..." text
- Placeholder for Turnstile token (Story 2.5 integration)
- Dark mode support via prefers-color-scheme
- 44px min-height button for touch accessibility
- Exported functions for Turnstile integration: setTurnstileToken, getTurnstileToken, resetForm

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Frontend form implemented | All tasks complete |

### File List

**New:**
(none - updated existing files)

**Modified:**
- dropmail-frontend/index.html (updated - semantic form structure)
- dropmail-frontend/src/main.ts (updated - form submission handler)
- dropmail-frontend/src/style.css (updated - form styles)

**Removed:**
- dropmail-frontend/src/counter.ts (removed - unused Vite template file)
- dropmail-frontend/src/typescript.svg (removed - unused Vite template file)
