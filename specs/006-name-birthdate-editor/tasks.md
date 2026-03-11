# Tasks: Click name to add, change, or remove birthdate

**Input**: Design documents from `/specs/006-name-birthdate-editor/`
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`

**Tests**: This feature benefits from targeted tests around date validation and JSON persistence, but exhaustive tests are not required for all tasks.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Ensure config and birthdate files are ready for interactive editing.

- [x] T001 Verify `birthdate_file_path` and `birthdate_mapping_path` defaults in `config.yaml` align with JSON contracts (`birthdates.json`, `birthdate-mapping.json`).
- [x] T002 [P] Add or update sample `birthdates.json` and `birthdate-mapping.json` in the repo root to match the existing contracts (for local dev and manual testing).
- [x] T003 [P] Confirm existing age-display feature (005) correctly loads and uses the configured birthdate files by running the kiosk and checking ages render as expected.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure enabling safe, atomic birthdate edits and mapping updates.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T004 Introduce helper functions in `internal/birthdate` to write updated birthdate and mapping maps back to JSON files atomically (temp file + rename).
- [x] T005 [P] Add functions in `internal/birthdate` to set or update a birthdate for a person identifier (ID or name) and update both in-memory maps and JSON files.
- [x] T006 [P] Add functions in `internal/birthdate` to remove a birthdate for a person identifier and keep maps and JSON files consistent.
- [x] T007 Add unit tests in `internal/birthdate` (e.g., `internal/birthdate/birthdate_edit_test.go`) for add/update/remove behavior, including invalid dates and edge cases (future dates, >120 years).
- [x] T008 Wire new `internal/birthdate` edit helpers into the existing loader lifecycle so that in-memory maps reflect changes immediately after a successful write.

**Checkpoint**: Birthdate JSON and mapping can be safely edited via Go helpers with tests in place.

---

## Phase 3: User Story 1 - Add a birthdate by clicking a person's name (Priority: P1) 🎯 MVP

**Goal**: Allow users to click a person’s name with no birthdate and add a valid date of birth that feeds the age-display feature.

**Independent Test**: On a photo where a person has no birthdate, click their name, enter a valid `YYYY-MM-DD` birthdate, save, and verify that their age appears on that and other photos using the existing age-display logic.

### Implementation for User Story 1

- [x] T009 [P] [US1] Add a JSON POST endpoint in `internal/server` (e.g., `internal/server/handlers_birthdate.go`) to handle “set birthdate” requests using the `setBirthdateRequest` and `setBirthdateResponse` contracts.
- [x] T010 [P] [US1] Implement request validation for the POST endpoint (parse `birthDate`, enforce format, reject future dates and implied ages >120 years) in `internal/server/handlers_birthdate.go`.
- [x] T011 [US1] Call the new `internal/birthdate` set/update helper from the POST handler and return localized success/error messages according to the contracts.
- [x] T012 [US1] Expose a route for the POST endpoint in the Echo router setup (e.g., `internal/server/router.go`) under an `/api/person-birthdate` path.
- [x] T013 [P] [US1] Update the relevant Templ component that renders person names (e.g., `web/templates/components/person_label.templ`) to render names as clickable elements with attributes needed to identify the person.
- [x] T014 [P] [US1] Add minimal client-side behavior (htmx or small JS) in `web/templates` / `web/assets` to POST to the new endpoint when a name is clicked and the user confirms adding a birthdate.
- [x] T015 [US1] Ensure i18n is used for any new labels, prompts, and error messages related to adding a birthdate.
- [x] T016 [US1] Add integration tests or a focused handler test in `internal/server` (e.g., `internal/server/handlers_birthdate_test.go`) that verifies a successful add flow updates the JSON files and returns the expected response.

**Checkpoint**: Users can add a birthdate by clicking a name; ages appear after the next refresh using the existing age-display feature.

---

## Phase 4: User Story 2 - Change an existing birthdate by clicking a person's name (Priority: P2)

**Goal**: Allow users to click a person’s name when a birthdate already exists and update it to a new, valid value.

**Independent Test**: For a person with an existing birthdate and visible age, click their name, change the date, save, and confirm that the age updates accordingly across photos.

### Implementation for User Story 2

- [x] T017 [P] [US2] Extend the POST “set birthdate” handler in `internal/server/handlers_birthdate.go` to handle updates when a birthdate already exists (including overwriting JSON and mapping safely).
- [x] T018 [P] [US2] Update the UI flow in the Templ component (e.g., `web/templates/components/person_label.templ`) so that clicking a name with an existing birthdate opens an edit form pre-populated with the current date.
- [x] T019 [US2] Ensure the edit flow reuses the same validation rules and error handling as the add flow, with user-friendly messages when validation fails.
- [x] T020 [US2] Add tests in `internal/server/handlers_birthdate_test.go` that cover a successful update, including JSON and in-memory map changes, and error cases (invalid new date).

**Checkpoint**: Users can correct or update existing birthdates via the same click-to-edit interaction, and ages reflect the updated date.

---

## Phase 5: User Story 3 - Remove a birthdate by clicking a person's name (Priority: P3)

**Goal**: Allow users to click a person’s name and remove their stored birthdate so that ages are no longer displayed.

**Independent Test**: For a person with a stored birthdate and age shown, click their name, choose remove, confirm, and verify that no age is displayed for that person on any photo afterward.

### Implementation for User Story 3

- [x] T021 [P] [US3] Add a JSON DELETE endpoint in `internal/server/handlers_birthdate.go` to handle `deleteBirthdateRequest`/`deleteBirthdateResponse` based on the contracts.
- [x] T022 [P] [US3] Implement the DELETE handler to call the `internal/birthdate` remove helper and treat missing entries as successful no-ops.
- [x] T023 [US3] Update the Templ UI (e.g., `web/templates/components/person_label.templ`) to provide a remove option in the click-to-edit flow, with a clear confirmation step.
- [x] T024 [US3] Ensure that after removal, any cached or in-memory state is updated so that ages stop displaying without requiring a full process restart.
- [x] T025 [US3] Add tests in `internal/server/handlers_birthdate_test.go` for the remove flow, including the no-op case and error handling.

**Checkpoint**: Users can remove birthdates via the UI, and ages stop displaying reliably after removal.

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Cross-story improvements and final validation.

- [x] T026 [P] Review all new user-facing strings (add/change/remove flows) for i18n coverage and consistency.
- [x] T027 Run `go test ./...` and address any regressions introduced by birthdate editing changes.
- [ ] T028 [P] Validate `quickstart.md` by walking through configuration and click-to-edit flows end-to-end.
- [x] T029 Perform a light code cleanup pass in `internal/birthdate`, `internal/server/handlers_birthdate.go`, and related Templ files to ensure naming and structure follow project conventions.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies – can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion – BLOCKS all user stories.
- **User Stories (Phase 3+)**: All depend on Foundational phase completion.
  - User stories can then proceed in parallel (if staffed).
  - Or sequentially in priority order (P1 → P2 → P3).
- **Polish (Final Phase)**: Depends on all desired user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) – no dependencies on other stories.
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) – builds on infrastructure from US1 but is independently testable.
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) – reuses shared handlers and UI from US1/US2 but remains independently testable.

### Within Each User Story

- Handlers and helpers should compile and pass tests before wiring into UI flows.
- UI changes should be verified manually in the kiosk for both success and error states.
- Story complete before moving to next priority for MVP-style delivery.

### Parallel Opportunities

- All tasks marked [P] can be worked on in parallel as long as they do not touch the same files.
- Once Phase 2 is complete, backend endpoint work (handlers) and frontend wiring (Templ/JS) for a given story can proceed in parallel.
- Different user stories can be implemented by different developers in parallel after Phase 2.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational (CRITICAL – blocks all stories).
3. Complete Phase 3: User Story 1 (add birthdate).
4. **STOP and VALIDATE**: Test User Story 1 independently by following `quickstart.md`.
5. Deploy/demo if ready.

### Incremental Delivery

1. Complete Setup + Foundational → foundation ready.
2. Add User Story 1 → test independently → deploy/demo (MVP).
3. Add User Story 2 → test independently → deploy/demo.
4. Add User Story 3 → test independently → deploy/demo.

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together.
2. Once Foundational is done:
   - Developer A: User Story 1 (handlers + UI).
   - Developer B: User Story 2 (update flows + tests).
   - Developer C: User Story 3 (remove flows + tests).
3. Stories complete and integrate independently.

---

## Notes

- Tasks marked [P] should avoid touching the same files concurrently.
- [Story] labels map tasks to specific user stories for traceability.
- Each user story can be independently validated using the acceptance scenarios in `spec.md`.
- Commit after each task or small group of related tasks, running `go test ./...` regularly.

