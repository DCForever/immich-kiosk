# Implementation Plan: Memories Progressive Time-Range Expansion

**Branch**: `008-memories-expand-range` | **Date**: 2026-03-13 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/008-memories-expand-range/spec.md`

## Summary

Extend the Memories feature with progressive time-range fallback (today → week → month) and change the photo count from up to 8 to 3–16. When no photos exist for the exact date (same month and day) in any past year, the system expands to the ISO week containing that date; if still none, expands to the same calendar month. Collage shows 3–16 photos (minimum 3 required; when fewer exist, treat as no memories). For each photo count (3–16), the system defines at least 3 distinct layout variants; one is chosen at random when rendering. Implementation extends `MemoriesCollage` in Immich and PhotoPrism providers with progressive search logic; adds `TimeRange` to `MemoriesCollage` for caption selection; extends collage CSS layouts for 3–16 photos with ≥3 variants per count; adds layout variant selection in routes/template; adds i18n keys for week/month captions.

## Technical Context

**Language/Version**: Go 1.26  
**Primary Dependencies**: Echo v5, Templ, Viper, existing internal packages (`config`, `source`, `photoprism`, `immich`, `templates`, `cache`)  
**Storage**: Existing cache layer for API responses; no new persistent storage  
**Testing**: `go test ./...` with table-driven tests; mock Immich/PhotoPrism APIs in tests  
**Target Platform**: Linux/Windows container or host; browser-based kiosk UI  
**Project Type**: Web service + HTML/Templ-based kiosk UI (HTMX, TypeScript)  
**Performance Goals**: Memories fetch and collage render within 3 seconds; no perceptible slideshow disruption  
**Constraints**: Kiosk-first unattended slideshow; memories duration fixed at 45s; i18n for all new strings; ISO week per spec  
**Scale/Scope**: 3–16 photos per collage; progressive search (exact date → week → month); PhotoPrism uses `after`/`before` for date ranges

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **Code Quality** (Go, golangci-lint, Templ) | Pass | Extends existing memories code; no new deps |
| **Testing Standards** | Pass | Unit tests for progressive search, ISO week, 3–16 bounds; mock providers |
| **User Experience** (kiosk-first, config-driven, i18n, a11y) | Pass | New i18n keys for week/month captions; existing lightbox unchanged |
| **Performance** (smooth slideshow, cache, no N+1) | Pass | Reuses existing cache; progressive search adds up to 3 queries when falling back |
| **Dependencies** | Pass | No new external packages; stdlib `time` for ISO week |

## Project Structure

### Documentation (this feature)

```text
specs/008-memories-expand-range/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
internal/
├── immich/
│   └── immich_memories.go    # Extend MemoriesCollage: progressive search (date→week→month); 3–16 assets; TimeRange
├── photoprism/
│   └── provider.go           # Extend MemoriesCollage: progressive search via after/before; 3–16 assets; TimeRange
├── source/
│   └── display.go            # Extend MemoriesCollage: add TimeRange field; MinCollageAssets=3, MaxCollageAssets=16
├── routes/
│   ├── routes_asset_helpers.go  # ProcessMemoriesCollage: pass TimeRange to caption; pick random LayoutVariant; validate 3–16
│   └── routes_memories.go    # No change; uses ProcessMemoriesCollage
├── templates/
│   └── components/
│       └── memories/
│           └── collage.templ  # Support collage-{n}-{variant} for n=3..16, ≥3 variants per n
└── i18n/ + locales/          # New keys: memories_caption_week, memories_caption_month

frontend/src/
└── css/
    └── collage.css           # Add ≥3 layout variants per count (collage-{n}-a, collage-{n}-b, collage-{n}-c, ...) for n=3..16
```

**Structure Decision**: Extend existing `internal/` layout. No new files; modify `MemoriesCollage` logic, `MemoriesCollage` struct, collage template, and CSS. Provider layer implements progressive search; routes and templates consume `TimeRange` for caption selection. ProcessMemoriesCollage selects a layout variant at random when building view data; template applies `collage-{n}-{variant}` class.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | — | — |
