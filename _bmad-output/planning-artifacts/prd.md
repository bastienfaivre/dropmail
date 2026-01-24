---
stepsCompleted: ["step-01-init", "step-02-discovery", "step-03-success", "step-04-journeys", "step-05-domain", "step-06-innovation", "step-07-project-type", "step-08-scoping", "step-09-functional", "step-10-nonfunctional", "step-11-polish"]
inputDocuments: ["/app/_bmad-output/planning-artifacts/product-brief-dropmail-2026-01-22.md"]
workflowType: 'prd'
briefCount: 1
researchCount: 0
brainstormingCount: 0
projectDocsCount: 0
classification:
  projectType: "Web Application (iframe-embeddable form + REST API backend)"
  domain: "General/Developer Tools"
  complexity: "Low"
  projectContext: "Greenfield"
---

# Product Requirements Document - dropmail

**Author:** Bastien
**Date:** 2026-01-22

## Executive Summary

dropmail is a self-hosted email collection tool built for developers and creators who want data sovereignty without feature bloat. It provides an ultra-minimalist embeddable form (email field + submit button) paired with a decoupled backend API, enabling clean integration across multiple properties (Notion pages, portfolios, project sites) while maintaining full control over collected data.

**Key Differentiators:**
- **Purpose-built minimalism**: Designed solely for gathering interest signals, no marketing campaign features
- **Decoupled multi-property architecture**: One backend serves forms across Notion, portfolios, and future projects
- **Portfolio showcase value**: Clean architecture and security hardening demonstrate technical capability
- **Zero vendor dependencies**: Self-hosted with no free tier limits, branding requirements, or terms-of-service changes
- **Developer-first philosophy**: Built by engineers for engineers who value control, simplicity, and infrastructure ownership

**Target Users:**
- **Primary**: Bastien (solo developer/creator deploying and managing the tool)
- **Secondary**: Visitors to Notion pages and portfolio sites (form users)

## Success Criteria

### User Success (Primary User - Bastien)

**Deployment Success:**
- **Notion embed "aha!" moment**: Create embed block in Notion, paste iframe URL, form renders and works immediately without configuration
- Deploy once, works everywhere: Forms embed cleanly across Notion, portfolio, and future sites without modification per property
- Clean data collection: Email submissions land in database correctly formatted with no corruption
- Spam protection works: Cloudflare Turnstile blocks bots while legitimate submissions get through (>95% legitimate)
- Portfolio quality: Codebase demonstrates technical capability through clean architecture and engineering
- Zero maintenance: Once deployed, system runs reliably without intervention

**Operational Success:**
- Form loads and submits successfully with 100% uptime after deployment
- Email validation works: Invalid emails rejected at submission time
- Email submissions accumulate over time (validation that people care about your work)
- Direct database access for viewing/exporting works without additional tools
- Source tracking works: Can identify which property (portfolio vs Notion page) drove each submission

### User Success (Secondary Users - Visitors)

**Submission Experience:**
- Fast submission: Enter email, click submit, complete in <5 seconds
- Clear feedback: Immediate confirmation that submission worked
- No friction: No complex captchas, no redirects, no account creation requirements
- Clean interface: Unbranded form that doesn't distract from content

### Business Success

**Personal Value Metrics:**
- **Stack is completely free**: Self-hosted deployment with zero vendor costs or subscription fees
- **Actively gathering emails**: System successfully collecting email addresses from interested visitors
- **3-month success indicator**: Forms deployed and working across multiple properties (Notion + portfolio at minimum), with email submissions accumulating

**Brand Alignment:**
- No vendor branding contradicting "I build my own tools" personal narrative
- Data sovereignty maintained with zero third-party dependencies
- Tool itself serves as portfolio demonstration piece

### Technical Success

**Code Quality:**
- **Most simple and lightweight solution possible**: No unnecessary complexity or feature bloat
- **Easily maintainable**: Clear code structure that's straightforward to update and modify
- **Clean architecture**: Proper separation of concerns (frontend iframe / backend API / data layer)
- **Test coverage**: Core functionality tested to prevent regressions
- **Clean documentation**: Clear setup instructions and architecture documentation

**Security & Performance:**
- HTTPS enforcement working
- Input sanitization prevents injection attacks
- Turnstile verification on all submissions
- Rate limiting prevents abuse
- CORS properly configured for authorized domains only

### Measurable Outcomes

