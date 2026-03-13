# Tasks: Memories Progressive Time-Range Expansion

**Input**: Design documents from `specs/008-memories-expand-range/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story (US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Types, constants, and helpers needed across all stories

- [ ] T001 [P] Add `TimeRange` type (ExactDate, Week, Month), `MinCollageAssets = 3`, `MaxCollageAssets = 16` to `internal/source/display.go`
- [ ] T002 Add `TimeRange` field to `MemoriesCollage` struct in `internal/source/display.go`; update comment for 3–16 assets (depends on T001)
- [ ] T003 [P] Add `LayoutVariant` field to `ViewData` in `internal/common/common.go` for collage layout variant (e.g. "a", "b", "c")
- [ ] T004 Add `isoWeekRange(t time.Time) (start, end time.Time)` helper in `internal/source/display.go` per research.md (ISO week Mon–Sun; keeps date logic with source types)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Provider progressive search and layout variant logic that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 Implement progressive search in `internal/immich/immich_memories.go`: (1) exact date; (2) week via MemoriesWithPastDays(7) + filter; (3) month via MemoriesWithPastDays(31) + filter; require ≥3 assets; set TimeRange; change maxCollageAssets to 16
- [ ] T006 Implement progressive search in `internal/photoprism/provider.go`: (1) taken:"YYYY-MM-DD"; (2) after/before for ISO week; (3) after/before for month; require ≥3 photos; set TimeRange; change memoriesCollageMaxAssets to 16; fetch 30+ for week/month
- [ ] T007 [P] Add `memories_caption_week` and `memories_caption_month` i18n keys to `locales/en.toml` and all locale files in `locales/`
- [ ] T008 Add layout variant selection: define `CollageLayoutVariants = []string{"a","b","c"}` and `PickLayoutVariant(n int) string` (or inline in ProcessMemoriesCollage) in `internal/routes/routes_asset_helpers.go` or shared package

**Checkpoint**: Providers return MemoriesCollage with TimeRange and 3–16 assets; i18n keys exist; layout variant logic ready

---

## Phase 3: User Story 1 - View memories with progressive fallback during slideshow (Priority: P1) 🎯 MVP

**Goal**: Memories collage appears in slideshow with progressive fallback (date → week → month), 3–16 photos, correct caption, and random layout variant

**Independent Test**: Run slideshow with memories=true; when no exact-date photos exist, confirm collage from week or month appears with appropriate caption; verify 3–16 photos and varied layouts

### Implementation for User Story 1

- [ ] T009 [US1] Update `ProcessMemoriesCollage` in `internal/routes/routes_asset_helpers.go`: read `collage.TimeRange`; select i18n key for caption (`memories_caption`, `memories_caption_week`, `memories_caption_month`); interpolate month name for month caption; pick random `LayoutVariant` for asset count; validate 3 ≤ len(Assets) ≤ 16; pass LayoutVariant to ViewData
- [ ] T010 [US1] Update `internal/templates/components/memories/collage.templ`: use `collage-{n}-{LayoutVariant}` class when LayoutVariant is set, else fallback to `collage-{n}` for backward compatibility; in 008 flow n is always 3–16 (provider returns empty when fewer than 3)
- [ ] T011 [US1] Add ≥3 CSS layout variants per count 3–16 in `frontend/src/css/collage.css`: `.collage-3-a`, `.collage-3-b`, `.collage-3-c` through `.collage-16-a`, `.collage-16-b`, `.collage-16-c` (42+ layout classes total)

**Checkpoint**: Slideshow shows memories with progressive fallback; caption reflects TimeRange; layout variants apply; 3–16 photos display correctly

---

## Phase 4: User Story 2 - View memories with progressive fallback via on-demand button (Priority: P2)

**Goal**: On-demand Memories button shows collage with same progressive fallback, caption, and layout behavior as slideshow

**Independent Test**: Tap memories button when no exact-date photos exist; confirm collage from week or month with correct caption; tap back to return

### Implementation for User Story 2

- [ ] T012 [US2] Verify `internal/routes/routes_memories.go` passes ViewData from ProcessMemoriesCollage (including LayoutVariant) to template; ensure empty/error states still show "No memories for this day" for 45s when <3 photos in all ranges

**Checkpoint**: On-demand button uses same ProcessMemoriesCollage; progressive fallback and captions work; empty state correct

---

## Phase 5: User Story 3 - Correct caption reflects time range used (Priority: P3)

**Goal**: Caption accurately indicates whether memories are from exact date, week, or month

**Independent Test**: Trigger memories with each range (mock or real data); confirm caption text matches TimeRange (exact: "X years ago"; week: "this week, X years ago"; month: "[month], X years ago")

### Implementation for User Story 3

- [ ] T013 [US3] Ensure month caption uses localized month name in `ProcessMemoriesCollage` in `internal/routes/routes_asset_helpers.go` (e.g. via i18n or time.Month.String); verify `memories_caption_month` interpolation in locales

**Checkpoint**: All three caption variants display correctly; month name localized

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validation, tests, and final verification

- [ ] T014 [P] Add unit tests for `isoWeekRange` in `internal/source/` (table-driven: various dates, verify Mon–Sun range)
- [ ] T015 [P] Add unit tests for `MemoriesCollage` progressive search in `internal/immich/` and `internal/photoprism/` (mock API; verify TimeRange and 3–16 bounds)
- [ ] T016 Run `go test ./...` and fix any regressions
- [ ] T017 Run quickstart.md verification checklist: progressive search, 3–16 photos, captions, layout variants, i18n, Immich/PhotoPrism support; verify lightbox, timer pause, dismiss, error handling unchanged (FR-006); verify memories fetch and collage render within 3 seconds

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational — core implementation
- **User Story 2 (Phase 4)**: Depends on US1 (ProcessMemoriesCollage) — verification only
- **User Story 3 (Phase 5)**: Depends on US1 — caption refinement
- **Polish (Phase 6)**: Depends on US1, US2, US3

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational — No dependencies on other stories
- **User Story 2 (P2)**: Uses ProcessMemoriesCollage from US1 — verification that on-demand route passes through correctly
- **User Story 3 (P3)**: Caption logic in ProcessMemoriesCollage — ensure month name localization

### Within Each Phase

- T001, T003, T004 can run in parallel (Setup); T002 depends on T001 (same file)
- T005 and T006 are sequential (different providers, no file conflict)
- T007 and T008 can run in parallel with T005/T006 (Foundational)
- T009–T011 in US1 are sequential (routes → template → CSS)
- T014 and T015 can run in parallel (Polish)

### Parallel Opportunities

- Setup: T001, T003 in parallel (T002 after T001)
- Foundational: T007, T008 in parallel with T005/T006
- Polish: T014, T015 in parallel

---

## Parallel Example: Phase 1 Setup

```bash
# Launch Setup tasks in parallel (T001, T003):
Task: "Add TimeRange type and constants to internal/source/display.go"
Task: "Add LayoutVariant field to ViewData in internal/common/common.go"
# Then T002 (depends on T001): Add TimeRange field to MemoriesCollage
# Then T004: Add isoWeekRange helper to internal/source/display.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Run slideshow; verify progressive fallback, 3–16 photos, captions, layout variants
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Providers return TimeRange and 3–16 assets
2. Add US1 → Slideshow with progressive fallback (MVP)
3. Verify US2 → On-demand button works (no code change if ProcessMemoriesCollage is shared)
4. Refine US3 → Month caption localization if needed
5. Polish → Tests and verification

### Task Count Summary

| Phase | Tasks | Story |
|-------|-------|-------|
| Phase 1: Setup | T001–T004 | — |
| Phase 2: Foundational | T005–T008 | — |
| Phase 3: US1 | T009–T011 | P1 |
| Phase 4: US2 | T012 | P2 |
| Phase 5: US3 | T013 | P3 |
| Phase 6: Polish | T014–T017 | — |
| **Total** | **17** | — |

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- US2 and US3 are largely satisfied by US1 implementation; verification/refinement tasks
- Collage layouts 1–2 from 007 may need handling: 008 uses 3–16 only; when 1–2 photos exist, treat as empty
- Commit after each task or logical group per project workflow
