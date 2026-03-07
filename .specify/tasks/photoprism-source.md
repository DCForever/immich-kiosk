# Tasks: PhotoPrism as media source

## Reference

- Spec: `.specify/specs/photoprism-source.md`
- Plan: `.specify/plans/photoprism-source.md`
- Constitution: `.specify/memory/constitution.md`

## Task list

Tasks must respect:

- **Code quality**: golangci-lint, Biome, project structure; no new third-party libs; prefer simplicity.
- **Testing**: Add/update tests for new behavior; mock external APIs; keep CI passing.
- **UX**: i18n for new strings, kiosk consistency, config-driven.
- **Performance**: Reuse cache/prefetch; bounded memory; thumbnails where appropriate.

| # | Task | Done |
|---|------|------|
| **Phase 1: Shared type and interface** |
| 1 | Define shared “display asset” type (thumb URL, original URL, type image/video, metadata, bucket/source id). Place in `internal/kiosk` or new `internal/source`; document which fields templates and routes need. | |
| 2 | Audit routes and `routes_asset_helpers.go` for all uses of `*immich.Asset` (methods and fields). Document the full list of operations the provider must support. | |
| 3 | Define the provider interface with only those operations (e.g. RandomAsset, AssetFromAlbum, PersonAssetCount, AlbumImageCount, ImagePreview, AssetInfo, tag/favourite/archive actions, etc.). No implementation yet. | |
| **Phase 2: Immich adapter and factory** |
| 4 | Add Immich adapter in `internal/immich`: type that wraps `*immich.Asset` and implements the new provider interface by delegating to existing methods and mapping to the shared display asset type. | |
| 5 | Add provider factory: given `config.Config`, return the provider. For now always return the Immich adapter (source not yet in config). Place factory where config is available (e.g. routes or a small `internal/source` helper). | |
| 6 | Refactor `routes_asset_helpers.go`: replace `*immich.Asset` with the interface; get provider from factory. Update `gatherAssetBuckets`, `retrieveImage`, `processAsset`, `fetchImagePreview` and any helpers to use interface + shared type. | |
| 7 | Refactor `routes_asset.go`, `routes_previous_asset.go`, `routes_url_builder.go`, `routes_webhooks.go`, and about view: replace `immich.New` with factory; use interface and shared type for asset operations. Update template data to use shared type where applicable. | |
| **Phase 3: Config for source and PhotoPrism** |
| 8 | Add config fields: `source` (immich | photoprism, default `immich`), `photoprism_url`, `photoprism_token`. Map from YAML/env; ensure existing configs without `source` behave as today (Immich). | |
| 9 | Config validation: when `source == "photoprism"` require `photoprism_url` and `photoprism_token`; when `source == "immich"` or default require `immich_url` and `immich_api_key`. Integrate into existing validation flow. | |
| 10 | Add unit tests for config: valid/invalid combinations (source+url+token for each backend). | |
| 11 | Update `config.example.yaml` and config schema (e.g. `config.schema.json`) with `source`, `photoprism_url`, `photoprism_token` and brief comments. | |
| **Phase 4: PhotoPrism client and provider** |
| 12 | Create `internal/photoprism`: minimal HTTP client (stdlib only: `net/http`, `encoding/json`). Structs for PhotoPrism API responses needed (e.g. photo list from GET /api/v1/photos). Auth via token header; document in code. | |
| 13 | Implement provider interface for PhotoPrism: random/recent photos (order=random or newest, count, offset), map response to shared display asset type. Implement ImagePreview (and thumb/original URLs using PhotoPrism’s thumbnail API + token). | |
| 14 | Add album support: album-scoped photos (e.g. `s` param or album API), album list/count for weighting. Map album UIDs to config album identifiers. | |
| 15 | Add labels (tags) and date-range support for PhotoPrism; add favorites if API supports it. Stub or no-op people/memories; document “not supported for PhotoPrism.” | |
| 16 | Reuse existing cache for PhotoPrism API responses (bounded keys/TTLs). Ensure prefetch path uses the interface so PhotoPrism assets can be prefetched. | |
| 17 | Implement any provider methods that mutate state (e.g. AddTag, FavouriteStatus) as no-ops or best-effort for PhotoPrism if API allows; otherwise document and no-op. | |
| **Phase 5: Wire and i18n** |
| 18 | Factory: when `config.Source == "photoprism"` build PhotoPrism client and return PhotoPrism provider; otherwise return Immich adapter. Ensure routes use factory everywhere. | |
| 19 | URL builder: support PhotoPrism (e.g. album UIDs, source param) so generated URLs work when source is PhotoPrism. Keep UI labels consistent. | |
| 20 | Add i18n keys for new user-facing strings (e.g. source name “PhotoPrism”, errors “PhotoPrism connection failed”). Add to locale files used by the project. | |
| **Phase 6: Tests and docs** |
| 21 | Unit tests for PhotoPrism client and response→display-asset mapping using `httptest.Server`; no real network. | |
| 22 | Confirm existing Immich tests pass with adapter + factory; add or adjust tests as needed. Config validation tests for source/PhotoPrism (covered in task 10). | |
| 23 | Document which features work per source (e.g. Memories: Immich only; People: Immich yes, PhotoPrism if supported). Update config example or docs accordingly. | |

## Definition of done (per task)

- [ ] Implements only what the task describes.
- [ ] Aligns with constitution (quality, tests, UX, performance).
- [ ] No new third-party libraries; prefer simplicity.
- [ ] Lint and tests pass for touched code.
