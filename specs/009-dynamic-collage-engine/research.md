# Research: Dynamic Collage Engine

**Feature**: 009-dynamic-collage-engine  
**Date**: 2026-03-13

## 1. Collage Generation Approach — Output Format Constraint

**Constraint (from spec clarification)**: Output MUST allow per-photo tap for lightbox. Single pre-rendered image (PNG) is out of scope.

**Implication**: Any approach that produces only a composite image (PNG, PDF, JPEG) is excluded. The engine must output HTML (or equivalent) with individual photo elements.

---

## 2. Approach Evaluation

### 2.1 photocollage (Python)

**Decision**: Rejected.

**Rationale**: photocollage produces poster images (PNG/PDF). It does not output HTML or per-photo coordinates in a format suitable for tappable elements. Integrating would require Python runtime, subprocess invocation, and a custom parser to extract layout positions — and the library does not expose layout metadata. Per-photo tap cannot be supported.

**Alternatives considered**: Use photocollage to compute positions and output a custom format — not supported by the library.

---

### 2.2 Pillow (PIL) with Custom Packing

**Decision**: Rejected for this project.

**Rationale**: Pillow is Python-based. The kiosk is Go-only; adding Python would require subprocess, runtime dependency, and cross-language coordination. Layout algorithms (treemap, justified rows) can be implemented in Go without image processing — we only need position math, not pixel manipulation. Pillow would add complexity without benefit.

**Alternatives considered**: Use Pillow via subprocess to compute layout and output JSON positions — possible but adds Python dependency; pure Go is simpler.

---

### 2.3 ImageMagick montage

**Decision**: Rejected.

**Rationale**: `montage` produces a composite image. Per-photo tap requires individual elements; ImageMagick cannot output HTML. Would need a separate step to compute per-photo coordinates, and montage does not expose layout metadata in a usable form.

---

### 2.4 CollageCreator (Go library)

**Decision**: Rejected for this project.

