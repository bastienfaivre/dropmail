// Cloudflare Turnstile integration

declare global {
  interface Window {
    turnstile: {
      render: (
        container: string | HTMLElement,
        options: TurnstileOptions
      ) => string;
      reset: (widgetId: string) => void;
      remove: (widgetId: string) => void;
    };
  }
}

interface TurnstileOptions {
  sitekey: string;
  callback: (token: string) => void;
  'error-callback'?: () => void;
  'expired-callback'?: () => void;
  theme?: 'light' | 'dark' | 'auto';
  size?: 'normal' | 'compact' | 'invisible';
}

type TurnstileCallback = (token: string) => void;
type TurnstileErrorCallback = () => void;

let widgetId: string | null = null;
let currentCallback: TurnstileCallback | null = null;
let currentErrorCallback: TurnstileErrorCallback | null = null;

// Get the site key from environment or use a placeholder for development
const SITE_KEY = import.meta.env.VITE_TURNSTILE_SITE_KEY || '1x00000000000000000000AA'; // Test key

// Check if Turnstile is loaded
function isTurnstileLoaded(): boolean {
  return typeof window.turnstile !== 'undefined';
}

// Wait for Turnstile to load
function waitForTurnstile(timeout = 5000): Promise<void> {
  return new Promise((resolve, reject) => {
    if (isTurnstileLoaded()) {
      resolve();
      return;
    }

    const startTime = Date.now();
    const interval = setInterval(() => {
      if (isTurnstileLoaded()) {
        clearInterval(interval);
        resolve();
      } else if (Date.now() - startTime > timeout) {
        clearInterval(interval);
        reject(new Error('Turnstile failed to load'));
      }
    }, 100);
  });
}

// Render the Turnstile widget
export async function renderTurnstile(
  onSuccess: TurnstileCallback,
  onError?: TurnstileErrorCallback
): Promise<void> {
  try {
    await waitForTurnstile();
  } catch (error) {
    console.error('Turnstile not available:', error);
    if (onError) {
      onError();
    }
    return;
  }

  // Remove existing widget if any
  if (widgetId !== null) {
    try {
      window.turnstile.remove(widgetId);
    } catch (e) {
      // Ignore removal errors
    }
    widgetId = null;
  }

  currentCallback = onSuccess;
  currentErrorCallback = onError || null;

  const container = document.getElementById('turnstile-container');
  if (!container) {
    console.error('Turnstile container not found');
    if (onError) {
      onError();
    }
    return;
  }

  try {
    widgetId = window.turnstile.render(container, {
      sitekey: SITE_KEY,
      callback: (token: string) => {
        if (currentCallback) {
          currentCallback(token);
        }
      },
      'error-callback': () => {
        console.error('Turnstile verification failed');
        if (currentErrorCallback) {
          currentErrorCallback();
        }
      },
      'expired-callback': () => {
        console.warn('Turnstile token expired');
        // Token expired, but we'll get a new one when user retries
      },
      theme: 'auto',
      size: 'normal',
    });
  } catch (error) {
    console.error('Failed to render Turnstile:', error);
    if (onError) {
      onError();
    }
  }
}

// Reset the widget for retry
export function resetTurnstile(): void {
  if (widgetId !== null && isTurnstileLoaded()) {
    try {
      window.turnstile.reset(widgetId);
    } catch (e) {
      // Ignore reset errors
    }
  }
}
