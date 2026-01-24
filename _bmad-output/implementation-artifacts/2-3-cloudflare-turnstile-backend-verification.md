# Story 2.3: Cloudflare Turnstile Backend Verification

Status: review

## Story

As a **developer (Bastien)**,
I want **the backend to verify Turnstile tokens before accepting submissions**,
So that **bot submissions are rejected**.

## Acceptance Criteria

1. **Valid Token Verification**
   - **Given** a submission includes a valid Turnstile token
   - **When** the backend verifies the token with Cloudflare's siteverify API
   - **Then** verification succeeds
   - **And** the submission proceeds to storage

2. **Invalid Token Rejection**
   - **Given** a submission includes an invalid or expired Turnstile token
   - **When** the backend verifies the token
   - **Then** verification fails
   - **And** the response status is 400 Bad Request
   - **And** the response body is `{"error": "Turnstile verification failed"}`
   - **And** no record is created in the database

3. **Missing Token Rejection**
   - **Given** a submission is missing the turnstile_token field
   - **When** the backend processes the request
   - **Then** the response status is 400 Bad Request
   - **And** the response body is `{"error": "Turnstile token required"}`

4. **Service Unavailable Handling**
   - **Given** Cloudflare's siteverify API is unavailable
   - **When** the backend attempts verification
   - **Then** the submission fails closed (rejected)
   - **And** the response body is `{"error": "Verification service unavailable"}`
   - **And** slog logs the error at error level

## Tasks / Subtasks

- [x] Task 1: Create Turnstile Verification Service (AC: #1, #2, #4)
  - [x] 1.1: Create internal/turnstile/verifier.go
  - [x] 1.2: Implement VerifyToken(token, remoteIP) method
  - [x] 1.3: Call Cloudflare siteverify API with TURNSTILE_SECRET
  - [x] 1.4: Parse response and return success/failure
  - [x] 1.5: Handle timeout and network errors (fail closed)

- [x] Task 2: Integrate with Submit Handler (AC: #1, #2, #3)
  - [x] 2.1: Add TurnstileVerifier interface to SubmitHandler
  - [x] 2.2: Verify token before storing submission
  - [x] 2.3: Return appropriate errors for invalid/missing tokens
  - [x] 2.4: Log verification failures

- [x] Task 3: Write Tests (AC: #1, #2, #3, #4)
  - [x] 3.1: Test valid token verification (mock API) - 6 verifier tests
  - [x] 3.2: Test invalid token rejection (via mock verifier)
  - [x] 3.3: Test API unavailable handling (timeout, invalid JSON, 500 errors)
  - [x] 3.4: Verify no database writes on verification failure (2 handler tests)

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**Turnstile Verification API:**
```
POST https://challenges.cloudflare.com/turnstile/v0/siteverify
Content-Type: application/x-www-form-urlencoded

secret=<TURNSTILE_SECRET>&response=<token>&remoteip=<client_ip>
```

**Response Format:**
```json
{
  "success": true|false,
  "error-codes": ["invalid-input-response", ...]
}
```

**Error Handling:**
- Fail closed: If verification service unavailable, reject submission
- Log all verification failures at appropriate levels

### Important Constraints

1. **Fail Closed**: Never accept submissions if Turnstile can't be verified
2. **Timeout**: Use reasonable timeout (5s) for Cloudflare API
3. **Secret Protection**: Never log the TURNSTILE_SECRET
4. **Client IP**: Include remoteip in verification request

### References

- [Cloudflare Turnstile Server-side Validation](https://developers.cloudflare.com/turnstile/get-started/server-side-validation/)
- [Source: architecture.md#Turnstile-Config] - Turnstile configuration
- [Source: epics.md#Story-2.3] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- 30 tests passing (6 turnstile + 6 middleware + 6 config + 12 model)
- Build verified with CGO_ENABLED=0
- Handler tests with mock verifier (9 tests with CGO)

### Completion Notes List

- Turnstile verifier service:
  - Calls Cloudflare siteverify API with form-encoded request
  - 5 second timeout for verification requests
  - Handles success/failure responses, timeouts, invalid JSON
  - Fails closed on any error (service unavailable, timeout, etc.)
  - slog logging at appropriate levels (warn for failures, error for service issues)
- TurnstileVerifier interface for dependency injection
- Updated SubmitHandler to verify tokens before storing
- Client IP extraction from X-Forwarded-For, X-Real-IP, or RemoteAddr
- 503 Service Unavailable returned when Turnstile API is down
- 400 Bad Request for invalid/missing tokens

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | Turnstile verification implemented | All tasks complete |

### File List

**New:**
- dropmail-backend/internal/turnstile/verifier.go (new - Turnstile API client)
- dropmail-backend/internal/turnstile/verifier_test.go (new - 6 tests with mock server)

**Modified:**
- dropmail-backend/internal/handler/submit.go (updated - TurnstileVerifier integration)
- dropmail-backend/internal/handler/submit_test.go (updated - mock verifier, 2 new tests)
- dropmail-backend/cmd/api/main.go (updated - inject Turnstile verifier)
