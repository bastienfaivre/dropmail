# Dropmail

Backend API for collecting email submissions with Cloudflare Turnstile protection, backed by SQLite.

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | No | `8080` | HTTP server port |
| `DB_PATH` | No | `/data/dropmail.db` | Path to the SQLite database file |
| `TURNSTILE_SECRET` | No | `1x00...AA` (test key) | Cloudflare Turnstile secret key |
| `CORS_ORIGINS` | No | `http://localhost:3000/` | Comma-separated list of allowed CORS origins |
| `VALID_SOURCES` | No | _(none)_ | Comma-separated list of allowed submission sources |
| `RESEND_API_KEY` | No | _(none)_ | Resend API key for sending emails |
| `RESEND_FROM` | No | _(none)_ | Sender address (e.g. `Dropmail <noreply@yourdomain.com>`) |
| `SUMMARY_EMAIL` | No | _(none)_ | Recipient email for the daily summary |
| `SUMMARY_SCHEDULE` | No | `08:00` | Time to send the daily summary (HH:MM, UTC) |

The daily summary email is opt-in. It only runs when `RESEND_API_KEY`, `RESEND_FROM`, and `SUMMARY_EMAIL` are all set.
