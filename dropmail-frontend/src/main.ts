import './style.css';
import { renderTurnstile, resetTurnstile } from './turnstile';
import { getSource } from './source';

interface SubmitRequest {
  email: string;
  source: string;
  turnstile_token: string;
}

interface SubmitResponse {
  success: boolean;
  message: string;
}

interface ErrorResponse {
  error: string;
}

// Validate email format
function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

// DOM Elements
const form = document.getElementById('subscribe-form') as HTMLFormElement;
const emailInput = document.getElementById('email') as HTMLInputElement;
const submitBtn = document.getElementById('submit-btn') as HTMLButtonElement;
const errorMessage = document.getElementById('email-error') as HTMLDivElement;
const successMessage = document.getElementById('success-message') as HTMLDivElement;

// State
let isSubmitting = false;

// Show error message with warning icon (not color alone per WCAG)
function showError(message: string): void {
  errorMessage.textContent = `⚠ ${message}`;
  emailInput.classList.add('invalid');
  emailInput.setAttribute('aria-invalid', 'true');
}

// Clear error message
function clearError(): void {
  errorMessage.textContent = '';
  emailInput.classList.remove('invalid');
  emailInput.removeAttribute('aria-invalid');
}

// Show success state
function showSuccess(): void {
  form.hidden = true;
  successMessage.hidden = false;
}

// Reset form to initial state
export function resetForm(): void {
  form.hidden = false;
  successMessage.hidden = true;
  emailInput.value = '';
  clearError();
  setSubmitting(false);
  resetTurnstile();
}

// Set submitting state with accessibility announcements
function setSubmitting(submitting: boolean): void {
  isSubmitting = submitting;
  submitBtn.disabled = submitting;
  submitBtn.textContent = submitting ? 'Submitting...' : 'Subscribe';

  // Set aria-busy on form for screen readers
  form.setAttribute('aria-busy', submitting ? 'true' : 'false');

  // Announce loading state to screen readers
  if (submitting) {
    submitBtn.setAttribute('aria-label', 'Submitting your email, please wait');
  } else {
    submitBtn.removeAttribute('aria-label');
  }
}

// Submit the form to the API with token
async function submitFormWithToken(token: string): Promise<void> {
  const email = emailInput.value.trim();

  const request: SubmitRequest = {
    email,
    source: getSource(),
    turnstile_token: token,
  };

  try {
    const response = await fetch('/api/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    const data = await response.json();

    if (!response.ok) {
      const errorData = data as ErrorResponse;
      showError(errorData.error || 'Submission failed');
      setSubmitting(false);
      resetTurnstile();
      return;
    }

    const successData = data as SubmitResponse;
    if (successData.success) {
      setSubmitting(false); // Reset aria-busy before showing success
      showSuccess();
    } else {
      showError('Submission failed');
      setSubmitting(false);
      resetTurnstile();
    }
  } catch {
    // Network errors are shown to user; no console logging in production
    showError('Network error. Please try again.');
    setSubmitting(false);
    resetTurnstile();
  }
}

// Handle form submission - validate then trigger Turnstile
async function handleSubmit(): Promise<void> {
  const email = emailInput.value.trim();

  // Validate email
  if (!email) {
    showError('Email is required');
    emailInput.focus();
    return;
  }

  if (!isValidEmail(email)) {
    showError('Please enter a valid email address');
    emailInput.focus();
    return;
  }

  setSubmitting(true);
  clearError();

  // Render Turnstile and wait for verification
  renderTurnstile(
    // Success callback - got token, submit form
    (token: string) => {
      submitFormWithToken(token);
    },
    // Error callback - verification failed
    () => {
      showError('Verification failed. Please try again.');
      setSubmitting(false);
    }
  );
}

// Event Listeners
form.addEventListener('submit', (e: Event) => {
  e.preventDefault();
  if (!isSubmitting) {
    handleSubmit();
  }
});

// Validate on blur
emailInput.addEventListener('blur', () => {
  const email = emailInput.value.trim();
  if (email && !isValidEmail(email)) {
    showError('Please enter a valid email address');
  } else {
    clearError();
  }
});

// Clear error on input
emailInput.addEventListener('input', () => {
  if (errorMessage.textContent) {
    clearError();
  }
});

// iframe embedding support: transparent background mode
// Add ?transparent=1 to URL for transparent background
function initializeEmbedMode(): void {
  const params = new URLSearchParams(window.location.search);
  if (params.get('transparent') === '1') {
    document.documentElement.style.backgroundColor = 'transparent';
    document.body.style.backgroundColor = 'transparent';
  }
}

initializeEmbedMode();
