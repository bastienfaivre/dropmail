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

const SITE_KEY = import.meta.env.VITE_TURNSTILE_SITE_KEY || '1x00000000000000000000BB';

function isTurnstileLoaded(): boolean {
  return typeof window.turnstile !== 'undefined';
}

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

export async function renderTurnstile(
  onSuccess: TurnstileCallback,
  onError?: TurnstileErrorCallback
): Promise<void> {
  try {
    await waitForTurnstile();
  } catch {
    if (onError) {
      onError();
    }
    return;
  }

  if (widgetId !== null) {
    try {
      window.turnstile.remove(widgetId);
    } catch { }
    widgetId = null;
  }

  currentCallback = onSuccess;
  currentErrorCallback = onError || null;

  const container = document.getElementById('turnstile-container');
  if (!container) {
    if (onError) {
      onError();
    }
    return;
  }

  container.style.display = '';

  try {
    widgetId = window.turnstile.render(container, {
      sitekey: SITE_KEY,
      callback: (token: string) => {
        if (currentCallback) {
          currentCallback(token);
        }
      },
      'error-callback': () => {
        if (currentErrorCallback) {
          currentErrorCallback();
        }
      },
      'expired-callback': () => { },
      theme: 'auto',
      size: 'normal',
    });
  } catch {
    if (onError) {
      onError();
    }
  }
}

export function removeTurnstile(): void {
  if (widgetId !== null && isTurnstileLoaded()) {
    try {
      window.turnstile.remove(widgetId);
    } catch { }
    widgetId = null;
  }

  const container = document.getElementById('turnstile-container');
  if (container) {
    container.innerHTML = '';
    container.style.display = 'none';
  }

  currentCallback = null;
  currentErrorCallback = null;
}
