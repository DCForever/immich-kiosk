# Plan: PhotoPrism people support (parity with Immich)

## Goal

Add people/subjects support for the PhotoPrism source so that it matches Immich's people behaviour: filter slideshow by person, show person names in metadata, support "all people" and excluded people in config and URL builder, and use people in asset weighting and retrieval.

## Context

- **Immich (current)**:
  - **List people**: `GET /api/people` (paginated), filter to named only for "all named people."
  - **Person asset count**: `GET /api/people/:id/statistics` (or sum over all named for `person: all`).
  - **Random asset of person**: `POST /api/search/random` with `personIds`, `type`, `withPeople`, etc.
  - **Random person from all**: fetch all named people, apply exclusions, pick random ID.
  - **Faces on asset**: `GET /api/faces?id=<assetId>` → convert to `People` on asset; `DisplayAsset.People` used in more-info and ShowPersonName/ShowPersonAge.
  - **Config**: `people`, `excluded_people`, `require_all_people`; URL builder uses `AllNamedPeople()` for dropdown.
  - **Flow**: `gatherAssetBuckets` uses `PersonAssetCount` per config person; `retrieveImage` for `SourcePerson` uses `RandomPersonFromAllPeople` when ID is `all`, then `RandomAssetOfPerson`; bucket ID is person ID (or `person@user` for multi-user).

- **PhotoPrism (current)**:
  - People-related methods are stubs: `RandomAssetOfPerson` and `RandomPersonFromAllPeople` return errors; `PersonAssetCount` returns 0; `AllNamedPeople` returns nil; `DisplayAsset.People` is nil.
  - Search API supports **person by name**: `q=person:"Jane Doe"` or `people:"Jane & John"` (exact subject names), and **by face ID**: `q=face:<face_uid>`.
  - Single subject: `GET /api/v1/subjects/:uid` returns one subject. A **list subjects** endpoint is not clearly documented; may exist as `GET /api/v1/subjects` (to be confirmed).
  - Photos response includes `Files[].Markers` (example showed `[]`); if Markers contain subject/face data, we can populate `DisplayAsset.People`.

## Constitution alignment

All work must comply with `.specify/memory/constitution.md`:

- **Code quality**: Lint-clean Go; reuse existing `source` types and provider interface; no new packages unless necessary.
- **Testing**: Unit tests for new PhotoPrism people logic (list/count/random), using httptest or mocks; existing Immich tests unchanged.
- **UX**: Same config keys and URL params for both sources where applicable; i18n for any new strings; docs updated (FEATURES-BY-SOURCE, config example).
- **Performance**: Reuse cache for subject list and person-scoped photo queries; bounded keys and TTLs.

## Phases / steps

### Phase 1: PhotoPrism API discovery and types

- **List subjects**: Determine whether PhotoPrism exposes a list/index for subjects (e.g. `GET /api/v1/subjects` with optional count/offset). If yes, define response type and use for `AllNamedPeople` and "random person from all." If no list endpoint exists, document that URL builder "People" dropdown stays empty for PhotoPrism and "all" people bucket is unsupported (or support only explicitly configured `people: ["Name"]`).
- **Subject/face in photo response**: Check if `GET /api/v1/photos` (or photo detail) returns subject/people data per photo (e.g. in `Files[].Markers` or a top-level `Subjects`). If present, add types and map to `source.Person` for `DisplayAsset.People` so "show person name" and more-info work.
- **Search behaviour**: Confirm that `q=person:"Name"` and `q=face:<uid>` work with `GET /api/v1/photos` (count, order=random) and that response headers (e.g. `X-Count`) give total count for `PersonAssetCount`.

### Phase 2: Implement people provider methods (PhotoPrism)

