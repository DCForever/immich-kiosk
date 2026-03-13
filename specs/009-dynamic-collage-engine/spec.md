# Feature Specification: Dynamic Collage Engine

**Feature Branch**: `009-dynamic-collage-engine`  
**Created**: 2026-03-13  
**Status**: Draft  
**Input**: User description: "Replace memories collage with dynamic collage engine that adapts layouts algorithmically, preserves aspect ratios, supports multiple layout styles, integrates with memories flow, includes research/evaluation phase, and delivers implementation"

## Clarifications

### Session 2026-03-13

- Q: When the collage engine produces output, how must lightbox-on-tap work? → A: Engine MUST output a format that allows per-photo tap (e.g. HTML with individual photo elements); PNG-only is out of scope.
- Q: When the collage engine fails or times out, what should the user see? → A: Simple grid (degraded but functional) — user still sees the photos in a basic layout.
- Q: How should "3 distinct layout variants per photo count" be interpreted? → A: Since the engine is dynamic and fluid each time, no requirement for 3 variants anymore; variety comes from the algorithmic nature of the engine.
- Q: What is the maximum acceptable time for collage generation before falling back to a simple grid? → A: 2 seconds.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View varied, aspect-ratio-preserving collages in memories (Priority: P1)

As a viewer, I see memories collages that adapt to the number of photos (3–16) and to the actual aspect ratios of the images. Photos may be cropped up to 10% from each end of the longer dimension to improve fit; no distortion is permitted. Layouts feel organic and varied rather than uniform grids.

**Why this priority**: Core value is replacing rigid grid layouts with visually pleasing, adaptive arrangements that respect each photo's proportions.

**Independent Test**: Trigger memories (slideshow or on-demand) with different photo counts and mixed aspect ratios; confirm collages preserve aspect ratios, avoid distortion, and produce varied layouts. Delivers value without changing how memories are accessed.

**Acceptance Scenarios**:

1. **Given** 3–16 photos are selected for a memories collage, **When** the collage is displayed, **Then** each photo maintains its aspect ratio within the allowed 10% crop range (no stretching; up to 10% cropping from each end of the longer dimension is permitted).
2. **Given** photos with mixed aspect ratios (portrait, landscape, square), **When** the collage is rendered, **Then** the layout accommodates them with appropriate sizing; gaps or spacing between images are acceptable.
3. **Given** the same number of photos on different occasions, **When** collages are rendered, **Then** layouts can vary (not always identical for the same count).

---

### User Story 2 - Varied layout styles from dynamic engine (Priority: P2)

As a viewer, I experience different visual styles when viewing memories collages. The dynamic engine produces varied arrangements (e.g., justified rows, treemap-style packing, hero-plus-cluster, or randomized organic layouts) that adapt each time; variety comes from the algorithmic, fluid nature of the engine rather than from a fixed set of per-count variants.

**Why this priority**: Increases visual variety and keeps the experience fresh across repeated views.

**Independent Test**: Trigger memories multiple times with the same photo count; confirm that different layout styles appear over time. Can be tested by observing several collage renders.

**Acceptance Scenarios**:

