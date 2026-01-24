/**
 * Source parameter extraction and validation
 * Captures the source of form submissions from URL query parameters
 */

const DEFAULT_SOURCE = 'direct';
const MAX_SOURCE_LENGTH = 50;

// Allowed characters: alphanumeric, hyphens, underscores
const SOURCE_PATTERN = /^[a-zA-Z0-9_-]+$/;

/**
 * Sanitize source parameter to prevent injection
 * Only allows alphanumeric characters, hyphens, and underscores
 */
function sanitizeSource(source: string): string {
  // Trim whitespace
  const trimmed = source.trim();

  // Check length
  if (trimmed.length === 0 || trimmed.length > MAX_SOURCE_LENGTH) {
    return DEFAULT_SOURCE;
  }

  // Validate characters
  if (!SOURCE_PATTERN.test(trimmed)) {
    return DEFAULT_SOURCE;
  }

  return trimmed;
}

/**
 * Extract and validate source from URL query parameters
 * @returns Sanitized source string or 'direct' as default
 */
export function getSource(): string {
  const params = new URLSearchParams(window.location.search);
  const source = params.get('source');

  if (!source) {
    return DEFAULT_SOURCE;
  }

  return sanitizeSource(source);
}
