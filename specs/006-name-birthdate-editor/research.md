# Research: Click name to add, change, or remove birthdate

## Decision 1: Where to persist edited birthdates

- **Decision**: Persist edits directly to the existing JSON birthdate file (`birthdates.json`) and reuse the existing mapping file (`birthdate-mapping.json`) rather than introducing a new store.
- **Rationale**: The `005-show-person-age` feature already defines these contracts, config keys, and loader behavior. Writing into the same files keeps a single source of truth for ages and avoids duplicating mapping logic.
- **Alternatives considered**:
  - Separate per-feature store (e.g. a new JSON file just for UI edits) – rejected because it would require additional reconciliation logic and risks drift from the age-display source.
  - Adding a database (SQLite/Postgres) – rejected as unnecessary complexity for a small, low-write dataset that already fits cleanly in JSON.

## Decision 2: How the UI exposes click-to-edit

- **Decision**: Make person names clickable only in interactive/setup views (where users are expected to interact) and keep the unattended slideshow loop non-interactive.
- **Rationale**: The constitution emphasizes kiosk-first, unattended slideshow behavior; interactive flows should not block or interfere with normal slideshow playback. Restricting click-to-edit to setup or an explicit “edit” mode keeps UX consistent.
- **Alternatives considered**:
  - Allow click-to-edit during slideshow playback – rejected because accidental clicks could interrupt the display and conflict with kiosk-first behavior.
  - Provide a completely separate admin UI outside the kiosk – deferred; may be added later, but reusing existing views keeps scope small for this feature.

## Decision 3: Validation rules for birthdates

- **Decision**: Reuse the same validation rules as the age-display feature: dates must parse as `YYYY-MM-DD`, not be in the future, and must not imply an age over 120 years at photo time.
- **Rationale**: Aligning validation with the age-display rules ensures that any saved birthdate can actually be used for age computation and avoids inconsistent edge handling.
- **Alternatives considered**:
  - Permitting any date and letting the display logic drop invalid ages – rejected because it would allow storing unusable data and confuse users when ages fail to appear.
  - Enforcing stricter rules (e.g. minimum age) – rejected as unnecessary for the current domain.

## Decision 4: How edits reach the loader used by PhotoPrism provider

- **Decision**: Implement helper functions in `internal/birthdate` to update the in-memory maps and write JSON atomically (write to temp file then rename), and invoke a reload or in-memory update after each successful edit.
- **Rationale**: The existing loader already encapsulates how JSON is parsed and validated. Keeping write logic adjacent, and reusing the same structs and normalization, ensures consistency and makes it easier to reason about cache vs. file state.
- **Alternatives considered**:
  - Writing JSON directly from handlers without touching the loader – rejected because it risks inconsistencies between loader expectations and file layout.
  - Forcing a full process restart to pick up changes – rejected as too heavy and disruptive for a kiosk.

## Decision 5: Transport between UI and backend

- **Decision**: Use simple JSON POST/DELETE endpoints under the existing Echo server (e.g. `/api/person-birthdate`) that accept a person identifier and `YYYY-MM-DD` string, and respond with success or validation errors; use htmx or minimal JS to wire the clickable names to these endpoints.
- **Rationale**: The project already uses Echo and Templ; adding thin JSON endpoints matches existing patterns, while letting the UI update without heavy client-side frameworks.
- **Alternatives considered**:
  - Full-page form submissions for each edit – rejected as clunky for quick, in-place edits.
  - Introducing a richer SPA-style frontend – rejected as overkill and contrary to the project’s lightweight frontend approach.

