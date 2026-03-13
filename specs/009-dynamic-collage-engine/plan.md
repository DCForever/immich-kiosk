# Implementation Plan: Dynamic Collage Engine

**Branch**: `009-dynamic-collage-engine` | **Date**: 2026-03-13 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/009-dynamic-collage-engine/spec.md`

## Summary

Replace the fixed CSS grid collage implementation with a **pure Go layout engine** that computes dynamic arrangements from photo aspect ratios. The engine outputs layout data (positions, sizes) consumed by the Templ collage template, producing HTML with per-photo elements for lightbox-on-tap. Layout algorithms (justified rows, treemap packing, hero-plus-cluster) are selected at random per render. Fallback to simple grid when engine fails or exceeds 2 seconds. See [research.md](./research.md) for approach evaluation.

## Technical Context

**Language/Version**: Go 1.26  
**Primary Dependencies**: Echo v5, Templ, Viper, existing internal packages (config, source, photoprism, templates)  
**Storage**: N/A (layout computed at request time)  
**Testing**: `go test ./...`, testify, table-driven tests; mock provider for collage flow  
**Target Platform**: Linux, Windows, macOS (same as kiosk)  
**Project Type**: Web application (Go backend + Templ/HTMX frontend)  
**Performance Goals**: Collage layout computation <2 seconds; layout math completes in milliseconds  
**Constraints**: No Python/ImageMagick; output must allow per-photo tap (HTML); 2s timeout triggers fallback  
**Scale/Scope**: 3–16 photos per collage; single kiosk instance

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| Go standards | Pass | New code in Go; golangci-lint, small packages |
| Dependencies | Pass | No new external deps; stdlib + existing |
| Frontend (TS/CSS) | Pass | Minimal change; CSS for object-fit; no new JS deps |
| Templ | Pass | Template consumes layout data; no heavy logic |
| Testing | Pass | Tests for layout engine, integration for collage flow |
| Kiosk-first UX | Pass | Unattended display; 45s timer, no blocking |
| Performance | Pass | Layout <2s; preload unchanged |
| i18n | Pass | No new user-facing strings from engine |

**No violations.**

## Project Structure

### Documentation (this feature)

```text
specs/009-dynamic-collage-engine/
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
├── collage/                    # NEW: layout engine
│   ├── layout.go              # LayoutResult, ComputeLayout entry point
│   ├── justified.go            # Justified-rows algorithm
│   ├── treemap.go              # Treemap packing algorithm
│   ├── hero.go                 # Hero-plus-cluster algorithm
│   └── layout_test.go          # Unit tests
├── routes/
│   └── routes_asset_helpers.go # ProcessMemoriesCollage: call engine, attach layout; fallback
├── templates/
│   └── components/
│       └── memories/
│           └── collage.templ   # Render from LayoutResult or fallback to collage-{n}
frontend/
├── src/
│   └── css/
│       └── collage.css         # object-fit: contain for aspect-ratio preservation
```

**Structure Decision**: Add `internal/collage` package for layout engine. Routes and templates are modified in place. No new top-level directories.

## Complexity Tracking

*No violations requiring justification.*