- **AllNamedPeople**: If list endpoint exists: GET list, filter to named (non-empty name), map to `[]source.Person` (ID = subject UID, Name = subject name, BirthDate if API provides it). If no list: return empty slice and document.
- **PersonAssetCount(personID)**:
  - If `personID == kiosk.PersonKeywordAll`: sum counts for all named people (from list) or return 0 if no list.
  - Else: treat `personID` as subject name or subject/face UID: build `q=person:"<name>"` for names, or `q=face:<uid>` for UID-like IDs; GET `/api/v1/photos?count=1&merged=true&primary=true&q=...`, read `X-Count` (or small count) for total. Use existing client and cache key pattern.
- **RandomAssetOfPerson(personID)**:
  - If `RequireAllPeople` and config has multiple people: build `people:"A & B"` (or equivalent) and fetch random photo with that query; otherwise single person filter.
  - Reuse `fetchPhotosWithCache` pattern: add person/face filter to `q`, same caching and batch behaviour. Set `Bucket` to `kiosk.SourcePerson` and `BucketID` to `personID` on the selected asset.
- **RandomPersonFromAllPeople**: If we have a subject list, filter to named, apply `ExcludedPeople`, pick random subject UID (or name) and return it. If no list, return error so "all" people bucket is skipped when source is PhotoPrism.

### Phase 3: DisplayAsset.People for PhotoPrism

- If photo response includes subject/marker data: extend `internal/photoprism` Photo/File types (e.g. `Markers` or `Subjects`) and map to `[]source.Person` in `DisplayAsset()`. If not available, leave `People` nil and document that "show person name" on asset may be empty for PhotoPrism until API provides it.

### Phase 4: Config and URL builder behaviour

- **Config**: No change to keys: `people`, `excluded_people`, `require_all_people` already exist. When `source == photoprism`, these apply to PhotoPrism subjects (by name or UID as decided in Phase 2). Validation: if source is PhotoPrism and people are configured, no extra validation beyond "subject list available or explicit names."
- **URL builder**: When source is PhotoPrism, call `AllNamedPeople`; if list is implemented, show people in dropdown and allow selecting/excluding; if not, people dropdown remains empty and doc says "configure people by name in config if needed."

### Phase 5: Documentation and tests

- **Docs**: Update `docs/FEATURES-BY-SOURCE.md`: set People row to supported for PhotoPrism with a short note (e.g. "by subject name; subject list if API supports it"). Update config example/schema comments for people when using PhotoPrism.
- **Tests**: Unit tests for: (1) PhotoPrism `AllNamedPeople` (with mock list response or empty); (2) `PersonAssetCount` with mock photos response and X-Count; (3) `RandomAssetOfPerson` with person filter in query; (4) `RandomPersonFromAllPeople` when list exists. Provider tests should not call real PhotoPrism.

## Success criteria

- [ ] PhotoPrism supports filtering by person (config `people`, URL param) using search filter `person:"Name"` (and optionally `face:<uid>` if we use subject UIDs).
- [ ] `PersonAssetCount` returns correct count for a given person/subject when PhotoPrism API is used (with cache).
- [ ] `RandomAssetOfPerson` returns a random photo containing that person and sets bucket to `SourcePerson` and bucket ID appropriately.
- [ ] If a subjects list API exists, `AllNamedPeople` and `RandomPersonFromAllPeople` work; URL builder shows people for PhotoPrism; "all" people bucket works. If not, these are documented and gracefully no-op or skip.
- [ ] `DisplayAsset.People` populated for PhotoPrism when API provides subject data per photo; otherwise nil and documented.
- [ ] No regressions for Immich; existing tests pass; new PhotoPrism people tests added.
- [ ] FEATURES-BY-SOURCE and config docs updated.

## Reference

- Spec: `.specify/specs/photoprism-source.md`
- Immich people: `internal/immich/immich_person.go`, `internal/immich/immich_faces.go`
- Provider interface: `internal/source/provider.go`
- PhotoPrism provider: `internal/photoprism/provider.go`
- Search filters: https://www.photoprism.app/kb/search-filters (person, people, face, subject)
- Constitution: `.specify/memory/constitution.md`
