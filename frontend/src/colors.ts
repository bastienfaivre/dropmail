function hexToRgb(hex: string): string | null {
  hex = hex.replace(/^#/, '');

  if (hex.length === 3) {
    hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];
  }

  if (hex.length !== 6) {
    return null;
  }

  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);

  if (isNaN(r) || isNaN(g) || isNaN(b)) {
    return null;
  }

  return `${r}, ${g}, ${b}`;
}

function getPerceivedBrightness(r: number, g: number, b: number): number {
  return (r * 299 + g * 587 + b * 114) / 1000;
}

function adjustColorBrightness(hex: string): string {
  hex = hex.replace(/^#/, '');

  if (hex.length === 3) {
    hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];
  }

  let r = parseInt(hex.substring(0, 2), 16);
  let g = parseInt(hex.substring(2, 4), 16);
  let b = parseInt(hex.substring(4, 6), 16);

  const brightness = getPerceivedBrightness(r, g, b);
  const adjustment = brightness < 128 ? 50 : -50;

  r = Math.min(255, Math.max(0, r + adjustment));
  g = Math.min(255, Math.max(0, g + adjustment));
  b = Math.min(255, Math.max(0, b + adjustment));

  const rr = r.toString(16).padStart(2, '0');
  const gg = g.toString(16).padStart(2, '0');
  const bb = b.toString(16).padStart(2, '0');

  return `#${rr}${gg}${bb}`;
}

export function initializeColors(): void {
  const params = new URLSearchParams(window.location.search);
  const primary = params.get('primary');
  const secondary = params.get('secondary');

  if (primary) {
    const primaryHex = primary.startsWith('#') ? primary : `#${primary}`;
    const primaryRgb = hexToRgb(primaryHex);

    if (primaryRgb) {
      document.documentElement.style.setProperty('--primary-color', primaryHex);
      document.documentElement.style.setProperty('--primary-rgb', primaryRgb);

      if (secondary) {
        const secondaryHex = secondary.startsWith('#') ? secondary : `#${secondary}`;
        document.documentElement.style.setProperty('--primary-hover', secondaryHex);
      } else {
        const hoverColor = adjustColorBrightness(primaryHex);
        document.documentElement.style.setProperty('--primary-hover', hoverColor);
      }
    }
  }
}
