# Tasks: Memories (On This Day)

**Input**: Design documents from `specs/007-memories-feature/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story (US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: i18n keys and shared types needed across all stories

- [x] T001 [P] Add i18n keys for memories in `locales/en.toml` and other locale files: `memories_caption`, `memories_no_photos`, `memories_error_load`, `memories_back`
- [x] T002 [P] Add `MemoriesCollage` type and `MemoryCaption` field to `internal/common/common.go` (or define in source package) per data-model.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Provider layer and collage infrastructure that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Add `MemoriesCollage(requestID, deviceID string) (MemoriesCollage, error)` to `ProviderOps` interface in `internal/source/provider.go`
- [x] T004 Implement `MemoriesCollage` in `internal/immich/immich_memories.go`: filter memories by `Data.Year`, pick random year with assets, return up to 8 assets; handle empty case
- [x] T005 Implement `MemoriesCollage` in `internal/immich/adapter.go` to delegate to asset.MemoriesCollage
- [x] T006 Implement `MemoriesCollage` in `internal/photoprism/provider.go`: return error "memories not supported" (or add PhotosOnThisDay if API supports)
- [x] T007 [P] Create CSS grid layouts `collage-1` through `collage-8` in `frontend/src/css/collage.css` per spec Collage Layout (Grid Design)
- [x] T008 Add `MemoryCaption` to `common.ViewData` (or pass via Config) and ensure `Layout` can be `"collage"` in `internal/common/common.go` and `internal/config/config.go` if needed

**Checkpoint**: Provider returns MemoriesCollage; CSS layouts exist; ViewData supports collage

---

## Phase 3: User Story 1 - View memories during the slideshow (Priority: P1) 🎯 MVP

**Goal**: Memories collage appears in the regular slideshow rotation with 45s duration and "Memories from X years ago" caption

**Independent Test**: Run slideshow with `memories=true`; confirm collage appears in rotation, displays 45 seconds, shows caption; advances to next item

### Implementation for User Story 1

- [x] T009 [US1] Add `ProcessMemoriesCollage` in `internal/routes/routes_asset_helpers.go`: call provider.MemoriesCollage, build ViewData with layout="collage", Duration=45, Assets (1–8), MemoryCaption; handle empty → return error so caller skips
- [x] T010 [US1] In `retrieveImage` in `internal/routes/routes_asset_helpers.go`, when bucket is `SourceMemories`, call `ProcessMemoriesCollage` instead of `RandomMemoryAsset`; on empty/error, retry with different bucket or skip
- [x] T011 [US1] Add `case "collage"` in `generateViewData` in `internal/routes/routes_asset_helpers.go`: call ProcessMemoriesCollage, set viewData with collage layout and 45s duration
- [x] T012 [US1] Create `internal/templates/components/memories/collage.templ`: grid with `collage-{n}` class, caption, `data-asset-id` on cells, `data-memories="true"` on container
- [x] T013 [US1] Add collage branch in `internal/templates/components/image/layout.templ`: when layout is collage, render collage component instead of single/splitview
- [x] T014 [US1] In `internal/routes/routes_asset.go` (NewAsset handler), when viewData has layout collage, render collage component via `imageComponent.Image` or new `memoriesComponent.Collage`; ensure kiosk-data includes `duration: 45`

**Checkpoint**: Memories collage appears in slideshow; 45s timer; caption displays

---

## Phase 4: User Story 2 - Open memories on demand from the top UI (Priority: P2)

**Goal**: User can tap a Memories button to see collage on demand; can dismiss before 45s via back button

**Independent Test**: Tap memories button; collage appears; tap back before 45s → returns to slideshow; or wait 45s → returns

### Implementation for User Story 2

- [x] T015 [US2] Create `internal/routes/routes_memories.go`: `GET/POST /memories` handler that calls ProcessMemoriesCollage (or equivalent), renders collage with back button; handle empty → show "No memories for this day" screen for 45s
- [x] T016 [US2] Register `/memories` route in `main.go`
- [x] T017 [US2] Add Memories button to `internal/templates/partials/menu.templ` with HTMX trigger to load `/memories` (or POST with query params)
- [x] T018 [US2] Add back/close button to memories view in `internal/templates/components/memories/memories_view.templ` that triggers `hx-post="/asset/new"` with current queries to return to slideshow
- [x] T019 [US2] Create "No memories for this day" template/state in `internal/templates/components/memories/empty.templ` for on-demand when no photos exist

**Checkpoint**: On-demand button works; back/dismiss returns; empty state shows for on-demand

---

## Phase 5: User Story 3 - View a single memory photo in a lightbox (Priority: P3)

**Goal**: Tapping a photo in the collage opens it in a lightbox; timer pauses; closing returns to collage with remaining time

**Independent Test**: Open collage (slideshow or button), tap photo → lightbox; close → collage with remaining time; timer paused while open

### Implementation for User Story 3

- [ ] T020 [US3] Create `frontend/src/ts/memories-lightbox.ts`: openLightbox(assetId), closeLightbox(), call pausePolling/resumePolling from `polling.ts`; render overlay with `<img src="/image/{id}">` and close control
- [ ] T021 [US3] Add click handler on `.collage-cell` in collage template or via `memories-lightbox.ts` init: on tap, open lightbox with `data-asset-id`
- [ ] T022 [US3] Add lightbox overlay HTML/CSS: full-screen overlay, close button, tap-outside-to-close; include in `frontend/src/css/collage.css` or new `memories-lightbox.css`
- [ ] T023 [US3] Integrate `memories-lightbox.ts` in `frontend/src/ts/kiosk.ts` or main entry: init when `data-memories="true"` present; ensure polling pauses on open, resumes on close with remaining time
- [ ] T024 [US3] Add `data-memories` and `memories: true` to kiosk-data JSON when rendering collage so frontend enables lightbox behavior

**Checkpoint**: Lightbox works; timer pauses; close returns to collage

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Error handling, loading states, PhotoPrism support, validation

- [ ] T025 [P] Add loading state for memories fetch: show spinner in `internal/templates/partials/spinner.templ` or collage placeholder while ProcessMemoriesCollage runs
- [ ] T026 Add error handling per FR-010: when ProcessMemoriesCollage fails, show "Couldn't load memories" message; slideshow skips; on-demand returns after few seconds
- [ ] T027 [P] Add PhotoPrism `PhotosOnThisDay` in `internal/photoprism/client.go` if API supports date filter; implement MemoriesCollage in provider; or document "not supported" and ensure graceful degradation
- [ ] T028 Run quickstart.md verification checklist: collage in slideshow, 45s, caption, on-demand, back, no-photos, lightbox, error, i18n

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — BLOCKS all user stories
- **Phase 3 (US1)**: Depends on Phase 2 — MVP
- **Phase 4 (US2)**: Depends on Phase 2, 3 (reuses collage component)
- **Phase 5 (US3)**: Depends on Phase 2, 3 (lightbox on collage)
- **Phase 6 (Polish)**: Depends on Phases 3–5

### User Story Dependencies

- **US1**: After Foundational — no other story dependency
- **US2**: After US1 (reuses collage template, ProcessMemoriesCollage)
- **US3**: After US1 (lightbox attaches to collage)

### Parallel Opportunities

- T001, T002 (Setup)
- T007, T008 (Foundational)
- T015–T019 within US2 (different files)
- T020, T022 (US3)
- T025, T027 (Polish)

---

## Parallel Example: Phase 2

```bash
# T007 and T008 can run in parallel:
Task: "Create CSS grid layouts collage-1 through collage-8 in frontend/src/css/collage.css"
Task: "Add MemoryCaption to common.ViewData..."
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Run slideshow with memories=true; verify collage, 45s, caption
5. Deploy/demo

### Incremental Delivery

1. Setup + Foundational → Provider and CSS ready
2. Add US1 → Test slideshow memories → MVP
3. Add US2 → Test on-demand button
4. Add US3 → Test lightbox
5. Add Polish → Error handling, PhotoPrism

### Task Count Summary

| Phase | Tasks | Story |
|-------|-------|-------|
| 1 Setup | 2 | — |
| 2 Foundational | 6 | — |
| 3 US1 | 6 | P1 |
| 4 US2 | 5 | P2 |
| 5 US3 | 5 | P3 |
| 6 Polish | 4 | — |
| **Total** | **28** | |
