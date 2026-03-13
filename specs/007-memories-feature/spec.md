# Feature Specification: Memories (On This Day)

**Feature Branch**: `007-memories-feature`  
**Created**: 2026-03-11  
**Status**: Draft  
**Input**: User description: "I want to add a memories feature. I want it to grab a random set of photos from a random year of the same month and day. For example, if today is 2026-03-11, it might grab 2007-03-11 photos. I want to support multiple photos in a collage. It should use an interesting grid depending on how many photos up to 8. The memory should have its own timer of 45 seconds. It should be a part of the regular slideshow and also available as its own button in the top UI. The UI in the bottom should say 'Memories from X years ago'. Tapping on any 1 photo in the collage will show it in a lightbox."

## Clarifications

### Session 2026-03-11

- Q: When no photos exist for the chosen date in any past year, should the system skip the memories slot or show a dedicated message screen? → A: For slideshow skip (do not show a memories screen; advance to next item). For on-demand button show a dedicated screen with a short message (e.g. "No memories for this day") for 45 seconds, then return.
- Q: While the user has a photo open in the lightbox, should the 45-second timer pause or keep running? → A: Pause the timer while the lightbox is open; on close show the collage again with the remaining time.
- Q: For "random past year," should the system consider all years with photos for that month/day or only the last N years? → A: All years that have at least one photo for that month/day; pick one of those years at random.
- Q: When memories are shown via the on-demand button, can the user dismiss before 45 seconds? → A: Yes; provide a way to dismiss (e.g. back or close) so the user can return before 45 seconds.
- Q: If the photo source is slow or unavailable when loading memories, what should happen? → A: Show a loading state while fetching; on failure show a short error message (e.g. "Couldn't load memories"), then skip (slideshow) or return (button) after a few seconds.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View memories during the slideshow (Priority: P1)

As a viewer, I see a "memories" moment as part of the normal slideshow. The system picks a random past year with the same month and day as today, selects up to 8 photos from that date, and shows them in a collage. The screen displays for 45 seconds and shows the caption "Memories from X years ago" at the bottom.

**Why this priority**: Core value is surfacing nostalgic content automatically in the main flow.

**Independent Test**: Run the slideshow; confirm a memories collage appears in rotation, displays for 45 seconds, and shows the correct caption. Delivers value without any extra user action.

**Acceptance Scenarios**:

1. **Given** the slideshow is running, **When** a memories moment is shown, **Then** the user sees a collage of up to 8 photos from the same month and day in a randomly chosen past year, with a bottom caption "Memories from X years ago" and the view advances after 45 seconds.
2. **Given** a memories moment is on screen, **When** 45 seconds elapse, **Then** the slideshow advances to the next item (e.g. next photo or next memory).

---

### User Story 2 - Open memories on demand from the top UI (Priority: P2)

As a viewer, I can tap a dedicated "Memories" (or equivalent) button in the top UI to immediately see a memories collage. The same rules apply: random past year, same month and day, up to 8 photos in a grid, 45-second timer, and "Memories from X years ago" in the bottom UI.

**Why this priority**: Lets users choose when to see memories without waiting for the slideshow.

**Independent Test**: Tap the memories button; confirm a memories collage appears with caption and 45-second duration. Can be tested without the slideshow.

**Acceptance Scenarios**:

1. **Given** the kiosk is on the main view or slideshow, **When** the user taps the memories button in the top UI, **Then** a memories collage is shown with the same behavior as in the slideshow (grid, caption, 45-second timer).
2. **Given** the memories view is showing (opened via the button), **When** 45 seconds elapse or the user dismisses (e.g. back/close), **Then** the user returns to the previous context (e.g. slideshow or main view).

---

### User Story 3 - View a single memory photo in a lightbox (Priority: P3)

As a viewer, I can tap any one photo in the memories collage to open it in a lightbox (full-screen or overlay) for closer viewing. Closing the lightbox returns to the collage (or the next slideshow step if the 45 seconds have passed).

**Why this priority**: Enhances engagement by allowing focus on one photo without changing the overall feature.

**Independent Test**: Open a memories collage (via slideshow or button), tap one photo; confirm it opens in a lightbox and can be closed to return. Delivers value independently of how the collage was opened.

**Acceptance Scenarios**:

1. **Given** a memories collage is on screen, **When** the user taps one photo, **Then** that photo is shown in a lightbox (full-screen or overlay).
2. **Given** a photo is open in the lightbox, **When** the user closes the lightbox (e.g. tap outside or close control), **Then** the user sees the memories collage again (if still within the 45-second window) or the next slideshow content.

---

### Edge Cases

- When no photos exist for the chosen date in any past year: **slideshow** — skip the memories slot and advance to the next item; **on-demand button** — show a dedicated screen with a short message (e.g. "No memories for this day") for 45 seconds, then return to previous context.
- When only one photo exists for that date: the collage shows a single photo in the grid (grid layout supports 1–8 photos).
- When 2–7 photos exist: the grid uses a layout appropriate to the count; see Collage Layout (Grid Design) for layout characteristics.
- When the user opens the lightbox: the 45-second timer pauses; when the user closes the lightbox, the collage is shown again with the remaining time.
- When loading memories is slow or the photo source fails: show a loading state while fetching; on failure show a short error message (e.g. "Couldn't load memories"), then skip the slot (slideshow) or return to previous context (on-demand button) after a few seconds.

