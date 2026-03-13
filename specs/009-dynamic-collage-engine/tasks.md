# Tasks: Dynamic Collage Engine

**Input**: Design documents from `/specs/009-dynamic-collage-engine/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `internal/` at repository root
- **Frontend**: `frontend/src/` for CSS
- **Templates**: `internal/templates/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the collage package structure

- [x] T001 Create `internal/collage/` package directory per plan.md structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core layout engine that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 [P] Define `LayoutResult` and `LayoutCell` types in `internal/collage/layout.go` per data-model.md
- [x] T003 Implement `ComputeLayout(assets, width, height) (LayoutResult, error)` entry point in `internal/collage/layout.go` with algorithm selection stub
- [x] T004 Implement justified-rows algorithm in `internal/collage/justified.go` (rows of similar height, aspect-ratio preservation)
- [x] T005 Add unit tests for layout types and justified algorithm in `internal/collage/layout_test.go`

**Checkpoint**: Layout engine with one algorithm ready; ProcessMemoriesCollage can call ComputeLayout

---

## Phase 3: User Story 1 - View varied, aspect-ratio-preserving collages (Priority: P1) 🎯 MVP

**Goal**: Replace fixed grid with dynamic layout; photos preserve aspect ratio; no cropping or distortion

**Independent Test**: Trigger memories (slideshow or on-demand); confirm collages preserve aspect ratios, avoid distortion; layouts adapt to photo count and aspect ratios

### Implementation for User Story 1

- [x] T006 [US1] Add `CollageLayout *collage.LayoutResult` field to `ViewData` in `internal/common/common.go`
- [x] T007 [US1] Update `ProcessMemoriesCollage` in `internal/routes/routes_asset_helpers.go`: call `collage.ComputeLayout` with 2-second timeout (context.WithTimeout); attach LayoutResult to ViewData on success; on failure/timeout leave nil for fallback
- [x] T008 [US1] Update `internal/templates/components/memories/collage.templ`: when `ViewData.CollageLayout != nil`, render cells with positions from layout (inline styles or CSS vars); when nil, use `collage-{n}` class (fallback)
- [x] T009 [US1] Update `frontend/src/css/collage.css`: change `.collage-cell img` from `object-fit: cover` to `object-fit: contain` for aspect-ratio preservation
- [x] T010 [US1] Update `internal/templates/components/memories/collage.templ`: ensure dynamic layout container uses `collage-grid collage-dynamic` and cells have `position: absolute` with computed left/top/width/height per contracts/memories-ui.md

**Checkpoint**: User Story 1 complete — dynamic collages display with aspect-ratio preservation; fallback works when engine fails

---

## Phase 4: User Story 2 - Varied layout styles from dynamic engine (Priority: P2)

**Goal**: Multiple layout algorithms (treemap, hero-plus-cluster); random selection per render yields variety

**Independent Test**: Trigger memories multiple times with same photo count; confirm different layout styles appear over time

### Implementation for User Story 2

- [x] T011 [P] [US2] Implement treemap-style packing algorithm in `internal/collage/treemap.go` (recursive subdivision, aspect-ratio-adjusted areas)
- [x] T012 [P] [US2] Implement hero-plus-cluster algorithm in `internal/collage/hero.go` (one large photo + cluster of rest)
- [x] T013 [US2] Update `ComputeLayout` in `internal/collage/layout.go`: randomly select from justified, treemap, hero; add tests for algorithm selection in `internal/collage/layout_test.go`

**Checkpoint**: User Story 2 complete — layouts vary per render across three algorithms

---

## Phase 5: User Story 3 - Seamless integration with existing memories flow (Priority: P3)

**Goal**: Timer, caption, lightbox, dismiss unchanged; only visual arrangement changes

**Independent Test**: Use memories in slideshow and via on-demand button; confirm 45s timer, caption, lightbox-on-tap, dismiss work as before

### Implementation for User Story 3

- [x] T014 [US3] Verify fallback path in `ProcessMemoriesCollage`: when `ComputeLayout` returns error or times out, `CollageLayout` is nil; template renders `collage-{n}` with `object-fit: contain`
- [x] T015 [US3] Verify `collage.templ` preserves `data-asset-id`, `role="button"`, `tabindex="0"` on cells for lightbox compatibility per contracts/memories-ui.md
- [ ] T016 [US3] Run integration test: memories in slideshow, on-demand button, lightbox tap, timer pause — no regressions

**Checkpoint**: User Story 3 complete — integration unchanged; lightbox, timer, caption work

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Tests, edge cases, validation

- [x] T017 Add table-driven tests in `internal/collage/layout_test.go` for edge cases: 3 photos, 16 photos, extreme aspect ratios (panorama, portrait), zero dimensions fallback
- [x] T018 Ensure `internal/collage` passes `golangci-lint`; run `go test ./internal/collage/...`
- [ ] T019 Run quickstart.md validation: manual test of dynamic layouts, fallback, aspect-ratio preservation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational — MVP deliverable
- **User Story 2 (Phase 4)**: Depends on US1 (needs integration in place)
- **User Story 3 (Phase 5)**: Depends on US1 (verification of integration)
- **Polish (Phase 6)**: Depends on US1–US3 complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational — No dependencies on other stories
- **User Story 2 (P2)**: Builds on US1; adds algorithms for variety
- **User Story 3 (P3)**: Verification of US1/US2 integration; can run in parallel with US2

### Within Each User Story

- Types/models before services
- Integration before verification
- Story complete before moving to next priority

### Parallel Opportunities

- T002, T004 can run in parallel (types vs algorithm)
- T011, T012 can run in parallel (treemap vs hero)
- T014, T015, T016 (US3) can overlap with T013 (US2) completion

---

## Parallel Example: User Story 2

```bash
# Launch both algorithms in parallel:
Task: "Implement treemap-style packing algorithm in internal/collage/treemap.go"
Task: "Implement hero-plus-cluster algorithm in internal/collage/hero.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (T002–T005)
3. Complete Phase 3: User Story 1 (T006–T010)
4. **STOP and VALIDATE**: Trigger memories; verify aspect-ratio preservation, dynamic layout
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Layout engine with justified algorithm
2. Add US1 → Test independently → MVP (dynamic collages with aspect-ratio preservation)
3. Add US2 → Treemap + hero algorithms → Variety across renders
4. Add US3 verification → Ensure integration unchanged
5. Polish → Tests, lint, quickstart validation

### Parallel Team Strategy

- Developer A: T002–T005 (Foundational)
- Once Foundational done:
  - Developer A: T006–T010 (US1)
  - Developer B: T011–T012 (US2 algorithms) in parallel
- After US1: Developer B integrates T013; Developer A runs US3 verification

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to user story for traceability
- User Story 4 (Research) is complete — research.md documents the approach
- Commit after each task or logical group per project workflow
- Container dimensions for ComputeLayout: use default (e.g. 1920×1080) or pass from viewport; plan defers to implementation
