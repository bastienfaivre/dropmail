# Story 4.1: Keyboard Navigation & Form Submission

Status: done

## Story

As a **visitor using keyboard navigation**,
I want **to navigate and submit the form using only my keyboard**,
So that **I can subscribe without needing a mouse**.

## Acceptance Criteria

1. **Tab Navigation**
   - **Given** I navigate to the form using Tab key
   - **When** I press Tab from the page
   - **Then** focus moves to the email input field first
   - **And** pressing Tab again moves focus to the submit button
   - **And** pressing Tab again moves focus out of the form (no keyboard trap)

2. **Enter Key Submission**
   - **Given** I have entered my email in the input field
   - **When** I press Enter
   - **Then** the form submits (same as clicking the submit button)
   - **And** Turnstile verification initiates

3. **No Keyboard Trap**
   - **Given** I am navigating with keyboard
   - **When** I interact with any form element
   - **Then** I never get trapped in a keyboard loop
   - **And** I can always Tab out of the form

4. **Turnstile Accessibility**
   - **Given** the Turnstile widget is rendered
   - **When** I interact via keyboard
   - **Then** the widget is accessible via keyboard (handled by Cloudflare)
   - **And** verification can complete without mouse interaction

## Tasks / Subtasks

- [x] Task 1: Verify Tab Order (AC: #1)
  - [x] 1.1: Natural tab order (input → button) via semantic HTML
  - [x] 1.2: No tabindex attributes used - natural order preserved
  - [x] 1.3: Focus moves out of form after button (no trap)

- [x] Task 2: Verify Enter Key (AC: #2)
  - [x] 2.1: Form submit event handles Enter (native behavior)
  - [x] 2.2: Submit triggers handleSubmit → Turnstile verification

- [x] Task 3: Verify No Keyboard Traps (AC: #3)
  - [x] 3.1: No focus loops - semantic HTML elements
  - [x] 3.2: Turnstile container empty until rendered, no trap

- [x] Task 4: Enhance Focus Visibility (AC: #4)
  - [x] 4.1: Added :focus-visible styles for keyboard users
  - [x] 4.2: Transparent outline for Windows High Contrast mode

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**Keyboard Accessibility:**
- Tab navigation through form elements
- Enter key submits form
- No keyboard traps
- Turnstile widget handles own keyboard accessibility

**Current Implementation:**
- Form uses semantic HTML (form, label, input, button)
- Form submit event listener handles submission
- No custom tabindex values that would break natural order

### Important Constraints

1. **Semantic HTML**: Use native form elements for built-in accessibility
2. **Turnstile**: Cloudflare manages widget's keyboard accessibility
3. **Focus Management**: Don't interfere with natural tab order

### References

- [Source: architecture.md#Accessibility] - Accessibility requirements
- [Source: epics.md#Story-4.1] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Frontend build: 6 modules, ~2.8KB CSS, ~3.6KB JS
- Focus-visible styles added for keyboard users

### Completion Notes List

- Tab order works naturally via semantic HTML:
  - Email input (first tabbable element)
  - Submit button (second tabbable element)
  - Tab out of form (no trap)
- Enter key submission:
  - Native form behavior with type="submit" button
  - Form submit event triggers handleSubmit()
  - Turnstile verification initiates on valid email
- Focus styles enhanced:
  - :focus styles for all focus methods (mouse, keyboard)
  - :focus-visible styles for keyboard-only users
  - outline: 2px solid transparent for High Contrast mode
- Turnstile keyboard accessibility:
  - Handled by Cloudflare's widget implementation
  - Container is empty until verification needed

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Focus-visible styles added | Enhanced keyboard focus indicators |
| 2026-01-24 | Code Review | Reviewed with minor suggestions |

## Senior Developer Review (AI)

**Reviewer:** Claude Opus 4.5
**Date:** 2026-01-24
**Outcome:** Approved

### Issues Found
1. **[LOW][Noted]** No skip-to-content link - Form is simple enough that this is optional

### Action Items
- [ ] [AI-Review][MEDIUM] Add automated keyboard navigation tests

### File List

**New:**
(none)

**Modified:**
- dropmail-frontend/src/style.css (focus-visible styles for input and button)
