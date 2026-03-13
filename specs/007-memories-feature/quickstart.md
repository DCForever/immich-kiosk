# Quickstart: Memories (On This Day)

**Feature**: 007-memories-feature  
**Date**: 2026-03-11

## Prerequisites

- Go 1.26
- Immich instance with memories enabled, or PhotoPrism (with date-search support)
- Existing kiosk config with `memories: true` (or equivalent)

## Implementation Order

1. **Provider layer: MemoriesCollage**
   - Immich: Extend `immich_memories.go` — add `MemoriesCollage` that filters by year, picks one, returns up to 8 assets.
   - PhotoPrism: Add `PhotosOnThisDay` or equivalent in `client.go`; implement `MemoriesCollage` in provider (or return error if API unsupported).

2. **Source interface**
   - Add `MemoriesCollage(requestID, deviceID) (MemoriesCollage, error)` to `ProviderOps` (or equivalent) if not using `RandomMemoryAsset` for collage.
   - Alternative: Extend `RandomMemoryAsset` to return multiple assets when in "collage mode" — requires interface change. Prefer new method `MemoriesCollage` for clarity.

3. **Routes**
   - `routes_asset_helpers.go`: In `retrieveImage`, when bucket is `SourceMemories`, call `MemoriesCollage`; build `ViewData` with `layout: "collage"`, `Duration: 45`, `Assets` from collage.
   - `routes_memories.go`: New handler for `GET/POST /memories`; same logic as asset path but for on-demand; render with back button.

4. **Templates**
   - `templates/components/memories/collage.templ`: Grid layout for 1–8 photos; caption; tap handlers.
   - `templates/components/memories/lightbox.templ`: Optional; lightbox may be client-only.
   - `templates/components/image/layout.templ`: Add branch for `layout == "collage"` → render collage component.

5. **Frontend**
   - `collage.css`: Grid layouts for `collage-1` through `collage-8`.
   - `memories-lightbox.ts`: Open/close lightbox; pause/resume polling; integrate with `polling.ts`.
   - `menu.templ`: Add memories button.

6. **i18n**
   - Add keys to `locales/*.toml` for `memories_caption`, `memories_no_photos`, `memories_error_load`, `memories_back`.

7. **Config**
   - Ensure `memories` and `past_memory_days` (Immich) are configurable; add `memory_duration` (default 45) if desired.

## Testing

- Unit: `MemoriesCollage` logic (year selection, asset count limit); mock Immich/PhotoPrism responses.
- Integration: Start kiosk with `memories=true`; verify collage appears in rotation; verify on-demand button.
- Manual: Tap photo → lightbox; close → timer resumes; 45s → next item.

## Verification Checklist

- [ ] Memories collage appears in slideshow when `memories=true` and photos exist.
- [ ] 45-second duration for memories view.
- [ ] "Memories from X years ago" caption.
- [ ] On-demand button opens memories view.
- [ ] Back/dismiss returns to slideshow.
- [ ] No photos: slideshow skips; on-demand shows "No memories for this day" (45s).
- [ ] Lightbox: tap photo → full view; close → back to collage; timer paused.
- [ ] Load error: loading spinner during fetch; "Couldn't load memories" (5s) on failure.
- [ ] i18n for all new strings (memories_caption, memories_no_photos, memories_error_load, memories_back, memories_button).
- [ ] PhotoPrism: shows "Couldn't load memories" gracefully (memories not supported).
