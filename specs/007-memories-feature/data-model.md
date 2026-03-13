# Data Model: Memories (On This Day)

**Feature**: 007-memories-feature  
**Date**: 2026-03-11

## Overview

The memories feature reuses existing `DisplayAsset` and `ViewImageData` types. New concepts are the **MemoriesCollage** (a collection of up to 8 assets for one year) and **MemoriesViewState** (client-side state for lightbox and timer).

## Entities

### Memory (Immich — existing)

- **Source**: `internal/immich/immich.go`
- **Fields**: `ID`, `MemoryAt`, `Type`, `Assets`, `Data.Year`
- **Use**: Each Memory represents "on this day" for a specific year. `Assets` holds photos for that date.

### MemoriesCollage (new — in-memory only)

Represents a single memories moment for display.

| Field | Type | Description |
|-------|------|-------------|
| `Year` | int | The chosen past year (e.g. 2007) |
| `Assets` | []DisplayAsset | Up to 8 photos from that year's month-day |
| `YearsAgo` | int | Current year minus Year (for caption "X years ago") |

- **Lifecycle**: Built per request; not persisted. Cached in view cache if prefetch is used.
- **Validation**: 1 ≤ len(Assets) ≤ 8; Year < current year.

### DisplayAsset (existing)

- **Source**: `internal/source/display.go`
- **Reuse**: No changes. Each collage cell displays one DisplayAsset.

### ViewImageData (existing)

- **Source**: `internal/common/common.go`
- **Reuse**: Each asset in the collage gets a `ViewImageData` (ImageData, Asset, etc.). Layout uses `viewData.Assets` with length 1–8.

### MemoriesViewState (client-side)

- **Location**: Frontend only (TypeScript)
- **Fields**: `lightboxOpen: boolean`, `remainingSeconds: number`, `paused: boolean`
- **Purpose**: Track lightbox visibility and timer pause for resume-on-close behavior.

## State Transitions

### Slideshow → Memories

1. `gatherAssetBuckets` includes memories when `Memories == true` and `MemoriesAssetsCount > 0`.
2. `retrieveImage` picks `SourceMemories` → calls `MemoriesCollage` (new method) instead of `RandomMemoryAsset`.
3. Provider returns up to 8 assets; `generateViewData` builds `ViewData` with `layout: "collage"`, `Duration: 45`, `Assets` with 1–8 items.
4. Template renders collage; frontend starts 45s timer.

### Memories → Next Item

1. Timer expires or user dismisses (on-demand) → `hx-post="/asset/new"` with current params.
2. Next asset (could be memories again or another source) is fetched.

### Collage → Lightbox → Collage

1. User taps photo → lightbox opens; `pausePolling()` called; timer state stored.
2. User closes lightbox → `resumePolling()`; remaining time used for next poll.
3. If remaining time ≤ 0, next poll immediately fetches next asset.

## Relationships

```
MemoriesCollage
  └── Assets[] → DisplayAsset (1–8)

ViewData (memories)
  └── Assets[] → ViewImageData (1–8)
  └── Duration = 45
  └── Layout = "collage"
  └── MemoryCaption = "Memories from X years ago"
```

## Data Flow

1. **Immich**: `Memories(for=date)` → `[]Memory`; filter by year; pick one; take up to 8 assets.
2. **PhotoPrism**: `Photos(search=date)` or equivalent → `[]Photo`; map to `DisplayAsset`; take up to 8.
3. **Provider**: New method `MemoriesCollage(requestID, deviceID) (MemoriesCollage, error)` → used by `retrieveImage` when bucket is `SourceMemories`.
4. **Routes**: `generateViewData` detects collage response → sets `Duration=45`, `Layout=collage`, `MemoryCaption`.