**Launch Readiness:**
- Docker Compose deployment succeeds on first try
- Form embeds successfully in Notion and portfolio
- First test submission stores correctly with source tracking
- Turnstile blocks bot submissions
- Rate limiting prevents spam floods

**Ongoing Success:**
- System uptime: 100% availability after deployment
- Spam prevention: >95% of submissions are legitimate
- Zero maintenance required: No interventions needed post-deployment
- Direct database queries work for email export

## Product Scope

### MVP Strategy & Philosophy

**MVP Approach:** Complete Utility Tool

dropmail's MVP is designed to be a fully functional, production-ready tool that solves the email collection problem completely. Unlike typical MVPs that are "minimum" stepping stones to a larger vision, dropmail's MVP is intentionally scoped to be the complete solution.

**Key Scoping Principles:**
- **No future feature planning during development**: Focus entirely on shipping a rock-solid MVP
- **Complete, not minimal**: Every feature in the MVP is there because it's essential for the tool to be useful
- **Quality over roadmap**: Portfolio-grade code quality, comprehensive tests, clear documentation
- **Pragmatic simplicity**: Choose the simplest technical solutions that meet requirements

**Philosophy:**
Build a tool that's so well-executed and complete that it requires zero maintenance or enhancement after launch. Future features aren't artificially prevented, but they're not planned or anticipated during development either.

### MVP Feature Set (Phase 1)

**Core User Journeys Supported:**
1. **Bastien's Deployment Journey**: Docker Compose deployment → Notion embed → form works
2. **Visitor's Submission Journey**: Discover form → submit email → instant confirmation
3. **Bastien's Operations Journey**: Direct database access → query submissions → export data

**Must-Have Capabilities:**

**Frontend (iframe-embeddable form):**
- Single HTML page with email input field and submit button
- Client-side email validation with immediate feedback
- Cloudflare Turnstile integration (spam protection - **non-negotiable requirement**)
- Clean, minimal CSS for beautiful UI
- Responsive design (mobile, tablet, desktop)
- Success/error message display
- Full WCAG 2.1 AA accessibility support

**Backend (REST API):**
- `POST /api/submit` endpoint for form submissions
- Server-side email validation and sanitization
- Cloudflare Turnstile token verification
- CORS middleware for cross-origin iframe embeds
- Rate limiting per IP (prevent DoS/spam)
- Database persistence layer

**Database:**
- **Simplest viable option**: SQLite (file-based, zero configuration, perfect for this use case)
- Schema: `email` (primary), `source` (tracking field - **required**), `created_at` (timestamp)
- Direct database access via standard SQLite client

**Deployment:**
- Docker Compose configuration
- Single-command deployment (`docker-compose up`)
- Automatic database initialization
- Environment variable management (Turnstile keys, CORS origins)
- HTTPS via reverse proxy

**Security (Required for Production):**
- HTTPS enforcement
- Input sanitization (prevent SQL injection, XSS)
- Cloudflare Turnstile on every submission
- Rate limiting per IP
- CORS whitelist for allowed origins
- Environment variables for secrets

**Source Tracking (Required):**
- URL parameter: `?source=portfolio-site` or `?source=notion-page-x`
- Extracted via JavaScript on form load
- Stored in database for analytics
- Essential for understanding which properties drive engagement

**Rationale for "Must-Have" Designation:**
- **Cloudflare Turnstile**: Without spam protection, the form becomes unusable due to bot submissions
- **Source tracking**: Essential for understanding which properties (portfolio vs Notion pages) drive interest
- **SQLite**: Simplest database that meets requirements - no PostgreSQL complexity needed for email collection
- **Docker Compose**: Single-command deployment is critical for "works on first try" success criteria
- **Accessibility**: Given form simplicity, full accessibility is achievable and demonstrates quality

### Post-MVP Features

**Approach:**
No features are planned for post-MVP during development. The focus is entirely on shipping a complete, well-executed tool.

**Possible Future Enhancements (Not Planned, Not Prioritized):**
If needs emerge after deployment based on real usage, potential quality-of-life improvements could include:
- CSV export API endpoint (if direct DB access becomes inconvenient)
- Email deduplication (if duplicate submissions become an issue)
- Basic submission statistics endpoint (if manual SQL queries are tedious)
- Migration from SQLite to PostgreSQL (if scale requires it)

