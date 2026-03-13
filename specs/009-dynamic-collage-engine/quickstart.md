# Quickstart: Dynamic Collage Engine

**Feature**: 009-dynamic-collage-engine  
**Date**: 2026-03-13

## Prerequisites

- Go 1.26
- Existing memories feature (specs 007, 008) implemented
- Immich or PhotoPrism as photo source (assets with `ExifImageWidth`, `ExifImageHeight`)

## Implementation Order

1. **Create `internal/collage` package**
   - `LayoutResult`, `LayoutCell` types
   - `ComputeLayout(assets, width, height) (LayoutResult, error)`
   - Implement justified-rows, treemap, hero-plus-cluster algorithms
   - Add `layout_test.go` with unit tests

2. **Integrate into `ProcessMemoriesCollage`**
   - Call `collage.ComputeLayout` with 2-second timeout (context.WithTimeout)
   - On success: attach `LayoutResult` to ViewData
   - On failure/timeout: leave LayoutResult nil (trigger fallback)

3. **Update `collage.templ`**
   - If `ViewData.CollageLayout != nil`: render cells with positions from layout
   - Else: render with `collage-{n}` class (fallback)

4. **Update `ViewData` in `internal/common/common.go`**
   - Add `CollageLayout *collage.LayoutResult` (or equivalent)

5. **Update `collage.css`**
   - Dynamic layout: ensure `object-fit: contain` for aspect-ratio preservation
   - Fallback: change `object-fit: cover` to `object-fit: contain` for `.collage-cell img`

## Key Files

| File | Change |
|------|--------|
| `internal/collage/layout.go` | New; entry point, types |
| `internal/collage/justified.go` | New; justified-rows algorithm |
| `internal/collage/treemap.go` | New; treemap packing |
| `internal/collage/hero.go` | New; hero-plus-cluster |
| `internal/routes/routes_asset_helpers.go` | Call engine, attach layout |
| `internal/common/common.go` | Add CollageLayout to ViewData |
| `internal/templates/components/memories/collage.templ` | Conditional render |
| `frontend/src/css/collage.css` | object-fit: contain |

## Testing

```bash
go test ./internal/collage/...
go test ./internal/routes/...
```

Manual: Trigger memories (slideshow or on-demand); verify dynamic layouts, aspect-ratio preservation, lightbox-on-tap, fallback when engine fails.
