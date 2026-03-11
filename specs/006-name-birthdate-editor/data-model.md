# Data Model: Click name to add, change, or remove birthdate

## Runtime Entities

### Person (existing)

- **Source**: `internal/source`
- **Fields**:
  - `ID` (string): Stable identifier for a person (from Immich/PhotoPrism).
  - `Name` (string): Display name shown in the UI.
  - `BirthDate` (string, `YYYY-MM-DD`): Optional date of birth, filled from the birthdate loader.
- **Notes**: This feature does not change the core `Person` struct; it relies on the same fields as the age-display feature.

### Birthdate Loader State (extended)

- **Package**: `internal/birthdate`
- **Existing Concepts**:
  - **Birthdate map**: `map[string]string` – birthdate-file key → `"YYYY-MM-DD"`.
  - **Mapping map**: `map[string]string` – person ID or normalized name → birthdate-file key.
- **New Behavior**:
  - Functions to:
    - Add or update a birthdate for a given person identifier (ID or name) by:
      - Resolving/creating the appropriate birthdate-file key.
      - Updating in-memory maps.
      - Writing back to `birthdates.json` atomically.
    - Remove a birthdate mapping and corresponding birthdate entry when requested.
  - All writes validate date format and ensure the resulting data remains usable by the age-display feature.

### Birthdate Source File (existing contract, now writable)

- **File**: `birthdates.json`
- **Schema**: As defined in `specs/005-show-person-age/contracts/birthdate-file.json`.
- **Behavior**:
  - On load: parsed into the birthdate map.
  - On edit: re-serialized via `internal/birthdate` helpers with atomic replace (e.g. write temp file then rename).

### Birthdate Mapping File (existing contract, potentially extended)

- **File**: `birthdate-mapping.json`
- **Schema**: As defined in `specs/005-show-person-age/contracts/birthdate-mapping.json`.
- **Behavior**:
  - On load: parsed into the mapping map.
  - On edit: when a new person is given a birthdate and has no existing mapping entry, this feature may add a new mapping entry using a deterministic key.

### HTTP API DTOs (planned)

- **SetBirthdateRequest**:
  - `personId` (string) – required.
  - `birthDate` (string, `YYYY-MM-DD`) – required for add/change.
- **SetBirthdateResponse**:
  - `status` (string, e.g. `"ok"` or `"error"`).
  - `message` (string, localized, for display if needed).
- **DeleteBirthdateRequest**:
  - `personId` (string) – required.
- **DeleteBirthdateResponse**:
  - Same shape as `SetBirthdateResponse`.

## Validation Rules

- `birthDate` must parse as `YYYY-MM-DD`.
- `birthDate` must not be in the future relative to “today”.
- Combined with photo taken date, implied age must not exceed 120 years; if it would, the date is rejected on save.
- For delete operations:
  - If there is no stored birthdate/mapping for the given person, deletion is a no-op but returns success to keep the UX simple.