## Collage Layout (Grid Design)

The memories collage uses an *interesting grid* — varied, asymmetrical layouts rather than uniform grids. Layouts are selected based on the number of photos (1–8). Reference examples illustrate the intended visual style.

### Layout Characteristics

- **Varied cell sizes**: Mix of large, medium, and small cells; some photos are prominent (hero), others grouped in smaller clusters.
- **Asymmetrical composition**: Layouts are not symmetric; they create visual interest through uneven distribution and varied proportions.
- **Mixed orientations**: Cells may be portrait (tall), landscape (wide), or square; the layout accommodates different aspect ratios.
- **Sub-grids**: For higher photo counts, groups of 2–4 smaller cells may form sub-grids (e.g. 2×2) within the overall collage.
- **Filled space**: No empty gaps; the collage fills the display area with thin borders or spacing between cells.
- **Layout per count**: Each photo count (1–8) has at least one defined layout; the system selects a layout appropriate to the number of photos retrieved.

### Reference

A design reference image (6 example layouts) illustrates the intended collage style. Cell counts in the reference may exceed 8; for memories, layouts are defined for 1–8 photos. The following table summarizes the layout styles:

| Style | Description |
|------|-------------|
| Asymmetrical with sub-grid | Tall/portrait cells, large squares, and a 2×2 sub-grid of smaller cells. |
| Balanced mix | Portrait and landscape blocks with stacked small squares and rectangles. |
| Landscape emphasis | Wide top block, medium cells, and a 2×2 sub-grid at bottom-right. |
| Vertical stack | Portrait stacks, landscape blocks, and a vertical stack of small squares. |
| Fewer, larger cells | Wide top span, prominent portrait, large landscape, vertical square stack (good for 5–8 photos). |
| Dense and varied | Many cells with portrait, landscape, and 2×2 sub-grid combinations. |

Implementation defines concrete layouts for each count (1–8) following these characteristics.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST select a random past year (same month and day as today) for each memories moment; the year MUST be chosen from all years that have at least one photo for that month and day (no cap).
- **FR-002**: System MUST retrieve up to 8 photos taken on that date (same month and day) from the photo source.
- **FR-003**: System MUST display those photos in a collage using an interesting grid layout that varies by the number of photos (1–8); see Collage Layout (Grid Design) for layout characteristics.
- **FR-004**: System MUST show each memories view for 45 seconds before advancing (when used in slideshow or from the button).
- **FR-005**: System MUST show the caption "Memories from X years ago" in the bottom UI during the memories view (X = current year minus selected year).
- **FR-006**: System MUST include memories as one of the items in the regular slideshow rotation.
- **FR-007**: System MUST provide a dedicated button in the top UI that, when activated, shows a memories collage (same behavior as in slideshow); when opened via this button, the user MUST be able to dismiss before 45 seconds (e.g. back or close) and return to the previous context.
- **FR-008**: System MUST allow the user to tap one photo in the collage to open it in a lightbox; the 45-second timer MUST pause while the lightbox is open; closing the lightbox returns the user to the collage with the remaining time (or next content if the full 45 seconds have elapsed after closing).
- **FR-009**: When no photos exist for the chosen date in any past year: system MUST skip the memories slot in the slideshow (advance to next item); when memories are opened via the top UI button, system MUST show a dedicated screen with a short message (e.g. "No memories for this day") for 45 seconds, then return to previous context.
- **FR-010**: When loading memories is slow or fails (e.g. photo source unavailable): system MUST show a loading state while fetching; on failure MUST show a short error message (e.g. "Couldn't load memories"), then skip the memories slot (slideshow) or return to previous context (on-demand button) after a few seconds.

### Key Entities

- **Memory**: A single memories moment: a chosen past date (same month and day as today), the set of photos for that date (up to 8), and the display duration (45 seconds).
- **Photo**: An asset from the photo source, with a capture date used to match "on this day" (month and day).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users see at least one memories moment during a typical slideshow session when photos exist for that day in a past year.
- **SC-002**: Each memories view remains on screen for 45 seconds before advancing, unless the user interacts (e.g. opens lightbox).
- **SC-003**: Users can open memories on demand via the top UI and see a collage with the correct "Memories from X years ago" caption.
- **SC-004**: Users can tap any photo in the collage to view it in a lightbox and return without error.
- **SC-005**: Collage layout adapts clearly and readably to 1–8 photos (no broken or overlapping layout).

## PhotoPrism Support

- **Memories are not supported** when using PhotoPrism as the photo source. The provider returns `ErrMemoriesNotSupported`; the UI shows "Couldn't load memories" and returns after 5 seconds. Graceful degradation per research.md: API date filter for "on this day" was not confirmed; fallback to "not supported" until API is verified.

## Assumptions

- "Same month and day" is based on the photo’s capture date (or the date used by the existing photo source for "on this day" logic).
- The slideshow already has a concept of "items" or "steps" so a memories moment can be one item with a 45-second duration.
- The top UI has space or can accommodate a new button for memories without removing required controls.
- Lightbox close action (tap outside, button, or gesture) follows existing kiosk patterns where applicable.
