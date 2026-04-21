import { useState, useRef } from 'react';
import type { Tag } from './lib/api';

export interface TagChipProps {
  tag: Tag;
  /** If true, render only value (omit "purpose:") */
  noPurpose?: boolean;
  /** If set, chip shows an × button that calls this */
  onRemove?: () => void;
  /** If set, chip is color-editable; emits #RRGGBBAA */
  onColorChange?: (color: string | null) => void;
}

/**
 * Compute a simple black-or-white text color based on perceived luminance.
 * Accepts #RRGGBB or #RRGGBBAA hex strings.
 */
function contrastColor(hexColor: string): string {
  const hex = hexColor.replace('#', '').slice(0, 6); // use RGB portion only
  if (hex.length < 6) return '#000000';
  const r = parseInt(hex.slice(0, 2), 16);
  const g = parseInt(hex.slice(2, 4), 16);
  const b = parseInt(hex.slice(4, 6), 16);
  // Perceived luminance (ITU-R BT.709)
  const luminance = (0.2126 * r + 0.7152 * g + 0.0722 * b) / 255;
  return luminance > 0.5 ? '#000000' : '#ffffff';
}

/**
 * Parse a #RRGGBBAA string into its RGB (#RRGGBB) and alpha (0–255) parts.
 * Falls back gracefully for #RRGGBB input.
 */
function parseColor(hex: string): { rgb: string; alpha: number } {
  const clean = hex.replace('#', '');
  const rgb = '#' + clean.slice(0, 6).padEnd(6, '0');
  const alphaHex = clean.slice(6, 8);
  const alpha = alphaHex.length === 2 ? parseInt(alphaHex, 16) : 255;
  return { rgb, alpha };
}

/** Compose a #RRGGBBAA string from an #RRGGBB color and an alpha 0–255. */
function composeColor(rgb: string, alpha: number): string {
  const alphaHex = Math.round(alpha).toString(16).padStart(2, '0');
  return rgb + alphaHex;
}

const NEUTRAL_BG = '#e5e7eb'; // tailwind gray-200 equivalent

export function TagChip({ tag, noPurpose, onRemove, onColorChange }: TagChipProps) {
  const bgColor = tag.color ?? NEUTRAL_BG;
  const textColor = contrastColor(bgColor);
  const label = noPurpose ? tag.value : `${tag.purpose}:${tag.value}`;

  const [showColorEditor, setShowColorEditor] = useState(false);
  const { rgb: initialRgb, alpha: initialAlpha } = parseColor(bgColor);
  const [pickerRgb, setPickerRgb] = useState(initialRgb);
  const [pickerAlpha, setPickerAlpha] = useState(initialAlpha);
  const colorEditorRef = useRef<HTMLDivElement>(null);

  function handleColorApply() {
    const composed = composeColor(pickerRgb, pickerAlpha);
    onColorChange?.(composed);
    setShowColorEditor(false);
  }

  function handleColorClear() {
    onColorChange?.(null);
    setShowColorEditor(false);
  }

  return (
    <span className="relative inline-flex items-center gap-1">
      {/* Main chip */}
      <span
        className="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium"
        style={{ backgroundColor: bgColor, color: textColor }}
      >
        <span
          className={onColorChange ? 'cursor-pointer' : undefined}
          onClick={onColorChange ? () => setShowColorEditor((v) => !v) : undefined}
          role={onColorChange ? 'button' : undefined}
          aria-label={onColorChange ? `Edit color for ${label}` : undefined}
        >
          {label}
        </span>

        {onRemove && (
          <button
            type="button"
            className="ml-0.5 inline-flex size-3.5 items-center justify-center rounded-full opacity-60 hover:opacity-100"
            style={{ color: textColor }}
            onClick={(e) => {
              e.stopPropagation();
              onRemove();
            }}
            aria-label={`Remove tag ${label}`}
          >
            ×
          </button>
        )}
      </span>

      {/* Inline color popover */}
      {onColorChange && showColorEditor && (
        <div
          ref={colorEditorRef}
          className="absolute left-0 top-full z-10 mt-1 flex flex-col gap-2 rounded-md border bg-white p-3 shadow-md"
          style={{ minWidth: '180px' }}
        >
          <div className="flex items-center gap-2">
            <label className="text-xs text-gray-600" htmlFor={`color-rgb-${tag.uuid}`}>
              Color
            </label>
            <input
              id={`color-rgb-${tag.uuid}`}
              type="color"
              value={pickerRgb}
              onChange={(e) => setPickerRgb(e.target.value)}
              className="h-6 w-10 cursor-pointer rounded border-0 p-0"
            />
          </div>
          <div className="flex items-center gap-2">
            <label className="text-xs text-gray-600" htmlFor={`color-alpha-${tag.uuid}`}>
              Alpha
            </label>
            <input
              id={`color-alpha-${tag.uuid}`}
              type="range"
              min={0}
              max={255}
              value={pickerAlpha}
              onChange={(e) => setPickerAlpha(Number(e.target.value))}
              className="flex-1"
            />
            <span className="w-6 text-right text-xs text-gray-500">{pickerAlpha}</span>
          </div>
          <div className="flex gap-1.5">
            <button
              type="button"
              className="flex-1 rounded bg-gray-900 px-2 py-1 text-xs text-white hover:bg-gray-700"
              onClick={handleColorApply}
            >
              Apply
            </button>
            <button
              type="button"
              className="flex-1 rounded border px-2 py-1 text-xs text-gray-600 hover:bg-gray-50"
              onClick={handleColorClear}
            >
              Clear
            </button>
            <button
              type="button"
              className="flex-1 rounded border px-2 py-1 text-xs text-gray-600 hover:bg-gray-50"
              onClick={() => setShowColorEditor(false)}
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </span>
  );
}
