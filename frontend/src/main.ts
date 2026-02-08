import './style.css';
import { renderTurnstile, removeTurnstile } from './turnstile';
import { getSource } from './source';
import { initializeColors } from './colors';

function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

const form = document.getElementById('subscribe-form') as HTMLFormElement;
const emailInput = document.getElementById('email') as HTMLInputElement;
const inputWrapper = document.querySelector('.input-wrapper') as HTMLDivElement;
const submitBtn = document.getElementById('submit-btn') as HTMLButtonElement;
const errorMessage = document.getElementById('email-error') as HTMLDivElement;
const successMessage = document.getElementById('success-message') as HTMLDivElement;

let isSubmitting = false;

function showError(message: string): void {
  errorMessage.textContent = `⚠ ${message}`;
  inputWrapper.classList.add('invalid');
  emailInput.setAttribute('aria-invalid', 'true');
}

function clearError(): void {
  errorMessage.textContent = '';
  inputWrapper.classList.remove('invalid');
  emailInput.removeAttribute('aria-invalid');
}

function showSuccess(): void {
  form.hidden = true;
  successMessage.hidden = false;
}

function setSubmitting(submitting: boolean): void {
  isSubmitting = submitting;
  submitBtn.disabled = submitting;

  if (submitting) {
    submitBtn.innerHTML = '<span class="spinner"></span>';
  } else {
    submitBtn.textContent = 'Subscribe';
  }

  form.setAttribute('aria-busy', submitting ? 'true' : 'false');

  if (submitting) {
    submitBtn.setAttribute('aria-label', 'Submitting your email, please wait');
  } else {
    submitBtn.removeAttribute('aria-label');
  }
}

async function submitFormWithToken(token: string): Promise<void> {
  const email = emailInput.value.trim();

  try {
    const response = await fetch(window.__ENV__?.VITE_API_ENDPOINT || import.meta.env.VITE_API_ENDPOINT || '/api/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        source: getSource(),
        turnstile_token: token,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      showError(data.error || 'Submission failed');
      setSubmitting(false);
      removeTurnstile();
      return;
    }

    if (data.success) {
      setSubmitting(false);
      showSuccess();
    } else {
      showError('Submission failed');
      setSubmitting(false);
      removeTurnstile();
    }
  } catch {
    showError('Network error. Please try again.');
    setSubmitting(false);
    removeTurnstile();
  }
}

async function handleSubmit(): Promise<void> {
  const email = emailInput.value.trim();

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

  renderTurnstile(
    (token: string) => {
      submitFormWithToken(token);
    },
    () => {
      showError('Verification failed. Please try again.');
      removeTurnstile();
      setSubmitting(false);
    }
  );
}

form.addEventListener('submit', (e: Event) => {
  e.preventDefault();
  if (!isSubmitting) {
    handleSubmit();
  }
});

emailInput.addEventListener('blur', () => {
  const email = emailInput.value.trim();
  if (email && !isValidEmail(email)) {
    showError('Please enter a valid email address');
  } else {
    clearError();
  }
});

emailInput.addEventListener('input', () => {
  if (errorMessage.textContent) {
    clearError();
  }
});

function initializeEmbedMode(): void {
  const params = new URLSearchParams(window.location.search);
  if (params.get('transparent') === '1') {
    document.documentElement.style.backgroundColor = 'transparent';
    document.body.style.backgroundColor = 'transparent';
  }
}

initializeColors();
initializeEmbedMode();
