# Spec: PhotoPrism as media source

## Summary

Add PhotoPrism as an alternative backend for the kiosk. Users can choose either Immich or PhotoPrism as the source of photos, albums, and all other filterable content. The kiosk will support a single source per configuration; behavior (slideshow, transitions, config-driven options) remains consistent regardless of which backend is selected.

## Context

- Immich Kiosk today is tightly coupled to the Immich API: config uses `immich_url` and `immich_api_key`, and all asset fetching (random, by person, album, tag, date, rating, memories) goes through the `internal/immich` package.
- Users who host PhotoPrism instead of (or in addition to) Immich have no way to use the kiosk with their library. PhotoPrism exposes a REST API (`/api/v1/photos`, albums, search) that can support the same conceptual model: photos/videos, albums, labels (tags), date filters, favorites, and optionally subjects (people).
- Introducing a second source improves reach and keeps the product focused on “kiosk for self‑hosted photo libraries” rather than Immich-only.

## Requirements

Requirements are aligned with the project constitution (code quality, testing, UX consistency, performance).

### 1. Source abstraction

- **Provider interface**: Introduce a media-source abstraction used by routes and asset helpers. Callers (e.g. `routes_asset.go`, `routes_asset_helpers.go`) must not depend on `immich` directly for “get next asset” or “get asset by source type”; they depend on an interface (or factory) that can be implemented by Immich and PhotoPrism.
- **Concrete types**: The kiosk’s internal representation of “an asset to display” (URLs for thumb/original, type image/video, EXIF-like metadata, album/person/tag info for UI) must be defined in a source-agnostic way (e.g. in `internal/kiosk` or a dedicated `internal/source` package). Immich-specific and PhotoPrism-specific code each translate their API responses into this representation.
- **Config**: Configuration must support selecting the source (e.g. `source: immich | photoprism`) and source-specific connection settings:
  - Immich: existing `immich_url`, `immich_api_key` (and any existing multi-user keys).
  - PhotoPrism: e.g. `photoprism_url`, `photoprism_token` (or session-based auth as per PhotoPrism API). Document which auth method(s) are supported (e.g. session token / app password).
- **Validation**: Config validation must require the correct URL and credentials for the selected source (e.g. if `source: photoprism` then `photoprism_url` and token are required; Immich fields can be optional when source is PhotoPrism).

### 2. PhotoPrism backend implementation

- **Photos and thumbnails**: Implement fetching of photos (and videos if supported by PhotoPrism and by kiosk config) using PhotoPrism’s API (e.g. `/api/v1/photos` with `merged`, `order`, `count`, `offset`). Thumbnail and original URLs must follow PhotoPrism’s thumbnail/original API and any required security tokens (e.g. from response headers or auth).
- **Albums**: Support album-based filtering using PhotoPrism album UIDs. Map config album identifiers to PhotoPrism’s album scope (e.g. `s` parameter or album endpoint). Support “favorites” if PhotoPrism exposes it (e.g. as a virtual album or filter).
- **Labels (tags)**: Map kiosk “tags” config to PhotoPrism labels/categories where applicable. Support inclusion and exclusion by label.
- **Date range**: Support date-range filtering using PhotoPrism’s search (e.g. `year`, `month` or query parameters) so that existing `dates` config behavior has an equivalent.
- **Rating / favorites**: If PhotoPrism supports quality/favorite flags, map kiosk rating/favorite options to the corresponding filters; otherwise document as unsupported for PhotoPrism and hide or no-op those options when source is PhotoPrism.
- **People / subjects**: If PhotoPrism exposes “subjects” or “people,” support person-based filtering and show person name in metadata where applicable; otherwise document as unsupported and ignore `people` / `excluded_people` when source is PhotoPrism.
- **Archive / private**: Map “show archived” and any private/hidden concepts to PhotoPrism’s visibility if the API supports it.

### 3. Feature parity and UX

