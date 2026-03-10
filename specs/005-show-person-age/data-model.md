# Data Model: Show age for people on photos

**Feature**: 005-show-person-age  
**Date**: 2026-03-09

## In-memory / runtime

### Person (existing, unchanged)

- **Source**: `internal/source/display.go`
- **Fields**: `ID`, `Name`, `BirthDate` (type `source.BirthDate`, string `YYYY-MM-DD`).
- **Usage**: Filled by each backend (Immich from API; PhotoPrism from markers then enriched by birthdate loader when available).

### Birthdate loader state

- **Package**: `internal/birthdate`
- **Concepts**:
  - **Birthdate map**: `map[string]string` — birthdate-file key → `"YYYY-MM-DD"`. Loaded from the birthdate JSON file.
  - **Mapping map**: `map[string]string` — Photoprism person ID or name (normalized) → birthdate-file key. Loaded from the mapping JSON file.
- **Lifecycle**: When config has birthdate paths set, one loader instance is created and `Load(birthdatePath, mappingPath)` is called (at startup or when the PhotoPrism provider is constructed). The same instance is made available to the provider for `Lookup` during asset build. On config reload, the loader is re-loaded so updated files are picked up (see FR-007). No TTL; file changes require reload.
- **Validation**: On load, reject malformed JSON; reject DOB not parseable as `2006-01-02`. Invalid entries are skipped and optionally logged. Duplicate keys in JSON: last wins (deterministic).

### DisplayAsset (existing, unchanged)

- **People**: `[]source.Person`. For PhotoPrism, each `Person` may have `BirthDate` set by the loader when a mapping and birthdate exist; otherwise `BirthDate` remains empty and no age is shown.

---

## Config (new fields)

- **BirthdateFilePath** (config key: `birthdate_file_path`, env: `KIOSK_BIRTHDATE_FILE_PATH`): Path to the JSON file containing key → DOB. Empty means “disabled”; no ages from external source.
- **BirthdateMappingPath** (config key: `birthdate_mapping_path`, env: `KIOSK_BIRTHDATE_MAPPING_PATH`): Path to the JSON file containing person ID/name → birthdate-file key. Empty means no mapping; no ages from external source.

Both are optional. If either is missing or unreadable, the loader logs a warning and lookup returns no birthdates.

---

## File formats (contracts)

See `contracts/` for:

- **birthdate-file.json**: Schema and example for the birthdate file.
- **birthdate-mapping.json**: Schema and example for the person → birthdate-key mapping file.

---

## State transitions

1. **Startup**: Config loaded → if paths set, loader loads both files; on failure, log warning and leave maps empty.
2. **Config reload**: Re-run loader; replace in-memory maps.
3. **Asset build (PhotoPrism)**: For each person in `photoMarkersToPeople(ph)`, lookup by `Person.ID` then by `Person.Name` in mapping; if found, lookup DOB in birthdate map; if valid, set `Person.BirthDate`.
4. **Display**: Existing template logic uses `Person.BirthDate` and `LocalDateTime` to compute and show age (with 120-year cap and months for &lt; 1 year).

---

## Validation rules (from spec)

- **FR-005**: Do not show age if photo taken date missing, birthdate missing, birthdate invalid, or computed age &gt; 120.
- **Duplicate mapping keys**: Deterministic (e.g. last wins in JSON decode).
- **Duplicate birthdate keys**: Deterministic (last wins).
- **Date format**: DOB must be `YYYY-MM-DD`; otherwise skip entry and do not set `Person.BirthDate`.
