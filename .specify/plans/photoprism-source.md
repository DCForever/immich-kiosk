# Plan: PhotoPrism as media source

## Goal

Deliver PhotoPrism as an alternative media source for the kiosk so users can choose Immich or PhotoPrism per configuration. Same slideshow behavior and config-driven UX for both. No new third-party libraries; keep the design as simple as possible.

## Constraints (from request)

- **No new libraries if possible**: Use only the Go standard library and existing project dependencies (e.g. `net/http`, `encoding/json`, existing `cache`, `config`, `kiosk`). Do not add a third-party PhotoPrism client; implement a minimal HTTP client against the public PhotoPrism REST API.
- **Prefer simplicity to complexity**: Prefer a narrow interface and thin adapters over a large abstraction. Avoid over-generalising for future backends; two implementations (Immich, PhotoPrism) are enough to shape the interface.

## Constitution alignment

All work must comply with `.specify/memory/constitution.md`:

- **Code quality**: Lint-clean Go/templ/TS; clear structure; no new lint issues.
- **Testing**: Tests for new or changed behavior; mock external APIs in tests; CI green.
- **UX**: Kiosk-first, config-driven, i18n for new strings, consistent design.
- **Performance**: Reuse existing cache and prefetch; bounded memory; use thumbnails where appropriate.

## Phases / steps

### Phase 1: Minimal shared asset type and provider interface

- Define a **single, source-agnostic “display asset” type** (e.g. in `internal/kiosk` or a small `internal/source` package). It holds only what the routes and templates need: thumb URL, original URL, asset type (image/video), basic metadata (e.g. date, title, album names, person names), and any fields needed for weighting/history (e.g. bucket/source id).
- Define a **narrow provider interface** (e.g. “get next random asset,” “get asset count by album,” “get asset count by person,” “get asset by album/id,” “get asset by person/id,” etc.) that matches how routes currently use `immich.Asset`. Keep the interface minimal: only the operations actually used by `routes_asset_helpers.go` and related code.
- Do **not** move or rewrite Immich internals yet; only introduce the type and the interface so that the next phase can implement them.

**Simplicity**: One shared struct, one small interface. No new packages unless necessary (e.g. if `internal/kiosk` would grow too large, use `internal/source` for the type and interface only).

### Phase 2: Immich implements the interface (adapter)

- Add an **adapter** in `internal/immich` that implements the new provider interface by delegating to existing `immich.Asset` and related functions. The adapter wraps `*immich.Asset` and translates between Immich types and the shared “display asset” type.
- Introduce a **factory** (e.g. in routes or a small place that has access to config): given `config.Config`, return the provider implementation (Immich adapter or, later, PhotoPrism). For now, always return the Immich adapter when `source` is missing or `immich`.
- **Refactor routes and asset helpers** so they depend only on the interface (and the shared type), not on `*immich.Asset`. Replace direct `immich.New(...)` usage with “get provider from factory, then call interface methods.” Keep all existing Immich behavior; this is a dependency inversion only.

**Simplicity**: Adapter pattern only; no rewrite of Immich API logic. Existing tests for Immich flows should still pass (possibly with small updates to construct the adapter).

### Phase 3: Config for source selection and PhotoPrism connection

- Add to **config**: `source` (e.g. `immich` | `photoprism`, default `immich` when omitted for backward compatibility), `photoprism_url`, `photoprism_token` (or a single token field; prefer one auth method for simplicity). Keep Immich fields as-is; when `source == "photoprism"`, Immich URL/key can be optional.
- **Validation**: If `source == "photoprism"`, require `photoprism_url` and `photoprism_token`; if `source == "immich"` (or default), require `immich_url` and `immich_api_key` as today. Add unit tests for valid/invalid combinations.
- **Schema and example**: Update `config.example.yaml` and config schema with the new keys and short comments.

**Simplicity**: No new dependencies; validation logic in existing config package.

### Phase 4: PhotoPrism client and provider (stdlib only)

- Add package **`internal/photoprism`**. Implement a **minimal client** using only `net/http` and `encoding/json`: no third-party PhotoPrism SDK. Define local structs for the PhotoPrism API responses you need (e.g. photo list from `GET /api/v1/photos`, album list if needed). Use `http.Client` with timeout; pass token via header (e.g. `Authorization: Bearer <token>` or PhotoPrism’s documented header).
- Implement **auth**: Prefer the simplest option that works (e.g. token in header from app password or session). Document in code and config example; avoid OAuth or multi-step auth in v1 if not strictly required.
- **Map PhotoPrism responses to the shared “display asset” type**: Implement the same provider interface as the Immich adapter. Start with the minimal set of operations: e.g. random/recent photos (via `order=random` or `order=newest` and pagination), and album-scoped photos (using PhotoPrism’s `s` or album API). Add **labels (tags)** and **date range** next; then **favorites** if the API supports it. Defer or stub **people/subjects** if the API is unclear; document “not supported for PhotoPrism” for memories and any other Immich-only features.
- **Caching**: Reuse the existing cache package (or same pattern as Immich) for PhotoPrism API responses; same bounded keys and TTLs. **Prefetch**: Reuse the same prefetch flow as for Immich, using the interface, so the next asset can be loaded in the background for PhotoPrism too.
- **Thumbnails/originals**: Use PhotoPrism’s thumbnail and download URLs with the required token (e.g. from response headers or auth); no new library.

**Simplicity**: One new package, minimal structs, no new deps. Implement only what’s needed for a usable slideshow (photos, albums, labels, dates, favorites if easy). Document gaps (e.g. people, memories) instead of over-engineering.

### Phase 5: Wire PhotoPrism into factory and routes

- **Factory**: When `config.Source == "photoprism"`, construct the PhotoPrism client and return the PhotoPrism implementation of the provider interface; otherwise return the Immich adapter. Routes already use the interface, so no route logic changes beyond using the factory.
- **URL builder**: If the URL builder currently exposes Immich-specific IDs (e.g. album IDs), extend it to support PhotoPrism (e.g. album UIDs) when `source == "photoprism"` so that generated URLs work for both. Keep the UI consistent (same labels, different param names/values per source if needed).
- **i18n**: Add keys for any new user-facing strings (e.g. source name “PhotoPrism”, errors like “PhotoPrism connection failed”). Use existing i18n system; no new dependencies.

### Phase 6: Tests and documentation

- **Unit tests**: PhotoPrism client and mapping from API response to shared asset type, using `httptest.Server` to fake the API; no real network in CI. Config validation tests for `source` and PhotoPrism URL/token. Ensure existing Immich tests still pass (adapter + factory).
- **Docs**: Update config example and schema with “which features work per source” (e.g. “Memories: Immich only”; “People: Immich yes; PhotoPrism if API supports subjects”). Keep the list short and accurate.

## Success criteria

- [ ] Spec requirements satisfied: single source per config, provider abstraction, PhotoPrism as second backend, config and validation, caching and prefetch, i18n.
- [ ] No new third-party libraries; PhotoPrism client uses only stdlib + existing deps.
- [ ] Simplicity: narrow interface, thin adapters, minimal new types; no over-abstraction.
- [ ] Constitution: lint-clean, tests for new/changed paths, UX and performance as above.
- [ ] Existing Immich configs keep working (default `source` or `source: immich`).

## Reference

- Spec: `.specify/specs/photoprism-source.md`
- Constitution: `.specify/memory/constitution.md`
