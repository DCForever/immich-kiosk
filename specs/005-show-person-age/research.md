# Research: Show age for people on photos

**Feature**: 005-show-person-age  
**Date**: 2026-03-09

## 1. External birthdate file format

**Decision**: Support JSON as the primary format for the birthdate file. Optional: add CSV in a follow-up if needed.

**Rationale**: JSON is easy to parse in Go (`encoding/json`), supports a clear key→DOB structure, and is common for small config-like data. The spec allows "e.g. JSON or CSV"; starting with JSON keeps the first implementation simple and avoids CSV parsing edge cases (quoting, headers).

**Alternatives considered**:
- CSV: Simple for non-technical admins to edit; deferred until user demand.
- YAML: Already used for main config; would mix “data” with “config”; JSON keeps birthdate data separate.

**Schema (see contracts/)**:
- JSON: `{ "<birthdate-file-key>": "YYYY-MM-DD", ... }`
- Keys are matched via the separate mapping file (Photoprism person ID or name → birthdate-file key).

---

## 2. Person–birthdate mapping format

**Decision**: Single JSON file that maps Photoprism person identifier to birthdate-file key. Config key: `birthdate_mapping_path` (and env `KIOSK_BIRTHDATE_MAPPING_PATH`).

**Rationale**: Spec requires "administrator-defined mapping (e.g. separate mapping file or table)". A separate file keeps mapping and birthdate data decoupled: one file lists all DOBs by logical key; the other maps backend person ID/name to that key. Supports Photoprism SubjUID or display name as the key.

**Alternatives considered**:
- Single file with both mapping and DOBs: Possible but mixes two concerns; two files are clearer.
- Inline in main config: Would bloat config; large mappings are better in dedicated files.

**Schema (see contracts/)**:
- JSON: `{ "<photoprism-person-id-or-name>": "<birthdate-file-key>", ... }`
- Lookup: normalize (trim, optional case-fold for name) then exact match; first match wins if multiple entries (deterministic).

---

## 3. Where to enrich DisplayAsset.People with birthdates

**Decision**: Enrich in the PhotoPrism provider when building `DisplayAsset`, immediately after `photoMarkersToPeople(ph)`. The birthdate loader is injected (or obtained from config) so the provider stays testable.

**Rationale**: PhotoPrism markers do not include DOB; the API Subject list can include BirthDate but the photo→markers path does not. Enriching at asset-build time keeps the rest of the app (templates, Immich) unchanged. Immich already provides BirthDate from the API; no enrichment for Immich.

**Alternatives considered**:
- Enrich in a central “display pipeline” for all sources: Would require a broader refactor; current per-provider build is sufficient.
- Enrich in templates: Would require passing the loader into templates and mixing concerns; rejected.

---

## 4. 120-year cap and under-1 display (months)

**Decision**: Enforce a 120-year maximum in the age calculation used for display (e.g. in `calculateAge` in `metadata.templ`). If computed age &gt; 120, return empty string (do not show age). For under 1 year, keep existing behavior: display age in whole months (e.g. "6 months"); spec confirms "in months".

**Rationale**: Spec (FR-005, Edge Cases) states computed age &gt; 120 is invalid; no display. Current template already supports months for &lt; 1 year; ensure it uses whole months only and no "0 years" for babies (use months). No change to Immich flow except this shared age logic.

**Alternatives considered**:
- Configurable max age: Spec fixed at 120; avoid config creep.
- Enforcing in loader: We still validate in one place (template/helper) so that any future source of birthdates (e.g. Immich) also respects the cap.

---

## 5. File missing or unreadable

**Decision**: If the birthdate file or mapping file is missing or unreadable at startup or reload: log a warning (slog), do not surface an error to the viewer, and do not show any ages. Photos and people names still display. Loader returns empty/nil lookup until files are valid.

**Rationale**: Matches spec FR-009 and clarification: "Show photos and people as usual; do not show any ages until the file is present and valid. Log a warning but do not surface an error to the viewer."

**Implementation**: Loader `Load()` or `Reload()` returns an error only for logging; callers treat missing file as “no birthdates available” and continue. Provider passes people through unchanged (no BirthDate set) when loader is unavailable or returns no match.

---

## 6. Time zone and age calculation

**Decision**: Use the photo’s taken date (already `time.Time` in the asset) and parse birthdate as date-only (YYYY-MM-DD) in a consistent zone (e.g. UTC or local) for age calculation. Ensure whole months/years are computed so that time-of-day does not change the displayed age (e.g. birth 2000-01-15, photo 2001-01-14 → 0 years; photo 2001-01-15 → 1 year).

**Rationale**: Spec: "Time zone differences ... do not cause the displayed age ... to be off by one unit." Existing `calculateAge` in metadata.templ already uses date-based year/month/day math; confirm it uses truncated/date-only comparison where needed.

---

## Summary table

| Topic | Decision |
|-------|----------|
| Birthdate file format | JSON: `{ "key": "YYYY-MM-DD", ... }` |
| Mapping file format | JSON: `{ "person-id-or-name": "birthdate-key", ... }` |
| Enrichment point | PhotoPrism provider, after `photoMarkersToPeople()` |
| 120-year cap | In age calculation; return empty if age &gt; 120 |
| Under 1 year | Display in whole months (existing behavior) |
| File missing | Log warning; show no ages; no user-facing error |
| Time zone | Date-only comparison for whole years/months |
