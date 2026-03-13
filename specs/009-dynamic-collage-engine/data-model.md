# Data Model: Dynamic Collage Engine

**Feature**: 009-dynamic-collage-engine  
**Date**: 2026-03-13

## Entities

### LayoutResult

Result of the collage layout engine. Passed to the template for rendering.

| Field | Type | Description |
|-------|------|-------------|
| Cells | []LayoutCell | Per-photo layout data; one cell per asset in order |
| Algorithm | string | Algorithm used: "justified", "treemap", "hero" |
| ContainerWidth | float64 | Target container width (px or logical units) |
| ContainerHeight | float64 | Target container height |

### LayoutCell

Layout data for a single photo in the collage.

| Field | Type | Description |
|-------|------|-------------|
| X | float64 | Left position (0–1 or px) |
| Y | float64 | Top position (0–1 or px) |
| Width | float64 | Cell width |
| Height | float64 | Cell height |
| AssetIndex | int | Index into original assets slice (for data-asset-id binding) |

**Note**: Positions and sizes can be in normalized units (0–1) or pixels; template applies appropriate CSS (e.g. `left: X%`, `width: W%` or `position: absolute; left: Xpx`).

### Input to Layout Engine

The engine receives:

| Input | Source | Description |
|-------|--------|-------------|
| Assets | []DisplayAsset | From provider.MemoriesCollage(); 3–16 items |
| Width | float64 | Container width (from viewport or config) |
| Height | float64 | Container height |
| Gap | float64 | Spacing between cells (optional) |

**Aspect ratios**: Derived from `asset.ExifInfo.ExifImageWidth` and `ExifInfo.ExifImageHeight`. When either is 0, use default (e.g. 1.0 for square).

## Relationships

- **LayoutResult** → **LayoutCell[]**: One-to-many; order matches assets.
- **LayoutCell** → **DisplayAsset**: Linked by index; template uses `Assets[AssetIndex]` for `data-asset-id` and image URL.

## State Transitions

- **ComputeLayout** (sync): Assets + dimensions → LayoutResult or error.
- **Fallback**: On error or timeout (>2s), LayoutResult is nil; template uses `collage-{n}` class.

## Validation Rules

- `len(Assets)` must be 3–16 (enforced by ProcessMemoriesCollage before engine call).
- Each `LayoutCell` must have Width > 0, Height > 0.
- Cells must fit within ContainerWidth × ContainerHeight (with gap allowance).
- AssetIndex must be in range [0, len(Assets)).
