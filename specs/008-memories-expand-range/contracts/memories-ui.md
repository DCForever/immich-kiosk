# Memories UI Contract (008 Extensions)

**Feature**: 008-memories-expand-range  
**Date**: 2026-03-13  
**Extends**: 007-memories-feature contracts

## Overview

This document defines the contract changes for the memories feature when progressive time-range expansion and 3–16 photos are used. The base contract remains in `specs/007-memories-feature/contracts/memories-ui.md`.

## Changes from 007

### Collage Container

- **Asset count**: 3–16 photos (was 1–8). When fewer than 3 exist in any range, show empty state.
- **Layout class**: `collage-{n}-{variant}` where n = 3..16 and variant = a, b, c, ... (≥3 variants per n).
- **Layout constraints**: Collage layouts MUST have at least 2 columns and 2 rows (no 1×N vertical stack or N×1 horizontal strip). Layouts MUST have varied proportions: either varied column widths, varied row heights, or spanning cells creating visual hierarchy (no uniform grids where every cell is equal).

```html
<div class="collage-grid collage-{n}-{variant}">
  <!-- n = 3..16, number of assets; variant = a, b, c, ... (≥3 per n) -->
  ...
</div>
```

### Caption (data-memory-caption)

Caption varies by `TimeRange`:

| TimeRange | data-memory-caption / i18n |
|-----------|----------------------------|
| ExactDate | "Memories from X years ago" |
| Week | "Memories from this week, X years ago" |
| Month | "Memories from [month], X years ago" |

- **New i18n keys**: `memories_caption_week`, `memories_caption_month` (or equivalent with interpolation).
- **Existing**: `memories_caption` for exact date (X years ago).

### kiosk-data (JSON)

No change. `duration: 45`, `memories: true` remain. Optionally add `timeRange: "exact"|"week"|"month"` if frontend needs it; otherwise caption is pre-rendered.

### Endpoints

No new endpoints. `GET/POST /asset/new` and `GET/POST /memories` behavior unchanged; response includes 3–16 assets and appropriate caption based on `TimeRange`.

## Accessibility

Unchanged from 007. Collage cells: `role="button"`, `tabindex="0"`. Lightbox: focus trap, Escape to close.
