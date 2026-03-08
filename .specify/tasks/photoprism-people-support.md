# Tasks: PhotoPrism people support

## Reference

- Plan: `.specify/plans/photoprism-people-support.md`
- Spec: `.specify/specs/photoprism-source.md`
- Constitution: `.specify/memory/constitution.md`

## Task list

Tasks must respect:

- **Code quality**: golangci-lint, project structure; reuse `source` types and provider interface; no new packages unless necessary.
- **Testing**: Add/update tests for new PhotoPrism people logic; mock API; keep CI passing.
- **UX**: Same config keys and URL params for both sources; i18n for new strings; update FEATURES-BY-SOURCE and config docs.
- **Performance**: Reuse cache for subject list and person-scoped queries; bounded keys and TTLs.

| # | Task | Done |
|---|------|------|
| **Phase 1: PhotoPrism API discovery and types** |
| 1 | Discover whether PhotoPrism exposes a list/index for subjects (e.g. `GET /api/v1/subjects` with count/offset). Try the endpoint; if it exists, document response shape and add a minimal Go type for subject list items (UID, Name, optional BirthDate). If no list exists, document that URL builder People dropdown and “all” people bucket will be unsupported for PhotoPrism. | |
| 2 | Check if `GET /api/v1/photos` (or photo detail) returns subject/people data per photo (e.g. `Files[].Markers` or top-level `Subjects`). If present, document the field names and add types in `internal/photoprism/types.go` for mapping to `DisplayAsset.People`. If absent, document that per-photo people will remain nil for PhotoPrism. | |
| 3 | Confirm search behaviour: call `GET /api/v1/photos?count=1&merged=true&primary=true&q=person:"<name>"` (and optionally `q=face:<uid>`) against a real or demo instance; verify response and that `X-Count` (or equivalent) returns total matches for use in `PersonAssetCount`. Document result. | |
| **Phase 2: Implement people provider methods (PhotoPrism)** |
| 4 | Implement `AllNamedPeople`: if list endpoint exists (from task 1), add client method to fetch subjects, filter to non-empty name, map to `[]source.Person` (ID = subject UID, Name = subject name, BirthDate if API provides). If no list, keep returning `nil, nil` and add a short code comment. Use cache with bounded key/TTL for the list. | |
| 5 | Implement `PersonAssetCount(personID)`: if `personID == kiosk.PersonKeywordAll`, sum counts for all named people from list (or return 0 if no list). Otherwise build `q=person:"<name>"` for name or `q=face:<uid>` for UID-like IDs; GET photos with count=1 and read `X-Count` for total. Reuse client and cache key pattern. Return (count, nil) or (0, err). | |
| 6 | Implement `RandomAssetOfPerson(personID)`: add person/face filter to `q` in photo fetch (single `person:"..."` or `people:"A & B"` when `RequireAllPeople` and multiple config people). Reuse `fetchPhotosWithCache` pattern (extend to accept person filter); set `Bucket` to `kiosk.SourcePerson` and `BucketID` to `personID` in `DisplayAsset` path. Remove the current “not supported” error return. | |
| 7 | Implement `RandomPersonFromAllPeople`: if subject list exists, filter to named, apply `ExcludedPeople`, pick random subject UID (or name) and return it. If no list, return a clear error so “all” people bucket is skipped when source is PhotoPrism. | |
| **Phase 3: DisplayAsset.People for PhotoPrism** |
| 8 | If photo response includes subject/marker data (from task 2): extend Photo/File types in `internal/photoprism/types.go`, and in `DisplayAsset()` map to `[]source.Person` and set `People` on the returned `source.DisplayAsset`. If data not available, leave `People` nil and add a comment that “show person name” may be empty for PhotoPrism until API provides it. | |
| **Phase 4: Config and URL builder** |
| 9 | Ensure config keys `people`, `excluded_people`, `require_all_people` are used when `source == photoprism` (no schema change). If any validation specifically blocked people for PhotoPrism, remove or relax it so configured names/UIDs are accepted. | |
| 10 | URL builder: when source is PhotoPrism, `AllNamedPeople` is already called; if list is implemented, people dropdown will show subjects. If not, ensure UI handles empty people list gracefully and add a brief doc note: “Configure people by name in config when using PhotoPrism if subject list is not available.” | |
| **Phase 5: Documentation and tests** |
| 11 | Update `docs/FEATURES-BY-SOURCE.md`: set People row to supported for PhotoPrism (e.g. “✅ (by subject name; subject list if API supports it)”). Update config example and/or schema comments for people when using PhotoPrism. | |
| 12 | Add unit tests for PhotoPrism people: (1) `AllNamedPeople` with mock list response or empty; (2) `PersonAssetCount` with mock photos response and X-Count; (3) `RandomAssetOfPerson` with person filter in query; (4) `RandomPersonFromAllPeople` when list exists (and error when not). Use httptest or in-memory mocks; no real PhotoPrism calls. | |
| 13 | Run full test suite and golangci-lint; fix any regressions. Confirm Immich people flows and existing PhotoPrism tests still pass. | |

## Definition of done (per task)

- [ ] Implements only what the task describes.
- [ ] Aligns with constitution (quality, tests, UX, performance).
- [ ] Lint and tests pass for touched code.
- [ ] No new third-party libraries; reuses existing cache and client patterns.
