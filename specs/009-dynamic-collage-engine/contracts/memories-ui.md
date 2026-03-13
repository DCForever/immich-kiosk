# Memories UI Contract (009 Extensions — Dynamic Collage)

**Feature**: 009-dynamic-collage-engine  
**Date**: 2026-03-13  
**Extends**: 007-memories-feature, 008-memories-expand-range contracts

## Overview

This document defines the contract changes when the dynamic collage engine is used. The base contract remains in `specs/007-memories-feature/contracts/memories-ui.md`. The 008 contract (`specs/008-memories-expand-range/contracts/memories-ui.md`) is superseded for layout structure; caption and endpoint behavior from 008 remain.

## Changes from 008

### Collage Container — Dynamic Layout

When the layout engine succeeds:
- **Layout class**: `collage-grid collage-dynamic` (or equivalent). No `collage-{n}-{variant}`.
- **Cell positioning**: Each `.collage-cell` has inline styles or CSS variables from `LayoutResult` (position, width, height). Layout is computed server-side; template renders positioned elements.
- **Aspect ratio**: Dynamic layout allows up to 10% cropping from each end of the longer dimension to improve fit. Cells use `object-fit: cover` and `object-position: center` so photos fill cells; the layout engine constrains cell aspect ratios to the allowed range. Gaps between cells are acceptable.

```html
<div class="collage-grid collage-dynamic" data-memories="true" data-layout="collage">
  <!-- Per-photo cells with computed positions -->
  <div class="collage-cell" data-asset-id="{id}" style="position:absolute;left:X%;top:Y%;width:W%;height:H%;" role="button" tabindex="0">
    <img src="/image/{id}" alt="Memory" loading="lazy" style="object-fit:cover;object-position:center;" />
  </div>
  ...
</div>
```

### Fallback (Engine Fails or >2s)

When the layout engine fails or times out:
- **Layout class**: `collage-grid collage-{n}` (simple grid, e.g. `collage-5`).
- **Aspect ratio**: Fallback MUST use `object-fit: contain` so photos are not cropped (update CSS from current `cover`).

### Data Attributes

Unchanged: `data-memories="true"`, `data-asset-id`, `data-duration`, `data-memory-caption`. Lightbox trigger remains `.collage-cell` click.

### Endpoints

No new endpoints. `GET/POST /asset/new` and `GET/POST /memories` behavior unchanged. Response includes HTML with either dynamic layout or fallback grid.

### ViewData Extension

`ViewData` gains optional field:
- **CollageLayout** (or equivalent): `*LayoutResult` — when non-nil, template uses dynamic layout; when nil, template uses `collage-{n}` fallback.

## Accessibility

Unchanged. Collage cells: `role="button"`, `tabindex="0"`. Lightbox: focus trap, Escape to close.
