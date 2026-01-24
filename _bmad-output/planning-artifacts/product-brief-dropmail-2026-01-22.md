---
stepsCompleted: [1, 2, 3, 4, 5]
inputDocuments: []
date: 2026-01-22
author: Bastien
---

# Product Brief: dropmail

<!-- Content will be appended sequentially through collaborative workflow steps -->

## Executive Summary

dropmail is a self-hosted email collection tool built for developers and creators who want data sovereignty without feature bloat. It provides an ultra-minimalist embeddable form (email field + submit button) paired with a decoupled backend API, enabling clean integration across multiple properties (Notion pages, portfolios, project sites) while maintaining full control over collected data. Unlike existing solutions that assume email marketing workflows, dropmail focuses solely on interest tracking and validation metrics for personal projects.

---

## Core Vision

### Problem Statement

Existing email collection tools suffer from fundamental misalignment: they're built for email marketing campaigns when many creators simply need a lightweight collection endpoint. Current solutions force users to accept branded free tiers with submission limits, give up data sovereignty to third-party platforms, or deploy heavyweight self-hosted marketing suites with features they'll never use. For technical creators building personal brands, using vendor-branded forms contradicts the "I build my own tools" identity they're showcasing.

### Problem Impact

**For Technical Creators:**
- **Brand Misalignment**: Vendor-branded forms undermine the "I build solutions" personal narrative
- **Data Sovereignty Loss**: Audience data locked in third-party platforms subject to terms changes and limits
- **Feature Bloat**: Paying (with money or brand exposure) for campaign management, A/B testing, and automation features that go unused
- **Artificial Constraints**: Free tier submission limits that cap growth arbitrarily
- **Aesthetic Pollution**: Cluttered, over-designed forms that clash with minimalist Notion pages and clean portfolio sites

**Current Workarounds:**
- Accepting branded footers on free tools (Tally, Typeform, Google Forms)
- Building quick-and-dirty Express endpoints that aren't portfolio-grade
- Over-provisioning with Listmonk/Mailtrain for simple collection needs
- Settling for ugly embeds that hurt user experience

### Why Existing Solutions Fall Short

**SaaS Tools (Tally, Typeform, Mailchimp, ConvertKit):**
- Built for marketing campaigns, not simple collection
- Monetize through feature upsells users don't need
- Free tiers impose submission limits and mandatory branding
- Data lives on vendor infrastructure with no export guarantees

**Google Forms:**
- Ugly, un-customizable embeds
- Google branding and data collection policies
- No control over presentation or data storage

**Self-Hosted Marketing Platforms (Listmonk, Mailtrain):**
- Designed for campaign management, not lightweight collection
- Heavy infrastructure requirements for unused features
- Overkill complexity for a simple use case

**DIY Solutions:**
- Quick Express + SQLite scripts work but aren't reusable or portfolio-worthy
- Lack anti-spam hardening for public-facing forms
- No decoupling for multi-property reuse

### Proposed Solution

dropmail provides a purpose-built solution for email collection without marketing bloat:

**Core Components:**
1. **Ultra-Minimalist Embeddable Form**: Lightweight web app containing only an email input field and submit button - nothing more. Designed for clean embedding in Notion pages, portfolios, and project sites.

2. **Decoupled Backend API**: Independent submission handler that can serve multiple frontend instances across different properties. One backend, infinite forms.

3. **Self-Hosted Architecture**: Deploy on your own infrastructure for complete data sovereignty and zero vendor lock-in.

4. **Anti-Spam Hardening**: Cloudflare Turnstile integration to prevent DoS attacks and spam submissions without user friction.

5. **Simple Metrics**: Submission count tracking to gauge interest and validate project ideas.

**Design Philosophy:**
- Build exactly what's needed, nothing more
- Prioritize reusability across multiple projects
- Portfolio-grade code quality that showcases engineering capabilities
- Privacy-first, no tracking, no bloat

### Key Differentiators

**1. Purpose-Built Minimalism**
Unlike marketing platforms retrofitted for collection, dropmail is designed solely for gathering interest signals. No campaign builders, no automation workflows, no feature creep.

**2. Decoupled Multi-Property Architecture**
One backend serves forms across Notion, portfolios, and future projects. Competitors lock you into single-use embeds or require separate accounts per property.

**3. Portfolio Showcase Value**
dropmail itself demonstrates technical capability: clean architecture, security hardening, reusable design. It's both a tool AND a portfolio piece.

