# Plan: Replace Immich with PhotoPrism in the UI when PhotoPrism is the chosen source

## Goal

When `source: photoprism` is configured, all user-facing UI text and branding that currently says "Immich" should show "PhotoPrism" (or a neutral alternative). Where PhotoPrism does not offer equivalent functionality, the UI should either hide those controls, show them as disabled with a clear explanation, or propose documented alternatives.

## Context

- **Already source-aware**: The "View in …" link (more_info.templ) and error-page hints (error.templ) already branch on `viewData.Source` and use `photoprism_link_btn` / `error_message_photoprism` when source is PhotoPrism. No change needed there.
- **Remaining Immich-only UI**:
  - **Page title and meta**: `head.templ` has hardcoded `<title>Immich Kiosk</title>` and meta description "… displaying Immich assets …".
  - **Console log**: `views_home.templ` logs "Immich Kiosk version" in the browser console.
  - **About page**: `views_about.templ` always shows a "kiosk" and an "immich" section (with Immich server version/status). When source is PhotoPrism, the backend section should be PhotoPrism (version/status), and the section label/logo should be PhotoPrism, not Immich.
- **Immich Frame (TS)**: `frontend/src/ts/immichframe.ts` and `sleep.ts` reference "Immich Frame" in class names and log messages. This refers to the physical Immich Frame device (localhost:53287). Recommendation: keep the integration name in code; in any user-visible strings (if we add them later), use a generic term like "Frame" or "Kiosk" when source is PhotoPrism, or leave as-is if the UI does not expose "Immich Frame" text.
- **Features PhotoPrism does not offer** (see `docs/FEATURES-BY-SOURCE.md`): Memories, Star rating. For **Like / Hide / Tag** actions, PhotoPrism has equivalents and should follow the same config-driven behaviour as Immich: **Favorite** for Like (same `like_button_action`); **Hide** uses the same `hide_button_action` — **tag** adds/removes a label in PhotoPrism (e.g. kiosk-skip), **archive** toggles archive status in PhotoPrism; **Label** for the Tag button (add/remove label). The plan should propose wiring these to PhotoPrism APIs and using the correct UI terms.

## Constitution alignment

All work must comply with `.specify/memory/constitution.md`:

- **Code quality**: Lint-clean Go/templ/TS; no new packages unless necessary.
- **Testing**: Tests for any new logic (e.g. About page source branching); existing tests remain green.
- **UX**: i18n for any new or parameterized strings; consistent behaviour for both sources.
- **Performance**: No new network or heavy work beyond optional PhotoPrism version fetch on About.

## Phases / steps

### Phase 1: Page title, meta description, and console log

- **head.templ**:
  - Make `<title>` depend on `viewData.Source`: when `source == "photoprism"` use a title that includes "PhotoPrism" (e.g. "PhotoPrism Kiosk" or use i18n key `kiosk_title` with source parameter). When source is Immich, keep "Immich Kiosk".
  - Make meta description source-aware: e.g. "… displaying Immich assets …" vs "… displaying PhotoPrism assets …" (or a single generic "… displaying your photo library …" to avoid maintaining two strings). Prefer i18n keys (e.g. `meta_description_immich`, `meta_description_photoprism` or one generic key) so locales can translate.
- **views_home.templ** (kioskData / console.log):
  - Change the console message to be source-aware or generic: e.g. "Kiosk version" or "PhotoPrism Kiosk version" / "Immich Kiosk version" based on `viewData.Source` (if kioskData is extended to include source), so that when source is PhotoPrism the log does not say "Immich".

**Deliverables**: Title and meta reflect selected source; console log does not say "Immich" when source is PhotoPrism.

### Phase 2: About page — show PhotoPrism when source is PhotoPrism

- **Branch by source**:
  - When `viewData.Source == "photoprism"`: show a **PhotoPrism** section instead of the Immich section (version + status + link to latest release). When source is Immich, keep current behaviour (Immich section with `getImmichStats`).
