# Tasks: Show age for people on photos

**Input**: Design documents from `/specs/005-show-person-age/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Go packages under `internal/` at repository root
- Config and schema at repo root: `config.example.yaml`, `config.schema.json`

---

## Phase 1: Setup

**Purpose**: Verify prerequisites and create the new package layout.

- [x] T001 Verify feature branch and design docs (plan.md, spec.md, data-model.md, research.md, contracts/) present in specs/005-show-person-age/
- [x] T002 Create internal/birthdate package directory and add package doc in internal/birthdate/loader.go (empty Load/Lookup stubs or minimal types so package compiles)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Config and birthdate loader that all user stories depend on.

**Independent Test**: Loader loads valid JSON birthdate and mapping files; Lookup(personID, personName) returns DOB when mapping and birthdate exist; returns empty when files missing (with logged warning).

- [x] T003 Add BirthdateFilePath and BirthdateMappingPath to Config struct in internal/config/config.go (yaml, mapstructure, env KIOSK_BIRTHDATE_FILE_PATH, KIOSK_BIRTHDATE_MAPPING_PATH)
- [x] T004 [P] Add birthdate_file_path and birthdate_mapping_path to config.example.yaml and config.schema.json at repo root
- [x] T005 Implement Load(birthdatePath, mappingPath) and Lookup(personID, personName) in internal/birthdate/loader.go; log warning via slog when either file missing or unreadable; parse JSON; validate DOB as YYYY-MM-DD; skip invalid entries
- [x] T006 Write unit tests for loader in internal/birthdate/loader_test.go (valid load, missing file, unreadable file, invalid date, lookup by ID then name, duplicate keys last-wins)
- [x] T007 Wire birthdate loader into application: when config has birthdate paths set, create one loader instance, call Load(birthdatePath, mappingPath), and make the loader available to the PhotoPrism provider (e.g. provider holds a reference or receives it from wherever the provider is constructed) in internal/photoprism/provider.go and/or the place that constructs the provider

**Checkpoint**: Foundation ready — loader and config can be used by provider.

---

## Phase 3: User Story 1 – See age for each person on a photo (Priority: P1) – MVP

**Goal**: Display each person’s age (years or months when &lt; 1) next to their name on photos when birthdates are provided by the external file and mapping.

**Independent Test**: Display a photo with one or more identified people and valid birthdates in the external source; verify each person’s age (years or months when under 1) matches the difference between photo taken date and birthdate.

- [x] T008 [US1] Add 120-year maximum age check in calculateAge in internal/templates/partials/metadata.templ (return empty string when computed age &gt; 120)
- [x] T009 [US1] Verify under-1 year displays as whole months (no "0 years") in internal/templates/partials/metadata.templ
- [x] T010 [US1] In internal/photoprism/provider.go, after photoMarkersToPeople(ph), enrich each Person with BirthDate from birthdate loader (lookup by Person.ID then Person.Name); only set when loader is available and Lookup returns non-empty valid DOB

**Checkpoint**: User Story 1 is testable — photos with mapped birthdates show correct ages (years or months).

---

## Phase 4: User Story 2 – Handle missing or partial birthdate data gracefully (Priority: P2)

**Goal**: When birthdate is missing or invalid, show no age and no error; rest of photo and names display normally.

**Independent Test**: Show photos for people without birthdate in the external source; confirm no age is displayed and no placeholder or error; other people on the same photo with valid birthdates still show ages.

- [x] T011 [US2] Add loader test cases for missing and unreadable birthdate/mapping files in internal/birthdate/loader_test.go (expect empty lookup and logged warning)
- [x] T012 [US2] Ensure in internal/photoprism/provider.go that when loader is nil or Lookup returns empty, Person.BirthDate is left unchanged and no error is returned to the viewer

**Checkpoint**: User Story 2 verified — missing/invalid birthdates never show an age or break the UI.

---

## Phase 5: User Story 3 – Respect configuration of external birthdate source (Priority: P3)

**Goal**: Administrator can configure paths and see ages after reload; changes to birthdate/mapping files are reflected after config reload.

**Independent Test**: Set birthdate_file_path and birthdate_mapping_path, add one person to the files, reload kiosk/config; confirm that person’s photos show age; edit the birthdate file and reload; confirm ages update.

- [ ] T013 [US3] Document birthdate_file_path and birthdate_mapping_path in config.example.yaml (comments) and ensure specs/005-show-person-age/quickstart.md references both config keys
- [ ] T014 [US3] Wire birthdate loader reload when config is reloaded (e.g. in config watcher or wherever config reload runs) so updated birthdate/mapping files are picked up without full restart

**Checkpoint**: User Story 3 verified — config and file changes take effect after reload.

---

## Phase 6: Polish & cross-cutting

**Purpose**: Schema, tests, and quickstart validation.

- [ ] T015 [P] Add description fields for birthdate_file_path and birthdate_mapping_path in config.schema.json
- [ ] T016 Run go test ./internal/birthdate/... ./internal/photoprism/... and fix any failures
- [ ] T017 Validate quickstart steps from specs/005-show-person-age/quickstart.md (manual or document result)
- [ ] T018 [P] Add table-driven tests for age calculation (120-year cap, under-1 months) in internal/templates/partials or a testable helper used by metadata.templ so calculateAge behavior is covered

---

## Dependencies & execution order

### Phase dependencies

- **Phase 1**: No dependencies.
- **Phase 2**: Depends on Phase 1 (package exists). Blocks all user stories.
- **Phase 3 (US1)**: Depends on Phase 2 (config + loader + loader wiring T007).
- **Phase 4 (US2)**: Depends on Phase 2 and Phase 3 (loader behavior and provider enrichment).
- **Phase 5 (US3)**: Depends on Phase 2 (config); reload wiring may depend on Phase 3.
- **Phase 6**: Depends on Phases 2–5.

### User story dependencies

- **US1 (P1)**: After Phase 2. No dependency on US2/US3.
- **US2 (P2)**: After Phase 2 and US1 (enrichment and “no age” path).
- **US3 (P3)**: After Phase 2; can overlap with US1/US2 (config and reload).

### Parallel opportunities

- T004 [P] can run in parallel with T003 (different files).
- T008 and T009 both touch metadata.templ but different logic; can be done in one pass or sequentially.
- T011 and T012 can run after T010 (different files: loader_test.go vs provider.go).
- T015 [P] can run in parallel with other polish tasks.
- T018 [P] can run in parallel with other polish tasks.

---

## Parallel example: Phase 2

```bash
# After T003:
# T004 (config.example + config.schema) can run in parallel with T005 (loader implementation)
```

---

## Implementation strategy

### MVP first (User Story 1 only)

1. Phase 1: Setup  
2. Phase 2: Foundational (config + loader + tests)  
3. Phase 3: US1 (120-year cap, under-1 months, provider enrichment + loader already wired in Phase 2)
4. **Stop and validate**: Manual test with PhotoPrism, one photo, one person with birthdate in file + mapping  

### Incremental delivery

1. Phase 1 + 2 → loader and config ready  
2. Phase 3 (US1) → ages on photos (MVP)  
3. Phase 4 (US2) → graceful missing/invalid  
4. Phase 5 (US3) → config docs and reload  
5. Phase 6 → schema, tests, quickstart  

### Suggested MVP scope

Phases 1–3 (T001–T010): config, birthdate loader, loader wiring, 120-year cap, under-1 months, PhotoPrism enrichment. Delivers “see age for each person on a photo” with external birthdate file and mapping.

---

## Notes

- [P] = different files, no shared state; [Story] = traceability to spec.
- Commit after each task or logical group; run tests after Phase 2 and Phase 3.
- Immich continues to use API birthdates; no changes to Immich adapter for this feature.
