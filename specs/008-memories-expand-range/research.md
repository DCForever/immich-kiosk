# Research: Memories Progressive Time-Range Expansion

**Feature**: 008-memories-expand-range  
**Date**: 2026-03-13

## 1. ISO Week Calculation in Go

**Decision**: Use `time` package with `ISOWeek()` and manual date arithmetic. Go's `time.Time` has `ISOWeek()` returning (year, week). For "week containing date D in year Y", compute the Monday and Sunday of that ISO week in year Y.

**Rationale**: Go 1.26 `time.Time.ISOWeek()` returns `(int, int)` for year and week number. To get the date range for a given ISO week in a past year: create `time.Date(year, month, day, 0, 0, 0, 0, time.UTC)`, call `ISOWeek()` to get the week number, then compute the Monday of that week (e.g. iterate back to Monday) and the following Sunday. Alternatively, use a helper: given a date, find the Monday of its ISO week: `d.AddDate(0, 0, -int(d.Weekday())+1)` when Sunday=0 (adjust for Go's Monday=1). Actually Go uses `time.Monday=1`, so `d.Weekday()` returns 1 for Monday. To get Monday: `d.AddDate(0, 0, -int(d.Weekday()-time.Monday))`. Then Sunday = Monday + 6 days.

**Implementation**: Add `isoWeekRange(t time.Time) (start, end time.Time)` in `internal/source/display.go` (keeps date logic with source types): given `t` (e.g. March 11, 2007), return the Monday and Sunday of that week. Use `t.Truncate(24*time.Hour)` and iterate to find Monday.

**Alternatives considered**:
- Third-party ISO week library — rejected; stdlib sufficient.
- Calendar week (Sun–Sat) — rejected; spec mandates ISO week.

---

## 2. PhotoPrism Date Range Queries (Week, Month)

**Decision**: Use `after:"YYYY-MM-DD"` and `before:"YYYY-MM-DD"` filters for week and month ranges. Per [Filter Reference](https://docs.photoprism.app/user-guide/search/filters/#filter-reference): `after` finds content on or after the date; `before` finds content before the date. For a week: `after:"2007-03-05" before:"2007-03-12"` (Mon–Sun, where before is exclusive of the day after Sunday). For a month: `after:"2007-03-01" before:"2007-04-01"` (all of March 2007).

**Rationale**: PhotoPrism supports `taken`, `after`, `before` as timestamp filters. Combining `after` and `before` yields a date range. Fetch count of 30+ to allow 16 photos after filtering; shuffle and take up to 16.

**Implementation**: Build search query per range:
- **Exact date**: `taken:"YYYY-MM-DD"` (existing)
- **Week**: `after:"YYYY-MM-DD" before:"YYYY-MM-DD"` (ISO week Mon–Sun)
- **Month**: `after:"YYYY-MM-01" before:"YYYY-MM+1-01"` (same month in past year)

**Alternatives considered**:
- `year`+`month`+`day` for week — not suitable; no day range. `after`/`before` is correct.
- Multiple single-day queries for week — inefficient; one range query is better.

---

## 3. Immich Memories API for Week/Month

**Decision**: Immich `/api/memories` returns memories for a single date via `for=<date>`. For week and month, we must call the API for each day in the range and aggregate, or use `pastDays` to fetch a window. Current implementation uses `for=<today>` or `pastDays` for multi-day. For progressive expansion: (1) try `for=<exact date>` with `pastDays=0` or equivalent; (2) if empty, fetch with `pastDays=7` and filter to ISO week; (3) if still empty, fetch with `pastDays=31` and filter to same month. Alternatively, call `Memories` once per day in the week/month — would be many calls. Better: use `MemoriesWithPastDays(requestID, deviceID, 7)` for week and `MemoriesWithPastDays(..., 31)` for month, then filter the aggregated response to the target ISO week or month. Filter by `MemoryAt` or asset date within range.

**Rationale**: Immich API returns `Memory` objects; each has `MemoryAt` and `Data.Year`. For week/month we need photos whose capture date falls in the range. The existing `MemoriesWithPastDays` fetches multiple days. We can fetch 7 days (or 31), then filter `Memory` objects to those whose `MemoryAt` (or asset `LocalDateTime`) falls within the ISO week or month of the target date in the chosen year.

**Implementation**: Extend `MemoriesCollage` in `immich_memories.go`: (1) Try exact date via `Memories`; filter memories by `Data.Year` and exact month/day; require ≥3 assets. (2) If none, compute ISO week range for a random past year; call `MemoriesWithPastDays(7)` or equivalent; filter to memories in that week with ≥3 assets; pick random year. (3) If none, compute month range; call `MemoriesWithPastDays(31)`; filter to same month in past year; pick random year. Return `MemoriesCollage` with `TimeRange` set (exact, week, month).

**Alternatives considered**:
- New Immich endpoint for date range — not available.
- Single call with wide pastDays — may work; filter results to target range.

---

## 4. Collage CSS Layouts for 3–16 Photos (≥3 Variants per Count)

**Decision**: For each photo count 3–16, define **at least 3 layout variants**. Use class names `collage-{n}-a`, `collage-{n}-b`, `collage-{n}-c` (and more where extra variety is desired). The template or routes pick one variant at random when rendering.

**Rationale**: Multiple layouts per count increase visual variety across slideshow cycles. The 007 spec's "interesting grid" characteristics (asymmetrical, sub-grids, varied cell sizes) can be applied across variants. Random selection ensures each memories moment feels fresh.

**Implementation**:
- In `collage.css`: For each n in 3..16, add ≥3 layout classes (e.g. `.collage-5-a`, `.collage-5-b`, `.collage-5-c`). Each variant uses a different grid arrangement (hero placement, sub-grids, aspect emphasis).
- In `ProcessMemoriesCollage` or collage template: When building view data, pick a random variant for the current asset count (e.g. `rand.IntN(3)` or similar) and pass it to the template as `LayoutVariant`.
- Add `LayoutVariant` field to view data so the backend selects the variant; template applies `collage-{n}-{variant}` class to the grid container.

**Alternatives considered**:
- Dynamic JS layout — overkill; CSS Grid sufficient.
- Single layout per count — rejected; user requirement for ≥3 variants per count.
