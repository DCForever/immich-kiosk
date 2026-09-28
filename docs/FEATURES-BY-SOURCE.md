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
