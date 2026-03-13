# Implementation Plan: Memories (On This Day)

**Branch**: `007-memories-feature` | **Date**: 2026-03-11 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/007-memories-feature/spec.md`

## Summary

Add a "memories" feature that displays a collage of up to 8 photos from a randomly chosen past year (same month and day as today). The collage appears in the regular slideshow and via a dedicated top-UI button, uses a 45-second timer, shows "Memories from X years ago", and supports tapping a photo to open it in a lightbox. Implementation extends the existing Immich memories API usage and adds PhotoPrism support for "on this day" queries; introduces a new collage layout and lightbox component; and adds a memories button plus duration override for memories views.

## Technical Context

**Language/Version**: Go 1.26  
**Primary Dependencies**: Echo v5, Templ, Viper, existing internal packages (`config`, `source`, `photoprism`, `immich`, `templates`, `cache`)  
**Storage**: Existing cache layer for API responses; no new persistent storage  
**Testing**: `go test ./...` with table-driven tests; mock Immich/PhotoPrism APIs in tests  
**Target Platform**: Linux/Windows container or host; browser-based kiosk UI  
**Project Type**: Web service + HTML/Templ-based kiosk UI (HTMX, TypeScript)  
**Performance Goals**: Memories fetch and collage render within 3 seconds; no perceptible slideshow disruption  
**Constraints**: Kiosk-first unattended slideshow; memories duration fixed at 45s (configurable later if needed); i18n for all new strings  
**Scale/Scope**: Up to 8 photos per collage; one memories slot per slideshow cycle; PhotoPrism "on this day" limited to existing search API

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **Code Quality** (Go, golangci-lint, Templ) | Pass | New code follows existing patterns; no new deps beyond stdlib + existing |
| **Testing Standards** | Pass | Unit tests for memories fetch logic, collage layout selection; mock providers |
| **User Experience** (kiosk-first, config-driven, i18n, a11y) | Pass | Memories in slideshow and on-demand; i18n for captions/errors; lightbox pauses timer |
| **Performance** (smooth slideshow, cache, no N+1) | Pass | Memories cached; collage preload; 45s duration avoids thrash |
| **Dependencies** | Pass | No new external packages |

## Project Structure

### Documentation (this feature)

```text
specs/007-memories-feature/
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
│   └── immich_memories.go    # Extend: MemoriesCollage() for up-to-8 assets from one year
├── photoprism/
│   ├── client.go             # Add: PhotosOnThisDay() if API supports
│   └── provider.go          # Extend: RandomMemoryAsset → MemoriesCollage or new method
├── source/
│   └── display.go            # No change; DisplayAsset reused
├── routes/
│   ├── routes_asset.go       # Add memories collage path; duration override
│   ├── routes_asset_helpers.go  # Add MemoriesCollage retrieval; gatherAssetBuckets memories
│   └── routes_memories.go    # New: /memories (on-demand) + /memories/back
├── templates/
│   ├── components/
│   │   ├── image/
│   │   │   ├── image.templ    # Branch for collage layout
│   │   │   └── layout.templ   # Add collage layout
│   │   └── memories/
│   │       ├── collage.templ  # New: collage grid (1–8 cells)
│   │       └── lightbox.templ # New: single-photo lightbox
│   └── partials/
│       └── menu.templ        # Add memories button
├── config/
│   └── config.go             # Add memory_duration (default 45) if configurable
└── i18n/ + locales/           # New strings: memories_caption, no_memories, error_load

frontend/src/
├── ts/
│   ├── kiosk.ts              # Handle memories duration; lightbox pause/resume
│   ├── polling.ts             # Pause on lightbox open
│   └── memories-lightbox.ts   # New: lightbox open/close, timer pause
└── css/
    └── collage.css           # New: grid layouts for 1–8 photos
```

**Structure Decision**: Extend existing `internal/` layout. New `routes_memories.go` for on-demand flow; new `templates/components/memories/` for collage and lightbox; extend `immich_memories.go` and `photoprism/provider.go` for collage retrieval.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | — | — |
