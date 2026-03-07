# Immich Kiosk — Project Constitution

Governing principles for the Immich Kiosk project. All specifications, plans, and implementation must adhere to these principles. They take precedence over other instructions.

---

## 1. Code Quality

- **Go standards**: Follow effective Go style and idioms. Use `golangci-lint` with the project’s `.golangci.yml`; all new and modified code must pass with no new issues. Prefer small, focused packages and clear boundaries between internal domains (routes, immich client, config, cache, templates).
- **Naming and structure**: Use descriptive names; avoid abbreviations except where idiomatic (e.g. `id`, `url`). Keep functions short and single-purpose; avoid deep nesting. Prefer returning errors and handling them at boundaries; use `log/slog` for logging (no deprecated `log`).
- **Dependencies**: Prefer the standard library and existing project dependencies. New dependencies require justification and must be maintained and compatible with the current Go module. No deprecated packages (e.g. `github.com/golang/protobuf`, `math/rand` without `/v2`).
- **Frontend (TypeScript/CSS)**: Lint with Biome; TypeScript must pass `tsc --noEmit`. Keep bundles minimal; use the existing build pipeline (esbuild, PostCSS). Prefer declarative markup and htmx where it fits; keep custom JS focused and type-safe.
- **Templ**: Keep templates readable and maintainable; avoid heavy logic in templates. Reuse components and partials for consistency.

---

## 2. Testing Standards

- **Backend (Go)**: Write tests for exported behavior and non-trivial logic (e.g. Immich client, config validation, caching, date/URL helpers). Prefer table-driven tests with clear cases and assertions. Use `testify` where it improves clarity. Mock external services (Immich API) in tests; avoid real network calls in CI.
- **Coverage**: New packages or major behavior changes should include tests. Critical paths (auth, config loading, asset fetching, URL building) must have tests. No requirement for a specific coverage percentage, but regressions in covered behavior must be caught by tests.
- **Frontend**: TypeScript is validated by `tsc --noEmit` as part of the build. Add unit or integration tests for non-trivial TS logic (e.g. URL builder, timing) when it reduces regression risk.
- **CI**: All tests must pass in CI (e.g. `go test ./...`, frontend build). Do not disable or skip tests without a documented reason and follow-up to fix or replace them.

---

## 3. User Experience Consistency

- **Kiosk-first**: The app is a configurable slideshow for Immich assets. UX must support unattended, long-running display: clear transitions, predictable timing, and no modal or blocking flows that require interaction during normal slideshow operation. Settings and setup can be more interactive.
- **Configuration-driven behavior**: Respect user configuration (display duration, order, filters, themes). Avoid hardcoding values that users can reasonably expect to configure. Document config options where they affect UX.
- **Internationalization (i18n)**: Use the project’s i18n system for all user-facing strings (UI labels, messages, errors). No hardcoded English-only copy in templates or frontend. Support RTL where applicable and use locale-aware formatting for dates and numbers.
- **Accessibility**: Ensure sufficient color contrast and focus visibility. Prefer semantic HTML and ARIA where it improves screen reader support. Don’t rely on color alone to convey information. Keep keyboard and pointer usage consistent.
- **Visual consistency**: Use the same design tokens (spacing, typography, colors) across views. Reuse existing templ components and CSS patterns. Theming (e.g. light/dark) must be consistent and respect user or system preference where configured.

---

## 4. Performance Requirements

- **Slideshow smoothness**: Transitions between assets must feel smooth (no unnecessary layout thrash or long freezes). Preload the next asset when possible so that switching is quick. Avoid blocking the main thread for long periods during the slideshow loop.
- **Memory and caching**: Use the existing cache layer for Immich responses and derived data where appropriate. Avoid unbounded growth (e.g. caches should have limits or TTLs). Be mindful of image size and count on low-memory or many-tab scenarios.
- **Asset loading**: Prefer lazy or progressive loading where it improves perceived performance. Prefer efficient image formats and sizes as configured (e.g. existing thumb/encoding usage). Minimize duplicate requests and respect cache headers when calling the Immich API.
- **Backend**: Handlers should respond quickly; heavy work (e.g. image processing, large list building) should be cached, backgrounded, or bounded. Avoid N+1 patterns when fetching related data from the API.
- **Frontend assets**: Keep JS and CSS bundle size in check; use the existing minification and tree-shaking. Avoid large dependencies for small features.

---

## Amendment and Compliance

- **Updates**: Changes to this constitution should be proposed and reviewed like other project changes. Principles can be amended to reflect new constraints or priorities.
- **Compliance**: Specs, task breakdowns, and code reviews should explicitly check alignment with these principles. When in doubt, favor the constitution over ad hoc preferences.
