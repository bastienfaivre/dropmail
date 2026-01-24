# Story 2.6: HTTPS Enforcement via Proxy

Status: review

## Story

As a **developer (Bastien)**,
I want **all traffic encrypted via HTTPS**,
So that **email submissions are secure in transit**.

## Acceptance Criteria

1. **HTTP to HTTPS Redirect**
   - **Given** the nginx proxy is configured
   - **When** a request arrives on port 80 (HTTP)
   - **Then** a 301 redirect is returned to the HTTPS version
   - **And** no content is served over HTTP

2. **HTTPS Termination**
   - **Given** HTTPS is configured with valid certificates
   - **When** a request arrives on port 443
   - **Then** TLS termination happens at the proxy
   - **And** internal traffic to backend/frontend is over HTTP (within Docker network)

3. **Browser Verification**
   - **Given** the proxy configuration is complete
   - **When** I access the form via HTTPS
   - **Then** the browser shows a secure connection
   - **And** mixed content warnings do not appear

4. **Security Headers**
   - **Given** security headers are configured
   - **When** any response is returned
   - **Then** appropriate headers are set (Content-Security-Policy, X-Frame-Options allowing iframe, etc.)

## Tasks / Subtasks

- [x] Task 1: Configure HTTPS (AC: #1, #2)
  - [x] 1.1: Update nginx.conf for HTTP to HTTPS redirect
  - [x] 1.2: Configure TLS settings (TLSv1.2/1.3, modern ciphers)
  - [x] 1.3: Add SSL certificate placeholder config

- [x] Task 2: Add Security Headers (AC: #4)
  - [x] 2.1: Add Content-Security-Policy header (frame-ancestors *)
  - [x] 2.2: Configure iframe embedding support (CSP frame-ancestors)
  - [x] 2.3: Add other security headers (X-Content-Type-Options, X-XSS-Protection, Referrer-Policy)

- [x] Task 3: Update Proxy Configuration (AC: #1, #2, #3)
  - [x] 3.1: Configure proxy_pass for backend and frontend
  - [x] 3.2: Set proper proxy headers (X-Real-IP, X-Forwarded-For, X-Forwarded-Proto)
  - [x] 3.3: Update docker-compose.yml with SSL volume mount

## Dev Notes

### Critical Architecture Requirements

**From Architecture Document [Source: architecture.md]:**

**Nginx SSL Configuration:**
```nginx
server {
    listen 80;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    # ... rest of config
}
```

**Security Headers:**
- Content-Security-Policy: Restrict resources
- X-Frame-Options: ALLOW-FROM for iframe embedding
- X-Content-Type-Options: nosniff
- Referrer-Policy: strict-origin-when-cross-origin

### Important Constraints

1. **Certificate Path**: SSL certs mounted at /etc/nginx/ssl
2. **iframe Support**: X-Frame-Options must allow embedding
3. **Internal HTTP**: Backend/frontend accessed via HTTP internally
4. **Docker Network**: All internal traffic on dropmail-internal network

### References

- [Source: architecture.md#Proxy-Config] - Nginx configuration
- [Source: epics.md#Story-2.6] - Original acceptance criteria

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- nginx.conf validated with proper HTTPS/redirect configuration
- Docker compose updated with SSL volume mount

### Completion Notes List

- HTTP to HTTPS redirect (301) on port 80
- HTTPS termination on port 443 with SSL/TLS
- Modern TLS configuration:
  - TLSv1.2 and TLSv1.3 protocols only
  - Strong cipher suite (ECDHE-ECDSA/RSA with AES-GCM)
  - SSL session cache (10m shared)
- Security headers configured:
  - X-Content-Type-Options: nosniff
  - X-XSS-Protection: 1; mode=block
  - Referrer-Policy: strict-origin-when-cross-origin
  - Content-Security-Policy: frame-ancestors * (for iframe embedding)
- Proxy configuration:
  - /api/ routes to backend:8080
  - /health endpoint for health checks
  - All other routes to frontend:80
  - Proper proxy headers (X-Real-IP, X-Forwarded-For, X-Forwarded-Proto)
- SSL certificates mounted from ./certs directory (read-only)
- Gzip compression enabled for text/CSS/JSON/JS

### Change Log
| Date | Change | Reason |
|------|--------|--------|
| 2026-01-23 | Story created | Created from epic definition |
| 2026-01-23 | HTTPS configuration complete | All tasks implemented |

### File List

**New:**
- docker/proxy/Dockerfile (nginx:alpine with SSL directory)

**Modified:**
- docker/proxy/nginx.conf (HTTPS, security headers, proxy config)
- docker-compose.yml (SSL volume mount for proxy)
