# Quickstart: Memories Progressive Time-Range Expansion

**Feature**: 008-memories-expand-range  
**Date**: 2026-03-13

## Prerequisites

- Go 1.26
- 007-memories-feature implemented (MemoriesCollage, collage template, lightbox)
- Immich or PhotoPrism with date-search support

## Implementation Order

1. **Source package: MemoriesCollage and TimeRange**
   - Add `TimeRange` type (ExactDate, Week, Month) to `internal/source/display.go`.
   - Add `MinCollageAssets = 3`, `MaxCollageAssets = 16` constants.
   - Add `TimeRange` field to `MemoriesCollage` struct.
   - Update validation: 3 ≤ len(Assets) ≤ 16.

2. **Provider: Immich**
   - In `immich_memories.go`, implement progressive search:
     - (1) Exact date: existing logic; require ≥3 assets; set `TimeRange=ExactDate`.
     - (2) Week: compute ISO week range for past years; use `MemoriesWithPastDays(7)` or iterate; filter to week; require ≥3; set `TimeRange=Week`.
     - (3) Month: use `MemoriesWithPastDays(31)`; filter to same month in past year; require ≥3; set `TimeRange=Month`.
   - Change `maxCollageAssets` from 8 to 16; require ≥3 before returning.
   - Add `isoWeekRange(t time.Time) (start, end time.Time)` helper in `internal/source/display.go`.

3. **Provider: PhotoPrism**
   - In `provider.go`, implement progressive search:
     - (1) Exact date: existing `taken:"YYYY-MM-DD"`; require ≥3; set `TimeRange=ExactDate`.
     - (2) Week: `after:"Mon" before:"Sun+1"` for ISO week in past year; require ≥3; set `TimeRange=Week`.
     - (3) Month: `after:"YYYY-MM-01" before:"YYYY-MM+1-01"`; require ≥3; set `TimeRange=Month`.
   - Change `memoriesCollageMaxAssets` to 16; fetch count 30+ for week/month; shuffle, take up to 16.
   - Require ≥3 photos before returning; else `ErrMemoriesEmpty`.

4. **Routes**
   - In `routes_asset_helpers.go`, `ProcessMemoriesCollage`: read `collage.TimeRange`; select i18n key for caption (`memories_caption`, `memories_caption_week`, `memories_caption_month`); pass to template.
   - Ensure `generateViewData` passes `TimeRange` or caption string to view data.

5. **Templates**
   - In `collage.templ`: support `collage-{n}-{variant}` class based on `len(Assets)` and `LayoutVariant` from view data.
   - Pass caption from view data (already interpolated).

6. **Routes**
   - In `ProcessMemoriesCollage`: for the current asset count (3–16), pick a random layout variant (e.g. "a", "b", "c") and set `LayoutVariant` in view data.

7. **CSS**
   - In `collage.css`: for each count 3–16, add **at least 3 layout variants** (e.g. `collage-5-a`, `collage-5-b`, `collage-5-c`). Each variant uses a different grid arrangement (hero placement, sub-grids, etc.).

8. **i18n**
   - Add `memories_caption_week` (e.g. "Memories from this week, %d years ago").
   - Add `memories_caption_month` (e.g. "Memories from %s, %d years ago" — month name + years).
   - Update `locales/*.toml` for all locales.

## Testing

- Unit: `MemoriesCollage` progressive search (exact → week → month); ISO week helper; 3–16 bounds; mock providers.
- Unit: Caption selection by `TimeRange`.
- Integration: Run kiosk; verify collage with 3–16 photos; verify week/month captions when exact date has no photos.
- Manual: Date with no exact-date photos → week caption; no week photos → month caption; <3 in all → empty.

## Verification Checklist

- [ ] Progressive search: exact date → week → month when <3 photos.
- [ ] Collage shows 3–16 photos; <3 triggers empty state.
- [ ] Caption reflects TimeRange (exact, week, month).
- [ ] ISO week used for "this week" (Mon–Sun).
- [ ] For each count 3–16, ≥3 layout variants exist; random variant selected per render.
- [ ] Collage layouts 3–16 with variants render correctly.
- [ ] i18n for `memories_caption_week`, `memories_caption_month`.
- [ ] Immich and PhotoPrism both support progressive search.
- [ ] Existing behavior (lightbox, 45s, back button) unchanged.