**Rationale**: [CollageCreator](https://github.com/schwerdf/CollageCreator) is a Go library for automatic collage generation with extensible layout and rendering. It provides:
- **Layout algorithms**: Random placement (images equidistant from neighbors) and Tile-in-order (justified rows, scaled to fit)
- **Output**: PNG/JPEG/TIFF raster, SVG (with file links), or ImageMagick shell script
- **ImageLayout API**: Exposes `PositionOf(img)`, `DimensionsOf(img)` — per-image positions are available

**Why rejected**:
1. **Input model mismatch**: Built-in `InputImageReader` expects local file paths (`InFiles`). Our assets come from Immich/PhotoPrism APIs (URLs, `DisplayAsset` with `ExifImageWidth`/`ExifImageHeight`). Adapting would require a custom `InputImageReader` and `ImageLayout` implementation to bridge our asset model — non-trivial.
2. **Output**: Raster and SVG reference local files. We need HTML with `/image/:id` URLs for per-photo tap. We could use `PositionCalculator` only and extract positions from `ImageLayout` to render our own HTML — but that still requires solving the input adapter.
3. **License**: LGPL-3.0 — more restrictive than typical permissive licenses.
4. **Maturity**: Very early (0 stars, 2 commits); API may change.
5. **Architecture**: Designed for CLI with `Parameters`, `ProgressMonitor`, file-based I/O. Fitting into request-time web flow adds glue code.

**Alternatives considered**: Implement custom `InputImageReader` + `ImageLayout` that wraps `DisplayAsset` and dimensions — possible but adds dependency and adapter complexity; pure in-house implementation is simpler and has no external coupling.

---

### 2.5 Pure Go Layout Algorithm

**Decision**: Selected.

**Rationale**:
- Outputs HTML with positioned photo elements — satisfies per-photo tap.
- No external runtime (no Python, no ImageMagick).
- Uses existing `DisplayAsset` with `ExifImageWidth` and `ExifImageHeight` from Immich/PhotoPrism.
- Layout computation is pure math (position/size calculation); completes in milliseconds, well under 2-second target.
- Integrates with existing Templ pipeline: engine returns layout data (positions, sizes); template renders HTML.
- Aligns with constitution: Go standards, minimal dependencies, existing project structure.

**Alternatives considered**:
- Python/Pillow subprocess — adds runtime dependency; rejected.
- Client-side JS layout (e.g. Masonry Grid) — would require new frontend dependency, layout shift/FOUC; backend control is preferred for consistency.

---

## 3. Layout Algorithm Strategy

**Decision**: Implement multiple layout algorithms in Go; select one at random per render to achieve variety.

**Algorithms to implement**:
1. **Justified rows** — Rows of similar height; photos scaled to fit row width while preserving aspect ratio. Gaps between rows acceptable.
2. **Treemap-style packing** — Recursive subdivision of space; photos get area proportional to aspect-ratio-adjusted size. Produces organic, varied layouts.
3. **Hero-plus-cluster** — One large photo (hero) plus remaining photos in a cluster (grid or justified). Randomly choose hero.

**Rationale**: These cover the spec's "justified rows, treemap-style packing, hero-plus-cluster, randomized organic" styles. Random selection per render yields variety without fixed per-count variants.

**Alternatives considered**:
- Single algorithm — less variety; rejected.
- More than 3 algorithms — increases implementation scope; 3 is sufficient for initial delivery.

---

## 4. Integration Point

**Decision**: Insert the collage engine between `ProcessMemoriesCollage` and the collage template. `ProcessMemoriesCollage` continues to fetch assets and build `ViewData`; a new component computes layout from assets and attaches layout data to `ViewData`; the template renders using that data instead of fixed CSS grid classes.

**Flow**:
1. `provider.MemoriesCollage()` → assets (3–16) with `ExifImageWidth`, `ExifImageHeight`
2. Collage engine: `ComputeLayout(assets, containerWidth, containerHeight) → LayoutResult`
3. `LayoutResult` contains per-photo `{x, y, width, height}` (or equivalent)
4. Template renders `collage-cell` divs with inline styles or CSS variables from `LayoutResult`
5. Fallback: if engine fails or exceeds 2s, use simple grid (existing `collage-{n}` fallback)

**Rationale**: Minimal change to existing flow; engine is a drop-in replacement for `PickLayoutVariant` + fixed grid.

---

## 5. Fallback Implementation

**Decision**: Reuse existing fixed grid CSS (`collage-{n}`) as fallback. When engine fails or times out, omit layout data; template detects missing layout and applies `collage-{n}` class (simple uniform grid).

**Rationale**: Spec requires "simple grid" fallback. The current `collage-{n}` classes provide a basic grid; we keep them for fallback. Fallback must also preserve aspect ratio — change `object-fit: cover` to `object-fit: contain` in fallback CSS so photos are not cropped.

---

## 6. Aspect Ratio Preservation in CSS

**Decision**: Dynamic layout uses `object-fit: cover` and `object-position: center` for collage cells. The layout engine constrains cell aspect ratios to an allowed range (up to 10% cropping from each end of the longer dimension). Fallback grid uses `object-fit: contain`.

**Rationale**: Allowing up to 10% crop improves packing and reduces gaps while keeping distortion minimal. The layout algorithms choose cell dimensions within the allowed aspect range; `cover` with centered crop displays the result. Fallback keeps `contain` since it uses fixed grid cells that may not match image ratios.

---

## Summary

| Decision | Choice |
|----------|--------|
| Approach | Pure Go layout algorithm |
| Output | HTML with positioned elements (via Templ) |
| Layout algorithms | Justified rows, treemap packing, hero-plus-cluster |
| Integration | Between ProcessMemoriesCollage and template |
| Fallback | Existing collage-{n} grid with object-fit: contain (no crop) |
| Performance | Layout computation in Go; target <2s (layout itself is ms) |

### Alternatives Evaluated (Including CollageCreator)

| Option | Outcome |
|--------|---------|
| photocollage (Python) | Rejected — PNG output, no per-photo tap |
| Pillow (PIL) | Rejected — Python dependency |
| ImageMagick montage | Rejected — composite image only |
| [CollageCreator](https://github.com/schwerdf/CollageCreator) (Go) | Rejected — input model (file paths) vs our URLs; LGPL-3.0; low maturity; would require custom InputImageReader/ImageLayout adapters |
| Pure Go (in-house) | **Selected** — no deps, full control, direct fit with DisplayAsset |
