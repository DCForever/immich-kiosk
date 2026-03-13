# Feature Specification: Memories Progressive Time-Range Expansion

**Feature Branch**: `008-memories-expand-range`  
**Created**: 2026-03-13  
**Status**: Draft  
**Input**: User description: "I want to update the Memories feature so that if there are no memories for today, it will expand to this week. If none for this week, it will expand to this month. I want to also update to now also be 3 to 16 photos for memories."

## Clarifications

### Session 2026-03-13

- Q: How should "this week" be defined — ISO week (Mon–Sun), calendar week (Sun–Sat), or a rolling 7-day window? → A: ISO week (Monday–Sunday): the week that contains the target date per ISO 8601.
- Q: When the week or month range has more photos than the max, how should the system select for the collage? → A: Random: pick up to 16 photos at random from the range.
- Q: Memories photo count: change from up to 8 to 3–16? → A: Yes. Memories collage shows between 3 and 16 photos. When fewer than 3 exist in the range, treat as no memories (skip/empty). When more than 16 exist, select 16 at random.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View memories with progressive fallback during slideshow (Priority: P1)

As a viewer, I see memories in the slideshow even when no photos exist for the exact date (same month and day) in any past year. The system first tries today's date, then expands to the week containing that date, then to the full month, so I still get nostalgic content instead of a skip or empty message.

**Why this priority**: Core value is maximizing the chance of showing memories rather than skipping or showing "no memories."

**Independent Test**: Run the slideshow on a date with no exact-date matches; confirm memories appear from the expanded range (week or month) and display correctly with an appropriate caption.

**Acceptance Scenarios**:

1. **Given** fewer than 3 photos exist for the exact date (same month and day) in any past year, **When** a memories moment is due in the slideshow, **Then** the system searches the week containing that date in past years and, if at least 3 photos exist, shows a memories collage (3–16 photos) from that week.
2. **Given** fewer than 3 photos exist for the exact date or the week containing that date in any past year, **When** a memories moment is due, **Then** the system searches the same month in past years and, if at least 3 photos exist, shows a memories collage (3–16 photos) from that month.
3. **Given** at least 3 photos exist for the exact date in at least one past year, **When** a memories moment is shown, **Then** the system uses the existing behavior (exact date only) and shows "Memories from X years ago" with 3–16 photos.

---

### User Story 2 - View memories with progressive fallback via on-demand button (Priority: P2)

As a viewer, I can tap the Memories button and see memories even when no photos exist for today's exact date. The same progressive expansion (today → week → month) applies, so I get content instead of "No memories for this day" whenever possible.

**Why this priority**: Ensures the on-demand experience matches the slideshow behavior and maximizes content shown.

**Independent Test**: Tap the memories button when no exact-date photos exist; confirm memories appear from the expanded range or the appropriate empty message if no photos exist in any range.

**Acceptance Scenarios**:

1. **Given** the user taps the memories button, **When** fewer than 3 photos exist for the exact date but at least 3 exist for the week or month, **Then** the system shows a memories collage (3–16 photos) from the expanded range with an appropriate caption.
2. **Given** fewer than 3 photos exist in any range (today, week, or month), **When** the user taps the memories button, **Then** the system shows the existing "No memories for this day" (or equivalent) message for 45 seconds, then returns.

---

### User Story 3 - Correct caption reflects time range used (Priority: P3)

As a viewer, I see a caption that reflects which time range was used (today, this week, or this month) so I understand the context of the memories shown.

**Why this priority**: Improves clarity and avoids confusion when expanded ranges are used.

**Independent Test**: Trigger memories with each range (exact date, week, month); confirm the caption accurately describes the range.

**Acceptance Scenarios**:

1. **Given** memories are shown from the exact date, **Then** the caption indicates "Memories from X years ago" (existing behavior).
2. **Given** memories are shown from the week range (no exact-date match), **Then** the caption indicates the expanded context (e.g. "Memories from this week, X years ago" or equivalent).
3. **Given** memories are shown from the month range (no week match), **Then** the caption indicates the month context (e.g. "Memories from [month], X years ago" or equivalent).

---

### Edge Cases

- When fewer than 3 photos exist in any range (today, week, or month): **slideshow** — skip the memories slot; **on-demand button** — show "No memories for this day" (or equivalent) for 45 seconds, then return.
- When the week or month range returns many photos: the system selects between 3 and 16 photos at random from the range for the collage (up to 16 when more than 16 exist).
- When the range returns 1 or 2 photos only: treat as insufficient; skip in slideshow or show "No memories" when opened via the on-demand button (same as when no photos exist).
- When multiple years have photos in the expanded range: the system picks one past year at random (same as existing behavior) and uses that year for the caption.
- "This week" is defined as the ISO week (Monday–Sunday) containing the target date per ISO 8601 (e.g. for March 11, the Mon–Sun week that contains March 11 in the chosen past year).
- "This month" is the same calendar month in the chosen past year (any day within that month).
- For each photo count (3–16), at least 3 layout variants exist; the system selects one at random per collage render.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST use a progressive search order for memories: (1) exact date (same month and day) in a random past year; (2) if none found, the week containing that date in a past year; (3) if none found, the same month in a past year.
- **FR-002**: System MUST select a random past year from years that have at least 3 photos in the chosen range (exact date, week, or month).
- **FR-003**: System MUST retrieve between 3 and 16 photos from the chosen range (selecting at random when more than 16 exist) and display them in the collage format. When fewer than 3 photos exist in the range, the system MUST treat it as no memories (skip in slideshow; show empty message when opened via on-demand button).
- **FR-004**: System MUST show a caption that reflects the time range used: exact date ("Memories from X years ago"), week ("Memories from this week, X years ago" or equivalent), or month ("Memories from [month], X years ago" or equivalent).
- **FR-005**: When fewer than 3 photos exist in any of the three ranges (or no photos at all), system MUST follow existing behavior: skip in slideshow; show "No memories for this day" for 45 seconds when opened via the on-demand button.
- **FR-006**: All existing memories behavior (lightbox, timer pause, dismiss, error handling) MUST remain unchanged.
- **FR-007**: For each photo count (3–16), the system MUST define at least 3 distinct collage layout variants and MUST select one at random when displaying a memories collage.

### Key Entities

- **Memory**: A memories moment with a chosen past year, time range (exact date, week, or month), 3–16 photos from that range, and display duration (45 seconds).
- **Time range**: The search scope — today (exact date), this week (7-day period containing the date), or this month (full calendar month).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users see memories in the slideshow more often when photos exist in the week or month range even when no exact-date match exists.
- **SC-002**: Each memories view remains on screen for 45 seconds before advancing, with behavior unchanged from the existing feature.
- **SC-003**: The caption accurately reflects which time range was used (exact date, week, or month).
- **SC-004**: When fewer than 3 photos exist in any range, the system behaves as before (skip or show empty message).

## Assumptions

- The existing Memories feature (spec 007) is implemented and deployed; this spec adds progressive expansion and changes the photo count from up to 8 to 3–16. Collage grid layouts must support 3–16 photos.
- "Week" is the ISO week (Monday–Sunday) containing the target date per ISO 8601.
- "Month" is the same calendar month in the chosen past year.
- Caption wording for week and month ranges can be localized; the spec provides examples for clarity.
