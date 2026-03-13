# Memories UI Contract

**Feature**: 007-memories-feature  
**Date**: 2026-03-11

## Overview

This document defines the contract between the backend (Go/Templ) and frontend (HTMX/TypeScript) for the memories feature.

## Endpoints

### GET/POST /asset/new (extended)

When the asset source is memories and the provider returns a collage:

- **Response**: HTML fragment (HTMX swap) containing the memories collage component.
- **Data attributes** (on the main container or a wrapper):
  - `data-memories="true"` — Indicates this is a memories view.
  - `data-duration="45"` — Display duration in seconds (injected via kiosk-data).
  - `data-memory-caption="Memories from X years ago"` — Caption for bottom UI.

### GET/POST /memories (new)

On-demand memories view.

- **Request**: Same query params as slideshow (device, user, etc.).
- **Response**: HTML fragment with memories collage, back button, and `data-memories="true"`.
- **Behavior**: Renders collage or "No memories for this day" / error state per spec.

### GET /image/:imageID (existing)

Used by lightbox to display full-size image when user taps a collage cell.

- **Request**: `imageID` from the tapped asset.
- **Response**: Image blob (JPEG/PNG).
- **No change** to existing contract.

## HTML Structure

### Memories collage container

```html
<div id="kiosk" data-memories="true" data-layout="collage" class="layout-collage" ...>
  <div class="collage-grid collage-{n}">
    <!-- n = 1..8, number of assets -->
    <div class="collage-cell" data-asset-id="{id}" role="button" tabindex="0">
      <img src="..." alt="..." />
    </div>
    ...
  </div>
  <div class="memories-caption">Memories from X years ago</div>
</div>
```

### Lightbox (client-rendered)

- **Trigger**: Click/tap on `.collage-cell`.
- **Structure**: Overlay with `class="memories-lightbox"`; contains `<img src="/image/{id}">` and close control.
- **Behavior**: Pause polling on open; resume on close. Timer state managed client-side.

## kiosk-data (JSON)

For memories view, the script tag `#kiosk-data` includes:

```json
{
  "duration": 45,
  "memories": true,
  ...
}
```

- `duration`: 45 for memories (overrides default).
- `memories`: true — Frontend uses this to enable lightbox behavior and caption display.

## i18n Keys

| Key | Purpose |
|-----|---------|
| `memories_caption` | "Memories from X years ago" (X is interpolated) |
| `memories_no_photos` | "No memories for this day" |
| `memories_error_load` | "Couldn't load memories" |
| `memories_back` | "Back" (dismiss button for on-demand) |

## Accessibility

- Collage cells: `role="button"`, `tabindex="0"`, keyboard activatable.
- Lightbox: Focus trap; Escape to close; `aria-modal="true"`.