- **Kiosk-first**: Slideshow behavior (duration, transitions, preload, history) must work the same for both sources. No backend-specific UI modes; only configuration may differ (e.g. which options are available per source).
- **Configuration-driven**: All source-specific options (albums, tags/labels, people/subjects, dates, rating, archive) must be driven by config and, where relevant, URL overrides. Options not supported by the selected source must be ignored or clearly documented (e.g. “Memories are not supported when using PhotoPrism”).
- **i18n**: Any new user-facing strings (e.g. “PhotoPrism” in about/settings, or errors like “PhotoPrism connection failed”) must use the project’s i18n system.
- **URL builder**: If the URL builder exposes source-specific parameters (e.g. album IDs), it must support both Immich and PhotoPrism (e.g. album UIDs for PhotoPrism vs IDs for Immich), with labels and behavior consistent for the selected source.

### 4. Code quality and structure

- **Packages**: Keep a clear boundary: Immich-specific code remains under `internal/immich` (or moves behind the new abstraction); add e.g. `internal/photoprism` for PhotoPrism client and mapping. Shared types and the provider interface live in a neutral package (e.g. `internal/source` or `internal/kiosk`) so that routes and asset helpers do not import `immich` or `photoprism` for “get asset” logic.
- **Linting**: All new and modified code must pass `golangci-lint` with the project’s `.golangci.yml`. No new lint issues.
- **Dependencies**: Prefer standard library and existing deps. If a third-party PhotoPrism client is used, it must be justified and maintained; otherwise implement a minimal client against the public PhotoPrism REST API.

### 5. Testing

- **Unit tests**: New code paths (PhotoPrism client, mapping from PhotoPrism responses to internal asset type, config validation for `source` and PhotoPrism URL/token) must have unit tests. Use mocks or fakes for the PhotoPrism HTTP API in tests; no real network calls in CI.
- **Existing behavior**: Immich-based flows must remain covered by existing (or updated) tests. No regressions in covered paths.
- **Config validation**: Tests for valid/invalid combinations (e.g. `source: photoprism` with missing `photoprism_url`, or `source: immich` with missing `immich_api_key`).

### 6. Performance

- **Caching**: Reuse the existing cache layer (or an equivalent) for PhotoPrism API responses and, where applicable, derived data (e.g. album list, search results). Cache keys and TTLs must be bounded and not cause unbounded growth.
- **Prefetch**: Prefetch behavior (e.g. “next asset” fetched in the background) must work for PhotoPrism-sourced assets so that transitions remain smooth.
- **Asset loading**: Use thumbnail/preview URLs where appropriate to avoid loading originals until needed; respect PhotoPrism’s token/security model for thumbnails and originals.

### 7. Documentation and schema

- **Config example**: Update `config.example.yaml` (and any schema) with `source`, `photoprism_url`, `photoprism_token` (or equivalent), and comments that explain Immich vs PhotoPrism options.
- **Docs**: Document which kiosk features are supported for each source (e.g. “Memories: Immich only”; “People: Immich yes, PhotoPrism if API supports subjects”).

## Out of scope

- **Multiple sources in one config**: One configuration uses one source (Immich or PhotoPrism). Mixing both in a single kiosk instance (e.g. “albums from Immich and PhotoPrism”) is not required for this spec.
- **PhotoPrism-only features**: No requirement to support PhotoPrism-specific concepts that have no Immich counterpart (e.g. stacks, places) beyond what is needed to implement the above mappings.
- **Migration tooling**: No tool to migrate or convert config from Immich to PhotoPrism; users can edit config manually.
- **Live Photos / motion**: If PhotoPrism’s model differs significantly from Immich’s (e.g. no separate “live photo video” ID), implement best-effort behavior (e.g. show primary image only) and document; full parity is out of scope if it requires large product changes.

## Open questions

- **PhotoPrism auth**: Confirm whether to support only session token (e.g. from `POST /api/v1/session`) or also app passwords / OAuth; document and implement one or both.
- **Subjects/people in PhotoPrism**: Confirm API support and field names (e.g. “subjects”) so that person-based filtering and “show person name” can be specified precisely.
- **Memories**: Decide whether to hide “memories” in UI when source is PhotoPrism or show a disabled state with tooltip; same for any other Immich-only feature.
- **Backward compatibility**: Whether `source` defaults to `immich` when omitted and Immich URL/key are set, so that existing configs keep working without change.

## Reference

- **PhotoPrism API (Swagger):** https://docs.photoprism.dev/
