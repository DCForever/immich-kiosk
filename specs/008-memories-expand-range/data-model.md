# Data Model: Memories Progressive Time-Range Expansion

**Feature**: 008-memories-expand-range  
**Date**: 2026-03-13

## Overview

This feature extends the existing Memories data model with a **TimeRange** enum and updates **MemoriesCollage** to support 3–16 assets and progressive search. No new entities; modifications to existing types.

## Entity Changes

### MemoriesCollage (modified)

| Field | Type | Description |
|-------|------|-------------|
| `Year` | int | The chosen past year (e.g. 2007) |
| `Assets` | []DisplayAsset | 3–16 photos from the chosen range |
| `YearsAgo` | int | Current year minus Year (for caption) |
| `TimeRange` | TimeRange | Which range was used: ExactDate, Week, or Month |

- **Validation**: 3 ≤ len(Assets) ≤ 16; Year < current year.
- **Lifecycle**: Built per request; not persisted. Cached in view cache if prefetch is used.

### TimeRange (new enum)

| Value | Description | Caption example |
|-------|-------------|-----------------|
| `ExactDate` | Same month and day as today | "Memories from X years ago" |
| `Week` | ISO week containing the date | "Memories from this week, X years ago" |
| `Month` | Same calendar month | "Memories from [month], X years ago" |

- **Location**: `internal/source/display.go`
- **Usage**: Provider sets `TimeRange` when building `MemoriesCollage`; routes use it to select i18n key and interpolate caption.

### Constants (new)

| Constant | Value | Location |
|----------|-------|----------|
| `MinCollageAssets` | 3 | `internal/source/display.go` or provider |
| `MaxCollageAssets` | 16 | `internal/source/display.go` or provider |

## State Transitions

### Progressive Search Flow

1. **Exact date**: Query for same month/day in past years. If ≥3 assets in any year → pick random year, return with `TimeRange=ExactDate`.
2. **Week fallback**: If exact date yields <3 assets across all years, query ISO week range in past years. If ≥3 assets in any year → pick random year, return with `TimeRange=Week`.
3. **Month fallback**: If week yields <3 assets, query same month in past years. If ≥3 assets in any year → pick random year, return with `TimeRange=Month`.
4. **Empty**: If all three ranges yield <3 assets → return `ErrMemoriesEmpty`.

### Asset Selection

- When range has 3–16 photos: use all.
- When range has >16 photos: shuffle, take 16 at random.
- When range has 1–2 photos: treat as empty (continue to next range or return empty).

## Relationships

```
MemoriesCollage
  ├── TimeRange → ExactDate | Week | Month
  └── Assets[] → DisplayAsset (3–16)

ViewData (memories)
  └── Assets[] → ViewImageData (3–16)
  └── MemoryCaption → string (from i18n based on TimeRange)
  └── LayoutVariant → string (e.g. "a", "b", "c" — random choice for collage-{n}-{variant})
  └── Duration = 45
  └── Layout = "collage"
```

### Layout Variant Selection

- **ProcessMemoriesCollage**: When building view data, for asset count n (3–16), pick a random layout variant (e.g. from `{"a","b","c"}` or more) and set `LayoutVariant`. Template applies `collage-{n}-{LayoutVariant}` class.

## Data Flow

1. **Immich**: Progressive search: `Memories` (exact) → `MemoriesWithPastDays(7)` (week) → `MemoriesWithPastDays(31)` (month). Filter by date range; require ≥3 assets per year; pick random year.
2. **PhotoPrism**: Progressive search: `taken:"YYYY-MM-DD"` → `after/before` (week) → `after/before` (month). Require ≥3 photos; pick random year; shuffle and take up to 16.
3. **Provider**: `MemoriesCollage` returns `MemoriesCollage` with `TimeRange` and 3–16 assets.
4. **Routes**: `ProcessMemoriesCollage` uses `TimeRange` to select caption i18n key; passes to template.
