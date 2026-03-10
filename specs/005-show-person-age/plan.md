# Implementation Plan: Show age for people on photos

**Branch**: `005-show-person-age` | **Date**: 2026-03-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-show-person-age/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Show each identified person’s age on photos (in years, or in months when under 1) based on the photo’s taken date and birthdates from an external local file. Use an administrator-defined mapping to link Photoprism person ID or name to the birthdate file. When the file or mapping is missing or invalid, show photos and people without ages and log a warning. Enrich PhotoPrism-sourced `DisplayAsset.People` with birthdates from the loader; Immich continues to use API birthdates. Enforce a 120-year maximum age and existing template age logic (months &lt; 1, years ≥ 1).

## Technical Context

**Language/Version**: Go 1.26  
**Primary Dependencies**: Echo v5, Templ, Viper, existing internal packages (config, source, photoprism, templates)  
**Storage**: Local files only (birthdate file + mapping file; paths from config/env). No new DB.  
**Testing**: `go test ./...`, table-driven tests for loader and age validation; mock file system where useful.  
**Target Platform**: Same as kiosk (Linux/Windows/macOS, browser-based UI).  
**Project Type**: Web application (Go backend, Templ + HTMX frontend).  
**Performance Goals**: Birthdate lookup O(1) per person after load; file read on startup and on config reload; no blocking of slideshow.  
**Constraints**: No new external services; birthdate/mapping files are optional (graceful degradation).  
**Scale/Scope**: Hundreds of people and thousands of photos; file sizes small (JSON/CSV &lt; 1MB typical).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|----------|--------|--------|
| **Code Quality** (Go, naming, log/slog, small packages) | Pass | New code in existing `internal` layout; use `log/slog` for loader warnings. |
| **Testing** (exported behavior, table-driven, mock external) | Pass | Tests for birthdate loader, mapping resolution, age validation (120 cap, &lt;1 in months). |
| **UX** (kiosk-first, config-driven, i18n, a11y) | Pass | Age display already gated by `ShowPersonAge`; use existing i18n for new strings if any. |
| **Performance** (smooth slideshow, cache, no N+1) | Pass | Loader cached in memory; enrichment in existing asset-build path. |
| **Dependencies** | Pass | Stdlib + existing deps; no new external services. |

## Project Structure

### Documentation (this feature)

```text
specs/005-show-person-age/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (birthdate file + mapping schema)
└── tasks.md             # Phase 2 output (/speckit.tasks – not created by /speckit.plan)
```

### Source Code (repository root)

```text
internal/
├── config/
│   └── config.go              # Add BirthdateFilePath, BirthdateMappingPath (yaml + env)
├── photoprism/
│   └── provider.go            # After photoMarkersToPeople(), enrich People via birthdate loader
├── source/
│   └── display.go             # Person already has BirthDate; no change
├── birthdate/                 # NEW: loader package
│   ├── loader.go              # Load birthdate file + mapping; lookup by person ID/name
│   └── loader_test.go
├── templates/
│   └── partials/
│       └── metadata.templ     # calculateAge: add 120-year cap; ensure <1 year in months
└── ...
```

**Structure Decision**: Add a new `internal/birthdate` package for file loading and lookup. Config fields live in existing `internal/config`. The PhotoPrism provider obtains the birthdate loader from application wiring: when config has birthdate paths set, one loader instance is created, `Load` is called (at startup or when the provider is constructed), and that same instance is reused for enrichment; on config reload the loader is re-loaded so updated files are picked up (see tasks T007 and T014). PhotoPrism provider calls the loader when building `DisplayAsset`. Template age logic is updated for 120-year cap and spec-compliant months display.

## Complexity Tracking

> No constitution violations. Leave table empty or omit section.
