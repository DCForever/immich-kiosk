# Features by source

The kiosk supports two media backends: **Immich** and **PhotoPrism**. Most behaviour (slideshow, transitions, config, URL builder) is the same; some options apply only to one source.

| Feature | Immich | PhotoPrism |
|--------|--------|------------|
| **Photos** | ✅ | ✅ |
| **Videos** | ✅ | ✅ (if returned by API) |
| **Albums** | ✅ (album IDs) | ✅ (album UIDs) |
| **Tags / labels** | ✅ | ✅ (labels) |
| **Date range** | ✅ | ✅ |
| **Favourites** | ✅ | ✅ (favourite filter) |
| **People** | ✅ | ✅ (by subject name; subject list if API supports it) |
| **Memories** | ✅ | ✅ (on this day, with week/month fallback and collage) |
| **Star rating** | ✅ | ❌ Not supported |
| **Like / Hide / Tag actions** | ✅ (writes to Immich) | No-op (does not write to PhotoPrism) |
| **Archive** | ✅ | No-op (does not write to PhotoPrism) |
| **URL Builder** | ✅ | ✅ (albums, labels) |
| **Caching** | ✅ | ✅ |
| **Prefetch** | ✅ | ✅ |
| **Offline mode** | ✅ | ✅ (same flow) |

When `source: photoprism`, star rating is not supported (the rating endpoint returns not implemented). Like, hide, tag, and archive actions are intentional no-ops and do not change PhotoPrism. People filtering uses PhotoPrism subject names (config `people`); if GET /api/v1/subjects is available, the URL builder shows people in the dropdown. Memories use PhotoPrism search filters (exact date, then week, then month).

The PhotoPrism client reads `GET /api/v1/config` once per provider and logs the `version` field when debug logging is on. A failed or missing version does not block the slideshow. List and detail JSON is decoded defensively (key case, empty dates, numeric booleans and strings, wrapped `{photos:[]}` collections, and preview-token header aliases). Thumbnail requests try `fit_720`, then `fit_1280`, then `fit_1920`. Fork-specific endpoints are not assumed until they are confirmed.
