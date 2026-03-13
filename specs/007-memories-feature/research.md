# Research: Memories (On This Day)

**Feature**: 007-memories-feature  
**Date**: 2026-03-11

## 1. Immich Memories API — Collage from One Year

**Decision**: Use existing Immich `/api/memories` with `for=<date>`; filter returned `Memory` objects by `Data.Year`; pick one year at random from those with assets; take up to 8 assets from that Memory.

**Rationale**: The Immich `Memory` struct has `Data.Year` (see `internal/immich/immich.go`). The API returns memories for "on this day" — each Memory corresponds to a specific year. We filter to memories with ≥1 asset, pick a random year, then take up to 8 assets from that Memory's `Assets` slice.

**Alternatives considered**:
- New Immich endpoint for "on this day" by year — not available; existing memories API suffices.
- Fetch all years and aggregate — would require multiple API calls; current API returns multiple Memory objects per call for the same month/day across years.

---

## 2. PhotoPrism "On This Day" Support

**Decision**: Add PhotoPrism support via `GET /api/v1/photos` with `taken_at` or `created_at` date filter (month and day). Query photos where `TakenAt` matches the chosen year's month-day; select up to 8 randomly.

**Rationale**: PhotoPrism Swagger (docs.photoprism.dev) documents search parameters. The photos API supports date-based filters. We construct a date range for the chosen year (e.g. 2007-03-11 00:00 to 2007-03-11 23:59) and fetch photos. If the API does not support exact-day filter, use a small range or `q` parameter if available.

**Alternatives considered**:
- Leave PhotoPrism without memories — spec implies multi-source support; Photoprism is a first-class source.
- Custom PhotoPrism endpoint — not available; use standard search.

**Implementation note**: Verify PhotoPrism `/api/v1/photos` or `/api/v1/search` parameters for date filtering before implementation. Fallback: return "memories not supported" for PhotoPrism until API is confirmed.

---

## 3. Collage Layout Implementation

**Decision**: Use CSS Grid with predefined layout classes for counts 1–8. Define one layout per count (e.g. `collage-1` through `collage-8`) following the spec's "interesting grid" characteristics (varied cell sizes, asymmetrical, sub-grids). Layouts are static CSS; no JS layout calculation.

**Rationale**: Keeps implementation simple; layouts are predictable and testable. The spec's reference image shows 6 styles — we map counts 1–8 to similar patterns (e.g. 1 = full bleed; 2 = split; 3–4 = asymmetric; 5–8 = grid with sub-grids).

**Alternatives considered**:
- Dynamic JS layout — overkill; CSS Grid sufficient.
- Masonry library — adds dependency; constitution prefers minimal deps.

---

## 4. Duration Override for Memories

**Decision**: When rendering a memories view, set `viewData.Duration = 45` (or `config.MemoryDuration` if added). The kiosk-data JSON passed to the frontend includes `duration`; the polling interval uses it. No change to the polling mechanism — only the value passed per response.

**Rationale**: The home template injects `viewData.Duration` into `kiosk-data`. The asset response (HTMX swap) includes a new kiosk container with its own duration. For memories, we render a view that sets `Duration: 45` in the view data.

**Alternatives considered**:
- Config option `memory_duration` — deferred; use 45 as constant initially; add config later if requested.
- Frontend override — would require passing a flag; backend override is cleaner.

---

## 5. Lightbox Implementation

**Decision**: Add a client-side lightbox (full-screen overlay) for the collage. On photo tap: show overlay with the tapped image, pause the HTMX polling (and thus the 45s timer). On close: hide overlay, resume polling with remaining time. Use existing `pausePolling`/`resumePolling` from `polling.ts`; add `memories-lightbox.ts` for open/close and timer state.

**Rationale**: The spec requires pausing the timer while the lightbox is open. The existing more-info overlay pauses polling via `handleOverlayToggle`. We replicate that pattern for the memories lightbox. The lightbox shows a single image (reuse `/image/:id` or inline the image URL from the collage asset).

**Alternatives considered**:
- Server-rendered lightbox (HTMX) — would require round-trip; client-side is snappier.
- Reuse more-info overlay — different purpose (metadata vs. full-size view); separate component clearer.

---

## 6. On-Demand Memories Route

**Decision**: Add `GET /memories` (or `POST` for HTMX) that renders the memories collage in the same main `#kiosk` area. Include a back/close button that triggers `POST /asset/new` (or equivalent) to return to the slideshow. The memories view is a full-page swap like the asset view.

**Rationale**: Keeps the flow consistent: the main content area swaps between asset view and memories view. The back button can use HTMX `hx-post="/asset/new"` with current query params to resume slideshow.

**Alternatives considered**:
- Modal overlay for memories — spec says "return to previous context"; full swap is simpler and matches asset flow.
- Separate route `/memories/view` — single `/memories` suffices; back goes to `/asset/new`.
