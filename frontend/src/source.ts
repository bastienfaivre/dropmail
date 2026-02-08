const DEFAULT_SOURCE = 'unknown';
const MAX_SOURCE_LENGTH = 50;
const SOURCE_PATTERN = /^[a-zA-Z0-9_-]+$/;

function sanitizeSource(source: string): string {
  const trimmed = source.trim();

  if (trimmed.length === 0 || trimmed.length > MAX_SOURCE_LENGTH) {
    return DEFAULT_SOURCE;
  }

  if (!SOURCE_PATTERN.test(trimmed)) {
    return DEFAULT_SOURCE;
  }

  return trimmed;
}

export function getSource(): string {
  const params = new URLSearchParams(window.location.search);
  const source = params.get('source');

  if (!source) {
    return DEFAULT_SOURCE;
  }

  return sanitizeSource(source);
}