**4. Zero Vendor Dependencies**
Self-hosted by design means no free tier limits, no branding requirements, no terms-of-service changes, no surprise shutdowns.

**5. Developer-First Philosophy**
Built by a software engineer for software engineers who value control, simplicity, and owning their infrastructure.

**Unfair Advantage**: Existing players can't pivot to this simplicity - their business models depend on upselling marketing features. dropmail's self-hosted model has no monetization pressure driving feature bloat.

## Target Users

### Primary User

**Bastien - Technical Creator & Solo Developer**

**Role & Context:**
- Software engineer building and showcasing personal projects
- Maintains presence across multiple properties: Notion pages, portfolio site, project showcases
- Values data sovereignty, clean code, and building custom tools over using third-party services
- Self-hosts infrastructure on owned hardware

**Current State:**
- Has no mechanism to track interest in work
- Needs lightweight email collection that works across all properties
- Wants portfolio-grade code that demonstrates technical capability
- Requires solution that aligns with "I build my own tools" personal brand

**Success Criteria:**
- Single backend instance serving forms across all web properties
- Clean email entries in database with minimal friction
- Direct database access for viewing/exporting (no dashboard complexity)
- Spam-protected without user experience friction
- Code quality worthy of portfolio showcase

**Usage Pattern:**
- Deploy dropmail once on self-hosted infrastructure
- Embed minimalist forms in Notion pages, portfolio site, and future properties
- Periodically query database to see interest levels and export emails when needed
- Use submission counts as validation metric for project ideas

### Secondary Users (Form Visitors)

**Profile:**
- People who discover Bastien's work through various channels
- Interested in staying updated on future projects
- Value simplicity - just want to submit email and move on

**Interaction:**
- Encounter embedded form on Notion page, portfolio, or project site
- Enter email address in single input field
- Click submit button
- Receive confirmation (optional visual feedback)
- No account creation, no verification emails, no friction

**Success Moment:**
- Form submission completes instantly without redirects or unnecessary steps
- Clean, unbranded interface that doesn't distract from content

### User Journey

**Bastien's Journey (Primary User):**

1. **Setup (One-time)**
   - Deploy dropmail backend on self-hosted server
   - Configure Cloudflare Turnstile for spam protection
   - Database automatically initialized

2. **Integration (Per Property)**
   - Embed lightweight form HTML/iframe in Notion page
   - Add same form to portfolio site
   - Future properties reuse same backend endpoint

3. **Ongoing Usage**
   - Check submission counts via direct database queries
   - Export email list when ready to send manual updates
   - Monitor for spam/abuse patterns (minimal due to Turnstile)

4. **Success Realization**
   - Form works seamlessly across all properties without modification
   - No vendor branding or limitations
   - Clean codebase serves as portfolio demonstration piece
   - Data sovereignty maintained with zero third-party dependencies

**Visitor Journey (Secondary Users):**

1. **Discovery**
   - Find Bastien's Notion page, portfolio, or project site
   - See clean, minimal email collection form

2. **Decision**
   - Decide to stay updated on future work
   - No concerns about spam (reputable creator, simple ask)

3. **Submission**
   - Enter email in single field
   - Complete Turnstile verification (invisible/minimal)
   - Click submit

4. **Confirmation**
   - Receive instant feedback (success message)
   - Return to browsing content immediately

