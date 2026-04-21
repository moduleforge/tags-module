/**
 * Parse a #RRGGBBAA (or #RRGGBB) hex string into its RGB and alpha parts.
 * Handles null/undefined input by returning a neutral gray at full opacity.
 */
export function parseColor(color: string | null | undefined): { rgb: string; alpha: number } {
  const clean = (color ?? '').replace('#', '');
  const rgb = '#' + clean.slice(0, 6).padEnd(6, '0');
  const alphaHex = clean.slice(6, 8);
  const alpha = alphaHex.length === 2 ? parseInt(alphaHex, 16) : 255;
  return { rgb, alpha };
}

/**
 * Combine an #RRGGBB hex color and an alpha value (0–255) into #RRGGBBAA.
 */
export function composeColor(rgbHex: string, alpha: number): string {
  return rgbHex + Math.round(alpha).toString(16).padStart(2, '0');
}
