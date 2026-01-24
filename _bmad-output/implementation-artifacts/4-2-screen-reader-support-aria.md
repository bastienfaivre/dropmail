# Story 4.2: Screen Reader Support & ARIA

Status: done

## Story

As a **visitor using a screen reader**,
I want **the form to be fully announced and understandable**,
So that **I can subscribe independently**.

## Acceptance Criteria

1. **Label Announcement**
   - **Given** a screen reader user loads the form
   - **When** they navigate to the email input
   - **Then** the label "Email address" is announced
   - **And** the input type (email) is announced
   - **And** any placeholder text is supplementary, not replacing the label

2. **Error Announcement**
   - **Given** I submit an invalid email
   - **When** validation fails
   - **Then** the error message is announced via `aria-live="polite"` region
   - **And** the error is associated with the input via `aria-describedby`
   - **And** focus moves to or remains on the email input

3. **Success Announcement**
   - **Given** I submit successfully
   - **When** the success message appears
   - **Then** the message "Thanks! You're subscribed." is announced via `aria-live="polite"`
   - **And** the announcement happens without requiring user action

4. **Loading State Announcement**
   - **Given** the form is in a loading state
   - **When** submission is in progress
   - **Then** the "Submitting..." state is announced
   - **And** users understand the form is processing

5. **Semantic HTML**
   - **Given** any form element
   - **When** inspected by assistive technology
   - **Then** semantic HTML is used (`<form>`, `<label>`, `<input>`, `<button>`)
   - **And** ARIA attributes supplement (not replace) native semantics

## Tasks / Subtasks

- [x] Task 1: Verify Label Association (AC: #1, #5)
  - [x] 1.1: Label for="email" properly associated with input id="email"
  - [x] 1.2: Placeholder is supplementary ("you@example.com")
  - [x] 1.3: Semantic HTML: form, label, input, button elements

- [x] Task 2: Verify Error Announcements (AC: #2)
  - [x] 2.1: Error container has role="alert" aria-live="polite"
  - [x] 2.2: Input has aria-describedby="email-error"
  - [x] 2.3: aria-invalid="true" set on error via JavaScript

- [x] Task 3: Verify Success Announcement (AC: #3)
  - [x] 3.1: Success message has role="status" aria-live="polite"
  - [x] 3.2: Message announced when shown (hidden attribute removed)

- [x] Task 4: Add Loading State Announcement (AC: #4)
  - [x] 4.1: Added aria-busy="true" on form during submission
  - [x] 4.2: Added aria-label on button during loading state

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**Screen Reader Support:**
- Semantic HTML for native accessibility
- ARIA attributes to supplement, not replace
- aria-live regions for dynamic content
- Proper label associations

**Current Implementation:**
- Label with `for="email"` associated with input `id="email"`
- `aria-describedby="email-error"` on input
- `role="alert" aria-live="polite"` on error message
- `role="status" aria-live="polite"` on success message

### Important Constraints

1. **ARIA Best Practices**: Supplement, don't replace native semantics
2. **Live Regions**: Use polite level for non-urgent updates
3. **Focus Management**: Keep focus predictable

### References

- [Source: architecture.md#Accessibility] - Accessibility requirements
- [Source: epics.md#Story-4.2] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Frontend build: 6 modules, ~2.8KB CSS, ~3.7KB JS
- ARIA attributes verified in HTML and JavaScript

### Completion Notes List

- Label association:
  - `<label for="email">Email address</label>`
  - `<input id="email" type="email">`
  - Placeholder "you@example.com" is supplementary
- Error announcements:
  - Error div has `role="alert" aria-live="polite"`
  - Input has `aria-describedby="email-error"`
  - JavaScript sets `aria-invalid="true"` on error
  - Focus moves to input on validation error
- Success announcements:
  - Success div has `role="status" aria-live="polite"`
  - Hidden attribute removed triggers announcement
- Loading state:
  - Added `aria-busy="true/false"` on form
  - Added `aria-label` on button during submission
  - Button text changes to "Submitting..."

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Loading state ARIA added | aria-busy and aria-label for screen readers |
| 2026-01-24 | Code Review Fix | Fixed aria-busy not reset on success path |

## Senior Developer Review (AI)

**Reviewer:** Claude Opus 4.5
**Date:** 2026-01-24
**Outcome:** Changes Requested → Fixed

### Issues Found & Fixed
1. **[HIGH][Fixed]** aria-busy not reset on success - Added `setSubmitting(false)` before `showSuccess()` in main.ts:113

### Action Items
- [ ] [AI-Review][MEDIUM] Add automated accessibility tests for ARIA attributes

### File List

**New:**
(none)

**Modified:**
- dropmail-frontend/src/main.ts (aria-busy on form, aria-label on button)