**Anti-Journey (What We're Avoiding):**
- No lengthy forms with name/company/phone fields
- No redirect to thank-you pages that break flow
- No "Powered by [Vendor]" footer diluting brand
- No vendor registration or account setup for Bastien
- No dashboard login just to see email list

## Success Metrics

### Primary User Success (Bastien)

**Deployment Success:**
- Deploy once, works everywhere: Forms embed cleanly across Notion, portfolio, and future sites without modification
- Clean data collection: Email submissions land in database correctly formatted with no corruption
- Spam protection: Turnstile blocks bots while legitimate submissions get through (target: >95% legitimate)
- Portfolio quality: Codebase is clean, well-architected, and demonstrates technical capability
- Zero maintenance: Once deployed, system runs reliably without intervention

**Operational Success:**
- Form loads and submits successfully (target: 100% uptime after deployment)
- Email validation works: Invalid emails rejected at submission time
- Email submissions accumulate over time (validation that people care about the work)
- Direct database access for viewing/exporting works without additional tools

**Success Realization Moment:**
- Form works seamlessly across all embedded properties without modification
- No vendor branding or limitations
- Clean codebase serves as portfolio demonstration piece
- Data sovereignty maintained with zero third-party dependencies

### Secondary User Success (Visitors)

**Submission Experience:**
- Fast submission: Enter email, click submit, complete in <5 seconds
- Clear feedback: Immediate confirmation that submission worked
- No friction: No complex captchas, no redirects, no account creation requirements
- Clean interface: Unbranded form that doesn't distract from content

**Success Indicator:**
- Form submission completes instantly without errors or unnecessary steps

### Anti-Goals (What We're NOT Optimizing For)

- **NOT** maximizing conversion rates with A/B testing or optimization
- **NOT** adding features or complexity beyond core email collection
- **NOT** measuring engagement metrics beyond submission count
- **NOT** building dashboards, analytics, or complex reporting
- **NOT** adding email verification workflows or double opt-in

**Success Definition:** dropmail stays simple, works reliably, and requires zero ongoing attention.

## MVP Scope

### Core Features

**Frontend - Embeddable Form:**
- Single HTML page with email input field and submit button
- Client-side email validation (format check before submission)
- Clean, minimal CSS suitable for embedding in Notion, portfolios, etc.
- Simple success message on successful submission
- Simple error message on failed submission
- Responsive design that works across devices

**Backend - API:**
- REST API endpoint to receive form submissions
- Server-side email validation (format verification)
- Cloudflare Turnstile integration for spam protection
- CORS support for cross-origin form embeds
- Rate limiting to prevent abuse
- Environment variable configuration (Turnstile keys, DB credentials)

**Database:**
- Simplest viable database (SQLite or PostgreSQL)
- Schema: email address (primary), form source (tracking), submission timestamp
- Direct database access for viewing/exporting (no admin UI needed)

**Deployment:**
- Docker Compose configuration for portable deployment
- Single-command deployment to self-hosted infrastructure
- Environment variable management for secrets
- Automatic database initialization on first run

**Security:**
- HTTPS enforcement
- Environment variables for sensitive configuration
- Input sanitization to prevent injection attacks
- Turnstile verification on all submissions

### Out of Scope for MVP

This is a complete, focused tool - the MVP is the final product. The following are intentionally excluded:

**Never Planned:**
- Email verification/double opt-in workflows
- Admin dashboard or web UI for viewing submissions
- Email campaign or sending functionality
- User accounts or authentication
- A/B testing or analytics beyond basic counts
- Multi-user support or team features
- CSV export endpoints (direct DB access is sufficient)
- Custom form field configuration
- Branding or white-labeling options

**Complexity Avoided:**
- No microservices architecture - single backend service
- No elaborate frontend framework - vanilla HTML/CSS/JS
- No complex deployment orchestration - Docker Compose only

### MVP Success Criteria

Since MVP = final product, success means:

**Functional Completeness:**
- Forms embed cleanly in Notion, portfolio, and future sites without modification
- Email submissions store correctly with source tracking and timestamps
- Turnstile blocks spam effectively (>95% legitimate submissions)
- Docker Compose deployment works on first try

**Operational Readiness:**
- Zero configuration after initial environment setup
- No maintenance required after deployment
- Direct database access works for viewing/exporting emails
- System runs reliably without intervention

**Portfolio Quality:**
- Clean, well-documented codebase
- Proper separation of concerns (frontend/backend/data)
- Security best practices implemented
- Demonstrates technical capability to employers/clients

**Decision Point:** Once deployed and collecting emails successfully, the project is complete. Future work is maintenance-only (security updates, dependency updates).

### Future Vision

**Maintenance Mode Only:**

This is a complete tool, not a growth product. Future work consists only of:

**Security Maintenance:**
- Dependency updates for security patches
- Turnstile API updates if Cloudflare changes their service

**Optional Quality-of-Life Improvements:**
- Simple CSV export endpoint (if direct DB access becomes inconvenient)
- Email deduplication (if duplicate submissions become an issue)
- Basic submission statistics endpoint (if manual DB queries are tedious)

**Explicitly NOT Planned:**
- Feature expansion
- Multi-tenancy or SaaS offering
- Marketing or monetization
- Community or open-source development

**Philosophy:** dropmail stays frozen at "simple email collection" forever. Any temptation to add features is scope creep, not product evolution.
