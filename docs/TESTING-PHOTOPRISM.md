# Manual testing with PhotoPrism

Use these steps to run the kiosk with PhotoPrism as the media source.

## 1. Get a PhotoPrism token

In PhotoPrism: **Settings → Account → Apps & Devices** → create an app password, or use **Settings → Account → API Clients** to create a token. Copy the token (you’ll use it as `photoprism_token`).

## 2. Configure the kiosk

Create or edit `config.yaml` in the project root (or set env vars).

**Minimal config for PhotoPrism:**

```yaml
source: photoprism
photoprism_url: "http://localhost:2342"   # or https://your-photoprism.example.com
photoprism_token: "YOUR_APP_PASSWORD_OR_TOKEN"

# Optional: leave empty to show all photos, or restrict to albums/labels
# albums:
#   - "album-uid-from-photoprism"
# tags: []
```

If you use **environment variables** instead:

```bash
export KIOSK_SOURCE=photoprism
export KIOSK_PHOTOPRISM_URL="http://localhost:2342"
export KIOSK_PHOTOPRISM_TOKEN="your-token"
```

No `immich_url` or `immich_api_key` are required when `source` is `photoprism`.

## 3. Run the kiosk

From the repo root:

```bash
go run .
```

Or build and run:

```bash
go build -o immich-kiosk .
./immich-kiosk
```

Default port is **3000**. Open:

- **http://localhost:3000** (or with password if you set `kiosk.password` in config)

## 4. What to try

- **Home** – Should load a random photo from PhotoPrism (or from albums/tags if you set them).
- **Next/Previous** – Navigate history.
- **URL Builder** – If enabled (`kiosk.enable_url_builder: true`), use it to pick albums/labels and build URLs; album IDs are PhotoPrism album UIDs.

## 5. Debugging

- **UI has no styles or looks broken:**  
  Build the frontend first so CSS/JS are embedded:  
  `cd frontend && bun run css && bun run js` (or `npm run css && npm run js`).  
  Then run the kiosk from the repo root again.
- **Log level:**  
  `KIOSK_LOG_LEVEL=debug` or `verbose` for more logs.
- **Config check:**  
  With a config file, the app validates `source` and required URL/token; it will exit with a clear error if PhotoPrism is selected but URL or token is missing.
- **Network:**  
  Ensure the kiosk can reach `photoprism_url` (same machine: `http://localhost:2342`; Docker: use the container hostname or host network as needed).

## 6. Running integration tests against live PhotoPrism

The `internal/photoprism` package includes optional tests that call a real PhotoPrism API to verify Bearer token authentication. They are skipped unless the required environment variables are set, so CI does not need a PhotoPrism instance.

1. **Set credentials** (use the same app password or token as for the kiosk, e.g. from **Settings → Account → Apps and Devices**):

   ```bash
   export KIOSK_PHOTOPRISM_URL="http://localhost:2342"
   export KIOSK_PHOTOPRISM_TOKEN="your-app-password-or-token"
   ```

2. **Run the live tests:**

   ```bash
   go test -v ./internal/photoprism/... -run Live
   ```

   - `TestClient_Live_ValidToken` – GET `/api/v1/photos?count=1` with your token; expects 200 and valid JSON.
   - `TestClient_Live_InvalidToken` – same endpoint with an invalid token; expects an error (e.g. 401).

Without `KIOSK_PHOTOPRISM_URL` and `KIOSK_PHOTOPRISM_TOKEN`, these tests are skipped.

## 7. Features per source

See [FEATURES-BY-SOURCE.md](FEATURES-BY-SOURCE.md) for a table of what works with Immich vs PhotoPrism (albums, tags, people, memories, rating, etc.).

## 8. Known limitations (PhotoPrism source)

- **People** – Not supported (filter/weighting skipped).
- **Memories** – Not supported.
- **Star rating** – Not supported.
- **Like/Hide/Tag** – Buttons work in the UI but do not change state in PhotoPrism (no-ops).
- **Videos** – Supported if PhotoPrism returns them from the photos API; playback uses PhotoPrism’s video endpoint.

If you see “no photos” or empty slides, check that PhotoPrism has indexed photos and that the token has access. Use `KIOSK_LOG_LEVEL=debug` to see API requests and any errors.
