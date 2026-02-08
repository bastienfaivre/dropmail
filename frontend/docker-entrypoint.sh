#!/bin/sh
# Generate runtime env config from container environment variables.
# This runs at container startup, before nginx serves the app.

cat <<EOF > /usr/share/nginx/html/env.js
window.__ENV__ = {
  VITE_API_ENDPOINT: "${VITE_API_ENDPOINT}",
  VITE_TURNSTILE_SITE_KEY: "${VITE_TURNSTILE_SITE_KEY}"
};
EOF

exec "$@"