1. **Given** a photo count between 3 and 16, **When** a collage is rendered, **Then** the dynamic engine produces an arrangement that adapts to the count and image aspect ratios.
2. **Given** multiple collage renders with the same photo count, **When** observed over time, **Then** layouts vary (the engine's algorithmic nature produces different arrangements per render).
3. **Given** layout styles such as justified rows, treemap packing, hero-plus-cluster, or organic arrangements, **When** the engine supports them, **Then** such variety is achievable through its dynamic layout logic.

---

### User Story 3 - Seamless integration with existing memories flow (Priority: P3)

As a viewer, I experience no change in how I access or interact with memories. The collage appears in the slideshow, via the on-demand button, with the 45-second timer, caption, and lightbox-on-tap behavior unchanged. Only the visual arrangement of photos changes.

**Why this priority**: Ensures the dynamic collage engine is a drop-in replacement without disrupting the existing user experience.

**Independent Test**: Use memories in slideshow and via on-demand button; confirm timer, caption, lightbox, and dismiss behavior work as before. Delivers value by preserving all existing interactions.

**Acceptance Scenarios**:

1. **Given** the slideshow is running, **When** a memories moment is shown, **Then** the dynamic collage displays with the same 45-second duration, caption, and advancement behavior as before.
2. **Given** the user taps the memories button, **When** a collage is shown, **Then** the on-demand flow (45-second timer, dismiss, caption) works identically to the previous implementation.
3. **Given** a collage is on screen, **When** the user taps a photo, **Then** the lightbox opens and the timer pauses; closing the lightbox returns to the collage with remaining time.

---

### User Story 4 - Research and evaluation before implementation (Priority: P4)

As a product owner, I want the project to evaluate different approaches for achieving dynamic collages before implementation. The evaluation must consider options that support dynamic layout, aspect-ratio preservation, and integration with the existing kiosk architecture. The chosen approach must be feasible for the project context.

**Why this priority**: Reduces risk by selecting an appropriate approach before committing to implementation.

**Independent Test**: Review the specification and planning artifacts; confirm that an evaluation phase is defined and that the chosen approach meets feasibility criteria (maturity, maintainability, output format, performance).

**Acceptance Scenarios**:

1. **Given** the feature is planned, **When** implementation begins, **Then** an evaluation of collage-generation approaches has been completed and a chosen approach documented.
2. **Given** the evaluation, **When** an approach is selected, **Then** it supports dynamic layout, aspect-ratio preservation, and integration with the kiosk (photo sources, web frontend).
3. **Given** the chosen approach, **When** assessed, **Then** it is feasible in terms of maturity, maintainability, per-photo-tappable output format (e.g. HTML), and performance.

---

### Edge Cases

- When fewer than 3 photos exist: the memories feature treats this as no memories (skip or empty message); the collage engine is not invoked.
- When exactly 3 photos exist: the collage engine produces a layout for 3 photos.
- When 16 photos exist: the collage engine produces a layout for 16 photos.
- When photos have extreme aspect ratios (e.g., very wide panoramas or very tall portraits): layouts accommodate them without distortion; gaps or spacing are acceptable.
- When the collage engine fails or exceeds 2 seconds: the system falls back to a simple grid layout so the user still sees the photos; the memories flow is not blocked.
- When the same memories are requested twice in quick succession: each render may produce a different layout (dynamic engine produces varied arrangements).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST replace the current fixed CSS grid collage implementation with a dynamic collage engine that generates arrangements algorithmically.
- **FR-002**: System MUST adapt collage layouts to the number of photos (3–16 per spec 008) and to the actual aspect ratios of the images.
- **FR-003**: System MUST preserve aspect ratios for all photos within an allowed range; up to 10% cropping from each end of the longer dimension is permitted; no distortion is permitted. Images MUST be scaled to fit while maintaining proportions within this range; gaps or spacing between images are acceptable.
- **FR-004**: System MUST produce varied layouts per collage render; the dynamic, algorithmic nature of the engine MUST yield different arrangements over time (no fixed per-count variant requirement).
- **FR-005**: System MUST support multiple layout styles (e.g., justified rows, treemap-style packing, hero-plus-cluster, randomized organic layouts) through its dynamic layout logic.
- **FR-006**: System MUST integrate with the existing memories flow: collage in slideshow, on-demand button, 45-second timer, caption, and lightbox-on-tap. User-facing behavior MUST remain unchanged except for the visual arrangement.
- **FR-006a**: System MUST output collage content in a format that allows per-photo tap for lightbox (e.g. HTML with individual photo elements). Single pre-rendered image (PNG) output is out of scope because it cannot support per-photo interaction.
- **FR-007**: System MUST include a research and evaluation phase before implementation. The evaluation MUST consider approaches that support dynamic layout, aspect-ratio preservation, per-photo-tappable output (e.g. HTML), integration with the kiosk, and generation within 2 seconds. The chosen approach MUST be feasible (maturity, maintainability, output format, performance).
- **FR-008**: System MUST deliver the chosen implementation and integrate it so collages are generated dynamically when memories are requested (slideshow or on-demand).
- **FR-009**: When the collage engine fails or cannot produce output within 2 seconds, system MUST fall back to a simple grid layout so the user still sees the photos in a basic arrangement; the memories flow MUST NOT be blocked by an error message or empty state.

### Key Entities

- **Collage**: A visual arrangement of 3–16 photos for display in the memories view, with layout determined by the dynamic engine.
- **Layout**: An arrangement produced by the dynamic engine for a given set of photos; layouts vary per render due to the engine's algorithmic nature.
- **Photo**: An asset from the photo source with dimensions and aspect ratio used by the collage engine for layout.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Collages preserve aspect ratios for all photos within the allowed 10% crop range; no distortion is observable.
- **SC-002**: Over multiple renders, different layouts appear; the dynamic engine produces varied arrangements.
- **SC-003**: Integration with the memories flow is unchanged: slideshow, on-demand button, 45-second timer, caption, and lightbox-on-tap work as before.
- **SC-004**: Layouts adapt to photo count and aspect ratios; collages feel organic and varied rather than uniform grids.
- **SC-005**: An evaluation of collage-generation approaches is completed before implementation, with a chosen approach documented and feasible for the project.
- **SC-006**: When the collage engine fails or exceeds 2 seconds, the memories flow continues with a simple grid fallback; the user always sees the photos.
- **SC-007**: Collage generation completes within 2 seconds under normal conditions; the evaluation phase MUST consider this performance target when selecting an approach.

## Assumptions

- The Memories feature (specs 007 and 008) is implemented; this spec replaces only the collage layout mechanism, not the memories selection, timer, caption, or lightbox logic.
- Photo count for memories is 3–16 per spec 008; the collage engine receives this range.
- The collage output MUST allow per-photo tap for lightbox; therefore output must be a format with individual photo elements (e.g. HTML/CSS), not a single pre-rendered image. The evaluation phase will select an approach that supports this.
- The kiosk backend can invoke the collage engine (e.g., via subprocess, service, or native integration); the evaluation phase will determine the integration pattern.
- Collage generation MUST complete within 2 seconds under normal conditions; exceeding this triggers the simple-grid fallback (FR-009). The evaluation phase MUST consider this performance target.