- **PhotoPrism version/status**:
  - **If** PhotoPrism exposes a session or info API that returns version (e.g. from session or config endpoint — confirm via https://docs.photoprism.dev/): add a small helper (e.g. `getPhotoPrismStats(photoprismURL, token)`) that returns version string and online status, and use it for the PhotoPrism section. Reuse existing cache/timeout patterns where appropriate.
  - **If** no version API exists: show the PhotoPrism section with version as "—" or "N/A" and status derived from a simple health/connectivity check (e.g. existing client ping or a lightweight request). Document in FEATURES-BY-SOURCE or code comment.
- **Assets and labels**:
  - Use a **PhotoPrism logo** for the PhotoPrism section: add `photoprism-logo.svg` under the same assets path as `immich-logo.svg` (e.g. `frontend/public/assets/images/` or project `assets/` as used by the About view), or reference an existing asset if the project already has one. Ensure `getLatestRelease` supports `repo == "photoprism"` (e.g. owner `photoprism-app`, repo `photoprism`) for the "Latest version" link.
  - Make the Stats section logo `alt` text dynamic (e.g. "{Service} logo" or use i18n) so it does not hardcode "Immich Logo" when the section is PhotoPrism.

**Deliverables**: About page shows "PhotoPrism" (with version/status when available) when source is PhotoPrism; Immich section only when source is Immich. All new or touched strings i18n-ready.

### Phase 3: Features PhotoPrism does not offer — propose alternatives

For each capability that PhotoPrism does not support (see FEATURES-BY-SOURCE.md), define a UI behaviour:

| Feature | Current (Immich) | PhotoPrism | Proposed UI alternative |
|--------|-------------------|------------|--------------------------|
| **Memories** | Config + URL param; slideshow includes memory assets | Not supported | When `source == photoprism`: hide "Memories" / "Show memories" in URL builder and config UI, or show as disabled with tooltip: "Not available for PhotoPrism." Config validation already ignores memories when source is PhotoPrism; ensure UI does not imply it works. |
| **Star rating** | Config + display of rating on assets | Not supported | When `source == photoprism`: hide rating filter in URL builder; hide rating display in more-info/metadata (or show only when source is Immich). Config option can remain but have no effect; prefer hiding over showing "N/A" everywhere. |
| **Like / Hide / Tag actions** | Like (favorite/album), Hide (config: `hide_button_action` = tag and/or archive — tag adds/removes e.g. kiosk-skip, archive toggles archive), Tag (add tag) | Same: **Favorite**, **Hide** (tag → label, archive → archive), **Label** | When `source == photoprism`: (1) **Like** → **Favorite**; respect `like_button_action`. (2) **Hide** → respect `hide_button_action`: **tag** = add/remove label (e.g. kiosk-skip) via PhotoPrism label API; **archive** = toggle archive via PhotoPrism archive API. Same form and routes; implement AddTag/RemoveTag and ArchiveStatus in photoprism. (3) **Tag** → **Label**; add/remove label via API. i18n for Favorite, Hide/Archive, Label. Same config and routes as Immich. |
| **Archive** | Filter by archived status | No-op | When `source == photoprism`: hide "Show archived" in URL builder and config UI, or show disabled with tooltip "Not available for PhotoPrism." |

Implementation steps:

- **URL builder**: When `viewData.Source == "photoprism"`, hide or disable Memories, Rating, and Show archived controls (and any Immich-only options). Use existing `viewData.Source` in the URL builder view.
- **More-info / metadata**: If rating is shown, only show when source is Immich; when source is PhotoPrism, do not render rating block (or render disabled with tooltip). **Like / Hide / Tag buttons**: when source is PhotoPrism, use the same config as Immich (`like_button_action`, `hide_button_action`) — label buttons as **Favorite** (like), **Hide** (tag: add/remove label e.g. kiosk-skip; archive: toggle archive), and **Label** (add/remove label); implement backend support in `internal/photoprism`: `FavouriteStatus`, `AddTag`/`RemoveTag` (for label, including hide-with-tag), and `ArchiveStatus` (for hide-with-archive) so Hide behaves like Immich. Use source-dependent i18n for button labels.
- **Config page / docs**: Document that Memories and Rating are ignored or no-op for PhotoPrism. Document that Like/Hide/Tag use the same config (`like_button_action`, `hide_button_action`) for both sources: for PhotoPrism, Like → Favorite, Hide → tag (label) and/or archive per config, Tag → Label. No need to remove from config schema.

**Deliverables**: No Immich-only feature is implied to work for PhotoPrism without an equivalent. Memories and Rating are hidden or disabled when source is PhotoPrism. Like/Hide/Tag use the same config-driven behaviour for both sources: PhotoPrism provider implements `FavouriteStatus`, `AddTag`/`RemoveTag` (label, including hide-with-tag), and `ArchiveStatus` (hide-with-archive) so Hide matches Immich. FEATURES-BY-SOURCE.md updated accordingly.

### Phase 4: i18n and optional Immich Frame wording

- **Locales**: Add or use keys for source-dependent title and meta (e.g. `kiosk_title` with placeholder, or `kiosk_title_immich` / `kiosk_title_photoprism`). Ensure `source_name_photoprism` (and if needed `source_name_immich`) are used where a display name is needed. Update all locale files that currently have `immich_link_btn` / `error_message_immich` to keep parity for any new keys.
- **Immich Frame**: If the product name "Immich Frame" appears in any user-facing UI (e.g. settings or about), consider showing "Frame" or "Kiosk Frame" when source is PhotoPrism; otherwise leave as-is. No change required for internal class names or logs that are developer-only.

## Success criteria

- [ ] When `source: photoprism`, page title and meta description show PhotoPrism (or neutral) wording, not "Immich".
- [ ] When `source: photoprism`, browser console log does not say "Immich Kiosk".
- [ ] About page shows a PhotoPrism section (version/status when available, correct logo and latest-release link) when source is PhotoPrism; Immich section only when source is Immich.
- [ ] Memories and Star rating are either hidden or clearly marked as unavailable when source is PhotoPrism. Like/Hide/Tag use the same config as Immich (`like_button_action`, `hide_button_action`); PhotoPrism provider implements Favorite, Hide (tag → label, archive → archive API), and Label so behaviour matches Immich.
- [ ] New or changed user-facing strings are i18n-ready; existing locales updated as needed.
- [ ] No regressions for Immich; tests pass; lint clean.

## Reference

- **Config source**: `internal/config/config.go` (`Source`, `SourceImmich`, `SourcePhotoPrism`).
- **ViewData**: `internal/common/common.go` (embeds `config.Config`, so `viewData.Source` available in templates).
- **Already source-aware**: `internal/templates/partials/more_info.templ` (View in Immich/PhotoPrism), `internal/templates/partials/error.templ` (error_message_immich / error_message_photoprism).
- **About view**: `internal/templates/views/views_about.templ`; Immich stats: `internal/immich/immich_server.go` (`AboutInfo`).
- **PhotoPrism API**: https://docs.photoprism.dev/ (for session/info/version if available).
- **Features matrix**: `docs/FEATURES-BY-SOURCE.md`.
- **Provider mutations**: `internal/source/provider.go` (`AddTag`, `RemoveTag`, `AddToKioskLikedAlbum`, `RemoveFromKioskLikedAlbum`, `FavouriteStatus`, `ArchiveStatus`). Immich: Hide uses `hide_button_action` (tag and/or archive), same routes and form (e.g. tagName: kiosk.TagSkip). PhotoPrism stubs in `internal/photoprism/provider.go`: implement `AddTag`/`RemoveTag` (label API, for Hide-with-tag and Tag button), `ArchiveStatus` (archive API, for Hide-with-archive), and `AddToKioskLikedAlbum`/`RemoveFromKioskLikedAlbum` or `FavouriteStatus` (Favorite) so behaviour matches Immich.
- **Constitution**: `.specify/memory/constitution.md`.