**Explicitly Out of Scope:**
- Email verification/double opt-in workflows
- Admin dashboard or web UI
- Email sending functionality
- User accounts or authentication
- Analytics beyond basic counts
- Multi-user support or team features
- Feature expansion for its own sake

### Risk Mitigation Strategy

**Database Choice: SQLite**
- **Decision**: Use SQLite as the production database for MVP
- **Rationale**: Simplest possible database setup (zero configuration, file-based), more than sufficient for email collection workload, direct query access via sqlite3 CLI
- **Risk**: What if submission volume exceeds SQLite capabilities?
- **Mitigation**: SQLite handles thousands of concurrent readers and moderate write loads well. For a single-user email collection tool, SQLite will handle hundreds of thousands of submissions without issue. Migration to PostgreSQL is straightforward if needed (extremely unlikely).

**Cloudflare Turnstile Dependency**
- **Decision**: Cloudflare Turnstile is a required dependency for spam protection
- **Risk**: External dependency introduces potential failure point
- **Mitigation**: Turnstile is a mature, reliable service from Cloudflare (high uptime SLA). Graceful degradation: If Turnstile is down, fail closed (don't accept submissions) with clear error message. Monitor Cloudflare status page for outages.

**Deployment Complexity**
- **Risk**: Docker Compose deployment might fail on different environments
- **Mitigation**: Comprehensive deployment documentation with tested steps, environment variable validation on startup (fail fast with clear errors), health check endpoints for containers, test deployment on multiple platforms (Linux, macOS, Windows WSL2)

### Resource Requirements

**Development:**
- Solo developer (Bastien)
- No external dependencies or team coordination
- Portfolio-grade code quality expected

**Infrastructure:**
- Self-hosted server (already owned)
- Minimal resource requirements:
  - Frontend container: nginx (< 10MB RAM)
  - Backend container: Node.js API (< 100MB RAM)
  - Database: SQLite file (< 1MB for thousands of emails)
- Total footprint: < 150MB RAM, minimal CPU

**No External Costs:**
- Self-hosted (no vendor fees)
- Cloudflare Turnstile (free tier)
- SQLite (free, open-source)
- Docker (free, open-source)

## User Journeys

### Journey 1: Bastien - First Deployment & Integration

**Opening Scene: The Vendor Branding Problem**

Bastien is working on his portfolio site, proud of the clean minimalist design he's built. He wants to add a simple email collection form so people can stay updated on his projects. He tries a few SaaS options - Tally, Typeform, Google Forms - but each one feels wrong. The "Powered by Tally" footer at the bottom of the form mocks him. Here he is, a software engineer building custom tools, and he's using someone else's branded solution. It contradicts everything his portfolio is supposed to demonstrate: "I build my own solutions."

He could hack together a quick Express + SQLite script, but that feels hacky - not something he'd want to show in a technical interview. The heavyweight self-hosted options like Listmonk are overkill; he doesn't need campaign management, A/B testing, or automation workflows. He just needs one thing: a clean form that collects emails without vendor bloat.

**Rising Action: Building dropmail**

Bastien decides to build exactly what he needs. He sketches out the architecture: an ultra-minimalist iframe-embeddable form (just email field + submit button) backed by a simple REST API. He'll self-host it, so there's complete data sovereignty. He adds Cloudflare Turnstile for spam protection, CORS for cross-origin embeds, and rate limiting to prevent abuse. The entire stack is designed around one principle: build exactly what's needed, nothing more.

He writes clean, portfolio-grade code - proper separation of concerns, comprehensive tests, clear documentation. This isn't just a tool; it's a demonstration of his technical capability.

**Climax: The Docker Compose Moment**

It's deployment time. Bastien runs `docker-compose up` on his self-hosted server. The containers spin up, the database initializes automatically, and within seconds, the backend API is running. He configures his environment variables (Turnstile keys, CORS allowed origins), and the system is live.

Now for the real test: the Notion embed. He opens his Notion page about his latest project, creates an embed block, and pastes the iframe URL. The form renders instantly - clean, minimal, unbranded. He types in a test email, clicks submit, and watches the network request complete. He checks the database: the email is there, correctly formatted, with source tracking showing "notion-dropmail-page" and a timestamp.

**Resolution: Freedom & Portfolio Value**

Bastien embeds the same form on his portfolio site. Same iframe URL, different source parameter. No modifications needed. It just works.

The freedom hits him: no vendor costs, no submission limits, no branding requirements, no terms of service changes to worry about. His data lives in his own database. He can query it directly, export emails whenever he wants, and the system requires zero maintenance.

But the real win is subtler: when he shows his portfolio to potential clients or employers, dropmail itself becomes a talking point. "I built this email collection tool because existing solutions didn't fit my needs" - it demonstrates problem-solving, clean architecture, and the ability to ship production-grade code. It's both a tool AND proof of capability.

**Emotional Arc:** Frustration with vendor lock-in → Determination to build it right → Pride in clean deployment → Freedom and validation

### Journey 2: Sarah - Form Visitor Email Submission

**Opening Scene: Discovery**

Sarah stumbles upon Bastien's Notion page while researching solutions to a problem she's facing. The page showcases a project he built that solves exactly what she needs. As she reads through the documentation and examples, she thinks "I wish I could stay updated on this guy's work."

At the bottom of the page, she sees a simple, clean form: just an email field and a submit button. No clutter, no branding, no lengthy fields asking for her name, company, phone number, or life story. Just "Stay updated on future projects" with a single input field.

**Rising Action: The Decision**

Sarah appreciates the simplicity. No "Powered by [Vendor]" footer that makes her question whether her email is being sold to third parties. The form fits seamlessly into the Notion page design - it doesn't scream "I'm a third-party embed."

She trusts Bastien; his work is solid, and he clearly respects user experience. She decides to subscribe.

**Climax: Frictionless Submission**

Sarah types her email address: `sarah.chen@example.com`. As she types, client-side validation gives her immediate feedback if the format is wrong. She clicks the submit button.

A brief, invisible Cloudflare Turnstile verification happens (no annoying captcha puzzles). Within 2 seconds, a simple success message appears: "Thanks! You're subscribed."

No redirect to a thank-you page that breaks her flow. No verification email she needs to click. No account creation or password setup. Just instant confirmation.

**Resolution: Return to Browsing**

Sarah continues reading the rest of the Notion page, her flow completely uninterrupted. The entire submission took less than 5 seconds, and she's confident she'll hear about Bastien's future projects without any commitment or friction.

**Emotional Arc:** Interest in the project → Appreciation for simplicity → Trust in the creator → Satisfaction with frictionless completion

### Journey 3: Bastien - Ongoing Operations & Data Export

**Opening Scene: Curiosity About Traction**

Three months after deploying dropmail, Bastien wonders: "Is anyone actually interested in my work?" He's been embedding the form across his portfolio, Notion pages, and project showcases, but he hasn't checked the submissions yet.

**Rising Action: Direct Database Access**

Bastien opens his terminal and connects directly to the dropmail database using a simple SQLite client command. No need to log into a dashboard, no need to navigate through a web UI. Just a direct query:

```sql
SELECT COUNT(*) FROM submissions;
```

Result: 247 submissions.

His heart jumps. Nearly 250 people have signed up to hear about his work. He runs another query to see the breakdown by source:

```sql
SELECT source, COUNT(*) FROM submissions GROUP BY source;
```

Results:
- portfolio-site: 156
- notion-dropmail-page: 68
- notion-project-x: 23

His portfolio is driving the most interest, but his Notion pages are pulling their weight too. The source tracking is working perfectly.

**Climax: Export for Email Send**

Bastien is ready to send a manual update about his latest project. He needs to export the email list. He runs a simple query:

```sql
SELECT email, source, created_at FROM submissions ORDER BY created_at DESC;
```

He pipes the output to a CSV file:

```bash
sqlite3 dropmail.db ".mode csv" ".headers on" ".output emails_export.csv" "SELECT email, source, created_at FROM submissions;" ".exit"
```

Done. He has a clean CSV file with all submissions, ready to use with whatever email tool he wants (or just a manual BCC list).

**Resolution: Validation & Control**

The numbers tell Bastien something important: people care about his work. The 247 submissions validate that his projects are resonating. He owns this data completely - it's not locked in a vendor platform subject to terms changes or export limits.

The entire operation took less than 5 minutes: connect to database, run queries, export data. No dashboard login, no UI navigation, no export limits. Just direct access to his own data.

dropmail has been running for three months with zero maintenance, zero downtime, and zero vendor costs. It's doing exactly what it was built to do: quietly collecting interest signals while staying completely out of the way.

**Emotional Arc:** Curiosity about traction → Excitement at seeing numbers → Satisfaction with data control → Validation of effort

## Web Application Specific Requirements

### Project-Type Overview

dropmail is built as an **ultra-lightweight web application** consisting of:
- **Frontend**: Static HTML page with minimal vanilla JavaScript for form handling and client-side validation
- **Backend**: REST API endpoint for form submission processing
- **Delivery**: iframe-embeddable for seamless integration into Notion pages, portfolio sites, and other properties

The architecture prioritizes simplicity and maintainability over complex frameworks. No build tooling overhead, no heavy dependencies - just clean, straightforward web technologies that work reliably across modern browsers.

### Technical Architecture Considerations

**Frontend Architecture:**
- Static HTML + vanilla JavaScript (no framework dependencies)
- Lightweight, modular CSS for beautiful UI components
- Client-side email validation before submission
- Cloudflare Turnstile widget integration for spam protection
- Asynchronous form submission (no page reload)
- Success/error message handling within the same page

**Backend Architecture:**
- RESTful API endpoint (`POST /api/submit`)
- Server-side email validation and sanitization
- Cloudflare Turnstile token verification
- CORS middleware for cross-origin iframe embeds
- Rate limiting middleware (per-IP throttling)
- Database persistence layer (SQLite)

**Communication Flow:**
1. User enters email in iframe-embedded form
2. Client-side validation provides immediate feedback
3. Cloudflare Turnstile verification (invisible challenge)
4. POST request to backend API with email + source + Turnstile token
5. Backend validates, verifies Turnstile, checks rate limits
6. Database insert with timestamp and source tracking
7. Success/error response returned to frontend
8. UI updates with feedback message

### Browser Support Matrix

**Supported Browsers:**
- Chrome/Chromium: Last 2 versions
- Firefox: Last 2 versions
- Safari: Last 2 versions
- Edge (Chromium): Last 2 versions

**Rationale:**
Modern browser support is sufficient for the target audience (developers and tech-savvy creators). This allows use of contemporary web APIs without polyfills or transpilation overhead.

**Explicit Non-Support:**
- Internet Explorer (all versions)
- Legacy Edge (pre-Chromium)
- Browsers older than 2 versions back

**Testing Coverage:**
- Desktop: Chrome, Firefox, Safari, Edge
- Mobile: Safari (iOS), Chrome (Android)

### Responsive Design Requirements

**Design Approach:**
Mobile-first responsive design to ensure form works seamlessly across all device sizes.

**Breakpoints:**
- Mobile: 320px - 768px (single column, full-width form)
- Tablet: 769px - 1024px (centered form, constrained width)
- Desktop: 1025px+ (centered form, optimal reading width)

**Responsive Considerations:**
- Touch-friendly input fields (minimum 44px tap targets on mobile)
- Readable font sizes across all devices (minimum 16px to prevent zoom on iOS)
- Flexible iframe dimensions (width adapts to parent container)
- Form button sized appropriately for touch interaction
- Proper spacing for mobile thumb zones

**Iframe Embedding:**
- Width: 100% of parent container (or configurable via URL parameter)
- Height: Fixed minimum (e.g., 150px) to prevent layout shift
- Responsive scaling within parent context

### Performance Targets

**Load Performance:**
- **Initial page load**: < 500ms (form renders and is interactive)
- **Time to Interactive (TTI)**: < 1 second
- **Form submission**: < 2 seconds (including Turnstile verification and API response)
- **Total page weight**: < 50KB (HTML + CSS + minimal JS)

**API Performance:**
- **API response time**: < 200ms (95th percentile)
- **Database write latency**: < 50ms
- **Turnstile verification**: < 500ms (external dependency, best effort)

**Optimization Strategies:**
- Minify CSS and JavaScript in production
- Inline critical CSS to eliminate render-blocking requests
- Lazy-load Cloudflare Turnstile widget (only when needed)
- Enable HTTPS/2 for multiplexed requests
- Implement aggressive HTTP caching headers for static assets
- Use CDN delivery for form page (if performance becomes critical)

**Uptime Target:**
- 100% availability after initial deployment (self-hosted, no external dependencies beyond Cloudflare Turnstile)

### SEO Strategy

**Approach:**
SEO is explicitly **not required** for the iframe-embedded form page itself. The parent pages (Notion, portfolio) handle SEO independently.

**Form Page SEO Posture:**
- `<meta name="robots" content="noindex, nofollow">` on form page
- No sitemap inclusion
- No schema.org markup
- No Open Graph tags

**Rationale:**
The form page is never accessed directly by users or search engines - it's only viewed within an iframe context. Parent pages provide all necessary SEO structure and content.

### Accessibility Requirements

**WCAG 2.1 AA Compliance:**
Given the simplicity of the form (single email input + submit button), full accessibility support is achievable and valuable.

**Keyboard Navigation:**
- Full keyboard accessibility (tab order: email field → submit button)
- Enter key submits form from email input
- Focus indicators clearly visible on all interactive elements
- No keyboard traps within iframe

**Screen Reader Support:**
- Proper semantic HTML (`<form>`, `<label>`, `<input type="email">`, `<button>`)
- ARIA labels where necessary for clarity
- Form validation errors announced to screen readers
- Success/error messages announced dynamically (`aria-live` regions)
- Cloudflare Turnstile widget accessible (Cloudflare handles this)

**Visual Accessibility:**
- Color contrast ratio ≥ 4.5:1 for text (WCAG AA standard)
- Focus indicators meet 3:1 contrast requirement
- Text resizable up to 200% without loss of functionality
- No reliance on color alone to convey information (e.g., error states use text + icons)

**Form Accessibility Specifics:**
- `<label>` explicitly associated with email `<input>` (for attribute + id)
- Clear placeholder text (supplementary, not replacing label)
- Descriptive submit button text ("Subscribe" or "Submit")
- Client-side validation errors displayed both visually and programmatically
- Error messages associated with input field (`aria-describedby`)

**Testing:**
- Automated: axe DevTools, Lighthouse accessibility audit
- Manual: Keyboard-only navigation testing
- Screen reader: NVDA (Windows) and VoiceOver (macOS) testing

### CSS Component Integration

**Approach:**
Support for beautiful, modern CSS component libraries while maintaining lightweight footprint.

**Integration Options:**
- **Utility-first CSS**: Tailwind CSS (via CDN or minimal custom build)
- **Component libraries**: Shadcn-style components (if desired)
- **Custom CSS**: Hand-crafted modular CSS for maximum control

**Design Constraints:**
- Minimal CSS footprint (target: < 20KB)
- No JavaScript dependencies for styling (CSS-only components preferred)
- Clean, modern aesthetic that embeds well in Notion and portfolio contexts
- Customizable theme variables (colors, fonts, spacing) via CSS custom properties

**Visual Design Priorities:**
- Clean, minimal interface
- Unbranded appearance (no dropmail branding)
- Subtle, professional styling that blends with parent page
- Clear visual feedback for interactions (hover, focus, success, error states)
- Smooth transitions and micro-interactions where appropriate

### Implementation Considerations

**Development Stack:**
- Static HTML page (index.html)
- Vanilla JavaScript (form.js - minimal, focused functionality)
- Modern CSS with custom properties (styles.css)
- No build tooling required (optional minification for production)

**Backend Technology:**
- Node.js + Express (or similar lightweight framework)
- SQLite for database
- Environment-based configuration
- CORS middleware
- Rate limiting middleware (e.g., express-rate-limit)
- Input validation and sanitization library

**Deployment:**
- Docker container for frontend (nginx serving static files)
- Docker container for backend API
- Docker container for database (or managed database service)
- Docker Compose orchestration
- HTTPS via reverse proxy (nginx or Caddy)

**Security Hardening:**
- HTTPS enforcement (redirect HTTP → HTTPS)
- Strict CORS policy (whitelist allowed origins via environment config)
- Content Security Policy (CSP) headers
- Input sanitization on backend (prevent SQL injection, XSS)
- Rate limiting per IP address (prevent spam/DoS)
- Cloudflare Turnstile verification on every submission
- Environment variables for sensitive configuration (Turnstile keys, database credentials)

**Source Tracking Implementation:**
- URL parameter: `?source=portfolio-site` or `?source=notion-page-x`
- Extracted via JavaScript on form load
- Included in API submission payload
- Stored in database `source` field for analytics

## Functional Requirements

### Form Submission & Collection

- FR1: Visitors can submit their email address via an embedded form
- FR2: Visitors can see immediate client-side validation feedback on email format
- FR3: Visitors can see an on-screen success message when their email is successfully submitted (no email confirmation sent)
- FR4: Visitors can see clear error messages when submission fails
- FR5: The system can validate email format on the server side before storing
- FR6: The system can reject invalid or malformed email addresses
- FR7: The system can store submitted emails with submission timestamp
- FR8: The system does NOT send verification emails or require email confirmation from visitors

### Spam Protection & Security

- FR9: The system can verify human users via Cloudflare Turnstile before accepting submissions
- FR10: The system can rate-limit submissions per IP address to prevent abuse
- FR11: The system can enforce HTTPS for all communications
- FR12: The system can sanitize inputs to prevent SQL injection and XSS attacks
- FR13: The system can restrict API access via CORS to authorized origins only

### Source Tracking & Analytics

- FR14: Bastien can embed forms with different source identifiers (e.g., portfolio-site, notion-page-x)
- FR15: The system can capture and store the source parameter with each submission
- FR16: Bastien can query submissions grouped by source to understand traffic sources
- FR17: Bastien can see submission timestamps for all collected emails

### Data Access & Export

- FR18: Bastien can access the database directly via standard SQLite client tools
- FR19: Bastien can query the total count of submissions
- FR20: Bastien can query submissions filtered by source
- FR21: Bastien can export email data to CSV format via SQL commands
- FR22: Bastien can view all submission data (email, source, timestamp) without additional tools

### Deployment & Configuration

- FR23: Bastien can deploy the entire system with a single Docker Compose command
- FR24: The system can initialize the database automatically on first run
- FR25: Bastien can configure Cloudflare Turnstile keys via environment variables
- FR26: Bastien can configure allowed CORS origins via environment variables
- FR27: The system can validate required environment variables on startup and fail fast with clear errors
- FR28: Bastien can access deployment health check endpoints to verify system status

### Multi-Property Embedding

- FR29: Bastien can embed the form as an iframe in Notion pages
- FR30: Bastien can embed the form as an iframe in portfolio websites
- FR31: Bastien can embed the form in multiple properties using the same backend
- FR32: The form can adapt its width to the parent container
- FR33: The form can maintain fixed minimum height to prevent layout shift

### Responsive Design & Device Support

- FR34: Visitors can use the form on mobile devices (320px-768px screens)
- FR35: Visitors can use the form on tablet devices (769px-1024px screens)
- FR36: Visitors can use the form on desktop devices (1025px+ screens)
- FR37: Visitors can interact with touch-friendly form controls on mobile (minimum 44px tap targets)
- FR38: The system can serve the form to modern browsers (Chrome, Firefox, Safari, Edge - last 2 versions)

### Accessibility & Keyboard Navigation

- FR39: Visitors can navigate the form using keyboard only (Tab key navigation)
- FR40: Visitors can submit the form by pressing Enter from the email input field
- FR41: Visitors using screen readers can understand form structure through semantic HTML
- FR42: Visitors using screen readers can hear validation errors and success messages
- FR43: Visitors can see clear focus indicators on all interactive elements
- FR44: Visitors can resize text up to 200% without loss of functionality
- FR45: Visitors can distinguish form states (normal, error, success) without relying on color alone

### Performance & Reliability

- FR46: The form can load and become interactive in under 1 second
- FR47: The system can respond to form submissions in under 2 seconds (including Turnstile verification)
- FR48: The system can maintain 100% uptime after initial deployment
- FR49: The system can operate with zero maintenance for extended periods (months)

## Non-Functional Requirements

### Performance

**Load Time Requirements:**
- **Time to Interactive (TTI)**: Form must load and become interactive within 1 second
- **Initial page load**: < 500ms for form rendering
- **Total page weight**: < 50KB (HTML + CSS + minimal JavaScript)

**Response Time Requirements:**
- **Form submission**: Complete within 2 seconds (including Cloudflare Turnstile verification and API response)
- **API response time**: < 200ms at 95th percentile (excluding Turnstile verification)
- **Database write latency**: < 50ms for email insertion

**Resource Efficiency:**
- **Frontend container**: < 10MB RAM usage (nginx serving static files)
- **Backend container**: < 100MB RAM usage (Node.js API)
- **Database**: < 1MB for thousands of email submissions (SQLite file-based)
- **Total system footprint**: < 150MB RAM, minimal CPU utilization

### Security

**Data Protection:**
- **HTTPS enforcement**: All communications over encrypted HTTPS (HTTP redirects to HTTPS)
- **Input sanitization**: All user inputs sanitized on server-side to prevent SQL injection and XSS attacks
- **Data at rest**: Database file protected with appropriate filesystem permissions
- **Environment variables**: All sensitive configuration (Turnstile keys, database credentials) stored in environment variables, never hardcoded

**Access Control:**
- **CORS policy**: Strict whitelist of allowed origins via environment configuration (only authorized domains can embed the form)
- **Rate limiting**: Per-IP rate limiting to prevent abuse and DoS attacks (maximum N submissions per IP per time window)
- **API authentication**: Backend API validates Cloudflare Turnstile tokens on every submission

**Abuse Prevention:**
- **Cloudflare Turnstile**: Required spam protection on all form submissions (non-negotiable requirement)
- **Rate limiting**: IP-based throttling prevents spam floods and automated submissions
- **Input validation**: Server-side email format validation prevents malformed data storage

**Privacy:**
- **Minimal data collection**: Only email, source, and timestamp stored (no tracking cookies, no analytics, no PII beyond email)
- **Data sovereignty**: All data stored in self-hosted database under Bastien's control
- **No third-party data sharing**: Zero external services receive submission data (except Cloudflare Turnstile for verification)

### Reliability

**Uptime:**
- **Target uptime**: 100% availability after initial deployment
- **Self-hosted stability**: System runs reliably on self-hosted infrastructure without external dependencies (except Cloudflare Turnstile)
- **Health checks**: Container health check endpoints for monitoring system status

**Zero Maintenance:**
- **Automated operations**: Database initialization automatic on first run
- **Stateless frontend**: Static HTML served via nginx requires no maintenance
- **Minimal backend complexity**: Simple REST API with minimal dependencies reduces maintenance surface area

**Error Handling:**
- **Graceful degradation**: If Cloudflare Turnstile is unavailable, fail closed with clear error message (don't accept unverified submissions)
- **Input validation errors**: Clear, user-friendly error messages for invalid email formats
- **API failures**: Proper HTTP status codes and error responses for all failure scenarios
- **Startup validation**: Environment variable validation on startup with fail-fast behavior and clear error messages

**Data Integrity:**
- **Database constraints**: Email field validation at database level
- **Transaction safety**: SQLite ACID compliance ensures data integrity on writes
- **No data loss**: All successful submissions persisted durably to disk

### Accessibility

**WCAG 2.1 AA Compliance:**
- **Full accessibility support**: Given form simplicity (single email input + submit button), full WCAG 2.1 AA compliance is achievable and required

**Keyboard Navigation:**
- **Full keyboard accessibility**: Complete form interaction via keyboard only (Tab navigation: email field → submit button)
- **Enter key submission**: Form submits when Enter pressed from email input field
- **Focus indicators**: Clear, visible focus indicators on all interactive elements (minimum 3:1 contrast ratio)
- **No keyboard traps**: Users can navigate into and out of iframe without getting stuck

**Screen Reader Support:**
- **Semantic HTML**: Proper use of `<form>`, `<label>`, `<input type="email">`, `<button>` elements
- **ARIA labels**: Appropriate ARIA attributes for clarity where needed
- **Error announcements**: Form validation errors announced to screen readers via `aria-live` regions
- **Success announcements**: Success messages announced dynamically to screen readers
- **Turnstile accessibility**: Cloudflare Turnstile widget accessible (handled by Cloudflare)

**Visual Accessibility:**
- **Color contrast**: Text color contrast ratio ≥ 4.5:1 for normal text (WCAG AA standard)
- **Focus contrast**: Focus indicators meet ≥ 3:1 contrast requirement
- **Text resizing**: Text resizable up to 200% without loss of functionality or content
- **No color-only information**: Form states (normal, error, success) communicated via text + icons, not color alone

**Form-Specific Accessibility:**
- **Explicit labels**: `<label>` element explicitly associated with email `<input>` via for/id attributes
- **Error association**: Validation error messages associated with input field via `aria-describedby`
- **Descriptive buttons**: Submit button has clear, descriptive text ("Subscribe" or "Submit")
- **Touch targets**: Minimum 44px tap targets on mobile for touch accessibility

**Testing Coverage:**
- **Automated**: axe DevTools and Lighthouse accessibility audits pass with no violations
- **Manual keyboard**: Full keyboard-only navigation testing
- **Screen readers**: Testing with NVDA (Windows) and VoiceOver (macOS)
